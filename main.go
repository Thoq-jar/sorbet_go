package sorbet

import (
	"strings"
)

type SorbetError struct {
	ErrorType string
	Message   string
}

func Parse(contents string) map[string]string {
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
				printError("Syntax", "Syntax error! Expected [key] => [value] at: "+line)
			}
			currentKey = strings.TrimSpace(parts[0])
			currentValue = strings.TrimSpace(parts[1])
		} else if len(trimmedLine) > 0 && trimmedLine[0] == '>' {
			if currentKey != "" {
				currentValue += "," + strings.TrimSpace(strings.TrimPrefix(trimmedLine, ">"))
			} else {
				printError("SyntaxException", "Continuation line without a key at: "+line)
			}
		}
	}

	if currentKey != "" {
		result[currentKey] = strings.TrimSpace(currentValue)
	}

	return result
}

func printError(errorType string, message string) {
	panic(errorType + ": " + message)
}
