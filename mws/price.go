package mws

import (

	"os"
	"log"
	"time"
	"net/http"
	"strconv"
	"strings"
	"github.com/Azunyan1111/amazonDatabaseUpdate/db"
	"github.com/Azunyan1111/gomws/mws/products"
	"github.com/Azunyan1111/gomws/gmws"
)

//*/
var client *products.Products



// APIの初期設定
func ApiInit() {
	// API key config
	config := gmws.MwsConfig{
		SellerId:  os.Getenv("SellerId"),
		AccessKey: os.Getenv("AccessKey"),
		SecretKey: os.Getenv("SecretKey"),
		Region:    "JP",
	}
	var err error
	// Create client
	client, err = products.NewClient(config)
	if err != nil {
		log.Println(err)
		return
	}
}

func main_() {
	ApiInit()
	db.DataBaseInit()
	GetPrice()
}

// go func only. 1 day
func GetPrice() {
	//*/ // ASINのリストをもらってくる
	allASIN, err := db.SelectAllForASINLimit864000()
	if err != nil {
		log.Println(err)
		return
	}
	//*/
	_ = []string{"000638420X",
		"0006388833",
		"0006754163",
		"0007133103",
		"0007145721",
		"0007152590",
		"0007171803",
		"0007177771",
		"0007178522",
		"0007179251",
		"0007183062",
		"0007248075",
		"0007267126",
		"000727596X",
		"0007279779",
		"0007301987",
		"0007305877",
		"000735634X",
		"000737139X",
		"0007371411",
		"000737142X",
		"0007382197",
		"0007395825",
		"0007417977",
		"000742325X",
		"0007423268",
		"0007423276",
		"0007440650",
		"0007460597",
		"0007460600",
		"0007460627",
		"0007462913",
		"0007464665",
		"000746519X",
		"0007466072",
		"0007481543",
		"0007485905",
		"0007493118",
		"0007498829",
		"0007499663",
		"0007505833",
		"000750621X",
		"0007522746",
		"000753812X",
		"0007544103",
		"0007548699",
		"0007554850",
		"0007556217",
		"0007562667",
		"0007562675"}

	// ASINリストを最大の20個づつに分けて変数に格納する
	var asinDoubleArray [][]string
	var tempArray []string
	for i, asin := range allASIN {
		if i%21 == 0 {
			if len(tempArray) == 0 {
				continue
			}
			asinDoubleArray = append(asinDoubleArray, tempArray)
			tempArray = []string{}
		} else {
			tempArray = append(tempArray, asin)
		}
	}

	// 20個セットをリクエストで送る。
	for _, asinArray := range asinDoubleArray {
		//*
		start := time.Now()
		response := client.GetLowestOfferListingsForASIN(asinArray, gmws.Parameters{"ItemCondition": "New"})
		if response.Error != nil || response.StatusCode != http.StatusOK {
			log.Println("http Status:" + string(response.StatusCode))
			log.Println(response.Error)
			return
		}

		// 返ってきたXMLをパースする。
		xmlNode, _ := gmws.GenerateXMLNode(response.Body)
		if gmws.HasErrors(xmlNode) {
			log.Println(gmws.GetErrors(xmlNode))
			continue
		}

		// 保存する商品在庫情報を格納する
		var saveProduct []db.ProductStock

		// 全ての商品の在庫情報をもらってくる
		products := xmlNode.FindByKey("GetLowestOfferListingsForASINResult")
		// 複数の商品情報から一つづつ商品ごとにさばいていく
		for _, product := range products {
			// 全ての在庫情報を取得する
			stocks := product.FindByPath("Product.LowestOfferListings.LowestOfferListing")
			// パース時の時間を取得する
			insertTime := time.Now().Unix()

			// 複数ある商品在庫情報を一つづつ捌く
			for _, stock := range stocks {
				amount,err := strconv.ParseInt(strings.Split(stock.FindByPath("Price.LandedPrice.Amount")[0].Value.(string),".")[0],10,64)
				if err != nil{
					log.Println(err)
					continue
				}
				temp := db.ProductStock{
					ASIN:         product.FindByPath("Product.Identifiers.MarketplaceASIN.ASIN")[0].Value.(string),
					Amount:       amount,
					Channel:      stock.FindByPath("Qualifiers.FulfillmentChannel")[0].Value.(string),
					Conditions:   stock.FindByPath("Qualifiers.ItemCondition")[0].Value.(string),
					ShippingTime: stock.FindByPath("Qualifiers.ShippingTime.Max")[0].Value.(string),
					InsertTime:   insertTime,
				}
				// 最安値以外も格納してしまう問題を解決する
				if !isInArray(saveProduct, temp){
					saveProduct = append(saveProduct, temp)
				}
			}
		}
		db.InsertProductPrice(saveProduct)
		log.Println(saveProduct)

		end := time.Now()
		// 2秒に一回のリクエストでAPI制限を受けない。
		if (end.Sub(start)).Seconds() < 2 {
			time.Sleep(2 * time.Second)
		}
		//*/
	}
}
// 配列の中にすでに商品次用法があるかを確認する（最安値は一番最初なので入ってない）
func isInArray(s []db.ProductStock, e db.ProductStock) bool {
	for _, v := range s {
		if e.ASIN == v.ASIN {
			return true
		}
	}
	return false
}