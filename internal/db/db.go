package db

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
)

// ex.
//
// dbDriver = postgreSQL
//
//	dsn = "host=127.0.0.1 port=5432 user=user password=password dbname=dbname sslmode=disable"
func Setup(dbDriver string, host string, port int, user string, password string, dbname string, sslmode string) (*sql.DB, error) {
	dsn := "host=" + host + " port=" + strconv.Itoa(port) + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		// 接続に失敗した場合は、確保したリソースを解放してからエラーを返す
		db.Close()
		return nil, err
	}
	return db, err
}
