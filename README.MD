# go-cli

#### 运行
1. `cd ./deploy && docker-compose build && docker-compose up -d`
2. 包含`mysql redis fastdfs nginx`镜像

#### 介绍
* 涉及:`golang` `gin` `gorm` `mysql` `redis` `viper` `websocket` `crons` `base64Captcha` `fastdfs` `docker-compose` `ratelimit` 


#### 架构
* `config`:配置文件存放
* `dao`:数据库连接  
* `model/dto`:结构体模型及其方法
* `model/entity`:请求参数模型
* `model/vo`：视图模型，即接口返回模型
* `routers`:路由及中间件
* `service`:业务逻辑
* `utils`:工具类包含如下  

    1. 阿里云短信/验证码生成
    2. 跨域中间件 
    3. JWTtoken生成/解析/验证/刷新
    4. 配置文件读取/解析/监听
    5. MD5密码
    6. 基础响应格式/可带总数页数状态码等
    7. 自定义的响应状态码及响应提示
    8. 定时器任务
    9. 常用工具函数(数据去重/移除指定数据/返回当天零点等,自行查看自行补充哦)
    10. tripleDES加解密
    11. 特定类型验证
    12. 图片验证码生成/验证  
    13. websocket长连接处理(含token)
    14. 公共常用请求参数验证，contxtverify.go
    15. 敏感词过滤
    16. `ratelimit`过载保护中间件


#### TIPS     

* `VScode`扩展安装
1. `$env:GO111MODULE="on"`
2. `$env:GOPROXY="https://goproxy.io"`
3. `go env -w GOPROXY=https://goproxy.io,direct`
4. `go env -w GOPRIVATE=*.corp.example.com`
5. 扩展管理--代理设置（默认override，改成on）
6. 重启，安装成功
7. 内存逃逸分析：`go build -gcflags '-m'`


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





