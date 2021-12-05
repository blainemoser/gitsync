package queue

import (
	"fmt"

	"github.com/blainemoser/gitsync/configs"
	"github.com/blainemoser/gitsync/repo"
	"github.com/blainemoser/gitsync/utils"
)

type Process struct {
	ready bool
	git   *repo.Git
	state interface{}
}

type Queue []*Process

// NewQueue returns a new instance of Queue
func NewProcess(git *repo.Git) *Process {
	return &Process{
		ready: true,
		git:   git,
		state: nil,
	}
}

// NewQueue returns a new *Queue
func NewQueue() *Queue {
	return &Queue{}
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
		}
		q.push(NewProcess(git))
	}
	if len(errs) > 0 {
		return q, utils.ParseErrors(errs)
	}
	return q, nil
}

// Git returns the Process' Git
func (p *Process) Git() *repo.Git {
	return p.git
}

// Process runs a queue process
func (p *Process) Process(result chan interface{}) {
	if !p.ready {
		result <- "running"
		return
	}
	p.ready = false
	p.state = p.updateState()
	p.ready = true
	result <- p.state
}

// ErrorState checks whether a returned state is an error
func (p *Process) GetState() (string, error) {
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

// ErrorState checks whether a returned state is an error
func (p *Process) ResultState(state interface{}) error {
	if err, ok := state.(error); ok {
		return err
	}
	return nil
}

func (q *Queue) push(p *Process) *Queue {
	*q = append(*q, p)
	return q
}

func (p *Process) updateState() interface{} {
	changed, err := p.Git().HasChanges()
	if err != nil {
		return err
	}
	if !changed {
		return "no changes"
	}
	sync, syncErr := p.Git().Sync()
	if syncErr != nil {
		return utils.ParseErrors([]error{
			fmt.Errorf(sync),
			syncErr,
		})
	} else {
		return sync
	}
}
