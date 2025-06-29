package models

import "time"

type Order struct {
	OrderUID          string
	TrackNumber       string
	Entry             string
	DeliveryInfo      DeliveryInfo
	Payment           Payment
	Items             []Item
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	Shardkey          string
	SmID              int64
	DateCreated       time.Time
	OofShard          string
}

type DeliveryInfo struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type Payment struct {
	Transaction  string
	RequestID    int64
	Currency     string
	Provider     string
	Amount       int
	PaymentDT    int64
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

type Item struct {
	ChrtID      int64
	TrackNumber string
	Price       int
	Rid         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int
	NmID        int64
	Brand       string
	Status      int
}
