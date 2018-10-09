package saveSvcLog

import (
	dbStruct "crudgen/structs/db"

	"github.com/astaxie/beego/orm"
)

// ISaveSvcLog - SaveSvcLog Logic Interface
type ISaveSvcLog interface {
	GetAllServiceLog() ([]dbStruct.ServiceLog, error)
	InsertServiceLog(orm.Ormer, interface{}) (int64, error)
}