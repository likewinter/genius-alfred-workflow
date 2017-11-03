package main

import (
	"net/http"
	"net/url"
	"fmt"
	"encoding/json"
	"regexp"
	"os"
	"os/exec"
	"log"
)

type Song struct {
	Title string
}

type GeniusResults struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Response struct {
		Sections []struct {
			Type string `json:"type"`
			Hits []struct {
				Result struct {
					URL string `json:"url"`
				} `json:"result"`
			} `json:"hits"`
		} `json:"sections"`
	} `json:"response"`
}

func (s Song) GetCleanTitle() string {
	versionRegexp := `(?i)([(\[]).*?(version|remastered|single|remix|(\d{2,4})).*?([)\]])`
	return regexp.MustCompile(versionRegexp).ReplaceAllString(s.Title, "")
}

func (s Song) SearchLyrics() (results *GeniusResults, err error) {
	apiURL := fmt.Sprintf("https://genius.com/api/search/multi?q=%s", url.QueryEscape(s.GetCleanTitle()))
	res, err := http.Get(apiURL)
	if err == nil {
		results = new(GeniusResults)
		json.NewDecoder(res.Body).Decode(results)
	}
	return
}

func (r GeniusResults) GetURL() (lyricsURL string, hasURL bool) {
	if len(r.Response.Sections) > 1 && len(r.Response.Sections[1].Hits) > 0 {
		lyricsURL = r.Response.Sections[1].Hits[0].Result.URL
		hasURL = true
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Pass the Title via first argument")
	}
	song := Song{Title: os.Args[1]}
	results, err := song.SearchLyrics()
	if lyricsURL, hasURL := results.GetURL(); err == nil && hasURL {
		cmd := exec.Command("open", lyricsURL)
		cmd.Start()
	}
}
