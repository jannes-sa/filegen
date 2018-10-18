package svclog

import (
	"logicmongogen/helper"
	dbBase "logicmongogen/models/db/pgsql"
	"logicmongogen/structs"
	dbStruct "logicmongogen/structs/db"
	lStruct "logicmongogen/structs/logic"
)

// GetAllServiceLog - GetAllServiceLog
func GetAllServiceLog(errCode *[]structs.TypeError) (rows []dbStruct.ServiceLog) {
	rows, err := DBSvclog.GetAllServiceLog()
	if err != nil {
		structs.ErrorCode.DatabaseError.String(errCode, err.Error(), "GetAllServiceLog ", logicName)
	}

	return
}

// InsertServiceLog - InsertServiceLog
func InsertServiceLog(
	contextStruct lStruct.ContextStruct,
	errCode *[]structs.TypeError,
) {
	var (
		nmFunc = "InsertServiceLog"
		row    dbStruct.ServiceLog
	)

	row.JobID = contextStruct.JobID
	row.Req = "req"
	row.Res = "res"
	row.Errcode = "ERRCODE"
	row.Type = "http"

	db := dbBase.Session()
	err := db.Begin()
	if err != nil {
		structs.ErrorCode.DatabaseError.String(errCode, err.Error(), nmFunc, logicName)
		return
	}

	_, err = DBSvclog.InsertServiceLog(db, &row)
	if err != nil {
		structs.ErrorCode.DatabaseError.String(errCode, err.Error(), nmFunc, logicName)
		err = db.Rollback()
		helper.CheckErr(nmFunc+" "+logicName, err)

		return
	}

	err = db.Commit()
	if err != nil {
		structs.ErrorCode.DatabaseError.String(errCode, err.Error(), nmFunc, logicName)
		err = db.Rollback()
		helper.CheckErr(nmFunc+" "+logicName, err)

		return
	}
}