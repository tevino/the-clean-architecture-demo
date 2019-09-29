package cui

import (
	"io"
	"testing"

	ui "github.com/gizak/termui/v3"

	"github.com/stretchr/testify/assert"

	"github.com/tevino/the-clean-architecture-demo/todo/cui/component/mock_component"

	"github.com/tevino/the-clean-architecture-demo/todo/cui/mock_cui"

	"github.com/tevino/the-clean-architecture-demo/todo/use/mock_use"

	"github.com/tevino/the-clean-architecture-demo/todo/cui/component"

	"github.com/tevino/the-clean-architecture-demo/todo/model"

	"github.com/golang/mock/gomock"
)

func newController(ctl *gomock.Controller) *Controller {
	lib := mock_cui.NewMockCUILib(ctl)
	return &Controller{
		CUI:       New(lib),
		CasesTask: mock_use.NewMockCasesTask(ctl),
		IO:        mock_cui.NewMockIO(ctl),
	}
}

func TestCatListAfterUpdateSetsTaskListParent(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockTaskList := mock_component.NewMockTaskList(ctl)
	mockCatList := mock_component.NewMockTaskList(ctl)
	task := &model.Task{ID: 42}

	c := newController(ctl)
	c.catList = mockCatList
	c.taskList = mockTaskList
	gomock.InOrder(
		mockCatList.EXPECT().GetSelectedTask().Return(task, true),
		mockTaskList.EXPECT().SetParentID(task.ID),
		mockCatList.EXPECT().IsActivated(),
	)
	c.handleCatListEvent(component.TaskListEvent{Type: component.TaskListEventAfterUpdate})
}

func TestInsertTaskWithOrder(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	// expectations
	var (
		parentID  int64  = 42
		order     uint64 = 4200
		taskInput        = "test title"
	)

	c := newController(ctl)

	mockList := mock_component.NewMockTaskList(ctl)
	c.taskList = mockList
	form, err := createFormAddTaskFromString(taskInput)
	assert.NoError(t, err)
	form.ParentID = parentID
	form.Order = order
	form.Type = model.TaskTypeTask
	gomock.InOrder(
		c.IO.(*mock_cui.MockIO).EXPECT().GetInputByLaunchingEditor().Return(taskInput, nil),
		mockList.EXPECT().ParentID().Return(parentID),
		c.CasesTask.(*mock_use.MockCasesTask).EXPECT().AddTask(form),
		c.CUILib.(*mock_cui.MockCUILib).EXPECT().Init(),
	)

	c.handleTaskListEvent(component.TaskListEvent{Type: component.EventInsertTaskWithOrder, Order: order})
}

func TestChangeStateEventHandled(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockText := mock_component.NewMockText(ctl)
	mockList := mock_component.NewMockTaskList(ctl)
	c := newController(ctl)
	c.stateBar = mockText
	task := &model.Task{}
	// OK
	gomock.InOrder(
		mockList.EXPECT().GetSelectedTask().Return(task, true),
		c.CasesTask.(*mock_use.MockCasesTask).EXPECT().ChangeTaskStateByID(task.ID, gomock.Any()).Return(nil),
	)
	c.changeTaskState(mockList)
	// Error
	gomock.InOrder(
		mockList.EXPECT().GetSelectedTask().Return(task, true),
		c.CasesTask.(*mock_use.MockCasesTask).EXPECT().ChangeTaskStateByID(task.ID, gomock.Any()).Return(io.EOF),
		mockText.EXPECT().Warn(gomock.Any()),
	)
	c.changeTaskState(mockList)
}

func TestToggleCompletedState(t *testing.T) {
	t.Parallel()
	normal := model.TaskStateNormal
	completed := model.TaskStateCompleted
	assert.Equal(t, completed, toggleCompletedState(normal))
	assert.Equal(t, normal, toggleCompletedState(completed))

	assert.Panics(t, func() {
		toggleCompletedState(99)
	})
}

func TestSetDescriptionByCurrentSelectedRow(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	c := newController(ctl)
	mockList := mock_component.NewMockTaskList(ctl)
	gomock.InOrder(
		mockList.EXPECT().IsActivated().Return(true),
		mockList.EXPECT().GetSelectedTask().Return(&model.Task{}, true),
	)
	c.setDescriptionByCurrentSelectedRow(mockList)
}

func TestControllerLoop(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	c := newController(ctl)
	gomock.InOrder(
		c.CUILib.(*mock_cui.MockCUILib).EXPECT().Init(),
		c.CUILib.(*mock_cui.MockCUILib).EXPECT().TerminalDimensions(),
		c.CUILib.(*mock_cui.MockCUILib).EXPECT().PollEvents().Return(uiEvents("q")),
		c.CUILib.(*mock_cui.MockCUILib).EXPECT().Close(),
	)
	ended := make(chan bool, 1)
	go func() {
		assert.NoError(t, c.Loop())
		ended <- true
	}()
	<-ended
}

func uiEvents(keys ...string) <-chan ui.Event {
	ch := make(chan ui.Event, len(keys)+1)
	for _, k := range keys {
		ch <- ui.Event{ID: k}
	}
	return ch
}
