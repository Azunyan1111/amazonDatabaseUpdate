package db


import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"github.com/pkg/errors"
)

var MyDB *sql.DB

func DataBaseInit() {
	dataSource := os.Getenv("DATABASE_URL")
	var err error
	MyDB, err = sql.Open("mysql", dataSource) //"root:@/my_database")
	if err != nil {
		panic(err)
	}
}

// get rank urls
// WANG this is 10 time second. only go func{}()
func SelectAllUrl() ([]string, error) {
	rows, err := MyDB.Query("SELECT URL FROM CategoryURL ORDER BY RAND()")
	if err != nil {
		return nil, err
	}
	// list append
	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func SelectAllForASINLimit864000() ([]string, error) {
	// new connection
	dataSource := os.Getenv("DATABASE_URL")
	myDB, err := sql.Open("mysql", dataSource) //"root:@/my_database")
	if err != nil {
		return nil, errors.New("Can not connection Database")
	}

	// query. API MAX 86500 / day
	rows, err := myDB.Query("SELECT ASIN FROM Items WHERE title IS NOT NULL ORDER BY RAND() LIMIT 864000")
	if err != nil {
		return nil, err
	}
	// list append
	var asins []string
	for rows.Next() {
		var asin string
		if err := rows.Scan(&asin); err != nil {
			return nil, err
		}
		asins = append(asins, asin)
	}
	defer myDB.Close()
	return asins, nil
}

func InsertProductPrice(asins []ProductStock) {
	for _, asin := range asins {
		_, err := MyDB.Exec("INSERT INTO Price(ASIN,Amount,Channel,Conditions,ShippingTime,InsertTime)"+
			" VALUES(?,?,?,?,?,?)", asin.ASIN, asin.Amount, asin.Channel, asin.Conditions, asin.ShippingTime, asin.InsertTime)
		if err != nil {
			continue
		}
	}
}
func InsertNewProduct(product Product) {
	_, err := MyDB.Exec("INSERT INTO Items(ASIN,Title,Image) VALUES(?,?,?)", product.ASIN, product.Title, product.Image)
	if err != nil {
		return
	}
}