# Today-Todo

一个把「任务管理 + 健康节律 + 防沉迷提醒」融合在一起的跨端效率系统。

普通 TODO 只记录任务状态；Today-Todo 更关注你的执行状态，目标是让你不仅列清单，还能真正完成。

## 项目亮点

### 亮点 1：行为驱动的专注评分（Focus Score）

系统会把任务完成、喝水、站立、短视频行为合并计算为专注评分，而不是只看任务勾选率。

- 高完成 + 低分心 -> 分数提升
- 短视频次数上升 -> 分数下降
- 结果可直接反馈到 Web/Android 端

### 亮点 2：今日智能洞察（新增）

新增 `GET /api/v1/insights/today`，根据当天行为给出：

- 风险等级：low / medium / high
- 节奏判断：稳定推进 / 节奏偏慢 / 注意力波动
- 建议动作：下一步最小可执行行动
- 推荐任务：自动挑选最多 3 个未完成任务（按优先级和时序）

这让系统从“记录工具”变成“策略助手”。

> 说明：当前洞察为规则/公式驱动（Rule-based），不是机器学习模型。
> 主要基于完成率、喝水进度、站立进度、短视频次数与专注评分阈值来判定风险等级和建议动作。

## 技术栈

- 后端：Go + Gin + GORM + SQLite
- Web：HTML + CSS + JavaScript + SSE
- Android：Kotlin + Jetpack Compose + Retrofit + WorkManager

## 目录结构

```text
Today-Todo/
├─ controllers/                  # API 控制器
│  ├─ todo_controller.go         # TODO 增删改查
│  ├─ health.go                  # 健康与每日统计
│  ├─ reminder_controller.go     # SSE 实时提醒
│  ├─ state.go                   # 状态机与调度器生命周期
│  └─ insight_controller.go      # 今日智能洞察（新增）
├─ models/                       # 数据模型与数据库初始化
├─ services/                     # 提醒调度、事件广播
├─ routers/                      # 路由注册与静态资源入口
├─ web/                          # Web 页面
├─ android-app/                  # Android 工程
├─ main.go
├─ go.mod
└─ README.md
```

## 核心能力

- TODO：新增、编辑、删除、完成
- 健康打卡：喝水、站立
- 防沉迷：短视频次数记录与提醒
- 实时提醒：SSE 推送 + 本地兜底定时器
- 每日统计：任务、喝水、站立、短视频、专注评分
- 今日洞察：风险等级 + 建议动作 + 推荐任务

## API 一览

基础路径：`/api/v1`

- `POST /todos`
- `GET /todos`
- `PUT /todos/:id`
- `DELETE /todos/:id`
- `POST /health/water`
- `POST /health/stand`
- `POST /health/short-video`
- `GET /health/daily-progress?user_id=1`
- `GET /insights/today?user_id=1` (新增)
- `GET /reminders/stream`
- `GET /reminders/config`

## 快速开始

### 1. 启动后端

```bash
go run .
```

### 2. 访问 Web

```text
http://localhost:8080/
```

### 3. 示例请求

新增任务：

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"写周报","description":"整理本周进展","priority":"high"}'
```

查询今日洞察：

```bash
curl "http://localhost:8080/api/v1/insights/today?user_id=1"
```

## Android 运行说明

- 用 Android Studio 打开 `android-app/`
- 模拟器访问后端：`http://10.0.2.2:8080/`
- 真机访问后端：请改成电脑局域网 IP（手机和电脑需同网段）

## 测试命令

已内置最小测试样例：

- 单元测试：`controllers/logic_test.go`
- 集成测试：`integration/api_integration_test.go`

执行命令：

```bash
# 全量测试
go test ./... -v

# 仅单元测试（controllers）
go test ./controllers -v

# 仅集成测试（HTTP + SQLite）
go test ./integration -v
```

## 工程化建议

- 容器化：增加 Dockerfile 与 docker-compose
- 可观测性：请求日志、错误追踪、性能指标

## 简历表达建议

可定位为：

- "跨端效率系统：Go API + Web + Android"
- "将健康行为与任务执行融合，构建可解释的专注评分与智能洞察"

建议在简历里附上：

- 1 张 Web 页面截图
- 1 段 30 秒演示视频
- 2-3 个核心接口示例
