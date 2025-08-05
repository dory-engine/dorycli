# dorycli changelog v1.7.0

**新特性:**

- dorycli console get 命令去掉 --items 参数，改为使用args
- dorycli console delete 命令去掉 --items 参数，改为使用args
- dorycli install 命令安装dory，新增以下构建环境: maven-jdk12 gradle-jdk12
- dorycli install 命令安装harbor，harbor证书从365天改为3650天
- dorycli admin get 命令支持 --filter 参数
- dorycli admin get 命令 支持 --name 参数，只显示名字
- dorycli install pull 命令不执行tag镜像操作，仅pull必需的镜像
- dorycli install script 命令自动输出自动配置nexus的脚本
- dorycli install script 命令自动输出自动配置gitlab或者gitea的脚本
- dorycli install script 命令自动输出自动配置sonarqube的脚本
- dorycli install script 命令自动输出自动配置harbor的脚本
- dorycli install 命令安装dory支持使用csi共享存储
- dorycli install pull 命令可以自动拉取harbor需要的镜像
- dorycli install 命令安装dory无需mount共享存储，自动通过kubectl cp命令把安装的初始化文件发送到共享存储
- dorycli install 命令升级 harbor 版本为 v2.13.1
- dorycli install 命令升级 gitea 版本为 1.24.3
- dorycli install 命令升级 gitlab 版本为 17.11.6-ce.0
- dorycli install 命令升级 nexus 版本为 3.82.0
- dorycli install 命令升级 sonarqube 版本为 10.1.0-community

- dory-engine 环境组件支持job和cronjob方式进行部署
- dory-engine 环境组件支持部署configmaps和secrets到kubernetes集群
- dory-engine 一个kubernetes集群，支持通过不同的namespace部署多个不同的dory-engine实例
- dory-engine checkDeploy步骤无论成功还是失败，都显示部署的event事件
- dory-engine 环境管理支持使用insecure方式连接kubernetes集群
- dory-engine scanImage步骤设置的时候，漏洞数量设置为-1表示不检测漏洞数量
- dory-engine packageImage步骤在推送镜像的时候，可以显示build context上传到docker daemon的进度
- dory-engine 开通项目的时候自动开通项目的npm以及pypi项目依赖库
- dory-engine 升级所有go mod依赖库为最新版本，提升稳定性
- dory-engine 控制台环境信息可以显示kubernetes环境的节点信息
- dory-engine kubernetes环境支持使用csi共享存储，内置支持csi-cephfs以及csi-nfs共享存储
- dory-engine kubernetes环境支持设置自定义csi的pv/pvc模板
- dory-engine 项目开通的时候，可以设置项目的harbor空间配额
- dory-engine 项目开通的时候，可以设置不允许使用那些名字作为项目的projectName
- dory-engine 项目开通的时候，可以设置检测对应的projectName是否在环境中已经存在对应的namespace，如果namespace存在那么不允许创建项目
- dory-engine 优化运行日志中表格的显示方式
- dory-engine gitPull、artifact步骤支持设置超时时间，超时可以直接终止流水线
- dory-engine 步骤支持设置超时时间，超时可以直接终止流水线
- dory-engine 流水线点击终止按钮，无论执行到哪个步骤，都可以即时终止流水线
- dory-engine 流水线执行结束后，执行日志写入到文件中，不再从redis读取流水线执行日志
- dory-engine 控制台的异步操作执行结束后，审计日志写入到文件中，不再从redis读取审计日志，大幅节约redis内存消耗
- dory-engine getRunSettings步骤支持设置是否在日志的表格中显示所有步骤概要
- dory-engine gitPull步骤进行git pull或者git clone的时候，步骤执行日志可以正常处理\r换行问题
- dory-engine 调试容器更新为 doryengine/debian-vnc-ssh，支持web方式访问vnc图形界面

- dory-console 开发空间和控制台项目列表页支持按照projectDesc搜索
- dory-console 环境组件以及调试组件页面的pod信息可以链接到kubernetes-dashboard查看pod信息
- dory-console deploy / checkDeloy / undo / syncImage 步骤执行记录的pod信息可以链接到kubernetes-dashboard查看pod信息
- dory-console 开发空间和控制台项目列表页项目团队支持下拉选择

**功能弃用:**

- dory-engine 移除kubernetes环境的glusterfs和rbd持久化存储的支持
- dorycli install 移除 run 命令
- dorycli install 移除 docker 模式的支持，仅支持把dory部署到kubernetes集群中

**问题修复:**

- dory-engine 修复程序异常: update runs not finish and not running error
- dory-engine 修复假如kubernetes环境无法访问的情况下，界面会出现卡死的问题
- dory-engine 修复假如kubernetes环境无法访问的情况下，无法删除环境的问题
- dory-engine 修复自定义步骤实际执行失败的情况下，步骤执行记录依然显示成功的问题
- dory-engine 项目定义假如没有修改，不再提示no change错误
- dory-engine checkDeploy步骤执行kubectl logs和kubectl describe执行异常情况下，不会结束步骤执行
- dory-engine 修复在项目控制台更新token的之后，没有使用最新的token更新kubernets环境中的项目secret配置的问题
- dory-engine 修复调试容器的proxy代理configmap的解析问题
- dory-engine 修复新增arm架构的kubernetes环境的时候，project-data-pod创建异常问题
- dory-engine 修复gitPull步骤过程中获取git diff处理逻辑会出现卡死的问题
- dory-engine 优化程序性能，getRunSettings 步骤执行时间从10秒下降到毫秒级
- dory-engine 优化程序性能，获取运行列表性能优化
- dory-engine 优化程序性能，查看容器部署定义生成的yaml接口执行时间从32秒下降到0.5秒
- dory-engine 优化程序性能，查看项目定义的历史记录
- dory-engine 优化程序性能，优化gitPull步骤过程中获取git diff处理逻辑，提升执行速度

- dory-console 修复开发空间页面点击搜索后，流水线定义的保存按钮更新的目标项目不正确的问题
- dory-console 修复环境管理页面无法进行分页的问题
- dory-console 修复commit提交记录页面修改每页显示多少记录操作无效的问题
