package opensubtitle

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
)

type Subtitle struct {
	IDMovie             string `xmlrpc:"IDMovie"`
	IDMovieImdb         string `xmlrpc:"IDMovieImdb"`
	IDSubMovieFile      string `xmlrpc:"IDSubMovieFile"`
	IDSubtitle          string `xmlrpc:"IDSubtitle"`
	IDSubtitleFile      string `xmlrpc:"IDSubtitleFile"`
	ISO639              string `xmlrpc:"ISO639"`
	LanguageName        string `xmlrpc:"LanguageName"`
	MatchedBy           string `xmlrpc:"MatchedBy"`
	MovieByteSize       string `xmlrpc:"MovieByteSize"`
	MovieFPS            string `xmlrpc:"MovieFPS"`
	MovieHash           string `xmlrpc:"MovieHash"`
	MovieImdbRating     string `xmlrpc:"MovieImdbRating"`
	MovieKind           string `xmlrpc:"MovieKind"`
	MovieName           string `xmlrpc:"MovieName"`
	MovieNameEng        string `xmlrpc:"MovieNameEng"`
	MovieReleaseName    string `xmlrpc:"MovieReleaseName"`
	MovieTimeMS         string `xmlrpc:"MovieTimeMS"`
	MovieYear           string `xmlrpc:"MovieYear"`
	MovieFileName       string `xmlrpc:"MovieName"`
	QueryNumber         string `xmlrpc:"QueryNumber"`
	SeriesEpisode       string `xmlrpc:"SeriesEpisode"`
	SeriesIMDBParent    string `xmlrpc:"SeriesIMDBParent"`
	SeriesSeason        string `xmlrpc:"SeriesSeason"`
	SubActualCD         string `xmlrpc:"SubActualCD"`
	SubAddDate          string `xmlrpc:"SubAddDate"`
	SubAuthorComment    string `xmlrpc:"SubAuthorComment"`
	SubAutoTranslation  string `xmlrpc:"SubAutoTranslation"`
	SubBad              string `xmlrpc:"SubBad"`
	SubComments         string `xmlrpc:"SubComments"`
	SubDownloadLink     string `xmlrpc:"SubDownloadLink"`
	SubDownloadsCnt     string `xmlrpc:"SubDownloadsCnt"`
	SubFeatured         string `xmlrpc:"SubFeatured"`
	SubFileName         string `xmlrpc:"SubFileName"`
	SubFormat           string `xmlrpc:"SubFormat"`
	SubHash             string `xmlrpc:"SubHash"`
	SubHD               string `xmlrpc:"SubHD"`
	SubHearingImpaired  string `xmlrpc:"SubHearingImpaired"`
	SubLanguageID       string `xmlrpc:"SubLanguageID"`
	SubRating           string `xmlrpc:"SubRating"`
	SubSize             string `xmlrpc:"SubSize"`
	SubSumCD            string `xmlrpc:"SubSumCD"`
	SubEncoding         string `xmlrpc:"SubEncoding"`
	SubForeignPartsOnly string `xmlrpc:"SubForeignPartsOnly"`
	SubFromTrusted      string `xmlrpc:"SubFromTrusted"`
	SubtitlesLink       string `xmlrpc:"SubtitlesLink"`
	UserID              string `xmlrpc:"UserID"`
	UserNickName        string `xmlrpc:"UserNickName"`
	UserRank            string `xmlrpc:"UserRank"`
	ZipDownloadLink     string `xmlrpc:"ZipDownloadLink"`
	subFilePath         string
}

type Subtitles []Subtitle

func (s Subtitle) Download(path string) (err error) {
	resp, err := http.Get(s.SubDownloadLink)
	if err != nil {
		return
	}

	var r io.Reader
	r, err = gzip.NewReader(resp.Body)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	f, err := os.Create(path)
	if err != nil {
		return
	}
	_, err = f.Write(resB.Bytes())
	if err != nil {
		return
	}
	return nil
}

func (s *Subtitles) ToMap() map[string]Subtitle {
	options := make(map[string]Subtitle)
	for _, sub := range *s {
		if sub.MatchedBy == "moviehash" {
			options["★"+sub.SubFileName] = sub
		} else {
			options[sub.SubFileName] = sub
		}
	}
	return options
}
