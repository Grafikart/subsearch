package opensubtitle

import "testing"

func TestSubtitlesToMap(t *testing.T) {
	s1 := Subtitle{SubFileName: "Demo"}
	s2 := Subtitle{SubFileName: "Demo2"}
	s := Subtitles{
		s1,
		s2,
	}
	m := s.ToMap()
	if *m["Demo"] != s2 {
		t.Errorf("%v != %v", s1, m["Demo"])
	}
}
