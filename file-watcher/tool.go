package filewatcher

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"

	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
)

// h : [md5.New() / sha1.New() / sha256.New()]
func FileHash(file string, h hash.Hash) string {

	if !fd.FileExists(file) {
		return ""
	}

	f, err := os.Open(file)
	lk.FailOnErr("%v", err)
	defer f.Close()
	_, err = io.Copy(h, f)
	lk.FailOnErr("%v", err)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SelfMD5() string {
	return FileHash(os.Args[0], md5.New())
}

func SelfSHA1() string {
	return FileHash(os.Args[0], sha1.New())
}

func SelfSHA256() string {
	return FileHash(os.Args[0], sha256.New())
}

func HasAnySuffix(s string, suffixGrp ...string) bool {
	for _, suffix := range suffixGrp {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func GetFileContentType(f *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
