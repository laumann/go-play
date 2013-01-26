package envmap

/**
 * Get the environment as a map[string]string
 */

import (
	"os"
	"strings"
)

func Env() map[string]string {
	envMap := make(map[string]string)

	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		envMap[kv[0]] = kv[1]
	}
	return envMap
}
