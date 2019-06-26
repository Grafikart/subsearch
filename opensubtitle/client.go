package opensubtitle

import (
	"github.com/kolo/xmlrpc"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	endpoint = "https://api.opensubtitles.org:443/xml-rpc"
)

type Client struct {
	Token string
	*xmlrpc.Client
}

func (c *Client) Search(path string) (subtitles Subtitles, err error) {
	var subsFromFile, subsFromName Subtitles
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		subsFromFile, err = c.searchFromFile(path)
		defer wg.Done()
	}()
	go func() {
		subsFromName, err = c.searchFromName(filepath.Base(path))
		defer wg.Done()
	}()

	wg.Wait()

	if err != nil {
		return
	}

	return append(subsFromFile, subsFromName...), nil
}

func (c *Client) login() (err error) {
	if c.Token != "" {
		return nil
	}
	res := struct {
		Token string `xmlrpc:"token"`
	}{}
	args := []interface{}{"", "", "en", "OpenSubtitlesPlayer v4.7"}
	err = c.Client.Call("LogIn", args, &res)
	if err != nil {
		return
	}
	c.Token = res.Token
	return nil
}

func (c *Client) searchFromFile(path string) (subtitles Subtitles, err error) {
	if err := c.login(); err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	fi, err := f.Stat()
	if err != nil {
		return
	}
	h, err := HashFile(f)
	if err != nil {
		return
	}
	params := []interface{}{
		c.Token,
		[]struct {
			Hash  string `xmlrpc:"moviehash"`
			Size  int64  `xmlrpc:"moviebytesize"`
			Langs string `xmlrpc:"sublanguageid"`
		}{{
			h,
			fi.Size(),
			"eng",
		}},
	}
	res := struct {
		Data Subtitles `xmlrpc:"data"`
	}{}

	if err := c.Call("SearchSubtitles", params, &res); err != nil {
		if strings.Contains(err.Error(), "type mismatch") {
			return nil, err
		}
	}

	return subtitles, nil
}

func (c *Client) searchFromName(name string) (subtitles Subtitles, err error) {
	if err := c.login(); err != nil {
		return nil, err
	}
	params := []interface{}{
		c.Token,
		[]struct {
			Query string `xmlrpc:"query"`
			Langs string `xmlrpc:"sublanguageid"`
		}{{
			name,
			"eng",
		}},
	}
	res := struct {
		Data Subtitles `xmlrpc:"data"`
	}{}
	if err := c.Call("SearchSubtitles", params, &res); err != nil {
		if strings.Contains(err.Error(), "type mismatch") {
			return nil, err
		}
	}
	return res.Data, nil
}

func NewClient() (*Client, error) {
	rpc, err := xmlrpc.NewClient(endpoint, nil)
	if err != nil {
		return nil, err
	}
	c := &Client{
		Client: rpc,
	}
	return c, nil
}
