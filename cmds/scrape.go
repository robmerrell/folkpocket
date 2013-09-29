package cmds

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/garyburd/redigo/redis"
)

var storyUrls []string

var startingUrl = "http://folklore.org/StoryView.py?project=Macintosh&story=I'll_Be_Your_Best_Friend.txt&sortOrder=Sort%20by%20Date&detail=medium"

var ScrapeDoc = `
Scrape all of the URLs from folklore.org starting at a story (usually the first) and then by finding the url
for the "next" button and following that until we have reached the end.
`

func ScrapeAction() error {
	storyUrls = make([]string, 0)

	// build up the storyUrls list
	storyUrls = append(storyUrls, startingUrl)
	err := getNextUrl(startingUrl)
	if err != nil {
		return err
	}

	return saveUrls()
}

func getNextUrl(currentUrl string) error {
	doc, err := goquery.NewDocument(currentUrl)
	if err != nil {
		return err
	}

	nextImage := doc.Find("img[src*='images/rightarrow.gif']")
	if nextImage.Length() == 0 {
		return nil
	}
	nextUrl, exists := nextImage.First().Parent().Attr("href")

	fmt.Println("Found", nextUrl)

	if exists {
		nextUrl = "http://folklore.org/" + nextUrl
		storyUrls = append(storyUrls, nextUrl)
		return getNextUrl(nextUrl)
	} else {
		return nil
	}
}

func saveUrls() error {
	// save the found urls in the database
	redisC, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer redisC.Close()

	// clear out the list if it already exists
	if _, err := redisC.Do("DEL", "cachedFolkloreUrls"); err != nil {
		return err
	}

	// save the urls
	for _, url := range storyUrls {
		if _, err := redisC.Do("RPUSH", "cachedFolkloreUrls", url); err != nil {
			return err
		}
	}

	return nil
}
