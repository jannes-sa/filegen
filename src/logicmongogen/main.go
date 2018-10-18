package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
)

var (
	table = flag.String("table", "", "table name")
	logic = flag.String("logic", "", "logic name")

	nmDBStruct   string
	dbStructFile string

	logicFile string

	dirInterfaceFile string
	interfaceFile    string

	dirModelDBFile string
	modelDBFile    string

	dirModelLogicFile string
	modelLogicFile    string

	dirModelStubFile string
	modelStubFile    string
)

func init() {
	err := os.Setenv("APPPATH", os.Getenv("GOPATH")+"/src/"+os.Getenv("GOAPP")+"/")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	flag.Parse()
	checkEnv()

	nmDBStruct = *table + ".go"
	logicFile = strcase.ToLowerCamel(*logic) + ".go"
	dbStructFile = filepath.Join(os.Getenv("APPPATH"), "structs/db/", nmDBStruct)

	dirInterfaceFile = filepath.Join(os.Getenv("APPPATH"), "models/db/interfaces/", strcase.ToLowerCamel(*logic))
	interfaceFile = filepath.Join(dirInterfaceFile, logicFile)

	dirModelDBFile = filepath.Join(os.Getenv("APPPATH"), "models/db/mongo/"+strcase.ToLowerCamel(*logic))
	modelDBFile = filepath.Join(dirModelDBFile, logicFile)

	dirModelLogicFile = filepath.Join(os.Getenv("APPPATH"), "models/logic/", strcase.ToLowerCamel(*logic))
	modelLogicFile = filepath.Join(dirModelLogicFile, logicFile)

	dirModelStubFile = filepath.Join(os.Getenv("APPPATH"), "models/stub/", strcase.ToLowerCamel(*logic))
	modelStubFile = filepath.Join(dirModelStubFile, logicFile)
	checkFile()
}

func checkEnv() {
	if os.Getenv("GOPATH") == "" {
		fmt.Println("SET GOPATH")
		os.Exit(1)
		return
	}
	if os.Getenv("GOAPP") == "" {
		fmt.Println("SET GOAPP")
		os.Exit(1)
		return
	}
	if *logic == "" {
		fmt.Println("Define -logic={logic_name}")
		os.Exit(1)
		return
	}
	if *table == "" {
		fmt.Println("Define -table={table_name}")
		os.Exit(1)
		return
	}

}

func checkFile() {
	_, err := readFile(dbStructFile)
	if err != nil {
		fmt.Println("FILE", dbStructFile, "Not Exists")
		os.Exit(1)
	}

	_, err = readFile(interfaceFile)
	if err == nil {
		fmt.Println("FILE", interfaceFile, "Already Exists")
		os.Exit(1)
	}

	_, err = readFile(modelDBFile)
	if err == nil {
		fmt.Println("FILE", modelDBFile, "Already Exists")
		os.Exit(1)
	}

	_, err = readFile(modelLogicFile)
	if err == nil {
		fmt.Println("FILE", modelLogicFile, "Already Exists")
		os.Exit(1)
	}

	_, err = readFile(modelStubFile)
	if err == nil {
		fmt.Println("FILE", modelStubFile, "Already Exists")
		os.Exit(1)
	}

}

func main() {
	CrudGen()
}

func CrudGen() {
	writeInterface()
	writeDB()
	writeLogic()
	writeStub()
}

func writeLogic() {
	pkgnm := strcase.ToLowerCamel(*logic)
	app := os.Getenv("GOAPP")
	flnm := strcase.ToCamel(*logic)
	flDB := strcase.ToCamel(*table)

	strBase := `package ` + pkgnm + `

import (
	"` + app + `/helper/constant"
	i` + flnm + ` "` + app + `/models/db/interfaces/` + pkgnm + `"
	mg` + flnm + ` "` + app + `/models/db/mongo/` + pkgnm + `"
	st` + flnm + ` "` + app + `/models/stub/` + pkgnm + `"
)

var (
	logicName = "@` + pkgnm + `"
	DB` + flnm + `  i` + flnm + `.I` + flnm + `
)

func init() {
	if constant.GOENV == constant.DEVCI {
		DB` + flnm + ` = new(st` + flnm + `.` + flnm + `)
	} else {
		DB` + flnm + ` = new(mg` + flnm + `.` + flnm + `)
	}
}`

	strBaseTest := `package ` + pkgnm + `

import (
	"` + app + `/helper/constant"
)

func init() {
	initialize()
}

func initialize() {
	constant.LoadEnv()
}`

	strLogic := `package ` + pkgnm + `

import (
	"` + app + `/helper"
	"` + app + `/structs"
	dbStruct "` + app + `/structs/db"
	lStruct "` + app + `/structs/logic"
)

// GetAll` + flDB + ` - GetAll` + flDB + `
func GetAll` + flDB + `(errCode *[]structs.TypeError) (rows []dbStruct.` + flDB + `) {
	rows, err := DB` + flnm + `.GetAll` + flDB + `()
	if err != nil {
		structs.ErrorCode.DatabaseError.String(errCode, err.Error(), "GetAll` + flDB + ` ", logicName)
	}

	return
}

// Insert` + flDB + ` - Insert` + flDB + `
func Insert` + flDB + `(
	contextStruct lStruct.ContextStruct,
	errCode *[]structs.TypeError,
) {
	var (
		nmFunc = "Insert` + flDB + `"
		row    dbStruct.` + flDB + `
	)

	row.JobID = contextStruct.JobID
	row.Req = "req"
	row.Res = "res"
	row.Errcode = "ERRCODE"
	row.Type = "http"

	if err != nil {
		structs.ErrorCode.DatabaseError.String(errCode, err.Error(), nmFunc, logicName)
		return
	}

	_, err = DB` + flnm + `.Insert` + flDB + `(db, &row)
	if err != nil {
		structs.ErrorCode.DatabaseError.String(errCode, err.Error(), nmFunc, logicName)
		helper.CheckErr(nmFunc+" "+logicName, err)

		return
	}
}`

	strLogicTest := `package ` + pkgnm + `

import (
	"` + app + `/structs"
	lStruct "` + app + `/structs/logic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInsert` + flDB + `(t *testing.T) {
	errCode := make([]structs.TypeError, 0)
	var ctxStruct lStruct.ContextStruct
	ctxStruct.JobID = "jobid1"

	Insert` + flDB + `(ctxStruct, &errCode)

	Convey("TestInsert` + flDB + `", t, func() {
		Convey("Should Success", func() {
			So(len(errCode), ShouldEqual, 0)
		})
	})
}

func TestGetAll` + flDB + `(t *testing.T) {
	errCode := make([]structs.TypeError, 0)
	rows := GetAll` + flDB + `(&errCode)

	Convey("TestGetAll` + flDB + `", t, func() {
		Convey("Should Success", func() {
			So(len(errCode), ShouldEqual, 0)
			So(len(rows), ShouldNotBeEmpty)
		})
	})
}`

	os.MkdirAll(dirModelLogicFile, os.ModePerm)
	createFile(strBase, filepath.Join(dirModelLogicFile, "base.go"))
	createFile(strBaseTest, filepath.Join(dirModelLogicFile, "base_test.go"))
	createFile(strLogic, filepath.Join(dirModelLogicFile, pkgnm+".go"))
	createFile(strLogicTest, filepath.Join(dirModelLogicFile, pkgnm+"_test.go"))
}

func writeStub() {
	pkgnm := strcase.ToLowerCamel(*logic)
	app := os.Getenv("GOAPP")
	flnm := strcase.ToCamel(*logic)
	flDB := strcase.ToCamel(*table)

	str := `package ` + pkgnm + `

import (
	"` + app + `/helper/constant/tablename"
	dbStruct "` + app + `/structs/db"

	"github.com/astaxie/beego/orm"
)

// ` + flnm + ` - Logic Struct DB
type ` + flnm + ` struct{}

var tbl` + flDB + ` = tablename.` + flDB + `

// GetAll` + flDB + ` - GetAll` + flDB + ` GetAll
func (d *` + flnm + `) GetAll` + flDB + `() (rows []dbStruct.` + flDB + `, err error) {
	row := dbStruct.` + flDB + `{
		JobID:   "job1",
		Req:     "req",
		Res:     "res",
		Errcode: "errcode",
		Type:    "http",
	}
	rows = append(rows, row)
	return
}

// Insert` + flDB + ` - Insert` + flDB + ` Insert
func (d *` + flnm + `) Insert` + flDB + `(o orm.Ormer, v interface{}) (cnt int64, err error) {

	return
}`
	os.MkdirAll(dirModelStubFile, os.ModePerm)
	createFile(str, modelStubFile)
}

func writeInterface() {
	str := `package ` + strcase.ToLowerCamel(*logic) + `

import (
	dbStruct "` + os.Getenv("GOAPP") + `/structs/db"

	"github.com/astaxie/beego/orm"
)

// I` + strcase.ToCamel(*logic) + ` - ` + strcase.ToCamel(*logic) + ` Logic Interface
type I` + strcase.ToCamel(*logic) + ` interface {
	GetAll` + strcase.ToCamel(*table) + `() ([]dbStruct.` + strcase.ToCamel(*table) + `, error)
	Insert` + strcase.ToCamel(*table) + `(orm.Ormer, interface{}) (int64, error)
}`

	os.MkdirAll(dirInterfaceFile, os.ModePerm)
	createFile(str, interfaceFile)
}

func writeDB() {
	str := `package ` + strcase.ToLowerCamel(*logic) + `

import (
	"` + os.Getenv("GOAPP") + `/helper"
	"` + os.Getenv("GOAPP") + `/helper/constant"
	"` + os.Getenv("GOAPP") + `/helper/constant/tablename"
	db "` + os.Getenv("GOAPP") + `/models/db/mongo"
	dbStruct "` + os.Getenv("GOAPP") + `/structs/db"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// ` + strcase.ToCamel(*logic) + ` - Logic Struct DB
type ` + strcase.ToCamel(*logic) + ` struct{}

func init() {
	var d ` + strcase.ToCamel(*logic) + `
	d.Index()
}

// GetColl - Get Collection ` + *table + `
func (d *` + strcase.ToCamel(*logic) + `) GetColl() (sess *mgo.Session, coll *mgo.Collection, err error) {
	sess, err = db.Connect()
	if err != nil {
		helper.CheckErr("Failed get collection ` + *table + `", err)
		return
	}

	coll = sess.DB(constant.GOAPP).C(tablename.` + strcase.ToCamel(*table) + `)

	return
}

// Index - Create Index
func (d *` + strcase.ToCamel(*logic) + `) Index() (err error) {
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

// GetAll` + strcase.ToCamel(*table) + ` - GetAll` + strcase.ToCamel(*table) + ` GetAll
func (d *` + strcase.ToCamel(*logic) + `) GetAll` + strcase.ToCamel(*table) + `() (rows []dbStruct.` + strcase.ToCamel(*table) + `, err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}

	err = coll.Find(bson.M{}).All(&rows)

	return
}

// GetOneByJobID` + strcase.ToCamel(*table) + ` - GetOneByJobID` + strcase.ToCamel(*table) + `
func (d *` + strcase.ToCamel(*logic) + `) GetOneByJobID` + strcase.ToCamel(*table) + `() (row dbStruct.` + strcase.ToCamel(*table) + `, err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}

	err = coll.Find(bson.M{"job_id": row.JobID}).One(&row)

	return
}

// UpdateByJobID` + strcase.ToCamel(*table) + ` - UpdateByJobID` + strcase.ToCamel(*table) + `
func (d *` + strcase.ToCamel(*logic) + `) UpdateByJobID` + strcase.ToCamel(*table) + `() (err error) {
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

// Insert` + strcase.ToCamel(*table) + ` - Insert` + strcase.ToCamel(*table) + `
func (d *` + strcase.ToCamel(*logic) + `) Insert` + strcase.ToCamel(*table) + `(v interface{}) (err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}

	err = coll.Insert(v)

	return
}

// DeleteByJobID` + strcase.ToCamel(*table) + ` - DeleteByJobID` + strcase.ToCamel(*table) + `
func (d *` + strcase.ToCamel(*logic) + `) DeleteByJobID` + strcase.ToCamel(*table) + `() (err error) {
	sess, coll, err := d.GetColl()
	defer sess.Close()
	if err != nil {
		return
	}
	selector := bson.M{"job_id": "xxxxxx"}
	err = coll.Remove(selector)

	return
}`

	os.MkdirAll(dirModelDBFile, os.ModePerm)
	createFile(str, modelDBFile)
}

func readFile(pathFile string) (string, error) {
	dat, err := ioutil.ReadFile(pathFile)

	return string(dat), err
}

func createFile(str string, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	fmt.Fprintf(file, str)
}

func removeFile(pathFile string) {
	err := os.Remove(pathFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
