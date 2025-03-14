package serphandler

import (
	"fmt"
	"strings"
)

func getFileExt(name string) (string, error) {
	// Use built-in strings.LastIndex to find the last period
	idx := strings.LastIndex(name, ".")
	if idx == -1 {
		return "csv", fmt.Errorf("no file extension found, defaulting to csv")
	}

	// Extract the extension without the period and convert to lowercase
	ext := strings.ToLower(name[idx+1:])

	// Validate extension is one we support
	switch ext {
	case "csv", "json", "txt":
		return ext, nil
	default:
		return "csv", fmt.Errorf("unsupported file extension: %s, defaulting to csv", ext)
	}
}
