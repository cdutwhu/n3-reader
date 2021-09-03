package n3reader

import (
	"regexp"

	"github.com/pkg/errors"
)

//
// checks provided nats topic only has alphanumeric & dot separators within the name
//
var topicRegex = regexp.MustCompile("^[A-Za-z0-9]([A-Za-z0-9.]*[A-Za-z0-9])?$")

//
// do regex check on topic names provided for nats
//
func ValidateNatsTopic(tName string) (bool, error) {
	valid := topicRegex.Match([]byte(tName))
	if valid {
		return valid, nil
	}
	return false, errors.New("Nats topic must be alphanumeric only, can also contain (but not start or end with) period ( . ) as token delimiter.")
}
