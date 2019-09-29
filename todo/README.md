# Todo

## Overview

This is the first attempt to put The Clean Architecture(TCA) into practice.

Some concepts of TCA are used as file and package names to help readers and writers like you to get a better understanding, however, in practice the names may vary.

I try to make this application usable in practical rather than a pure demo, since some problems may not be revealed in the case of a relatively simple demo.

## Features

- Category/Project
- Task with due and description
- Edit with your favorite editor
- An interactive console user interface
- Vi-like key map

## TODO

- [ ] Implement `?` for help
- [ ] Make template a basic tutorial
- [ ] Implement `dd` to let user move tasks
- [ ] Implement edit of existing tasks
- [ ] Humanize due dates
- [ ] Add `.` and `..` to the task list so that the navigation of nested tasks is possible
- [ ] Implement the file system based storage
- [ ] Complete the concept of View(a set of conditions to filter tasks, there may be views like `Today`, `This Week` etc)
- [ ] Implement a storage that interacts with existing TODO applications like OmniFocus or Todoist
- [ ] Add a web interface
- [ ] Add HTTP API support
- [ ] Tests to add more coverage
- [ ] Add TravisCI to make sure PR passes the test
