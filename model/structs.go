package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Quote struct
type Quote struct {
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	Customer   *Customer
	CustomerID bson.ObjectId `bson:"customerID" json:"customerID"`
	Discount   *Discount     `bson:"discount" json:"discount"`
	Features   string        `bson:"features" bson:"features"`
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Invoiced   bool          `bson:"invoiced" bson:"invoiced"`
	ItemCosts  *ItemCosts    `bson:"itemCosts" json:"itemCosts"`
	ItemIds    *ItemIds      `bson:"items"`
	Items      *Items        `bson:"is"`
	JobsheetID bson.ObjectId `bson:"jobsheetID" json:"jobsheetID"`
	Number     int           `bson:"number" json:"number"`
	Payments   []*Payment    `json:"payments"`
	Price      *Price        `bson:"quotePrice" bson:"quotePrice"`
	Revision   int           `bson:"version" bson:"version"`
	UpdatedAt  time.Time     `bson:"updatedAt" json:"updatedAt"`
}

// Address struct
type Address struct {
	Associate  string        `bson:"associate" json:"associate"`
	City       string        `bson:"city" json:"city"`
	CustomerID bson.ObjectId `bson:"customerID" json:"customerID"`
	PostalCode string        `bson:"postalCode" json:"postalCode"`
	Province   string        `bson:"provinceCode" json:"province"`
	Street1    string        `bson:"street1" json:"street1"`
	Type       string        `bson:"type" json:"type"`
}

// Customer struct
type Customer struct {
	Email string `bson:"email" json:"email"`
	Name  struct {
		First  string `bson:"first" json:"first"`
		Last   string `bson:"last" json:"last"`
		Spouse string `bson:"spouse" json:"spouse"`
	}
	Phones   []*Phone
	Address  *Address
	PhoneMap map[string]string
}

// Dim struct
type Dim struct {
	Decimal  float64 `bson:"decimal" json:"decimal"`
	Fraction string  `bson:"fraction" json:"fraction"`
	Inch     int     `bson:"inch" json:"inch"`
}

// Dims struct
type Dims struct {
	Height *Dim `bson:"height" json:"height"`
	Width  *Dim `bson:"width" json:"width"`
}

// Discount struct
type Discount struct {
	Description string  `bson:"description" json:"description"`
	Discount    float64 `bson:"discount" json:"discount"`
	Subtotal    float64 `bson:"subtotal" json:"subtotal"`
	Total       float64 `bson:"total" json:"total"`
}

// Group struct
type Group struct {
	Costs *GroupCosts    `bson:"costs" json:"costs"`
	Dims  *Dims          `bson:"dims" json:"dims"`
	Items []*GroupWindow `bson:"items" json:"items"`
	Specs bson.M         `bson:"specs" json:"specs"`
	Qty   int            `bson:"qty" json:"qty"`
	Rooms []string       `bson:"rooms" json:"rooms"`
}

// GroupWindow struct
type GroupWindow struct {
	Costs   *ItemCost `bson:"costs" json:"costs"`
	Dims    *Dims     `bson:"dims" json:"dims"`
	Product bson.M    `bson:"product" json:"product"`
	Qty     int       `bson:"qty" json:"qty"`
	Specs   bson.M    `bson:"specs" json:"specs"`
}

// GroupCosts struct
type GroupCosts struct {
	Windows     float64 `bson:"windows" json:"windows"`
	Trim        float64 `bson:"trim" json:"trim"`
	Options     float64 `bson:"options" json:"options"`
	InstallType float64 `bson:"installType" json:"installType"`
	Install     float64 `bson:"install" json:"install"`
	Unit        float64 `bson:"extendUnit" json:"extendUnit"`
	Total       float64 `bson:"extendTotal" json:"extendTotal"`
}

// ItemCost struct
type ItemCost struct {
	Total float64 `bson:"extendTotal" json:"extendTotal"`
	Unit  float64 `bson:"extendUnit" json:"extendUnit"`
}

// ItemCosts struct
type ItemCosts struct {
	Group    float64 `bson:"group" json:"group"`
	Other    float64 `bson:"other" json:"other"`
	Subtotal float64 `bson:"subtotal" json:"subtotal"`
	Window   float64 `bson:"window" json:"window"`
}

// ItemIds struct
type ItemIds struct {
	Group  []string `bson:"group" bson:"group"`
	Other  []string `bson:"other" bson:"other"`
	Window []string `bson:"window" bson:"window"`
}

// Items struct
type Items struct {
	Group  []*Group
	Other  []*Other
	Window []*Window
}

// Other struct
type Other struct {
	Costs       *ItemCost `bson:"costs" json:"costs"`
	Description string    `bson:"description" json:"description"`
	Qty         int       `bson:"qty" json:"qty"`
	Rooms       []string  `bson:"rooms" json:"rooms"`
	Specs       struct {
		Options  string `bson:"options" json:"options"`
		Location string `bson:"location" json:"location"`
	}
}

// Payment struct
type Payment struct {
	Amount float64   `bson:"amount" json:"amount"`
	Date   time.Time `bson:"createdAt" json:"date"`
	Type   string    `bson:"type" json:"type"`
}

// Payments struct
// type Payments []*Payment

// Phone struct
type Phone struct {
	CountryCode string `bson:"countryCode" json:"countryCode"`
	Number      string `bson:"number" json:"number"`
	Type        string `bson:"_id" json:"type"`
}

// Price struct
type Price struct {
	Outstanding float64 `bson:"outstanding" bson:"outstanding"`
	Payments    float64 `bson:"payments" bson:"payments"`
	Subtotal    float64 `bson:"subtotal" bson:"subtotal"`
	Tax         float64 `bson:"tax" bson:"tax"`
	Total       float64 `bson:"total" bson:"total"`
}

// Product struct
type Product struct {
	Name string `bson:"name" json:"name"`
}

// Window struct
type Window struct {
	Costs       *WindowCosts  `bson:"costs" json:"costs"`
	Dims        *Dims         `bson:"dims" json:"dims"`
	Qty         int           `bson:"qty" json:"qty"`
	ProductID   bson.ObjectId `bson:"productID" json:"productID"`
	ProductName string
	Rooms       []string `bson:"rooms" json:"rooms"`
	Specs       bson.M   `bson:"specs" json:"specs"`
}

// WindowCosts struct
type WindowCosts struct {
	Window      float64 `bson:"window" json:"window"`
	Trim        float64 `bson:"trim" json:"trim"`
	Options     float64 `bson:"options" json:"options"`
	InstallType float64 `bson:"installType" json:"installType"`
	Install     float64 `bson:"install" json:"install"`
	Unit        float64 `bson:"extendUnit" json:"extendUnit"`
	Total       float64 `bson:"extendTotal" json:"extendTotal"`
}
