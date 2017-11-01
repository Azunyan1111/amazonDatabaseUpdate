package main

import (
	rss "github.com/jteeuwen/go-pkg-rss"
	"regexp"
	"log"
	"./db"
	"./mws"
	"time"
)

var titleR = regexp.MustCompile(`.*: `)
var asinR = regexp.MustCompile(`/[A-Z0-9]{10}/`)
var urlR = regexp.MustCompile(`src="http.*\.jpg`)
var rssUrlR = regexp.MustCompile(`/gp/bestsellers/`)


var product chan db.Product

// 初期設定
var timeout = 5

func main(){

	// 初期設定
	product = make(chan db.Product)
	go insertCh(product)

	mws.ApiInit()
	db.DataBaseInit()

	//*/
	endCh := make(chan int)

	// TODO: GetPriceは1日1回起動しないといけない曲者
	go func() {
		mws.GetPrice()
		time.Sleep(24 * time.Hour)
	}()
	// TODO: 一時間に1回でOK!
	go func(){
		getProduct()
		time.Sleep(1 * time.Hour)
	}()

	<-endCh
	//*/
}

func getProduct(){
	allUrl, err := db.SelectAllUrl()
	if err != nil{
		log.Println(err)
		return
	}

	allRssUrl := []string{}
	bar := rssUrlR.Copy()
	for _, foo := range allUrl{
		bestsellers := bar.ReplaceAllString(foo,"/gp/rss/bestsellers/")
		moversAndShakers := bar.ReplaceAllString(foo,"/gp/rss/movers-and-shakers/")
		newReleases := bar.ReplaceAllString(foo,"/gp/rss/new-releases/")
		mostWishedFor := bar.ReplaceAllString(foo,"/gp/rss/most-wished-for/")
		mostGifted := bar.ReplaceAllString(foo,"/gp/rss/most-gifted/")
		allRssUrl = append(allRssUrl,bestsellers)
		allRssUrl = append(allRssUrl,moversAndShakers)
		allRssUrl = append(allRssUrl,newReleases)
		allRssUrl = append(allRssUrl,mostWishedFor)
		allRssUrl = append(allRssUrl,mostGifted)
	}

	// 最大のgoroutineの数を制限する
	c := make(chan bool, 2000)
	for _,foo := range allRssUrl {
		// もしcが一杯ならこの行で待たされる
		c <- true
		go func(url string) {
			err := getASIN(url)
			if err != nil{
				log.Println(err)
			}
			// cから読みだしてその値を捨てる -- ほかのgoroutineのための空きを作る
			defer func() { <-c }()
		}(foo)
	}
}

func getASIN(url string)error{
	// RSS ゲッター
	feed := rss.New(timeout, true, chanHandler, itemHandler)
	// RSS を取ってくる
	err := feed.Fetch(url, nil)
	if err != nil {
		return err
	}
	return nil
}

// RSS チャンネルのハンドラ(自作RSSとかの複数サイトの情報が出た時にチャンネルが分かれる(Amazonはひとつ))
func chanHandler(feed *rss.Feed, newChannels []*rss.Channel) {
	//fmt.Printf("%d new channel(s) in %s\n", len(newChannels), feed.Url)
}

// RSSのアイテムを色々するハンドラ
func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	// アイテム数を出力
	//fmt.Printf("%d new item(s) in %s\n", len(newItems), feed.Url)
	if len(newItems) == 1 {
		return
	}
	// アイテム一つひとつに対する操作
	for _, item := range newItems {

		if item.Links[0].Href == "http://www.amazon.co.jp"{
			continue
		}

		titleRe := titleR.Copy()
		asinRe := asinR.Copy()
		urlRe := urlR.Copy()
		title := titleRe.ReplaceAllString(item.Title,"")
		ASIN := asinRe.FindStringSubmatch(item.Links[0].Href)[0][1:11]

		// 画像が存在しない場合
		var imageUrl string
		if len(urlRe.FindStringSubmatch(item.Description)) == 0{
			imageUrl = "http://ec1.images-amazon.com/images/G/09/nav2/dp/no-image-no-ciu.gif"
		}else {
			imageUrl = urlRe.FindStringSubmatch(item.Description)[0][5:]
		}

		var p db.Product
		p.Title = title
		p.Image = imageUrl
		p.ASIN = ASIN
		product <- p
	}
}

func insertCh(myCh chan db.Product){
	//count := 0
	for out := range myCh {
		db.InsertNewProduct(out)
		//log.Println(count)
		//count++
	}
}
