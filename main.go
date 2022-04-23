package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nesneros/gitinject/gitx"
	"golang.org/x/mod/semver"
)

//go:generate go run main.go -cmd=gen

//go:embed .gitinject/*
var gitInjectFs embed.FS

const (
	genDirDefault = ".gitinject"
)

func readGitInfo(fs embed.FS, fallbackVer string) *gitInfo {
	v, _ := fs.ReadFile("version")
	ver := string(v)
	if ver == "" {
		ver = fallbackVer
	}
	sha, _ := fs.ReadFile("sha")
	return &gitInfo{ver: ver, sha: string(sha)}
}

var GitInfo *gitInfo

func usageError(errMsg string) {
	fmt.Fprintf(os.Stderr, "%s\n", errMsg)
	usage()
	os.Exit(1)
}

func reportError(err error) {
	if err == nil {
		return
	}
	usageError(err.Error())
}

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Version: %s (sha: %s)\nUsage of %s:\n\nOptions:\n", GitInfo.ver, GitInfo.sha, os.Args[0])
	flag.PrintDefaults()
}

func getGitInfo(repo string, fallback string) *gitInfo {
	result := gitInfo{ver: fallback}
	sha, err := gitx.GitSha(repo)
	reportError(err)
	result.sha = sha
	tag, err := gitx.GitTag(repo)
	if err == nil && semver.IsValid(tag) {
		result.ver = tag
	}
	return &result
}

type gitInfo struct {
	ver string
	sha string
}

// Example how to resolve a revision into its commit counterpart
func main() {
	GitInfo = readGitInfo(gitInjectFs, "<dev>")
	cmd := flag.String("cmd", "help", "Command to execute")
	repo := flag.String("repo", ".", "Git repository. Default to current directory")
	genDir := flag.String("genDir", genDirDefault, "Directory to generate files with git info. Default is "+genDirDefault)

	flag.Usage = usage
	flag.Parse()

	switch *cmd {
	case "help":
		usage()
	case "show":
		gitInfo := getGitInfo(*repo, "<dev>")
		fmt.Printf("sha: %s\nver: %s\n", gitInfo.sha, gitInfo.ver)
	case "gen":
		gitInfo := getGitInfo(*repo, "<dev>")
		generate(gitInfo, *genDir)
	default:
		usageError("Invalid command: " + *cmd)
	}
}

func generate(gitInfo *gitInfo, genDir string) {
	if strings.HasPrefix(genDir, "/") || genDir == "" {
		usageError("Invalid gendir: " + genDir)
	}
	shaFile := genDir + "/sha"
	verFile := genDir + "/version"

	err := os.MkdirAll(genDir, 0755)
	reportError(err)
	err = os.WriteFile(shaFile, []byte(gitInfo.sha), 0644)
	reportError(err)
	err = os.WriteFile(verFile, []byte(gitInfo.ver), 0644)
	reportError(err)
}
