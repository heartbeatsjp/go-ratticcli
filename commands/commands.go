package commands

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/urfave/cli.v1"
)

var (
	// SearchFlags are flags of sub command Search
	SearchFlags []cli.Flag
	// ShowFlags are flags of sub command Show
	ShowFlags []cli.Flag
	// ReloadFlags are flags of sub command Reload
	ReloadFlags []cli.Flag
)

func init() {
	SearchFlags = []cli.Flag{}
	ShowFlags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "fileds",
			Usage: "username,password",
		}}
	ReloadFlags = []cli.Flag{}
}

/*
SearchAction do HTTP request to search Cred url
*/
func SearchAction(c *cli.Context) error {
	cachePath := c.String("cache-path")
	if CacheExpired(cachePath) {
		err := ReloadAction(c)
		if err != nil {
			return err
		}
	}

	// print cached Creds from Creds Bucket
	creds := GetCachedCreds(cachePath)
	for _, cred := range creds {
		fmt.Println(fmt.Sprintf("%s %s", cred[0], cred[1]))
	}

	return nil
}

/*
ShowAction do HTTP request to fetch Cred details
*/
func ShowAction(c *cli.Context) error {
	return nil
}

/*
ReloadAction do re-auth, update token, discard local cache
*/
func ReloadAction(c *cli.Context) error {
	var err error
	cachePath := c.String("cache-path")
	// Authrize and refresh token
	// store Token to Config Bucket

	//TODO

	// do HTTP search request
	// store search results to Creds Bucket

	//TODO
	token := GetCachedToken(cachePath)
	credLines := GetWebCreds(c.String("endpoint"), c.String("user"), token)
	//store creds
	err = StoreCreds(cachePath, credLines)
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
		log.Fatalln(err)
	}
	defer db.Close()

	// compare
	var lastUpdated time.Time
	err = db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("Config"))
		v := b.Get([]byte("LastUpdated"))
		lastUpdated, err = time.Parse(time.RFC1123Z, string(v))
		if err != nil {
			log.Fatalln(err)
		}

		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return lastUpdated.Add(86400 * time.Second).After(time.Now())
}

/*
GetCreds return creds ( slice of [][]byte{key, value} ) from cache
*/
func GetCachedCreds(cachePath string) [][][]byte {
	db, err := bolt.Open(cachePath, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var creds [][][]byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Creds"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			creds = append(creds, [][]byte{k, v})
		}

		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return creds
}

/*
GetCachedToken
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
GetWebCreds return `Cred.id Cred.title` line from RatticWeb
*/
func GetWebCreds(cachePath string) []string {
	var credLines []string

	//TODO

	return credLines
}

/*
BuildHTTPRequest
*/
func BuildHTTPRequest(user, token, endpoint string) string {
	//TODO
	client := &http.Client{Timeout: time.Duration(10) * time.Second}
}

/*
StoreCreds
*/
func StoreCreds(cachePath string, lines []string) error {
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

		b := tx.CreateBucket([]byte("Creds"))
		for _, line := range lines {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				log.Printf("ERROR: Parse response line failed. %s\n", line)
				continue
			}
			err = b.Put([]byte(parts[0]), []byte(parts[1]))
			if err != nil {
				log.Printf("ERROR: Put to Bucket failed. %v\n", err)
			}

		}

		return nil
	})

	return err
}
