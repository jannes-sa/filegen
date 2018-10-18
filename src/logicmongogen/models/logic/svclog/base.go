package svclog

import (
	"logicmongogen/helper/constant"
	iSvclog "logicmongogen/models/db/interfaces/svclog"
	mgSvclog "logicmongogen/models/db/mongo/svclog"
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
		DBSvclog = new(mgSvclog.Svclog)
	}
}