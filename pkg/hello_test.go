package hello

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pamrulla/gagster-feed/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HelloTestSuite struct {
	suite.Suite
	endpoint string
	ts       *httptest.Server
	router   *chi.Mux
}

func (hts *HelloTestSuite) SetupTest() {
	hts.router = helpers.CreateNewRouter()
	hts.router.Get("/api/hello", Hello)
	hts.ts = httptest.NewServer(hts.router)
}

func (hts *HelloTestSuite) TearDownTest() {
	hts.ts.Close()
}

func (hts *HelloTestSuite) TestStatusCodeShouldBeEqual200() {
	// Arrange
	hts.endpoint = "/api/hello/"

	// Act
	resp, _ := helpers.RunRequest("GET", hts.ts, hts.endpoint, nil)

	// Assert
	assert.Equal(hts.T(), 200, resp.StatusCode)
}
func TestHelloTestSuite(t *testing.T) {
	suite.Run(t, new(HelloTestSuite))
}
