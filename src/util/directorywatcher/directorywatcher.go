// import "util/directorywatcher"
package directorywatcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type directoryWatcher struct {
	Interval  uint64 // interval in ms
	Recursive bool   // Use filepath.Walk or filepath.Glob?
	Pattern   string // glob pattern

	// Internal details
	path      string                 // the path being watched
	files     map[string]os.FileInfo // Map of files watched
	ticker    *time.Ticker
	observers []observer // List of observers
}

// Type of observer function - adding an observer means adding a function of this type
type observer chan []Event

/**
 * API
 *
 * Ruby's directory watcher supports the following:
 *  - preload: suppress the initial added events by preloading files
 *  - '** /*.go' means all .go files found in the subtree rooted at the path
 */
func New(path string) (*directoryWatcher, error) {
	if stat, err := os.Stat(path); err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("Provided path is not a directory: %s", path)
	}

	return &directoryWatcher{
		Interval:  2000,
		observers: []observer{},
		path:      path,
		files:     make(map[string]os.FileInfo),
	}, nil
}

// Set up some sort of looping mechanism: Maybe use a goroutine with two
// channels: 'walk' and 'stop' - receiving a message on 'walk' kicks off the
// walking routine, passing back a list of changes on the same channel.
// Receiving a message on 'stop' will stop the iteration
func (dw *directoryWatcher) Start() {
	if dw.ticker != nil {
		return
	}

	dw.ticker = time.NewTicker(time.Duration(dw.Interval) * time.Millisecond)

	go func() {
		for now := range dw.ticker.C {
			dw.scan(now)
		}
	}()
}

func (dw *directoryWatcher) Stop() {
	dw.ticker.Stop()
	dw.ticker = nil
}

func (dw *directoryWatcher) AddObserver(obs observer) {
	dw.observers = append(dw.observers, obs)
}

// The actual walking function
func (dw *directoryWatcher) scan(at time.Time) {
	var changed []Event                 // The new events
	var touched = make(map[string]bool) // path names of the files seen in a pass

	if dw.Recursive {
		filepath.Walk(dw.path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || !matches(dw.Pattern, info.Name()) {
				return nil
			}
			if ev, ok := dw.hasChange(path, info); ok {
				dw.files[path] = info
				changed = append(changed, ev)
			}
			touched[path] = true
			return err
		})
	} else {
		matches, _ := filepath.Glob(filepath.Join(dw.path, dw.Pattern))
		for _, path := range matches {
			info, err := os.Stat(path)
			if err != nil || info.IsDir() {
				continue
			}
			if ev, ok := dw.hasChange(path, info); ok {
				dw.files[path] = info
				changed = append(changed, ev)
			}
			touched[path] = true
		}
	}
	for path, info := range dw.files {
		if !touched[path] {
			changed = append(changed, Event{Deleted, path, info})
			delete(dw.files, path)
		}
	}

	// Notify observers if anything changed
	if len(changed) > 0 {
		for _, c := range dw.observers {
			c <- changed
		}
	}

}

func matches(pattern, name string) bool {
	matched, err := filepath.Match(pattern, name)
	return err == nil && matched
}

/**
 * This tells us if a given file has been changed or added.
 *
 * Uses the comma-ok style to indicate whether or not a given file actually changed.
 */
func (dw *directoryWatcher) hasChange(path string, info os.FileInfo) (Event, bool) {
	if oldInfo, ok := dw.files[path]; ok {
		if info.ModTime().After(oldInfo.ModTime()) {
			return Event{Changed, path, info}, true
		}
		return Event{}, false
	}
	return Event{Added, path, info}, true
}
