package svclog

import (
	"template/structs"
	lStruct "template/structs/logic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertServiceLog(t *testing.T) {
	errCode := make([]structs.TypeError, 0)
	var ctxStruct lStruct.ContextStruct
	ctxStruct.JobID = "jobid1"

	InsertServiceLog(ctxStruct, &errCode)

	Convey("TestInsertServiceLog", t, func() {
		Convey("Should Success", func() {
			So(len(errCode), ShouldEqual, 0)
		})
	})
}

func TestGetAllServiceLog(t *testing.T) {
	errCode := make([]structs.TypeError, 0)
	rows := GetAllServiceLog(&errCode)

	Convey("TestGetAllServiceLog", t, func() {
		Convey("Should Success", func() {
			So(len(errCode), ShouldEqual, 0)
			So(len(rows), ShouldNotBeEmpty)
		})
	})
}