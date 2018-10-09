package saveSvcLog

import (
	"crudgen/helper/constant"
	iSaveSvcLog "crudgen/models/db/interfaces/saveSvcLog"
	pgSaveSvcLog "crudgen/models/db/pgsql/saveSvcLog"
	stSaveSvcLog "crudgen/models/stub/saveSvcLog"
)

var (
	logicName = "@saveSvcLog"
	DBSaveSvcLog  iSaveSvcLog.ISaveSvcLog
)

func init() {
	if constant.GOENV == constant.DEVCI {
		DBSaveSvcLog = new(stSaveSvcLog.SaveSvcLog)
	} else {
		DBSaveSvcLog = new(pgSaveSvcLog.SaveSvcLog)
	}
}