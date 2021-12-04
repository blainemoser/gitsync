package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Git struct {
	repo string
	cwd  string
}

// NewGit creates a new instance of Git
func NewGit() *Git {
	return &Git{}
}

func (g *Git) SetRepo(repo string) (*Git, error) {
	g.repo = repo
	err := g.checkRepo()
	return g, err
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
