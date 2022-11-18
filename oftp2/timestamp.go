package oftp2

import (
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func NewTimeStamp(c []byte) (Timestamp, error) {
	year, err := strconv.Atoi(string(c[0:4]))
	if err != nil {
		return Timestamp{}, err
	}
	month, err := strconv.Atoi(string(c[4:6]))
	if err != nil {
		return Timestamp{}, err
	}
	day, err := strconv.Atoi(string(c[6:8]))
	if err != nil {
		return Timestamp{}, err
	}
	hour, err := strconv.Atoi(string(c[8:10]))
	if err != nil {
		return Timestamp{}, err
	}
	minute, err := strconv.Atoi(string(c[10:12]))
	if err != nil {
		return Timestamp{}, err
	}
	second, err := strconv.Atoi(string(c[12:14]))
	if err != nil {
		return Timestamp{}, err
	}
	milli, err := strconv.Atoi(string(c[14:18]))
	if err != nil {
		return Timestamp{}, err
	}

	return Timestamp{
		Time: time.Date(year, time.Month(month), day, hour, minute, second, milli, time.UTC),
	}, nil
}

func (t Timestamp) ToString() string {
	month, _ := fillUpInt(int(t.Month()), 2)
	day, _ := fillUpInt(t.Day(), 2)
	hour, _ := fillUpInt(t.Hour(), 2)
	minute, _ := fillUpInt(t.Minute(), 2)
	second, _ := fillUpInt(t.Second(), 2)
	milli, _ := fillUpInt(t.Nanosecond(), 4)
	return strconv.Itoa(t.Year()) +
		month +
		day +
		hour +
		minute +
		second +
		milli
}
