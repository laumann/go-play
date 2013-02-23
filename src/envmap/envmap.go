package envmap

import (
	"os"
	"strings"
)

// Get the environment as a map[string]string
func Env() map[string]string {
	envMap := make(map[string]string)

	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		envMap[kv[0]] = kv[1]
	}
	return envMap
}

// Only compute this map once
var __env map[string]string

// Should rename the package to 'env' and this to 'Map()', so one would call it
// with 'env.Map()'
func Map() map[string]string {
	if __env != nil {
		return __env
	}

	__env = make(map[string]string)
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		__env[kv[0]] = kv[1]
	}
	return __env
}
