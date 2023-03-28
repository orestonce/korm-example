package korm_example

import "time"

type test05User_D struct {
	Id   string
	Age  int16
	Name string
}

type test05UserAge_V struct {
	_   struct{} `korm:"view:test05User_D"` // 主表
	Id  string
	Age int16
}

type test05Group_D struct {
	Id   string
	Name string
}

type test05UserGroup_D struct {
	UserId     string         `korm:"primary"`
	GroupId    string         `korm:"primary"`
	User       *test05User_D  `korm:"join:this.UserId==other.Id"`
	Group      *test05Group_D `korm:"join:this.GroupId==other.Id"`
	CreateTime time.Time
}

type test05UserInGroup_V struct {
	_          struct{} `korm:"view:test05UserGroup_D"` // 主表
	UserId     string
	GroupId2   string `korm:"path:GroupId"`
	UserAge    uint16 `korm:"path:User.Age"`
	UserName   string `korm:"path:User.Name"`
	GroupName  string `korm:"path:Group.Name"`
	CreateTime time.Time
}

func Test05View(db *OrmAll) {
	db.test05User_D().MustSet(test05User_D{
		Id:   "uid1",
		Age:  10,
		Name: "uname1",
	})
	db.test05Group_D().MustSet(test05Group_D{
		Id:   "group1",
		Name: "群组1",
	})
	db.test05Group_D().MustSet(test05Group_D{
		Id:   "group2",
		Name: "群组1",
	})
	db.test05UserGroup_D().MustSet(test05UserGroup_D{
		UserId:     "uid1",
		GroupId:    "group1",
		CreateTime: time.Now(),
	})
	db.test05UserGroup_D().MustSet(test05UserGroup_D{
		UserId:     "uid1",
		GroupId:    "group2",
		CreateTime: time.Now(),
	})
	cnt := db.test05UserInGroup_V().Select().MustRun_Count()
	assert(cnt == 2)
	info, ok := db.test05UserInGroup_V().Select().Where_GroupId2().Equal("group1").MustRun_ResultOne2()
	assert(ok && info.UserId == "uid1" && info.UserName == "uname1" && info.UserAge == 10)

	list := db.test05UserInGroup_V().Select().OrderBy_GroupId2().ASC().MustRun_ResultList()
	assert(len(list) == 2 && list[0].GroupId2 == "group1")

	list, match := db.test05UserInGroup_V().Select().OrderBy_GroupId2().ASC().LimitOffset(1, 0).MustRun_ResultListWithTotalMatch()
	assert(len(list) == 1 && list[0].GroupId2 == "group1" && match == 2)

	uidGidMap := db.test05UserInGroup_V().Select().MustRun_ResultMap()
	assert(len(uidGidMap) == 1)
	info2 := uidGidMap["uid1"]
	assert(len(info2) == 2)
	_, ok1 := info2["group1"]
	_, ok2 := info2["group1"]
	assert(ok1 && ok2)

	ageMap := db.test05UserAge_V().Select().MustRun_ResultMap()
	assert(len(ageMap) == 1 && ageMap["uid1"].Age == 10)
}
