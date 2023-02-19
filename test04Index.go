package korm_example

type test04User_D struct {
	Id   int
	Name string `korm:"index:Name;index:Id,Name"`
	Key  string `korm:"index:Key"`
}
