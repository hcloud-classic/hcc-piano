package syscheck

import (
	"errors"
	"os"
)

// CheckRoot : Check root permission
func CheckRoot() error {
	if os.Geteuid() != 0 {
		return errors.New("Please run as root authority.")
	}

	return nil
}
