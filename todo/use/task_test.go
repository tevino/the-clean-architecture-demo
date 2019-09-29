package use

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tevino/the-clean-architecture-demo/todo/model"
	"github.com/tevino/the-clean-architecture-demo/todo/use/mock_use"

	"github.com/tevino/the-clean-architecture-demo/todo/entity"

	"github.com/golang/mock/gomock"
)

var i64 int64

func newTask(ctl *gomock.Controller) *TaskInteractor {
	return &TaskInteractor{
		Presenter: mock_use.NewMockPresenter(ctl),
		Storage:   mock_use.NewMockStorage(ctl),
	}
}

type itemMatcher struct {
	item entity.Item
}

// Matches returns whether x is a match.
func (m itemMatcher) Matches(x interface{}) bool {
	target := x.(*entity.Item)
	// skip the comparision of created and updated at
	m.item.CreatedAt = target.CreatedAt
	m.item.UpdatedAt = target.UpdatedAt
	return gomock.Eq(m.item).Matches(*target)
}

// String describes what the matcher matches.
func (m itemMatcher) String() string { return fmt.Sprintf("%v", &m.item) }

func TestAddTaskSavesAllFields(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	tt := newTask(ctl)
	form := &model.FormAddTask{
		Title:       "test",
		Due:         time.Now(),
		Description: "desc",
		Type:        model.TaskTypeCategory,
		ParentID:    entity.RootID + 1,
		Order:       42,
	}
	item := &entity.Item{
		Title:        form.Title,
		Due:          form.Due,
		Description:  form.Description,
		Order:        form.Order,
		State:        entity.ItemStateNormal,
		Type:         taskTypeToItemType(form.Type),
		ParentItemID: form.ParentID,
	}

	gomock.InOrder(
		tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemByID(form.ParentID).Return(&entity.Item{}, nil),
		tt.Storage.(*mock_use.MockStorage).EXPECT().SaveItem(itemMatcher{*item}).Return(i64, io.EOF),
	)

	err := tt.AddTask(form)
	assert.Error(t, err)
}

func TestAddTaskUpdatedsOrder(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	tt := newTask(ctl)
	var parentID int64 = entity.RootID
	gomock.InOrder(
		tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemByID(parentID).Return(&entity.Item{}, nil),
		tt.Storage.(*mock_use.MockStorage).EXPECT().SaveItem(gomock.Any()),
		tt.Storage.(*mock_use.MockStorage).EXPECT().IncreaseOrderAfter(gomock.Any()),
		tt.Presenter.(*mock_use.MockPresenter).EXPECT().ShowTaskAdded(gomock.Any()),
	)
	err := tt.AddTask(&model.FormAddTask{Title: "test", ParentID: entity.RootID})
	assert.NoError(t, err)
}

func TestAddTaskValidationErrors(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tt := newTask(ctl)

	for _, c := range []struct {
		form model.FormAddTask
		err  error
	}{
		{model.FormAddTask{Title: ""}, ErrEmptyTitle},
		{model.FormAddTask{Title: " "}, ErrEmptyTitle},
	} {
		err := tt.AddTask(&c.form)
		if c.err != nil {
			assert.True(t, errors.Is(err, c.err))
		} else {
			assert.NoError(t, err)
		}
	}

	// case with storage mock
	tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemByID(gomock.Any()).Return(nil, io.EOF)
	err := tt.AddTask(&model.FormAddTask{Title: "x", Due: time.Now(), ParentID: 42})
	assert.Error(t, err)
}

func TestAddTaskNonValidationErrors(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tt := newTask(ctl)

	// SaveItem error
	gomock.InOrder(
		tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemByID(gomock.Any()).Return(nil, nil),
		tt.Storage.(*mock_use.MockStorage).EXPECT().SaveItem(gomock.Any()).Return(i64, io.EOF),
	)

	err := tt.AddTask(&model.FormAddTask{Title: "x", Due: time.Now(), ParentID: 42})
	assert.Error(t, err)

	// IncreaseOrderAfter error
	gomock.InOrder(
		tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemByID(gomock.Any()).Return(nil, nil),
		tt.Storage.(*mock_use.MockStorage).EXPECT().SaveItem(gomock.Any()).Return(i64, nil),
		tt.Storage.(*mock_use.MockStorage).EXPECT().IncreaseOrderAfter(gomock.Any()).Return(io.EOF),
	)
	err = tt.AddTask(&model.FormAddTask{Title: "x", Due: time.Now(), ParentID: 42})
	assert.Error(t, err)
}

func TestChangeTaskStateByID(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tt := newTask(ctl)

	gomock.InOrder(
		tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemByID(gomock.Any()).Return(&entity.Item{}, nil),
		tt.Storage.(*mock_use.MockStorage).EXPECT().SaveItem(gomock.Any()).Return(i64, nil),
	)
	err := tt.ChangeTaskStateByID(1, model.TaskStateCompleted)
	assert.NoError(t, err)
}

func TestChangeTaskStateByIDErrors(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tt := newTask(ctl)

	// GetItemByID error
	tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemByID(gomock.Any()).Return(nil, io.EOF)
	err := tt.ChangeTaskStateByID(1, model.TaskStateCompleted)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, io.EOF))

	// SaveItem error
	gomock.InOrder(
		tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemByID(gomock.Any()).Return(&entity.Item{}, nil),
		tt.Storage.(*mock_use.MockStorage).EXPECT().SaveItem(gomock.Any()).Return(i64, io.EOF),
	)
	err = tt.ChangeTaskStateByID(1, model.TaskStateCompleted)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, io.EOF))
}

func TestListTasksByParentID(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tt := newTask(ctl)

	tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemsByParentID(gomock.Any()).Return([]*entity.Item{
		{Type: entity.ItemTypeCategory},
		{Type: entity.ItemTypeTask},
	}, nil)
	tt.Presenter.(*mock_use.MockPresenter).EXPECT().ShowTasksOfParentID(gomock.Any(), gomock.Any())
	err := tt.ListTasksByParentID(0)
	assert.NoError(t, err)
}

func TestListTasksByParentIDErrors(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tt := newTask(ctl)

	// GetItemsByParentID error
	tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemsByParentID(gomock.Any()).Return(nil, io.EOF)
	err := tt.ListTasksByParentID(0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, io.EOF))

	gomock.InOrder(
		tt.Storage.(*mock_use.MockStorage).EXPECT().GetItemsByParentID(gomock.Any()).Return(nil, nil),
		tt.Presenter.(*mock_use.MockPresenter).EXPECT().ShowTasksOfParentID(gomock.Any(), gomock.Any()).Return(io.EOF),
	)
	err = tt.ListTasksByParentID(0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, io.EOF))
}

func TestItemTypeToTaskTypeEqualLength(t *testing.T) {
	t.Parallel()
	assert.Equal(t, len(itemTypeToTaskTypeMap), len(taskTypeToItemTypeMap))
}
