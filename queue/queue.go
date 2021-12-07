package queue

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/blainemoser/gitsync/configs"
	"github.com/blainemoser/gitsync/logging"
	"github.com/blainemoser/gitsync/repo"
	"github.com/blainemoser/gitsync/utils"
	"github.com/fsnotify/fsnotify"
)

type Process struct {
	ready bool
	git   *repo.Git
	state interface{}
	cond  *sync.Cond
	log   *logging.Log
	wait  chan struct{}
}

type Queue []*Process

// NewQueue returns a new instance of Queue
func NewProcess(git *repo.Git) (*Process, error) {
	mutex := &sync.Mutex{}
	log, err := logging.NewLog()
	if err != nil {
		return nil, err
	}
	return &Process{
		ready: false,
		git:   git,
		state: nil,
		cond:  sync.NewCond(mutex),
		log:   log,
		wait:  make(chan struct{}, 1),
	}, nil
}

// NewQueue returns a new *Queue
func NewQueue() *Queue {
	return &Queue{}
}

// Git returns the Process' Git
func (p *Process) Git() *repo.Git {
	return p.git
}

func (p *Process) Log(message, level string) {
	err := p.log.Write(message, level)
	if err != nil {
		log.Println(err.Error())
	}
}

// Walk goes through the config's directories and adds each as Git instance to the Queue
func (q *Queue) Walk(c *configs.Configs) (*Queue, error) {
	var git *repo.Git
	var err error
	errs := make([]error, 0)
	for _, dir := range c.GetDirectories() {
		git, err = repo.NewGit().SetRepo(dir)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		process, err := NewProcess(git)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		q.push(process)
	}
	if len(errs) > 0 {
		return q, utils.ParseErrors(errs)
	}
	return q, nil
}

// StandbyAll sets all repos to watch for changes
func (q *Queue) StandbyAll() {
	for _, process := range *q {
		process.Standby()
	}
}

func (q *Queue) WaitForAll() {
	for _, process := range *q {
		<-process.wait
	}
}

func (p *Process) Standby() {
	go p.awaitEvent()
	p.ready = true
}

func (q *Queue) ResetWaits() {
	for _, process := range *q {
		process.ResetWait()
	}
}

func (p *Process) ResetWait() {
	close(p.wait)
	p.wait = make(chan struct{}, 1)
}

func (p *Process) pushWait() {
	if len(p.wait) > 0 {
		p.ResetWait()
	}
	p.wait <- struct{}{}
}

func (p *Process) awaitEvent() {
	var result string
	var err error
	p.Log(fmt.Sprintf("%s is waiting", p.Git().GetRepo()), "INFO")
	select {
	case event := <-p.Git().Watcher().Events:
		p.ready = false
		p.Log(fmt.Sprintf("%s: %s", event.Name, event.Op.String()), "INFO")
		result, err = p.process(event)
		p.completeProcess(err, result)
	case err = <-p.Git().Watcher().Errors:
		p.completeProcess(err, "")
	}
}

func (p *Process) completeProcess(err error, result string) {
	if err != nil {
		p.Log(err.Error(), "ERROR")
	}
	if len(result) > 0 {
		p.Log(result, "INFO")
	}
	p.cond.Signal()
	p.pushWait()
	p.Standby()
}

// Process runs a queue process
func (p *Process) process(event fsnotify.Event) (string, error) {
	changed, err, state := p.Git().HasChanges()
	p.Log("checked state: "+strings.ReplaceAll(state, "\n", " "), "INFO")
	if err != nil {
		p.Log(err.Error(), "ERROR")
		return "", err
	}
	if changed {
		p.syncEvent(event)
	}
	return p.getState()
}

func (p *Process) syncEvent(event fsnotify.Event) {
	errChan := make(chan error, 1)
	resultChan := make(chan string, 1)
	p.Git().HandleEvent(errChan, resultChan, event)
	err := <-errChan
	result := <-resultChan
	p.Log(result, "INFO")
	if err != nil {
		p.Log(err.Error(), "ERROR")
	}
	close(errChan)
	close(resultChan)
}

func (p *Process) getState() (string, error) {
	if p.state == nil {
		return "", nil
	}
	if err, ok := p.state.(error); ok {
		return "", err
	} else if state, ok := p.state.(string); ok {
		return state, nil
	}
	return "", nil
}

func (q *Queue) push(p *Process) *Queue {
	*q = append(*q, p)
	return q
}
