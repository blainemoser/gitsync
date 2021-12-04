package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Git struct {
	repo string
	cwd  string
}

// NewGit creates a new instance of Git
func NewGit() *Git {
	return &Git{}
}

// SetRepo sets this instance's repo
func (g *Git) SetRepo(repo string) (*Git, error) {
	g.repo = repo
	err := g.checkRepo()
	return g, err
}

// Sync syncs the current repo by running pull then push
func (g *Git) Sync() (string, error) {
	stage, stageErr := g.Stage()
	commit, commitErr := g.Commit()
	pull, pullErr := g.Pull()
	push, pushErr := g.Push()
	result := fmt.Sprintf(
		"Stage Result:\n%s\nCommit Result:\n%s\nPull Result:\n%s\nPush Result:\n%s\n",
		stage,
		commit,
		push,
		pull,
	)
	return result, ParseErrors([]error{
		stageErr,
		commitErr,
		pullErr,
		pushErr,
	})
}

// Stage runs git stage
func (g *Git) Stage() (string, error) {
	return g.action([]string{"stage", "."})
}

// Commit runs git commit
func (g *Git) Commit() (string, error) {
	return g.action([]string{"commit", ".", "-m", g.commitMessage("Automatic Commit by GitSync")})
}

// Commit runs git push
func (g *Git) Push() (string, error) {
	return g.action([]string{"push"})
}

// Commit runs git pull
func (g *Git) Pull() (string, error) {
	return g.action([]string{"pull"})
}

func (g *Git) action(args []string) (string, error) {
	err := os.Chdir(g.repo)
	if err != nil {
		return "", err
	}
	defer g.back()
	result, err := exec.Command("git", args...).CombinedOutput()
	return string(result), err
}

func (g *Git) commitMessage(message string) string {
	timeStamp := time.Now()
	return fmt.Sprintf("[%s]: %s [%d]", timeStamp.Format(time.RFC3339), message, timeStamp.Unix())
}

func (g *Git) back() {
	err := os.Chdir(g.cwd)
	if err != nil {
		panic(err)
	}
}

func (g *Git) init() error {
	result, err := g.action([]string{"init"})
	fmt.Println(result)
	return err
}

func (g *Git) checkRepo() error {
	err := g.setCWD()
	if err != nil {
		return err
	}
	isRepo, err := g.isRepo()
	if err != nil {
		return err
	}
	if !isRepo {
		return g.init()
	}
	return nil
}

func (g *Git) setCWD() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	g.cwd = cwd
	return nil
}

func (g *Git) isRepo() (bool, error) {
	result, err := g.action([]string{"status"})
	if err != nil {
		if len(result) > 0 && strings.Contains(result, "fatal: not a git repository") {
			return false, nil
		}
		return false, err
	}
	return !strings.Contains(result, "fatal: not a git repository"), nil
}

func (g *Git) GetRepo() string {
	return g.repo
}
