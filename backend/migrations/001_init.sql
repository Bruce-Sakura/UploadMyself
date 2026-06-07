-- 001_init.sql — UploadMyself initial schema
-- Replaces GORM AutoMigrate. Idempotent (IF NOT EXISTS) so it can run on every boot.

-- Skill — thinking framework clone
CREATE TABLE IF NOT EXISTS skills (
    id         TEXT PRIMARY KEY,
    name       VARCHAR(128) NOT NULL,
    corpus     TEXT,
    status     VARCHAR(32) NOT NULL DEFAULT 'pending', -- pending | processing | done | failed
    result     TEXT,                                   -- generated SKILL.md
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Voice — voice clone
CREATE TABLE IF NOT EXISTS voices (
    id             TEXT PRIMARY KEY,
    name           VARCHAR(128) NOT NULL,
    audio_path     VARCHAR(512),
    duration       DOUBLE PRECISION NOT NULL DEFAULT 0, -- seconds
    model_path     VARCHAR(512),
    ref_audio_path VARCHAR(512),
    status         VARCHAR(32) NOT NULL DEFAULT 'pending',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Avatar — 2D/3D virtual avatar
CREATE TABLE IF NOT EXISTS avatars (
    id          TEXT PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    type        VARCHAR(8) NOT NULL,        -- 2d | 3d
    photo_path  VARCHAR(512),
    style       VARCHAR(32),                -- realistic | cartoon | anime
    status      VARCHAR(32) NOT NULL DEFAULT 'pending',
    result      TEXT,                       -- output JSON (cartoon/views/glb/obj URLs)
    output_path VARCHAR(512),               -- quick-access output file
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Task — async task tracking
CREATE TABLE IF NOT EXISTS tasks (
    id         TEXT PRIMARY KEY,
    type       VARCHAR(32) NOT NULL,        -- skill_process | voice_train | avatar_process
    ref_id     VARCHAR(64),                 -- related entity ID
    status     VARCHAR(32) NOT NULL DEFAULT 'pending',
    progress   INTEGER NOT NULL DEFAULT 0,  -- 0-100
    error      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_tasks_ref_id ON tasks (ref_id);
CREATE INDEX IF NOT EXISTS idx_tasks_type ON tasks (type);

-- FileUpload — uploaded file metadata
CREATE TABLE IF NOT EXISTS file_uploads (
    id            TEXT PRIMARY KEY,
    original_name VARCHAR(256) NOT NULL,
    stored_path   VARCHAR(512) NOT NULL,
    size          BIGINT NOT NULL DEFAULT 0,
    mime_type     VARCHAR(128),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Message — agent conversation history
CREATE TABLE IF NOT EXISTS messages (
    id              BIGSERIAL PRIMARY KEY,
    conversation_id TEXT,
    role            VARCHAR(32),            -- system | user | assistant | tool
    content         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages (conversation_id);
