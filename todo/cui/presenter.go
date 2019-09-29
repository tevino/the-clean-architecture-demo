package cui

import (
	"fmt"

	"github.com/tevino/the-clean-architecture-demo/todo/model"
)

type Presenter struct {
	*CUI
}

func (p *Presenter) ShowTasksOfParentID(parentID int64, tasks []*model.Task) error {
	var err error
	switch parentID {
	case p.catList.ParentID():
		p.catList.UpdateTasks(tasks)
	case p.taskList.ParentID():
		p.taskList.UpdateTasks(tasks)
	}
	return err
}

func (p *Presenter) ShowTaskAdded(task *model.Task) error {
	p.stateBar.Info(fmt.Sprintf("Task Added: %s", task.Title))
	return nil
}
