-- 001_init.sql — UploadMyself 初始 schema (SQLite)
-- 由 main.go 启动时 go:embed 执行，幂等。
-- 可空文本列统一 DEFAULT '' 以便扫描进 Go string。

CREATE TABLE IF NOT EXISTS skills (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    corpus     TEXT NOT NULL DEFAULT '',
    status     TEXT NOT NULL DEFAULT 'pending', -- pending | processing | done | failed
    result     TEXT NOT NULL DEFAULT '',        -- generated SKILL.md
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS voices (
    id             TEXT PRIMARY KEY,
    name           TEXT NOT NULL,
    audio_path     TEXT NOT NULL DEFAULT '',
    duration       REAL NOT NULL DEFAULT 0,
    model_path     TEXT NOT NULL DEFAULT '',
    ref_audio_path TEXT NOT NULL DEFAULT '',
    status         TEXT NOT NULL DEFAULT 'pending',
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS avatars (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    type        TEXT NOT NULL,                  -- 2d | 3d
    photo_path  TEXT NOT NULL DEFAULT '',
    style       TEXT NOT NULL DEFAULT '',       -- realistic | cartoon | anime
    status      TEXT NOT NULL DEFAULT 'pending',
    result      TEXT NOT NULL DEFAULT '',       -- output JSON (cartoon/views/glb/obj URLs)
    output_path TEXT NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
    id         TEXT PRIMARY KEY,
    type       TEXT NOT NULL,                   -- skill_process | voice_train | avatar_process
    ref_id     TEXT NOT NULL DEFAULT '',
    status     TEXT NOT NULL DEFAULT 'pending',
    progress   INTEGER NOT NULL DEFAULT 0,      -- 0-100
    error      TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_tasks_ref_id ON tasks (ref_id);
CREATE INDEX IF NOT EXISTS idx_tasks_type ON tasks (type);

CREATE TABLE IF NOT EXISTS file_uploads (
    id            TEXT PRIMARY KEY,
    original_name TEXT NOT NULL,
    stored_path   TEXT NOT NULL,
    size          INTEGER NOT NULL DEFAULT 0,
    mime_type     TEXT NOT NULL DEFAULT '',
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    conversation_id TEXT NOT NULL DEFAULT '',
    role            TEXT NOT NULL DEFAULT '',   -- system | user | assistant | tool
    content         TEXT NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages (conversation_id);
