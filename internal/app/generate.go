package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"gxcommit/internal/config"
)

// Commit represents a logical git commit suggested by the LLM.
type Commit struct {
	Files   []string `json:"files"`
	Message string   `json:"message"`
}

// GenerateScript reads git diff, sends it to Groq,
// parses the response, and generates a bash script.
func GenerateScript(jira string) string {
	diff, err := getDiff()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("No changes to commit")
		return ""
	}

	commits, err := generateCommits(diff)
	if err != nil {
		fmt.Println("Groq error:", err)
		return ""
	}

	if len(commits) == 0 {
		fmt.Println("No commits generated")
		return ""
	}

	return generateScript(commits, jira)
}

// ---------------- Git helpers ----------------

func getDiff() (string, error) {
	cmd := exec.Command("git", "diff", "HEAD")
	out, err := cmd.Output()
	if err == nil {
		return string(out), nil
	}

	// Handle initial commit
	if exec.Command("git", "rev-parse", "HEAD").Run() != nil {
		untracked, _ := exec.Command(
			"git", "ls-files", "--others", "--exclude-standard",
		).Output()

		if len(untracked) == 0 {
			return "", nil
		}

		return "Initial commit with new files:\n" + string(untracked), nil
	}

	return "", fmt.Errorf("failed to read git diff")
}

// ---------------- Groq integration ----------------

func generateCommits(diff string) ([]Commit, error) {
	apiKey, _ := config.Get("groq", "GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY is not set in ~/.gxconfig. Run: gxcommit config set GROQ_API_KEY=<your-key>")
	}

	prompt := `Analyze the following git diff and suggest multiple logical commits.

Rules:
- Group related changes
- Assume full files (no partial hunks)
- Use conventional commit messages
- Return ONLY valid JSON

Conventional commit types:
- feat: A new feature
- fix: A bug fix
- docs: Documentation only changes
- style: Code style changes (formatting, whitespace)
- refactor: Code refactoring
- perf: Performance improvements
- test: Adding or updating tests
- build: Build system changes
- ci: CI configuration changes
- chore: Maintenance tasks
- revert: Revert a previous commit

JSON format:
{
  "commits": [
    {
      "files": ["file1.go", "file2.go"],
      "message": "feat: add feature"
    }
  ]
}

Diff:
` + diff

	payload := map[string]any{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Return ONLY valid JSON. No markdown. No explanations.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.6,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("groq http %d: %s", resp.StatusCode, string(b))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseGroqResponse(raw)
}

func parseGroqResponse(raw []byte) ([]Commit, error) {
	var envelope struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, err
	}

	if len(envelope.Choices) == 0 {
		return nil, fmt.Errorf("empty groq response")
	}

	content := envelope.Choices[0].Message.Content

	var parsed struct {
		Commits []Commit `json:"commits"`
	}

	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("invalid JSON from groq: %w\n%s", err, content)
	}

	return parsed.Commits, nil
}

// ---------------- Script generation ----------------

func generateScript(commits []Commit, jira string) string {
	var b strings.Builder

	b.WriteString("#!/bin/bash\n\n")
	b.WriteString("set -e\n\n")

	for _, c := range commits {
		if len(c.Files) == 0 || c.Message == "" {
			continue
		}

		b.WriteString("git add ")
		b.WriteString(strings.Join(c.Files, " "))
		b.WriteString("\n")

		msg := c.Message
		if jira != "" {
			msg = fmt.Sprintf("[%s] %s", jira, msg)
		}

		b.WriteString(fmt.Sprintf("git commit -m \"%s\"\n\n", escape(msg)))
	}

	return b.String()
}

func escape(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
