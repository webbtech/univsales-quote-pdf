package mongo

import (
	"os"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/pulpfree/univsales-pdf/config"
	"github.com/pulpfree/univsales-pdf/model"
	"github.com/stretchr/testify/suite"
)

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
	cfg *config.Config
	db  *DB
	q   *model.Quote
}

const (
	defaultsFP = "../../config/defaults.yml"
	// quoteID    = "5b2d7e1de730a4ea81b6ad34"
	// quoteID = "5c86bbfbd30c27fb399c71ae"
	quoteID = "5b19e0c62aac0409e37ec013" // quote with payment
)

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := suite.cfg.Load()
	suite.NoError(err)

	// setup db
	s, err := mgo.Dial(suite.cfg.GetMongoConnectURL())
	suite.NoError(err)

	suite.db = &DB{
		session: s,
		dbName:  suite.cfg.DBName,
	}
	suite.q = &model.Quote{}
}

// TestNewDB method
func (suite *IntegSuite) TestNewDB() {
	_, err := NewDB(suite.cfg.GetMongoConnectURL(), suite.cfg.DBName)
	suite.NoError(err)
}

func (suite *IntegSuite) TestGetQuote() {
	err := suite.db.getQuote(suite.q, quoteID)
	suite.NoError(err)
	suite.NotNil(suite.q.Number)
}

func (suite *IntegSuite) TestGetPayments() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getPayments(suite.q)
	suite.NoError(err)
	suite.True(len(suite.q.Payments) > 0)
}

func (suite *IntegSuite) TestGetGroupItems() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getGroupItems(suite.q)
	suite.NoError(err)
	suite.True(len(suite.q.Items.Group) > 0)
}

func (suite *IntegSuite) TestGetWindowItems() {
	_ = suite.db.getQuote(suite.q, quoteID)
	err := suite.db.getWindowItems(suite.q)
	suite.NoError(err)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
