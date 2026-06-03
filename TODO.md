# UploadMyself — Development TODO

> **核心理念：不是工具集合，是住在电脑里的「你」**
> 
> SKILL.md = 你的思维 | Tool = 你的能力 | Voice = 你的声音 | Avatar = 你的外貌

## Phase 1: Agent Core（核心对话引擎）🔵

| # | Task | Status | Notes |
|---|------|--------|-------|
| 1.1 | POST /api/v1/agent/chat 端点 | ⬜ | 核心：接收消息 → 加载 SKILL.md → LLM 推理 → 返回回复 |
| 1.2 | SKILL.md 作为 system prompt 加载 | ⬜ | 从数据库读取 skill result，注入 LLM 对话 |
| 1.3 | LLM Provider 接口 | ⬜ | 支持 OpenAI/Qwen/Ollama API 调用 |
| 1.4 | 对话历史管理 | ⬜ | 多轮对话上下文，存数据库 |
| 1.5 | WebSocket 实时流 | ⬜ | 流式输出文字，支持中断 |

## Phase 2: 工具系统（技能/能力）⬜

| # | Task | Status | Notes |
|---|------|--------|-------|
| 2.1 | Tool 注册接口 | ⬜ | POST /api/v1/agent/tools/register |
| 2.2 | Shell 执行工具 | ⬜ | 运行命令，返回 stdout/stderr |
| 2.3 | 文件操作工具 | ⬜ | 读写文件，列目录 |
| 2.4 | 浏览器控制工具 | ⬜ | 搜索/浏览/截图（Playwright） |
| 2.5 | 代码执行工具 | ⬜ | 运行 Python/JS 代码片段 |
| 2.6 | 工具调用编排 | ⬜ | LLM function calling → 工具执行 → 结果回传 LLM |

## Phase 3: 思维框架生成（SKILL.md）🔵

| # | Task | Status | Notes |
|---|------|--------|-------|
| 3.1 | 语料分析 → SKILL.md 生成 | ✅ | analyze_corpus.py + LLM |
| 3.2 | SKILL.md 模板（仿女娲） | ⬜ | 心智模型 + 决策启发式 + 表达DNA |
| 3.3 | 质量验证 | ⬜ | 已知测试/边缘测试/风格测试 |

## Phase 4: 语音克隆 🔵

| # | Task | Status | Notes |
|---|------|--------|-------|
| 4.1 | 音频预处理 | ✅ | preprocess_audio.py |
| 4.2 | 声音训练 | ✅ | voice_clone_train.py |
| 4.3 | TTS 合成 | ✅ | voice_synthesize.py |
| 4.4 | Agent 回复 → 语音输出 | ⬜ | chat 回复自动调用 TTS |

## Phase 5: 虚拟形象 🔵

| # | Task | Status | Notes |
|---|------|--------|-------|
| 5.1 | 人脸检测 | ✅ | detect_face.py |
| 5.2 | 2D 形象生成 | ⬜ | LivePortrait |
| 5.3 | 3D 形象重建 | ⬜ | InstantMesh/TripoSR |
| 5.4 | Avatar 口型同步 | ⬜ | 音频驱动面部动画 |

## Phase 6: 前端对话界面 ⬜

| # | Task | Status | Notes |
|---|------|--------|-------|
| 6.1 | 对话页面 | ⬜ | 聊天气泡 + 输入框 + 发送 |
| 6.2 | 语音输入 | ⬜ | 录音 → STT → 文字 |
| 6.3 | 语音输出 | ⬜ | Agent 回复 → TTS → 播放 |
| 6.4 | Avatar 展示 | ⬜ | 3D 模型实时渲染 + 口型同步 |
| 6.5 | 工具调用可视化 | ⬜ | 显示 Agent 正在执行什么工具 |

## Phase 7: 集成打磨 ⬜

| # | Task | Status | Notes |
|---|------|--------|-------|
| 7.1 | 完整对话流程 | ⬜ | 文字输入 → 思考 → 工具调用 → 文字+语音+形象输出 |
| 7.2 | 自定义工具注册 | ⬜ | 用户可添加自己的技能 |
| 7.3 | Docker 生产部署 | ✅ | docker-compose |
| 7.4 | CI/CD | ⬜ | GitHub Actions |

---

## 开发优先级

```
Week 1:   Phase 1 (Agent Core) — 让它能对话
Week 2:   Phase 2 (Tool System) — 让它能做事
Week 3:   Phase 3+4 (SKILL + Voice) — 让它像你
Week 4:   Phase 6 (Frontend) — 让你能看到它
Week 5-6: Phase 5+7 (Avatar + Polish) — 完善
```

## 关键设计决策

| 决策 | 选择 | 理由 |
|------|------|------|
| Agent 核心 | **LLM function calling** | 标准化工具调用，兼容所有主流 LLM |
| SKILL.md 用途 | **System prompt** | 不是文档，是 Agent 的大脑 |
| 工具执行 | **Go subprocess** | 安全隔离，支持任何语言的工具 |
| 前端通信 | **WebSocket** | 实时流式输出，低延迟 |
| 语音输出 | **流式 TTS** | Agent 回复一段就播一段，不用等全部生成 |
