package grpc

import (
	"encoding/json"
	"grpcroutergen/helper"
	pb "grpcroutergen/proto"
	"grpcroutergen/structs"
	rpcStructs "grpcroutergen/structs/api/grpc"
)

// RPCtrlValidate - RPCtrlValidate Controllers
func RPCtrlValidate(
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
		helper.CheckErr("failed unmarshal @RPCtrlValidate", err)
		structs.ErrorCode.UnexpectedError.String(&errRPCCode.Error)
		return
	}

	res.ID = req.ID
	res.Res = "response"
	resBy, err := json.Marshal(res)
	if err != nil {
		helper.CheckErr("failed marshal &RPCtrlValidate", err)
		structs.ErrorCode.UnexpectedError.String(&errRPCCode.Error)
		return
	}

	*body = resBy
}