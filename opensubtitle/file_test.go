package opensubtitle

import (
	"strings"
	"testing"
)

func TestHashFile(t *testing.T) {

	t.Run("block if content is too small", func(t *testing.T) {
		_, err := hashFile(strings.NewReader("a"))
		if err == nil {
			t.Errorf("error expected")
		}
	})

	t.Run("hash file correctly", func(t *testing.T) {
		hashes := map[string]string{
			"hello this is a test":      "6a80582a627a80c7",
			"Another short test to see": "ece6b95fa5caf17a",
			"a":                         "58585858585a4000",
		}

		for k, want := range hashes {
			got, err := hashFile(genStringReader(k))
			if err != nil {
				t.Errorf("error unexpected %v", err)
			}
			if got != want {
				t.Errorf("expected %q, got %q", want, got)
			}
		}
	})
}

func genStringReader(s string) *strings.Reader {
	for len(s) <= ChunkSize {
		s += s
	}
	return strings.NewReader(s)
}
