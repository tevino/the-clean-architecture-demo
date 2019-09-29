package cui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/tevino/the-clean-architecture-demo/todo/cui/component"
	"github.com/tevino/the-clean-architecture-demo/todo/cui/io"
)

// CUI represents the Console User Interface, the instance of this struct is shared by Presenter and Controller.
type CUI struct {
	io.CUILib
	grid       *component.GridComponent
	taskList   component.TaskList
	catList    component.TaskList
	stateBar   component.Text
	descBox    component.Text
	components []component.Component
}

// New creates a new CUI.
func New(lib io.CUILib) *CUI {
	catList := component.NewListComponent("Categories")
	taskList := component.NewListComponent("Tasks")
	stateBar := component.NewTextComponent("State")
	descBox := component.NewTextComponent("Description")
	c := &CUI{
		CUILib: lib,
		grid: component.NewGridComponent(
			map[string]component.InteractiveComponent{
				component.DefaultActivated: catList,
				"h":                        catList,
				"<Left>":                   catList,
				"l":                        taskList,
				"<Right>":                  taskList,
			},
			ui.NewRow(9.0/10,
				ui.NewCol(2.0/10, catList),
				ui.NewCol(8.0/10,
					ui.NewRow(5.0/10, taskList),
					ui.NewRow(5.0/10, descBox),
				),
			),
			ui.NewRow(1.0/10,
				ui.NewCol(1.0/1, stateBar),
			),
		),
		taskList: taskList,
		catList:  catList,
		stateBar: stateBar,
		descBox:  descBox,
		components: []component.Component{
			taskList,
			catList,
		},
	}
	return c
}
