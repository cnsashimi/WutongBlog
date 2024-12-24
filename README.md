# 梧桐树下blog程序 Go版

### 简介
原来就是一个写php 和 java的马楼，尝试一下go，所以有了这个。
### 更新历史
-  **2024/12/21**   V0.6.0 上线 https://www.imzixuan.cn/ 
-  **2024/12/20**   V0.5.8 前后台规整
-  **2024/12/11**   V0.5.6 使用Redis对一些重复读取数据直接缓冲 
-  **2024/12/10**   V0.5.5 基本增删改实现，首页分页。
-  **2024/11/29**   V0.5.3 上传图片接口规整，目前还不具备完整的blog功能。 权当学习go语言中，佛系更新。
-  **2024/11/19**   V0.5beta helloworld 大概的框架结构确认， 目前不具备任何实用功能。 
 
### 软件架构
#### 后端程序环境： go1.22.8 + fiber v2.52.5 + mysql(版本>5.7.8 ) + redis
#### blog前端: Logbook Bootstrap  https://themefisher.com/products/logbook-bootstrap   MIT license<BR>
#### 后台管理页前端: 修改于webman-admin https://github.com/webman-php/admin    MIT license<BR>   Tags标签插件: https://gitee.com/cshaptx4869/input-tag
#### layui https://gitee.com/layui/layui  MIT license<BR>
### 安装教程

##### 1.  mysql导入mysqldb.sql
##### 2.  修改/config.yml里相关设置
##### 3.  运行： go run main.go
##### 4.  编译： go build main.go

### 使用说明
http://127.0.0.1:3000/   

### 后台
http://127.0.0.1:3000/login    用户名 wutong  密码 123456


 