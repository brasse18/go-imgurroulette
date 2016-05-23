package main

import (
	"fmt"
	"net/http"
	"html/template"
	"gitlab.com/Niesch/go-imgurroulette/imgur"
)

type Result struct {
	Link  string
	Tries int
}

func indexHandler(w http.ResponseWriter, r *http.Request, maxtries int, minlength int, maxlength int, dbg bool) {
	i := imgur.New("https://imgur.com/", "https://i.imgur.com/",".png", maxtries, minlength, maxlength, dbg)

	link, tries, err := i.FindValidGalleryLink()
	if err != nil {
		i.ErrorLogger.Println(err)
		fmt.Fprintf(w, "Something went wrong. The error has been logged.\n")
		return
	}
	ilink := i.BuildImageLink(link)
	t := template.Must(template.New("index.html").ParseFiles("assets/index.html"))
	result := Result{Link: ilink, Tries: tries}
	err = t.Execute(w, result)
	if err != nil {
		i.ErrorLogger.Println(err)
		fmt.Fprintf(w, "Something went wrong. The error has been logged.\n")
	}
}
