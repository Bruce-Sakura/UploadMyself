"""音频预处理"""

import subprocess
from pathlib import Path


class AudioPreprocessor:
    """音频预处理：降噪、切片、格式转换"""

    def __init__(self, sample_rate: int = 22050):
        self.sample_rate = sample_rate

    async def process(self, input_path: Path, output_dir: Path) -> list[Path]:
        """完整预处理流程"""
        output_dir.mkdir(parents=True, exist_ok=True)

        # 1. 转换为 WAV
        wav_path = output_dir / "converted.wav"
        await self._convert_to_wav(input_path, wav_path)

        # 2. 降噪
        clean_path = output_dir / "clean.wav"
        await self._denoise(wav_path, clean_path)

        # 3. VAD 切片
        segments = await self._vad_segment(clean_path, output_dir)

        return segments

    async def _convert_to_wav(self, input_path: Path, output_path: Path):
        """统一转为 WAV 格式"""
        cmd = [
            "ffmpeg", "-i", str(input_path),
            "-ar", str(self.sample_rate),
            "-ac", "1", "-f", "wav",
            str(output_path), "-y"
        ]
        subprocess.run(cmd, capture_output=True, check=True)

    async def _denoise(self, input_path: Path, output_path: Path):
        """降噪处理"""
        # TODO: 使用 noisereduce 库
        import shutil
        shutil.copy(input_path, output_path)

    async def _vad_segment(self, input_path: Path, output_dir: Path) -> list[Path]:
        """VAD 语音活动检测切片"""
        # TODO: 使用 Silero-VAD 切片
        # 临时：直接返回整段
        return [input_path]
