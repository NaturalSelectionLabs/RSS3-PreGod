package timex

import (
	"encoding/json"
	"fmt"
	"time"
)

var (
	_ json.Marshaler   = &Time{}
	_ json.Unmarshaler = &Time{}
)

type Time time.Time

func (t *Time) UnmarshalJSON(bytes []byte) error {
	internalTime, err := time.Parse(ISO8601, string(bytes))
	if err != nil {
		return err
	}

	*t = Time(internalTime)

	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(*t).Format(ISO8601))), nil
}

func (t *Time) Time() time.Time {
	return time.Time(*t)
}
