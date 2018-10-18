package svclog

import (
	"logicmongogen/helper/constant/tablename"
	dbStruct "logicmongogen/structs/db"

	"github.com/astaxie/beego/orm"
)

// Svclog - Logic Struct DB
type Svclog struct{}

var tblServiceLog = tablename.ServiceLog

// GetAllServiceLog - GetAllServiceLog GetAll
func (d *Svclog) GetAllServiceLog() (rows []dbStruct.ServiceLog, err error) {
	row := dbStruct.ServiceLog{
		JobID:   "job1",
		Req:     "req",
		Res:     "res",
		Errcode: "errcode",
		Type:    "http",
	}
	rows = append(rows, row)
	return
}

// InsertServiceLog - InsertServiceLog Insert
func (d *Svclog) InsertServiceLog(o orm.Ormer, v interface{}) (cnt int64, err error) {

	return
}