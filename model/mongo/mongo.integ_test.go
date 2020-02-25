package mongo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pulpfree/univsales-quote-pdf/config"
	"github.com/pulpfree/univsales-quote-pdf/model"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
	cfg *config.Config
	db  *MDB
	q   *model.Quote
}

const (
	defaultsFP = "../../config/defaults.yml"
	// quoteID    = "5cd09cb3bc1f5e721f382e63"
	quoteID = "5b8169aacf2585816cec1e64"
)

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := suite.cfg.Load()
	suite.NoError(err)
	fmt.Printf("suite.cfg.GetMongoConnectURL() %s\n", suite.cfg.GetMongoConnectURL())

	// Set client options
	clientOptions := options.Client().ApplyURI(suite.cfg.GetMongoConnectURL())

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	suite.NoError(err)

	suite.db = &MDB{
		client: client,
		dbName: suite.cfg.DBName,
		db:     client.Database(suite.cfg.DBName),
	}
	suite.q = &model.Quote{}
	suite.q.Items = &model.Items{}
}

// TestNewDB method
func (suite *IntegSuite) TestNewDB() {
	_, err := NewDB(suite.cfg.GetMongoConnectURL(), suite.cfg.DBName)
	suite.NoError(err)
}

func (suite *IntegSuite) TestGetQuote() {
	err := suite.db.getQuote(suite.q, quoteID)
	suite.NoError(err)
	suite.True(suite.q.Number > 0)
}

func (suite *IntegSuite) TestGetPayments() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getPayments(suite.q)
	suite.NoError(err)
	if quoteID == "5b19e0c62aac0409e37ec013" {
		suite.True(len(suite.q.Payments) > 0)
	}
}

func (suite *IntegSuite) TestGetGroupItems() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getGroupItems(suite.q)
	suite.NoError(err)
	// fmt.Printf("suite.q.Items.Group %+v\n", suite.q.Items.Group[0])
	// suite.True(len(suite.q.Items.Group) > 0)
}

func (suite *IntegSuite) TestGetWindowItems() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getWindowItems(suite.q)
	suite.NoError(err)
	// suite.True(len(suite.q.Items.Window) > 0)
}

func (suite *IntegSuite) TestGetOtherItems() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getOtherItems(suite.q)
	suite.NoError(err)
	// suite.True(len(suite.q.Items.Other) > 0)
}

func (suite *IntegSuite) TestGetCustomer() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getCustomer(suite.q)
	suite.NoError(err)
	suite.True(suite.q.Customer.Name.First != "")
	suite.True(len(suite.q.Customer.PhoneMap) > 0)
	suite.True(suite.q.Customer.Address.Associate == "customer")
}

/* func (suite *IntegSuite) TestGetJobsheetFeatures() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getJobsheetFeatures(suite.q)
	suite.NoError(err)
	suite.True(suite.q.Features != "")
} */

func (suite *IntegSuite) TestFetchQuote() {
	q, err := suite.db.FetchQuote(quoteID)
	fmt.Printf("q %+v\n", q.ID)
	suite.NoError(err)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
