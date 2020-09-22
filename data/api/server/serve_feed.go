package server

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"fullRssAPI/model/db"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/html"
)

func ServeFeed(feedInfo *db.FeedInfo) error {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feedInfo.URL)
	newItems := []*feeds.Item{}

	oldFile, err := os.Open(fmt.Sprintf("./rss/%s.xml", feedInfo.ID.String()))
	defer oldFile.Close()
	if err != nil {
		if os.IsNotExist(err) {
			createItems(&newItems, feed, feedInfo)
		} else {
			return err
		}
	} else {
		oldFeed, _ := fp.Parse(oldFile)
		updateItems(&newItems, feed, oldFeed, feedInfo)
	}
	newFeed := createFeed(feed)
	newFeed.Items = newItems

	err = writeFeed(newFeed, feedInfo.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetTitle(feedInfo *db.FeedInfo) (string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedInfo.URL)
	if err != nil {
		return "", err
	}
	return feed.Title, nil
}

func getFullText(url string, XPath string) string {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Fatal(err)
	}
	content := htmlquery.FindOne(doc, XPath)
	fullText := htmlquery.OutputHTML(content, true)
	fullText = replaceImgur(content, fullText)
	fullText = replaceFont(content, fullText)
	return "<![CDATA[" + fullText + "]]>"
}

func replaceFont(content *html.Node ,fullText string) string {
	fontContents := htmlquery.Find(content, "//font")
	convFont := map[string]string {
		"1": "10px",
		"2": "13px",
		"3": "16px",
		"4": "18px",
		"5": "24px",
		"6": "32px",
		"7": "48px",
	}
	for _, fontContent := range fontContents {
		var color, size string 
		text := htmlquery.OutputHTML(fontContent, false)
		colorNode := htmlquery.FindOne(fontContent, "//@color")
		if colorNode != nil {
			color = "color: " + htmlquery.InnerText(colorNode) + ";"
		}
		sizeNode := htmlquery.FindOne(fontContent, "//@size")
		if sizeNode != nil {
			size = "font-size: " + convFont[htmlquery.InnerText(sizeNode)] + ";"
		}
		html := fmt.Sprintf("<span style=\"%s%s\">%s</span>", color, size, text)
		fullText = strings.Replace(fullText, htmlquery.OutputHTML(fontContent, true), html, 1)
	}
	return fullText
}

// imgurが存在するときは画像を埋め込む
func replaceImgur(content *html.Node, fullText string) string {
	imgurPath := "//blockquote[@class='imgur-embed-pub']"
	imgurContents := htmlquery.Find(content, imgurPath)
	for _, imgurContent := range imgurContents {
		html, err := getImgurHTML(imgurContent)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fullText = strings.Replace(fullText, htmlquery.OutputHTML(imgurContent, true), html, 1)
	}
	return fullText
}

func getImgurHTML(content *html.Node) (string, error) {
	urlNode := htmlquery.FindOne(content, "//a/@href")
	if urlNode == nil {
		return "", errors.New("Not Found")
	}
	
	doc, err := htmlquery.LoadURL("https:" + htmlquery.InnerText(urlNode))
	if err != nil {
		return "", err
	}

	imgurPaths := map[string]string {
		"imgXPath": "//link[@rel='image_src']/@href",
		"gifXPath": "//meta[@property='og:image']/@content",
		"videoXPath": "//meta[@property='og:video']/@content",
	}

	var imgurURLNode *html.Node
	var html string
	for key, value := range imgurPaths {
		imgurURLNode = htmlquery.FindOne(doc, value)
		if imgurURLNode != nil {
			url := strings.Split(htmlquery.InnerText(imgurURLNode), "?")[0]
			if key == "videoXPath" {
				html = fmt.Sprintf("<video src=%s controls></video>", url)
				break
			} 
			html = fmt.Sprintf("<img src=%s />", url)
			break
		}
	}
	if imgurURLNode == nil {
		return "", errors.New("Not Found")
	}
	return html, nil
}

func createFeed(feed *gofeed.Feed) *feeds.Feed {
	newFeed := &feeds.Feed{
		Title:       feed.Title,
		Link:        &feeds.Link{Href: feed.Link},
		Description: feed.Description,
	}
	return newFeed
}

func updateItems(newItems *[]*feeds.Item, feed *gofeed.Feed, oldFeed *gofeed.Feed, feedInfo *db.FeedInfo) {
	var newItem *feeds.Item
	for _, item := range feed.Items {
		itemExisted := false
		for _, oldItem := range oldFeed.Items {
			if item.Link == oldItem.Link {
				itemExisted = true
				newItem = convToFeedsItem(oldItem)
				break
			}
		}
		if itemExisted {
			*newItems = append(*newItems, newItem)
		} else {
			newItem = createItem(item, feedInfo.XPath)
			*newItems = append(*newItems, newItem)
		}
	}
}

func convToFeedsItem(item *gofeed.Item) *feeds.Item {
	return &feeds.Item{
		Title:       item.Title,
		Link:        &feeds.Link{Href: item.Link},
		Description: item.Description,
		Created:     *item.PublishedParsed,
	}
}

func createItems(newItems *[]*feeds.Item, feed *gofeed.Feed, feedInfo *db.FeedInfo) {
	for _, item := range feed.Items {
		newItem := createItem(item, feedInfo.XPath)
		*newItems = append(*newItems, newItem)
	}
}

func createItem(item *gofeed.Item, XPath string) *feeds.Item {
	fullText := getFullText(item.Link, XPath)
	feeditem := &feeds.Item{
		Title:       item.Title,
		Link:        &feeds.Link{Href: item.Link},
		Description: fullText,
		Created:     *item.PublishedParsed,
	}
	return feeditem
}

func writeFeed(feed *feeds.Feed, id uuid.UUID) error {
	rss, err := feed.ToRss()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(fmt.Sprintf("./rss/%s.xml", id.String()), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Fprintln(file, rss)
	return nil
}