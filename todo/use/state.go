package use

import (
	"github.com/tevino/the-clean-architecture-demo/todo/entity"
	"github.com/tevino/the-clean-architecture-demo/todo/model"
)

var taskStateToItemStateMap = map[model.TaskState]entity.ItemState{
	model.TaskStateNormal:    entity.ItemStateNormal,
	model.TaskStateCompleted: entity.ItemStateCompleted,
}

var itemStateToTaskStateMap = map[entity.ItemState]model.TaskState{}

func init() {
	for k, v := range taskStateToItemStateMap {
		itemStateToTaskStateMap[v] = k
	}
}

func taskStateToItemState(s model.TaskState) entity.ItemState {
	return taskStateToItemStateMap[s]
}

func itemStateToTaskState(s entity.ItemState) model.TaskState {
	return itemStateToTaskStateMap[s]
}
