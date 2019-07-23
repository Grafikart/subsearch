package opensubtitle

import (
	"github.com/kolo/xmlrpc"
	"path/filepath"
	"sync"
)

const (
	endpoint = "https://api.opensubtitles.org:443/xml-rpc"
)

type xmlRpcClient interface {
	Call(method string, args interface{}, res interface{}) error
}

type loginResponse struct {
	Token string `xmlrpc:"token"`
}

type dataResponse struct {
	Data Subtitles `xmlrpc:"data"`
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

type Client struct {
	Token  string
	Client xmlRpcClient
}

func (c *Client) Search(f File) (subtitles Subtitles, err error) {
	var subsFromFile, subsFromName Subtitles
	var wg sync.WaitGroup
	if err := c.login(); err != nil {
		return nil, err
	}
	wg.Add(2)
	go func() {
		subsFromFile, err = c.searchFromFile(f)
		defer wg.Done()
	}()
	go func() {
		subsFromName, err = c.searchFromName(filepath.Base(f.Name()))
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
	res := loginResponse{}
	args := []interface{}{"", "", "en", "OpenSubtitlesPlayer v4.7"}
	err = c.Client.Call("LogIn", args, &res)
	if err != nil {
		return
	}
	c.Token = res.Token
	return nil
}

func (c *Client) searchFromFile(f File) (subtitles Subtitles, err error) {
	h, err := hashFile(f)
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
			f.Size(),
			"eng",
		}},
	}
	res := dataResponse{}

	if err := c.Client.Call("SearchSubtitles", params, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) searchFromName(name string) (subtitles Subtitles, err error) {
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
	res := dataResponse{}
	if err := c.Client.Call("SearchSubtitles", params, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}
