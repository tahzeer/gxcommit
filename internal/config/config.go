package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const configFileName = ".gxconfig"

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}

func Set(section, key, value string) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	var content string
	if data, err := os.ReadFile(configPath); err == nil {
		content = string(data)
	}

	lines := []string{}
	inSection := false
	sectionFound := false
	keyFound := false
	targetLine := fmt.Sprintf("\t%s=%s", key, value)

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if inSection && !keyFound {
				lines = append(lines, targetLine)
			}
			inSection = (trimmed == "["+section+"]")
			if inSection {
				sectionFound = true
			}
			lines = append(lines, line)
			continue
		}

		if inSection && strings.HasPrefix(trimmed, key+"=") {
			lines = append(lines, targetLine)
			keyFound = true
			continue
		}

		lines = append(lines, line)
	}

	if inSection && !keyFound {
		lines = append(lines, targetLine)
	}

	if !sectionFound {
		if content != "" {
			lines = append(lines, "")
		}
		lines = append(lines, "["+section+"]")
		lines = append(lines, targetLine)
	}

	final := strings.Join(lines, "\n")
	if !strings.HasSuffix(final, "\n") {
		final += "\n"
	}

	return os.WriteFile(configPath, []byte(final), 0644)
}

func Get(section, key string) (string, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", nil
	}

	inSection := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inSection = (trimmed == "["+section+"]")
			continue
		}

		if inSection && strings.HasPrefix(trimmed, key+"=") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				return parts[1], nil
			}
		}
	}

	return "", nil
}
