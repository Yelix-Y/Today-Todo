# Today-Todo

创新型跨端 TODO 系统，集成任务管理、健康提醒、防沉迷提醒与每日统计。

## 项目目录结构

```text
Today-Todo/
├─ controllers/
│  ├─ todo_controller.go          # TODO 增删改查
│  ├─ health.go                   # 喝水/站立/短视频记录 + 每日统计
│  ├─ reminder_controller.go      # SSE 提醒流 + 提醒配置
│  └─ state.go                    # 状态机接口 + 调度器生命周期
├─ models/
│  ├─ setup.go                    # 数据库连接与自动迁移
│  ├─ todo.go                     # TODO 数据模型
│  ├─ health.go                   # 健康与统计数据模型
│  └─ state.go                    # 用户状态机模型
├─ routers/
│  └─ setup.go                    # API 路由与 Web 静态资源入口
├─ services/
│  ├─ scheduler.go                # 并发提醒调度器（喝水/站立/防沉迷）
│  └─ reminder_hub.go             # SSE 广播中心
├─ web/
│  ├─ index.html                  # Web 页面结构
│  ├─ style.css                   # 响应式样式 + 动画
│  └─ app.js                      # 前端状态、API 请求、提醒逻辑
├─ android-app/
│  ├─ settings.gradle.kts
│  ├─ build.gradle.kts
│  ├─ gradle.properties
│  └─ app/
│     ├─ build.gradle.kts
│     └─ src/main/
│        ├─ AndroidManifest.xml
│        ├─ java/com/todaytodo/app/
│        │  ├─ MainActivity.kt
│        │  ├─ data/              # Retrofit API + 数据模型 + Repository
│        │  ├─ ui/                # Compose 页面与 ViewModel
│        │  └─ reminder/          # WorkManager 通知提醒
│        └─ res/values/
│           ├─ strings.xml
│           └─ themes.xml
├─ main.go
├─ go.mod
└─ README.md
```

## 技术栈

- 后端：Go + Gin + GORM + SQLite
- Web：HTML + CSS + JavaScript + SSE
- Android：Kotlin + Jetpack Compose + Retrofit + WorkManager + Notification

## 核心功能

- TODO：添加、删除、编辑、完成任务
- 健康提醒：喝水记录与提醒、每小时站立/拉伸提醒
- 防沉迷：短视频次数记录 + 动效提醒弹窗
- 每日统计：
  - 完成任务数/总任务数
  - 喝水总量与打卡次数
  - 站立总分钟与站立次数
  - 短视频刷屏次数
  - 专注评分 `focus_score`

## 模块简述

- `controllers`：把业务动作暴露成 REST API。
- `services`：负责后台定时提醒、SSE 广播。
- `models`：定义持久化结构与统计结构。
- `web`：桌面+手机浏览器可用的单页应用，含提醒动效。
- `android-app`：原生 Android 实现，支持周期通知与动效提醒卡片。

## API 概览（`/api/v1`）

- `POST /todos`
- `GET /todos`
- `PUT /todos/:id`
- `DELETE /todos/:id`
- `POST /health/water`
- `POST /health/stand`
- `POST /health/short-video`
- `GET /health/daily-progress?user_id=1`
- `GET /reminders/stream`（SSE）
- `GET /reminders/config`

## 使用示例

### 1) 启动后端

```bash
go run .
```

### 2) 打开 Web

- 浏览器访问 `http://localhost:8080/`

### 3) 创建任务

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"写技术方案","description":"补齐提醒动画","priority":"high"}'
```

### 4) 记录喝水

```bash
curl -X POST http://localhost:8080/api/v1/health/water \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"amount":250}'
```

### 5) 查询每日统计

```bash
curl "http://localhost:8080/api/v1/health/daily-progress?user_id=1"
```

## Android 运行

- 使用 Android Studio 打开 `android-app/` 目录并执行 Sync + Run。
- Android 模拟器访问后端地址默认使用 `http://10.0.2.2:8080/`。
- 真机调试时请改为你电脑局域网 IP。

## UI 动画设计建议与交互方案

- 喝水提醒：水滴脉冲 + 蓝绿色高亮边框，按钮「喝水 +250ml」。
- 站立提醒：橙色呼吸光效 + 倒计时文案「站立 5 分钟」。
- 防沉迷提醒：红色闪烁边框 + 轻微缩放动效 + 快速记录按钮。
- 统计面板：环形进度展示喝水/站立，短视频次数用高对比警示色。

## 简单布局示意

```text
[Header: 今日目标 + 刷新]
[Stats: 任务 | 喝水环 | 站立环 | 短视频]
[Todo 输入 + Todo 列表]
[健康打卡按钮区]
[提醒弹窗层（SSE/本地定时触发）]
```

## 可扩展创新功能

- 番茄节律模式：25/5 自动切换，提醒样式跟状态联动。
- 能量曲线：每 2~3 小时自评精力，叠加任务完成趋势。
- 智能建议：短视频次数连续上升时，自动推荐低门槛任务。
