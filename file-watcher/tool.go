package filewatcher

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
)

func SelfMD5() string {
	f, err := os.Open(os.Args[0])
	if err != nil {
		panic(err)
	}
	h := md5.New() // sha1.New() // sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HasAnySuffix(s string, suffixGrp ...string) bool {
	for _, suffix := range suffixGrp {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}
