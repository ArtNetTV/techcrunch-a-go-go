package main

/*

Questions to answer:
 - what is the average time an article is on the techcrunch homepage
 - what is the average time an article is above the fold
 - authors ordered by article count
 - author with most minutes on the homepage

{
    time: "2012-11-15T16:30:00",
    posts: [
        {
            Id: "asdf",
            Title: "xxx",
            Author: "me",
            Url: "http://xxx",
            Shares: 0,
            Likes: 0
            Comments: 0
        },
        ...
    ]
}

*/

import (
       "fmt"
       //"foobar"
       //"exp/html"
       "net/http"
       "strings"
       "io/ioutil"
       "github.com/PuerkitoBio/goquery"
       "encoding/json"
       "time"
       "os"
)

//comment
type Post struct {
     Id string
     Title string
     Author string
     Url string
     Shares uint16
     Likes uint16
     Comments uint16
}

//comment
type CaptureEvent struct {
     CapturedAt time.Time
     Posts []Post
}

func parsePosts() {
    var document *goquery.Document
    var err error

    if document, err = goquery.NewDocument("http://www.techcrunch.com"); err != nil {
        fmt.Println("sheeeet")
    }


    matches := document.Find("div .post")
    var captureEvent CaptureEvent
    captureEvent.Posts = make([]Post, matches.Length())
	captureEvent.CapturedAt = time.Now()

    matches.Each(func(i int, s *goquery.Selection) {

        //TODO: Defer?
        var post *Post = &captureEvent.Posts[i];

        //TODO :Errors
        post.Id, _ = s.Attr("id")
        post.Author = s.Find("a span.name").Text();
        if post.Author == "" {
            post.Author = s.Find("div.by-line div.by-line").Text()
        }

        post.Title = strings.TrimSpace(s.Find("h2.headline a").Text());
        post.Url, _ = s.Find("h2.headline a").Attr("href")

        response, _ := http.Get("https://graph.facebook.com/fql?q=SELECT%20share_count,like_count,comment_count,click_count%20FROM%20link_stat%20WHERE%20url='" + post.Url + "'")
        //    fmt.Println('whoops')
        //}

        xx, _ := ioutil.ReadAll(response.Body)

        var fbData interface{}
        _ = json.Unmarshal(xx, &fbData);

        // We can create types, but we don't need to since we can just do type assertions instead
        // of declaring explicit types
        zz := fbData.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})

        post.Shares = uint16(zz["share_count"].(float64))
        post.Likes = uint16(zz["like_count"].(float64))
        post.Comments = uint16(zz["comment_count"].(float64))

        //fmt.Println(post)
    })

    //fmt.Println(captureEvent)

    //TODO: Add time to front of data
    var file *os.File
    if file, err = os.Create("./data/" + captureEvent.CapturedAt.Format(time.RFC3339)  + ".json"); err != nil {
        fmt.Println("Failed to create file")
    }
    defer file.Close()

    var b []byte
    if b, err = json.Marshal(captureEvent); err != nil {
        fmt.Println("Failed to marshal")
    }

    //TODO: Is this the best way
    if _, err := file.WriteString(string(b)); err != nil {
        fmt.Println("Failed to write")
    }

	/*
    var file2 *os.File
    if file2, err = os.Open("./data/something.json"); err != nil {
        fmt.Println("Could not open")
    }

    var fi os.FileInfo
    fi, _ = file2.Stat()
    fmt.Println(fi.Size())
	*/
    /*
    var d []byte
    d = make([]byte, fi.Size())
    _, _ = file2.Read(d)
    var readCE CaptureEvent
    _ = json.Unmarshal(d, &readCE)

    fmt.Println(readCE.CapturedAt)
    */
    //fmt.Println("All done")
    //TODO: Write output to disk
}

func main() {
         parsePosts()
	 c := time.Tick(60 * 5 * time.Second);
	 for _ = range c {
	     fmt.Println("ticking:", time.Now().Format(time.RFC3339))
     	     parsePosts()	  
	 }
}