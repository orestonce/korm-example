package korm_example

import (
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/go-sql-driver/mysql"
)

type DownloadCache_D struct {
	Url     string `korm:"primary"`
	Content []byte
}

type ActorInfo_D struct {
	Name string `korm:"primary"`
	Uid  uint16 `korm:"primary"`
}
