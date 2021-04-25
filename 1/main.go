package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var Db *sql.DB

func init() {
	database, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
}

func queryData(sqlstr string, dest ...interface{}) (interface{}, error) {
	res, err := queryDataNotEmpty(sqlstr, dest)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, nil
	}
	return res, err
}

func queryDataNotEmpty(sql string, dest interface{}) (interface{}, error) {
	err := Db.QueryRow(sql).Scan(&dest)
	if err != nil {
		//包装err
		return nil, errors.Wrapf(err, sql+" query data is null")
	}
	return dest, nil
}

func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func main() {
	if Db == nil {
		os.Exit(1)
	}
	defer Db.Close()

	var name string
	var sqlstr = "select name from t_user where id=1"
	res, err := queryData(sqlstr, &name)
	if err != nil {
		fmt.Printf("FATAL: %v", err)
	} else {
		if res != nil {
			fmt.Println("queryData : ", B2S(res.([]uint8)))
		} else {
			fmt.Println("queryData : ", "data is empty")
		}
	}

	//数据为空时提示异常
	notEmptyRes, err2 := queryDataNotEmpty(sqlstr, name)
	if err2 != nil {
		fmt.Printf("FATAL: %+v", err2)
	} else {
		if res != nil {
			fmt.Println("queryData : ", B2S(notEmptyRes.([]uint8)))
		} else {
			fmt.Println("queryData : ", "data is empty")
		}
	}
}
