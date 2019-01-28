package types

import (
	"encoding/json"
	"time"
)

type UnixTime time.Time

func (t UnixTime) Time() time.Time {
	return time.Time(t)
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time().Unix())
}

func (t *UnixTime) UnmarshalJSON(bytes []byte) error {
	var sec int64
	if err := json.Unmarshal(bytes, &sec); err != nil {
		return err
	}
	*t = UnixTime(time.Unix(sec, 0))
	return nil
}
