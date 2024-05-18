package helper

import (
	"math"
	"strconv"
)

func MatchRole(nip, role int64) bool {
	nipLen := len(strconv.FormatInt(nip, 10))

	skipDigit := int64(math.Pow10(nipLen - 3))

	return (nip / skipDigit) == role
}
