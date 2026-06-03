#!/bin/bash
# UploadMyself 模型下载脚本
# 需要先安装: pip install huggingface_hub

set -e

MODEL_DIR="${MODEL_DIR:-./ml/models}"
mkdir -p "$MODEL_DIR"

echo "=== UploadMyself 模型下载 ==="
echo "目标目录: $MODEL_DIR"
echo ""

# 1. GPT-SoVITS (语音克隆)
echo "[1/5] 下载 GPT-SoVITS v2..."
if [ ! -d "$MODEL_DIR/voice/gptsovits" ]; then
    mkdir -p "$MODEL_DIR/voice/gptsovits"
    # TODO: 从 HuggingFace 下载权重
    echo "  请手动下载: https://github.com/RVC-Boss/GPT-SoVITS"
    echo "  解压到: $MODEL_DIR/voice/gptsovits/"
else
    echo "  已存在，跳过"
fi

# 2. CosyVoice2 (语音合成)
echo "[2/5] 下载 CosyVoice2..."
if [ ! -d "$MODEL_DIR/voice/cosyvoice" ]; then
    mkdir -p "$MODEL_DIR/voice/cosyvoice"
    # TODO: 从 HuggingFace 下载权重
    echo "  请手动下载: https://github.com/FunAudioLLM/CosyVoice"
    echo "  解压到: $MODEL_DIR/voice/cosyvoice/"
else
    echo "  已存在，跳过"
fi

# 3. LivePortrait (2D 驱动)
echo "[3/5] 下载 LivePortrait..."
if [ ! -d "$MODEL_DIR/avatar_2d/liveportrait" ]; then
    mkdir -p "$MODEL_DIR/avatar_2d/liveportrait"
    echo "  请手动下载: https://github.com/KwaiVGI/LivePortrait"
    echo "  解压到: $MODEL_DIR/avatar_2d/liveportrait/"
else
    echo "  已存在，跳过"
fi

# 4. SadTalker (音频驱动)
echo "[4/5] 下载 SadTalker..."
if [ ! -d "$MODEL_DIR/avatar_2d/sadtalker" ]; then
    mkdir -p "$MODEL_DIR/avatar_2d/sadtalker"
    echo "  请手动下载: https://github.com/OpenTalker/SadTalker"
    echo "  解压到: $MODEL_DIR/avatar_2d/sadtalker/"
else
    echo "  已存在，跳过"
fi

# 5. InstantMesh (3D 重建)
echo "[5/5] 下载 InstantMesh..."
if [ ! -d "$MODEL_DIR/avatar_3d/instantmesh" ]; then
    mkdir -p "$MODEL_DIR/avatar_3d/instantmesh"
    echo "  请手动下载: https://github.com/TencentARC/InstantMesh"
    echo "  解压到: $MODEL_DIR/avatar_3d/instantmesh/"
else
    echo "  已存在，跳过"
fi

echo ""
echo "=== 下载完成 ==="
echo "注意: 部分模型需要手动下载，请参考各模型的 GitHub 页面"
echo "模型目录结构:"
find "$MODEL_DIR" -maxdepth 2 -type d | sort
