# UploadMyself — Development TODO

> Generated: 2026-06-03 | Last Updated: 2026-06-03

## Legend

- ⬜ Not started
- 🔵 In progress
- ✅ Done
- 🔴 Blocked

---

## Phase 0: Project Foundation ✅

| # | Task | Status | Notes |
|---|------|--------|-------|
| 0.1 | Git repo + remote setup | ✅ | github.com/Bruce-Sakura/UploadMyself |
| 0.2 | Go backend skeleton (Gin) | ✅ | |
| 0.3 | React frontend skeleton | ✅ | |
| 0.4 | Docker Compose (Redis/PG/MinIO) | ✅ | |
| 0.5 | PROJECT_PLAN.md | ✅ | |

---

## Phase 1: Voice Cloning Engine 🔵

### 1.1 Audio Preprocessing Pipeline

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 1.1.1 | FFmpeg format conversion (any→WAV) | ⬜ | `os/exec` ffmpeg | 22050Hz mono WAV |
| 1.1.2 | Noise reduction | ⬜ | Python `noisereduce` via subprocess | Go calls Python script |
| 1.1.3 | VAD segmentation | ⬜ | **Silero-VAD v5** (2025) | ONNX runtime in Go, or Python subprocess |
| 1.1.4 | Speaker embedding extraction | ⬜ | **Resemblyzer** / **pyannote-audio 3.x** | Extract d-vector / speaker embedding |
| 1.1.5 | Audio quality validation | ⬜ | SNR check + clipping detection | Reject low-quality samples |

### 1.2 Voice Cloning — Local Model

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 1.2.1 | **GPT-SoVITS v2** integration | ⬜ | [RVC-Boss/GPT-SoVITS](https://github.com/RVC-Boss/GPT-SoVITS) | 1-min few-shot, best CN quality |
| 1.2.2 | GPT-SoVITS training worker | ⬜ | Asynq task → Python subprocess | Fine-tune with user audio |
| 1.2.3 | GPT-SoVITS inference worker | ⬜ | API or subprocess | Text → cloned voice audio |
| 1.2.4 | **CosyVoice2-0.5B** integration | ⬜ | [FunAudioLLM/CosyVoice](https://github.com/FunAudioLLM/CosyVoice) | Streaming TTS, 30-50% fewer errors than v1 |
| 1.2.5 | **Fish Speech v1.5** integration | ⬜ | [fishaudio/fish-speech](https://github.com/fishaudio/fish-speech) | DualAR, 500M params, 80+ languages, zero-shot |
| 1.2.6 | Model hot-swap (switch at runtime) | ⬜ | ModelRegistry + Provider pattern | User picks which engine to use |

### 1.3 Voice Cloning — Cloud API

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 1.3.1 | **ElevenLabs** API integration | ⬜ | REST API | Premium quality, $$$ |
| 1.3.2 | **Fish Audio** cloud API | ⬜ | REST API | Cheaper, S2 Pro model |
| 1.3.3 | **MiniMax TTS** API | ⬜ | REST API | Good CN quality |
| 1.3.4 | Unified cloud provider interface | ⬜ | Go `provider.CloudProvider` | Abstract API differences |

### 1.4 Voice Synthesis API

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 1.4.1 | POST /api/v1/voice/upload | ⬜ | Gin + MinIO | Store audio, trigger preprocess |
| 1.4.2 | POST /api/v1/voice/train | ⬜ | Asynq async task | Returns task_id |
| 1.4.3 | POST /api/v1/voice/synthesize | ⬜ | Sync or async | Text → audio file |
| 1.4.4 | GET /api/v1/voice/:id/samples | ⬜ | | Preview samples |
| 1.4.5 | WebSocket streaming TTS | ⬜ | `gorilla/websocket` | Real-time audio stream |

---

## Phase 2: 2D Avatar Engine 🔵

### 2.1 Photo Processing

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 2.1.1 | Face detection + alignment | ⬜ | **InsightFace** (ONNX in Go) or Python subprocess | Detect, crop, align |
| 2.1.2 | Face quality scoring | ⬜ | Blur/lighting/occlusion check | Reject bad photos |
| 2.1.3 | Background removal | ⬜ | **RMBG-2.0** / `rembg` | Transparent PNG output |
| 2.1.4 | Style transfer (optional) | ⬜ | Realistic → Cartoon/Anime | SD img2img or dedicated model |

### 2.2 Face Animation — Local

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 2.2.1 | **LivePortrait** integration | ⬜ | [KlingAIResearch/LivePortrait](https://github.com/KlingAIResearch/LivePortrait) | Expression retargeting, stitching control |
| 2.2.2 | **SadTalker** integration | ⬜ | [OpenTalker/SadTalker](https://github.com/OpenTalker/SadTalker) | Audio→3DMM→face, CVPR 2023 |
| 2.2.3 | **HunyuanVideo-Avatar** integration | ⬜ | [Tencent/HunyuanVideo-Avatar](https://github.com/Tencent/HunyuanVideo-Avatar) | **NEWEST (2025)**: full-body, multi-char, singing support |
| 2.2.4 | **Sonic** integration | ⬜ | [jixiaozhong/Sonic](https://github.com/jixiaozhong/Sonic) | Global audio perception, real-time |
| 2.2.5 | Audio→lip sync pipeline | ⬜ | Wav2Lip or MuseTalk | Dedicated lip-sync quality |

### 2.3 Face Animation — Cloud API

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 2.3.1 | **HeyGen** API | ⬜ | REST API | Enterprise-grade, expensive |
| 2.3.2 | **D-ID** API | ⬜ | REST API | Good quality, moderate price |
| 2.3.3 | **Synthesia** API | ⬜ | REST API | Video generation platform |

### 2.4 2D Avatar API

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 2.4.1 | POST /api/v1/avatar/2d/upload | ⬜ | | Photo upload + face detect |
| 2.4.2 | POST /api/v1/avatar/2d/generate | ⬜ | | Style + generate avatar |
| 2.4.3 | POST /api/v1/avatar/2d/animate | ⬜ | | Audio → talking video |
| 2.4.4 | Video output streaming | ⬜ | | Chunked response or WebSocket |

---

## Phase 3: 3D Avatar Engine 🔵

### 3.1 3D Reconstruction — Local

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 3.1.1 | **InstantMesh** integration | ⬜ | [TencentARC/InstantMesh](https://github.com/TencentARC/InstantMesh) | Single image → 3D mesh, sparse-view |
| 3.1.2 | **TripoSR** integration | ⬜ | [VAST-AI-Research/TripoSR](https://github.com/VAST-AI-Research/TripoSR) | MIT license, <0.5s on A100 |
| 3.1.3 | **SF3D** (Stable Fast 3D) | ⬜ | CVPR 2025 | 0.5s reconstruction with UV texture |
| 3.1.4 | **3DGS Head Avatar** | ⬜ | CVPR 2025 TensorialGaussianAvatar | Real-time rendering, expressive |
| 3.1.5 | FLAME/SMPL-X parametric model | ⬜ | SMPL-X for body, FLAME for face | Standard parametric models |

### 3.2 3D Reconstruction — Cloud API

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 3.2.1 | **Tripo3D v2.5** API | ⬜ | REST API (fal.ai / direct) | Best cloud 3D gen, <0.5s |
| 3.2.2 | **Rodin** API | ⬜ | REST API | By Deemos, good quality |
| 3.2.3 | **Meshy** API | ⬜ | REST API | Text/Image → 3D |

### 3.3 3D Rendering (Frontend)

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 3.3.1 | GLB/VRM loader | ⬜ | Three.js + `@pixiv/three-vrm` | Load 3D models in browser |
| 3.3.2 | Real-time BlendShape animation | ⬜ | Three.js morphTargets | Audio → facial expression |
| 3.3.3 | Orbit controls + interaction | ⬜ | Three.js OrbitControls | Rotate/zoom/pan |
| 3.3.4 | VRM avatar customization | ⬜ | VRoid SDK | Change clothes/hair/color |
| 3.3.5 | WebGPU renderer (optional) | ⬜ | Three.js r164+ WebGPU | Better perf on supported browsers |

### 3.4 3D Avatar API

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 3.4.1 | POST /api/v1/avatar/3d/upload | ⬜ | | Multi-photo upload |
| 3.4.2 | POST /api/v1/avatar/3d/reconstruct | ⬜ | | Photo → 3D model |
| 3.4.3 | GET /api/v1/avatar/3d/:id/model | ⬜ | | Download GLB/VRM |
| 3.4.4 | GET /api/v1/avatar/3d/:id/preview | ⬜ | | Three.js render data |

---

## Phase 4: Skill Cloning Engine (仿女娲) ⬜

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 4.1 | Corpus collector (text cleaning) | ⬜ | Go stdlib | UTF-8, dedup, chunk |
| 4.2 | LLM analysis (mind model extraction) | ⬜ | **Qwen2.5-72B** / **DeepSeek-V3** / **GPT-4o** | 3-7 mental models |
| 4.3 | Decision heuristics extraction | ⬜ | LLM | 5-10 rules |
| 4.4 | Expression DNA analysis | ⬜ | LLM | Style/syntax/humor |
| 4.5 | SKILL.md synthesizer | ⬜ | Template + LLM | Nuwa format output |
| 4.6 | Quality validator (3 test types) | ⬜ | LLM-based evaluation | Sanity/edge/voice check |
| 4.7 | Local LLM provider (Ollama/vLLM) | ⬜ | HTTP API call | Qwen2.5 / DeepSeek local |
| 4.8 | Cloud LLM provider | ⬜ | OpenAI/Anthropic/Qwen API | Fallback chain |

---

## Phase 5: Model Distillation ⬜

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 5.1 | KD pipeline (teacher→student) | ⬜ | PyTorch + HuggingFace | KL + CE loss |
| 5.2 | LLM distillation (72B→7B) | ⬜ | NeurIPS 2025 few-shot KD | Counterfactual approach |
| 5.3 | Voice model distillation | ⬜ | Feature-based KD | CosyVoice Large → Small |
| 5.4 | NVIDIA TensorRT optimization | ⬜ | TensorRT Model Optimizer | Pruning + quantization |
| 5.5 | ONNX export + INT8 quantization | ⬜ | `optimum` + `onnxruntime` | Cross-platform deploy |
| 5.6 | Benchmark (accuracy/speed/size) | ⬜ | Custom eval harness | Before vs After |

---

## Phase 6: Frontend UI ⬜

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 6.1 | Home page / onboarding | ⬜ | React + Ant Design | Upload flow |
| 6.2 | Skill clone page | ⬜ | | Text upload → SKILL.md preview |
| 6.3 | Voice clone page | ⬜ | | Audio upload → train → preview |
| 6.4 | 2D avatar page | ⬜ | | Photo → animated video player |
| 6.5 | 3D avatar page | ⬜ | React Three Fiber | Interactive 3D viewer |
| 6.6 | Playground (combined) | ⬜ | | Chat + voice + 3D avatar |
| 6.7 | Settings (provider/model) | ⬜ | | Switch local/cloud per module |
| 6.8 | Task queue dashboard | ⬜ | | Real-time progress |

---

## Phase 7: Integration & Polish ⬜

| # | Task | Status | Tech | Notes |
|---|------|--------|------|-------|
| 7.1 | User auth (JWT) | ⬜ | Go JWT middleware | Optional for local use |
| 7.2 | PostgreSQL models + migrations | ⬜ | GORM AutoMigrate + Atlas | User/Skill/Voice/Avatar |
| 7.3 | MinIO file storage integration | ⬜ | minio-go | Upload/download |
| 7.4 | WebSocket real-time updates | ⬜ | gorilla/websocket | Task progress push |
| 7.5 | Docker production compose | ⬜ | Multi-stage build | Go binary + Node + Python ML |
| 7.6 | CI/CD (GitHub Actions) | ⬜ | | Lint + test + build + push |
| 7.7 | API documentation (Swagger) | ⬜ | `swaggo/swag` | Auto-gen from Go comments |
| 7.8 | E2E tests | ⬜ | Playwright (frontend) + Go test | Full flow test |

---

## Priority Order (Recommended)

```
Week 1-2:  Phase 1 (Voice) — most standalone, clear I/O
Week 3-4:  Phase 2 (2D Avatar) — photo → video, high visual impact
Week 5-6:  Phase 4 (Skill) — LLM-heavy, depends on provider setup
Week 7-8:  Phase 3 (3D Avatar) — complex but incremental
Week 9-10: Phase 6 (Frontend) — wire up all backends
Week 11-12: Phase 5 (Distill) + Phase 7 (Polish)
```

---

## Key Tech Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Voice primary model | **GPT-SoVITS v2** | Best CN few-shot, 1-min training |
| Voice alternative | **Fish Speech v1.5** | 500M params, 80+ lang, zero-shot |
| 2D animation primary | **HunyuanVideo-Avatar** | Newest (2025), full-body, singing |
| 2D animation fallback | **LivePortrait** + **SadTalker** | Mature, well-tested |
| 3D reconstruction | **TripoSR** (local) + **Tripo3D v2.5** (cloud) | MIT + fastest cloud |
| LLM for Skill | **Qwen2.5-72B** (local) / **GPT-4o** (cloud) | Best CN reasoning |
| Distillation | **NVIDIA TensorRT Model Optimizer** | Production-grade pruning+KD |
| Go ↔ Python bridge | **os/exec** → Python subprocess | Simplest, decoupled |
