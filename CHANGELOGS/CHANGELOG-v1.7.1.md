# dorycli changelog v1.7.1

**新特性:**

- dory-engine 优化docker exec执行过程中日志输出逻辑
- dory-engine .dockerignore默认增加默认配置
- dory-engine packageImage镜像打包步骤支持输出build context文件数量进度
- dory-engine 自定义步骤设置的时候，租户管理员不允许使用gitRepo，避免租户管理员可以拉取任意代码
- dory-engine 容器部署定义和环境组件设置的时候，service可以选择的类型支持loadbalancer
- dory-engine 调整流水线gitPull.timeout默认为0，默认代码拉取不设置超时
- dory-engine 支持设置租户可以分配的环境总资源配额，如果不设置表示不限制
- dory-engine 租户新增、修改项目、给项目分配新环境的时候检查租户环境配额
- dory-engine 支持管理员创建自定义资源模板
- dory-engine 项目管理员可以通过项目控制台在kubernetes环境创建自定义资源，自定义资源类似helm模板，用户设置模板参数，在kubernetes环境中创建对应的自定义资源对象

- dory-console 项目定义界面支持模块名字过滤，可以过滤显示对应模块的定义
- dory-console 运行日志支持下载到文件
- dory-console 审计日志支持下载到文件

- dorycli install pull 镜像拉取与镜像构建新增platform参数
- dorycli project get新增ignore-demo参数，output新增name参数
- dorycli admin/console/def/pipeline/project/run get -o 支持name参数

**问题修复:**

- dory-engine 修复运行日志中输出的表格错乱问题
- dory-engine 解决trivy扫描的时候，因为db读写冲突导致不能并行扫描的问题
- dory-engine 修复流水线日志写入残缺不完整的问题
- dory-engine 修复使用csi作为环境的共享存储的时候，创建项目的pv路径异常问题

- dorycli 修复表格输出对齐问题
- dorycli 修复admin get custom-step重复显示问题
- dorycli 修复trivy漏洞库下载操作指引错误问题
