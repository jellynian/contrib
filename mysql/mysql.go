package mysql

import (
	"fmt"
	"sync"

	"database/sql"

	"bitbucket.org/jellynian/contrib/config"

	"bitbucket.org/jellynian/contrib/log"

	_ "github.com/go-sql-driver/mysql"
)

var onceForDefault sync.Once
var dbForLabchan *sql.DB

func Default() *sql.DB {
	onceForDefault.Do(func() {
		var conf = config.Default()
		var err error
		user := conf.Get("mysql.user").String()
		passwd := conf.Get("mysql.passwd").String()
		addr := conf.Get("mysql.addr").String()
		dbname := conf.Get("mysql.dbname").String()

		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", user, passwd, addr, dbname)
		log.Println(dataSourceName)
		dbForLabchan, err = sql.Open("mysql", dataSourceName)
		idleConn := conf.Get("mysql.idle").Int()

		dbForLabchan.SetMaxIdleConns(idleConn)
		maxConn := conf.Get("mysql.max").Int()
		dbForLabchan.SetMaxOpenConns(maxConn)

		err = dbForLabchan.Ping()
		if err != nil {
			log.Panic("connect mysql failed due to:", err.Error())
		}
	})
	return dbForLabchan
}
