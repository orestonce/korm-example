package korm_example

import "time"

func Test06FnLength(db *OrmAll) {
	db.test05User_D().Delete().MustRun()
	db.test05Group_D().Delete().MustRun()
	db.test05UserGroup_D().Delete().MustRun()

	db.test05User_D().MustSet(test05User_D{
		Id:   "uid1",
		Age:  10,
		Name: "uname1",
	})
	db.test05Group_D().MustSet(test05Group_D{
		Id:   "group1",
		Name: "g1",
	})
	db.test05Group_D().MustSet(test05Group_D{
		Id:   "group2",
		Name: "",
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

	list := db.test05User_D().Select().Where_Name().Length().Equal(6).MustRun_ResultList()
	if len(list) != 1 {
		panic(len(list))
	}

	list2 := db.test05Group_D().Select().Where_Name().Length().Equal(0).MustRun_ResultList()
	if len(list2) != 1 {
		panic(len(list2))
	}

	list2 = db.test05Group_D().Select().Where_Name().Length().Equal(1).MustRun_ResultList()
	if len(list2) != 0 {
		panic(len(list2))
	}

	count := db.test05UserInGroup_V().Select().Where_GroupName().Length().Equal(0).MustRun_Count()
	if count != 1 {
		panic(count)
	}
}
