package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ctrl = flag.String("ctrl", "", "controller name")

	pathGRPCCtrl   string
	pathGRPCRouter string
	pathGRPCTest   string
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

	pathGRPCCtrl = os.Getenv("APPPATH") + "controllers/grpc/" + *ctrl + ".go"
	pathGRPCRouter = os.Getenv("APPPATH") + "routers/grpc/router.go"
	pathGRPCTest = os.Getenv("APPPATH") + "routers/componenttest/grpc/" + *ctrl + "_test.go"
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
	writeCtrl()
	readRouterAndWrite()
	writeComponentTest()
}

func writeComponentTest() {
	str := `package grpc

import (
	"encoding/json"
	"` + os.Getenv("GOAPP") + `/helper"
	"` + os.Getenv("GOAPP") + `/helper/timetn"
	"` + os.Getenv("GOAPP") + `/structs"
	structsAPI "` + os.Getenv("GOAPP") + `/structs/api"
	structsRPC "` + os.Getenv("GOAPP") + `/structs/api/grpc"
	"` + os.Getenv("GOAPP") + `/thirdparty/rpc"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test` + strings.Title(*ctrl) + `(t *testing.T) {
	reqID := helper.GenJobID()

	var errorHeader structs.TypeGRPCError
	header := structsRPC.TypeHeaderRPC{
		ReqID:       reqID,
		Date:        timetn.Now(),
		ContentType: "application/grpc",
		RoundTrip:   "",
		Error:       errorHeader,
	}
	headerByte, _ := json.Marshal(header)

	var req structsRPC.ReqTest
	req.ID = 1
	req.Data = "requestdata"
	reqBy, _ := json.Marshal(req)

	var tracer structsAPI.HeaderTracer
	tracer.ParSpanID = "ParSpanID"
	tracer.SpanID = "SpanID"
	tracer.TraceID = "TraceID"
	tracer.XReqID = "XReqID"

	resp, err := rpc.SendGRPCComponentTest(
		prefix+"/` + *ctrl + `",
		host,
		reqBy,
		headerByte,
		reqID,
		tracer,
	)

	var resHeader structsRPC.TypeHeaderRPC
	json.Unmarshal(resp.Header, &resHeader)

	var resBody structsRPC.ResTest
	json.Unmarshal(resp.Body, &resBody)

	Convey("Test` + strings.Title(*ctrl) + `", t, func() {
		Convey("Should Success", func() {
			So(err, ShouldEqual, nil)
			So(len(resHeader.Error.Error), ShouldEqual, 0)
		})
	})
}`
	createFile(str, pathGRPCTest)
}

func readRouterAndWrite() {
	str, err := readFile(pathGRPCRouter)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	newStrRouter := `prefix + "/` + *ctrl + `": ctrl.RPCtrl` + strings.Title(*ctrl) + `,`
	newStr := inBetweenAddString(
		str,
		"/*:STARTGRPC*/",
		"/*:ENDGRPC*/",
		newStrRouter,
	)

	createFile(newStr, pathGRPCRouter)
}

func writeCtrl() {
	_, err := readFile(pathGRPCCtrl)
	if err == nil {
		fmt.Println("FILE ALREADY EXISTS", *ctrl)
		os.Exit(1)
		return
	}

	str := `package grpc

import (
	"encoding/json"
	"` + os.Getenv("GOAPP") + `/helper"
	pb "` + os.Getenv("GOAPP") + `/proto"
	"` + os.Getenv("GOAPP") + `/structs"
	rpcStructs "` + os.Getenv("GOAPP") + `/structs/api/grpc"
)

// RPCtrl` + strings.Title(*ctrl) + ` - RPCtrl` + strings.Title(*ctrl) + ` Controllers
func RPCtrl` + strings.Title(*ctrl) + `(
	in *pb.DoReq,
	errRPCCode *structs.TypeGRPCError,
	body *[]byte,
) {

	var (
		req rpcStructs.ReqTest
		res rpcStructs.ResTest
	)

	err := json.Unmarshal(in.GetBody(), &req)
	if err != nil {
		helper.CheckErr("failed unmarshal @RPCtrl` + strings.Title(*ctrl) + `", err)
		structs.ErrorCode.UnexpectedError.String(&errRPCCode.Error)
		return
	}

	res.ID = req.ID
	res.Res = "response"
	resBy, err := json.Marshal(res)
	if err != nil {
		helper.CheckErr("failed marshal &RPCtrl` + strings.Title(*ctrl) + `", err)
		structs.ErrorCode.UnexpectedError.String(&errRPCCode.Error)
		return
	}

	*body = resBy
}`

	createFile(str, pathGRPCCtrl)
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
