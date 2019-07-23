package opensubtitle

import (
	"reflect"
	"strings"
	"testing"
)

type fakeFile struct {
	*strings.Reader
	name string
}

func (f fakeFile) Name() string {
	return f.name
}

type fakeXmlRpcClient struct {
	t     *testing.T
	calls map[string]interface{}
}

func (c fakeXmlRpcClient) Call(method string, args interface{}, res interface{}) error {
	if c.calls[method] == nil {
		c.t.Errorf("unexpected call to method %q", method)
	}
	var data interface{}
	data = c.calls[method]
	reflect.ValueOf(res).Elem().Set(reflect.ValueOf(data))
	return nil
}

func newFakeClient(t *testing.T, subtitles Subtitles) Client {
	fakeDataResponse := dataResponse{Data: subtitles}
	fakeLoginResponse := loginResponse{"Fake"}
	c := fakeXmlRpcClient{t, map[string]interface{}{
		"LogIn":           fakeLoginResponse,
		"SearchSubtitles": fakeDataResponse,
	}}
	return Client{
		Client: c,
	}
}

func TestClient(t *testing.T) {

	t.Run("Search work as expected", func(t *testing.T) {
		s1 := Subtitle{IDMovie: "a1"}
		s2 := Subtitle{IDMovie: "a2"}
		client := newFakeClient(t, Subtitles{s1, s2})
		file := fakeFile{genStringReader("fake conntent"), "filename.mkv"}
		subtitles, err := client.Search(file)
		if err != nil {
			t.Errorf("expected no errors, got %v", err)
		}
		if len(subtitles) != 4 {
			t.Errorf("expected %d subtitles", 4)
		}
	})

}
