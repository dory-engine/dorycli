# dorycli changelog v1.6.6

**新特性:**

- dory-engine 支持在流水线上设置cpu架构，一个项目可以通过不同的流水线，实现跨x86/arm64架构的编译、打包、部署
- dory-engine 环境管理可以支持一个k8s集群有x86/arm64混合的节点，可以根据nodeSelector自动识别集群节点所支持的cpu架构，并且自动识别默认使用的cpu架构
- dory-engine 把cpu架构信息从项目上移出，调整到在流水线上配置
- dory-engine 流水线的步骤新增cpu架构信息
- dory-engine 所有步骤执行记录新增cpu架构信息，可以看到编译、打包、部署步骤使用的cpu架构
- dory-engine 环境组件部署可以选择cpu架构，自动识别环境中是否有可用的cpu架构的节点
- dory-engine 环境调试组件无需设置cpu架构，自动根据k8s环境的默认cpu架构部署对应的调试组件
- dory-engine 容器镜像打包定义的Dockerfile中支持根据流水线的cpu架构动态设置来源镜像
- dory-engine 流水线、运行记录、步骤执行记录可以显示并过滤cpu架构信息
- dory-engine 运行、步骤的度量统计信息支持使用cpu架构作为维度归类统计信息，也支持使用cpu架构信息过滤统计数据
- dory-engine 新建项目以及为项目分配新nodePort端口段现在支持手工设置使用哪个nodePort端口段
- dory-engine 控制台的项目查看页支持显示环境的cpu架构信息，管理控制台的环境管理页面支持显示环境的cpu架构信息

**问题修复:**

- dory-engine 修复新建项目提示制品仓库类型不能为空的问题
- dory-engine 修复新建项目harbor的存储空间配额设置提示错误的问题
- dory-engine 修复 管理控制台 - 租户管理 按照租户编码搜索过滤有问题
- dorycli 修复 admin apply custom-step提示错误的问题
- dorycli 修复 def get pipeline 不要显示ops流水线的问题
- dorycli 支持显示cpu架构信息

