package storage

import (
	"sort"
	"testing"
	"time"

	"github.com/tevino/the-clean-architecture-demo/todo/use"

	"github.com/stretchr/testify/assert"
	"github.com/tevino/the-clean-architecture-demo/todo/entity"
)

func foreachImplementations(t *testing.T, test func(use.Storage)) {
	for _, imp := range []use.Storage{
		NewMemory(),
		// NewFileSystem("temp/path")
	} {
		t.Logf("Testing storage implementation: %T", imp)
		test(imp)
	}
}

func TestSaveItemIDNotZero(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		id, err := s.SaveItem(&entity.Item{})
		assert.NoError(t, err)
		assert.NotZero(t, id)
	})
}

func TestSaveItemNilErr(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		_, err := s.SaveItem(nil)
		assert.Equal(t, ErrNilItem, err)
	})
}

func TestSaveItemUpdatesExistingItem(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		item := addTestingItems(t, s)[0]
		item.Description = "updated description"
		item.Title = "updated title"
		item.Order++
		item.State = entity.ItemStateCompleted
		item.Due = item.Due.Add(time.Second)

		oldUpdatedAt := item.UpdatedAt
		_, err := s.SaveItem(item)
		assert.NoError(t, err)
		newItem, err := s.GetItemByID(item.ID)
		assert.NoError(t, err)
		assert.True(t, newItem.UpdatedAt.After(oldUpdatedAt))
		assert.Equal(t, item, newItem)
	})
}

func TestGetItemsByParentIDReturnsNoRootID(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		addTestingItems(t, s)
		items, err := s.GetItemsByParentID(entity.RootID)
		assert.NoError(t, err)
		assert.NotEmpty(t, items)
		for _, it := range items {
			assert.NotEqual(t, it.ID, entity.RootID)
		}
	})
}

func TestIncreasesOrderAfter(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		// Add duplicate items to break the order
		item := addTestingItems(t, s)[0]
		addTestingItems(t, s)

		err := s.IncreaseOrderAfter(item)
		assert.NoError(t, err)

		// test if order of items after the inserted one are updated
		items, err := s.GetItemsByParentID(item.ParentItemID)
		assert.NoError(t, err)
		assert.True(t, len(items) > 1)

		for _, it := range items {
			if it.ID == item.ID {
				continue
			}
			assert.True(t, it.Order > item.Order)
		}
	})
}

func TestGetItemsByParentIDReturnsSortedASC(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		// Add duplicate items to break the order
		addTestingItems(t, s)
		addTestingItems(t, s)

		items, err := s.GetItemsByParentID(entity.RootID)
		assert.NoError(t, err)
		assert.True(t, sort.SliceIsSorted(items, func(i int, j int) bool {
			return items[i].Order < items[j].Order
		}))
	})
}

func TestGetItemByIDReturnsRootItem(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		it, err := s.GetItemByID(entity.RootID)
		assert.NoError(t, err)
		assert.Equal(t, entity.RootItem, it)
	})
}

func TestGetItemByID(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		items := addTestingItems(t, s)
		it, err := s.GetItemByID(items[1].ID)
		assert.NoError(t, err)
		assert.Equal(t, items[1], it)
	})
}

func TestGetItemByIDErrItemNotFound(t *testing.T) {
	foreachImplementations(t, func(s use.Storage) {
		_, err := s.GetItemByID(42)
		assert.Equal(t, ErrItemNotFound, err)
	})
}

func addTestingItems(t *testing.T, s use.Storage) []*entity.Item {
	t1 := &entity.Item{Title: "top1", ParentItemID: entity.RootID, Order: 1, Type: entity.ItemTypeCategory}
	t1ID, err := s.SaveItem(t1)
	assert.NoError(t, err)

	t2 := &entity.Item{Title: "top2", ParentItemID: entity.RootID, Order: 2, Type: entity.ItemTypeCategory}
	t2ID, err := s.SaveItem(t2)
	assert.NoError(t, err)

	s1 := &entity.Item{Title: "sub1", ParentItemID: t1ID, Order: 1, Type: entity.ItemTypeTask}
	_, err = s.SaveItem(s1)
	assert.NoError(t, err)

	s2 := &entity.Item{Title: "sub2", ParentItemID: t2ID, Order: 1, Type: entity.ItemTypeTask}
	_, err = s.SaveItem(s2)
	assert.NoError(t, err)

	return []*entity.Item{t1, t2, s1, s2}
}
