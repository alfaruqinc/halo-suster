package validation

import (
	"health-record/internal/domain"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func NIP(fl validator.FieldLevel) bool {
	nip := fl.Field().Int()
	role := fl.Param()

	nipStr := strconv.FormatInt(nip, 10)

	// skip last three to last digits
	digitSkip := len(nipStr) - 10
	nip = nip / int64(math.Pow10(digitSkip))

	// check month
	month := nip % 100
	if month < 1 || month > 12 {
		return false
	}
	nip = nip / 100

	// check year
	year := nip % 10000
	if year < 2000 || year > int64(time.Now().Year()) {
		return false
	}
	nip = nip / 10000

	// check gender
	gender := nip % 10
	if gender < 1 || gender > 2 {
		return false
	}
	nip = nip / 10

	// pass if role is not provided
	if role == "" {
		return true
	}
	// check role
	itRole := role == "it" && nip == domain.ITRole
	nurseRole := role == "nurse" && nip == domain.NurseRole
	if itRole {
		return true
	}

	if nurseRole {
		return true
	}

	return false
}

func URL(fl validator.FieldLevel) bool {
	urlString := fl.Field().String()
	if !strings.Contains(urlString, ".") {
		return false
	}
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	u, err := url.Parse(urlString)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func IntLen(fl validator.FieldLevel) bool {
	field := fl.Field().Int()
	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}

	n := 0
	for ; field > 0; field /= 10 {
		n++
	}

	return n == param
}

func ISO8601(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	_, err := time.Parse(time.RFC3339, field)

	return err == nil
}
