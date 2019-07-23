package opensubtitle

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubtitlesToMap(t *testing.T) {

	t.Run("subitles can be converted to map", func(t *testing.T) {
		s1 := Subtitle{SubFileName: "Demo"}
		s2 := Subtitle{SubFileName: "Demo2"}
		s := Subtitles{
			s1,
			s2,
		}
		m := s.ToMap()
		if *m["Demo"] != s2 {
			t.Errorf("expected %v, got %v", s1, *m["Demo"])
		}
	})

	t.Run("download is working correctly", func(t *testing.T) {
		fakeConntent := "Hello World"
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			gz := gzip.NewWriter(res)
			defer gz.Close()
			_, err := gz.Write([]byte(fakeConntent))
			if err != nil {
				t.Errorf("unexpected error while writing test data")
			}
		}))
		defer testServer.Close()
		w := new(bytes.Buffer)
		s := Subtitle{SubDownloadLink: testServer.URL}
		err := s.Download(w)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if w.String() != fakeConntent {
			t.Errorf("expected %q, got %q", fakeConntent, w.String())
		}
	})
}
