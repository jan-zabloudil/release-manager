package util

import "strings"

// ExplodeFileKey splits the file key into bucket and file path.
// The file key is in the format of `bucket/file-path`.
func ExplodeFileKey(key string) (bucket, filePath string) {
	parts := strings.Split(key, "/")
	return parts[0], strings.Join(parts[1:], "/")
}
