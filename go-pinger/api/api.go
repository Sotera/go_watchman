package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	mgo "gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"

	"path/filepath"
	"time"

	"strings"

	"gopkg.in/gin-gonic/gin.v1"
)

// Stats are various metrics
type Stats struct {
	QueryType string  `form:"query_type" json:"query_type"`
	RespTime  int     `form:"resp_time" json:"resp_time"`
	Count     float64 `form:"count" json:"count"`
	Created   int     `form:"created" json:"created"`
}

var db *mgo.Session

// match cache-busting names with local filenames
func cacheBustHandler(c *gin.Context) {
	// match route like /app/(js|css)/script.123456.(js|css)
	compiled, err := regexp.Compile(`(?i)([a-z0-9_-]+)\.\d+\.([a-z]+)`)
	if err != nil {
		fmt.Println(err)
		return // ok to return here?
	}
	param := c.Param("filename")
	if match := compiled.FindStringSubmatch(param); match != nil {
		// rm cache bust from param
		f := strings.Join([]string{match[1], match[2]}, ".")
		// create relative path
		p := filepath.Join("../app", match[2], f)
		fp, err := filepath.Abs(p)
		if err != nil {
			panic(err)
		}
		c.File(fp)
	} else {
		// not sure if this is correct way to continue route matching
		c.Abort()
	}
}

func main() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "mongo:27017"
	}

	db, err = mgo.Dial(dbHost)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Optional. Switch the session to a monotonic behavior.
	// dbSession.SetMode(mgo.Monotonic, true)

	r := gin.Default()

	indexTmpl, err := filepath.Abs("../app/index.tmpl")
	if err != nil {
		panic(err)
	}

	r.LoadHTMLFiles(indexTmpl)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dash")
	})

	r.GET("/dash", func(c *gin.Context) {
		var ts int64
		if m := os.Getenv("GIN_MODE"); m == "release" {
			ts = time.Now().Unix()
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"cache_bust": ts,
		})
	})

	r.GET("/app/js/:filename", cacheBustHandler)
	r.GET("/app/css/:filename", cacheBustHandler)

	// non-vendored assets go thru cachebust handler
	r.Static("/app/vendor", "../app/vendor")
	r.StaticFile("/favicon.ico", "../app/favicon.ico")

	test := r.Group("/test")
	{
		test.GET("/ping/:v", func(c *gin.Context) {
			v := c.Param("v")
			c.String(http.StatusOK, "%s", v)
		})
	}

	api := r.Group("/api")
	{
		api.POST("/watchman", func(c *gin.Context) {
			var json Stats
			if err := c.BindJSON(&json); err == nil {
				fmt.Printf("%+v\n", json)

				saveToDb(json)
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			} else {
				fmt.Println(err)
			}
		})

		api.GET("/watchman", func(c *gin.Context) {
			results := []Stats{}

			coll := db.DB("app_stats").C("app")
			err := coll.Find(nil).Limit(60000).Sort("-created").All(&results)
			if err != nil {
				log.Fatal(err)
			}

			c.JSON(http.StatusOK, &results)
		})

		// allow get request for ease of use
		api.GET("/watchman/drop", func(c *gin.Context) {
			coll := db.DB("app_stats").C("app")
			err := coll.DropCollection()
			if err != nil {
				fmt.Println(err)
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func saveToDb(stats Stats) {
	coll := db.DB("app_stats").C("app")
	err := coll.Insert(stats)
	if err != nil {
		log.Fatal(err)
	}
}
