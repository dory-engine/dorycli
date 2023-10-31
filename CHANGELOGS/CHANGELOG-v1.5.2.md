# dorycli changelog v1.5.2

**新特性:**

- dorycli install print 命令默认不安装代码仓库、镜像仓库、制品仓库、代码扫描仓库

- dorycli install print 命令支持--full参数，full参数表示完整安装代码仓库、镜像仓库、制品仓库、代码扫描仓库

- dorycli install 不再自动下载trivy漏洞库，需要手工下载

- dorycli install pull 命令在默认安装情况下，不进行镜像的拉取、构建和推送到内部镜像仓库

- dorycli admin 命令的类型参数名字更新

- dory-engine v2.5.2 数据结构升级

- dory-engine 支持使用外部制品仓库，支持ftp sftp http方式上传制品

- dory-engine 支持制品仓库功能，可以把制品保存在DORY中

- dory-engine 支持不设置制品扫描仓库，不设置情况下不启用代码扫描功能

- dory-engine 支持在已有的代码仓库中创建演示项目代码和演示配置

- dory-engine 支持设置项目的演示代码目录信息

