package svclog

import (
	dbStruct "logicmongogen/structs/db"

	"github.com/astaxie/beego/orm"
)

// ISvclog - Svclog Logic Interface
type ISvclog interface {
	GetAllServiceLog() ([]dbStruct.ServiceLog, error)
	InsertServiceLog(orm.Ormer, interface{}) (int64, error)
}