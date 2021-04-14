package mysql

import (
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"testing"
)

func Test_DB_Prepare(t *testing.T) {

	t.Skip()

	err := logger.Init()
	if err != nil {
		t.Fatal("Failed to init logger!")
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	err = Init()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		_ = Db.Close()
	}()
}
