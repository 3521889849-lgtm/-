package ticket

import (
	"context"
	"example_shop/common/cache"
	"example_shop/common/db"
	"example_shop/internal/gateway/http/dto"
	model2 "example_shop/internal/model"
	"sort"
	"strings"
	"time"
)

var stationListCache = cache.NewTTLCache[[]string]()

func normalizeStationKeyword(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "")
	return strings.ToLower(s)
}

func stationAlias(keyword string) string {
	switch normalizeStationKeyword(keyword) {
	case "沪", "sh", "shanghai", "shangh":
		return "上海"
	case "京", "bj", "beijing":
		return "北京"
	case "广", "gz", "guangzhou":
		return "广州"
	case "深", "sz", "shenzhen":
		return "深圳"
	default:
		return ""
	}
}

func loadStations(ctx context.Context) ([]string, error) {
	if v, ok := stationListCache.Get("all"); ok && len(v) > 0 {
		return v, nil
	}

	type row struct {
		Name string
	}
	stations := make(map[string]struct{})

	var all []row
	if err := db.ReadDB().Model(&model2.TrainStationPass{}).Select("DISTINCT station_name AS name").Scan(&all).Error; err != nil {
		return nil, err
	}
	for _, r := range all {
		name := strings.TrimSpace(r.Name)
		if name != "" {
			stations[name] = struct{}{}
		}
	}

	items := make([]string, 0, len(stations))
	for s := range stations {
		items = append(items, s)
	}
	sort.Strings(items)

	stationListCache.Set("all", items, 10*time.Minute)
	return items, nil
}

// StationSuggest 站点联想用例：从“途经站点表”提取候选站点。
//
// 使用 DISTINCT station_name 汇总全站点集合并做 TTL 缓存，避免每次输入联想都打 DB。
func (s *Service) StationSuggest(ctx context.Context, req dto.StationSuggestHTTPReq) Result {
	limit := req.Limit
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	kw := normalizeStationKeyword(req.Keyword)
	if len(kw) > 64 {
		return Result{Status: 400, Body: dto.BaseHTTPResp{Code: 400, Msg: "keyword过长"}}
	}

	stations, err := loadStations(ctx)
	if err != nil {
		return Result{Status: 500, Body: dto.BaseHTTPResp{Code: 500, Msg: "加载站点失败: " + err.Error()}}
	}

	alias := stationAlias(req.Keyword)
	res := make([]string, 0, limit)
	seen := make(map[string]struct{}, limit)

	if alias != "" {
		seen[alias] = struct{}{}
		res = append(res, alias)
	}

	for _, st := range stations {
		if len(res) >= limit {
			break
		}
		if _, ok := seen[st]; ok {
			continue
		}
		if strings.Contains(normalizeStationKeyword(st), kw) || strings.Contains(st, req.Keyword) {
			seen[st] = struct{}{}
			res = append(res, st)
		}
	}

	return Result{Status: 200, Body: dto.StationSuggestHTTPResp{Code: 200, Msg: "success", Items: res}}
}

