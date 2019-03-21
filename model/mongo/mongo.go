package mongo

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pulpfree/univsales-pdf/model"

	log "github.com/sirupsen/logrus"
)

// DB and Table constants
const (
	colAddress    = "addresses"
	colCustomer   = "customers"
	colGroupTypes = "group-types"
	colJS         = "jobsheets"
	colJSGroups   = "jobsheet-win-grps"
	colJSOther    = "jobsheet-other"
	colJSWindows  = "jobsheet-wins"
	colPayments   = "payments"
	colProducts   = "products"
	colQuotes     = "quotes"
	Users         = "users"
)

// DB struct
type DB struct {
	session *mgo.Session
	dbName  string
}

// NewDB sets up new MongoDB struct
func NewDB(connection string, dbNm string) (model.DBHandler, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}

	return &DB{
		session: s,
		dbName:  dbNm,
	}, err
}

// FetchQuote method
func (db *DB) FetchQuote(quoteID string) (*model.Quote, error) {

	// Initialize
	q := &model.Quote{}

	// Fetch quote
	err := db.getQuote(q, quoteID)
	if err != nil {
		return q, err
	}

	// Fetch payments
	/*err = db.getPayments(q)
	if err != nil {
		return q, err
	}

	// Fetch Group items
	err = db.getGroupItems(q)
	if err != nil {
		return q, err
	}

	// Fetch Windows
	err = db.getWindowItems(q)
	if err != nil {
		return q, err
	}

	// Fetch Other items
	err = db.getOtherItems(q)
	if err != nil {
		return q, err
	}

	// Fetch Jobsheet features
	err = db.getJobsheetFeatures(q)
	if err != nil {
		return q, err
	}

	// Fetch customer data
	err = db.getCustomer(q)
	if err != nil {
		return q, err
	}
	*/
	return q, nil
}

func (db *DB) getQuote(q *model.Quote, quoteID string) error {

	if quoteID == "" {
		return errors.New("Missing quoteID string")
	}
	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(db.dbName).C(colQuotes)
	err := col.FindId(bson.ObjectIdHex(quoteID)).One(&q)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getPayments(q *model.Quote) error {

	s := db.getFreshSession()
	defer s.Close()

	q.Payments = []*model.Payment{}
	col := s.DB(db.dbName).C(colPayments)
	err := col.Find(bson.M{"quoteID": q.ID}).All(&q.Payments)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getGroupItems(q *model.Quote) error {

	s := db.getFreshSession()
	defer s.Close()

	q.Items = &model.Items{}
	col := s.DB(db.dbName).C(colJSGroups)
	for _, i := range q.ItemIds.Group {
		item := &model.Group{}
		if false == bson.IsObjectIdHex(i) {
			log.WithFields(log.Fields{"groupID": i}).Fatal("Invalid groupID")
			return errors.New("Invalid groupID")
		}
		if err := col.FindId(bson.ObjectIdHex(i)).One(&item); err != nil {
			log.Fatalf("Group not found. Error: %s", err)
			return err
		}
		// Fetch group type
		colGT := s.DB(db.dbName).C(colGroupTypes)
		ret := bson.M{}
		if err := colGT.FindId(item.Specs["groupType"]).One(&ret); err != nil {
			log.Fatalf("Group Type not found. Error: %s", err)
			return err
		}
		// Set group type name
		item.Specs["groupTypeName"] = ret["name"]
		q.Items.Group = append(q.Items.Group, item)
	}

	return nil
}

/* func (db *DB) getWindowItems(q *model.Quote) error {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(db.dbName).C(colJobsheetWindows)
	for _, i := range q.ItemIds.Window {
		item := &model.Window{}
		if false == bson.IsObjectIdHex(i) {
			log.WithFields(log.Fields{"windowID": i}).Fatal("Invalid windowID")
			return errors.New("Invalid windowID")
		}
		if err := col.FindId(bson.ObjectIdHex(i)).One(&item); err != nil {
			log.Fatalf("Window not found. Error: %s", err)
			return err
		}
		// fetch product info
		prod, err := db.getProductName(item.ProductID)
		if err != nil {
			return err
		}
		item.ProductName = prod.Name
		q.Items.Window = append(q.Items.Window, item)
	}

	return nil
} */

/* func (db *DB) getOtherItems(q *model.Quote) error {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(db.dbName).C(colJobsheetOther)
	for _, i := range q.ItemIds.Other {
		item := &model.Other{}
		if false == bson.IsObjectIdHex(i) {
			log.WithFields(log.Fields{"windowID": i}).Fatal("Invalid windowID")
			return errors.New("Invalid windowID")
		}
		if err := col.FindId(bson.ObjectIdHex(i)).One(&item); err != nil {
			log.Fatalf("Other not found. Error: %s", err)
			return err
		}
		q.Items.Other = append(q.Items.Other, item)
	}

	return nil
} */

/* func (db *DB) getProductName(productID bson.ObjectId) (*model.Product, error) {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(db.dbName).C(colProducts)
	p := &model.Product{}
	if err := col.FindId(productID).One(&p); err != nil {
		return p, err
	}
	return p, nil
} */

/* func (db *DB) getCustomer(q *model.Quote) error {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(db.dbName).C(colCustomer)
	err := col.FindId(q.CustomerID).One(&q.Customer)
	if err != nil {
		return err
	}
	phoneMap := map[string]string{}
	for _, v := range q.Customer.Phones {
		phoneMap[v.Type] = v.Number
	}
	q.Customer.PhoneMap = phoneMap

	// Fetch customer address data
	col = s.DB(dbName).C(colAddress)
	err = col.Find(bson.M{"customerID": q.CustomerID, "associate": "customer"}).One(&q.Customer.Address)
	if err != nil {
		return err
	}

	return nil
} */

/* func (db *DB) getJobsheetFeatures(q *model.Quote) error {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(dbName).C(colJobsheet)
	jobsheet := map[string]string{}
	if err := col.FindId(q.JobsheetID).One(&jobsheet); err != nil {
		log.Fatalf("Jobsheet not found. Error: %s", err)
		return err
	}
	q.Features = jobsheet["features"]

	return nil
} */

// Close method
func (db *DB) Close() {
	db.session.Close()
}

// ============================== Helper methods ==============================

func (db *DB) getFreshSession() *mgo.Session {
	return db.session.Copy()
}
