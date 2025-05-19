package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type BlameResult struct {
	mu      sync.Mutex
	Commits map[string]*CommitInfo
}

type CommitInfo struct {
	Author    string
	LineCount int
	Files     map[string]bool
}

func BlameFiles(files []string, repoPath, revision string, useCommitter bool, errOut io.Writer) (*BlameResult, error) {
	result := &BlameResult{
		Commits: make(map[string]*CommitInfo),
	}

	const workerCount = 8
	fileCh := make(chan string, len(files))
	errCh := make(chan error, workerCount)

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for range workerCount {
		go func() {
			defer wg.Done()
			for file := range fileCh {
				if err := blameOneFile(file, repoPath, revision, useCommitter, result); err != nil {
					errCh <- fmt.Errorf("blame failed for file %s: %w", file, err)
					return
				}
			}
		}()
	}

	for _, f := range files {
		fileCh <- f
	}
	close(fileCh)

	wg.Wait()
	close(errCh)

	for e := range errCh {
		return nil, e
	}

	return result, nil
}

func blameOneFile(file, repoPath, revision string, useCommitter bool, result *BlameResult) error {
	cmd := exec.Command("git", "blame", "--porcelain", revision, "--", file)
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git blame error: %v (stderr: %s)", err, stderr.String())
	}

	if stdout.Len() == 0 {
		return handleEmptyFile(repoPath, revision, file, useCommitter, result)
	}

	return parsePorcelainBlame(stdout.Bytes(), file, useCommitter, result)
}

func handleEmptyFile(repoPath, revision, file string, useCommitter bool, result *BlameResult) error {
	commitHash, err := lastCommitForFile(repoPath, revision, file)
	if err != nil {
		return fmt.Errorf("failed to find last commit for empty file %q: %w", file, err)
	}
	if commitHash == "" {
		return nil
	}

	authorName, err := getAuthorOrCommitter(repoPath, commitHash, useCommitter)
	if err != nil {
		return fmt.Errorf("failed to get author for commit %s (file %s): %w", commitHash, file, err)
	}
	addCommitInfo(result, commitHash, authorName, file, 0)
	return nil
}

func lastCommitForFile(repoPath, revision, file string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%H", revision, "--", file)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getAuthorOrCommitter(repoPath, commitHash string, useCommitter bool) (string, error) {
	var formatArg string
	if useCommitter {
		formatArg = "--format=%cn"
	} else {
		formatArg = "--format=%an"
	}
	cmd := exec.Command("git", "show", "-s", formatArg, commitHash)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func parsePorcelainBlame(data []byte, file string, useCommitter bool, result *BlameResult) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))

	ccommitAmount := make(map[string]int)
	commitAuthor := make(map[string]string)

	var whoTag string
	if useCommitter {
		whoTag = "committer "
	} else {
		whoTag = "author "
	}

	var currentHash string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) == 4 && isSha(fields[0]) {
			currentHash = fields[0]
			countStr := fields[3]

			count, err := strconv.Atoi(countStr)
			if err == nil {
				ccommitAmount[currentHash] += count
			}
			continue
		}

		if strings.HasPrefix(line, whoTag) {
			if currentHash != "" {
				if _, ok := commitAuthor[currentHash]; !ok {
					name := strings.TrimPrefix(line, whoTag)
					commitAuthor[currentHash] = name
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for sha, linesCount := range ccommitAmount {
		authorName := commitAuthor[sha]
		addCommitInfo(result, sha, authorName, file, linesCount)
	}

	return nil
}

func addCommitInfo(result *BlameResult, commitHash, authorName, file string, lines int) {
	if commitHash == "" {
		return
	}
	result.mu.Lock()
	defer result.mu.Unlock()

	ci, ok := result.Commits[commitHash]
	if !ok {
		ci = &CommitInfo{
			Author:    authorName,
			LineCount: 0,
			Files:     make(map[string]bool),
		}
		result.Commits[commitHash] = ci
	}
	ci.LineCount += lines
	ci.Files[file] = true
}

func isSha(s string) bool {
	if len(s) < 4 {
		return false
	}
	for _, r := range s {
		if !((r >= '0' && r <= '9') ||
			(r >= 'a' && r <= 'f') ||
			(r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}
