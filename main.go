package sorbet

import (
	"fmt"
	"strings"
)

type SorbetError struct {
	ErrorType string
	Message   string
}

func (e *SorbetError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorType, e.Message)
}

func Parse(contents string) (map[string]string, error) {
	result := make(map[string]string)
	var currentKey string
	currentValue := ""

	lines := strings.Split(contents, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.Contains(line, "=>") {
			if currentKey != "" {
				result[currentKey] = strings.TrimSpace(currentValue)
				currentValue = ""
			}

			parts := strings.Split(line, "=>")
			if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
				return nil, &SorbetError{
					ErrorType: "Syntax",
					Message:   "Expected [key] => [value] at: " + line,
				}
			}
			currentKey = strings.TrimSpace(parts[0])
			currentValue = strings.TrimSpace(parts[1])
		} else if len(trimmedLine) > 0 && trimmedLine[0] == '>' {
			if currentKey != "" {
				currentValue += "," + strings.TrimSpace(strings.TrimPrefix(trimmedLine, ">"))
			} else {
				return nil, &SorbetError{
					ErrorType: "SyntaxException",
					Message:   "Continuation line without a key at: " + line,
				}
			}
		}
	}

	if currentKey != "" {
		result[currentKey] = strings.TrimSpace(currentValue)
	}

	return result, nil
}
