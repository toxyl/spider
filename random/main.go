package random

import (
	"math/rand"
	"time"

	gostringgenerator "github.com/toxyl/go-string-generator"
	"github.com/toxyl/go-string-generator/tokens"
)

var (
	stringGenerator *tokens.RandomStringGenerator
	DataDir         = ""
	Taunts          = []string{}
)

func Taunt() string {
	if stringGenerator == nil {
		stringGenerator = gostringgenerator.NewGenerator(DataDir, func(err error) {})
	}
	return stringGenerator.Generate(StringFromList(Taunts...))
}

func StringFromList(strings ...string) string {
	if len(strings) <= 0 {
		return ""
	}
	var i int = Int(0, len(strings)-1)
	return strings[i]
}

func Linebreak() string {
	return StringFromList("\r\n", "\r", "\n")
}

func Int(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	n := max - min + 1
	if n <= 0 {
		return min
	}
	return rand.Intn(n) + min
}

// GenerateGarbage produces a string (length is randomly chosen between 1 and `n`)
// consisting of random (non)-printable characters.
func GenerateGarbage(n int) string {
	garbage := make([]byte, Int(1, n))
	rand.Read(garbage)
	return string(garbage)
}
