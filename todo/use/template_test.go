package use

import (
	"errors"
	"io"
	"testing"

	"github.com/tevino/the-clean-architecture-demo/todo/use/mock_use"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestGetLeadingSpace(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, getLeadingSpace(""))
	assert.Equal(t, 1, getLeadingSpace(" "))
	assert.Equal(t, 4, getLeadingSpace("    "))
	assert.Equal(t, 4, getLeadingSpace("    + test"))
	assert.Equal(t, 8, getLeadingSpace("        [ ]"))
}

func TestTaskInteractor_AddTemplate(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tt := newTask(ctl)

	tt.Storage.(*mock_use.MockStorage).EXPECT().SaveItem(gomock.Any()).AnyTimes()
	err := tt.AddTemplate()

	assert.NoError(t, err)
}

func TestTaskInteractor_AddTemplate_Errors(t *testing.T) {
	t.Parallel()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tt := newTask(ctl)

	// SaveItem error
	tt.Storage.(*mock_use.MockStorage).EXPECT().SaveItem(gomock.Any()).Return(i64, io.EOF)
	err := tt.AddTemplate()

	assert.Error(t, err)
	assert.True(t, errors.Is(err, io.EOF))
}
