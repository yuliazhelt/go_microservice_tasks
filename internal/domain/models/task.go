package models

type Task struct {
	Id	string
	Author	string
	Title	string
	Description	string
	Approves	[]*Approve
	IsCancelled	bool
}
