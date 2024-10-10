package config

import (
	"github.com/sistemakreditasi/backend-akreditasi/helper"
)

var MongoString string = GetEnv("MONGOSTRING")

var mongoinfo = helper.DBInfo{
	DBString: MongoString,
	DBName:   "db_akreditasi",
}

var Mongoconn, ErrorMongoconn = helper.MongoConnect(mongoinfo)
