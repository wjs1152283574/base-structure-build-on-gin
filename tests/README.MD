##### 单元测试
* 切到目录:`cd tests`
* 全部测试: `go test -v`
* 测试单个文件: `go test -v date_util_test.go`,其中`date_util_test.go`为需要运行的测试文件
* 测试指定函数: `go test -v -run TestLocalNow`,其中`TestLocalNow`为需要测试的函数
* `test -v `: 单元测试中的日志只有不通过时才会输出，加上`test -v`或`-test.v`来输出测试具体日志
