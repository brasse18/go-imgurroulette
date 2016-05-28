package main

import (
	"flag"
	"fmt"
	"gitlab.com/Niesch/go-imgurroulette/imgur"
	"net/http"
	"os"
)

func main() {
	var maxtries, minlength, maxlength, cachesize, workers int
	maxtries = 50
	minlength = 5
	maxlength = 7
	cachesize = 50
	workers = 25

	flag.IntVar(&maxtries, "maxtries", 50, "how many attempts should be made while finding a valid URL")
	flag.IntVar(&minlength, "minlength", 5, "minimum length of imgur URLs")
	flag.IntVar(&maxlength, "maxlength", 7, "maximum length of imgur URLs")
	flag.IntVar(&cachesize, "cachesize", 50, "the amount of items to try to keep cached")
	flag.IntVar(&workers, "workers", 25, "the amount of workers to work on keeping cache filled")
	debug := flag.Bool("debug", false, "debug to stdout")
	license := flag.Bool("license", false, "show license information")
	flag.Parse()

	if *license {
		fmt.Println(`    go-imgurroulette
    Copyright (C) 2016  Niesch

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.`)
		os.Exit(0)
	}

	i := imgur.New(&imgur.Config{DefaultFileExtension: ".png", AlbumBaseUrl: "https://imgur.com/", DirectBaseUrl: "https://i.imgur.com/", MaxTries: maxtries, MinLength: minlength, MaxLength: maxlength, CacheSize: cachesize, Debug: *debug})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, i)
	})

	go func(client *imgur.ImgurAnonymousClient) {
		for j := 0; j < workers; j++ {
			if client.Cfg.Debug {
				client.DebugLogger.Printf("Starting worker %d out of %d\n", j, workers)
			}
			go func(client *imgur.ImgurAnonymousClient) {
				for {
					link, tries, err := client.FindValidGalleryLink()
					if err != nil {
						client.ErrorLogger.Println(err)
					}
					ilink := client.BuildImageLink(link)
					client.CacheChan <- &imgur.ImgurResult{Link: ilink, Tries: tries}
				}
			}(client)
		}
	}(i)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}
