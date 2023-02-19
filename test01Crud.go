package korm_example

import "fmt"

type test01Crud_D struct {
	Id   int
	Name string
}

func Test01Crud_Create(db *OrmAll) {
	db.test01Crud_D().Delete().MustRun()
	cnt := db.test01Crud_D().Select().MustRun_Count()
	assert(cnt == 0)

	db.test01Crud_D().MustInsert(test01Crud_D{
		Id:   10,
		Name: "10",
	})
	value := db.test01Crud_D().Select().Where_Id().Equal(10).MustRun_ResultOne()
	assert(value.Id == 10 && value.Name == "10")

	err := panicToError(func() {
		db.test01Crud_D().MustInsert(test01Crud_D{
			Id:   10,
			Name: "name",
		})
	})
	assert(err != nil)

	db.test01Crud_D().MustSet(test01Crud_D{
		Id:   12,
		Name: "12",
	})
	db.test01Crud_D().MustSet(test01Crud_D{
		Id:   12,
		Name: "12",
	})
	cnt = db.test01Crud_D().Select().MustRun_Count()
	assert(cnt == 2)

	value, ok := db.test01Crud_D().Select().Where_Id().Equal(12).MustRun_ResultOne2()
	assert(ok)
	assert(value.Id == 12 && value.Name == "12")
}

func Test01Crud_Read(db *OrmAll) {
	db.test01Crud_D().Delete().MustRun()

	cnt := db.test01Crud_D().Select().MustRun_Count()
	assert(cnt == 0)

	db.test01Crud_D().MustSet(test01Crud_D{
		Id:   10,
		Name: "name10",
	})
	db.test01Crud_D().MustSet(test01Crud_D{
		Id:   15,
		Name: "name15",
	})

	cnt = db.test01Crud_D().Select().Where_Id().GreaterOrEqual(10).MustRun_Count()
	assert(cnt == 2)

	cnt = db.test01Crud_D().Select().Where_Id().GreaterOrEqual(11).MustRun_Count()
	assert(cnt == 1)

	value, ok := db.test01Crud_D().Select().Where_Id().Equal(15).MustRun_ResultOne2()
	assert(ok && value.Id == 15)

	list := db.test01Crud_D().Select().OrderBy_Id().ASC().MustRun_ResultList()
	assert(len(list) == 2)
	assert(list[0].Id == 10)
	assert(list[1].Name == "name15")
}

func Test01Crud_Update(db *OrmAll) {
	db.test01Crud_D().Delete().MustRun()

	db.test01Crud_D().MustSet(test01Crud_D{
		Id:   10,
		Name: "10",
	})

	db.test01Crud_D().MustUpdateBy_Id(test01Crud_D{
		Id:   10,
		Name: "11",
	})
	value := db.test01Crud_D().Select().Where_Id().Equal(10).MustRun_ResultOne()
	assert(value.Id == 10 && value.Name == "11")

	_, ok := db.test01Crud_D().Select().Where_Id().Equal(11).MustRun_ResultOne2()
	assert(ok == false)

	af := db.test01Crud_D().Update().Where_Id().Equal(10).Set_Name("10").MustRun()
	assert(af == 1)
	value, ok = db.test01Crud_D().Select().Where_Name().Equal("10").MustRun_ResultOne2()
	assert(ok)
	assert(value.Id == 10 && value.Name == "10")

	cnt := db.test01Crud_D().Select().MustRun_Count()
	assert(cnt == 1)
}

func Test01Crud_Delete(db *OrmAll) {
	db.test01Crud_D().Delete().MustRun()

	af := db.test01Crud_D().MustUpdateBy_Id(test01Crud_D{
		Id:   10,
		Name: "10",
	})
	assert(af == 0)
	db.test01Crud_D().MustSet(test01Crud_D{
		Id:   11,
		Name: "11",
	})
	cnt := db.test01Crud_D().Select().MustRun_Count()
	assert(cnt == 1)

	af = db.test01Crud_D().Delete().MustRun()
	assert(af == 1)

	af = db.test01Crud_D().Delete().MustRun()
	assert(af == 0)
}

func assert(ok bool) {
	if ok == false {
		panic(`assert failed`)
	}
}

func panicToError(fn func()) (err error) {
	defer func() {
		v := recover()
		if v != nil {
			err = fmt.Errorf("panicToError %v", v)
		}
	}()
	fn()
	return err
}
