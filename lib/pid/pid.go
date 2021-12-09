package pid

import (
	"hcc/piano/lib/fileutil"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

var pianoPIDFileLocation = "/var/run"
var pianoPIDFile = "/var/run/piano.pid"

// IsPianoRunning : Check if piano is running
func IsPianoRunning() (running bool, pid int, err error) {
	if _, err := os.Stat(pianoPIDFile); os.IsNotExist(err) {
		return false, 0, nil
	}

	pidStr, err := ioutil.ReadFile(pianoPIDFile)
	if err != nil {
		return false, 0, err
	}

	pianoPID, _ := strconv.Atoi(string(pidStr))

	proc, err := os.FindProcess(pianoPID)
	if err != nil {
		return false, 0, err
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true, pianoPID, nil
	}

	return false, 0, nil
}

// WritePianoPID : Write piano PID to "/var/run/piano.pid"
func WritePianoPID() error {
	pid := os.Getpid()

	err := fileutil.CreateDirIfNotExist(pianoPIDFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(pianoPIDFile, strconv.Itoa(pid))
	if err != nil {
		return err
	}

	return nil
}

// DeletePianoPID : Delete the piano PID file
func DeletePianoPID() error {
	err := fileutil.DeleteFile(pianoPIDFile)
	if err != nil {
		return err
	}

	return nil
}
