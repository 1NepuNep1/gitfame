package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func ListFiles(repoPath, revision string) ([]string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", revision)
	cmd.Dir = repoPath
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git ls-tree failed: %w (stderr: %s)", err, stderr.String())
	}

	lines := strings.Split(stdout.String(), "\n")
	var files []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		files = append(files, line)
	}

	return files, nil

}
