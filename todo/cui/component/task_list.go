package component

import (
	"fmt"

	"github.com/dustin/go-humanize"

	"github.com/tevino/the-clean-architecture-demo/todo/use"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/tevino/the-clean-architecture-demo/todo/model"
)

//go:generate mockgen -destination mock_component/task_list_mock.go github.com/tevino/the-clean-architecture-demo/todo/cui/component TaskList

// TaskList represents a list of tasks.
type TaskList interface {
	InteractiveComponent
	ParentID() int64
	SetParentID(parentID int64)
	UpdateTasks(tasks []*model.Task)
	GetSelectedTask() (*model.Task, bool)
	SetEventHandler(func(TaskListEvent))
}

// TaskListComponent displays a list of tasks.
type TaskListComponent struct {
	*widgets.List
	parentID int64
	tasks    []*model.Task
	use.CasesTask
	previousKey string
	isActivated bool
	handleEvent func(TaskListEvent)
}

// NewListComponent creates a TaskListComponent with given title.
func NewListComponent(title string) *TaskListComponent {
	list := widgets.NewList()
	if title != "" {
		list.Title = title
		list.TitleStyle.Modifier = ui.ModifierBold
	}
	return &TaskListComponent{
		List:        list,
		isActivated: false,
		handleEvent: func(TaskListEvent) {},
	}
}

// HandleEvent handles keyboard events.
func (l *TaskListComponent) HandleEvent(e ui.Event) error {
	switch e.ID {
	case "j", "<Down>":
		l.selectTaskAt(l.SelectedRow + 1)
	case "k", "<Up>":
		l.selectTaskAt(l.SelectedRow - 1)
	// case "<C-d>":
	// 	l.ScrollHalfPageDown()
	// case "<C-u>":
	// 	l.ScrollHalfPageUp()
	// case "<C-f>":
	// 	l.ScrollPageDown()
	// case "<C-b>":
	// 	l.ScrollPageUp()
	case "g":
		if l.previousKey != "g" {
			break
		}
		l.previousKey = ""
		fallthrough
	case "<Home>":
		// l.ScrollTop()
		l.selectTaskAt(0)
	case "G", "<End>":
		// l.ScrollBottom()
		l.selectTaskAt(len(l.Rows) - 1)
	case "o", "O":
		var order uint64
		switch e.ID {
		case "o":
			order = 1
		case "O":
			order = 0
		}
		t, ok := l.GetSelectedTask()
		if ok {
			order += t.Order
		} else {
			order = 0
		}
		l.handleEvent(TaskListEvent{Type: EventInsertTaskWithOrder, Order: order})
	case "<Space>":
		l.handleEvent(TaskListEvent{Type: EventChangeTaskState})
	}

	l.previousKey = e.ID
	return nil
}

func (l *TaskListComponent) SetEventHandler(handle func(TaskListEvent)) {
	l.handleEvent = handle
}

func (l *TaskListComponent) selectTaskAt(idx int) {
	if idx >= 0 && idx < len(l.Rows) {
		l.SelectedRow = idx
	}
}

func (l *TaskListComponent) GetSelectedTask() (*model.Task, bool) {
	if l.SelectedRow < len(l.tasks) {
		return l.tasks[l.SelectedRow], true
	}
	return nil, false
}

// SetActivate highlights selected row.
func (l *TaskListComponent) SetActivate(yes bool) {
	l.isActivated = yes
	modifier := ui.ModifierUnderline
	if yes {
		modifier = ui.ModifierReverse | ui.ModifierBold
	}
	l.SelectedRowStyle.Modifier = modifier
}

func (l *TaskListComponent) IsActivated() bool {
	return l.isActivated
}

func formatTaskRow(t *model.Task, width int) string {
	var row string
	switch t.Type {
	case model.TaskTypeCategory:
		row = fmt.Sprintf("+ %s", t.Title)
	case model.TaskTypeTask:
		var x string
		switch t.State {
		case model.TaskStateCompleted:
			x = "[x]"
		case model.TaskStateNormal:
			x = "[ ]"

		}
		due := humanize.Time(t.Due)
		// the 1s are the count of spaces in the formatting string
		titleLength := width - len(x) - 1 - len(due) - 1
		title := t.Title
		if len(title) > titleLength {
			title = title[:titleLength-3]
			title = title + "..."
		}
		format := fmt.Sprintf("%%s %%-%ds %%10s", titleLength)
		row = fmt.Sprintf(format, x, title, due)
	}

	return row
}

var emptyRows = []string{"<Empty>"}

func (l *TaskListComponent) Update() error {
	rowMatched := 0
	rows := make([]string, len(l.tasks))

	for i, t := range l.tasks {
		rows[i] = formatTaskRow(t, l.Inner.Dx())
		if selectedTask, ok := l.GetSelectedTask(); ok {
			if t.ID == selectedTask.ID {
				rowMatched = i
			}
		}
	}
	l.SelectedRow = rowMatched
	if len(rows) == 0 {
		rows = emptyRows
	}
	l.Rows = rows
	l.handleEvent(TaskListEvent{Type: TaskListEventAfterUpdate})
	return nil
}

// ParentID returns the ParentID of all tasks within the List.
func (l *TaskListComponent) ParentID() int64 {
	return l.parentID
}

// SetParentID sets the parentID of this list.
func (l *TaskListComponent) SetParentID(parentID int64) {
	l.parentID = parentID
}

// UpdateTasks replaces tasks displayed with given slice.
func (l *TaskListComponent) UpdateTasks(tasks []*model.Task) {
	l.tasks = tasks
}
