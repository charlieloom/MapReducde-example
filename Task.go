package main

type Task struct {
	Param    string `json:"param,omitempty"`
	Taskname string `json:"taskname,omitempty"`
	Group    string `json:"group,omitempty"`
}
