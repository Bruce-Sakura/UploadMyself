"""语料采集与清洗"""

from pathlib import Path


class CorpusCollector:
    """从用户上传的文本中采集和清洗语料"""

    async def collect(self, file_path: Path) -> str:
        """读取并清洗文本"""
        content = Path(file_path).read_text(encoding="utf-8")
        return self._clean(content)

    def _clean(self, text: str) -> str:
        """基础清洗：去除多余空白、统一编码"""
        lines = text.strip().splitlines()
        cleaned = [line.strip() for line in lines if line.strip()]
        return "\n".join(cleaned)

    def chunk(self, text: str, max_chars: int = 4000) -> list[str]:
        """将长文本分段，用于 LLM 分析"""
        paragraphs = text.split("\n\n")
        chunks = []
        current = ""
        for para in paragraphs:
            if len(current) + len(para) > max_chars:
                if current:
                    chunks.append(current.strip())
                current = para
            else:
                current += "\n\n" + para
        if current.strip():
            chunks.append(current.strip())
        return chunks
