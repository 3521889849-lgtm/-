// Package ticketapp 承载“票务域”的应用层用例（车次详情/余票统计等）。
package ticketapp

import (
	"context"
	"errors"
	"example_shop/common/db"
	"example_shop/internal/ticket_service/app/shared"
	"example_shop/internal/ticket_service/model"
	kitexuser "example_shop/kitex_gen/user"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Service struct{}

// New 创建票务域应用服务。
func New() *Service {
	return &Service{}
}

// GetTrainDetail 车次详情用例：返回车次基础信息、各席别余票与最低价。
func (s *Service) GetTrainDetail(ctx context.Context, req *kitexuser.GetTrainDetailReq) (*kitexuser.GetTrainDetailResp, error) {
	if strings.TrimSpace(req.TrainId) == "" {
		return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 400, Msg: "train_id不能为空"}}, nil
	}

	var t model.TrainInfo
	if err := db.ReadDB().Where("train_id = ?", req.TrainId).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 404, Msg: "车次不存在"}}, nil
		}
		return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "查询失败: " + err.Error()}}, nil
	}

	fromStation := t.DepartureStation
	toStation := t.ArrivalStation
	if req.DepartureStation != nil && strings.TrimSpace(*req.DepartureStation) != "" {
		fromStation = strings.TrimSpace(*req.DepartureStation)
	}
	if req.ArrivalStation != nil && strings.TrimSpace(*req.ArrivalStation) != "" {
		toStation = strings.TrimSpace(*req.ArrivalStation)
	}

	var depStop model.TrainStationPass
	if err := db.ReadDB().Where("train_id = ? AND station_name = ?", t.ID, fromStation).First(&depStop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 404, Msg: "上车站不在该车次途经站列表"}}, nil
		}
		return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "查询上车站失败: " + err.Error()}}, nil
	}
	var arrStop model.TrainStationPass
	if err := db.ReadDB().Where("train_id = ? AND station_name = ?", t.ID, toStation).First(&arrStop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 404, Msg: "下车站不在该车次途经站列表"}}, nil
		}
		return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "查询下车站失败: " + err.Error()}}, nil
	}
	if depStop.Sequence >= arrStop.Sequence {
		return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 400, Msg: "上车站必须早于下车站"}}, nil
	}

	depTime := t.DepartureTime
	if depStop.DepartureTime.Valid {
		depTime = depStop.DepartureTime.Time
	} else if depStop.ArrivalTime.Valid {
		depTime = depStop.ArrivalTime.Time
	}
	arrTime := t.ArrivalTime
	if arrStop.ArrivalTime.Valid {
		arrTime = arrStop.ArrivalTime.Time
	} else if arrStop.DepartureTime.Valid {
		arrTime = arrStop.DepartureTime.Time
	}

	var seatTypesRaw []string
	if err := db.ReadDB().Table("seat_infos").Select("DISTINCT seat_type").Where("train_id = ? AND deleted_at IS NULL", t.ID).Scan(&seatTypesRaw).Error; err != nil {
		return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "查询座席类型失败: " + err.Error()}}, nil
	}

	seatTypes := make([]*kitexuser.SeatTypeRemain, 0, len(seatTypesRaw))
	for _, st := range seatTypesRaw {
		st = strings.TrimSpace(st)
		if st == "" {
			continue
		}
		remain, err := countRemainingSegment(t.ID, st, depStop.Sequence, arrStop.Sequence)
		if err != nil {
			return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "座席统计失败: " + err.Error()}}, nil
		}
		minPrice, err := minSeatPrice(t.ID, st)
		if err != nil {
			return &kitexuser.GetTrainDetailResp{BaseResp: &kitexuser.BaseResp{Code: 500, Msg: "票价统计失败: " + err.Error()}}, nil
		}
		seatTypes = append(seatTypes, &kitexuser.SeatTypeRemain{SeatType: st, Remaining: int32(remain), MinPrice: shared.Money2(minPrice)})
	}

	trainCode := ""
	if strings.TrimSpace(t.TrainCode) != "" {
		trainCode = t.TrainCode
	}
	runtime := int32(arrTime.Sub(depTime).Minutes())
	return &kitexuser.GetTrainDetailResp{
		BaseResp:          &kitexuser.BaseResp{Code: 200, Msg: "success"},
		TrainId:           t.ID,
		TrainCode:         trainCode,
		TrainType:         t.TrainType,
		DepartureStation:  fromStation,
		ArrivalStation:    toStation,
		DepartureTimeUnix: depTime.Unix(),
		ArrivalTimeUnix:   arrTime.Unix(),
		RuntimeMinutes:    runtime,
		SeatTypes:         seatTypes,
	}, nil
}

func countRemainingSegment(trainID, seatType string, fromSeq, toSeq uint32) (int64, error) {
	if fromSeq == 0 || toSeq == 0 || fromSeq >= toSeq {
		return 0, nil
	}
	now := time.Now()
	var count int64
	err := db.ReadDB().Table("seat_infos si").
		Where("si.train_id = ? AND si.seat_type = ? AND si.status = ? AND si.deleted_at IS NULL", trainID, seatType, "AVAILABLE").
		Where(
			"NOT EXISTS (SELECT 1 FROM seat_segment_occupancies o WHERE o.deleted_at IS NULL AND o.train_id = si.train_id AND o.seat_id = si.seat_id AND o.status IN ('SOLD','LOCKED') AND (o.status <> 'LOCKED' OR o.lock_expire_time > ?) AND o.from_seq < ? AND o.to_seq > ?)",
			now, toSeq, fromSeq,
		).
		Count(&count).Error
	return count, err
}

func minSeatPrice(trainID, seatType string) (float64, error) {
	var price float64
	err := db.ReadDB().Model(&model.SeatInfo{}).
		Select("MIN(seat_price)").
		Where("train_id = ? AND seat_type = ?", trainID, seatType).
		Scan(&price).Error
	return price, err
}
