package db

type Product struct {
	Title string
	ASIN string
	Image string
}

type ProductStock struct {
	ASIN         string
	Amount       int64
	Channel      string
	Conditions   string
	ShippingTime string
	InsertTime   int64
}