package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Person struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Name      string        `json:"name"`
		Age       uint8         `json:"age"`
		Timestamp time.Time     `json:"-"`
	}
)

var (
	MongoConn *mgo.Session
)

const (
	databaseName = "mongoadmin"
	people       = "people"
)

func main() {
	e := echo.New()
	database()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello k8s")
	})

	e.GET("/users", func(c echo.Context) error {
		var peoples []*Person
		err := MongoConn.DB(databaseName).C(people).Find(bson.M{}).All(&peoples)
		if err != nil {
			return c.JSON(http.StatusFound, err)
		}
		return c.JSON(http.StatusOK, peoples)
	})

	e.POST("/users", func(c echo.Context) (err error) {
		var p *Person
		if err := c.Bind(&p); err != nil {
			return err
		}
		MongoConn.DB(databaseName).C(people).Insert(&Person{
			Name: p.Name, Age: p.Age, Timestamp: time.Now()},
		)
		return c.JSON(http.StatusOK, p)
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func database() *mgo.Session {
	// MongoDB
	MongoConfig := &mgo.DialInfo{
		//Addrs:     []string{"mongo-server:27017"},
		Addrs:     []string{"localhost:27017"},
		Timeout:   10 * time.Second,
		PoolLimit: 10,
		Username:  "mongoadmin",
		Password:  "mongoadmin",
		Database:  "mongoadmin",
	}

	session, err := mgo.DialWithInfo(MongoConfig)
	if err != nil {
		log.Fatalf("session -> %s", err)
		os.Exit(1)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	MongoConn = session
	return session
}
