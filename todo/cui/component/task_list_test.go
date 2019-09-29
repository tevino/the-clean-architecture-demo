package component

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tevino/the-clean-architecture-demo/todo/model"

	ui "github.com/gizak/termui/v3"
)

func TestTaskListComponentHandleInsertTaskEvent(t *testing.T) {
	tasks := []*model.Task{
		{Order: 1},
		{Order: 2},
	}

	for _, c := range []struct {
		selectedRow int
		key         string
		expectOrder uint64
	}{
		{selectedRow: 0, key: "o", expectOrder: tasks[0].Order + 1},
		{selectedRow: 0, key: "O", expectOrder: tasks[0].Order},
		{selectedRow: 1, key: "o", expectOrder: tasks[1].Order + 1},
		{selectedRow: 1, key: "O", expectOrder: tasks[1].Order},
	} {
		run := 0
		l := NewListComponent("")
		l.tasks = tasks
		l.SetEventHandler(func(e TaskListEvent) {
			assert.Equal(t, c.expectOrder, e.Order)
			run += 1
		})
		l.SelectedRow = c.selectedRow
		err := l.HandleEvent(ui.Event{ID: c.key})
		assert.NoError(t, err)
		assert.Equal(t, 1, run)
	}
}

func TestSelectedTask(t *testing.T) {
	for _, c := range []struct {
		tasks           []*model.Task
		selectedRow     int
		expectTaskIndex int
	}{
		{selectedRow: 0, expectTaskIndex: 0, tasks: []*model.Task{{ID: 1}}},
		{selectedRow: 1, expectTaskIndex: 1, tasks: []*model.Task{{ID: 1}, {ID: 2}}},
		{selectedRow: 1, expectTaskIndex: -1, tasks: []*model.Task{}},
		{selectedRow: 1, expectTaskIndex: -1, tasks: nil},
		{selectedRow: 4, expectTaskIndex: -1, tasks: []*model.Task{{ID: 1}, {ID: 2}}},
	} {
		l := NewListComponent("")
		l.tasks = c.tasks
		l.SelectedRow = c.selectedRow
		task, ok := l.GetSelectedTask()
		if c.expectTaskIndex >= 0 {
			assert.True(t, ok)
			assert.Equal(t, l.tasks[c.expectTaskIndex], task)
		} else {
			assert.False(t, ok)
		}
	}
}
