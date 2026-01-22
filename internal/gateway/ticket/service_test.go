package ticket

import (
	"testing"
	"time"
)

func TestNormalizeStationKeyword(t *testing.T) {
	if got := normalizeStationKeyword(" Shang Hai "); got != "shanghai" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := normalizeStationKeyword("北 京"); got != "北京" {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestStationAlias(t *testing.T) {
	if got := stationAlias("sh"); got != "上海" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := stationAlias("BJ"); got != "北京" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := stationAlias("x"); got != "" {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestValidateSeatType(t *testing.T) {
	if v, msg := validateSeatType(""); v != "" || msg != "" {
		t.Fatalf("unexpected: v=%q msg=%q", v, msg)
	}
	if v, msg := validateSeatType("二等座"); v != "二等座" || msg != "" {
		t.Fatalf("unexpected: v=%q msg=%q", v, msg)
	}
	if v, msg := validateSeatType("bad"); v != "" || msg == "" {
		t.Fatalf("unexpected: v=%q msg=%q", v, msg)
	}
}

func TestValidateTextField(t *testing.T) {
	if v, msg := validateTextField(" 北京 ", "出发地"); v != "北京" || msg != "" {
		t.Fatalf("unexpected: v=%q msg=%q", v, msg)
	}
	if v, msg := validateTextField("", "出发地"); v != "" || msg == "" {
		t.Fatalf("unexpected: v=%q msg=%q", v, msg)
	}
	if v, msg := validateTextField("北京@", "出发地"); v != "" || msg == "" {
		t.Fatalf("unexpected: v=%q msg=%q", v, msg)
	}
}

func TestCursorEncodeDecode(t *testing.T) {
	c := trainQueryCursor{DepartureTime: time.Date(2026, 1, 16, 9, 0, 0, 0, time.Local), TrainID: "T1"}
	s := encodeCursor(c)
	got, ok := decodeCursor(s)
	if !ok {
		t.Fatalf("expected ok")
	}
	if got.TrainID != c.TrainID || !got.DepartureTime.Equal(c.DepartureTime) {
		t.Fatalf("unexpected decoded: %#v", got)
	}
	if _, ok := decodeCursor("bad"); ok {
		t.Fatalf("expected not ok")
	}
}

func TestParseTravelDate(t *testing.T) {
	start, end, err := parseTravelDate("2026-01-16")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !end.After(start) || end.Sub(start) != 24*time.Hour {
		t.Fatalf("unexpected range: %v - %v", start, end)
	}
	if _, _, err := parseTravelDate("2026/01/16"); err == nil {
		t.Fatalf("expected err")
	}
}

func TestParseTimeOfDay(t *testing.T) {
	if d, ok := parseTimeOfDay(""); ok || d != 0 {
		t.Fatalf("unexpected")
	}
	if d, ok := parseTimeOfDay("09:30"); !ok || d != 9*time.Hour+30*time.Minute {
		t.Fatalf("unexpected: %v %v", d, ok)
	}
	if _, ok := parseTimeOfDay("99:99"); ok {
		t.Fatalf("expected not ok")
	}
}

func TestApplyDepartTimeWindow_UsesDayStartForEnd(t *testing.T) {
	dayStart := time.Date(2026, 1, 16, 0, 0, 0, 0, time.Local)
	dayEnd := dayStart.Add(24 * time.Hour)

	start, end, msg := applyDepartTimeWindow(dayStart, dayEnd, "09:00", "12:00")
	if msg != "" {
		t.Fatalf("unexpected msg: %q", msg)
	}
	if !start.Equal(time.Date(2026, 1, 16, 9, 0, 0, 0, time.Local)) {
		t.Fatalf("unexpected start: %v", start)
	}
	if !end.Equal(time.Date(2026, 1, 16, 12, 0, 0, 0, time.Local)) {
		t.Fatalf("unexpected end: %v", end)
	}
}

func TestApplyDepartTimeWindow_EndBeforeStart(t *testing.T) {
	dayStart := time.Date(2026, 1, 16, 0, 0, 0, 0, time.Local)
	dayEnd := dayStart.Add(24 * time.Hour)

	_, _, msg := applyDepartTimeWindow(dayStart, dayEnd, "18:00", "09:00")
	if msg == "" {
		t.Fatalf("expected msg")
	}
}

