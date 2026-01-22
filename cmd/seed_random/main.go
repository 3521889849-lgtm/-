package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"example_shop/common/config"
	"example_shop/common/db"
	"example_shop/internal/ticket_service/model"

	"gorm.io/gorm/clause"
)

type seedOptions struct {
	StartOffsetDays int
	Days            int
	TrainsPerDay    int
	Routes          int
	MinDepartHour   int
	MaxDepartHour   int
	Seed            int64
	BatchSize       int
	SeatDensity     float64
	FlushTrains     int
	FlushSeats      int
	FlushStops      int
}

type route struct {
	Dep string
	Arr string
	Km  float64
}

type station struct {
	Name string
	Lat  float64
	Lon  float64
}

type plan struct {
	TrainCode        string
	TrainType        string
	DepartureStation string
	ArrivalStation   string
	DistanceKm       float64
	DepartMinute     int
	RuntimeMinutes   uint32
}

func main() {
	op := parseFlags()

	if err := config.ViperInit(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	if err := db.MysqlInit(); err != nil {
		log.Fatalf("mysql初始化失败: %v", err)
	}

	rng := rand.New(rand.NewSource(op.Seed))

	stations := stationPool()
	routes := buildRoutes(rng, stations, op.Routes)
	plans := buildPlans(rng, routes, op.TrainsPerDay, op.MinDepartHour, op.MaxDepartHour)

	createdTrains := 0
	createdSeats := 0
	startDate := dateOnly(time.Now().AddDate(0, 0, op.StartOffsetDays))

	trainBuf := make([]model.TrainInfo, 0, op.FlushTrains)
	seatBuf := make([]model.SeatInfo, 0, op.FlushSeats)
	stopBuf := make([]model.TrainStationPass, 0, op.FlushStops)
	flush := func() {
		if len(trainBuf) > 0 {
			if err := db.MysqlDB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(trainBuf, op.BatchSize).Error; err != nil {
				log.Fatalf("插入车次失败: %v", err)
			}
			trainBuf = trainBuf[:0]
		}
		if len(seatBuf) > 0 {
			if err := db.MysqlDB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(seatBuf, op.BatchSize).Error; err != nil {
				log.Fatalf("插入座位失败: %v", err)
			}
			seatBuf = seatBuf[:0]
		}
		if len(stopBuf) > 0 {
			if err := db.MysqlDB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(stopBuf, op.BatchSize).Error; err != nil {
				log.Fatalf("插入途经站失败: %v", err)
			}
			stopBuf = stopBuf[:0]
		}
	}

	for d := 0; d < op.Days; d++ {
		serviceDate := startDate.AddDate(0, 0, d)
		for _, p := range plans {
			departureTime := serviceDate.Add(time.Duration(p.DepartMinute) * time.Minute)
			arrivalTime := departureTime.Add(time.Duration(p.RuntimeMinutes) * time.Minute)
			trainID := makeTrainID(p.TrainCode, serviceDate, departureTime)
			seatLayout, seats := buildSeatsForTrain(rng, trainID, p.TrainType, p.DistanceKm, op.SeatDensity)
			stops := buildStationPasses(rng, trainID, departureTime, arrivalTime, p.DepartureStation, p.ArrivalStation, stationNames(stations))

			trainBuf = append(trainBuf, model.TrainInfo{
				ID:               trainID,
				TrainCode:        p.TrainCode,
				ServiceDate:      model.NullDate(serviceDate),
				TrainType:        p.TrainType,
				DepartureStation: p.DepartureStation,
				ArrivalStation:   p.ArrivalStation,
				DepartureTime:    departureTime,
				ArrivalTime:      arrivalTime,
				RuntimeMinutes:   p.RuntimeMinutes,
				SeatLayout:       &seatLayout,
				Status:           "NORMAL",
			})
			createdTrains++

			seatBuf = append(seatBuf, seats...)
			createdSeats += len(seats)
			stopBuf = append(stopBuf, stops...)

			if len(trainBuf) >= op.FlushTrains || len(seatBuf) >= op.FlushSeats || len(stopBuf) >= op.FlushStops {
				flush()
			}
		}
	}
	flush()

	fmt.Printf("seed_random 完成：days=%d trains/day=%d routes=%d\n", op.Days, op.TrainsPerDay, op.Routes)
	fmt.Printf("写入（尝试）：trains=%d seats=%d\n", createdTrains, createdSeats)
	fmt.Printf("可查询日期范围：%s ~ %s\n", startDate.Format("2006-01-02"), startDate.AddDate(0, 0, op.Days-1).Format("2006-01-02"))
}

func parseFlags() seedOptions {
	op := seedOptions{}
	flag.IntVar(&op.StartOffsetDays, "start-offset-days", 1, "从今天起偏移多少天开始写入（默认明天）")
	flag.IntVar(&op.Days, "days", 7, "写入多少天的车次")
	flag.IntVar(&op.TrainsPerDay, "trains-per-day", 200, "每天写入多少条车次（同一批 train_code 会在每天重复出现，时刻固定）")
	flag.IntVar(&op.Routes, "routes", 60, "线路数（出发站/到达站组合）")
	flag.IntVar(&op.MinDepartHour, "min-depart-hour", 6, "最早发车小时")
	flag.IntVar(&op.MaxDepartHour, "max-depart-hour", 22, "最晚发车小时")
	flag.Int64Var(&op.Seed, "seed", time.Now().UnixNano(), "随机种子")
	flag.IntVar(&op.BatchSize, "batch-size", 500, "批量写入大小")
	flag.Float64Var(&op.SeatDensity, "seat-density", 0.4, "座位密度缩放（0.1~1.0），用于控制每趟车生成的座位数量")
	flag.IntVar(&op.FlushTrains, "flush-trains", 500, "车次缓冲区达到多少条后批量写入")
	flag.IntVar(&op.FlushSeats, "flush-seats", 20000, "座位缓冲区达到多少条后批量写入")
	flag.IntVar(&op.FlushStops, "flush-stops", 5000, "途经站缓冲区达到多少条后批量写入")
	flag.Parse()

	if op.Days <= 0 {
		op.Days = 1
	}
	if op.TrainsPerDay <= 0 {
		op.TrainsPerDay = 1
	}
	if op.Routes <= 0 {
		op.Routes = 10
	}
	if op.MinDepartHour < 0 {
		op.MinDepartHour = 0
	}
	if op.MaxDepartHour > 23 {
		op.MaxDepartHour = 23
	}
	if op.MaxDepartHour < op.MinDepartHour {
		op.MaxDepartHour = op.MinDepartHour
	}
	if op.BatchSize <= 0 {
		op.BatchSize = 500
	}
	if op.SeatDensity <= 0 {
		op.SeatDensity = 0.4
	}
	if op.SeatDensity > 1 {
		op.SeatDensity = 1
	}
	if op.FlushTrains <= 0 {
		op.FlushTrains = 500
	}
	if op.FlushSeats <= 0 {
		op.FlushSeats = 20000
	}
	if op.FlushStops <= 0 {
		op.FlushStops = 5000
	}
	return op
}

func stationPool() []station {
	return []station{
		{Name: "北京", Lat: 39.9042, Lon: 116.4074},
		{Name: "上海", Lat: 31.2304, Lon: 121.4737},
		{Name: "广州", Lat: 23.1291, Lon: 113.2644},
		{Name: "深圳", Lat: 22.5431, Lon: 114.0579},
		{Name: "杭州", Lat: 30.2741, Lon: 120.1551},
		{Name: "南京", Lat: 32.0603, Lon: 118.7969},
		{Name: "苏州", Lat: 31.2989, Lon: 120.5853},
		{Name: "成都", Lat: 30.5728, Lon: 104.0668},
		{Name: "重庆", Lat: 29.5630, Lon: 106.5516},
		{Name: "武汉", Lat: 30.5928, Lon: 114.3055},
		{Name: "西安", Lat: 34.3416, Lon: 108.9398},
		{Name: "郑州", Lat: 34.7466, Lon: 113.6254},
		{Name: "长沙", Lat: 28.2278, Lon: 112.9389},
		{Name: "青岛", Lat: 36.0671, Lon: 120.3826},
		{Name: "厦门", Lat: 24.4798, Lon: 118.0894},
		{Name: "福州", Lat: 26.0745, Lon: 119.2965},
		{Name: "天津", Lat: 39.0842, Lon: 117.2000},
		{Name: "沈阳", Lat: 41.8057, Lon: 123.4315},
		{Name: "大连", Lat: 38.9140, Lon: 121.6147},
		{Name: "哈尔滨", Lat: 45.8038, Lon: 126.5349},
		{Name: "济南", Lat: 36.6512, Lon: 117.1201},
		{Name: "合肥", Lat: 31.8206, Lon: 117.2272},
		{Name: "昆明", Lat: 25.0389, Lon: 102.7183},
		{Name: "南宁", Lat: 22.8170, Lon: 108.3669},
		{Name: "贵阳", Lat: 26.6470, Lon: 106.6302},
		{Name: "兰州", Lat: 36.0611, Lon: 103.8343},
		{Name: "乌鲁木齐", Lat: 43.8256, Lon: 87.6168},
		{Name: "石家庄", Lat: 38.0428, Lon: 114.5149},
		{Name: "太原", Lat: 37.8706, Lon: 112.5489},
		{Name: "南昌", Lat: 28.6820, Lon: 115.8579},
		{Name: "宁波", Lat: 29.8683, Lon: 121.5440},
		{Name: "无锡", Lat: 31.4912, Lon: 120.3119},
		{Name: "常州", Lat: 31.8107, Lon: 119.9741},
		{Name: "佛山", Lat: 23.0215, Lon: 113.1214},
		{Name: "东莞", Lat: 23.0207, Lon: 113.7518},
		{Name: "珠海", Lat: 22.2707, Lon: 113.5767},
		{Name: "中山", Lat: 22.5159, Lon: 113.3928},
		{Name: "温州", Lat: 27.9949, Lon: 120.6994},
		{Name: "泉州", Lat: 24.8739, Lon: 118.6759},
		{Name: "南通", Lat: 31.9802, Lon: 120.8943},
	}
}

func buildRoutes(rng *rand.Rand, stations []station, n int) []route {
	if n < 1 {
		n = 1
	}

	seen := make(map[string]struct{}, n)
	items := make([]route, 0, n)
	for len(items) < n {
		dep := stations[rng.Intn(len(stations))]
		arr := stations[rng.Intn(len(stations))]
		if dep.Name == arr.Name {
			continue
		}
		k := dep.Name + "->" + arr.Name
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		items = append(items, route{Dep: dep.Name, Arr: arr.Name, Km: haversineKm(dep.Lat, dep.Lon, arr.Lat, arr.Lon)})
	}
	return items
}

func buildTrainCodes(rng *rand.Rand, n int) []string {
	if n < 1 {
		n = 1
	}
	prefixes := []string{"G", "D", "K"}
	seen := make(map[string]struct{}, n)
	items := make([]string, 0, n)
	for len(items) < n {
		p := prefixes[rng.Intn(len(prefixes))]
		num := rng.Intn(9000) + 1000
		code := fmt.Sprintf("%s%d", p, num)
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		items = append(items, code)
	}
	sort.Strings(items)
	return items
}

func buildPlans(rng *rand.Rand, routes []route, trainsPerDay int, minHour, maxHour int) []plan {
	if trainsPerDay < 1 {
		trainsPerDay = 1
	}

	slots := scheduleSlots(minHour, maxHour)
	if len(slots) == 0 {
		slots = []int{minHour * 60}
	}

	trainCodes := buildTrainCodes(rng, trainsPerDay)
	sort.Slice(routes, func(i, j int) bool { return routes[i].Km < routes[j].Km })

	plans := make([]plan, 0, trainsPerDay)
	for i := 0; i < trainsPerDay; i++ {
		r := routes[i%len(routes)]
		departMinute := slots[(i*7)%len(slots)]
		trainType := chooseTrainTypeByDistance(rng, r.Km)
		runtime := runtimeByDistanceKm(rng, trainType, r.Km)
		plans = append(plans, plan{
			TrainCode:        trainCodes[i],
			TrainType:        trainType,
			DepartureStation: r.Dep,
			ArrivalStation:   r.Arr,
			DistanceKm:       r.Km,
			DepartMinute:     departMinute,
			RuntimeMinutes:   runtime,
		})
	}
	return plans
}

func dateOnly(t time.Time) time.Time {
	loc := t.Location()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func makeTrainID(trainCode string, serviceDate time.Time, departure time.Time) string {
	return fmt.Sprintf("%s_%s_%s", trainCode, serviceDate.Format("20060102"), departure.Format("1504"))
}

func buildSeatsForTrain(rng *rand.Rand, trainID string, trainType string, distanceKm float64, density float64) (model.JSON, []model.SeatInfo) {
	layout := map[string]int{}
	seats := make([]model.SeatInfo, 0, 1024)
	carriage := 1

	addSeat := func(seatType, carriageNum, seatNum string, price float64) {
		seats = append(seats, model.SeatInfo{
			ID:          fmt.Sprintf("%s_%s_%s", trainID, carriageNum, seatNum),
			TrainID:     trainID,
			CarriageNum: carriageNum,
			SeatNum:     seatNum,
			SeatType:    seatType,
			SeatPrice:   round2(price),
			Status:      "AVAILABLE",
		})
		layout[seatType]++
	}

	if trainType == "高铁" || trainType == "动车" {
		businessCars := 1
		firstCars := 2
		secondCars := 6
		bizRows := scaledRows(10, density)
		firstRows := scaledRows(12, density)
		secondRows := scaledRows(15, density)

		for i := 0; i < businessCars; i++ {
			carriageNum := fmt.Sprintf("%02d", carriage)
			for row := 1; row <= bizRows; row++ {
				for _, letter := range []string{"A", "C"} {
					seatNum := fmt.Sprintf("%02d%s", row, letter)
					addSeat("商务座", carriageNum, seatNum, priceByDistance(rng, trainType, "商务座", distanceKm))
				}
			}
			carriage++
		}
		for i := 0; i < firstCars; i++ {
			carriageNum := fmt.Sprintf("%02d", carriage)
			for row := 1; row <= firstRows; row++ {
				for _, letter := range []string{"A", "C", "D", "F"} {
					seatNum := fmt.Sprintf("%02d%s", row, letter)
					addSeat("一等座", carriageNum, seatNum, priceByDistance(rng, trainType, "一等座", distanceKm))
				}
			}
			carriage++
		}
		for i := 0; i < secondCars; i++ {
			carriageNum := fmt.Sprintf("%02d", carriage)
			for row := 1; row <= secondRows; row++ {
				for _, letter := range []string{"A", "B", "C", "D", "F"} {
					seatNum := fmt.Sprintf("%02d%s", row, letter)
					addSeat("二等座", carriageNum, seatNum, priceByDistance(rng, trainType, "二等座", distanceKm))
				}
			}
			carriage++
		}
	} else {
		hardSeatCars := 6
		hardSleeperCars := 4
		softSleeperCars := 2
		hardSeatRows := scaledRows(18, density)
		hardSleeperCompart := scaledRows(10, density)
		softSleeperCompart := scaledRows(8, density)

		for i := 0; i < hardSeatCars; i++ {
			carriageNum := fmt.Sprintf("%02d", carriage)
			for row := 1; row <= hardSeatRows; row++ {
				for _, letter := range []string{"A", "B", "C", "D", "F"} {
					seatNum := fmt.Sprintf("%02d%s", row, letter)
					addSeat("硬座", carriageNum, seatNum, priceByDistance(rng, trainType, "硬座", distanceKm))
				}
			}
			carriage++
		}
		for i := 0; i < hardSleeperCars; i++ {
			carriageNum := fmt.Sprintf("%02d", carriage)
			for comp := 1; comp <= hardSleeperCompart; comp++ {
				for _, seatNum := range []string{"A上", "A中", "A下", "B上", "B中", "B下"} {
					addSeat("硬卧", carriageNum, fmt.Sprintf("%02d%s", comp, seatNum), priceByDistance(rng, trainType, "硬卧", distanceKm))
				}
			}
			carriage++
		}
		for i := 0; i < softSleeperCars; i++ {
			carriageNum := fmt.Sprintf("%02d", carriage)
			for comp := 1; comp <= softSleeperCompart; comp++ {
				for _, seatNum := range []string{"A上", "A下", "B上", "B下"} {
					addSeat("软卧", carriageNum, fmt.Sprintf("%02d%s", comp, seatNum), priceByDistance(rng, trainType, "软卧", distanceKm))
				}
			}
			carriage++
		}
	}

	seatLayout, _ := model.ToJSON(layout)
	return seatLayout, seats
}

func layoutTotal(m map[string]int) int {
	t := 0
	for _, v := range m {
		t += v
	}
	return t
}

func scheduleSlots(minHour, maxHour int) []int {
	base := []int{6*60 + 30, 7*60 + 5, 7*60 + 40, 8*60 + 10, 8*60 + 45, 9*60 + 20, 10*60 + 0, 10*60 + 35, 11*60 + 10, 11*60 + 45, 12*60 + 20, 13*60 + 0, 13*60 + 35, 14*60 + 10, 14*60 + 45, 15*60 + 20, 16*60 + 0, 16*60 + 35, 17*60 + 10, 17*60 + 45, 18*60 + 20, 19*60 + 0, 19*60 + 35, 20*60 + 10, 20*60 + 45, 21*60 + 20, 22*60 + 0}
	min := minHour * 60
	max := maxHour*60 + 55
	items := make([]int, 0, len(base))
	for _, m := range base {
		if m >= min && m <= max {
			items = append(items, m)
		}
	}
	return items
}

func chooseTrainTypeByDistance(rng *rand.Rand, km float64) string {
	if km >= 900 {
		return weightedPick(rng, []string{"高铁", "动车", "普速"}, []int{40, 30, 30})
	}
	if km >= 500 {
		return weightedPick(rng, []string{"动车", "高铁", "普速"}, []int{50, 40, 10})
	}
	return weightedPick(rng, []string{"高铁", "动车"}, []int{70, 30})
}

func runtimeByDistanceKm(rng *rand.Rand, trainType string, km float64) uint32 {
	speed := 200.0
	switch trainType {
	case "高铁":
		speed = 280
	case "动车":
		speed = 200
	case "普速":
		speed = 120
	}
	hours := km / speed
	baseMin := hours*60 + 20
	jitter := 0.9 + rng.Float64()*0.2
	minutes := int(math.Round(baseMin * jitter))
	if minutes < 45 {
		minutes = 45
	}
	return uint32(minutes)
}

func priceByDistance(rng *rand.Rand, trainType, seatType string, km float64) float64 {
	perKm := 0.2
	switch trainType {
	case "高铁":
		switch seatType {
		case "二等座":
			perKm = 0.48
		case "一等座":
			perKm = 0.78
		case "商务座":
			perKm = 1.35
		}
	case "动车":
		switch seatType {
		case "二等座":
			perKm = 0.40
		case "一等座":
			perKm = 0.65
		case "商务座":
			perKm = 1.10
		}
	default:
		switch seatType {
		case "硬座":
			perKm = 0.18
		case "硬卧":
			perKm = 0.30
		case "软卧":
			perKm = 0.45
		}
	}
	factor := 0.9 + rng.Float64()*0.2
	price := km*perKm*factor + 10
	if price < 20 {
		price = 20
	}
	return price
}

func scaledRows(base int, density float64) int {
	v := int(math.Round(float64(base) * density))
	if v < 6 {
		v = 6
	}
	if v > base {
		v = base
	}
	return v
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func weightedPick(rng *rand.Rand, items []string, weights []int) string {
	sum := 0
	for _, w := range weights {
		sum += w
	}
	r := rng.Intn(sum)
	acc := 0
	for i, w := range weights {
		acc += w
		if r < acc {
			return items[i]
		}
	}
	return items[len(items)-1]
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func stationNames(stations []station) []string {
	items := make([]string, 0, len(stations))
	for _, s := range stations {
		items = append(items, s.Name)
	}
	return items
}

func buildStationPasses(rng *rand.Rand, trainID string, depTime, arrTime time.Time, depStation, arrStation string, pool []string) []model.TrainStationPass {
	if depStation == arrStation {
		return nil
	}

	runtimeMin := int(arrTime.Sub(depTime).Minutes())
	if runtimeMin < 30 {
		runtimeMin = 30
	}

	stopCount := 0
	if runtimeMin >= 240 {
		stopCount = rng.Intn(4) + 2
	} else if runtimeMin >= 120 {
		stopCount = rng.Intn(3) + 1
	} else {
		stopCount = rng.Intn(2)
	}

	seen := map[string]struct{}{depStation: {}, arrStation: {}}
	mids := make([]string, 0, stopCount)
	for len(mids) < stopCount {
		s := pool[rng.Intn(len(pool))]
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		mids = append(mids, s)
	}

	stations := make([]string, 0, 2+len(mids))
	stations = append(stations, depStation)
	stations = append(stations, mids...)
	stations = append(stations, arrStation)

	segments := len(stations) - 1
	if segments <= 0 {
		return nil
	}
	baseSeg := runtimeMin / segments
	remain := runtimeMin - baseSeg*segments

	cur := depTime
	items := make([]model.TrainStationPass, 0, len(stations))
	for i, name := range stations {
		seq := uint32(i + 1)
		if i == 0 {
			items = append(items, model.TrainStationPass{
				TrainID:       trainID,
				StationName:   name,
				Sequence:      seq,
				ArrivalTime:   sql.NullTime{Valid: false},
				DepartureTime: model.NullDate(cur),
				StopMinutes:   0,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			})
			continue
		}

		segMin := baseSeg
		if remain > 0 {
			segMin++
			remain--
		}
		if segMin < 20 {
			segMin = 20
		}
		arr := cur.Add(time.Duration(segMin) * time.Minute)
		if i == len(stations)-1 {
			items = append(items, model.TrainStationPass{
				TrainID:       trainID,
				StationName:   name,
				Sequence:      seq,
				ArrivalTime:   model.NullDate(arr),
				DepartureTime: sql.NullTime{Valid: false},
				StopMinutes:   0,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			})
			break
		}

		stopMin := uint32(rng.Intn(7) + 2)
		dep := arr.Add(time.Duration(stopMin) * time.Minute)
		items = append(items, model.TrainStationPass{
			TrainID:       trainID,
			StationName:   name,
			Sequence:      seq,
			ArrivalTime:   model.NullDate(arr),
			DepartureTime: model.NullDate(dep),
			StopMinutes:   stopMin,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
		cur = dep
	}

	return items
}
