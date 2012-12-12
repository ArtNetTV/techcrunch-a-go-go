package main

import (
	   "fmt"
	   "net/http"
	   "strings"
	   "io/ioutil"
	   "github.com/PuerkitoBio/goquery"
	   "encoding/json"
	   "time"
	   "os"
	   "errors"
)

type Post struct {
	 Id string
	 Title string
	 Author string
	 Url string
	 Shares uint16
	 Likes uint16
	 Comments uint16
}

type CaptureEvent struct {
	 CapturedAt time.Time
	 Posts []Post
}

func parsePosts() error {
	var document *goquery.Document
	var err error

	if document, err = goquery.NewDocument("http://www.techcrunch.com"); err != nil {
	   return errors.New("Failed to fetch techcrunch article")
	}

	matches := document.Find("div .post")
	if matches == nil {
	   return nil
	}

	var captureEvent CaptureEvent
	captureEvent.Posts = make([]Post, matches.Length())
	captureEvent.CapturedAt = time.Now()

	matches.Each(func(i int, s *goquery.Selection) {

		var post *Post = &captureEvent.Posts[i];

		post.Id, _ = s.Attr("id")
		if post.Author = s.Find("a span.name").Text(); post.Author == "" {
			post.Author = s.Find("div.by-line div.by-line").Text()
		}
		post.Title = strings.TrimSpace(s.Find("h2.headline a").Text());
		post.Url, _ = s.Find("h2.headline a").Attr("href")

		response, err := http.Get("https://graph.facebook.comxs/fql?q=SELECT%20share_count,like_count,comment_count,click_count%20FROM%20link_stat%20WHERE%20url='" + post.Url + "'")
		if err != nil {
			return
		}
		defer response.Body.Close()

		fbJson, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}

		var fbData interface{}
		if err = json.Unmarshal(fbJson, &fbData); err != nil {
			return
		}

		fqlRoot, ok := fbData.(map[string]interface{})
		if !ok || fqlRoot["data"] == nil {
			return
		}

		fqlDataItems, ok := fqlRoot["data"].([]interface{})
		if !ok || len(fqlDataItems) == 0 {
			return
		}

		fqlFields, ok := fqlDataItems[0].(map[string]interface{})
		if !ok {
			return
		}

		post.Shares = uint16(fqlFields["share_count"].(float64))
		post.Likes = uint16(fqlFields["like_count"].(float64))
		post.Comments = uint16(fqlFields["comment_count"].(float64))
	})

	var file *os.File
	if file, err = os.Create("./data/" + captureEvent.CapturedAt.Format(time.RFC3339)  + ".json"); err != nil {
		return errors.New("Failed to create file")
	}
	defer file.Close()

	var b []byte
	if b, err = json.Marshal(captureEvent); err != nil {
		return errors.New("Failed to marshal JSON")
	}

	if _, err := file.WriteString(string(b)); err != nil {
		return errors.New("Failed to write JSON")
	}

	return nil
}

func main() {
	var err error
	parsePosts()

	c := time.Tick(60 * 5 * time.Second)
	for _ = range c {
		fmt.Println("ticking:", time.Now().Format(time.RFC3339))
		if err = parsePosts(); err != nil {
			fmt.Println("Error processing post: %s", err)
		}
	}
}