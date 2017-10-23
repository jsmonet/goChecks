package checkify

import "errors"

func Port(port int) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("Port out of range")

	}

	return true, errors.New("no errors")

}
