package main

import (
       "net/http"
       "fmt"
       "io/ioutil"
       "html"
)

func main() {

     //var ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/537.17 (KHTML, like Gecko) Chrome/24.0.1309.0 Safari/537.17"
     response, error := http.Get("http://www.techcrunch.com")
     //fmt.Println(error, response)

     if error != nil {
        fmt.Println("something went wrong")
        return
     }

     defer response.Body.Close()
     body, error := ioutil.ReadAll(response.Body)
     //fmt.Println(body)

     //stringBody := string(body)
     //fmt.Println(stringBody)

     t := html.NewTokenizer(body)
     fmt.Println(t)
/*
     for {

         tt := t.Next()
         //if tt == html.ErrorToken {
         //   return
         //}

         switch tt {
             case ErrorToken:
                 fmt.Println("Found the end")
                 return t.Err()
         }
     }*/
}