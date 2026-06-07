// Command upme is a CLI client for the UploadMyself backend REST API.
//
// Usage:
//
//	upme skill list
//	upme skill get    -id <id>
//	upme skill import -url <github/raw url> [-name <name>]
//	upme skill import -file <path>          [-name <name>]
//	upme skill new    -name <name> -corpus <text|@file>
//	upme skill rm     -id <id>
//	upme chat         -skill <id> -m "<message>" [-conv <id>]
//	upme health
//
// Server defaults to http://localhost:8000 (override with -server or $UPME_SERVER).
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func defaultServer() string {
	if s := os.Getenv("UPME_SERVER"); s != "" {
		return s
	}
	return "http://localhost:8000"
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "skill":
		skillCmd(os.Args[2:])
	case "chat":
		chatCmd(os.Args[2:])
	case "health":
		healthCmd(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Print(`upme — UploadMyself CLI

Commands:
  skill list                                  列出所有思维框架
  skill get    -id <id>                       查看单个 skill（含 SKILL.md 内容）
  skill import -url <url> [-name <n>]         从 URL/GitHub 下载导入 skill
  skill import -file <path> [-name <n>]       从本地文件导入 skill
  skill new    -name <n> -corpus <text|@file> 用语料生成 skill（触发 LLM）
  skill rm     -id <id>                        删除 skill
  chat         -skill <id> -m "<msg>" [-conv <id>]  与分身对话
  health                                       健康检查

Global:
  -server <url>   后端地址（默认 $UPME_SERVER 或 http://localhost:8000）
`)
}

// ---- skill subcommands ----

func skillCmd(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	action, rest := args[0], args[1:]

	switch action {
	case "list":
		fs := flag.NewFlagSet("skill list", flag.ExitOnError)
		server := fs.String("server", defaultServer(), "backend url")
		_ = fs.Parse(rest)
		printJSON(apiGet(*server, "/api/v1/skills"))

	case "get":
		fs := flag.NewFlagSet("skill get", flag.ExitOnError)
		server := fs.String("server", defaultServer(), "backend url")
		id := fs.String("id", "", "skill id")
		_ = fs.Parse(rest)
		requireFlag("id", *id)
		printJSON(apiGet(*server, "/api/v1/skills/"+*id))

	case "import":
		fs := flag.NewFlagSet("skill import", flag.ExitOnError)
		server := fs.String("server", defaultServer(), "backend url")
		url := fs.String("url", "", "source url (github blob/raw)")
		file := fs.String("file", "", "local SKILL.md path")
		name := fs.String("name", "", "skill name (optional)")
		_ = fs.Parse(rest)
		body := map[string]string{"name": *name}
		switch {
		case *url != "":
			body["url"] = *url
		case *file != "":
			b, err := os.ReadFile(*file)
			fatalIf(err)
			body["content"] = string(b)
		default:
			fmt.Fprintln(os.Stderr, "need -url or -file")
			os.Exit(2)
		}
		printJSON(apiPost(*server, "/api/v1/skills/import", body))

	case "new":
		fs := flag.NewFlagSet("skill new", flag.ExitOnError)
		server := fs.String("server", defaultServer(), "backend url")
		name := fs.String("name", "", "skill name")
		corpus := fs.String("corpus", "", "corpus text, or @file to read from file")
		_ = fs.Parse(rest)
		requireFlag("name", *name)
		requireFlag("corpus", *corpus)
		corpusText := *corpus
		if strings.HasPrefix(corpusText, "@") {
			b, err := os.ReadFile(corpusText[1:])
			fatalIf(err)
			corpusText = string(b)
		}
		created := apiPost(*server, "/api/v1/skills", map[string]string{"name": *name, "corpus": corpusText})
		var sk struct {
			ID string `json:"id"`
		}
		_ = json.Unmarshal(created, &sk)
		if sk.ID == "" {
			printJSON(created)
			return
		}
		fmt.Printf("created skill %s, processing...\n", sk.ID)
		printJSON(apiPost(*server, "/api/v1/skills/"+sk.ID+"/process", nil))

	case "rm":
		fs := flag.NewFlagSet("skill rm", flag.ExitOnError)
		server := fs.String("server", defaultServer(), "backend url")
		id := fs.String("id", "", "skill id")
		_ = fs.Parse(rest)
		requireFlag("id", *id)
		printJSON(apiDelete(*server, "/api/v1/skills/"+*id))

	default:
		fmt.Fprintf(os.Stderr, "unknown skill action: %s\n", action)
		os.Exit(2)
	}
}

func chatCmd(args []string) {
	fs := flag.NewFlagSet("chat", flag.ExitOnError)
	server := fs.String("server", defaultServer(), "backend url")
	skill := fs.String("skill", "", "skill id (system prompt)")
	msg := fs.String("m", "", "message")
	conv := fs.String("conv", "", "conversation id (optional)")
	_ = fs.Parse(args)
	requireFlag("m", *msg)

	resp := apiPost(*server, "/api/v1/agent/chat", map[string]string{
		"message":         *msg,
		"skill_id":        *skill,
		"conversation_id": *conv,
	})
	var r struct {
		Reply          string `json:"reply"`
		ConversationID string `json:"conversation_id"`
	}
	if json.Unmarshal(resp, &r) == nil && r.Reply != "" {
		fmt.Println(r.Reply)
		fmt.Fprintf(os.Stderr, "\n[conversation_id: %s]\n", r.ConversationID)
		return
	}
	printJSON(resp)
}

func healthCmd(args []string) {
	fs := flag.NewFlagSet("health", flag.ExitOnError)
	server := fs.String("server", defaultServer(), "backend url")
	_ = fs.Parse(args)
	printJSON(apiGet(*server, "/health"))
}

// ---- HTTP helpers ----

var httpClient = &http.Client{Timeout: 120 * time.Second}

func apiGet(server, path string) []byte    { return doReq(http.MethodGet, server+path, nil) }
func apiDelete(server, path string) []byte { return doReq(http.MethodDelete, server+path, nil) }

func apiPost(server, path string, body any) []byte {
	var buf io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewReader(b)
	}
	return doReq(http.MethodPost, server+path, buf)
}

func doReq(method, url string, body io.Reader) []byte {
	req, err := http.NewRequest(method, url, body)
	fatalIf(err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "HTTP %d: %s\n", resp.StatusCode, strings.TrimSpace(string(out)))
		os.Exit(1)
	}
	return out
}

func printJSON(b []byte) {
	var v any
	if json.Unmarshal(b, &v) == nil {
		pretty, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(pretty))
		return
	}
	fmt.Println(string(b))
}

func requireFlag(name, val string) {
	if val == "" {
		fmt.Fprintf(os.Stderr, "missing required flag: -%s\n", name)
		os.Exit(2)
	}
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
