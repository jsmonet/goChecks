package validify

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Port checks for a valid port number
func Port(port int) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("Port out of range")

	}
	return true, errors.New("port: valid")
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

// Host just makes sure the host flag isn't empty before moving on
// func Host(address string) (bool, error) {

// }
