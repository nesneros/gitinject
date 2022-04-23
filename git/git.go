package git

import (
	"os/exec"
	"strings"
)

var gitBin = "git"

type GitInfo struct {
	tag string
	sha string
}

func GitSha(repo string) (string, error) {
	return trim(execGit("-C", repo, "rev-parse", "--verify", "HEAD"))
}

func GitTag(repo string) (string, error) {
	return trim(execGit("-C", repo, "describe", "--tags", "HEAD"))
}

func trim(s string, err error) (string, error) {
	if err != nil {
		return s, err
	}
	return strings.TrimSpace(s), nil
}

func execGit(args ...string) (string, error) {
	cmd := exec.Command(gitBin, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
