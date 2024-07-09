package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	// "gopkg.in/yaml.v2"

	"github.com/gocolly/colly/v2"
	"github.com/tanishqtrivedi27/media-scraper/internal/storage"
)

const basedir = "/home/tanishq/projects/media-scraper/reservoir"
const websiteurl = "https://news.google.com/"

func main() {
	redisConfig := storage.RedisConfig{
		Address:  "localhost:6379",
		Password: "",
		DBName:   0,
	}
	rdb, _ := storage.NewRedisStorage(redisConfig)
	redisStore, _ := storage.NewStorage(rdb)

	dbconfig := storage.PostgresConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		DBName:   "mediascraper",
		Port:     5432,
	}
	db, _ := storage.NewPostgreSQLStorage(dbconfig)
	postgreStore, _ := storage.NewDBStorage(db)

	scrapeWebsite(redisStore, postgreStore)
}

func scrapeWebsite(redisStore *storage.Storage, postgreStore *storage.DBStorage) {
	c := colly.NewCollector(colly.MaxDepth(2))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.66 Safari/537.36")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", websiteurl)
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		imgURL := e.Attr("src")
		if imgURL != "" {
			added, err := redisStore.StoreUrl(imgURL)
			if err != nil {
				fmt.Printf("Failed to store URL: %v\n", err)
			} else if added {
				// fmt.Printf("Image: %v\n", imgURL)
				downloadImage(imgURL, postgreStore)
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Scraped:", r.Request.URL)
	})

	err := c.Visit(websiteurl)
	if err != nil {
		log.Fatal(err)
	}
}

func buildFileName(imgURL string) string {
	fileUrl, err := url.Parse(imgURL)
	if err != nil {
		panic(err)
	}
	path := fileUrl.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	fmt.Println(fileName)
	return fileName
}

func downloadImage(imgURL string, postgreStore *storage.DBStorage) {
	response, err := http.Get(imgURL)
	if err != nil {
		fmt.Print(err)
	}

	defer response.Body.Close()
	fileName := buildFileName(imgURL)
	fileName = filepath.Join(basedir, fileName)

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}

	postgreStore.InsertTuple(imgURL, fileName)
}
