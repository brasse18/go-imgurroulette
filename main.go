package main

import (
	"net/http"
	"flag"
)

func main(){
	var maxtries, minlength, maxlength int
	maxtries = 50
	minlength = 5
	maxlength = 7
	flag.IntVar(&maxtries, "maxtries", 50, "how many attempts should be made while finding a valid URL")
	flag.IntVar(&minlength, "minlength", 5, "minimum length of imgur URLs")
	flag.IntVar(&maxlength, "maxlength", 7, "maximum length of imgur URLs")
	debug := flag.Bool("debug", false, "debug to stdout")
	flag.Parse()
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, maxtries, minlength, maxlength, *debug)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}
