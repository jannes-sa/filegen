package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"log"
	"os"
)

var (
	tableName         = flag.String("table", "", "table Name")
	pathDBStruct      string
	pathFileTableName string
	driver            = "postgres"
)

func init() {
	err := os.Setenv("APPPATH", os.Getenv("GOPATH")+"/src/"+os.Getenv("GOAPP")+"/")
	if err != nil {
		checkErr(err)
		panic(1)
		return
	}

	err = godotenv.Load(os.Getenv("APPPATH") + "conf/env")
	if err != nil {
		checkErr(err)
		panic(1)
		return
	}

	flag.Parse()
	check()

	pathDBStruct = os.Getenv("APPPATH") + "structs/db/" + *tableName + ".go"
	pathFileTableName = os.Getenv("APPPATH") + "helper/constant/tablename/" + *tableName + ".go"
}

func check() {
	if *tableName == "" {
		fmt.Println("Define -table params")
		os.Exit(1)
		return
	}
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
}

// postgres://postgres:root@127.0.0.1:5432/postgres?sslmode=disable
func main() {

	generateSQLIntoStruct(os.Getenv("CRED_PGSQL"))
	removeFile()
	createConstantTableName()

}

func createConstantTableName() {
	str := `package tablename

const (
	` + strcase.ToCamel(*tableName) + ` string = "` + *tableName + `"
)`

	createFile(str, pathFileTableName)
}

func generateSQLIntoStruct(sqlCred string) {
	cmd := `bee generate appcode -driver="` + driver + `" -tables="` + *tableName + `" -conn="` + sqlCred + `" -level=1 && Y`
	exec.Command(
		"sh",
		"-c",
		cmd).Output()

	str := appendStructFile(readFile())
	createFile(str, pathDBStruct)
}

func appendStructFile(str string) string {
	str += `
func (u *` + strcase.ToCamel(*tableName) + `) TableName() string {
	return "` + *tableName + `"
}

func init() {
	orm.RegisterModel(new(` + strcase.ToCamel(*tableName) + `))
}`

	str = `package db

import "github.com/astaxie/beego/orm"
` + str

	return str
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createFile(str string, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	fmt.Fprintf(file, str)
}

func readFile() string {
	dat, err := ioutil.ReadFile("models/" + *tableName + ".go")
	checkErr(err)
	if err != nil {
		panic(1)
	}

	newStr := strings.Replace(string(dat), "package models", "", -1)
	return newStr
}

func removeFile() {
	err := os.Remove("models/" + *tableName + ".go")
	if err != nil {
		panic(err)
	}
}

func moveFile() {
	err := os.Rename("models/"+*tableName+".go", pathDBStruct)
	if err != nil {
		panic(err)
	}
}
