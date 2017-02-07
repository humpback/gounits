package rand

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const shortLen = 12

var validShortID = regexp.MustCompile("^[a-z0-9]{12}$")

func IsShortID(id string) bool {

	return validShortID.MatchString(id)
}

func TruncateID(id string) string {

	if i := strings.IndexRune(id, ':'); i >= 0 {
		id = id[i+1:]
	}
	if len(id) > shortLen {
		id = id[:shortLen]
	}
	return id
}

// GenerateRandomID returns a 32-bit unique id
func GenerateRandomID() string {

	b := make([]byte, 32)
	for {
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			panic(err)
		}
		id := hex.EncodeToString(b)
		if _, err := strconv.ParseInt(TruncateID(id), 10, 64); err == nil {
			continue
		}
		return id
	}
}
