package grpc

import (
	"encoding/json"
	"grpcroutergen/helper"
	"grpcroutergen/helper/timetn"
	"grpcroutergen/structs"
	structsAPI "grpcroutergen/structs/api"
	structsRPC "grpcroutergen/structs/api/grpc"
	"grpcroutergen/thirdparty/rpc"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidate(t *testing.T) {
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
		prefix+"/validate",
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

	Convey("TestValidate", t, func() {
		Convey("Should Success", func() {
			So(err, ShouldEqual, nil)
			So(len(resHeader.Error.Error), ShouldEqual, 0)
		})
	})
}