package main

import (
	"log"

	"github.com/tevino/the-clean-architecture-demo/todo/cui/io"

	"github.com/tevino/the-clean-architecture-demo/todo/cui"
	"github.com/tevino/the-clean-architecture-demo/todo/storage"
	"github.com/tevino/the-clean-architecture-demo/todo/use"
)

func main() {
	ui := cui.New(&io.TermUI{})
	mem := storage.NewMemory()
	presenter := &cui.Presenter{ui}
	cases := &use.TaskInteractor{
		Presenter: presenter,
		Storage:   mem,
	}
	ctl := &cui.Controller{
		CUI:       ui,
		IO:        &io.UnixLikeIO{},
		CasesTask: cases,
	}
	cases.AddTemplate()
	if err := ctl.Loop(); err != nil {
		log.Fatal(err)
	}
}
