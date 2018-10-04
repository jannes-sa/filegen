package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

var (
	pathHTTPRouter  string
	pathHTTPCtrl    string
	pathHTTPCtrlOri string
	pathHTTPTest    string

	ctrl = flag.String("ctrl", "", "controller name")
)

func init() {
	err := os.Setenv("APPPATH", os.Getenv("GOPATH")+"/src/"+os.Getenv("GOAPP")+"/")
	if err != nil {
		log.Fatal(err)
		panic(1)
		return
	}

	err = godotenv.Load(os.Getenv("APPPATH") + "conf/env")
	if err != nil {
		log.Fatal(err)
		panic(1)
		return
	}

	flag.Parse()
	check()

	pathHTTPRouter = os.Getenv("APPPATH") + "routers/http/router.go"
	pathHTTPCtrl = os.Getenv("APPPATH") + "controllers/http/" + *ctrl + ".go"
	pathHTTPCtrlOri = os.Getenv("APPPATH") + "controllers/" + *ctrl + ".go"
	pathHTTPTest = os.Getenv("APPPATH") + "routers/componenttest/http/" + *ctrl + "_test.go"
}

func check() {
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
	if *ctrl == "" {
		fmt.Println("Define -ctrl={controller name}")
		os.Exit(1)
		return
	}
}

func main() {
	RouterGen()
}

func RouterGen() {
	readWriteCtrl()
	readRouterAndWrite()
	writeComponentTest()
}

func writeComponentTest() {
	str := `package http

import (
	"encoding/json"
	"strconv"
	"` + os.Getenv("GOAPP") + `/helper/constant"
	"` + os.Getenv("GOAPP") + `/routers/componenttest"
	"` + os.Getenv("GOAPP") + `/structs"
	httpStructs "` + os.Getenv("GOAPP") + `/structs/api/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/astaxie/beego"
)

func Test` + strings.Title(*ctrl) + `Success(t *testing.T) {

	var req structs.ReqData
	var reqBody httpStructs.ReqTest
	reqBody.ID = 1
	req.ReqBody = reqBody
	by, _ := json.Marshal(req)

	res := componenttest.SendHTTP(
		"POST",
		host+"/"+constant.DOMAINNAME+"/v1/` + *ctrl + `",
		by,
	)

	Convey("Test` + strings.Title(*ctrl) + `Success", t, func() {
		Convey("Should Success", func() {
			So(len(res.Error), ShouldEqual, 0)
		})
	})
}`
	createFile(str, pathHTTPTest)
}

func readWriteCtrl() {
	// Check File Exists //
	_, err := readFile(pathHTTPCtrl)
	if err == nil {
		fmt.Println("FILE ALREADY EXISTS", *ctrl)
		os.Exit(1)
	}
	// Check File Exists //

	// Generate Controller //
	cmd := `bee generate controller ` + *ctrl
	_, err = exec.Command(
		"sh",
		"-c",
		cmd).Output()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// Generate Controller //

	str, err := readFile(pathHTTPCtrlOri)

	sp := `func (c *` + strings.Title(*ctrl) + `Controller)`

	// Append POST Method //
	spPOST := sp + ` Post() {`
	newStrCtrlPOST := `
	errCode := make([]structs.TypeError, 0)
	
	var (
		reqInterface structsAPI.TestInterface
		req          structsAPI.ReqTest
		res          structsAPI.ResTest
	)

	rqBodyByte := helper.GetRqBodyRev(c.Ctx, &errCode)
	if len(errCode) > 0 {
		SendOutput(c.Ctx, c.Data["json"], errCode)
		return
	}

	err := json.Unmarshal(rqBodyByte, &reqInterface)
	if err != nil {
		structs.ErrorCode.UnexpectedError.String(&errCode)
		SendOutput(c.Ctx, c.Data["json"], errCode)
		return
	}

	reqInterface.ValidateRequest(&req, &errCode)
	if len(errCode) > 0 {
		SendOutput(c.Ctx, c.Data["json"], errCode)
		return
	}

	res.ID = req.ID

	c.Data["json"] = res

	SendOutput(c.Ctx, res, errCode)`
	str = addStringAfter(str, spPOST, newStrCtrlPOST)
	// Append POST Method //

	// Append GetOne Method //
	spGetOne := sp + ` GetOne() {`
	newStrCtrlGetOne := `
	errCode := make([]structs.TypeError, 0)
	
	id := c.Ctx.Input.Param(":id")
	beego.Debug(id)

	SendOutput(c.Ctx, c.Data["json"], errCode)
	`
	str = addStringAfter(str, spGetOne, newStrCtrlGetOne)
	// Append GetOne Method //

	// Replace TOP //
	str = strings.Replace(str, `package controllers`, "", -1)
	str = `package http

import (
	"encoding/json"
	"` + os.Getenv("GOAPP") + `/helper"
	"` + os.Getenv("GOAPP") + `/structs"
	structsAPI "` + os.Getenv("GOAPP") + `/structs/api/http"
)
` + str
	// Replace TOP //

	createFile(str, pathHTTPCtrl)
	removeFile(pathHTTPCtrlOri)
}

func addStringAfter(
	str string,
	sp string,
	strApp string,
) string {
	ar1 := strings.Split(str, sp)
	newStr := ar1[0] + sp + strApp + ar1[1]
	return newStr
}

func readRouterAndWrite() {
	str, err := readFile(pathHTTPRouter)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	newStrRouter := `beego.NSNamespace("/` + *ctrl + `",
			beego.NSInclude(
				&ctrl.` + strings.Title(*ctrl) + `Controller{},
			),
		),`
	newStr := inBetweenAddString(
		str,
		"/*:STARTHTTP*/",
		"/*:ENDHTTP*/",
		newStrRouter,
	)

	createFile(newStr, pathHTTPRouter)
}

func inBetweenAddString(
	str string,
	sp1 string,
	sp2 string,
	strApp string,
) string {
	res := strings.Split(str, sp1)
	res2 := strings.Split(res[1], sp2)

	strUp := res[0]
	strIn := res2[0]
	strDo := res2[1]

	newStrRouter := strApp + `
		` + sp2

	strIn += newStrRouter
	strIn = sp1 + strIn

	newStr := strUp + strIn + strDo
	return newStr
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
