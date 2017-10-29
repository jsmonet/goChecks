package validify

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Port checks for a valid port number
func Port(port int) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("Port out of range")

	}
	return true, errors.New("port is valid")
}

// Neorole just checks that the user input a valid choice. This is just to give a friendly error message
func Neorole(role string) (bool, error) {
	lowerCaseRole := strings.ToLower(role)
	if lowerCaseRole != "master" && lowerCaseRole != "slave" {
		errString := fmt.Sprintf("This switch, checking Neo4j, only accepts the following arguments: master, slave. You entered %v", role)
		return false, errors.New(errString)
	}
	return true, errors.New("role: valid")
}

// Authb64 tries to decode the auth string and bombs if you messed up. this deliberately throws away decoded data
func Authb64(authstring string) (bool, error) {
	_, err := base64.StdEncoding.DecodeString(authstring)
	if err != nil {
		return false, errors.New("String does not decode, please enter a valid base64 auth string")
	}
	return true, errors.New("auth: valid")
}

// Percentages returns a bool and error message
func Percentages(warn int, crit int) (bool, error) {
	if warn > crit {
		return false, errors.New("You can't have a warn value greater than your crit value")
	} else if warn > 99 || crit > 99 {
		return false, errors.New("One of your values exceeds 99, the highest applicable percentage for this check")
	}
	return true, errors.New("")
}

// RmqQueueExists grabs a status code and returns a bool, and also doesn't QUITE work the way I want right now
func RmqQueueExists(target string) (result bool, err error) {
	res, resErr := http.Get(target)
	if res.StatusCode != 200 {
		return false, resErr
	}
	return false, resErr
}
