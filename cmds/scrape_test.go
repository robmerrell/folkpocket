package cmds

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/robmerrell/folkpocket/config"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

//--------------------------------------------
// Test for scraping urls
//--------------------------------------------
type urlSuite struct {
	server1   *httptest.Server
	server2   *httptest.Server
	server3   *httptest.Server
	server4   *httptest.Server
	badServer *httptest.Server
}

var _ = Suite(&urlSuite{})

// create test servers
func (s *urlSuite) SetUpSuite(c *C) {
	s.badServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<tml><y><div>empty</div></bdy></hml>")
	}))

	s.server1 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body><div>empty</div></body></html>")
	}))

	s.server2 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		doc := "<html><body><a href='%s'><img src='images/rightarrow.gif'></a></body></html>"
		formatted := fmt.Sprintf(doc, s.server1.URL)
		fmt.Fprintln(w, formatted)
	}))

	s.server3 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		doc := "<html><body><a href='%s'><img src='images/rightarrow.gif'></a></body></html>"
		formatted := fmt.Sprintf(doc, s.server2.URL)
		fmt.Fprintln(w, formatted)
	}))

	s.server4 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		doc := "<html><body><a href='%s'><img src='images/rightarrow.gif'></a></body></html>"
		formatted := fmt.Sprintf(doc, s.server3.URL)
		fmt.Fprintln(w, formatted)
	}))
}

// close the test servers
func (s *urlSuite) TearDownSuite(c *C) {
	s.server1.Close()
	s.server2.Close()
	s.server3.Close()
	s.server4.Close()
	s.badServer.Close()
}

// test for extracting the next url when it exists
func (s *urlSuite) TestDocumentWithUrl(c *C) {
	urls, _ := getNextUrl(s.server4.URL, []string{})

	c.Check(urls, HasLen, 3)

	// check that all 3 urls were found
	c.Check(urls[0], Equals, s.server3.URL)
	c.Check(urls[1], Equals, s.server2.URL)
	c.Check(urls[2], Equals, s.server1.URL)
}

// test for returning gracefully if there is no next url
func (s *urlSuite) TestDocumentWithoutUrl(c *C) {
	urls, _ := getNextUrl(s.server1.URL, []string{})
	c.Check(urls, HasLen, 0)
}

// test returning an error
func (s *urlSuite) TestBadDocument(c *C) {
	s.badServer.Close()
	_, err := getNextUrl(s.badServer.URL, []string{})
	c.Check(err, Not(Equals), nil)
}

//--------------------------------------------
// Test saving a list of urls to Redis
//--------------------------------------------
type saveUrlSuite struct{}

var _ = Suite(&saveUrlSuite{})

func (s *saveUrlSuite) SetUpSuite(c *C) {
	os.Setenv("FOLKPOCKET_ENV", "test")
	config.LoadConfigFile("../config.toml")
}

func (s *saveUrlSuite) TestSavingAList(c *C) {
	key := config.Env().Get("cacheKey").(string)
	saveUrls([]string{"one", "two", "three"}, key)

	redisC, _ := connectToRedis(config.Env().Get("redishost").(string))
	defer redisC.Close()

	urls, _ := redis.Strings(redisC.Do("LRANGE", key, 0, -1))
	c.Check(urls[0], Equals, "one")
	c.Check(urls[1], Equals, "two")
	c.Check(urls[2], Equals, "three")

	// check that they values are overwritten
	saveUrls([]string{"more", "values"}, key)
	urls, _ = redis.Strings(redisC.Do("LRANGE", key, 0, -1))
	c.Check(urls[0], Equals, "more")
	c.Check(urls[1], Equals, "values")
}
