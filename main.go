package korm_example

type DownloadCache_D struct {
	Url     string `korm:"primary"`
	Content []byte
}

type ActorInfo_D struct {
	Name string `korm:"primary"`
	Uid  uint16 `korm:"primary"`
}
