package use

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tevino/the-clean-architecture-demo/todo/entity"
)

const template = `
+ Inbox

Inbox is the place to dump your thoughts into.

    [ ] Welcome!
    today
    This is the description

    [ ] Press ? to show help
    today

    [ ] Use j/k to move down/up
    today
+ Projects
`

func (t *TaskInteractor) AddTemplate() error {
	levelInfoMap := make(map[int]*levelInfo)
	previousLineIsItem := false
	var item *entity.Item
	for i, line := range strings.Split(template, "\n") {
		if line == "" {
			continue
		}
		level := getLeadingSpace(line) / 4
		if _, ok := levelInfoMap[level]; !ok {
			levelInfoMap[level] = &levelInfo{ParentID: entity.RootID}
		}
		if _, ok := levelInfoMap[level+1]; !ok {
			levelInfoMap[level+1] = &levelInfo{}
		}
		line = strings.TrimSpace(line)
		switch line[0] {
		case '+':
			previousLineIsItem = false
			// category
			levelInfoMap[level].Order++

			item = &entity.Item{Title: line[2:], Type: entity.ItemTypeCategory,
				Order: levelInfoMap[level].Order, ParentItemID: levelInfoMap[level].ParentID}
		case '[':
			// task
			previousLineIsItem = true
			levelInfoMap[level].Order++
			item = &entity.Item{Title: line[4:], Type: entity.ItemTypeTask,
				Order: levelInfoMap[level].Order, ParentItemID: levelInfoMap[level].ParentID}
		default:
			// item detail
			if previousLineIsItem {
				// this line is due
				due, err := _parseDue(line)
				if err != nil {
					return fmt.Errorf("parse due of template[L%d]: %w", i, err)
				}
				item.Due = due
			} else {
				// this line is description
				item.Description += (strings.TrimSpace(line) + "\n")
			}
			previousLineIsItem = false
		}

		if item != nil {
			id, err := t.Storage.SaveItem(item)
			if err != nil {
				return fmt.Errorf("save item: %w", err)
			}
			levelInfoMap[level+1].ParentID = id
		}
	}
	return nil
}

func getLeadingSpace(l string) int {
	for i, c := range l {
		if c != ' ' {
			return i
		}
	}
	return len(l)
}

var errInvalidDue = errors.New("invalid due")

func _parseDue(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	var d time.Time
	switch s {
	case "", "today":
		d = time.Now()
	case "tom", "tomorrow":
		d = time.Now().Add(time.Hour * 24)
	default:
		return d, errInvalidDue
	}
	// TODO: use regexp to implement more, like +2d, next week, month
	return d, nil
}

type levelInfo struct {
	ParentID int64
	Order    uint64
}
