package main

func hoge(){

}
/*
import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"log"
	"regexp"
)
var titleR = regexp.MustCompile(`.*: `)
var asinR = regexp.MustCompile(`/[A-Z0-9]{10}/`)
var urlR = regexp.MustCompile(`src="http.*\.jpg`)

func main() {
	// 初期設定"
	timeout := 5
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	// RSS を取ってくる
	uri := "https://www.amazon.co.jp/gp/rss/bestsellers/videogames/2494235051"
	err := feed.Fetch(uri, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
// RSS チャンネルのハンドラ(自作RSSとかの複数サイトの情報が出た時にチャンネルが分かれる(Amazonはひとつ))
func chanHandler(feed *rss.Feed, newChannels []*rss.Channel) {
	//fmt.Printf("%d new channel(s) in %s\n", len(newChannels), feed.Url)
}

// RSSのアイテムを色々するハンドラ
func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	// アイテム数を出力
	//fmt.Printf("%d new item(s) in %s\n", len(newItems), feed.Url)

	// アイテム一つひとつに対する操作
	for _, item := range newItems {

		titleRe := titleR.Copy()
		asinRe := asinR.Copy()
		urlRe := urlR.Copy()
		title := titleRe.ReplaceAllString(item.Title,"")
		ASIN := asinRe.FindStringSubmatch(item.Links[0].Href)[0][1:11]
		imageUrl := urlRe.FindStringSubmatch(item.Description)[0][5:]

		fmt.Println(title)
		fmt.Println(ASIN)
		fmt.Println(imageUrl)
	}
}