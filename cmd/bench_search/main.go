package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var (
		url         = flag.String("url", "http://127.0.0.1:5200/api/v1/train/search?departure_station=上海&arrival_station=北京&travel_date=2026-01-16&seat_type=硬座&limit=20", "请求URL")
		concurrency = flag.Int("c", 50, "并发数")
		duration    = flag.Duration("d", 10*time.Second, "压测时长")
	)
	flag.Parse()

	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(*duration)

	var okCount int64
	var errCount int64
	latCh := make(chan time.Duration, 200000)

	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(deadline) {
				begin := time.Now()
				resp, err := client.Get(*url)
				lat := time.Since(begin)
				if err != nil {
					atomic.AddInt64(&errCount, 1)
					continue
				}
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					atomic.AddInt64(&okCount, 1)
					select {
					case latCh <- lat:
					default:
					}
				} else {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	wg.Wait()
	close(latCh)

	lats := make([]time.Duration, 0, len(latCh))
	for v := range latCh {
		lats = append(lats, v)
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

	p := func(q float64) time.Duration {
		if len(lats) == 0 {
			return 0
		}
		idx := int(float64(len(lats)-1) * q)
		return lats[idx]
	}

	total := okCount + errCount
	rps := float64(total) / duration.Seconds()
	fmt.Printf("total=%d ok=%d err=%d rps=%.1f p50=%s p99=%s\n", total, okCount, errCount, rps, p(0.50), p(0.99))
}

