package svclog

import (
	"logicmongogen/helper/constant"
)

func init() {
	initialize()
}

func initialize() {
	constant.LoadEnv()
}