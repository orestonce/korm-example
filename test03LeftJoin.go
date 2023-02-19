package korm_example

import "time"

type test03User_D struct {
	Id           string
	Name         string
	PasswordHash string
	CreateTime   time.Time
}

type test03Group_D struct {
	Id         uint64
	Name       string
	CreateTime time.Time
}

type test03UserInGroup_D struct {
	UserId  string
	GroupId uint64
	User    *test03User_D  `korm:"join:this.UserId==other.Id"`
	Group   *test03Group_D `korm:"join:this.GroupId==other.Id"`
}

func Test03LeftJoin(db *OrmAll) {
	db.test03User_D().Delete().MustRun()
	db.test03Group_D().Delete().MustRun()
	db.test03UserInGroup_D().Delete().MustRun()

	db.test03User_D().MustSet(test03User_D{
		Id:         "uid",
		Name:       "uname",
		CreateTime: time.Now(),
	})
	db.test03Group_D().MustSet(test03Group_D{
		Id:         10,
		Name:       "gname",
		CreateTime: time.Now(),
	})
	db.test03UserInGroup_D().MustSet(test03UserInGroup_D{
		UserId:  "uid",
		GroupId: 10,
	})

	query := db.test03UserInGroup_D().Select()
	query.LeftJoin_User()
	query.LeftJoin_Group()
	value, ok := query.MustRun_ResultOne2()
	assert(ok)
	assert(value.UserId == "uid")
	assert(value.GroupId == 10)
	assert(value.User.Name == "uname")
	assert(value.Group.Name == "gname")
}
