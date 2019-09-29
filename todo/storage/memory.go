package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/tevino/the-clean-architecture-demo/todo/entity"
)

// Errors
var (
	ErrNilItem      = errors.New("Item could not be nil")
	ErrItemNotFound = errors.New("Item not found")
)

// Memory is a memory based volatile storage.
type Memory struct {
	id    int64
	items []*entity.Item
}

// NewMemory creates a Memory.
func NewMemory() *Memory {
	return &Memory{1, make([]*entity.Item, 0)}
}

// SaveItem saves an item into memory, return its id.
func (m *Memory) SaveItem(item *entity.Item) (int64, error) {
	if item == nil {
		return -1, ErrNilItem
	}
	// update timestamps
	now := time.Now().UTC()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	if item.ID > 0 {
		for i, it := range m.items {
			if it.ID == item.ID {
				m.items[i] = copyItem(item)
				return item.ID, nil
			}
		}
	}

	// item does not exist
	item.ID += m.id
	m.items = append(m.items, copyItem(item))
	m.id++

	return item.ID, nil
}

// IncreaseOrderAfter increases order by one for items after given one.
func (m *Memory) IncreaseOrderAfter(item *entity.Item) error {
	items, _ := m.GetItemsByParentID(item.ParentItemID)
	for _, it := range items {
		if it.ID == item.ID {
			continue
		}
		if it.Order >= item.Order {
			it.Order++
		}
	}
	return nil
}

// GetItemsByParentID returns items of given parent.
func (m *Memory) GetItemsByParentID(parentID int64) ([]*entity.Item, error) {
	items := []*entity.Item{}
	for _, it := range m.items {
		if it.ParentItemID == parentID {
			items = append(items, it)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Order < items[j].Order
	})
	return items, nil
}

// GetItemByID returns items of given ID.
func (m *Memory) GetItemByID(id int64) (*entity.Item, error) {
	if id == entity.RootID {
		return entity.RootItem, nil
	}

	for _, it := range m.items {
		if it.ID == id {
			return it, nil
		}
	}

	return nil, ErrItemNotFound
}

func copyItem(it *entity.Item) *entity.Item {
	buf, err := json.Marshal(it)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Item: %s", err))
	}
	var clone entity.Item
	err = json.Unmarshal(buf, &clone)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal Item: %s", err))
	}
	return &clone
}
