# UploadMyself — 核心架构重定义

## 本质

**UploadMyself 不是工具集合，是一个住在电脑里的「你」。**

```
用户上传数据
    ↓
┌─────────────────────────────────────┐
│         UploadMyself Agent          │
│                                     │
│  ┌───────────┐  ┌───────────────┐  │
│  │ SKILL.md  │  │   LLM 推理    │  │
│  │ (你的思维) │→│  (你的大脑)    │  │
│  └───────────┘  └───────┬───────┘  │
│                         │           │
│  ┌──────────────────────┼────────┐ │
│  │              Tool Layer       │ │
│  │  ┌─────┐ ┌─────┐ ┌─────┐    │ │
│  │  │Shell│ │Browser│ │File │    │ │
│  │  └─────┘ └─────┘ └─────┘    │ │
│  │  ┌─────┐ ┌─────┐ ┌─────┐    │ │
│  │  │Email│ │Code │ │API  │    │ │
│  │  └─────┘ └─────┘ └─────┘    │ │
│  └──────────────────────┬────────┘ │
│                         │           │
│  ┌──────────┐  ┌───────┴───────┐  │
│  │ Voice    │  │   Avatar      │  │
│  │ (你的声音)│  │  (你的外貌)   │  │
│  └──────────┘  └───────────────┘  │
└─────────────────────────────────────┘
```

## 四个核心模块

### 1. 🧠 SKILL.md = 你的思维方式（系统提示词）
- 输入：用户的文本语料（聊天记录/文章/笔记）
- 输出：仿女娲格式的 SKILL.md
- 用途：**作为 Agent 的 system prompt**，定义「你」怎么思考、怎么说话、怎么做决策

### 2. 🛠️ Tool Layer = 你的能力（技能系统）
- **Shell 执行**：运行命令、操控电脑
- **文件操作**：读写文件、管理数据
- **浏览器控制**：搜索、浏览、自动化操作
- **代码执行**：写代码、跑脚本
- **API 调用**：外部服务集成
- 用户可以自定义工具（就像你学会新技能）

### 3. 🎤 Voice = 你的声音（语音输出）
- 输入：用户的语音样本
- 输出：克隆的 TTS 引擎
- 用途：Agent 回复时用你的声音说出来

### 4. 🖼️ Avatar = 你的外貌（虚拟形象）
- 输入：用户的照片
- 输出：2D/3D 可驱动的虚拟形象
- 用途：Agent 对话时展示你的形象，口型同步

## 对话流程

```
用户输入（文字/语音）
    ↓
Agent 加载 SKILL.md 作为 system prompt
    ↓
LLM 推理（用你的思维方式思考）
    ↓
判断是否需要工具
    ├── 需要 → 调用 Tool → 获取结果 → 继续推理
    └── 不需要 → 直接生成回复
    ↓
回复输出
    ├── 文字 → 渲染到聊天界面
    ├── 语音 → TTS 引擎用你的声音说出来
    └── 形象 → Avatar 口型同步 + 表情驱动
```

## 技术栈

| 模块 | 技术 | 说明 |
|------|------|------|
| Agent Core | Go (Gin) + LLM API | 编排层，分层模块结构 |
| 数据访问 | pgxpool (github.com/jackc/pgx/v5/pgxpool) | 连接池，取代旧的 GORM |
| 思维框架 | LLM (OpenAI 兼容客户端) | 语料分析 → SKILL.md |
| 工具系统 | Go plugin / subprocess | 可扩展工具注册 |
| 语音克隆 | GPT-SoVITS / CosyVoice2 | TTS |
| 虚拟形象 | LivePortrait / Three.js | 2D/3D 驱动 |
| 前端 | React + WebSocket | 实时对话界面 |

## 后端分层结构

后端已从 GORM 单层结构重构为**分层模块结构**（不再使用 gorm）：

```
backend/
├── main.go                  # 入口：依赖注入，pgxpool 连接 + 启动时执行 go:embed 的 migrations
├── migrations/001_init.sql  # DDL（通过 go:embed 嵌入，启动时按语句执行）
├── internal/llm/client.go   # OpenAI 兼容 LLM 客户端
├── pkg/                     # 每张表一个模块，每模块六层
│   ├── tasks/               # ← 模块模板（异步任务追踪）
│   ├── skills/              # 思维框架 / SKILL.md 生成
│   ├── voices/              # 语音克隆（训练 / 合成）
│   ├── avatars/             # 2D/3D 形象处理
│   ├── file_uploads/        # 文件 / 语料上传（OCR/PDF/Word）
│   └── messages/            # Agent 对话、会话消息、工具注册
├── cmd/upme/                # upme CLI 客户端（调用 REST API）
└── middleware/              # CORS 等
```

**每个 `pkg/<模块>` 含六层**（以 `backend/pkg/tasks` 为范例）：
`entity`（表结构）→ `dto`（请求/响应对象）→ `mapper`（pgx 数据访问）→
`service`（业务接口）→ `service/impl`（接口实现）→ `handler`（HTTP 处理 + 路由 `Register`）。

依赖在 `main.go` 中按 `mapper → service → handler` 顺序注入。
旧的 `backend/agent`、`backend/handler`、`backend/model` 包已删除。

## Skill 文件存储

思维框架（生成的或导入的）以**文件目录**形式持久化，而非存进数据库：

```
<SKILLS_DIR>/<id>/
├── SKILL.md     # 思维框架正文
└── meta.json    # 元信息（id / name / source / imported_at）
```

- 数据库 `skills` 表只做**元数据索引**（id、name、status、corpus 等）；`SKILL.md`
  正文从磁盘读取（`Get`/`List` 时回填到 `result` 字段）。
- 配置项 `SKILLS_DIR` 默认 `./skills`，启动时自动创建；Docker 中挂载为 `skills`
  卷，使 skill 在容器重建后保留。

## CLI (upme)

`backend/cmd/upme` 是后端 REST API 的命令行客户端，构建：
`cd backend && go build -o upme ./cmd/upme`。后端地址默认 `http://localhost:8000`，
可用 `-server` 或环境变量 `$UPME_SERVER` 覆盖。

| 命令 | 说明 |
|------|------|
| `upme health` | 健康检查 |
| `upme skill list` | 列出所有 skill |
| `upme skill get -id <id>` | 查看单个 skill（含 SKILL.md） |
| `upme skill import -url <url> [-name <n>]` | 从 URL/GitHub 导入 |
| `upme skill import -file <path> [-name <n>]` | 从本地文件导入 |
| `upme skill new -name <n> -corpus <text\|@file>` | 用语料生成（触发 LLM） |
| `upme skill rm -id <id>` | 删除 |
| `upme chat -skill <id> -m "<msg>" [-conv <id>]` | 与分身对话 |

## API 端点

所有路由挂载在 `/api/v1` 分组下，由各模块 handler 的 `Register` 注册。

```
# Agent（messages 模块）
POST /api/v1/agent/chat            # 核心对话端点：加载 SKILL.md 作为 system prompt，
                                   #   LLM 推理 + 工具调用，返回文字 + 音频 + Avatar 数据
GET  /api/v1/agent/tools           # 列出可用工具

# Skills（思维框架）
POST   /api/v1/skills              # 创建
POST   /api/v1/skills/import       # 导入现成 SKILL.md（{url|content, name?}）
GET    /api/v1/skills              # 列表
GET    /api/v1/skills/:id          # 详情
PUT    /api/v1/skills/:id          # 更新
DELETE /api/v1/skills/:id          # 删除
POST   /api/v1/skills/:id/process  # 触发语料分析生成 SKILL.md

# Voices（语音克隆）
POST   /api/v1/voices              # 创建
GET    /api/v1/voices              # 列表
GET    /api/v1/voices/:id          # 详情
DELETE /api/v1/voices/:id          # 删除
POST   /api/v1/voices/:id/train       # 训练
POST   /api/v1/voices/:id/synthesize  # 合成

# Avatars（虚拟形象）
POST   /api/v1/avatars             # 创建
GET    /api/v1/avatars             # 列表
GET    /api/v1/avatars/:id         # 详情
DELETE /api/v1/avatars/:id         # 删除
POST   /api/v1/avatars/:id/process # 触发处理

# Tasks（异步任务追踪）
GET    /api/v1/tasks               # 列表（可按 type / ref_id 过滤）
GET    /api/v1/tasks/:id           # 详情

# File uploads（文件 / 语料上传）
POST   /api/v1/upload              # 通用文件上传
POST   /api/v1/upload-corpus       # 语料上传（OCR/PDF/Word）
GET    /api/v1/files/:id           # 文件下载

# 其它
GET    /health                     # 健康检查
GET    /uploads/*                  # 静态文件服务
```

### 导入 Skill：`POST /api/v1/skills/import`

无需重新调用 LLM 即可导入现成的 `SKILL.md`：

- 请求体 `{url|content, name?}`，`url` 与 `content` 二选一。
- GitHub `blob` 网页地址自动转换为 `raw.githubusercontent.com` 原始地址。
- `name` 缺省时从 YAML frontmatter 的 `name:` 字段解析（Claude Agent Skills 格式），
  兜底为 `imported-skill`。

## Avatar 处理（2D / 3D）

`avatars` 模块按类型走不同的 ML 流程：

- **2D**：调用 `ml_service` 时传 `mode=2d`，**跳过 3D 重建**，输出卡通图
  （`cartoon_image`），`output_path` 取该图。
- **3D**：传 `mode=3d`，进行三维重建；返回的 `glb_model` / `obj_model` / `views`
  路径统一归一化为可访问的 `/uploads` URL，`output_path` 取 `glb_model`。
