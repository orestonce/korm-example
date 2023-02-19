# korm-example

## korm的例子 使用方法
1. 解决依赖
```bash
go mod tidy
```
2. 预编译korm提供的函数
```bash
go run before-build/main.go
```
3. 运行c1/main.go, 测试mysql和sqlite的适配效果
```bash
go run c1/main.go
```