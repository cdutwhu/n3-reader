package n3reader

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func SetIfValidStr(s *string, val string, f func(subject string) (bool, error)) error {
	ok, err := f(val)
	if ok {
		*s = val
		return nil
	}
	return err
}

func SetIfNotEmpty(s *string, val, dVal string) error {
	switch strings.Trim(val, " \t") {
	case "":
		*s = dVal
	default:
		*s = val
	}
	return nil
}

func SetIfNotZero(n *int, val, dVal int) error {
	switch val {
	case 0:
		*n = dVal
	default:
		*n = val
	}
	return nil
}

//
// checks provided nats topic only has alphanumeric & dot separators within the name
//
var subjectRegex = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9\-.]*[A-Za-z0-9])?$`)

//
// do regex check on topic names provided for nats
//
func ValidateNatsSubject(subject string) (bool, error) {
	valid := subjectRegex.Match([]byte(subject))
	if valid {
		return valid, nil
	}
	return false, errors.New("NATS topic must be alphanumeric only, can also contain (but not start or end with) period ( . ) as token delimiter.")
}
