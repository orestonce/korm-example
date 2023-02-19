package korm_example

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/orestonce/korm"
	"strconv"
	"strings"
	"time"
)

type OrmAll struct {
	db   *sql.DB // db, tx任选其一
	tx   *sql.Tx
	mode string // sqlite, mysql
}

func (this *OrmAll) ExecRawQuery(query string, args ...interface{}) (*sql.Rows, error) {
	if this.db != nil {
		return this.db.Query(query, args...)
	} else if this.tx != nil {
		return this.tx.Query(query, args...)
	}
	panic("ExecRawQuery: OrmAll must include db or tx")
}

func OrmAllNew(db *sql.DB, mode string) *OrmAll {
	return &OrmAll{
		db:   db,
		mode: mode,
	}
}

func (this *OrmAll) ExecRaw(query string, args ...interface{}) (sql.Result, error) {
	if this.db != nil {
		return this.db.Exec(query, args...)
	} else if this.tx != nil {
		return this.tx.Exec(query, args...)
	}
	panic("ExecRaw: OrmAll must include db or tx")
}

func (this *OrmAll) MustTxCallback(cb func(db *OrmAll)) {
	if this.tx != nil {
		panic("MustSingleThreadTxCallback repeat call")
	}
	tx, err := this.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	cb(&OrmAll{
		tx:   tx,
		mode: this.mode,
	})
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

type KORM_MustNewDbMysqlReq struct {
	Addr        string
	UserName    string
	Password    string
	UseDatabase string
}

func KORM_MustNewDbMysql(req KORM_MustNewDbMysqlReq) (db *sql.DB) {
	var err error

	db, err = sql.Open("mysql", req.UserName+":"+req.Password+"@tcp("+req.Addr+")/")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + req.UseDatabase)
	if err != nil {
		panic(err)
	}

	_ = db.Close()
	db, err = sql.Open("mysql", req.UserName+":"+req.Password+"@tcp("+req.Addr+")/"+req.UseDatabase+"?charset=utf8")
	if err != nil {
		panic(err)
	}

	return db
}
func KORM_MustInitTableAll(db *sql.DB, mode string) {
	var err error
	err = korm.InitTable(korm.InitTableReq{
		Mode:      mode,
		TableName: "DownloadCache_D",
		FieldList: []korm.FieldSqlDefine{
			{
				Type:         korm.SqlTypeChar255,
				Name:         "Url",
				IsPrimaryKey: true,
			}, {
				Type:         korm.SqlTypeLongBlob,
				Name:         "Content",
				IsPrimaryKey: false,
			},
		},
		Db: db,
	})
	if err != nil {
		panic(err)
	}

	err = korm.InitTable(korm.InitTableReq{
		Mode:      mode,
		TableName: "test01Crud_D",
		FieldList: []korm.FieldSqlDefine{
			{
				Type:         korm.SqlTypeBigInt,
				Name:         "Id",
				IsPrimaryKey: true,
			}, {
				Type:         korm.SqlTypeLongBlob,
				Name:         "Name",
				IsPrimaryKey: false,
			},
		},
		Db: db,
	})
	if err != nil {
		panic(err)
	}

	err = korm.InitTable(korm.InitTableReq{
		Mode:      mode,
		TableName: "test02MultiplePk_D",
		FieldList: []korm.FieldSqlDefine{
			{
				Type:         korm.SqlTypeChar255,
				Name:         "UserId",
				IsPrimaryKey: true,
			}, {
				Type:         korm.SqlTypeChar255,
				Name:         "GroupId",
				IsPrimaryKey: true,
			}, {
				Type:         korm.SqlTypeChar255,
				Name:         "CreateTime",
				IsPrimaryKey: false,
			},
		},
		Db: db,
	})
	if err != nil {
		panic(err)
	}

	err = korm.InitTable(korm.InitTableReq{
		Mode:      mode,
		TableName: "test03User_D",
		FieldList: []korm.FieldSqlDefine{
			{
				Type:         korm.SqlTypeChar255,
				Name:         "Id",
				IsPrimaryKey: true,
			}, {
				Type:         korm.SqlTypeLongBlob,
				Name:         "Name",
				IsPrimaryKey: false,
			}, {
				Type:         korm.SqlTypeLongBlob,
				Name:         "PasswordHash",
				IsPrimaryKey: false,
			}, {
				Type:         korm.SqlTypeChar255,
				Name:         "CreateTime",
				IsPrimaryKey: false,
			},
		},
		Db: db,
	})
	if err != nil {
		panic(err)
	}

	err = korm.InitTable(korm.InitTableReq{
		Mode:      mode,
		TableName: "test03Group_D",
		FieldList: []korm.FieldSqlDefine{
			{
				Type:         korm.SqlTypeBigInt,
				Name:         "Id",
				IsPrimaryKey: true,
			}, {
				Type:         korm.SqlTypeLongBlob,
				Name:         "Name",
				IsPrimaryKey: false,
			}, {
				Type:         korm.SqlTypeChar255,
				Name:         "CreateTime",
				IsPrimaryKey: false,
			},
		},
		Db: db,
	})
	if err != nil {
		panic(err)
	}

	err = korm.InitTable(korm.InitTableReq{
		Mode:      mode,
		TableName: "test03UserInGroup_D",
		FieldList: []korm.FieldSqlDefine{
			{
				Type:         korm.SqlTypeChar255,
				Name:         "UserId",
				IsPrimaryKey: true,
			}, {
				Type:         korm.SqlTypeBigInt,
				Name:         "GroupId",
				IsPrimaryKey: false,
			},
		},
		Db: db,
	})
	if err != nil {
		panic(err)
	}

	err = korm.InitTable(korm.InitTableReq{
		Mode:      mode,
		TableName: "test04User_D",
		FieldList: []korm.FieldSqlDefine{
			{
				Type:         korm.SqlTypeBigInt,
				Name:         "Id",
				IsPrimaryKey: true,
			}, {
				Type:         korm.SqlTypeLongBlob,
				Name:         "Name",
				IsPrimaryKey: false,
			}, {
				Type:         korm.SqlTypeLongBlob,
				Name:         "Key",
				IsPrimaryKey: false,
			},
		},
		IndexList: [][]string{
			{"Name"},
			{"Id", "Name"},
			{"Key"},
		},
		Db: db,
	})
	if err != nil {
		panic(err)
	}

}

type korm_scan_resp struct {
	argList    []interface{}
	argParseFn []func(v string) (err error)
}

func korm_DownloadCache_D_scan(joinNode *korm_leftJoinNode, info *DownloadCache_D) (resp korm_scan_resp) {
	for _, one := range joinNode.selectFieldNameList {
		switch one {
		default:
			panic("GetScanInfoCode error" + strconv.Quote(one))
		case "Url":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.Url = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "Content":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.Content = []byte(v)

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		}
	}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("DownloadCache_D")
		}
	}
	return resp
}
func korm_test01Crud_D_scan(joinNode *korm_leftJoinNode, info *test01Crud_D) (resp korm_scan_resp) {
	for _, one := range joinNode.selectFieldNameList {
		switch one {
		default:
			panic("GetScanInfoCode error" + strconv.Quote(one))
		case "Id":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				vi, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				info.Id = int(vi)

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "Name":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.Name = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		}
	}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("test01Crud_D")
		}
	}
	return resp
}
func korm_test02MultiplePk_D_scan(joinNode *korm_leftJoinNode, info *test02MultiplePk_D) (resp korm_scan_resp) {
	for _, one := range joinNode.selectFieldNameList {
		switch one {
		default:
			panic("GetScanInfoCode error" + strconv.Quote(one))
		case "UserId":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.UserId = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "GroupId":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.GroupId = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "CreateTime":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				vt, err := time.Parse(time.RFC3339Nano, v)
				if err != nil {
					return err
				}
				info.CreateTime = vt

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		}
	}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("test02MultiplePk_D")
		}
	}
	return resp
}
func korm_test03User_D_scan(joinNode *korm_leftJoinNode, info *test03User_D) (resp korm_scan_resp) {
	for _, one := range joinNode.selectFieldNameList {
		switch one {
		default:
			panic("GetScanInfoCode error" + strconv.Quote(one))
		case "Id":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.Id = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "Name":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.Name = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "PasswordHash":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.PasswordHash = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "CreateTime":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				vt, err := time.Parse(time.RFC3339Nano, v)
				if err != nil {
					return err
				}
				info.CreateTime = vt

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		}
	}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("test03User_D")
		}
	}
	return resp
}
func korm_test03Group_D_scan(joinNode *korm_leftJoinNode, info *test03Group_D) (resp korm_scan_resp) {
	for _, one := range joinNode.selectFieldNameList {
		switch one {
		default:
			panic("GetScanInfoCode error" + strconv.Quote(one))
		case "Id":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				vu, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return err
				}
				info.Id = uint64(vu)

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "Name":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.Name = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "CreateTime":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				vt, err := time.Parse(time.RFC3339Nano, v)
				if err != nil {
					return err
				}
				info.CreateTime = vt

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		}
	}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("test03Group_D")
		}
	}
	return resp
}
func korm_test03UserInGroup_D_scan(joinNode *korm_leftJoinNode, info *test03UserInGroup_D) (resp korm_scan_resp) {
	for _, one := range joinNode.selectFieldNameList {
		switch one {
		default:
			panic("GetScanInfoCode error" + strconv.Quote(one))
		case "UserId":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.UserId = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "GroupId":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				vu, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return err
				}
				info.GroupId = uint64(vu)

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		}
	}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("test03UserInGroup_D")
		case "User":
			info.User = &test03User_D{}
			resp0 := korm_test03User_D_scan(sub, info.User)
			resp.argList = append(resp.argList, resp0.argList...)
			resp.argParseFn = append(resp.argParseFn, resp0.argParseFn...)
		case "Group":
			info.Group = &test03Group_D{}
			resp0 := korm_test03Group_D_scan(sub, info.Group)
			resp.argList = append(resp.argList, resp0.argList...)
			resp.argParseFn = append(resp.argParseFn, resp0.argParseFn...)
		}
	}
	return resp
}
func korm_test04User_D_scan(joinNode *korm_leftJoinNode, info *test04User_D) (resp korm_scan_resp) {
	for _, one := range joinNode.selectFieldNameList {
		switch one {
		default:
			panic("GetScanInfoCode error" + strconv.Quote(one))
		case "Id":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				vi, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				info.Id = int(vi)

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "Name":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.Name = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		case "Key":
			resp.argList = append(resp.argList, new(sql.NullString))
			tmpFn := func(v string) (err error) {
				info.Key = v

				return nil
			}
			resp.argParseFn = append(resp.argParseFn, tmpFn)
		}
	}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("test04User_D")
		}
	}
	return resp
}

type korm_leftJoinRootInfo struct {
	bufLeftJoin  bytes.Buffer
	tableNameIdx int
}

type korm_leftJoinNode struct {
	thisLevelJoinList   []*korm_leftJoinNode
	fieldName           string
	tableName           string // _0, _1, _2
	root                *korm_leftJoinRootInfo
	selectFieldNameList []string
}

func (this *korm_leftJoinNode) writeToBuf(buf2 *bytes.Buffer) {
	for idx, fieldName := range this.selectFieldNameList {
		if idx > 0 {
			buf2.WriteString(", ")
		}
		buf2.WriteString(this.tableName + ".`" + fieldName + "`")
	}
	buf2.WriteString(" ")
}
func (this *korm_leftJoinNode) fillSelect(buf2 *bytes.Buffer, isFirst bool) {
	for idx, fieldName := range this.selectFieldNameList {
		if idx > 0 || isFirst == false {
			buf2.WriteString(", ")
		}
		buf2.WriteString(this.tableName + ".`" + fieldName + "`")
	}
	for _, sub := range this.thisLevelJoinList {
		sub.fillSelect(buf2, false)
	}
}
func (this *korm_leftJoinNode) addLeftJoin(rightTableName string, thisFieldName string, leftFieldName string, rightFieldName string) *korm_leftJoinNode {
	for _, one := range this.thisLevelJoinList {
		if one.fieldName == thisFieldName {
			return one
		}
	}
	newNode := &korm_leftJoinNode{
		fieldName: thisFieldName,
		tableName: "_" + strconv.Itoa(this.root.tableNameIdx),
		root:      this.root,
	}
	this.root.tableNameIdx++
	this.thisLevelJoinList = append(this.thisLevelJoinList, newNode)
	this.root.bufLeftJoin.WriteString("LEFT JOIN `" + rightTableName + "` " + newNode.tableName + " ON " + this.tableName + ".`" + leftFieldName + "` = " + newNode.tableName + ".`" + rightFieldName + "` ")
	return newNode
}

type KORM_DownloadCache_D struct {
	supper *OrmAll
}

func (this *OrmAll) DownloadCache_D() *KORM_DownloadCache_D {
	return &KORM_DownloadCache_D{supper: this}
}
func korm_fillSelectFieldNameList_DownloadCache_D(joinNode *korm_leftJoinNode) {
	joinNode.selectFieldNameList = []string{"Url", "Content"}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("korm_fillSelectFieldNameList_DownloadCache_D" + strconv.Quote(sub.fieldName))
		}
	}
}
func (this *KORM_DownloadCache_D) MustInsert(info DownloadCache_D) {
	var err error
	_, err = this.supper.ExecRaw("INSERT INTO `DownloadCache_D`(`Url` ,`Content` ) VALUES(?,?)", info.Url, string(info.Content))
	if err != nil {
		panic(err)
	}

	return
}
func (this *KORM_DownloadCache_D) MustSet(info DownloadCache_D) {
	var err error
	_, err = this.supper.ExecRaw("REPLACE INTO `DownloadCache_D`(`Url` ,`Content` ) VALUES(?,?)", info.Url, string(info.Content))
	if err != nil {
		panic(err)
	}

	return
}

// Select DownloadCache_D
type KORM_DownloadCache_D_SelectObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
	joinNode    *korm_leftJoinNode
	rootInfo    korm_leftJoinRootInfo
}

func (this *KORM_DownloadCache_D) Select() *KORM_DownloadCache_D_SelectObj {
	one := &KORM_DownloadCache_D_SelectObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
	one.joinNode = &korm_leftJoinNode{tableName: "_0"}
	one.joinNode.root = &one.rootInfo
	one.rootInfo.tableNameIdx = 1
	return one
}

type KORM_DownloadCache_D_SelectObj_OrderByObj struct {
	fieldName string
	supper    *KORM_DownloadCache_D_SelectObj
}

func (this *KORM_DownloadCache_D_SelectObj_OrderByObj) ASC() *KORM_DownloadCache_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_DownloadCache_D_SelectObj_OrderByObj) DESC() *KORM_DownloadCache_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_DownloadCache_D_SelectObj) OrderBy_Url() *KORM_DownloadCache_D_SelectObj_OrderByObj {
	return &KORM_DownloadCache_D_SelectObj_OrderByObj{
		fieldName: "Url",
		supper:    this,
	}
}

func (this *KORM_DownloadCache_D_SelectObj) LimitOffset(limit int, offset int) *KORM_DownloadCache_D_SelectObj {
	this.limit = limit
	this.offset = offset
	return this
}

// pageSize: [1, n)
// pageNo:   [1,n)
func (this *KORM_DownloadCache_D_SelectObj) SetPageLimit(pageSize int, pageNo int) *KORM_DownloadCache_D_SelectObj {
	if pageSize <= 0 || pageNo <= 0 {
		panic("KORM_DownloadCache_D_SelectObj SetPageLimit error param")
	}
	this.limit = pageSize
	this.offset = pageSize * (pageNo - 1)
	return this
}

func (this *KORM_DownloadCache_D_SelectObj) MustRun_Count() (cnt int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `DownloadCache_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_DownloadCache_D_SelectObj not found.")
	}
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt
}

func (this *KORM_DownloadCache_D_SelectObj) MustRun_Exist() (exist bool) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `DownloadCache_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	buf2.WriteString(" LIMIT 1 ")

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_DownloadCache_D_SelectObj not found.")
	}
	var cnt int64
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func (this *KORM_DownloadCache_D_SelectObj) MustRun_ResultOne2() (info DownloadCache_D, ok bool) {
	this.limit = 1
	korm_fillSelectFieldNameList_DownloadCache_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `DownloadCache_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		return info, false
	}
	resp := korm_DownloadCache_D_scan(this.joinNode, &info)
	err = result.Scan(resp.argList...)
	if err != nil {
		panic(err)
	}

	for idx, a := range resp.argList {
		v := a.(*sql.NullString).String
		if v == "" {
			continue
		}
		err = resp.argParseFn[idx](v)
		if err != nil {
			panic(err)
		}

	}

	return info, true
}

func (this *KORM_DownloadCache_D_SelectObj) MustRun_ResultOne() (info DownloadCache_D) {
	info, _ = this.MustRun_ResultOne2()
	return info
}
func (this *KORM_DownloadCache_D_SelectObj) MustRun_ResultList() (list []DownloadCache_D) {
	korm_fillSelectFieldNameList_DownloadCache_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `DownloadCache_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info DownloadCache_D
		korm_fillSelectFieldNameList_DownloadCache_D(this.joinNode)
		resp := korm_DownloadCache_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	return list
}
func (this *KORM_DownloadCache_D_SelectObj) MustRun_ResultMap() (m map[string]DownloadCache_D) {
	korm_fillSelectFieldNameList_DownloadCache_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `DownloadCache_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info DownloadCache_D
		korm_fillSelectFieldNameList_DownloadCache_D(this.joinNode)
		resp := korm_DownloadCache_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		if m == nil {
			m = map[string]DownloadCache_D{}
		}
		m[info.Url] = info

	}
	return m
}
func (this *KORM_DownloadCache_D_SelectObj) MustRun_ResultListWithTotalMatch() (list []DownloadCache_D, totalMatch int64) {
	var err error
	korm_fillSelectFieldNameList_DownloadCache_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	if this.supper.mode == korm.InitModeMysql {
		buf2.WriteString("SQL_CALC_FOUND_ROWS  ")
	}
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `DownloadCache_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}
	var conn *sql.Conn
	var result *sql.Rows
	if this.supper.db != nil {
		var err error
		conn, err = this.supper.db.Conn(context.Background())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		result, err = conn.QueryContext(context.Background(), buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		result, err = this.supper.tx.Query(buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	}

	defer result.Close()
	for result.Next() {
		var info DownloadCache_D
		korm_fillSelectFieldNameList_DownloadCache_D(this.joinNode)
		resp := korm_DownloadCache_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	result.Close()
	if this.limit > 0 {
		nextQuery := ""
		if this.supper.mode == korm.InitModeMysql {
			nextQuery = "select FOUND_ROWS()"
		} else if this.supper.mode == korm.InitModeSqlite {
			buf2.Reset()
			buf2.WriteString("SELECT COUNT(1) ")
			buf2.WriteString("FROM `DownloadCache_D` " + this.joinNode.tableName + " ")
			buf2.WriteString(this.rootInfo.bufLeftJoin.String())
			buf2.WriteString(this.bufWhere.String())
			nextQuery = buf2.String()
		} else {
			panic("not support")
		}
		var result2 *sql.Rows
		if conn != nil {
			result2, err = conn.QueryContext(context.Background(), nextQuery)
		} else {
			result2, err = this.supper.tx.Query(nextQuery)
		}
		if err != nil {
			panic(err)
		}

		defer result2.Close()

		if result2.Next() == false {
			panic("MustRun_ResultListWithPageInfo ")
		}
		err = result2.Scan(&totalMatch)
		if err != nil {
			panic(err)
		}

	}
	return list, totalMatch
}

type KORM_Where_KORM_DownloadCache_D_SelectObj_Url struct {
	supper *KORM_DownloadCache_D_SelectObj
}

func (this *KORM_DownloadCache_D_SelectObj) Where_Url() *KORM_Where_KORM_DownloadCache_D_SelectObj_Url {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Url` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_DownloadCache_D_SelectObj_Url{supper: this}
}
func (this *KORM_Where_KORM_DownloadCache_D_SelectObj_Url) Equal(Url string) *KORM_DownloadCache_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_SelectObj_Url) NotEqual(Url string) *KORM_DownloadCache_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_SelectObj_Url) GreaterOrEqual(Url string) *KORM_DownloadCache_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_SelectObj_Url) Less(Url string) *KORM_DownloadCache_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_SelectObj_Url) LessOrEqual(Url string) *KORM_DownloadCache_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_SelectObj_Url) In(vList []string) *KORM_DownloadCache_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_DownloadCache_D_SelectObj) CondMultOpBegin_AND() *KORM_DownloadCache_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_DownloadCache_D_SelectObj) CondMultOpBegin_OR() *KORM_DownloadCache_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_DownloadCache_D_SelectObj) CondMultOpEnd() *KORM_DownloadCache_D_SelectObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

// Update DownloadCache_D
type KORM_DownloadCache_D_UpdateObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_DownloadCache_D) Update() *KORM_DownloadCache_D_UpdateObj {
	return &KORM_DownloadCache_D_UpdateObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_DownloadCache_D) MustUpdateBy_Url(info DownloadCache_D) (rowsAffected int64) {
	rowsAffected = this.Update().Where_Url().Equal(info.Url).Set_Content(info.Content).MustRun()
	return rowsAffected
}
func (this *KORM_DownloadCache_D_UpdateObj) Set_Content(Content []byte) *KORM_DownloadCache_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`Content` = ? ")
	this.argsSet = append(this.argsSet, string(Content))
	return this
}
func (this *KORM_DownloadCache_D_UpdateObj) MustRun() (RowsAffected int64) {
	if len(this.argsSet) == 0 {
		panic("len(this.argsSet) == 0")
	}
	var buf2 bytes.Buffer
	buf2.WriteString("UPDATE `DownloadCache_D` ")
	buf2.WriteString(this.bufSet.String())
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error

	result, err = this.supper.ExecRaw(buf2.String(), append(this.argsSet, this.argsWhere...)...)
	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_DownloadCache_D_UpdateObj_Url struct {
	supper *KORM_DownloadCache_D_UpdateObj
}

func (this *KORM_DownloadCache_D_UpdateObj) Where_Url() *KORM_Where_KORM_DownloadCache_D_UpdateObj_Url {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Url` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_DownloadCache_D_UpdateObj_Url{supper: this}
}
func (this *KORM_Where_KORM_DownloadCache_D_UpdateObj_Url) Equal(Url string) *KORM_DownloadCache_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_UpdateObj_Url) NotEqual(Url string) *KORM_DownloadCache_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_UpdateObj_Url) GreaterOrEqual(Url string) *KORM_DownloadCache_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_UpdateObj_Url) Less(Url string) *KORM_DownloadCache_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_UpdateObj_Url) LessOrEqual(Url string) *KORM_DownloadCache_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_UpdateObj_Url) In(vList []string) *KORM_DownloadCache_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_DownloadCache_D_UpdateObj) CondMultOpBegin_AND() *KORM_DownloadCache_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_DownloadCache_D_UpdateObj) CondMultOpBegin_OR() *KORM_DownloadCache_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_DownloadCache_D_UpdateObj) CondMultOpEnd() *KORM_DownloadCache_D_UpdateObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_DownloadCache_D_DeleteObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_DownloadCache_D) Delete() *KORM_DownloadCache_D_DeleteObj {
	return &KORM_DownloadCache_D_DeleteObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_DownloadCache_D_DeleteObj) MustRun() (RowsAffected int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("DELETE FROM DownloadCache_D ")
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error
	result, err = this.supper.ExecRaw(buf2.String(), this.argsWhere...)

	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_DownloadCache_D_DeleteObj_Url struct {
	supper *KORM_DownloadCache_D_DeleteObj
}

func (this *KORM_DownloadCache_D_DeleteObj) Where_Url() *KORM_Where_KORM_DownloadCache_D_DeleteObj_Url {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Url` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_DownloadCache_D_DeleteObj_Url{supper: this}
}
func (this *KORM_Where_KORM_DownloadCache_D_DeleteObj_Url) Equal(Url string) *KORM_DownloadCache_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_DeleteObj_Url) NotEqual(Url string) *KORM_DownloadCache_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_DeleteObj_Url) GreaterOrEqual(Url string) *KORM_DownloadCache_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_DeleteObj_Url) Less(Url string) *KORM_DownloadCache_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_DeleteObj_Url) LessOrEqual(Url string) *KORM_DownloadCache_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Url)
	return this.supper
}
func (this *KORM_Where_KORM_DownloadCache_D_DeleteObj_Url) In(vList []string) *KORM_DownloadCache_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_DownloadCache_D_DeleteObj) CondMultOpBegin_AND() *KORM_DownloadCache_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_DownloadCache_D_DeleteObj) CondMultOpBegin_OR() *KORM_DownloadCache_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_DownloadCache_D_DeleteObj) CondMultOpEnd() *KORM_DownloadCache_D_DeleteObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test01Crud_D struct {
	supper *OrmAll
}

func (this *OrmAll) test01Crud_D() *KORM_test01Crud_D {
	return &KORM_test01Crud_D{supper: this}
}
func korm_fillSelectFieldNameList_test01Crud_D(joinNode *korm_leftJoinNode) {
	joinNode.selectFieldNameList = []string{"Id", "Name"}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("korm_fillSelectFieldNameList_test01Crud_D" + strconv.Quote(sub.fieldName))
		}
	}
}
func (this *KORM_test01Crud_D) MustInsert(info test01Crud_D) {
	var err error
	_, err = this.supper.ExecRaw("INSERT INTO `test01Crud_D`(`Id` ,`Name` ) VALUES(?,?)", info.Id, info.Name)
	if err != nil {
		panic(err)
	}

	return
}
func (this *KORM_test01Crud_D) MustSet(info test01Crud_D) {
	var err error
	_, err = this.supper.ExecRaw("REPLACE INTO `test01Crud_D`(`Id` ,`Name` ) VALUES(?,?)", info.Id, info.Name)
	if err != nil {
		panic(err)
	}

	return
}

// Select test01Crud_D
type KORM_test01Crud_D_SelectObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
	joinNode    *korm_leftJoinNode
	rootInfo    korm_leftJoinRootInfo
}

func (this *KORM_test01Crud_D) Select() *KORM_test01Crud_D_SelectObj {
	one := &KORM_test01Crud_D_SelectObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
	one.joinNode = &korm_leftJoinNode{tableName: "_0"}
	one.joinNode.root = &one.rootInfo
	one.rootInfo.tableNameIdx = 1
	return one
}

type KORM_test01Crud_D_SelectObj_OrderByObj struct {
	fieldName string
	supper    *KORM_test01Crud_D_SelectObj
}

func (this *KORM_test01Crud_D_SelectObj_OrderByObj) ASC() *KORM_test01Crud_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_test01Crud_D_SelectObj_OrderByObj) DESC() *KORM_test01Crud_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_test01Crud_D_SelectObj) OrderBy_Id() *KORM_test01Crud_D_SelectObj_OrderByObj {
	return &KORM_test01Crud_D_SelectObj_OrderByObj{
		fieldName: "Id",
		supper:    this,
	}
}
func (this *KORM_test01Crud_D_SelectObj) OrderBy_Name() *KORM_test01Crud_D_SelectObj_OrderByObj {
	return &KORM_test01Crud_D_SelectObj_OrderByObj{
		fieldName: "Name",
		supper:    this,
	}
}

func (this *KORM_test01Crud_D_SelectObj) LimitOffset(limit int, offset int) *KORM_test01Crud_D_SelectObj {
	this.limit = limit
	this.offset = offset
	return this
}

// pageSize: [1, n)
// pageNo:   [1,n)
func (this *KORM_test01Crud_D_SelectObj) SetPageLimit(pageSize int, pageNo int) *KORM_test01Crud_D_SelectObj {
	if pageSize <= 0 || pageNo <= 0 {
		panic("KORM_test01Crud_D_SelectObj SetPageLimit error param")
	}
	this.limit = pageSize
	this.offset = pageSize * (pageNo - 1)
	return this
}

func (this *KORM_test01Crud_D_SelectObj) MustRun_Count() (cnt int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test01Crud_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test01Crud_D_SelectObj not found.")
	}
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt
}

func (this *KORM_test01Crud_D_SelectObj) MustRun_Exist() (exist bool) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test01Crud_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	buf2.WriteString(" LIMIT 1 ")

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test01Crud_D_SelectObj not found.")
	}
	var cnt int64
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func (this *KORM_test01Crud_D_SelectObj) MustRun_ResultOne2() (info test01Crud_D, ok bool) {
	this.limit = 1
	korm_fillSelectFieldNameList_test01Crud_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test01Crud_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		return info, false
	}
	resp := korm_test01Crud_D_scan(this.joinNode, &info)
	err = result.Scan(resp.argList...)
	if err != nil {
		panic(err)
	}

	for idx, a := range resp.argList {
		v := a.(*sql.NullString).String
		if v == "" {
			continue
		}
		err = resp.argParseFn[idx](v)
		if err != nil {
			panic(err)
		}

	}

	return info, true
}

func (this *KORM_test01Crud_D_SelectObj) MustRun_ResultOne() (info test01Crud_D) {
	info, _ = this.MustRun_ResultOne2()
	return info
}
func (this *KORM_test01Crud_D_SelectObj) MustRun_ResultList() (list []test01Crud_D) {
	korm_fillSelectFieldNameList_test01Crud_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test01Crud_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test01Crud_D
		korm_fillSelectFieldNameList_test01Crud_D(this.joinNode)
		resp := korm_test01Crud_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	return list
}
func (this *KORM_test01Crud_D_SelectObj) MustRun_ResultMap() (m map[int]test01Crud_D) {
	korm_fillSelectFieldNameList_test01Crud_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test01Crud_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test01Crud_D
		korm_fillSelectFieldNameList_test01Crud_D(this.joinNode)
		resp := korm_test01Crud_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		if m == nil {
			m = map[int]test01Crud_D{}
		}
		m[info.Id] = info

	}
	return m
}
func (this *KORM_test01Crud_D_SelectObj) MustRun_ResultListWithTotalMatch() (list []test01Crud_D, totalMatch int64) {
	var err error
	korm_fillSelectFieldNameList_test01Crud_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	if this.supper.mode == korm.InitModeMysql {
		buf2.WriteString("SQL_CALC_FOUND_ROWS  ")
	}
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test01Crud_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}
	var conn *sql.Conn
	var result *sql.Rows
	if this.supper.db != nil {
		var err error
		conn, err = this.supper.db.Conn(context.Background())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		result, err = conn.QueryContext(context.Background(), buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		result, err = this.supper.tx.Query(buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	}

	defer result.Close()
	for result.Next() {
		var info test01Crud_D
		korm_fillSelectFieldNameList_test01Crud_D(this.joinNode)
		resp := korm_test01Crud_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	result.Close()
	if this.limit > 0 {
		nextQuery := ""
		if this.supper.mode == korm.InitModeMysql {
			nextQuery = "select FOUND_ROWS()"
		} else if this.supper.mode == korm.InitModeSqlite {
			buf2.Reset()
			buf2.WriteString("SELECT COUNT(1) ")
			buf2.WriteString("FROM `test01Crud_D` " + this.joinNode.tableName + " ")
			buf2.WriteString(this.rootInfo.bufLeftJoin.String())
			buf2.WriteString(this.bufWhere.String())
			nextQuery = buf2.String()
		} else {
			panic("not support")
		}
		var result2 *sql.Rows
		if conn != nil {
			result2, err = conn.QueryContext(context.Background(), nextQuery)
		} else {
			result2, err = this.supper.tx.Query(nextQuery)
		}
		if err != nil {
			panic(err)
		}

		defer result2.Close()

		if result2.Next() == false {
			panic("MustRun_ResultListWithPageInfo ")
		}
		err = result2.Scan(&totalMatch)
		if err != nil {
			panic(err)
		}

	}
	return list, totalMatch
}

type KORM_Where_KORM_test01Crud_D_SelectObj_Id struct {
	supper *KORM_test01Crud_D_SelectObj
}

func (this *KORM_test01Crud_D_SelectObj) Where_Id() *KORM_Where_KORM_test01Crud_D_SelectObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test01Crud_D_SelectObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Id) Equal(Id int) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Id) NotEqual(Id int) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Id) GreaterOrEqual(Id int) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Id) Less(Id int) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Id) LessOrEqual(Id int) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Id) In(vList []int) *KORM_test01Crud_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test01Crud_D_SelectObj_Name struct {
	supper *KORM_test01Crud_D_SelectObj
}

func (this *KORM_test01Crud_D_SelectObj) Where_Name() *KORM_Where_KORM_test01Crud_D_SelectObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test01Crud_D_SelectObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Name) Equal(Name string) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Name) NotEqual(Name string) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Name) GreaterOrEqual(Name string) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Name) Less(Name string) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Name) LessOrEqual(Name string) *KORM_test01Crud_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_SelectObj_Name) In(vList []string) *KORM_test01Crud_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test01Crud_D_SelectObj) CondMultOpBegin_AND() *KORM_test01Crud_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test01Crud_D_SelectObj) CondMultOpBegin_OR() *KORM_test01Crud_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test01Crud_D_SelectObj) CondMultOpEnd() *KORM_test01Crud_D_SelectObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

// Update test01Crud_D
type KORM_test01Crud_D_UpdateObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test01Crud_D) Update() *KORM_test01Crud_D_UpdateObj {
	return &KORM_test01Crud_D_UpdateObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test01Crud_D) MustUpdateBy_Id(info test01Crud_D) (rowsAffected int64) {
	rowsAffected = this.Update().Where_Id().Equal(info.Id).Set_Name(info.Name).MustRun()
	return rowsAffected
}
func (this *KORM_test01Crud_D_UpdateObj) Set_Name(Name string) *KORM_test01Crud_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`Name` = ? ")
	this.argsSet = append(this.argsSet, Name)
	return this
}
func (this *KORM_test01Crud_D_UpdateObj) MustRun() (RowsAffected int64) {
	if len(this.argsSet) == 0 {
		panic("len(this.argsSet) == 0")
	}
	var buf2 bytes.Buffer
	buf2.WriteString("UPDATE `test01Crud_D` ")
	buf2.WriteString(this.bufSet.String())
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error

	result, err = this.supper.ExecRaw(buf2.String(), append(this.argsSet, this.argsWhere...)...)
	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test01Crud_D_UpdateObj_Id struct {
	supper *KORM_test01Crud_D_UpdateObj
}

func (this *KORM_test01Crud_D_UpdateObj) Where_Id() *KORM_Where_KORM_test01Crud_D_UpdateObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test01Crud_D_UpdateObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Id) Equal(Id int) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Id) NotEqual(Id int) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Id) GreaterOrEqual(Id int) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Id) Less(Id int) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Id) LessOrEqual(Id int) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Id) In(vList []int) *KORM_test01Crud_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test01Crud_D_UpdateObj_Name struct {
	supper *KORM_test01Crud_D_UpdateObj
}

func (this *KORM_test01Crud_D_UpdateObj) Where_Name() *KORM_Where_KORM_test01Crud_D_UpdateObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test01Crud_D_UpdateObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Name) Equal(Name string) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Name) NotEqual(Name string) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Name) GreaterOrEqual(Name string) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Name) Less(Name string) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Name) LessOrEqual(Name string) *KORM_test01Crud_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_UpdateObj_Name) In(vList []string) *KORM_test01Crud_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test01Crud_D_UpdateObj) CondMultOpBegin_AND() *KORM_test01Crud_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test01Crud_D_UpdateObj) CondMultOpBegin_OR() *KORM_test01Crud_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test01Crud_D_UpdateObj) CondMultOpEnd() *KORM_test01Crud_D_UpdateObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test01Crud_D_DeleteObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test01Crud_D) Delete() *KORM_test01Crud_D_DeleteObj {
	return &KORM_test01Crud_D_DeleteObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test01Crud_D_DeleteObj) MustRun() (RowsAffected int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("DELETE FROM test01Crud_D ")
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error
	result, err = this.supper.ExecRaw(buf2.String(), this.argsWhere...)

	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test01Crud_D_DeleteObj_Id struct {
	supper *KORM_test01Crud_D_DeleteObj
}

func (this *KORM_test01Crud_D_DeleteObj) Where_Id() *KORM_Where_KORM_test01Crud_D_DeleteObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test01Crud_D_DeleteObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Id) Equal(Id int) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Id) NotEqual(Id int) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Id) GreaterOrEqual(Id int) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Id) Less(Id int) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Id) LessOrEqual(Id int) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Id) In(vList []int) *KORM_test01Crud_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test01Crud_D_DeleteObj_Name struct {
	supper *KORM_test01Crud_D_DeleteObj
}

func (this *KORM_test01Crud_D_DeleteObj) Where_Name() *KORM_Where_KORM_test01Crud_D_DeleteObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test01Crud_D_DeleteObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Name) Equal(Name string) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Name) NotEqual(Name string) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Name) GreaterOrEqual(Name string) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Name) Less(Name string) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Name) LessOrEqual(Name string) *KORM_test01Crud_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test01Crud_D_DeleteObj_Name) In(vList []string) *KORM_test01Crud_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test01Crud_D_DeleteObj) CondMultOpBegin_AND() *KORM_test01Crud_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test01Crud_D_DeleteObj) CondMultOpBegin_OR() *KORM_test01Crud_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test01Crud_D_DeleteObj) CondMultOpEnd() *KORM_test01Crud_D_DeleteObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test02MultiplePk_D struct {
	supper *OrmAll
}

func (this *OrmAll) test02MultiplePk_D() *KORM_test02MultiplePk_D {
	return &KORM_test02MultiplePk_D{supper: this}
}
func korm_fillSelectFieldNameList_test02MultiplePk_D(joinNode *korm_leftJoinNode) {
	joinNode.selectFieldNameList = []string{"UserId", "GroupId", "CreateTime"}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("korm_fillSelectFieldNameList_test02MultiplePk_D" + strconv.Quote(sub.fieldName))
		}
	}
}
func (this *KORM_test02MultiplePk_D) MustInsert(info test02MultiplePk_D) {
	var err error
	vCreateTime := info.CreateTime.UTC().Format(time.RFC3339Nano)
	_, err = this.supper.ExecRaw("INSERT INTO `test02MultiplePk_D`(`UserId` ,`GroupId` ,`CreateTime` ) VALUES(?,?,?)", info.UserId, info.GroupId, vCreateTime)
	if err != nil {
		panic(err)
	}

	return
}
func (this *KORM_test02MultiplePk_D) MustSet(info test02MultiplePk_D) {
	var err error
	vCreateTime := info.CreateTime.UTC().Format(time.RFC3339Nano)
	_, err = this.supper.ExecRaw("REPLACE INTO `test02MultiplePk_D`(`UserId` ,`GroupId` ,`CreateTime` ) VALUES(?,?,?)", info.UserId, info.GroupId, vCreateTime)
	if err != nil {
		panic(err)
	}

	return
}

// Select test02MultiplePk_D
type KORM_test02MultiplePk_D_SelectObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
	joinNode    *korm_leftJoinNode
	rootInfo    korm_leftJoinRootInfo
}

func (this *KORM_test02MultiplePk_D) Select() *KORM_test02MultiplePk_D_SelectObj {
	one := &KORM_test02MultiplePk_D_SelectObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
	one.joinNode = &korm_leftJoinNode{tableName: "_0"}
	one.joinNode.root = &one.rootInfo
	one.rootInfo.tableNameIdx = 1
	return one
}

type KORM_test02MultiplePk_D_SelectObj_OrderByObj struct {
	fieldName string
	supper    *KORM_test02MultiplePk_D_SelectObj
}

func (this *KORM_test02MultiplePk_D_SelectObj_OrderByObj) ASC() *KORM_test02MultiplePk_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_test02MultiplePk_D_SelectObj_OrderByObj) DESC() *KORM_test02MultiplePk_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_test02MultiplePk_D_SelectObj) OrderBy_UserId() *KORM_test02MultiplePk_D_SelectObj_OrderByObj {
	return &KORM_test02MultiplePk_D_SelectObj_OrderByObj{
		fieldName: "UserId",
		supper:    this,
	}
}
func (this *KORM_test02MultiplePk_D_SelectObj) OrderBy_GroupId() *KORM_test02MultiplePk_D_SelectObj_OrderByObj {
	return &KORM_test02MultiplePk_D_SelectObj_OrderByObj{
		fieldName: "GroupId",
		supper:    this,
	}
}
func (this *KORM_test02MultiplePk_D_SelectObj) OrderBy_CreateTime() *KORM_test02MultiplePk_D_SelectObj_OrderByObj {
	return &KORM_test02MultiplePk_D_SelectObj_OrderByObj{
		fieldName: "CreateTime",
		supper:    this,
	}
}

func (this *KORM_test02MultiplePk_D_SelectObj) LimitOffset(limit int, offset int) *KORM_test02MultiplePk_D_SelectObj {
	this.limit = limit
	this.offset = offset
	return this
}

// pageSize: [1, n)
// pageNo:   [1,n)
func (this *KORM_test02MultiplePk_D_SelectObj) SetPageLimit(pageSize int, pageNo int) *KORM_test02MultiplePk_D_SelectObj {
	if pageSize <= 0 || pageNo <= 0 {
		panic("KORM_test02MultiplePk_D_SelectObj SetPageLimit error param")
	}
	this.limit = pageSize
	this.offset = pageSize * (pageNo - 1)
	return this
}

func (this *KORM_test02MultiplePk_D_SelectObj) MustRun_Count() (cnt int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test02MultiplePk_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test02MultiplePk_D_SelectObj not found.")
	}
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt
}

func (this *KORM_test02MultiplePk_D_SelectObj) MustRun_Exist() (exist bool) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test02MultiplePk_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	buf2.WriteString(" LIMIT 1 ")

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test02MultiplePk_D_SelectObj not found.")
	}
	var cnt int64
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func (this *KORM_test02MultiplePk_D_SelectObj) MustRun_ResultOne2() (info test02MultiplePk_D, ok bool) {
	this.limit = 1
	korm_fillSelectFieldNameList_test02MultiplePk_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test02MultiplePk_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		return info, false
	}
	resp := korm_test02MultiplePk_D_scan(this.joinNode, &info)
	err = result.Scan(resp.argList...)
	if err != nil {
		panic(err)
	}

	for idx, a := range resp.argList {
		v := a.(*sql.NullString).String
		if v == "" {
			continue
		}
		err = resp.argParseFn[idx](v)
		if err != nil {
			panic(err)
		}

	}

	return info, true
}

func (this *KORM_test02MultiplePk_D_SelectObj) MustRun_ResultOne() (info test02MultiplePk_D) {
	info, _ = this.MustRun_ResultOne2()
	return info
}
func (this *KORM_test02MultiplePk_D_SelectObj) MustRun_ResultList() (list []test02MultiplePk_D) {
	korm_fillSelectFieldNameList_test02MultiplePk_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test02MultiplePk_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test02MultiplePk_D
		korm_fillSelectFieldNameList_test02MultiplePk_D(this.joinNode)
		resp := korm_test02MultiplePk_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	return list
}
func (this *KORM_test02MultiplePk_D_SelectObj) MustRun_ResultMap() (m map[string]map[string]test02MultiplePk_D) {
	korm_fillSelectFieldNameList_test02MultiplePk_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test02MultiplePk_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test02MultiplePk_D
		korm_fillSelectFieldNameList_test02MultiplePk_D(this.joinNode)
		resp := korm_test02MultiplePk_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		if m == nil {
			m = map[string]map[string]test02MultiplePk_D{}
		}
		if m[info.UserId] == nil {
			m[info.UserId] = map[string]test02MultiplePk_D{}
		}
		m[info.UserId][info.GroupId] = info

	}
	return m
}
func (this *KORM_test02MultiplePk_D_SelectObj) MustRun_ResultListWithTotalMatch() (list []test02MultiplePk_D, totalMatch int64) {
	var err error
	korm_fillSelectFieldNameList_test02MultiplePk_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	if this.supper.mode == korm.InitModeMysql {
		buf2.WriteString("SQL_CALC_FOUND_ROWS  ")
	}
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test02MultiplePk_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}
	var conn *sql.Conn
	var result *sql.Rows
	if this.supper.db != nil {
		var err error
		conn, err = this.supper.db.Conn(context.Background())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		result, err = conn.QueryContext(context.Background(), buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		result, err = this.supper.tx.Query(buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	}

	defer result.Close()
	for result.Next() {
		var info test02MultiplePk_D
		korm_fillSelectFieldNameList_test02MultiplePk_D(this.joinNode)
		resp := korm_test02MultiplePk_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	result.Close()
	if this.limit > 0 {
		nextQuery := ""
		if this.supper.mode == korm.InitModeMysql {
			nextQuery = "select FOUND_ROWS()"
		} else if this.supper.mode == korm.InitModeSqlite {
			buf2.Reset()
			buf2.WriteString("SELECT COUNT(1) ")
			buf2.WriteString("FROM `test02MultiplePk_D` " + this.joinNode.tableName + " ")
			buf2.WriteString(this.rootInfo.bufLeftJoin.String())
			buf2.WriteString(this.bufWhere.String())
			nextQuery = buf2.String()
		} else {
			panic("not support")
		}
		var result2 *sql.Rows
		if conn != nil {
			result2, err = conn.QueryContext(context.Background(), nextQuery)
		} else {
			result2, err = this.supper.tx.Query(nextQuery)
		}
		if err != nil {
			panic(err)
		}

		defer result2.Close()

		if result2.Next() == false {
			panic("MustRun_ResultListWithPageInfo ")
		}
		err = result2.Scan(&totalMatch)
		if err != nil {
			panic(err)
		}

	}
	return list, totalMatch
}

type KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId struct {
	supper *KORM_test02MultiplePk_D_SelectObj
}

func (this *KORM_test02MultiplePk_D_SelectObj) Where_UserId() *KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`UserId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId) Equal(UserId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId) NotEqual(UserId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId) GreaterOrEqual(UserId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId) Less(UserId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId) LessOrEqual(UserId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_UserId) In(vList []string) *KORM_test02MultiplePk_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId struct {
	supper *KORM_test02MultiplePk_D_SelectObj
}

func (this *KORM_test02MultiplePk_D_SelectObj) Where_GroupId() *KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`GroupId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId) Equal(GroupId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId) NotEqual(GroupId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId) GreaterOrEqual(GroupId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId) Less(GroupId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId) LessOrEqual(GroupId string) *KORM_test02MultiplePk_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_GroupId) In(vList []string) *KORM_test02MultiplePk_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test02MultiplePk_D_SelectObj_CreateTime struct {
	supper *KORM_test02MultiplePk_D_SelectObj
}

func (this *KORM_test02MultiplePk_D_SelectObj) Where_CreateTime() *KORM_Where_KORM_test02MultiplePk_D_SelectObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_SelectObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_CreateTime) Equal(CreateTime time.Time) *KORM_test02MultiplePk_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_CreateTime) Less(CreateTime time.Time) *KORM_test02MultiplePk_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_SelectObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test02MultiplePk_D_SelectObj) CondMultOpBegin_AND() *KORM_test02MultiplePk_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test02MultiplePk_D_SelectObj) CondMultOpBegin_OR() *KORM_test02MultiplePk_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test02MultiplePk_D_SelectObj) CondMultOpEnd() *KORM_test02MultiplePk_D_SelectObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

// Update test02MultiplePk_D
type KORM_test02MultiplePk_D_UpdateObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test02MultiplePk_D) Update() *KORM_test02MultiplePk_D_UpdateObj {
	return &KORM_test02MultiplePk_D_UpdateObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test02MultiplePk_D) MustUpdateBy_UserId_GroupId(info test02MultiplePk_D) (rowsAffected int64) {
	rowsAffected = this.Update().Where_UserId().Equal(info.UserId).Where_GroupId().Equal(info.GroupId).Set_CreateTime(info.CreateTime).MustRun()
	return rowsAffected
}
func (this *KORM_test02MultiplePk_D_UpdateObj) Set_CreateTime(CreateTime time.Time) *KORM_test02MultiplePk_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`CreateTime` = ? ")
	this.argsSet = append(this.argsSet, CreateTime.UTC().Format(time.RFC3339Nano))
	return this
}
func (this *KORM_test02MultiplePk_D_UpdateObj) MustRun() (RowsAffected int64) {
	if len(this.argsSet) == 0 {
		panic("len(this.argsSet) == 0")
	}
	var buf2 bytes.Buffer
	buf2.WriteString("UPDATE `test02MultiplePk_D` ")
	buf2.WriteString(this.bufSet.String())
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error

	result, err = this.supper.ExecRaw(buf2.String(), append(this.argsSet, this.argsWhere...)...)
	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId struct {
	supper *KORM_test02MultiplePk_D_UpdateObj
}

func (this *KORM_test02MultiplePk_D_UpdateObj) Where_UserId() *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`UserId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId) Equal(UserId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId) NotEqual(UserId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId) GreaterOrEqual(UserId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId) Less(UserId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId) LessOrEqual(UserId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_UserId) In(vList []string) *KORM_test02MultiplePk_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId struct {
	supper *KORM_test02MultiplePk_D_UpdateObj
}

func (this *KORM_test02MultiplePk_D_UpdateObj) Where_GroupId() *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`GroupId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId) Equal(GroupId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId) NotEqual(GroupId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId) GreaterOrEqual(GroupId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId) Less(GroupId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId) LessOrEqual(GroupId string) *KORM_test02MultiplePk_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_GroupId) In(vList []string) *KORM_test02MultiplePk_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test02MultiplePk_D_UpdateObj_CreateTime struct {
	supper *KORM_test02MultiplePk_D_UpdateObj
}

func (this *KORM_test02MultiplePk_D_UpdateObj) Where_CreateTime() *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_UpdateObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_CreateTime) Equal(CreateTime time.Time) *KORM_test02MultiplePk_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_CreateTime) Less(CreateTime time.Time) *KORM_test02MultiplePk_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_UpdateObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test02MultiplePk_D_UpdateObj) CondMultOpBegin_AND() *KORM_test02MultiplePk_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test02MultiplePk_D_UpdateObj) CondMultOpBegin_OR() *KORM_test02MultiplePk_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test02MultiplePk_D_UpdateObj) CondMultOpEnd() *KORM_test02MultiplePk_D_UpdateObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test02MultiplePk_D_DeleteObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test02MultiplePk_D) Delete() *KORM_test02MultiplePk_D_DeleteObj {
	return &KORM_test02MultiplePk_D_DeleteObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test02MultiplePk_D_DeleteObj) MustRun() (RowsAffected int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("DELETE FROM test02MultiplePk_D ")
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error
	result, err = this.supper.ExecRaw(buf2.String(), this.argsWhere...)

	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId struct {
	supper *KORM_test02MultiplePk_D_DeleteObj
}

func (this *KORM_test02MultiplePk_D_DeleteObj) Where_UserId() *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`UserId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId) Equal(UserId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId) NotEqual(UserId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId) GreaterOrEqual(UserId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId) Less(UserId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId) LessOrEqual(UserId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_UserId) In(vList []string) *KORM_test02MultiplePk_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId struct {
	supper *KORM_test02MultiplePk_D_DeleteObj
}

func (this *KORM_test02MultiplePk_D_DeleteObj) Where_GroupId() *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`GroupId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId) Equal(GroupId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId) NotEqual(GroupId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId) GreaterOrEqual(GroupId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId) Less(GroupId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId) LessOrEqual(GroupId string) *KORM_test02MultiplePk_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_GroupId) In(vList []string) *KORM_test02MultiplePk_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test02MultiplePk_D_DeleteObj_CreateTime struct {
	supper *KORM_test02MultiplePk_D_DeleteObj
}

func (this *KORM_test02MultiplePk_D_DeleteObj) Where_CreateTime() *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test02MultiplePk_D_DeleteObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_CreateTime) Equal(CreateTime time.Time) *KORM_test02MultiplePk_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_CreateTime) Less(CreateTime time.Time) *KORM_test02MultiplePk_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test02MultiplePk_D_DeleteObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test02MultiplePk_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test02MultiplePk_D_DeleteObj) CondMultOpBegin_AND() *KORM_test02MultiplePk_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test02MultiplePk_D_DeleteObj) CondMultOpBegin_OR() *KORM_test02MultiplePk_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test02MultiplePk_D_DeleteObj) CondMultOpEnd() *KORM_test02MultiplePk_D_DeleteObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test03User_D struct {
	supper *OrmAll
}

func (this *OrmAll) test03User_D() *KORM_test03User_D {
	return &KORM_test03User_D{supper: this}
}
func korm_fillSelectFieldNameList_test03User_D(joinNode *korm_leftJoinNode) {
	joinNode.selectFieldNameList = []string{"Id", "Name", "PasswordHash", "CreateTime"}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("korm_fillSelectFieldNameList_test03User_D" + strconv.Quote(sub.fieldName))
		}
	}
}
func (this *KORM_test03User_D) MustInsert(info test03User_D) {
	var err error
	vCreateTime := info.CreateTime.UTC().Format(time.RFC3339Nano)
	_, err = this.supper.ExecRaw("INSERT INTO `test03User_D`(`Id` ,`Name` ,`PasswordHash` ,`CreateTime` ) VALUES(?,?,?,?)", info.Id, info.Name, info.PasswordHash, vCreateTime)
	if err != nil {
		panic(err)
	}

	return
}
func (this *KORM_test03User_D) MustSet(info test03User_D) {
	var err error
	vCreateTime := info.CreateTime.UTC().Format(time.RFC3339Nano)
	_, err = this.supper.ExecRaw("REPLACE INTO `test03User_D`(`Id` ,`Name` ,`PasswordHash` ,`CreateTime` ) VALUES(?,?,?,?)", info.Id, info.Name, info.PasswordHash, vCreateTime)
	if err != nil {
		panic(err)
	}

	return
}

// Select test03User_D
type KORM_test03User_D_SelectObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
	joinNode    *korm_leftJoinNode
	rootInfo    korm_leftJoinRootInfo
}

func (this *KORM_test03User_D) Select() *KORM_test03User_D_SelectObj {
	one := &KORM_test03User_D_SelectObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
	one.joinNode = &korm_leftJoinNode{tableName: "_0"}
	one.joinNode.root = &one.rootInfo
	one.rootInfo.tableNameIdx = 1
	return one
}

type KORM_test03User_D_SelectObj_OrderByObj struct {
	fieldName string
	supper    *KORM_test03User_D_SelectObj
}

func (this *KORM_test03User_D_SelectObj_OrderByObj) ASC() *KORM_test03User_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_test03User_D_SelectObj_OrderByObj) DESC() *KORM_test03User_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_test03User_D_SelectObj) OrderBy_Id() *KORM_test03User_D_SelectObj_OrderByObj {
	return &KORM_test03User_D_SelectObj_OrderByObj{
		fieldName: "Id",
		supper:    this,
	}
}
func (this *KORM_test03User_D_SelectObj) OrderBy_Name() *KORM_test03User_D_SelectObj_OrderByObj {
	return &KORM_test03User_D_SelectObj_OrderByObj{
		fieldName: "Name",
		supper:    this,
	}
}
func (this *KORM_test03User_D_SelectObj) OrderBy_PasswordHash() *KORM_test03User_D_SelectObj_OrderByObj {
	return &KORM_test03User_D_SelectObj_OrderByObj{
		fieldName: "PasswordHash",
		supper:    this,
	}
}
func (this *KORM_test03User_D_SelectObj) OrderBy_CreateTime() *KORM_test03User_D_SelectObj_OrderByObj {
	return &KORM_test03User_D_SelectObj_OrderByObj{
		fieldName: "CreateTime",
		supper:    this,
	}
}

func (this *KORM_test03User_D_SelectObj) LimitOffset(limit int, offset int) *KORM_test03User_D_SelectObj {
	this.limit = limit
	this.offset = offset
	return this
}

// pageSize: [1, n)
// pageNo:   [1,n)
func (this *KORM_test03User_D_SelectObj) SetPageLimit(pageSize int, pageNo int) *KORM_test03User_D_SelectObj {
	if pageSize <= 0 || pageNo <= 0 {
		panic("KORM_test03User_D_SelectObj SetPageLimit error param")
	}
	this.limit = pageSize
	this.offset = pageSize * (pageNo - 1)
	return this
}

func (this *KORM_test03User_D_SelectObj) MustRun_Count() (cnt int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test03User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test03User_D_SelectObj not found.")
	}
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt
}

func (this *KORM_test03User_D_SelectObj) MustRun_Exist() (exist bool) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test03User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	buf2.WriteString(" LIMIT 1 ")

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test03User_D_SelectObj not found.")
	}
	var cnt int64
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func (this *KORM_test03User_D_SelectObj) MustRun_ResultOne2() (info test03User_D, ok bool) {
	this.limit = 1
	korm_fillSelectFieldNameList_test03User_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		return info, false
	}
	resp := korm_test03User_D_scan(this.joinNode, &info)
	err = result.Scan(resp.argList...)
	if err != nil {
		panic(err)
	}

	for idx, a := range resp.argList {
		v := a.(*sql.NullString).String
		if v == "" {
			continue
		}
		err = resp.argParseFn[idx](v)
		if err != nil {
			panic(err)
		}

	}

	return info, true
}

func (this *KORM_test03User_D_SelectObj) MustRun_ResultOne() (info test03User_D) {
	info, _ = this.MustRun_ResultOne2()
	return info
}
func (this *KORM_test03User_D_SelectObj) MustRun_ResultList() (list []test03User_D) {
	korm_fillSelectFieldNameList_test03User_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test03User_D
		korm_fillSelectFieldNameList_test03User_D(this.joinNode)
		resp := korm_test03User_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	return list
}
func (this *KORM_test03User_D_SelectObj) MustRun_ResultMap() (m map[string]test03User_D) {
	korm_fillSelectFieldNameList_test03User_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test03User_D
		korm_fillSelectFieldNameList_test03User_D(this.joinNode)
		resp := korm_test03User_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		if m == nil {
			m = map[string]test03User_D{}
		}
		m[info.Id] = info

	}
	return m
}
func (this *KORM_test03User_D_SelectObj) MustRun_ResultListWithTotalMatch() (list []test03User_D, totalMatch int64) {
	var err error
	korm_fillSelectFieldNameList_test03User_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	if this.supper.mode == korm.InitModeMysql {
		buf2.WriteString("SQL_CALC_FOUND_ROWS  ")
	}
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}
	var conn *sql.Conn
	var result *sql.Rows
	if this.supper.db != nil {
		var err error
		conn, err = this.supper.db.Conn(context.Background())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		result, err = conn.QueryContext(context.Background(), buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		result, err = this.supper.tx.Query(buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	}

	defer result.Close()
	for result.Next() {
		var info test03User_D
		korm_fillSelectFieldNameList_test03User_D(this.joinNode)
		resp := korm_test03User_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	result.Close()
	if this.limit > 0 {
		nextQuery := ""
		if this.supper.mode == korm.InitModeMysql {
			nextQuery = "select FOUND_ROWS()"
		} else if this.supper.mode == korm.InitModeSqlite {
			buf2.Reset()
			buf2.WriteString("SELECT COUNT(1) ")
			buf2.WriteString("FROM `test03User_D` " + this.joinNode.tableName + " ")
			buf2.WriteString(this.rootInfo.bufLeftJoin.String())
			buf2.WriteString(this.bufWhere.String())
			nextQuery = buf2.String()
		} else {
			panic("not support")
		}
		var result2 *sql.Rows
		if conn != nil {
			result2, err = conn.QueryContext(context.Background(), nextQuery)
		} else {
			result2, err = this.supper.tx.Query(nextQuery)
		}
		if err != nil {
			panic(err)
		}

		defer result2.Close()

		if result2.Next() == false {
			panic("MustRun_ResultListWithPageInfo ")
		}
		err = result2.Scan(&totalMatch)
		if err != nil {
			panic(err)
		}

	}
	return list, totalMatch
}

type KORM_Where_KORM_test03User_D_SelectObj_Id struct {
	supper *KORM_test03User_D_SelectObj
}

func (this *KORM_test03User_D_SelectObj) Where_Id() *KORM_Where_KORM_test03User_D_SelectObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_SelectObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Id) Equal(Id string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Id) NotEqual(Id string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Id) GreaterOrEqual(Id string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Id) Less(Id string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Id) LessOrEqual(Id string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Id) In(vList []string) *KORM_test03User_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_SelectObj_Name struct {
	supper *KORM_test03User_D_SelectObj
}

func (this *KORM_test03User_D_SelectObj) Where_Name() *KORM_Where_KORM_test03User_D_SelectObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_SelectObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Name) Equal(Name string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Name) NotEqual(Name string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Name) GreaterOrEqual(Name string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Name) Less(Name string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Name) LessOrEqual(Name string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_Name) In(vList []string) *KORM_test03User_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_SelectObj_PasswordHash struct {
	supper *KORM_test03User_D_SelectObj
}

func (this *KORM_test03User_D_SelectObj) Where_PasswordHash() *KORM_Where_KORM_test03User_D_SelectObj_PasswordHash {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`PasswordHash` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_SelectObj_PasswordHash{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_PasswordHash) Equal(PasswordHash string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_PasswordHash) NotEqual(PasswordHash string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_PasswordHash) GreaterOrEqual(PasswordHash string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_PasswordHash) Less(PasswordHash string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_PasswordHash) LessOrEqual(PasswordHash string) *KORM_test03User_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_PasswordHash) In(vList []string) *KORM_test03User_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_SelectObj_CreateTime struct {
	supper *KORM_test03User_D_SelectObj
}

func (this *KORM_test03User_D_SelectObj) Where_CreateTime() *KORM_Where_KORM_test03User_D_SelectObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_SelectObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_CreateTime) Equal(CreateTime time.Time) *KORM_test03User_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test03User_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test03User_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_CreateTime) Less(CreateTime time.Time) *KORM_test03User_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_SelectObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test03User_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test03User_D_SelectObj) CondMultOpBegin_AND() *KORM_test03User_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03User_D_SelectObj) CondMultOpBegin_OR() *KORM_test03User_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03User_D_SelectObj) CondMultOpEnd() *KORM_test03User_D_SelectObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

// Update test03User_D
type KORM_test03User_D_UpdateObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test03User_D) Update() *KORM_test03User_D_UpdateObj {
	return &KORM_test03User_D_UpdateObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test03User_D) MustUpdateBy_Id(info test03User_D) (rowsAffected int64) {
	rowsAffected = this.Update().Where_Id().Equal(info.Id).Set_Name(info.Name).Set_PasswordHash(info.PasswordHash).Set_CreateTime(info.CreateTime).MustRun()
	return rowsAffected
}
func (this *KORM_test03User_D_UpdateObj) Set_Name(Name string) *KORM_test03User_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`Name` = ? ")
	this.argsSet = append(this.argsSet, Name)
	return this
}
func (this *KORM_test03User_D_UpdateObj) Set_PasswordHash(PasswordHash string) *KORM_test03User_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`PasswordHash` = ? ")
	this.argsSet = append(this.argsSet, PasswordHash)
	return this
}
func (this *KORM_test03User_D_UpdateObj) Set_CreateTime(CreateTime time.Time) *KORM_test03User_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`CreateTime` = ? ")
	this.argsSet = append(this.argsSet, CreateTime.UTC().Format(time.RFC3339Nano))
	return this
}
func (this *KORM_test03User_D_UpdateObj) MustRun() (RowsAffected int64) {
	if len(this.argsSet) == 0 {
		panic("len(this.argsSet) == 0")
	}
	var buf2 bytes.Buffer
	buf2.WriteString("UPDATE `test03User_D` ")
	buf2.WriteString(this.bufSet.String())
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error

	result, err = this.supper.ExecRaw(buf2.String(), append(this.argsSet, this.argsWhere...)...)
	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test03User_D_UpdateObj_Id struct {
	supper *KORM_test03User_D_UpdateObj
}

func (this *KORM_test03User_D_UpdateObj) Where_Id() *KORM_Where_KORM_test03User_D_UpdateObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_UpdateObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Id) Equal(Id string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Id) NotEqual(Id string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Id) GreaterOrEqual(Id string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Id) Less(Id string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Id) LessOrEqual(Id string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Id) In(vList []string) *KORM_test03User_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_UpdateObj_Name struct {
	supper *KORM_test03User_D_UpdateObj
}

func (this *KORM_test03User_D_UpdateObj) Where_Name() *KORM_Where_KORM_test03User_D_UpdateObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_UpdateObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Name) Equal(Name string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Name) NotEqual(Name string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Name) GreaterOrEqual(Name string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Name) Less(Name string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Name) LessOrEqual(Name string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_Name) In(vList []string) *KORM_test03User_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash struct {
	supper *KORM_test03User_D_UpdateObj
}

func (this *KORM_test03User_D_UpdateObj) Where_PasswordHash() *KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`PasswordHash` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash) Equal(PasswordHash string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash) NotEqual(PasswordHash string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash) GreaterOrEqual(PasswordHash string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash) Less(PasswordHash string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash) LessOrEqual(PasswordHash string) *KORM_test03User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_PasswordHash) In(vList []string) *KORM_test03User_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_UpdateObj_CreateTime struct {
	supper *KORM_test03User_D_UpdateObj
}

func (this *KORM_test03User_D_UpdateObj) Where_CreateTime() *KORM_Where_KORM_test03User_D_UpdateObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_UpdateObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_CreateTime) Equal(CreateTime time.Time) *KORM_test03User_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test03User_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test03User_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_CreateTime) Less(CreateTime time.Time) *KORM_test03User_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_UpdateObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test03User_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test03User_D_UpdateObj) CondMultOpBegin_AND() *KORM_test03User_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03User_D_UpdateObj) CondMultOpBegin_OR() *KORM_test03User_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03User_D_UpdateObj) CondMultOpEnd() *KORM_test03User_D_UpdateObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test03User_D_DeleteObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test03User_D) Delete() *KORM_test03User_D_DeleteObj {
	return &KORM_test03User_D_DeleteObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test03User_D_DeleteObj) MustRun() (RowsAffected int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("DELETE FROM test03User_D ")
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error
	result, err = this.supper.ExecRaw(buf2.String(), this.argsWhere...)

	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test03User_D_DeleteObj_Id struct {
	supper *KORM_test03User_D_DeleteObj
}

func (this *KORM_test03User_D_DeleteObj) Where_Id() *KORM_Where_KORM_test03User_D_DeleteObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_DeleteObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Id) Equal(Id string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Id) NotEqual(Id string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Id) GreaterOrEqual(Id string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Id) Less(Id string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Id) LessOrEqual(Id string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Id) In(vList []string) *KORM_test03User_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_DeleteObj_Name struct {
	supper *KORM_test03User_D_DeleteObj
}

func (this *KORM_test03User_D_DeleteObj) Where_Name() *KORM_Where_KORM_test03User_D_DeleteObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_DeleteObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Name) Equal(Name string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Name) NotEqual(Name string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Name) GreaterOrEqual(Name string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Name) Less(Name string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Name) LessOrEqual(Name string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_Name) In(vList []string) *KORM_test03User_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash struct {
	supper *KORM_test03User_D_DeleteObj
}

func (this *KORM_test03User_D_DeleteObj) Where_PasswordHash() *KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`PasswordHash` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash) Equal(PasswordHash string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash) NotEqual(PasswordHash string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash) GreaterOrEqual(PasswordHash string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash) Less(PasswordHash string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash) LessOrEqual(PasswordHash string) *KORM_test03User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_PasswordHash) In(vList []string) *KORM_test03User_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_DeleteObj_CreateTime struct {
	supper *KORM_test03User_D_DeleteObj
}

func (this *KORM_test03User_D_DeleteObj) Where_CreateTime() *KORM_Where_KORM_test03User_D_DeleteObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03User_D_DeleteObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_CreateTime) Equal(CreateTime time.Time) *KORM_test03User_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test03User_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test03User_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_CreateTime) Less(CreateTime time.Time) *KORM_test03User_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_DeleteObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test03User_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test03User_D_DeleteObj) CondMultOpBegin_AND() *KORM_test03User_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03User_D_DeleteObj) CondMultOpBegin_OR() *KORM_test03User_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03User_D_DeleteObj) CondMultOpEnd() *KORM_test03User_D_DeleteObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test03Group_D struct {
	supper *OrmAll
}

func (this *OrmAll) test03Group_D() *KORM_test03Group_D {
	return &KORM_test03Group_D{supper: this}
}
func korm_fillSelectFieldNameList_test03Group_D(joinNode *korm_leftJoinNode) {
	joinNode.selectFieldNameList = []string{"Id", "Name", "CreateTime"}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("korm_fillSelectFieldNameList_test03Group_D" + strconv.Quote(sub.fieldName))
		}
	}
}
func (this *KORM_test03Group_D) MustInsert(info test03Group_D) {
	var err error
	vCreateTime := info.CreateTime.UTC().Format(time.RFC3339Nano)
	_, err = this.supper.ExecRaw("INSERT INTO `test03Group_D`(`Id` ,`Name` ,`CreateTime` ) VALUES(?,?,?)", info.Id, info.Name, vCreateTime)
	if err != nil {
		panic(err)
	}

	return
}
func (this *KORM_test03Group_D) MustSet(info test03Group_D) {
	var err error
	vCreateTime := info.CreateTime.UTC().Format(time.RFC3339Nano)
	_, err = this.supper.ExecRaw("REPLACE INTO `test03Group_D`(`Id` ,`Name` ,`CreateTime` ) VALUES(?,?,?)", info.Id, info.Name, vCreateTime)
	if err != nil {
		panic(err)
	}

	return
}

// Select test03Group_D
type KORM_test03Group_D_SelectObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
	joinNode    *korm_leftJoinNode
	rootInfo    korm_leftJoinRootInfo
}

func (this *KORM_test03Group_D) Select() *KORM_test03Group_D_SelectObj {
	one := &KORM_test03Group_D_SelectObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
	one.joinNode = &korm_leftJoinNode{tableName: "_0"}
	one.joinNode.root = &one.rootInfo
	one.rootInfo.tableNameIdx = 1
	return one
}

type KORM_test03Group_D_SelectObj_OrderByObj struct {
	fieldName string
	supper    *KORM_test03Group_D_SelectObj
}

func (this *KORM_test03Group_D_SelectObj_OrderByObj) ASC() *KORM_test03Group_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_test03Group_D_SelectObj_OrderByObj) DESC() *KORM_test03Group_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_test03Group_D_SelectObj) OrderBy_Name() *KORM_test03Group_D_SelectObj_OrderByObj {
	return &KORM_test03Group_D_SelectObj_OrderByObj{
		fieldName: "Name",
		supper:    this,
	}
}
func (this *KORM_test03Group_D_SelectObj) OrderBy_CreateTime() *KORM_test03Group_D_SelectObj_OrderByObj {
	return &KORM_test03Group_D_SelectObj_OrderByObj{
		fieldName: "CreateTime",
		supper:    this,
	}
}

func (this *KORM_test03Group_D_SelectObj) LimitOffset(limit int, offset int) *KORM_test03Group_D_SelectObj {
	this.limit = limit
	this.offset = offset
	return this
}

// pageSize: [1, n)
// pageNo:   [1,n)
func (this *KORM_test03Group_D_SelectObj) SetPageLimit(pageSize int, pageNo int) *KORM_test03Group_D_SelectObj {
	if pageSize <= 0 || pageNo <= 0 {
		panic("KORM_test03Group_D_SelectObj SetPageLimit error param")
	}
	this.limit = pageSize
	this.offset = pageSize * (pageNo - 1)
	return this
}

func (this *KORM_test03Group_D_SelectObj) MustRun_Count() (cnt int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test03Group_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test03Group_D_SelectObj not found.")
	}
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt
}

func (this *KORM_test03Group_D_SelectObj) MustRun_Exist() (exist bool) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test03Group_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	buf2.WriteString(" LIMIT 1 ")

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test03Group_D_SelectObj not found.")
	}
	var cnt int64
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func (this *KORM_test03Group_D_SelectObj) MustRun_ResultOne2() (info test03Group_D, ok bool) {
	this.limit = 1
	korm_fillSelectFieldNameList_test03Group_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03Group_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		return info, false
	}
	resp := korm_test03Group_D_scan(this.joinNode, &info)
	err = result.Scan(resp.argList...)
	if err != nil {
		panic(err)
	}

	for idx, a := range resp.argList {
		v := a.(*sql.NullString).String
		if v == "" {
			continue
		}
		err = resp.argParseFn[idx](v)
		if err != nil {
			panic(err)
		}

	}

	return info, true
}

func (this *KORM_test03Group_D_SelectObj) MustRun_ResultOne() (info test03Group_D) {
	info, _ = this.MustRun_ResultOne2()
	return info
}
func (this *KORM_test03Group_D_SelectObj) MustRun_ResultList() (list []test03Group_D) {
	korm_fillSelectFieldNameList_test03Group_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03Group_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test03Group_D
		korm_fillSelectFieldNameList_test03Group_D(this.joinNode)
		resp := korm_test03Group_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	return list
}
func (this *KORM_test03Group_D_SelectObj) MustRun_ResultMap() (m map[uint64]test03Group_D) {
	korm_fillSelectFieldNameList_test03Group_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03Group_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test03Group_D
		korm_fillSelectFieldNameList_test03Group_D(this.joinNode)
		resp := korm_test03Group_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		if m == nil {
			m = map[uint64]test03Group_D{}
		}
		m[info.Id] = info

	}
	return m
}
func (this *KORM_test03Group_D_SelectObj) MustRun_ResultListWithTotalMatch() (list []test03Group_D, totalMatch int64) {
	var err error
	korm_fillSelectFieldNameList_test03Group_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	if this.supper.mode == korm.InitModeMysql {
		buf2.WriteString("SQL_CALC_FOUND_ROWS  ")
	}
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03Group_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}
	var conn *sql.Conn
	var result *sql.Rows
	if this.supper.db != nil {
		var err error
		conn, err = this.supper.db.Conn(context.Background())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		result, err = conn.QueryContext(context.Background(), buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		result, err = this.supper.tx.Query(buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	}

	defer result.Close()
	for result.Next() {
		var info test03Group_D
		korm_fillSelectFieldNameList_test03Group_D(this.joinNode)
		resp := korm_test03Group_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	result.Close()
	if this.limit > 0 {
		nextQuery := ""
		if this.supper.mode == korm.InitModeMysql {
			nextQuery = "select FOUND_ROWS()"
		} else if this.supper.mode == korm.InitModeSqlite {
			buf2.Reset()
			buf2.WriteString("SELECT COUNT(1) ")
			buf2.WriteString("FROM `test03Group_D` " + this.joinNode.tableName + " ")
			buf2.WriteString(this.rootInfo.bufLeftJoin.String())
			buf2.WriteString(this.bufWhere.String())
			nextQuery = buf2.String()
		} else {
			panic("not support")
		}
		var result2 *sql.Rows
		if conn != nil {
			result2, err = conn.QueryContext(context.Background(), nextQuery)
		} else {
			result2, err = this.supper.tx.Query(nextQuery)
		}
		if err != nil {
			panic(err)
		}

		defer result2.Close()

		if result2.Next() == false {
			panic("MustRun_ResultListWithPageInfo ")
		}
		err = result2.Scan(&totalMatch)
		if err != nil {
			panic(err)
		}

	}
	return list, totalMatch
}

type KORM_Where_KORM_test03Group_D_SelectObj_Id struct {
	supper *KORM_test03Group_D_SelectObj
}

func (this *KORM_test03Group_D_SelectObj) Where_Id() *KORM_Where_KORM_test03Group_D_SelectObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_SelectObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Id) Equal(Id uint64) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Id) NotEqual(Id uint64) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Id) GreaterOrEqual(Id uint64) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Id) Less(Id uint64) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Id) LessOrEqual(Id uint64) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Id) In(vList []uint64) *KORM_test03Group_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03Group_D_SelectObj_Name struct {
	supper *KORM_test03Group_D_SelectObj
}

func (this *KORM_test03Group_D_SelectObj) Where_Name() *KORM_Where_KORM_test03Group_D_SelectObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_SelectObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Name) Equal(Name string) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Name) NotEqual(Name string) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Name) GreaterOrEqual(Name string) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Name) Less(Name string) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Name) LessOrEqual(Name string) *KORM_test03Group_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_Name) In(vList []string) *KORM_test03Group_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03Group_D_SelectObj_CreateTime struct {
	supper *KORM_test03Group_D_SelectObj
}

func (this *KORM_test03Group_D_SelectObj) Where_CreateTime() *KORM_Where_KORM_test03Group_D_SelectObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_SelectObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_CreateTime) Equal(CreateTime time.Time) *KORM_test03Group_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test03Group_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test03Group_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_CreateTime) Less(CreateTime time.Time) *KORM_test03Group_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_SelectObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test03Group_D_SelectObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test03Group_D_SelectObj) CondMultOpBegin_AND() *KORM_test03Group_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03Group_D_SelectObj) CondMultOpBegin_OR() *KORM_test03Group_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03Group_D_SelectObj) CondMultOpEnd() *KORM_test03Group_D_SelectObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

// Update test03Group_D
type KORM_test03Group_D_UpdateObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test03Group_D) Update() *KORM_test03Group_D_UpdateObj {
	return &KORM_test03Group_D_UpdateObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test03Group_D) MustUpdateBy_Id(info test03Group_D) (rowsAffected int64) {
	rowsAffected = this.Update().Where_Id().Equal(info.Id).Set_Name(info.Name).Set_CreateTime(info.CreateTime).MustRun()
	return rowsAffected
}
func (this *KORM_test03Group_D_UpdateObj) Set_Name(Name string) *KORM_test03Group_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`Name` = ? ")
	this.argsSet = append(this.argsSet, Name)
	return this
}
func (this *KORM_test03Group_D_UpdateObj) Set_CreateTime(CreateTime time.Time) *KORM_test03Group_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`CreateTime` = ? ")
	this.argsSet = append(this.argsSet, CreateTime.UTC().Format(time.RFC3339Nano))
	return this
}
func (this *KORM_test03Group_D_UpdateObj) MustRun() (RowsAffected int64) {
	if len(this.argsSet) == 0 {
		panic("len(this.argsSet) == 0")
	}
	var buf2 bytes.Buffer
	buf2.WriteString("UPDATE `test03Group_D` ")
	buf2.WriteString(this.bufSet.String())
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error

	result, err = this.supper.ExecRaw(buf2.String(), append(this.argsSet, this.argsWhere...)...)
	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test03Group_D_UpdateObj_Id struct {
	supper *KORM_test03Group_D_UpdateObj
}

func (this *KORM_test03Group_D_UpdateObj) Where_Id() *KORM_Where_KORM_test03Group_D_UpdateObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_UpdateObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Id) Equal(Id uint64) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Id) NotEqual(Id uint64) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Id) GreaterOrEqual(Id uint64) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Id) Less(Id uint64) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Id) LessOrEqual(Id uint64) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Id) In(vList []uint64) *KORM_test03Group_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03Group_D_UpdateObj_Name struct {
	supper *KORM_test03Group_D_UpdateObj
}

func (this *KORM_test03Group_D_UpdateObj) Where_Name() *KORM_Where_KORM_test03Group_D_UpdateObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_UpdateObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Name) Equal(Name string) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Name) NotEqual(Name string) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Name) GreaterOrEqual(Name string) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Name) Less(Name string) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Name) LessOrEqual(Name string) *KORM_test03Group_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_Name) In(vList []string) *KORM_test03Group_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03Group_D_UpdateObj_CreateTime struct {
	supper *KORM_test03Group_D_UpdateObj
}

func (this *KORM_test03Group_D_UpdateObj) Where_CreateTime() *KORM_Where_KORM_test03Group_D_UpdateObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_UpdateObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_CreateTime) Equal(CreateTime time.Time) *KORM_test03Group_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test03Group_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test03Group_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_CreateTime) Less(CreateTime time.Time) *KORM_test03Group_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_UpdateObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test03Group_D_UpdateObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test03Group_D_UpdateObj) CondMultOpBegin_AND() *KORM_test03Group_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03Group_D_UpdateObj) CondMultOpBegin_OR() *KORM_test03Group_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03Group_D_UpdateObj) CondMultOpEnd() *KORM_test03Group_D_UpdateObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test03Group_D_DeleteObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test03Group_D) Delete() *KORM_test03Group_D_DeleteObj {
	return &KORM_test03Group_D_DeleteObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test03Group_D_DeleteObj) MustRun() (RowsAffected int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("DELETE FROM test03Group_D ")
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error
	result, err = this.supper.ExecRaw(buf2.String(), this.argsWhere...)

	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test03Group_D_DeleteObj_Id struct {
	supper *KORM_test03Group_D_DeleteObj
}

func (this *KORM_test03Group_D_DeleteObj) Where_Id() *KORM_Where_KORM_test03Group_D_DeleteObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_DeleteObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Id) Equal(Id uint64) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Id) NotEqual(Id uint64) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Id) GreaterOrEqual(Id uint64) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Id) Less(Id uint64) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Id) LessOrEqual(Id uint64) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Id) In(vList []uint64) *KORM_test03Group_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03Group_D_DeleteObj_Name struct {
	supper *KORM_test03Group_D_DeleteObj
}

func (this *KORM_test03Group_D_DeleteObj) Where_Name() *KORM_Where_KORM_test03Group_D_DeleteObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_DeleteObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Name) Equal(Name string) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Name) NotEqual(Name string) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Name) GreaterOrEqual(Name string) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Name) Less(Name string) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Name) LessOrEqual(Name string) *KORM_test03Group_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_Name) In(vList []string) *KORM_test03Group_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03Group_D_DeleteObj_CreateTime struct {
	supper *KORM_test03Group_D_DeleteObj
}

func (this *KORM_test03Group_D_DeleteObj) Where_CreateTime() *KORM_Where_KORM_test03Group_D_DeleteObj_CreateTime {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`CreateTime` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03Group_D_DeleteObj_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_CreateTime) Equal(CreateTime time.Time) *KORM_test03Group_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_CreateTime) NotEqual(CreateTime time.Time) *KORM_test03Group_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test03Group_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_CreateTime) Less(CreateTime time.Time) *KORM_test03Group_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_DeleteObj_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test03Group_D_DeleteObj {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, vCreateTime)
	return this.supper
}
func (this *KORM_test03Group_D_DeleteObj) CondMultOpBegin_AND() *KORM_test03Group_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03Group_D_DeleteObj) CondMultOpBegin_OR() *KORM_test03Group_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03Group_D_DeleteObj) CondMultOpEnd() *KORM_test03Group_D_DeleteObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test03UserInGroup_D struct {
	supper *OrmAll
}

func (this *OrmAll) test03UserInGroup_D() *KORM_test03UserInGroup_D {
	return &KORM_test03UserInGroup_D{supper: this}
}
func korm_fillSelectFieldNameList_test03UserInGroup_D(joinNode *korm_leftJoinNode) {
	joinNode.selectFieldNameList = []string{"UserId", "GroupId"}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("korm_fillSelectFieldNameList_test03UserInGroup_D" + strconv.Quote(sub.fieldName))
		case "User":
			korm_fillSelectFieldNameList_test03User_D(sub)
		case "Group":
			korm_fillSelectFieldNameList_test03Group_D(sub)
		}
	}
}
func (this *KORM_test03UserInGroup_D) MustInsert(info test03UserInGroup_D) {
	var err error
	_, err = this.supper.ExecRaw("INSERT INTO `test03UserInGroup_D`(`UserId` ,`GroupId` ) VALUES(?,?)", info.UserId, info.GroupId)
	if err != nil {
		panic(err)
	}

	return
}
func (this *KORM_test03UserInGroup_D) MustSet(info test03UserInGroup_D) {
	var err error
	_, err = this.supper.ExecRaw("REPLACE INTO `test03UserInGroup_D`(`UserId` ,`GroupId` ) VALUES(?,?)", info.UserId, info.GroupId)
	if err != nil {
		panic(err)
	}

	return
}

// Select test03UserInGroup_D
type KORM_test03UserInGroup_D_SelectObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
	joinNode    *korm_leftJoinNode
	rootInfo    korm_leftJoinRootInfo
}

func (this *KORM_test03UserInGroup_D) Select() *KORM_test03UserInGroup_D_SelectObj {
	one := &KORM_test03UserInGroup_D_SelectObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
	one.joinNode = &korm_leftJoinNode{tableName: "_0"}
	one.joinNode.root = &one.rootInfo
	one.rootInfo.tableNameIdx = 1
	return one
}

type KORM_test03UserInGroup_D_SelectObj_OrderByObj struct {
	fieldName string
	supper    *KORM_test03UserInGroup_D_SelectObj
}

func (this *KORM_test03UserInGroup_D_SelectObj_OrderByObj) ASC() *KORM_test03UserInGroup_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_test03UserInGroup_D_SelectObj_OrderByObj) DESC() *KORM_test03UserInGroup_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_test03UserInGroup_D_SelectObj) OrderBy_UserId() *KORM_test03UserInGroup_D_SelectObj_OrderByObj {
	return &KORM_test03UserInGroup_D_SelectObj_OrderByObj{
		fieldName: "UserId",
		supper:    this,
	}
}

func (this *KORM_test03UserInGroup_D_SelectObj) LimitOffset(limit int, offset int) *KORM_test03UserInGroup_D_SelectObj {
	this.limit = limit
	this.offset = offset
	return this
}

// pageSize: [1, n)
// pageNo:   [1,n)
func (this *KORM_test03UserInGroup_D_SelectObj) SetPageLimit(pageSize int, pageNo int) *KORM_test03UserInGroup_D_SelectObj {
	if pageSize <= 0 || pageNo <= 0 {
		panic("KORM_test03UserInGroup_D_SelectObj SetPageLimit error param")
	}
	this.limit = pageSize
	this.offset = pageSize * (pageNo - 1)
	return this
}

func (this *KORM_test03UserInGroup_D_SelectObj) MustRun_Count() (cnt int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test03UserInGroup_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test03UserInGroup_D_SelectObj not found.")
	}
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt
}

func (this *KORM_test03UserInGroup_D_SelectObj) MustRun_Exist() (exist bool) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test03UserInGroup_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	buf2.WriteString(" LIMIT 1 ")

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test03UserInGroup_D_SelectObj not found.")
	}
	var cnt int64
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func (this *KORM_test03UserInGroup_D_SelectObj) MustRun_ResultOne2() (info test03UserInGroup_D, ok bool) {
	this.limit = 1
	korm_fillSelectFieldNameList_test03UserInGroup_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03UserInGroup_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		return info, false
	}
	resp := korm_test03UserInGroup_D_scan(this.joinNode, &info)
	err = result.Scan(resp.argList...)
	if err != nil {
		panic(err)
	}

	for idx, a := range resp.argList {
		v := a.(*sql.NullString).String
		if v == "" {
			continue
		}
		err = resp.argParseFn[idx](v)
		if err != nil {
			panic(err)
		}

	}

	return info, true
}

func (this *KORM_test03UserInGroup_D_SelectObj) MustRun_ResultOne() (info test03UserInGroup_D) {
	info, _ = this.MustRun_ResultOne2()
	return info
}
func (this *KORM_test03UserInGroup_D_SelectObj) MustRun_ResultList() (list []test03UserInGroup_D) {
	korm_fillSelectFieldNameList_test03UserInGroup_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03UserInGroup_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test03UserInGroup_D
		korm_fillSelectFieldNameList_test03UserInGroup_D(this.joinNode)
		resp := korm_test03UserInGroup_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	return list
}
func (this *KORM_test03UserInGroup_D_SelectObj) MustRun_ResultMap() (m map[string]test03UserInGroup_D) {
	korm_fillSelectFieldNameList_test03UserInGroup_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03UserInGroup_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test03UserInGroup_D
		korm_fillSelectFieldNameList_test03UserInGroup_D(this.joinNode)
		resp := korm_test03UserInGroup_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		if m == nil {
			m = map[string]test03UserInGroup_D{}
		}
		m[info.UserId] = info

	}
	return m
}
func (this *KORM_test03UserInGroup_D_SelectObj) MustRun_ResultListWithTotalMatch() (list []test03UserInGroup_D, totalMatch int64) {
	var err error
	korm_fillSelectFieldNameList_test03UserInGroup_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	if this.supper.mode == korm.InitModeMysql {
		buf2.WriteString("SQL_CALC_FOUND_ROWS  ")
	}
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test03UserInGroup_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}
	var conn *sql.Conn
	var result *sql.Rows
	if this.supper.db != nil {
		var err error
		conn, err = this.supper.db.Conn(context.Background())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		result, err = conn.QueryContext(context.Background(), buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		result, err = this.supper.tx.Query(buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	}

	defer result.Close()
	for result.Next() {
		var info test03UserInGroup_D
		korm_fillSelectFieldNameList_test03UserInGroup_D(this.joinNode)
		resp := korm_test03UserInGroup_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	result.Close()
	if this.limit > 0 {
		nextQuery := ""
		if this.supper.mode == korm.InitModeMysql {
			nextQuery = "select FOUND_ROWS()"
		} else if this.supper.mode == korm.InitModeSqlite {
			buf2.Reset()
			buf2.WriteString("SELECT COUNT(1) ")
			buf2.WriteString("FROM `test03UserInGroup_D` " + this.joinNode.tableName + " ")
			buf2.WriteString(this.rootInfo.bufLeftJoin.String())
			buf2.WriteString(this.bufWhere.String())
			nextQuery = buf2.String()
		} else {
			panic("not support")
		}
		var result2 *sql.Rows
		if conn != nil {
			result2, err = conn.QueryContext(context.Background(), nextQuery)
		} else {
			result2, err = this.supper.tx.Query(nextQuery)
		}
		if err != nil {
			panic(err)
		}

		defer result2.Close()

		if result2.Next() == false {
			panic("MustRun_ResultListWithPageInfo ")
		}
		err = result2.Scan(&totalMatch)
		if err != nil {
			panic(err)
		}

	}
	return list, totalMatch
}

func (this *KORM_test03UserInGroup_D_SelectObj) LeftJoin_User() *KORM_test03User_D_BeLeftJoin {
	node := this.joinNode.addLeftJoin("test03User_D", "User", "UserId", "Id")
	return &KORM_test03User_D_BeLeftJoin{
		joinNode:     node,
		bufWhere:     &this.bufWhere,
		argsWhereP:   &this.argsWhere,
		isLinkBeginP: &this.isLinkBegin,
		linkOpListP:  &this.linkOpList,
		orderByP:     &this.orderBy,
	}
}

type KORM_test03User_D_BeLeftJoin struct {
	joinNode     *korm_leftJoinNode
	bufWhere     *bytes.Buffer
	argsWhereP   *[]interface{}
	isLinkBeginP *bool
	linkOpListP  *[]string
	orderByP     *[]string
}
type KORM_Where_KORM_test03User_D_BeLeftJoin_Id struct {
	supper *KORM_test03User_D_BeLeftJoin
}

func (this *KORM_test03User_D_BeLeftJoin) Where_Id() *KORM_Where_KORM_test03User_D_BeLeftJoin_Id {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Id` ")
	*(this.isLinkBeginP) = false
	return &KORM_Where_KORM_test03User_D_BeLeftJoin_Id{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Id) Equal(Id string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Id) NotEqual(Id string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("!=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Id) GreaterOrEqual(Id string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString(">=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Id) Less(Id string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Id) LessOrEqual(Id string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Id) In(vList []string) *KORM_test03User_D_BeLeftJoin {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_BeLeftJoin_Name struct {
	supper *KORM_test03User_D_BeLeftJoin
}

func (this *KORM_test03User_D_BeLeftJoin) Where_Name() *KORM_Where_KORM_test03User_D_BeLeftJoin_Name {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Name` ")
	*(this.isLinkBeginP) = false
	return &KORM_Where_KORM_test03User_D_BeLeftJoin_Name{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Name) Equal(Name string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Name) NotEqual(Name string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("!=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Name) GreaterOrEqual(Name string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString(">=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Name) Less(Name string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Name) LessOrEqual(Name string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_Name) In(vList []string) *KORM_test03User_D_BeLeftJoin {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash struct {
	supper *KORM_test03User_D_BeLeftJoin
}

func (this *KORM_test03User_D_BeLeftJoin) Where_PasswordHash() *KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`PasswordHash` ")
	*(this.isLinkBeginP) = false
	return &KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash) Equal(PasswordHash string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash) NotEqual(PasswordHash string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("!=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash) GreaterOrEqual(PasswordHash string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString(">=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash) Less(PasswordHash string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash) LessOrEqual(PasswordHash string) *KORM_test03User_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), PasswordHash)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_PasswordHash) In(vList []string) *KORM_test03User_D_BeLeftJoin {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03User_D_BeLeftJoin_CreateTime struct {
	supper *KORM_test03User_D_BeLeftJoin
}

func (this *KORM_test03User_D_BeLeftJoin) Where_CreateTime() *KORM_Where_KORM_test03User_D_BeLeftJoin_CreateTime {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`CreateTime` ")
	*(this.isLinkBeginP) = false
	return &KORM_Where_KORM_test03User_D_BeLeftJoin_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_CreateTime) Equal(CreateTime time.Time) *KORM_test03User_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_CreateTime) NotEqual(CreateTime time.Time) *KORM_test03User_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test03User_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_CreateTime) Less(CreateTime time.Time) *KORM_test03User_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03User_D_BeLeftJoin_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test03User_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_test03User_D_BeLeftJoin) CondMultOpBegin_AND() *KORM_test03User_D_BeLeftJoin {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	(*this.linkOpListP) = append([]string{"AND"}, (*this.linkOpListP)...)
	*(this.isLinkBeginP) = true
	return this
}
func (this *KORM_test03User_D_BeLeftJoin) CondMultOpBegin_OR() *KORM_test03User_D_BeLeftJoin {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	(*this.linkOpListP) = append([]string{"OR"}, (*this.linkOpListP)...)
	*(this.isLinkBeginP) = true
	return this
}
func (this *KORM_test03User_D_BeLeftJoin) CondMultOpEnd() *KORM_test03User_D_BeLeftJoin {
	if *(this.isLinkBeginP) {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	(*this.linkOpListP) = (*this.linkOpListP)[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test03User_D_BeLeftJoin_OrderByObj struct {
	fieldName string
	supper    *KORM_test03User_D_BeLeftJoin
}

func (this *KORM_test03User_D_BeLeftJoin_OrderByObj) ASC() *KORM_test03User_D_BeLeftJoin {
	(*this.supper.orderByP) = append((*this.supper.orderByP), this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_test03User_D_BeLeftJoin_OrderByObj) DESC() *KORM_test03User_D_BeLeftJoin {
	(*this.supper.orderByP) = append((*this.supper.orderByP), this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_test03User_D_BeLeftJoin) OrderBy_Id() *KORM_test03User_D_BeLeftJoin_OrderByObj {
	return &KORM_test03User_D_BeLeftJoin_OrderByObj{
		fieldName: "Id",
		supper:    this,
	}
}
func (this *KORM_test03User_D_BeLeftJoin) OrderBy_Name() *KORM_test03User_D_BeLeftJoin_OrderByObj {
	return &KORM_test03User_D_BeLeftJoin_OrderByObj{
		fieldName: "Name",
		supper:    this,
	}
}
func (this *KORM_test03User_D_BeLeftJoin) OrderBy_PasswordHash() *KORM_test03User_D_BeLeftJoin_OrderByObj {
	return &KORM_test03User_D_BeLeftJoin_OrderByObj{
		fieldName: "PasswordHash",
		supper:    this,
	}
}
func (this *KORM_test03User_D_BeLeftJoin) OrderBy_CreateTime() *KORM_test03User_D_BeLeftJoin_OrderByObj {
	return &KORM_test03User_D_BeLeftJoin_OrderByObj{
		fieldName: "CreateTime",
		supper:    this,
	}
}
func (this *KORM_test03UserInGroup_D_SelectObj) LeftJoin_Group() *KORM_test03Group_D_BeLeftJoin {
	node := this.joinNode.addLeftJoin("test03Group_D", "Group", "GroupId", "Id")
	return &KORM_test03Group_D_BeLeftJoin{
		joinNode:     node,
		bufWhere:     &this.bufWhere,
		argsWhereP:   &this.argsWhere,
		isLinkBeginP: &this.isLinkBegin,
		linkOpListP:  &this.linkOpList,
		orderByP:     &this.orderBy,
	}
}

type KORM_test03Group_D_BeLeftJoin struct {
	joinNode     *korm_leftJoinNode
	bufWhere     *bytes.Buffer
	argsWhereP   *[]interface{}
	isLinkBeginP *bool
	linkOpListP  *[]string
	orderByP     *[]string
}
type KORM_Where_KORM_test03Group_D_BeLeftJoin_Id struct {
	supper *KORM_test03Group_D_BeLeftJoin
}

func (this *KORM_test03Group_D_BeLeftJoin) Where_Id() *KORM_Where_KORM_test03Group_D_BeLeftJoin_Id {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Id` ")
	*(this.isLinkBeginP) = false
	return &KORM_Where_KORM_test03Group_D_BeLeftJoin_Id{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Id) Equal(Id uint64) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Id) NotEqual(Id uint64) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("!=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Id) GreaterOrEqual(Id uint64) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString(">=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Id) Less(Id uint64) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Id) LessOrEqual(Id uint64) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Id)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Id) In(vList []uint64) *KORM_test03Group_D_BeLeftJoin {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03Group_D_BeLeftJoin_Name struct {
	supper *KORM_test03Group_D_BeLeftJoin
}

func (this *KORM_test03Group_D_BeLeftJoin) Where_Name() *KORM_Where_KORM_test03Group_D_BeLeftJoin_Name {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Name` ")
	*(this.isLinkBeginP) = false
	return &KORM_Where_KORM_test03Group_D_BeLeftJoin_Name{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Name) Equal(Name string) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Name) NotEqual(Name string) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("!=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Name) GreaterOrEqual(Name string) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString(">=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Name) Less(Name string) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Name) LessOrEqual(Name string) *KORM_test03Group_D_BeLeftJoin {
	this.supper.bufWhere.WriteString("<=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), Name)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_Name) In(vList []string) *KORM_test03Group_D_BeLeftJoin {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03Group_D_BeLeftJoin_CreateTime struct {
	supper *KORM_test03Group_D_BeLeftJoin
}

func (this *KORM_test03Group_D_BeLeftJoin) Where_CreateTime() *KORM_Where_KORM_test03Group_D_BeLeftJoin_CreateTime {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`CreateTime` ")
	*(this.isLinkBeginP) = false
	return &KORM_Where_KORM_test03Group_D_BeLeftJoin_CreateTime{supper: this}
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_CreateTime) Equal(CreateTime time.Time) *KORM_test03Group_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_CreateTime) NotEqual(CreateTime time.Time) *KORM_test03Group_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("!=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_CreateTime) GreaterOrEqual(CreateTime time.Time) *KORM_test03Group_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString(">=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_CreateTime) Less(CreateTime time.Time) *KORM_test03Group_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_Where_KORM_test03Group_D_BeLeftJoin_CreateTime) LessOrEqual(CreateTime time.Time) *KORM_test03Group_D_BeLeftJoin {
	vCreateTime := CreateTime.UTC().Format(time.RFC3339Nano)
	this.supper.bufWhere.WriteString("<=? ")
	*(this.supper.argsWhereP) = append(*(this.supper.argsWhereP), vCreateTime)
	return this.supper
}
func (this *KORM_test03Group_D_BeLeftJoin) CondMultOpBegin_AND() *KORM_test03Group_D_BeLeftJoin {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	(*this.linkOpListP) = append([]string{"AND"}, (*this.linkOpListP)...)
	*(this.isLinkBeginP) = true
	return this
}
func (this *KORM_test03Group_D_BeLeftJoin) CondMultOpBegin_OR() *KORM_test03Group_D_BeLeftJoin {
	if this.bufWhere.Len() > 0 {
		if *(this.isLinkBeginP) == false {
			this.bufWhere.WriteString((*this.linkOpListP)[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	(*this.linkOpListP) = append([]string{"OR"}, (*this.linkOpListP)...)
	*(this.isLinkBeginP) = true
	return this
}
func (this *KORM_test03Group_D_BeLeftJoin) CondMultOpEnd() *KORM_test03Group_D_BeLeftJoin {
	if *(this.isLinkBeginP) {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	(*this.linkOpListP) = (*this.linkOpListP)[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test03Group_D_BeLeftJoin_OrderByObj struct {
	fieldName string
	supper    *KORM_test03Group_D_BeLeftJoin
}

func (this *KORM_test03Group_D_BeLeftJoin_OrderByObj) ASC() *KORM_test03Group_D_BeLeftJoin {
	(*this.supper.orderByP) = append((*this.supper.orderByP), this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_test03Group_D_BeLeftJoin_OrderByObj) DESC() *KORM_test03Group_D_BeLeftJoin {
	(*this.supper.orderByP) = append((*this.supper.orderByP), this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_test03Group_D_BeLeftJoin) OrderBy_Name() *KORM_test03Group_D_BeLeftJoin_OrderByObj {
	return &KORM_test03Group_D_BeLeftJoin_OrderByObj{
		fieldName: "Name",
		supper:    this,
	}
}
func (this *KORM_test03Group_D_BeLeftJoin) OrderBy_CreateTime() *KORM_test03Group_D_BeLeftJoin_OrderByObj {
	return &KORM_test03Group_D_BeLeftJoin_OrderByObj{
		fieldName: "CreateTime",
		supper:    this,
	}
}

type KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId struct {
	supper *KORM_test03UserInGroup_D_SelectObj
}

func (this *KORM_test03UserInGroup_D_SelectObj) Where_UserId() *KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`UserId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId{supper: this}
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId) Equal(UserId string) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId) NotEqual(UserId string) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId) GreaterOrEqual(UserId string) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId) Less(UserId string) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId) LessOrEqual(UserId string) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_UserId) In(vList []string) *KORM_test03UserInGroup_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId struct {
	supper *KORM_test03UserInGroup_D_SelectObj
}

func (this *KORM_test03UserInGroup_D_SelectObj) Where_GroupId() *KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`GroupId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId{supper: this}
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId) Equal(GroupId uint64) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId) NotEqual(GroupId uint64) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId) GreaterOrEqual(GroupId uint64) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId) Less(GroupId uint64) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId) LessOrEqual(GroupId uint64) *KORM_test03UserInGroup_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_SelectObj_GroupId) In(vList []uint64) *KORM_test03UserInGroup_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test03UserInGroup_D_SelectObj) CondMultOpBegin_AND() *KORM_test03UserInGroup_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03UserInGroup_D_SelectObj) CondMultOpBegin_OR() *KORM_test03UserInGroup_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03UserInGroup_D_SelectObj) CondMultOpEnd() *KORM_test03UserInGroup_D_SelectObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

// Update test03UserInGroup_D
type KORM_test03UserInGroup_D_UpdateObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test03UserInGroup_D) Update() *KORM_test03UserInGroup_D_UpdateObj {
	return &KORM_test03UserInGroup_D_UpdateObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test03UserInGroup_D) MustUpdateBy_UserId(info test03UserInGroup_D) (rowsAffected int64) {
	rowsAffected = this.Update().Where_UserId().Equal(info.UserId).Set_GroupId(info.GroupId).MustRun()
	return rowsAffected
}
func (this *KORM_test03UserInGroup_D_UpdateObj) Set_GroupId(GroupId uint64) *KORM_test03UserInGroup_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`GroupId` = ? ")
	this.argsSet = append(this.argsSet, GroupId)
	return this
}
func (this *KORM_test03UserInGroup_D_UpdateObj) MustRun() (RowsAffected int64) {
	if len(this.argsSet) == 0 {
		panic("len(this.argsSet) == 0")
	}
	var buf2 bytes.Buffer
	buf2.WriteString("UPDATE `test03UserInGroup_D` ")
	buf2.WriteString(this.bufSet.String())
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error

	result, err = this.supper.ExecRaw(buf2.String(), append(this.argsSet, this.argsWhere...)...)
	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId struct {
	supper *KORM_test03UserInGroup_D_UpdateObj
}

func (this *KORM_test03UserInGroup_D_UpdateObj) Where_UserId() *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`UserId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId{supper: this}
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId) Equal(UserId string) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId) NotEqual(UserId string) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId) GreaterOrEqual(UserId string) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId) Less(UserId string) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId) LessOrEqual(UserId string) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_UserId) In(vList []string) *KORM_test03UserInGroup_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId struct {
	supper *KORM_test03UserInGroup_D_UpdateObj
}

func (this *KORM_test03UserInGroup_D_UpdateObj) Where_GroupId() *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`GroupId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId{supper: this}
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId) Equal(GroupId uint64) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId) NotEqual(GroupId uint64) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId) GreaterOrEqual(GroupId uint64) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId) Less(GroupId uint64) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId) LessOrEqual(GroupId uint64) *KORM_test03UserInGroup_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_UpdateObj_GroupId) In(vList []uint64) *KORM_test03UserInGroup_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test03UserInGroup_D_UpdateObj) CondMultOpBegin_AND() *KORM_test03UserInGroup_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03UserInGroup_D_UpdateObj) CondMultOpBegin_OR() *KORM_test03UserInGroup_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03UserInGroup_D_UpdateObj) CondMultOpEnd() *KORM_test03UserInGroup_D_UpdateObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test03UserInGroup_D_DeleteObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test03UserInGroup_D) Delete() *KORM_test03UserInGroup_D_DeleteObj {
	return &KORM_test03UserInGroup_D_DeleteObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test03UserInGroup_D_DeleteObj) MustRun() (RowsAffected int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("DELETE FROM test03UserInGroup_D ")
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error
	result, err = this.supper.ExecRaw(buf2.String(), this.argsWhere...)

	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId struct {
	supper *KORM_test03UserInGroup_D_DeleteObj
}

func (this *KORM_test03UserInGroup_D_DeleteObj) Where_UserId() *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`UserId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId{supper: this}
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId) Equal(UserId string) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId) NotEqual(UserId string) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId) GreaterOrEqual(UserId string) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId) Less(UserId string) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId) LessOrEqual(UserId string) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, UserId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_UserId) In(vList []string) *KORM_test03UserInGroup_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId struct {
	supper *KORM_test03UserInGroup_D_DeleteObj
}

func (this *KORM_test03UserInGroup_D_DeleteObj) Where_GroupId() *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`GroupId` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId{supper: this}
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId) Equal(GroupId uint64) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId) NotEqual(GroupId uint64) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId) GreaterOrEqual(GroupId uint64) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId) Less(GroupId uint64) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId) LessOrEqual(GroupId uint64) *KORM_test03UserInGroup_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, GroupId)
	return this.supper
}
func (this *KORM_Where_KORM_test03UserInGroup_D_DeleteObj_GroupId) In(vList []uint64) *KORM_test03UserInGroup_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test03UserInGroup_D_DeleteObj) CondMultOpBegin_AND() *KORM_test03UserInGroup_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03UserInGroup_D_DeleteObj) CondMultOpBegin_OR() *KORM_test03UserInGroup_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test03UserInGroup_D_DeleteObj) CondMultOpEnd() *KORM_test03UserInGroup_D_DeleteObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test04User_D struct {
	supper *OrmAll
}

func (this *OrmAll) test04User_D() *KORM_test04User_D {
	return &KORM_test04User_D{supper: this}
}
func korm_fillSelectFieldNameList_test04User_D(joinNode *korm_leftJoinNode) {
	joinNode.selectFieldNameList = []string{"Id", "Name", "Key"}
	for _, sub := range joinNode.thisLevelJoinList {
		switch sub.fieldName {
		default:
			panic("korm_fillSelectFieldNameList_test04User_D" + strconv.Quote(sub.fieldName))
		}
	}
}
func (this *KORM_test04User_D) MustInsert(info test04User_D) {
	var err error
	_, err = this.supper.ExecRaw("INSERT INTO `test04User_D`(`Id` ,`Name` ,`Key` ) VALUES(?,?,?)", info.Id, info.Name, info.Key)
	if err != nil {
		panic(err)
	}

	return
}
func (this *KORM_test04User_D) MustSet(info test04User_D) {
	var err error
	_, err = this.supper.ExecRaw("REPLACE INTO `test04User_D`(`Id` ,`Name` ,`Key` ) VALUES(?,?,?)", info.Id, info.Name, info.Key)
	if err != nil {
		panic(err)
	}

	return
}

// Select test04User_D
type KORM_test04User_D_SelectObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
	joinNode    *korm_leftJoinNode
	rootInfo    korm_leftJoinRootInfo
}

func (this *KORM_test04User_D) Select() *KORM_test04User_D_SelectObj {
	one := &KORM_test04User_D_SelectObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
	one.joinNode = &korm_leftJoinNode{tableName: "_0"}
	one.joinNode.root = &one.rootInfo
	one.rootInfo.tableNameIdx = 1
	return one
}

type KORM_test04User_D_SelectObj_OrderByObj struct {
	fieldName string
	supper    *KORM_test04User_D_SelectObj
}

func (this *KORM_test04User_D_SelectObj_OrderByObj) ASC() *KORM_test04User_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` ASC ")
	return this.supper
}
func (this *KORM_test04User_D_SelectObj_OrderByObj) DESC() *KORM_test04User_D_SelectObj {
	this.supper.orderBy = append(this.supper.orderBy, this.supper.joinNode.tableName+".`"+this.fieldName+"` DESC ")
	return this.supper
}
func (this *KORM_test04User_D_SelectObj) OrderBy_Id() *KORM_test04User_D_SelectObj_OrderByObj {
	return &KORM_test04User_D_SelectObj_OrderByObj{
		fieldName: "Id",
		supper:    this,
	}
}
func (this *KORM_test04User_D_SelectObj) OrderBy_Name() *KORM_test04User_D_SelectObj_OrderByObj {
	return &KORM_test04User_D_SelectObj_OrderByObj{
		fieldName: "Name",
		supper:    this,
	}
}
func (this *KORM_test04User_D_SelectObj) OrderBy_Key() *KORM_test04User_D_SelectObj_OrderByObj {
	return &KORM_test04User_D_SelectObj_OrderByObj{
		fieldName: "Key",
		supper:    this,
	}
}

func (this *KORM_test04User_D_SelectObj) LimitOffset(limit int, offset int) *KORM_test04User_D_SelectObj {
	this.limit = limit
	this.offset = offset
	return this
}

// pageSize: [1, n)
// pageNo:   [1,n)
func (this *KORM_test04User_D_SelectObj) SetPageLimit(pageSize int, pageNo int) *KORM_test04User_D_SelectObj {
	if pageSize <= 0 || pageNo <= 0 {
		panic("KORM_test04User_D_SelectObj SetPageLimit error param")
	}
	this.limit = pageSize
	this.offset = pageSize * (pageNo - 1)
	return this
}

func (this *KORM_test04User_D_SelectObj) MustRun_Count() (cnt int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test04User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test04User_D_SelectObj not found.")
	}
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt
}

func (this *KORM_test04User_D_SelectObj) MustRun_Exist() (exist bool) {
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT COUNT(1) FROM `test04User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	buf2.WriteString(" LIMIT 1 ")

	result, err := this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		panic("KORM_test04User_D_SelectObj not found.")
	}
	var cnt int64
	err = result.Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt > 0
}

func (this *KORM_test04User_D_SelectObj) MustRun_ResultOne2() (info test04User_D, ok bool) {
	this.limit = 1
	korm_fillSelectFieldNameList_test04User_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test04User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()
	if result.Next() == false {
		return info, false
	}
	resp := korm_test04User_D_scan(this.joinNode, &info)
	err = result.Scan(resp.argList...)
	if err != nil {
		panic(err)
	}

	for idx, a := range resp.argList {
		v := a.(*sql.NullString).String
		if v == "" {
			continue
		}
		err = resp.argParseFn[idx](v)
		if err != nil {
			panic(err)
		}

	}

	return info, true
}

func (this *KORM_test04User_D_SelectObj) MustRun_ResultOne() (info test04User_D) {
	info, _ = this.MustRun_ResultOne2()
	return info
}
func (this *KORM_test04User_D_SelectObj) MustRun_ResultList() (list []test04User_D) {
	korm_fillSelectFieldNameList_test04User_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test04User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test04User_D
		korm_fillSelectFieldNameList_test04User_D(this.joinNode)
		resp := korm_test04User_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	return list
}
func (this *KORM_test04User_D_SelectObj) MustRun_ResultMap() (m map[int]test04User_D) {
	korm_fillSelectFieldNameList_test04User_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test04User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}

	var result *sql.Rows
	var err error

	result, err = this.supper.ExecRawQuery(buf2.String(), this.argsWhere...)
	if err != nil {
		panic(err)
	}

	defer result.Close()

	for result.Next() {
		var info test04User_D
		korm_fillSelectFieldNameList_test04User_D(this.joinNode)
		resp := korm_test04User_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		if m == nil {
			m = map[int]test04User_D{}
		}
		m[info.Id] = info

	}
	return m
}
func (this *KORM_test04User_D_SelectObj) MustRun_ResultListWithTotalMatch() (list []test04User_D, totalMatch int64) {
	var err error
	korm_fillSelectFieldNameList_test04User_D(this.joinNode)
	var buf2 bytes.Buffer
	buf2.WriteString("SELECT ")
	if this.supper.mode == korm.InitModeMysql {
		buf2.WriteString("SQL_CALC_FOUND_ROWS  ")
	}
	this.joinNode.fillSelect(&buf2, true)
	buf2.WriteString(" ")
	buf2.WriteString("FROM `test04User_D` " + this.joinNode.tableName + " ")
	buf2.WriteString(this.rootInfo.bufLeftJoin.String())
	buf2.WriteString(this.bufWhere.String())
	if len(this.orderBy) > 0 {
		buf2.WriteString("ORDER BY " + strings.Join(this.orderBy, ",") + " ")
	}
	if this.limit != 0 {
		buf2.WriteString("LIMIT " + strconv.Itoa(this.limit) + " ")
	}
	if this.offset != 0 {
		buf2.WriteString("OFFSET " + strconv.Itoa(this.offset) + " ")
	}
	var conn *sql.Conn
	var result *sql.Rows
	if this.supper.db != nil {
		var err error
		conn, err = this.supper.db.Conn(context.Background())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		result, err = conn.QueryContext(context.Background(), buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		result, err = this.supper.tx.Query(buf2.String(), this.argsWhere...)
		if err != nil {
			panic(err)
		}
	}

	defer result.Close()
	for result.Next() {
		var info test04User_D
		korm_fillSelectFieldNameList_test04User_D(this.joinNode)
		resp := korm_test04User_D_scan(this.joinNode, &info)
		err = result.Scan(resp.argList...)
		if err != nil {
			panic(err)
		}

		for idx, a := range resp.argList {
			v := a.(*sql.NullString).String
			if v == "" {
				continue
			}
			err = resp.argParseFn[idx](v)
			if err != nil {
				panic(err)
			}

		}

		list = append(list, info)
	}
	result.Close()
	if this.limit > 0 {
		nextQuery := ""
		if this.supper.mode == korm.InitModeMysql {
			nextQuery = "select FOUND_ROWS()"
		} else if this.supper.mode == korm.InitModeSqlite {
			buf2.Reset()
			buf2.WriteString("SELECT COUNT(1) ")
			buf2.WriteString("FROM `test04User_D` " + this.joinNode.tableName + " ")
			buf2.WriteString(this.rootInfo.bufLeftJoin.String())
			buf2.WriteString(this.bufWhere.String())
			nextQuery = buf2.String()
		} else {
			panic("not support")
		}
		var result2 *sql.Rows
		if conn != nil {
			result2, err = conn.QueryContext(context.Background(), nextQuery)
		} else {
			result2, err = this.supper.tx.Query(nextQuery)
		}
		if err != nil {
			panic(err)
		}

		defer result2.Close()

		if result2.Next() == false {
			panic("MustRun_ResultListWithPageInfo ")
		}
		err = result2.Scan(&totalMatch)
		if err != nil {
			panic(err)
		}

	}
	return list, totalMatch
}

type KORM_Where_KORM_test04User_D_SelectObj_Id struct {
	supper *KORM_test04User_D_SelectObj
}

func (this *KORM_test04User_D_SelectObj) Where_Id() *KORM_Where_KORM_test04User_D_SelectObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_SelectObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Id) Equal(Id int) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Id) NotEqual(Id int) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Id) GreaterOrEqual(Id int) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Id) Less(Id int) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Id) LessOrEqual(Id int) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Id) In(vList []int) *KORM_test04User_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test04User_D_SelectObj_Name struct {
	supper *KORM_test04User_D_SelectObj
}

func (this *KORM_test04User_D_SelectObj) Where_Name() *KORM_Where_KORM_test04User_D_SelectObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_SelectObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Name) Equal(Name string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Name) NotEqual(Name string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Name) GreaterOrEqual(Name string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Name) Less(Name string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Name) LessOrEqual(Name string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Name) In(vList []string) *KORM_test04User_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test04User_D_SelectObj_Key struct {
	supper *KORM_test04User_D_SelectObj
}

func (this *KORM_test04User_D_SelectObj) Where_Key() *KORM_Where_KORM_test04User_D_SelectObj_Key {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`" + this.joinNode.tableName + "`.")
	this.bufWhere.WriteString("`Key` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_SelectObj_Key{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Key) Equal(Key string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Key) NotEqual(Key string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Key) GreaterOrEqual(Key string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Key) Less(Key string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Key) LessOrEqual(Key string) *KORM_test04User_D_SelectObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_SelectObj_Key) In(vList []string) *KORM_test04User_D_SelectObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test04User_D_SelectObj) CondMultOpBegin_AND() *KORM_test04User_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test04User_D_SelectObj) CondMultOpBegin_OR() *KORM_test04User_D_SelectObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test04User_D_SelectObj) CondMultOpEnd() *KORM_test04User_D_SelectObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

// Update test04User_D
type KORM_test04User_D_UpdateObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	limit       int
	offset      int
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test04User_D) Update() *KORM_test04User_D_UpdateObj {
	return &KORM_test04User_D_UpdateObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test04User_D) MustUpdateBy_Id(info test04User_D) (rowsAffected int64) {
	rowsAffected = this.Update().Where_Id().Equal(info.Id).Set_Name(info.Name).Set_Key(info.Key).MustRun()
	return rowsAffected
}
func (this *KORM_test04User_D_UpdateObj) Set_Name(Name string) *KORM_test04User_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`Name` = ? ")
	this.argsSet = append(this.argsSet, Name)
	return this
}
func (this *KORM_test04User_D_UpdateObj) Set_Key(Key string) *KORM_test04User_D_UpdateObj {
	if len(this.argsSet) > 0 {
		this.bufSet.WriteString(", ")
	} else {
		this.bufSet.WriteString("SET ")
	}
	this.bufSet.WriteString("`Key` = ? ")
	this.argsSet = append(this.argsSet, Key)
	return this
}
func (this *KORM_test04User_D_UpdateObj) MustRun() (RowsAffected int64) {
	if len(this.argsSet) == 0 {
		panic("len(this.argsSet) == 0")
	}
	var buf2 bytes.Buffer
	buf2.WriteString("UPDATE `test04User_D` ")
	buf2.WriteString(this.bufSet.String())
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error

	result, err = this.supper.ExecRaw(buf2.String(), append(this.argsSet, this.argsWhere...)...)
	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test04User_D_UpdateObj_Id struct {
	supper *KORM_test04User_D_UpdateObj
}

func (this *KORM_test04User_D_UpdateObj) Where_Id() *KORM_Where_KORM_test04User_D_UpdateObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_UpdateObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Id) Equal(Id int) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Id) NotEqual(Id int) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Id) GreaterOrEqual(Id int) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Id) Less(Id int) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Id) LessOrEqual(Id int) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Id) In(vList []int) *KORM_test04User_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test04User_D_UpdateObj_Name struct {
	supper *KORM_test04User_D_UpdateObj
}

func (this *KORM_test04User_D_UpdateObj) Where_Name() *KORM_Where_KORM_test04User_D_UpdateObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_UpdateObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Name) Equal(Name string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Name) NotEqual(Name string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Name) GreaterOrEqual(Name string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Name) Less(Name string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Name) LessOrEqual(Name string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Name) In(vList []string) *KORM_test04User_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test04User_D_UpdateObj_Key struct {
	supper *KORM_test04User_D_UpdateObj
}

func (this *KORM_test04User_D_UpdateObj) Where_Key() *KORM_Where_KORM_test04User_D_UpdateObj_Key {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Key` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_UpdateObj_Key{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Key) Equal(Key string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Key) NotEqual(Key string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Key) GreaterOrEqual(Key string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Key) Less(Key string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Key) LessOrEqual(Key string) *KORM_test04User_D_UpdateObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_UpdateObj_Key) In(vList []string) *KORM_test04User_D_UpdateObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test04User_D_UpdateObj) CondMultOpBegin_AND() *KORM_test04User_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test04User_D_UpdateObj) CondMultOpBegin_OR() *KORM_test04User_D_UpdateObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test04User_D_UpdateObj) CondMultOpEnd() *KORM_test04User_D_UpdateObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}

type KORM_test04User_D_DeleteObj struct {
	supper      *OrmAll
	bufWhere    bytes.Buffer
	argsWhere   []interface{}
	bufSet      bytes.Buffer
	argsSet     []interface{}
	orderBy     []string
	linkOpList  []string
	isLinkBegin bool
}

func (this *KORM_test04User_D) Delete() *KORM_test04User_D_DeleteObj {
	return &KORM_test04User_D_DeleteObj{supper: this.supper, linkOpList: []string{"AND"}, isLinkBegin: true}
}
func (this *KORM_test04User_D_DeleteObj) MustRun() (RowsAffected int64) {
	var buf2 bytes.Buffer
	buf2.WriteString("DELETE FROM test04User_D ")
	buf2.WriteString(this.bufWhere.String())
	var result sql.Result
	var err error
	result, err = this.supper.ExecRaw(buf2.String(), this.argsWhere...)

	if err != nil {
		panic(err)
	}

	RowsAffected, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return RowsAffected
}

type KORM_Where_KORM_test04User_D_DeleteObj_Id struct {
	supper *KORM_test04User_D_DeleteObj
}

func (this *KORM_test04User_D_DeleteObj) Where_Id() *KORM_Where_KORM_test04User_D_DeleteObj_Id {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Id` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_DeleteObj_Id{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Id) Equal(Id int) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Id) NotEqual(Id int) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Id) GreaterOrEqual(Id int) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Id) Less(Id int) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Id) LessOrEqual(Id int) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Id)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Id) In(vList []int) *KORM_test04User_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test04User_D_DeleteObj_Name struct {
	supper *KORM_test04User_D_DeleteObj
}

func (this *KORM_test04User_D_DeleteObj) Where_Name() *KORM_Where_KORM_test04User_D_DeleteObj_Name {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Name` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_DeleteObj_Name{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Name) Equal(Name string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Name) NotEqual(Name string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Name) GreaterOrEqual(Name string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Name) Less(Name string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Name) LessOrEqual(Name string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Name)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Name) In(vList []string) *KORM_test04User_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}

type KORM_Where_KORM_test04User_D_DeleteObj_Key struct {
	supper *KORM_test04User_D_DeleteObj
}

func (this *KORM_test04User_D_DeleteObj) Where_Key() *KORM_Where_KORM_test04User_D_DeleteObj_Key {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("`Key` ")
	this.isLinkBegin = false
	return &KORM_Where_KORM_test04User_D_DeleteObj_Key{supper: this}
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Key) Equal(Key string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Key) NotEqual(Key string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("!=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Key) GreaterOrEqual(Key string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString(">=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Key) Less(Key string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Key) LessOrEqual(Key string) *KORM_test04User_D_DeleteObj {
	this.supper.bufWhere.WriteString("<=? ")
	this.supper.argsWhere = append(this.supper.argsWhere, Key)
	return this.supper
}
func (this *KORM_Where_KORM_test04User_D_DeleteObj_Key) In(vList []string) *KORM_test04User_D_DeleteObj {
	if len(vList) == 0 {
		this.supper.bufWhere.WriteString("= '' AND 0 ") // 什么都不存在, 直接返回
		return this.supper
	}
	this.supper.bufWhere.WriteString("IN (")
	for idx, v := range vList {
		if idx > 0 {
			this.supper.bufWhere.WriteString(", ")
		}
		this.supper.bufWhere.WriteString("?")
		this.supper.argsWhere = append(this.supper.argsWhere, v)
	}
	this.supper.bufWhere.WriteString(") ")
	return this.supper
}
func (this *KORM_test04User_D_DeleteObj) CondMultOpBegin_AND() *KORM_test04User_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"AND"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test04User_D_DeleteObj) CondMultOpBegin_OR() *KORM_test04User_D_DeleteObj {
	if this.bufWhere.Len() > 0 {
		if this.isLinkBegin == false {
			this.bufWhere.WriteString(this.linkOpList[0] + " ")
		}
	} else {
		this.bufWhere.WriteString("WHERE ")
	}
	this.bufWhere.WriteString("( ")
	this.linkOpList = append([]string{"OR"}, this.linkOpList...)
	this.isLinkBegin = true
	return this
}
func (this *KORM_test04User_D_DeleteObj) CondMultOpEnd() *KORM_test04User_D_DeleteObj {
	if this.isLinkBegin {
		panic("() is not allowed in sql statement") // bad sql: SELECT * FROM u where ()
	}
	this.linkOpList = this.linkOpList[1:]
	this.bufWhere.WriteString(") ")
	return this
}
