package korm_example

import (
	"time"
)

type test02MultiplePk_D struct {
	UserId     string `korm:"primary"`
	GroupId    string `korm:"primary"`
	CreateTime time.Time
}

func Test02MultiplePk(db *OrmAll) {
	db.test02MultiplePk_D().Delete().MustRun()
	db.test02MultiplePk_D().MustSet(test02MultiplePk_D{
		UserId:     "uid",
		GroupId:    "gid",
		CreateTime: time.Now(),
	})
	db.test02MultiplePk_D().MustSet(test02MultiplePk_D{
		UserId:     "uid",
		GroupId:    "gid2",
		CreateTime: time.Now(),
	})

	cnt := db.test02MultiplePk_D().Select().MustRun_Count()
	assert(cnt == 2)

	cnt = db.test02MultiplePk_D().Select().Where_GroupId().Equal("gid2").MustRun_Count()
	assert(cnt == 1)

	cnt = db.test02MultiplePk_D().Select().Where_GroupId().Equal("gid").MustRun_Count()
	assert(cnt == 1)

	value, ok := db.test02MultiplePk_D().Select().Where_GroupId().Equal("gid").MustRun_ResultOne2()
	assert(ok)
	assert(value.UserId == "uid" && value.GroupId == "gid")

	ra := db.test02MultiplePk_D().Delete().Where_GroupId().Equal("gid2").MustRun()
	assert(ra == 1)

	ok = db.test02MultiplePk_D().Select().MustRun_Exist()
	assert(ok)

	ra = db.test02MultiplePk_D().Delete().MustRun()
	assert(ra == 1)

	ok = db.test02MultiplePk_D().Select().MustRun_Exist()
	assert(ok == false)

	db.test02MultiplePk_D().MustSet(test02MultiplePk_D{
		UserId:     "uid",
		GroupId:    "gid",
		CreateTime: time.Now(),
	})
}
