package use

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tevino/the-clean-architecture-demo/todo/entity"
	"github.com/tevino/the-clean-architecture-demo/todo/model"
)

// TaskInteractor implementations the use case of task.
type TaskInteractor struct {
	Presenter
	Storage
}

// errors
var (
	ErrEmptyTitle = errors.New("Task title could not be empty")
)

func (t *TaskInteractor) validateAddTask(f *model.FormAddTask) error {
	if f.Title == "" || strings.TrimSpace(f.Title) == "" {
		return ErrEmptyTitle
	}
	if _, err := t.Storage.GetItemByID(f.ParentID); err != nil {
		return fmt.Errorf("getting parent item: %w", err)
	}
	return nil
}

func (t *TaskInteractor) AddTask(f *model.FormAddTask) error {
	if err := t.validateAddTask(f); err != nil {
		return fmt.Errorf("validating task: %w", err)
	}
	newTask := &entity.Item{
		Title:        f.Title,
		Due:          f.Due,
		Description:  f.Description,
		Order:        f.Order,
		Type:         taskTypeToItemType(f.Type),
		State:        entity.ItemStateNormal,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ParentItemID: f.ParentID,
	}
	taskID, err := t.Storage.SaveItem(newTask)
	if err != nil {
		return fmt.Errorf("saving task: %w", err)
	}
	newTask.ID = taskID

	err = t.Storage.IncreaseOrderAfter(newTask)
	if err != nil {
		return fmt.Errorf("changing order: %w", err)
	}

	return t.Presenter.ShowTaskAdded(itemToTask(newTask))
}

func (t *TaskInteractor) ChangeTaskStateByID(taskID int64, s model.TaskState) error {
	item, err := t.Storage.GetItemByID(taskID)
	if err != nil {
		return fmt.Errorf("getting item: %w", err)
	}
	item.State = taskStateToItemState(s)
	_, err = t.Storage.SaveItem(item)
	if err != nil {
		return fmt.Errorf("saving item: %w", err)
	}
	// TODO: Presenter.ShowTaskUpdated()?
	return nil
}

// ListTasksByParentID lists sub tasks of a given parent.
func (t *TaskInteractor) ListTasksByParentID(parentID int64) error {
	items, err := t.Storage.GetItemsByParentID(parentID)
	if err != nil {
		return fmt.Errorf("getting tasks from storage: %w", err)
	}
	tasks := make([]*model.Task, len(items))
	for i, it := range items {
		tasks[i] = itemToTask(it)
	}
	err = t.ShowTasksOfParentID(parentID, tasks)
	if err != nil {
		return fmt.Errorf("showing task of parent[%d]: %w", parentID, err)
	}
	return nil
}

// TODO: implement task search
// // SearchInTasks returns tasks with given keyword.
// func (t *TaskInteractor) SearchInTasks(keyword string) ([]*model.Task, error) {
// 	return nil, nil
// }

func itemToTask(it *entity.Item) *model.Task {
	return &model.Task{
		ID:          it.ID,
		Title:       it.Title,
		Due:         it.Due,
		Type:        itemTypeToTaskType(it.Type),
		State:       itemStateToTaskState(it.State),
		Description: it.Description,
		Order:       it.Order,
	}
}

//go:generate mockgen -destination mock_use/task_mock.go github.com/tevino/the-clean-architecture-demo/todo/use CasesTask,Presenter,Storage

// CasesTask represents the Input Port of the task Interactor.
type CasesTask interface {
	AddTask(*model.FormAddTask) error
	ListTasksByParentID(int64) error
	ChangeTaskStateByID(int64, model.TaskState) error
}

// Presenter represents the Output Port of Interactor.
type Presenter interface {
	ShowTaskAdded(*model.Task) error
	ShowTasksOfParentID(int64, []*model.Task) error
}

// Storage represents the entity gateway.
type Storage interface {
	SaveItem(*entity.Item) (int64, error)
	GetItemsByParentID(parentID int64) ([]*entity.Item, error)
	GetItemByID(int64) (*entity.Item, error)
	IncreaseOrderAfter(item *entity.Item) error
}
