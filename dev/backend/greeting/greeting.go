package greeting

import (
	"fmt"
	"strings"
)

func Hello(name string) string { // Exported Function
	return fmt.Sprintf("Olá, %s!", capitalize(name))
}

func capitalize(s string) string { // Unexported Function
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
