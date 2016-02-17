package main

type Line struct {
	name   string
	status string
	text   string
}

func (l *Line) ToString() string {
	return l.name + l.status + l.text
}
