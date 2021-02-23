# go-cli

#### 介绍
golang gin gorm mysql redis

#### 架构
* `config`:配置文件存放
* `dao`:数据库连接   
* `database`:结构体模型
* `routers`:路由
* `service`:业务逻辑
* `utils`:工具类(跨域/JWT/配置文件解析/MD5/验证器/websocket)


#### TIPS     

* `VScode`扩展安装
1. `$env:GO111MODULE="on"`
2. `$env:GOPROXY="https://goproxy.io"`
3. `go env -w GOPROXY=https://goproxy.io,direct`
4. `go env -w GOPRIVATE=*.corp.example.com`
5. 扩展管理--代理设置（默认override，改成on）
6. 重启，安装成功


#### Commit 规范

###### 格式如下
* 例：`fead(type)`:本次提交概述
* `type`: 本次 commit 的类型，诸如 bugfix docs style 等，参考如下:  

    * `fead`：添加新功能
    * `fix`：修补缺陷
    * `docs`：修改文档
    * `style`：修改格式
    * `refactor`：重构
    * `perf`：优化
    * `test`：增加测试
    * `chore`：构建过程或辅助工具的变动
    * `revert`：回滚到上一个版本

* `scope`: 本次 `commit` 波及的范围
* `subject`: 简明扼要的阐述下本次 `commit` 的主旨，在原文中特意强调了几点：

    1. 使用祈使句，是不是很熟悉又陌生的一个词
    2. 首字母不要大写
    3. 结尾无需添加标点





