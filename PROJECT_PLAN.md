# UploadMyself — 项目总规划

> **克隆你自己**：输入照片 + 文本语料 + 语音样本 → 生成你的数字分身（思维框架 + 声音 + 2D/3D 虚拟形象）

---

## 一、项目概述

UploadMyself 是一个**全栈数字人克隆平台**，用户只需提供：

1. **文本语料**（文章/聊天记录/博客）→ 生成思维框架 Skill（仿女娲）
2. **语音样本**（3-10 分钟音频）→ 克隆声音
3. **一张正面照片** → 生成 2D 动态形象 / 3D 可交互虚拟人

支持**本地模型推理**和**云端 API 调用**两种模式，可选**模型蒸馏**降低部署成本。

---

## 二、技术栈选型

### 后端（Python）

| 层级 | 技术 | 理由 |
|------|------|------|
| Web 框架 | **FastAPI** | 异步、高性能、自带 OpenAPI 文档 |
| 任务队列 | **Celery + Redis** | 重量级推理任务异步处理 |
| 数据库 | **PostgreSQL + SQLAlchemy** | 用户数据、任务状态、生成历史 |
| 缓存 | **Redis** | 会话缓存、任务状态 |
| 存储 | **MinIO / 本地 FS** | 音频/图片/模型文件存储 |
| ML 框架 | **PyTorch** | 统一推理引擎，所有模型共用 |

### ML / AI 模块

| 功能 | 模型方案 | 备选 |
|------|---------|------|
| **思维框架蒸馏** | LLM 调用（仿女娲流程） | Qwen2.5 / DeepSeek / GPT-4o |
| **语音克隆** | **GPT-SoVITS v2**（少样本，中文优） | CosyVoice2、Fish Speech、OpenVoice |
| **语音合成(TTS)** | GPT-SoVITS / CosyVoice2 | ChatTTS、Edge-TTS（云端） |
| **2D 数字人驱动** | **LivePortrait**（照片→动态脸） | SadTalker、MuseTalk |
| **音频驱动口型** | SadTalker / MuseTalk | Wav2Lip |
| **3D 人脸重建** | **InstantMesh / Wonder3D** | Tripo3D API（云端） |
| **3D 全身生成** | **SMPL-X + 3D Gaussian Splatting** | Rodin API（云端） |
| **3D 渲染** | Three.js（前端） | Babylon.js |
| **模型蒸馏** | 自定义 KD pipeline | Hugging Face Optimum |

### 前端

| 层级 | 技术 | 理由 |
|------|------|------|
| 框架 | **React 18 + TypeScript** | 生态成熟、组件丰富 |
| 3D 渲染 | **Three.js + React Three Fiber** | Web 端 3D 虚拟人渲染 |
| 2D 渲染 | **HTML5 Canvas / Video** | 视频流播放 |
| UI 库 | **Ant Design 5** | 企业级组件、中文友好 |
| 状态管理 | **Zustand** | 轻量 |
| 构建 | **Vite** | 快速 |

---

## 三、项目目录结构

```
UploadMyself/
├── README.md                           # 项目说明
├── LICENSE                             # MIT License
├── pyproject.toml                      # Python 项目配置 (uv/pip)
├── docker-compose.yml                  # 一键启动
├── Makefile                            # 常用命令
├── .env.example                        # 环境变量模板
│
├── backend/                            # 后端服务
│   ├── __init__.py
│   ├── main.py                         # FastAPI 入口
│   ├── config.py                       # 配置管理
│   ├── api/                            # API 路由
│   │   ├── __init__.py
│   │   ├── auth.py                     # 认证
│   │   ├── skill.py                    # 思维框架 Skill 生成 API
│   │   ├── voice.py                    # 语音克隆 API
│   │   ├── avatar_2d.py                # 2D 数字人 API
│   │   ├── avatar_3d.py                # 3D 数字人 API
│   │   ├── distill.py                  # 模型蒸馏 API
│   │   └── tasks.py                    # 异步任务状态查询
│   │
│   ├── core/                           # 核心业务逻辑
│   │   ├── __init__.py
│   │   ├── skill_engine/               # 思维框架引擎（仿女娲）
│   │   │   ├── __init__.py
│   │   │   ├── collector.py            # 语料采集与清洗
│   │   │   ├── analyzer.py             # 思维模式分析
│   │   │   ├── extractor.py            # 心智模型提取
│   │   │   ├── synthesizer.py          # 框架合成
│   │   │   ├── validator.py            # 质量验证
│   │   │   └── templates/              # Skill 输出模板
│   │   │       ├── skill_template.md
│   │   │       └── persona_template.md
│   │   │
│   │   ├── voice_engine/               # 语音克隆引擎
│   │   │   ├── __init__.py
│   │   │   ├── trainer.py              # 声音模型训练/微调
│   │   │   ├── synthesizer.py          # 语音合成推理
│   │   │   ├── preprocessor.py         # 音频预处理(降噪/切片/VAD)
│   │   │   └── speaker.py              # 说话人特征提取
│   │   │
│   │   ├── avatar_engine/              # 虚拟形象引擎
│   │   │   ├── __init__.py
│   │   │   ├── photo_processor.py      # 照片预处理(人脸检测/对齐)
│   │   │   ├── avatar_2d.py            # 2D 形象生成与驱动
│   │   │   ├── avatar_3d.py            # 3D 形象重建
│   │   │   ├── animator.py             # 动画驱动(口型/表情/动作)
│   │   │   └── renderer.py             # 渲染输出
│   │   │
│   │   └── distill_engine/             # 模型蒸馏引擎
│   │       ├── __init__.py
│   │       ├── teacher.py              # 教师模型管理
│   │       ├── student.py              # 学生模型构建
│   │       ├── pipeline.py             # 蒸馏训练流程
│   │       └── evaluator.py            # 蒸馏效果评估
│   │
│   ├── models/                         # 数据模型
│   │   ├── __init__.py
│   │   ├── user.py
│   │   ├── skill.py
│   │   ├── voice.py
│   │   ├── avatar.py
│   │   └── task.py
│   │
│   ├── services/                       # 服务层
│   │   ├── __init__.py
│   │   ├── storage.py                  # 文件存储服务
│   │   ├── queue.py                    # 任务队列
│   │   ├── model_registry.py           # 模型注册与管理
│   │   └── provider/                   # 模型提供者(本地/云端切换)
│   │       ├── __init__.py
│   │       ├── base.py                 # 抽象基类
│   │       ├── local_provider.py       # 本地模型推理
│   │       └── cloud_provider.py       # 云端 API 调用
│   │
│   └── workers/                        # Celery 异步任务
│       ├── __init__.py
│       ├── celery_app.py
│       ├── skill_worker.py
│       ├── voice_worker.py
│       ├── avatar_worker.py
│       └── distill_worker.py
│
├── ml/                                 # ML 模型与脚本
│   ├── models/                         # 预训练模型存放
│   │   ├── voice/
│   │   │   ├── gptsovits/              # GPT-SoVITS 权重
│   │   │   └── cosyvoice/              # CosyVoice 权重
│   │   ├── avatar_2d/
│   │   │   ├── liveportrait/           # LivePortrait 权重
│   │   │   └── sadtalker/              # SadTalker 权重
│   │   ├── avatar_3d/
│   │   │   ├── instantmesh/            # InstantMesh 权重
│   │   │   └── gaussian_splatting/     # 3DGS 权重
│   │   └── base/
│   │       └── llm/                    # 基座 LLM (Qwen/DeepSeek)
│   │
│   ├── scripts/                        # ML 辅助脚本
│   │   ├── download_models.sh          # 模型下载脚本
│   │   ├── preprocess_audio.py         # 音频批量预处理
│   │   ├── preprocess_image.py         # 图片批量预处理
│   │   ├── export_onnx.py              # 模型导出 ONNX
│   │   ├── benchmark.py                # 推理性能基准测试
│   │   └── distill_voice.py            # 语音模型蒸馏
│   │
│   └── configs/                        # 模型配置文件
│       ├── gptsovits_v2.yaml
│       ├── cosyvoice2.yaml
│       ├── liveportrait.yaml
│       └── instantmesh.yaml
│
├── frontend/                           # 前端应用
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   ├── index.html
│   ├── public/
│   └── src/
│       ├── main.tsx                    # 入口
│       ├── App.tsx                     # 路由
│       ├── api/                        # API 客户端
│       │   ├── client.ts
│       │   ├── skill.ts
│       │   ├── voice.ts
│       │   └── avatar.ts
│       ├── components/                 # 通用组件
│       │   ├── Layout/
│       │   ├── Upload/
│       │   ├── AudioPlayer/
│       │   └── ModelSelector/
│       ├── pages/                      # 页面
│       │   ├── Home/                   # 首页/引导页
│       │   ├── SkillClone/             # 思维框架克隆
│       │   ├── VoiceClone/             # 语音克隆
│       │   ├── Avatar2D/               # 2D 数字人
│       │   ├── Avatar3D/               # 3D 数字人
│       │   ├── Distill/                # 模型蒸馏
│       │   ├── Playground/             # 综合体验(对话+语音+形象)
│       │   └── Settings/               # 模型/Provider 配置
│       ├── three/                      # Three.js 3D 相关
│       │   ├── Scene.tsx
│       │   ├── AvatarModel.tsx
│       │   ├── controls.ts
│       │   └── shaders/
│       ├── stores/                     # Zustand 状态
│       └── utils/
│
├── skills/                             # 生成的 Skill 存放
│   └── {user_id}/
│       ├── persona/
│       │   └── SKILL.md                # 仿女娲格式的思维框架
│       ├── voice/
│       │   ├── reference.wav           # 参考音频
│       │   └── model/                  # 微调后的声音模型
│       └── avatar/
│           ├── source.png              # 原始照片
│           ├── avatar_2d/              # 2D 生成结果
│           └── avatar_3d/              # 3D 模型文件(.glb/.vrm)
│
├── docs/                               # 文档
│   ├── architecture.md                 # 架构设计
│   ├── api.md                          # API 文档
│   ├── deployment.md                   # 部署指南
│   ├── models.md                       # 模型说明
│   └── contributing.md                 # 贡献指南
│
└── tests/                              # 测试
    ├── unit/
    ├── integration/
    └── e2e/
```

---

## 四、功能模块详细设计

### 模块 1：思维框架克隆（仿女娲）

**输入**：用户的文本语料（文章、聊天记录、博客、笔记）

**流程**：

```
用户上传文本语料
    ↓
Phase 1: 语料采集与清洗 (collector.py)
  - 文本分段、去重、格式化
  - 提取关键表达模式
    ↓
Phase 2: 思维模式分析 (analyzer.py)
  - LLM 分析：核心论点、决策模式、表达风格
  - 提取高频词、句式偏好、价值观信号
    ↓
Phase 3: 心智模型提取 (extractor.py)
  - 三重验证：跨域复现 / 生成力 / 排他性
  - 筛选 3-7 个核心心智模型
    ↓
Phase 4: 框架合成 (synthesizer.py)
  - 生成 SKILL.md（仿女娲格式）
  - 含：身份卡 / 心智模型 / 决策启发式 / 表达DNA / 诚实边界
    ↓
Phase 5: 质量验证 (validator.py)
  - 已知测试 / 边缘测试 / 风格测试
  - 输出质量报告
```

**Provider 切换**：
- 本地：Qwen2.5-72B / DeepSeek-V3（通过 vLLM/Ollama）
- 云端：OpenAI GPT-4o / Claude / Qwen API

### 模块 2：语音克隆

**输入**：3-10 分钟语音样本

**流程**：

```
用户上传音频
    ↓
预处理 (preprocessor.py)
  - 降噪 (noisereduce)
  - VAD 切片 (Silero-VAD)
  - 采样率统一 (16kHz/22.05kHz/44.1kHz)
    ↓
说话人特征提取 (speaker.py)
  - 提取音色 embedding
  - 分析语速/音高/停顿模式
    ↓
声音模型训练 (trainer.py)
  - GPT-SoVITS: 参考音频 + 文本 → 微调
  - CosyVoice2: Speaker embedding 注入
  - 支持 few-shot (3-10 句即可)
    ↓
语音合成 (synthesizer.py)
  - 输入文本 → 输出克隆声音音频
  - 支持流式输出
```

**本地模型**：GPT-SoVITS v2 / CosyVoice2
**云端 API**：ElevenLabs / Fish Audio / MiniMax TTS

### 模块 3：2D 虚拟形象

**输入**：一张正面照片

**流程**：

```
用户上传照片
    ↓
照片预处理 (photo_processor.py)
  - 人脸检测 (InsightFace)
  - 人脸对齐 & 裁剪
  - 质量评估（光照/遮挡/角度）
    ↓
2D 形象生成 (avatar_2d.py)
  - 风格化：写实 / 卡通 / 动漫（可选）
  - 背景替换/移除
    ↓
动态驱动 (animator.py)
  - LivePortrait: 驱动源 → 目标表情迁移
  - SadTalker/MuseTalk: 音频 → 口型同步
  - 输出：带音频的说话视频
```

**本地模型**：LivePortrait + SadTalker
**云端 API**：HeyGen / D-ID / Synthesia

### 模块 4：3D 虚拟形象

**输入**：一张正面照片（可选多角度）

**流程**：

```
用户上传照片
    ↓
3D 人脸重建 (avatar_3d.py)
  - InstantMesh: 单图 → 3D Mesh
  - FLAME/SMPL-X: 参数化人脸/人体模型
    ↓
纹理生成
  - 纹理贴图从照片投影
  - 可选 AI 纹理增强
    ↓
骨骼绑定 & 动画
  - VRM 标准骨骼
  - BlendShape 面部表情
  - 支持 VRM/VRoid 格式
    ↓
Web 端渲染 (Three.js)
  - React Three Fiber 加载 .glb/.vrm
  - 实时口型同步（音频驱动 BlendShape）
  - 用户可交互（旋转/缩放/换装）
```

**本地模型**：InstantMesh + Gaussian Splatting
**云端 API**：Tripo3D / Rodin / Meshy

### 模块 5：模型蒸馏

**目标**：将大模型压缩为轻量版，降低部署成本

**支持场景**：

| 蒸馏类型 | 教师 → 学生 | 用途 |
|---------|------------|------|
| LLM 蒸馏 | Qwen-72B → Qwen-7B | 思维框架推理加速 |
| 语音蒸馏 | CosyVoice-Large → CosyVoice-Small | TTS 推理加速 |
| 2D 模型蒸馏 | LivePortrait-Base → Lite | 端侧部署 |
| 知识蒸馏 | 自定义 pipeline | 通用 |

**流程**：

```
选择教师模型 & 学生模型架构
    ↓
准备蒸馏数据集
  - 教师模型推理输出作为 soft label
  - 原始数据作为 hard label
    ↓
蒸馏训练 (pipeline.py)
  - Loss = α × KL(teacher||student) + (1-α) × CE(student, label)
  - 温度参数 T 调节 soft label 平滑度
    ↓
评估 (evaluator.py)
  - 精度对比 (teacher vs student)
  - 推理速度对比
  - 模型大小对比
    ↓
导出 & 部署
  - ONNX / TorchScript 导出
  - 量化 (INT8/FP16)
```

---

## 五、Provider 架构（本地/云端切换）

```python
# 统一接口设计
class VoiceProvider(ABC):
    @abstractmethod
    async def clone_and_synthesize(self, ref_audio, text, **kwargs) -> AudioResult:
        ...

class LocalVoiceProvider(VoiceProvider):
    """GPT-SoVITS / CosyVoice 本地推理"""
    ...

class CloudVoiceProvider(VoiceProvider):
    """ElevenLabs / Fish Audio API"""
    ...
```

用户可在 Settings 页面自由切换：
- **本地模式**：所有推理在本机 GPU 运行，隐私优先
- **云端模式**：调用第三方 API，免部署
- **混合模式**：核心模型本地 + 增强功能云端

---

## 六、API 设计概览

```
POST   /api/v1/skill/create              # 创建思维框架 Skill
POST   /api/v1/skill/upload-corpus       # 上传语料
GET    /api/v1/skill/{id}/result          # 获取 Skill 结果
POST   /api/v1/voice/upload               # 上传语音样本
POST   /api/v1/voice/train                # 启动声音训练
POST   /api/v1/voice/synthesize           # 语音合成
POST   /api/v1/avatar/2d/upload           # 上传照片
POST   /api/v1/avatar/2d/generate         # 生成 2D 形象
POST   /api/v1/avatar/2d/animate          # 驱动动画(音频→口型)
POST   /api/v1/avatar/3d/upload           # 上传照片
POST   /api/v1/avatar/3d/reconstruct      # 3D 重建
GET    /api/v1/avatar/3d/{id}/model        # 下载 3D 模型
POST   /api/v1/distill/start              # 启动蒸馏任务
GET    /api/v1/distill/{id}/status         # 蒸馏进度
GET    /api/v1/tasks/{id}                  # 通用任务状态查询
WS     /api/v1/ws/stream                  # WebSocket 实时流(语音/视频)
```

---

## 七、开发路线图

### Phase 1: MVP（4 周）
- [x] 项目骨架搭建（FastAPI + React + 目录结构）
- [ ] 思维框架克隆（基于 LLM 的 Skill 生成）
- [ ] 语音克隆（GPT-SoVITS 集成）
- [ ] 2D 数字人（LivePortrait 驱动）
- [ ] 基础前端 UI（上传 + 结果展示）

### Phase 2: 增强（3 周）
- [ ] 3D 数字人重建与渲染
- [ ] 音频驱动口型同步
- [ ] Provider 切换框架（本地/云端）
- [ ] 异步任务队列
- [ ] 用户系统 & 历史记录

### Phase 3: 蒸馏 & 优化（3 周）
- [ ] 模型蒸馏 pipeline
- [ ] 量化部署（INT8/FP16）
- [ ] 流式推理优化
- [ ] WebSocket 实时交互

### Phase 4: 打磨 & 发布（2 周）
- [ ] 综合 Playground（对话 + 语音 + 形象一体化）
- [ ] Docker 一键部署
- [ ] 文档完善
- [ ] 性能优化 & 测试

---

## 八、部署方案

### 开发环境
```bash
# 后端
cd backend && pip install -e ".[dev]"
uvicorn backend.main:app --reload

# 前端
cd frontend && npm install && npm run dev

# 依赖服务
docker-compose up -d redis postgres minio
```

### 生产环境
```bash
# Docker Compose 一键部署
docker-compose -f docker-compose.prod.yml up -d

# GPU 节点需要 nvidia-docker
# 推荐配置：RTX 4090 × 1（本地推理）或 A100 × 1（蒸馏训练）
```

### 硬件需求

| 场景 | 最低配置 | 推荐配置 |
|------|---------|---------|
| 纯云端 API | CPU 4核, 8GB RAM | CPU 8核, 16GB RAM |
| 本地推理(单模块) | RTX 3060 12GB | RTX 4090 24GB |
| 全模块本地推理 | RTX 4090 24GB | A100 80GB |
| 模型蒸馏训练 | A100 40GB | A100 80GB × 2 |

---

## 九、关键依赖

```toml
[project]
dependencies = [
    "fastapi>=0.110",
    "uvicorn[standard]>=0.29",
    "celery[redis]>=5.4",
    "sqlalchemy[asyncio]>=2.0",
    "asyncpg>=0.29",
    "redis>=5.0",
    "minio>=7.2",
    "torch>=2.3",
    "torchaudio>=2.3",
    "torchvision>=0.18",
    "numpy>=1.26",
    "librosa>=0.10",
    "soundfile>=0.12",
    "insightface>=0.7",
    "onnxruntime-gpu>=1.17",
    "diffusers>=0.28",
    "transformers>=4.40",
    "accelerate>=0.30",
    "pydantic>=2.7",
    "python-multipart>=0.0.9",
]
```

---

## 十、许可证

MIT License — 开源自由使用

---

*Generated by UploadMyself Project Planner*
