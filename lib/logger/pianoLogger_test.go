package logger

import (
	"testing"

	"innogrid.com/hcloud-classic/hcc_errors"
)

func Test_CreateDirIfNotExist(t *testing.T) {
	err := CreateDirIfNotExist("/var/log/" + LogName)
	if err != nil {
		t.Fatal("Failed to create dir!")
	}
}

func Test_Logger_Prepare(t *testing.T) {
	err := Init()
	if err != nil {
		hcc_errors.SetErrLogger(Logger)
		hcc_errors.NewHccError(hcc_errors.PianoInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	hcc_errors.SetErrLogger(Logger)

	defer func() {
		End()
	}()
}
