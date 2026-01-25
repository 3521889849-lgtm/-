package ticket

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type trainQueryCursor struct {
	DepartureTime time.Time `json:"departure_time"`
	TrainID       string    `json:"train_id"`
}

func encodeCursor(c trainQueryCursor) string {
	b, _ := json.Marshal(c)
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeCursor(s string) (trainQueryCursor, bool) {
	if strings.TrimSpace(s) == "" {
		return trainQueryCursor{}, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return trainQueryCursor{}, false
	}
	var c trainQueryCursor
	if err := json.Unmarshal(raw, &c); err != nil {
		return trainQueryCursor{}, false
	}
	if c.TrainID == "" || c.DepartureTime.IsZero() {
		return trainQueryCursor{}, false
	}
	return c, true
}

