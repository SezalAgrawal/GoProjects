package utils

import (
	"fmt"
	"time"
)

const RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"

// TimeISOMicro is a Rounded implementation of RFC3339Micro microformat
// This type serves to have WKTs(Well-Known-Types) in as many places as we can
// Note:: t = time.Now(); TimeISOMicro(t).String() may not always equal t.Format(RFC3339Micro)
// due to the rounding off
type TimeISOMicro time.Time

func (t TimeISOMicro) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.String())
	return []byte(stamp), nil
}

func (t TimeISOMicro) String() string {
	// Rounding off has been made the default behaviour from a pragmatic view
	// The initial approach was to implement this such that .String() equals .Format(RFC3339Micro)
	// but in most of the places this was used to represent timestamps that were persisted using postgres
	// hence resulted in a microsecond difference. Adding the round method, results in much cleaner
	// and mentally less taxing code elsewhere.
	return time.Time(t).Round(time.Microsecond).Format(RFC3339Micro)
}

func (t *TimeISOMicro) AsTime() *time.Time {
	if t == nil {
		return nil
	}
	tt := time.Time(*t)
	return &tt
}

func TimeISOMicroPtr(t *time.Time) *TimeISOMicro {
	if t == nil {
		return nil
	}
	tm := TimeISOMicro(*t)
	return &tm
}
