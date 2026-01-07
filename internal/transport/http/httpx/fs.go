package httpx

import (
	"os"
	"path/filepath"
	"strings"
)

func SafeJoin(base string, parts ...string) (string, bool) {
	path := filepath.Clean(filepath.Join(append([]string{base}, parts...)...))
	if !strings.HasPrefix(path, base+string(os.PathSeparator)) {
		return "", false
	}
	return path, true
}
