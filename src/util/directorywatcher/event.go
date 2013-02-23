package directorywatcher

import (
	"os"
	"fmt"
)

type eventType int

const (
	Added eventType = iota
	Changed
	Deleted
)

// Mapping event types to a string, for implementing Stringer interface
var eventNames = map[eventType]string{
	Added:   "Added",
	Changed: "Changed",
	Deleted: "Deleted",
}

// Implement Stringer
func (e Event) String() string {
	return fmt.Sprintf("%s %s", eventNames[e.Type], e.Name())
}

// An event contains its type and the file involved.
type Event struct {
	Type eventType
	os.FileInfo
}
