package imgur

import (
	"net/http"
	"math/rand"
	"time"
	"fmt"
	"log"
	"os"
)

type ImgurAnonymousClient struct {
	httpClient            http.Client
	InfoLogger            *log.Logger
	DebugLogger           *log.Logger
	ErrorLogger           *log.Logger
	CacheChan             chan *ImgurResult
	cfg                   *Config
}

type Config struct {
	DefaultFileExtension string
	AlbumBaseUrl         string
	DirectBaseUrl        string
	MaxTries             int
	MinLength            int
	MaxLength            int
	CacheSize            int
	Debug                bool
}

type ImgurResult struct {
	Link     string
	Tries    int
}


// Returns a brand spankin' new ImgurAnonymousClient 
func New(cfg *Config) (*ImgurAnonymousClient) {
	il := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	dl := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	el := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	c := make(chan *ImgurResult, cfg.CacheSize)
	return &ImgurAnonymousClient{httpClient: http.Client{}, InfoLogger: il, DebugLogger: dl, ErrorLogger: el, CacheChan: c, cfg: cfg}
} 

// Returns the first found valid gallery link, and the amount of tries.
func (client *ImgurAnonymousClient) FindValidGalleryLink() (string, int, error) {
	var i int
	for i = 0; i <= client.cfg.MaxTries; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		// We check against the "album" URL and not the direct file, since accessing a direct removed file will return 200 OK and removed.png.
		l := (rand.Intn(client.cfg.MaxLength+1 - client.cfg.MinLength) + client.cfg.MinLength)
		s := randomString(l)
		// yeah, don't hate me 
		url := client.cfg.AlbumBaseUrl + s
		if client.CheckLink(url) == nil {
			if client.cfg.Debug {
				client.DebugLogger.Printf("Found valid image URL: %s\n", url)
			}
			return client.cfg.DirectBaseUrl + s, i, nil
		}
	}
	return "", i, fmt.Errorf("Failed to find valid URL")
}

// Returns the direct image link from string "gallery". 
func (client *ImgurAnonymousClient) BuildImageLink(gallery string) (string) {
	return gallery + client.cfg.DefaultFileExtension
}

// Checks whether url exists or not. Returns nil upon success. 
func (client *ImgurAnonymousClient) CheckLink(url string) error {
	resp, err := client.httpClient.Head(url)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		if client.cfg.Debug {
			client.DebugLogger.Println("Found URL: %s\n", url)
		}
		return nil
	} else {
		if client.cfg.Debug {
			client.DebugLogger.Printf("Got status code %d on %s\n", resp.StatusCode, url)
		}
		return fmt.Errorf("error: Non-200 status code %d on URL %s", resp.StatusCode, url)
	}
}

// Returns a random string of length "length".
// imgur links are 5-7 characters from what I can tell.
func randomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
