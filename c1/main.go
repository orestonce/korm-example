package main

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/go-sql-driver/mysql"
	"github.com/orestonce/korm"
	korm_example "korm-example"
	"os"
)

func main() {
	for idx, db := range []*korm_example.OrmAll{
		initDbSqlite(),
		initTableMysql(),
	} {
		korm_example.Test01Crud_Create(db)
		korm_example.Test01Crud_Read(db)
		korm_example.Test01Crud_Update(db)
		korm_example.Test01Crud_Delete(db)
		korm_example.Test02MultiplePk(db)
		korm_example.Test03LeftJoin(db)
		korm_example.Test05View(db)
		fmt.Println("done", idx)
	}
}

func initDbSqlite() *korm_example.OrmAll {
	const f = "D:/1234.sqlite3"
	os.Remove(f)
	db, err := sql.Open("sqlite", f)
	if err != nil {
		panic(err)
	}
	korm_example.KORM_MustInitTableAll(db, korm.InitModeSqlite)
	return korm_example.OrmAllNew(db, korm.InitModeSqlite)
}

func initTableMysql() *korm_example.OrmAll {
	db := korm_example.KORM_MustNewDbMysql(korm_example.KORM_MustNewDbMysqlReq{
		Addr:        "127.0.0.1:3306",
		UserName:    "root",
		Password:    "",
		UseDatabase: "test",
	})
	korm_example.KORM_MustInitTableAll(db, korm.InitModeMysql)
	return korm_example.OrmAllNew(db, korm.InitModeMysql)
}
