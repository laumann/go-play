package pathutil

import (
	"bytes"
	"os"
	"path/filepath"
)

type Path struct {
	orig     string   // Original 
	segments []string // The different parts of the path
	root     string   // The root
	abs      string
}

func (p *Path) String() string {
	return p.orig
}

func (p *Path) NameCount() int {
	return len(p.segments)
}

/* Doesn't handle input like '~/Downloads' properly */
func (p *Path) ToAbs() string {
	if p.abs != "" {
		return p.abs
	}

	toAbs := p.orig
	if toAbs[0] == '~' {
		if toAbs[1] == '/' {
			home := os.ExpandEnv("$HOME")
			toAbs = home + toAbs[1:]
		} else {
			/* String starts with ~user (something) */
		}
	}
	p.abs, _ = filepath.Abs(filepath.Clean(toAbs))

	return p.abs
}

func Get(p string) *Path {
	var segment bytes.Buffer
	var segments []string

	for _, r := range p {
		if r == '/' {
			if segment.Len() > 0 {
				segments = append(segments, segment.String())
				segment.Reset()
			}
			continue
		}
		segment.WriteRune(r)
	}

	if segment.Len() > 0 {
		segments = append(segments, segment.String())
	}

	return &Path{p, segments, "/", ""}
}
