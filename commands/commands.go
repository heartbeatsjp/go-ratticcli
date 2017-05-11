package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/urfave/cli.v1"
)

var (
	// ListFlags are flags of sub command List
	ListFlags []cli.Flag
	// ShowFlags are flags of sub command Show
	ShowFlags []cli.Flag
	// ReloadFlags are flags of sub command Reload
	ReloadFlags []cli.Flag
)

func init() {
	ListFlags = []cli.Flag{}
	ShowFlags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "fields",
			Usage: "username,password",
		},
		cli.StringFlag{
			Name:  "id",
			Value: "-",
			Usage: "Cred.id ( - is STDIN)",
		},
	}
	ReloadFlags = []cli.Flag{}
}

// ListResponse HTTP Response
type ListResponse struct {
	Meta    ListResponseMeta   `json:"meta"`
	Objects []ListResponseCred `json:"objects"`
}

// ListResponseMeta HTTP Response
type ListResponseMeta struct {
	Next   string `json:"next"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// ListResponseCred HTTP Response
type ListResponseCred struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// ShowResponseCred HTTP Response
type ShowResponseCred struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

/*
ListAction do HTTP request to list Cred url
*/
func ListAction(c *cli.Context) error {
	cachePath := c.GlobalString("cache-path")
	if CacheExpired(cachePath) {
		err := ReloadAction(c)
		if err != nil {
			return err
		}
	}

	// print cached Creds from Creds Bucket
	creds := GetCachedCreds(cachePath)
	for _, cred := range creds {
		//fmt.Println(fmt.Sprintf("%s %s", cred[0], cred[1]))
		fmt.Println(cred)
	}

	return nil
}

/*
ShowAction do HTTP request to fetch Cred details
*/
func ShowAction(c *cli.Context) error {
	cachePath := c.GlobalString("cache-path")
	if CacheExpired(cachePath) {
		err := ReloadAction(c)
		if err != nil {
			return err
		}
	}

	fields := c.StringSlice("fields")

	// do HTTP list request
	token := c.GlobalString("token")
	if token == "" {
		token = GetCachedToken(cachePath)
	}

	var idString string
	if c.String("id") == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		idString = strings.SplitN(string(b), " ", 2)[0]
	} else {
		idString = c.String("id")
	}
	id, _ := strconv.Atoi(idString)

	cred := GetCred(c.GlobalString("endpoint"), c.GlobalString("user"), token, id)

	for _, field := range fields {
		if strings.ToUpper(field) == "ID" {
			fmt.Println(cred.ID)
		}
		if strings.ToUpper(field) == "PASSWORD" {
			fmt.Println(cred.Password)
		}
	}
	return nil
}

/*
ReloadAction do re-auth, update token, discard local cache
*/
func ReloadAction(c *cli.Context) error {
	var err error
	cachePath := c.GlobalString("cache-path")

	//TODO
	// Authrize and refresh token
	// store Token to Config Bucket

	// do HTTP list request
	token := c.GlobalString("token")
	if token == "" {
		token = GetCachedToken(cachePath)
	}

	now := time.Now()
	creds := GetCreds(c.GlobalString("endpoint"), c.GlobalString("user"), token)

	//store creds
	err = StoreCreds(cachePath, creds, now)
	if err != nil {
		return err
	}
	return nil
}

/*
CacheExpired return cache is expired or not
*/
func CacheExpired(cachePath string) bool {
	// get LastUpdated from Config Bucket
	db, err := bolt.Open(cachePath, 0600, nil)
	if err != nil {
		log.Println(err, "cache-path:", cachePath)
	}
	defer db.Close()

	// compare
	var lastUpdated time.Time
	err = db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("Config"))
		if b == nil {
			return errors.New("Bucket (and cache) does not exist.")
		}
		v := b.Get([]byte("LastUpdated"))
		lastUpdated, err = time.Parse(time.RFC1123Z, string(v))
		if err != nil {
			log.Fatalln(err)
		}

		return nil
	})
	if err != nil {
		return true
	}
	log.Println("LastUpdated:", lastUpdated, " and expired at", lastUpdated.Add(86400*time.Second))
	return lastUpdated.Add(86400 * time.Second).Before(time.Now())
}

/*
GetCachedCreds return creds from cache
*/
func GetCachedCreds(cachePath string) []string {
	db, err := bolt.Open(cachePath, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var creds []string
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Creds"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			creds = append(creds, fmt.Sprintf("%s %s", string(k), string(v)))
		}

		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return creds
}

/*
GetCachedToken return token
*/
func GetCachedToken(cachePath string) string {
	db, err := bolt.Open(cachePath, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var token string
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Config"))
		v := b.Get([]byte("Token"))
		token = string(v)

		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	return token
}

/*
GetCreds return ... from RatticWeb
*/
func GetCreds(endpoint, user, token string) []ListResponseCred {
	var creds []ListResponseCred

	var limit int
	var offset int
	hasNext := true
	for hasNext {
		req, err := BuildHTTPListRequest(endpoint, user, token, limit, offset)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := DoHTTPRequest(req)
		if err != nil {
			log.Fatalln(err)
		}

		// parse
		var listResponse ListResponse
		err = json.Unmarshal(body, &listResponse)
		if err != nil {
			log.Fatalln(err)
		}
		creds = append(creds, listResponse.Objects...)

		if listResponse.Meta.Next == "" {
			hasNext = false
		} else {
			limit = listResponse.Meta.Limit
			offset = listResponse.Meta.Offset + listResponse.Meta.Limit
		}
	}

	return creds
}

/*
GetCred return ... from RatticWeb
*/
func GetCred(endpoint, user, token string, id int) ShowResponseCred {
	var cred ShowResponseCred

	req, err := BuildHTTPShowRequest(endpoint, user, token, id)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := DoHTTPRequest(req)
	if err != nil {
		log.Fatalln(err)
	}

	// parse
	err = json.Unmarshal(body, &cred)
	if err != nil {
		log.Fatalln(err)
	}

	return cred
}

/*
BuildHTTPRequest builds HTTP Request
*/
func BuildHTTPRequest(endpoint, user, token, path string, queryParams map[string]string) (*http.Request, error) {

	ratticWebURL := fmt.Sprintf("%s/api/v1/%s", strings.Trim(endpoint, "/"), path)
	req, err := http.NewRequest("GET", ratticWebURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("ApiKey %s:%s", user, token))
	req.Header.Add("Accept", "application/json")

	values := url.Values{}
	for k, v := range queryParams {
		values.Add(k, v)
	}

	req.URL.RawQuery = values.Encode()

	return req, nil
}

/*
BuildHTTPListRequest builds HTTP Request to list Creds
*/
func BuildHTTPListRequest(endpoint, user, token string, limit, offset int) (*http.Request, error) {

	queryParams := make(map[string]string)

	if limit < 0 {
		// default
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	queryParams["limit"] = strconv.Itoa(limit)
	queryParams["offset"] = strconv.Itoa(offset)
	return BuildHTTPRequest(endpoint, user, token, "cred/", queryParams)
}

/*
BuildHTTPShowRequest builds HTTP Request to list Creds
*/
func BuildHTTPShowRequest(endpoint, user, token string, id int) (*http.Request, error) {

	return BuildHTTPRequest(endpoint, user, token, fmt.Sprintf("cred/%d/", id), make(map[string]string))
}

/*
DoHTTPRequest do HTTP Request and return response body
*/
func DoHTTPRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	log.Println(req, res.Status) //FIXME debug
	if res.StatusCode != 200 {
		return []byte{}, errors.New(res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}

/*
StoreCreds store creds to cache
*/
func StoreCreds(cachePath string, creds []ListResponseCred, lastUpdated time.Time) error {
	var err error

	db, err := bolt.Open(cachePath, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		if tx.Bucket([]byte("Creds")) != nil {
			tx.DeleteBucket([]byte("Creds"))
		}

		b, err := tx.CreateBucket([]byte("Creds"))
		if err != nil {
			log.Fatalln(err)
		}
		for _, cred := range creds {
			err = b.Put([]byte(strconv.Itoa(cred.ID)), []byte(cred.Title))
			if err != nil {
				log.Printf("ERROR: Put to Bucket failed. %v\n", err)
			}
		}

		b, err = tx.CreateBucketIfNotExists([]byte("Config"))
		b.Put([]byte("LastUpdated"), []byte(lastUpdated.Format(time.RFC1123Z)))

		return nil
	})

	return err
}
