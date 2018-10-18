package svclog

import (
	"logicmongogen/helper"
	"logicmongogen/helper/constant"
	"logicmongogen/helper/constant/tablename"
	db "logicmongogen/models/db/mongo"
	dbStruct "logicmongogen/structs/db"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Svclog - Logic Struct DB
type Svclog struct{}

func init() {
	var d Svclog
	d.Index()
}

// GetColl - Get Collection service_log
func (d *Svclog) GetColl() (sess *mgo.Session, coll *mgo.Collection, err error) {
	sess, err = db.Connect()
	if err != nil {
		helper.CheckErr("Failed get collection service_log", err)
		return
	}

	coll = sess.DB(constant.GOAPP).C(tablename.ServiceLog)

	return
}

// Index - Create Index
func (d *Svclog) Index() (err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}

	index := mgo.Index{
		Key:        []string{"job_id"},
		Unique:     true,  // Prevent two documents from having the same index key
		DropDups:   false, // Drop documents with the same index key as a previously indexed one
		Background: false, // Build index in background and return immediately
		Sparse:     false, // Only index documents containing the Key fields
	}

	err = coll.EnsureIndex(index)

	return
}

// GetAllServiceLog - GetAllServiceLog GetAll
func (d *Svclog) GetAllServiceLog() (rows []dbStruct.ServiceLog, err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}

	err = coll.Find(bson.M{}).All(&rows)

	return
}

// GetOneByJobIDServiceLog - GetOneByJobIDServiceLog
func (d *Svclog) GetOneByJobIDServiceLog() (row dbStruct.ServiceLog, err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}

	err = coll.Find(bson.M{"job_id": row.JobID}).One(&row)

	return
}

// UpdateByJobIDServiceLog - UpdateByJobIDServiceLog
func (d *Svclog) UpdateByJobIDServiceLog() (err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}

	selector := bson.M{"job_id": "xxxxxx"}
	update := bson.M{
		"$set": bson.M{
			"res": "yyyy",
		},
	}

	err = coll.Update(selector, update)

	return
}

// InsertServiceLog - InsertServiceLog
func (d *Svclog) InsertServiceLog(v interface{}) (err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}

	err = coll.Insert(v)

	return
}

// DeleteByJobIDServiceLog - DeleteByJobIDServiceLog
func (d *Svclog) DeleteByJobIDServiceLog() (err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}
	selector := bson.M{"job_id": "xxxxxx"}
	err = coll.Remove(selector)

	return
}