package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/segmentio/ksuid"
)

func BoolPtr(b bool) *bool {
	return &b
}

func ConvertToString(v interface{}) string {
	if v == nil {
		return ""
	}

	stringV := fmt.Sprintf("%v", v)
	return stringV
}

func StringContains(s string, strArray []string, caseSensitive bool) bool {
	for _, str := range strArray {
		if caseSensitive {
			if s == str {
				return true
			}
		} else {
			if strings.EqualFold(str, s) {
				return true
			}
		}
	}
	return false
}

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// NewKSUID generates k-sorted random unique IDs
func NewKSUID() string {
	kid, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		panic(err)
	}

	return kid.String()
}

func NewUUID() string {
	return uuid.New().String()
}

func StringToBool(str string) bool {
	return strings.EqualFold(str, "true")
}
