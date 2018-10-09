package saveSvcLog

import (
	"crudgen/helper/constant"
	db "crudgen/models/db/pgsql"

	_ "github.com/lib/pq"
)

func init() {
	initialize()
}

func initialize() {
	constant.LoadEnv()
	db.RegisterPGSQL()
}