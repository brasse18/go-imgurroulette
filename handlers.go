package main

import (
	"fmt"
	"net/http"
	"html/template"
	"gitlab.com/Niesch/go-imgurroulette/imgur"
)

type Result struct {
	Link        string
	Tries       int
	CacheLength int
}

func indexHandler(w http.ResponseWriter, r *http.Request, i *imgur.ImgurAnonymousClient) {
	
	imgurResult := <- i.CacheChan
	result :=  Result{Link: imgurResult.Link, Tries: imgurResult.Tries, CacheLength: len(i.CacheChan)}
	t := template.Must(template.New("index.html").ParseFiles("assets/index.html"))
	err := t.Execute(w, result)
	if err != nil {
		i.ErrorLogger.Println(err)
		fmt.Fprintf(w, "Something went wrong. The error has been logged.\n")
	}
}
