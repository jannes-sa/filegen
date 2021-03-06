package svclog

import (
	"logicmongogen/helper"
	"logicmongogen/helper/constant"
	"logicmongogen/helper/constant/tablename"
	dbStruct "logicmongogen/structs/db"

	"github.com/astaxie/beego/orm"
)

// Svclog - Logic Struct DB
type Svclog struct{}

var tblServiceLog = tablename.ServiceLog

// GetAllServiceLog - GetAllServiceLog GetAll
func (d *Svclog) GetAllServiceLog() (rows []dbStruct.ServiceLog, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(tblServiceLog).All(&rows)
	return
}

// InsertServiceLog - InsertServiceLog Insert
func (d *Svclog) InsertServiceLog(o orm.Ormer, v interface{}) (cnt int64, err error) {
	cnt, err = o.Insert(v)

	if err.Error() != constant.ExceptionLastInsertID {
		helper.CheckErr("Failed Inserted", err)
		return
	}

	return cnt, nil
}

/*
// GetOneByJobIDServiceLog - GetOneByJobIDServiceLog GetOne
func (d *Svclog) GetOneByJobIDServiceLog(r dbStruct.ServiceLog) (row dbStruct.ServiceLog, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(tblServiceLog).Filter("job_id", r.JobID).One(&row)
	return
}

// UpdateByJobIDServiceLog - UpdateByJobIDServiceLog Update
func (d *Svclog) UpdateByJobIDServiceLog(
	o orm.Ormer,
	row dbStruct.ServiceLog,
) (err error) {

	_, err = o.QueryTable(tblServiceLog).Filter("job_id", row.JobID).Update(orm.Params{
		"req": row.Req,
	})

	return
}

// UpdateReturnByJobIDServiceLog - UpdateReturnByJobIDServiceLog UpdateReturn
func (d *Svclog) UpdateReturnByJobIDServiceLog(
	o orm.Ormer,
	row dbStruct.ServiceLog,
) (rows []dbStruct.ServiceLog, err error) {
	q := []string{
		"UPDATE", tblServiceLog,
		"SET req = ?",
		"WHERE job_id = ?",
		"RETURNING type, job_id, req, res, errcode",
	}
	sql := strings.Join(q, " ")
	_, err = o.Raw(sql, row.Req, row.JobID).QueryRows(&rows)

	return
}

// DeleteByJobIDServiceLog - DeleteByJobIDServiceLog Delete
func (d *Svclog) DeleteByJobIDServiceLog(
	o orm.Ormer,
	row dbStruct.ServiceLog,
) (err error) {
	_, err = o.QueryTable(tblServiceLog).Filter("job_id", row.JobID).Delete()
	return

	return
}*/