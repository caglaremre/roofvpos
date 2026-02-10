package utils

import (
	"strings"
)

func TransformToken(token string) string {
	t1 := strings.ReplaceAll(token, ".", "+")
	t2 := strings.ReplaceAll(t1, "_", "/")
	t3 := strings.ReplaceAll(t2, "-", "=")
	return string(t3)
}
