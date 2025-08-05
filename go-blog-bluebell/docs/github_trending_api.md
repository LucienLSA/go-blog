# GitHub热点数据API文档

## 功能概述

本项目新增了GitHub热点数据功能，包括：

1. **自动数据获取**：每小时自动从GitHub API获取热点仓库数据
2. **Redis缓存**：将数据存储在Redis中，提高访问速度
3. **API接口**：提供RESTful API接口供前端调用
4. **定时任务**：使用Go的定时器实现每小时数据刷新

## API接口

### 1. 获取GitHub热点数据

**接口地址：** `GET /api/v2/github/trending`

**请求参数：**
- `language` (可选): 编程语言筛选，如：javascript, python, go, java, rust
- `since` (可选): 时间范围，可选值：daily, weekly, monthly

**请求示例：**
```bash
# 获取所有语言的热点数据
GET /api/v2/github/trending

# 获取JavaScript语言的热点数据
GET /api/v2/github/trending?language=javascript

# 获取Python语言的每日热点数据
GET /api/v2/github/trending?language=python&since=daily
```

**响应示例：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "language": "javascript",
    "since": "daily",
    "data": [
      {
        "id": 123456789,
        "name": "awesome-project",
        "full_name": "user/awesome-project",
        "description": "这是一个很棒的项目",
        "language": "JavaScript",
        "stars": 1000,
        "forks": 100,
        "url": "https://github.com/user/awesome-project",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1
  }
}
```

### 2. 获取所有GitHub热点数据

**接口地址：** `GET /api/v2/github/trending/all`

**请求示例：**
```bash
GET /api/v2/github/trending/all
```

**响应示例：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "javascript:daily": [
      {
        "id": 123456789,
        "name": "awesome-project",
        "full_name": "user/awesome-project",
        "description": "这是一个很棒的项目",
        "language": "JavaScript",
        "stars": 1000,
        "forks": 100,
        "url": "https://github.com/user/awesome-project",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "python:weekly": [
      // ... 更多数据
    ]
  }
}
```

### 3. 手动刷新GitHub热点数据

**接口地址：** `POST /api/v2/github/trending/refresh`

**请求示例：**
```bash
POST /api/v2/github/trending/refresh
```

**响应示例：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "message": "GitHub热点数据刷新任务已启动"
  }
}
```

## 技术实现

### 1. 数据存储

- **Redis Key设计**：
  - 单个数据：`bluebell:github:trending:{language}:{since}`
  - 哈希表：`bluebell:github:trending:hash`
- **过期时间**：1小时
- **数据格式**：JSON序列化

### 2. 定时任务

- **执行频率**：每小时执行一次
- **执行内容**：获取所有支持的语言和时间范围的数据
- **支持语言**：javascript, python, go, java, rust
- **支持时间范围**：daily, weekly, monthly

### 3. 数据源

使用GitHub Trending API：
- **API地址**：https://github-trending-api.now.sh/repositories
- **请求方式**：GET
- **参数支持**：language, since

### 4. 错误处理

- API请求失败时记录错误日志
- Redis操作失败时记录错误日志
- 数据解析失败时记录错误日志
- 服务异常时返回统一的错误响应

## 部署说明

### 1. 环境要求

- Go 1.16+
- Redis 6.0+
- MySQL 8.0+

### 2. 配置说明

确保Redis配置正确，相关配置在`settings/config.go`中：

```go
type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     string `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
    PoolSize int    `mapstructure:"pool_size"`
}
```

### 3. 启动服务

```bash
# 编译项目
go build -o main .

# 启动服务
./main
```

服务启动后会自动：
1. 初始化Redis连接
2. 启动GitHub热点数据定时任务
3. 立即执行一次数据获取
4. 每小时自动刷新数据

### 4. 监控日志

可以通过日志监控定时任务的执行情况：

```bash
# 查看应用日志
tail -f logs/app.log
```

关键日志信息：
- `github trending data scheduler started` - 定时任务启动
- `starting github trending data refresh task` - 开始刷新任务
- `save github trending data success` - 数据保存成功
- `github trending data refresh task completed` - 刷新任务完成

## 注意事项

1. **API限制**：GitHub API有请求频率限制，建议合理使用
2. **数据准确性**：数据来源于第三方API，可能存在延迟或不准确的情况
3. **缓存策略**：数据缓存1小时，如需实时数据可调用手动刷新接口
4. **错误恢复**：服务重启时会自动重新获取数据
5. **资源消耗**：定时任务会消耗一定的网络和存储资源 