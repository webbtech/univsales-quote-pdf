package pdf

import (
	"os"
	"regexp"
	"testing"

	"github.com/pulpfree/univsales-pdf/config"
	"github.com/pulpfree/univsales-pdf/model"
	"github.com/pulpfree/univsales-pdf/model/mongo"
	"github.com/stretchr/testify/suite"
)

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
	cfg *config.Config
	db  model.DBHandler
	q   *model.Quote
	p   *PDF
}

const (
	defaultsFP = "../config/defaults.yml"
	quoteID    = "5cd16f18699e0300c7b10d30"
	invoiceID  = "5cc7b5410765188d68d191ca"
)

// SetupTest method
func (suite *IntegSuite) SetupTest() {

	// req := &Request{
	// 	QuoteID: quoteID,
	// 	DocType: "quote",
	// }
	// setup config
	os.Setenv("Stage", "test")
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := suite.cfg.Load()
	suite.NoError(err)

	suite.db, err = mongo.NewDB(suite.cfg.GetMongoConnectURL(), suite.cfg.DBName)
	suite.NoError(err)

	// suite.q, err = suite.db.FetchQuote(quoteID)
	// suite.NoError(err)

	// suite.p = New(req, suite.q, suite.cfg)
	// suite.NoError(err)
}

func (suite *IntegSuite) TestTypes() {
	suite.True(suite.q.Number > 0)
	suite.IsType(&model.Quote{}, suite.q)
	suite.IsType(&PDF{}, suite.p)
}

func (suite *IntegSuite) TestSetFileNameOutput() {

	// start with quote
	req := &Request{
		QuoteID: quoteID,
		DocType: "quote",
	}
	suite.p = New(req, suite.q, suite.cfg)
	suite.p.setOutputFileName()
	r, _ := regexp.Compile("^quote\\/qte-([0-9]+)-r([0-9]+)\\.pdf?")
	suite.True(r.MatchString(suite.p.outputFileName))

	// Now with invoice
	req = &Request{
		QuoteID: quoteID,
		DocType: "invoice",
	}
	suite.p = New(req, suite.q, suite.cfg)
	suite.p.setOutputFileName()
	r, _ = regexp.Compile("^invoice\\/inv-([0-9]+)\\.pdf?")
	suite.True(r.MatchString(suite.p.outputFileName))
}

func (suite *IntegSuite) TestQuoteOutputToDisk() {

	req := &Request{
		QuoteID: quoteID,
		DocType: "quote",
	}
	suite.cfg.SetStageEnv("test")
	suite.p = New(req, suite.q, suite.cfg)

	err := suite.p.Quote()
	suite.NoError(err)

	err = suite.p.OutputToDisk()
	suite.NoError(err)
}

func (suite *IntegSuite) TestInvoiceOutputToDisk() {

	quoteID := invoiceID
	// fmt.Printf("quoteID %s\n", quoteID)
	suite.q, _ = suite.db.FetchQuote(quoteID)

	req := &Request{
		QuoteID: invoiceID,
		DocType: "invoice",
	}
	suite.cfg.SetStageEnv("test")
	suite.p = New(req, suite.q, suite.cfg)

	err := suite.p.Invoice()
	suite.NoError(err)

	err = suite.p.OutputToDisk()
	suite.NoError(err)
}

func (suite *IntegSuite) TestQuoteSaveToS3() {
	req := &Request{
		QuoteID: quoteID,
		DocType: "quote",
	}
	suite.cfg.SetStageEnv("test")
	suite.p = New(req, suite.q, suite.cfg)

	err := suite.p.Quote()
	suite.NoError(err)

	location, err := suite.p.SaveToS3()
	suite.NoError(err)
	suite.True(location != "")
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
