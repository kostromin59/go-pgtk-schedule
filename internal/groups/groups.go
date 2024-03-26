package groups

import (
	"sync"
)

type Group struct {
	Name      string
	Value     string
	Subgroups []string
	mu        sync.Mutex
}

func NewGroup(name, value string) *Group {
	return &Group{Name: name, Value: value}
}

// TODO: Сделать метод
func (g *Group) ParseSubgroups(studyYearId, semester string, weekNumber int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	return nil
}
