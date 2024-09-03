# dorycli changelog v1.6.5

**新特性:**

- dory-engine 社区版支持设置项目的代码仓库、镜像仓库、依赖与制品仓库、代码扫描仓库
- dory-engine 社区版限制项目开通、中间件部署、调试组件部署只能选择amd64架构
- dory-engine 社区版限制禁用自定义资源配额的功能
- dory-engine 新开通的演示项目默认禁用制品打包和制品主机部署步骤，默认只开启gin-demo模块
- dory-engine 项目开通支持设置harbor镜像仓库的空间配额
- dory-engine 可以设置流水线超过多长时间没有输入自动终止流水线
- dory-engine 自定义步骤假如从代码仓库拉取代码，那么自动创建一个gitPullCustomStep步骤
- dory-engine 项目定义页面和运行查看页面新增执行OPS批处理的菜单按钮
- dory-engine 控制台项目管理列表页面支持使用环境名过滤项目
- dory-engine 控制台页面支持搜索排序
- dory-engine 提高步骤消耗的容器cpu和内存资源的可读性
- dory-engine k8s环境支持使用storageClass来为项目动态分配PV和PVC
- dory-engine k8s环境列表可以展示可用的storageClass，也可以展示pv的状态
- dory-engine 用户管理列表页支持过滤不属于任何项目成员的用户
- dory-engine 新增 /api/console/project/:projectName/minimal 接口，用于dorycli console子命令

- dorycli 新增console子命令，支持通过命令行设置项目控制台信息，包括：项目成员、流水线、流水线触发器、项目主机、项目数据库、调试组件、项目组件等，需要项目维护者权限

**问题修复:**

- dory-engine 修复度量统计因为时区存在+8小时偏差导致统计数据异常的问题
- dory-engine 修复按照时间进行运行列表过滤因为时区存在+8小时偏差导致数据显示的问题
- dory-engine 重启dory-engine的时候，自动清理等待输入的流水线数据
- dory-engine 修复并行执行多个自定义步骤的时候，从代码仓库拉取自定义步骤代码会存在写入冲突的问题
