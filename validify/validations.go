package validify

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Port requires an int and outputs both a bool as well as error content. This test checks that the int falls within the range of acceptable TCP ports, from 1 to 65535.
func Port(port int) (bool, error) {
	if port <= 1 || port >= 65535 {
		return false, errors.New("Port out of range")

	}
	return true, errors.New("port is valid")
}

// Neorole just checks that the user input a valid choice. It case-insensitively checks that the string you feed it contains either 'master' or 'slave'. This pertains specifically to Neo4j HA cluster roles, but not necessarily to Causal cluster roles.
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
	return true, err
}

// Percentages takes two int values and returns bool + error. The first int is always the "warning" value, the second is always the "critical" value. This feeds a test where you can supply arbitrary percentages between 1 and 99 to denote warning and critical consumption values of a resource. 100+ is invalidated by `else if warn > 99 || crit > 99`
func Percentages(warn int, crit int) (bool, error) {
	if warn > crit {
		return false, errors.New("You can't have a warn value greater than your crit value")
	} else if warn > 99 || crit > 99 {
		return false, errors.New("One of your values exceeds 99, the highest applicable percentage for this check")
	}
	return true, errors.New("")
}

// removed rmq queue exists checker. I didn't like it and it isn't worth testing. Let it throw a larger exception
