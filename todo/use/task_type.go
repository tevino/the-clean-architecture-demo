package use

import (
	"github.com/tevino/the-clean-architecture-demo/todo/entity"
	"github.com/tevino/the-clean-architecture-demo/todo/model"
)

var itemTypeToTaskTypeMap = map[entity.ItemType]model.TaskType{
	entity.ItemTypeCategory: model.TaskTypeCategory,
	entity.ItemTypeTask:     model.TaskTypeTask,
}

var taskTypeToItemTypeMap = make(map[model.TaskType]entity.ItemType, len(itemTypeToTaskTypeMap))

func init() {
	for it, tt := range itemTypeToTaskTypeMap {
		taskTypeToItemTypeMap[tt] = it
	}
}

func itemTypeToTaskType(it entity.ItemType) model.TaskType {
	return itemTypeToTaskTypeMap[it]
}

func taskTypeToItemType(tt model.TaskType) entity.ItemType {
	return taskTypeToItemTypeMap[tt]
}
