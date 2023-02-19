package main

import (
	"github.com/orestonce/korm"
)

func main() {
	korm.MustCreateCode(korm.MustCreateCode_Req{
		ModelPkgDir:      "./",
		ModelPkgFullPath: "korm_example",
		ModelNameList:    []string{"DownloadCache_D", "test01Crud_D", "test02MultiplePk_D", "test03User_D", "test03Group_D", "test03UserInGroup_D", "test04User_D"},
		OutputFileName:   `./generated_korm.go`,
		GenMustFn:        true,
	})
}
