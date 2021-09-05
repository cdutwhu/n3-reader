package n3reader

import (
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func Check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func SetIfValidStr(s *string, val string, f func(subject string) (bool, error)) error {
	ok, err := f(val)
	if ok {
		*s = val
		return nil
	}
	return err
}

func SetIfNotEmpty(s *string, val, dVal string) {
	switch strings.Trim(val, " \t") {
	case "":
		*s = dVal
	default:
		*s = val
	}
}

func SetIfNotZero(n *int, val, dVal int) {
	switch val {
	case 0:
		*n = dVal
	default:
		*n = val
	}
}

//
// checks provided nats topic only has alphanumeric & dot separators within the name
//
var topicRegex = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9\-.]*[A-Za-z0-9])?$`)

//
// do regex check on topic names provided for nats
//
func ValidateNatsSubject(subject string) (bool, error) {
	valid := topicRegex.Match([]byte(subject))
	if valid {
		return valid, nil
	}
	return false, errors.New("Nats topic must be alphanumeric only, can also contain (but not start or end with) period ( . ) as token delimiter.")
}
