package model

import (
	"database/sql"
	"time"
)

func NullDate(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}

