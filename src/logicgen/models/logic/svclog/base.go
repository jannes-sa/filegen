package svclog

import (
	"logicmongogen/helper/constant"
	iSvclog "logicmongogen/models/db/interfaces/svclog"
	pgSvclog "logicmongogen/models/db/pgsql/svclog"
	stSvclog "logicmongogen/models/stub/svclog"
)

var (
	logicName = "@svclog"
	DBSvclog  iSvclog.ISvclog
)

func init() {
	if constant.GOENV == constant.DEVCI {
		DBSvclog = new(stSvclog.Svclog)
	} else {
		DBSvclog = new(pgSvclog.Svclog)
	}
}