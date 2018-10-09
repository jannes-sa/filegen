package saveSvcLog

import (
	"crudgen/helper/constant/tablename"
	dbStruct "crudgen/structs/db"

	"github.com/astaxie/beego/orm"
)

// SaveSvcLog - Logic Struct DB
type SaveSvcLog struct{}

var tblServiceLog = tablename.ServiceLog

// GetAllServiceLog - GetAllServiceLog GetAll
func (d *SaveSvcLog) GetAllServiceLog() (rows []dbStruct.ServiceLog, err error) {
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
func (d *SaveSvcLog) InsertServiceLog(o orm.Ormer, v interface{}) (cnt int64, err error) {

	return
}