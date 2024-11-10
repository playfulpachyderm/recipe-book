package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t Timestamp) Value() (driver.Value, error) {
	return t.UnixMilli(), nil
}

func (t *Timestamp) Scan(src interface{}) error {
	val, is_ok := src.(int64)
	if !is_ok {
		return fmt.Errorf("Incompatible type for Timestamp: %#v", src)
	}
	*t = Timestamp{time.UnixMilli(val)}
	return nil
}

func TimestampFromUnix(num int64) Timestamp {
	return Timestamp{time.Unix(num, 0)}
}
func TimestampFromUnixMilli(num int64) Timestamp {
	return Timestamp{time.UnixMilli(num)}
}
