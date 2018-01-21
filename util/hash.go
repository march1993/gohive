package util

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func Hash(content string) string {
	h := sha256.New()
	io.WriteString(h, content)
	return fmt.Sprintf("%x", h.Sum(nil))
}
