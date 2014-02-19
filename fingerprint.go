package monk

import (
	"crypto/md5"
	"fmt"
	"time"
)

type fingerprint struct {
	modtime time.Time
	hash    string
}

var fingerprintCache = make(map[string]*fingerprint)

// This function is not safe to call concurrently.
func GenerateFingerprint(fs fileSystem, path string) (string, error) {
	if fp, ok := fingerprintCache[path]; ok {
		return fp.hash, nil
	}

	content, err := fs.ReadFile(path)

	if err != nil {
		return "", err
	}

	hash := md5.New()
	_, err = hash.Write(content)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
