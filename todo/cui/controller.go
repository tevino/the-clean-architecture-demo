package cui

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tevino/the-clean-architecture-demo/todo/cui/io"

	"github.com/tevino/the-clean-architecture-demo/todo/cui/component"

	ui "github.com/gizak/termui/v3"
	"github.com/tevino/the-clean-architecture-demo/todo/model"
	"github.com/tevino/the-clean-architecture-demo/todo/use"
)

// Controller is the controller of CUI.
type Controller struct {
	*CUI
	io.IO
	CasesTask use.CasesTask
}

func (c *Controller) handleCatListEvent(e component.TaskListEvent) {
	switch e.Type {
	case component.TaskListEventAfterUpdate:
		if t, ok := c.catList.GetSelectedTask(); ok {
			c.taskList.SetParentID(t.ID)
		}
		fallthrough
	default:
		c.handleGenericTaskListEvent(c.catList, e)
	}
}

func (c *Controller) handleTaskListEvent(e component.TaskListEvent) {
	c.handleGenericTaskListEvent(c.taskList, e)
}

func (c *Controller) handleGenericTaskListEvent(l component.TaskList, e component.TaskListEvent) {
	switch e.Type {
	case component.TaskListEventAfterUpdate:
		c.setDescriptionByCurrentSelectedRow(l)
	case component.EventInsertTaskWithOrder:
		c.insertTaskWithOrder(l, e.Order)
	case component.EventChangeTaskState:
		c.changeTaskState(l)
	}
}

func (c *Controller) changeTaskState(l component.TaskList) {
	t, ok := l.GetSelectedTask()
	if ok {
		err := c.CasesTask.ChangeTaskStateByID(t.ID, toggleCompletedState(t.State))
		if err != nil {
			c.stateBar.Warn(fmt.Errorf("changing task[%d] state: %w", t.ID, err))
		}
	}
}

func toggleCompletedState(s model.TaskState) model.TaskState {
	switch s {
	case model.TaskStateNormal:
		return model.TaskStateCompleted
	case model.TaskStateCompleted:
		return model.TaskStateNormal
	default:
		panic(fmt.Sprintf("unexpected task state: %v", s))
	}
}

func (c *Controller) setDescriptionByCurrentSelectedRow(l component.TaskList) {
	if !l.IsActivated() {
		return
	}
	t, ok := l.GetSelectedTask()
	if ok {
		c.descBox.Plain(t.Description)
	}
}

func (c *Controller) insertTaskWithOrder(l component.TaskList, order uint64) {
	defer func() {
		err := c.CUILib.Init()
		if err != nil {
			panic(fmt.Sprintf("failed to initialize CUI: %s", err))
		}
	}()
	buf, err := c.GetInputByLaunchingEditor()
	if err != nil {
		c.stateBar.Warn(err)
		return
	}
	form, err := createFormAddTaskFromString(buf)
	if err != nil {
		return
	}
	switch l {
	case c.catList:
		form.Type = model.TaskTypeCategory
	case c.taskList:
		form.Type = model.TaskTypeTask
	}
	form.ParentID = l.ParentID()
	form.Order = order
	err = c.CasesTask.AddTask(form)
	if err != nil {
		// TODO: launch editor again to let use re-edit the file.
		c.stateBar.Warn(fmt.Errorf("adding task: %w", err))
		return
	}
}

func (c *Controller) init() error {
	c.catList.SetEventHandler(c.handleCatListEvent)
	c.taskList.SetEventHandler(c.handleTaskListEvent)

	if err := c.CUILib.Init(); err != nil {
		return fmt.Errorf("initializing CUILib: %w", err)
	}

	c.stateBar.Plain("Good day!")
	termWidth, termHeight := c.CUILib.TerminalDimensions()
	c.grid.SetRect(0, 0, termWidth, termHeight)
	return nil
}

func (c *Controller) close() {
	c.CUILib.Close()
}

// Loop starts rendering the CUI.
func (c *Controller) Loop() error {
	err := c.init()
	if err != nil {
		return fmt.Errorf("initializing: %w", err)
	}
	defer c.close()
	uiEvents := c.CUILib.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			quit := c.handleEvent(e)
			if quit {
				return nil
			}
		default:
			err := c.Update()
			if err != nil {
				c.stateBar.Warn(err)
			}
		}
	}
}

func (c *Controller) handleEvent(e ui.Event) bool {
	switch e.ID {
	case "q", "<C-c>":
		return true
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		c.grid.SetRect(0, 0, payload.Width, payload.Height)
		c.CUILib.Clear()
	default:
		err := c.grid.HandleEvent(e)
		if err != nil {
			c.stateBar.Warn(err)
		}
	}
	return false
}

func (c *Controller) Update() error {
	c.CUILib.Render(c.grid)
	for _, r := range c.components {
		// render
		err := r.Update()
		if err != nil {
			return err
		}
		// Update tasks
		if l, ok := r.(component.TaskList); ok {
			err := c.CasesTask.ListTasksByParentID(l.ParentID())
			if err != nil {
				return fmt.Errorf("get tasks of parent[%d]: %w", l.ParentID(), err)
			}
		}
		c.CUILib.Render(r.(ui.Drawable))
	}
	return nil
}

var ErrInvalidDue = errors.New("invalid due")

func parseDue(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	var d time.Time
	switch s {
	case "", "today":
		d = time.Now()
	case "tom", "tomorrow":
		d = time.Now().Add(time.Hour * 24)
	default:
		return d, ErrInvalidDue
	}
	// TODO: use regexp to implement more, like +2d, next week, month
	return d, nil
}

var errEmptyInput = errors.New("empty input")

func createFormAddTaskFromString(s string) (*model.FormAddTask, error) {
	if s == "" {
		return nil, errEmptyInput
	}

	var title, desc string
	var due time.Time
	var skipDue = false
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		isLineEmpty := line == "" || trimmedLine == ""
		if title == "" && !isLineEmpty {
			title = trimmedLine
		} else if !skipDue && due.IsZero() && !isLineEmpty {
			parsedDue, err := parseDue(line)
			if err != nil {
				// if the second non-empty line after the title is not a valid due, we assume that no due provided.
				skipDue = true
				desc += line + "\n"
			} else {
				due = parsedDue
			}
		} else {
			desc += line + "\n"
		}
	}
	return &model.FormAddTask{
		Title:       title,
		Due:         due,
		Description: desc,
	}, nil
}
