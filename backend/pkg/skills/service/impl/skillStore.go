package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// maxSkillDownload caps a downloaded SKILL.md at 5MB.
const maxSkillDownload = 5 << 20

// ---- file storage: <SkillsDir>/<id>/SKILL.md (+ meta.json, assets/) ----

func (s *SkillServiceImpl) skillDir(id string) string {
	return filepath.Join(s.cfg.SkillsDir, id)
}

func (s *SkillServiceImpl) skillFile(id string) string {
	return filepath.Join(s.skillDir(id), "SKILL.md")
}

func (s *SkillServiceImpl) writeSkillFile(id, content string) error {
	if err := os.MkdirAll(s.skillDir(id), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.skillFile(id), []byte(content), 0o644)
}

func (s *SkillServiceImpl) readSkillFile(id string) string {
	b, err := os.ReadFile(s.skillFile(id))
	if err != nil {
		return ""
	}
	return string(b)
}

func (s *SkillServiceImpl) removeSkillDir(id string) error {
	return os.RemoveAll(s.skillDir(id))
}

func (s *SkillServiceImpl) writeSkillMeta(id, name, source string) error {
	meta := map[string]string{
		"id":          id,
		"name":        name,
		"source":      source,
		"imported_at": time.Now().Format(time.RFC3339),
	}
	b, _ := json.MarshalIndent(meta, "", "  ")
	return os.WriteFile(filepath.Join(s.skillDir(id), "meta.json"), b, 0o644)
}

// ---- download + parsing helpers ----

// downloadSkill fetches a SKILL.md from a URL, rewriting GitHub blob URLs to raw.
func downloadSkill(ctx context.Context, rawURL string) (string, error) {
	u := toRawURL(strings.TrimSpace(rawURL))
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return "", fmt.Errorf("unsupported url scheme: %s", rawURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxSkillDownload))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// toRawURL converts a GitHub "blob" web URL to its raw.githubusercontent.com form.
//
//	https://github.com/owner/repo/blob/main/path/SKILL.md
//	-> https://raw.githubusercontent.com/owner/repo/main/path/SKILL.md
func toRawURL(u string) string {
	const ghPrefix = "https://github.com/"
	if strings.HasPrefix(u, ghPrefix) && strings.Contains(u, "/blob/") {
		u = strings.Replace(u, ghPrefix, "https://raw.githubusercontent.com/", 1)
		u = strings.Replace(u, "/blob/", "/", 1)
	}
	return u
}

// parseFrontmatterName extracts `name:` from a leading YAML frontmatter block
// (the format used by Claude Agent Skills). Returns "" if absent.
func parseFrontmatterName(content string) string {
	if !strings.HasPrefix(content, "---") {
		return ""
	}
	// Find the closing fence after the opening one.
	rest := content[3:]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return ""
	}
	fm := rest[:end]
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		}
	}
	return ""
}
