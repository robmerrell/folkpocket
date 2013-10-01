package cmds

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/garyburd/redigo/redis"
	"github.com/robmerrell/folkpocket/config"
	"net/url"
)

var startingUrl = "http://folklore.org/StoryView.py?project=Macintosh&story=I'll_Be_Your_Best_Friend.txt&sortOrder=Sort%20by%20Date&detail=medium"

var ScrapeDoc = `
Scrape all of the URLs from folklore.org starting at a story (usually the first) and then by finding the url
for the "next" button and following that until we have reached the end.
`

func ScrapeAction() error {
	// build up the storyUrls list
	urls, err := getNextUrl(startingUrl, []string{startingUrl})
	if err != nil {
		return err
	}

	return saveUrls(urls, config.Env().Get("cacheKey").(string))
}

func getNextUrl(currentUrl string, urlAccumulator []string) ([]string, error) {
	doc, err := goquery.NewDocument(currentUrl)
	if err != nil {
		return urlAccumulator, err
	}

	nextImage := doc.Find("img[src*='images/rightarrow.gif']")
	if nextImage.Length() == 0 {
		return urlAccumulator, nil
	}
	nextUrl, exists := nextImage.First().Parent().Attr("href")

	if exists {
		u, _ := url.Parse(nextUrl)
		if !u.IsAbs() {
			u.Host = "folklore.org"
			u.Scheme = "http"
			u.Path = "/" + u.Path
		}

		nextUrl = u.String()
		fmt.Println("Found", nextUrl)
		return getNextUrl(nextUrl, append(urlAccumulator, nextUrl))
	} else {
		return urlAccumulator, nil
	}
}

func saveUrls(storyUrls []string, redisKey string) error {
	// save the found urls in the database
	redisC, err := connectToRedis(config.Env().Get("redishost").(string))
	if err != nil {
		return err
	}
	defer redisC.Close()

	// clear out the list if it already exists
	if _, err := redisC.Do("DEL", redisKey); err != nil {
		return err
	}

	// save the urls
	for _, url := range storyUrls {
		if _, err := redisC.Do("RPUSH", redisKey, url); err != nil {
			return err
		}
	}

	return nil
}

func connectToRedis(addr string) (redis.Conn, error) {
	return redis.Dial("tcp", addr)
}
