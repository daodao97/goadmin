package util

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

const DateLayout = "2006-01-02"

const DateTimeLayout = "2006-01-02 15:04:05"

const HMLayout = "15:04"

func NewTime(t time.Time) Time {
	return Time(t)
}

type Time time.Time

func (t Time) Date() string {
	return time.Time(t).Format(DateLayout)
}

func (t Time) DateTomorrow() string {
	return time.Time(t).AddDate(0, 0, 1).Format(DateLayout)
}

func (t Time) DateYesToday() string {
	return time.Time(t).AddDate(0, 0, -1).Format(DateLayout)
}

func (t Time) Datetime() string {
	return time.Time(t).Format(DateTimeLayout)
}

func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format(DateTimeLayout))
	return []byte(formatted), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(`"`+DateTimeLayout+`"`, string(b))
	if err != nil {
		return err
	}
	*t = Time(tm)
	return nil
}

type Timestamp uint64 // 10位

func (t Timestamp) MarshalJSON() ([]byte, error) {
	tmp := time.Unix(cast.ToInt64(int64(t)), 0)
	formatted := fmt.Sprintf("\"%s\"", tmp.Format(DateTimeLayout))
	return []byte(formatted), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	tmp := cast.ToUint64(string(b))
	*t = Timestamp(tmp)
	return nil
}
