package svclog

import (
	"logicmongogen/helper/constant"
	db "logicmongogen/models/db/pgsql"

	_ "github.com/lib/pq"
)

func init() {
	initialize()
}

func initialize() {
	constant.LoadEnv()
	db.RegisterPGSQL()
}