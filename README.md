# 基于七米老师的go-web开发的bluebell基本项目，仅供自学

## 版本
 v1: 基于gin-mysql-redis的博客后端项目
 v2: 重构基于gorm框架

### 功能
注册用户
用户名/密码登录、邮箱登录。用户单点登录
邮箱绑定和解绑
用户头像上传和修改
创建社区、查询社区
发布帖子、查询帖子详情
帖子投票
帖子热度点赞或者时间排序

TO：
帖子浏览次数记录redis
mysql读写分离
context上下文信息传递和超时处理
每个层进行责任划分（定义接口实现方法）
Jager分布式链路追踪
opentelemetry收集检测数据

### 技术栈
1. gin框架
2. gorm框架
3. zap日志库
4. viper配置管理
5. swagger生成文档
6. JWT认证鉴权
7. CORS跨域
8. 令牌桶限流
9. 雪花算法生成用户ID
10. go操作redis
11. go操作mysql
12. air实现热重载
13. mysql读写分离
14. Jager分布式链路追踪
15. opentelemetry
16. docker部署
