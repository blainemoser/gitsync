package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/blainemoser/gitsync/logging"
	"github.com/blainemoser/gitsync/utils"
	"github.com/fsnotify/fsnotify"
)

type Git struct {
	repo    string
	cwd     string
	watcher *fsnotify.Watcher
}

// NewGit creates a new instance of Git
func NewGit() *Git {
	return &Git{}
}

func (g *Git) SetWatcher() error {
	err := g.newWatcher()
	if err != nil {
		return err
	}
	errs := make([]error, 0)
	if err := filepath.WalkDir(g.GetRepo(), g.watchDir); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return utils.ParseErrors(errs)
	}
	return nil
}

// Watcher returns the watcher instance
func (g *Git) Watcher() *fsnotify.Watcher {
	return g.watcher
}

func (g *Git) HandleEvent(errChan chan error, resultChan chan string, event fsnotify.Event) {
	if strings.Contains(event.Name, ".git") {
		resultChan <- ""
		errChan <- nil
		return
	}
	errs := make([]error, 0)
	result, err := g.Sync()
	if err != nil {
		errs = append(errs, err)
	}
	err = g.SetWatcher()
	if err != nil {
		errs = append(errs, err)
	}
	resultChan <- result
	errChan <- utils.ParseErrors(errs)
}

func (g *Git) newWatcher() error {
	if g.watcher != nil {
		g.watcher.Close()
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	g.watcher = watcher
	return nil
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func (g *Git) watchDir(path string, di os.DirEntry, err error) error {
	if strings.Contains(path, ".git") {
		return nil
	}
	if di.IsDir() {
		return g.watcher.Add(path)
	}

	return nil
}

// SetRepo sets this instance's repo
func (g *Git) SetRepo(repo string) (*Git, error) {
	g.repo = repo
	err := g.checkRepo()
	if err != nil {
		return g, err
	}
	err = g.SetWatcher()
	return g, err
}

// Status runs git status
func (g *Git) Status() (string, error) {
	return g.action([]string{"status"})
}

func (g *Git) HasChanges() (bool, error, string) {
	status, err := g.Status()
	if err != nil {
		return false, err, status
	}
	return strings.Contains(status, "Changes not staged for commit") ||
		strings.Contains(status, "Untracked files"), nil, status
}

// Sync syncs the current repo by running stage, commit, pull then finally push
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
	return result, utils.ParseErrors([]error{
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
	logging.StaticWrite(result, "INFO")
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
