package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	//"net/url"

	"github.com/gin-gonic/gin"
)

var DB *sql.DB

type Link struct {
	LINK       string    `json:"link"`
	TIME_STAMP time.Time `json:"time_stamp"`
}

func main() {
	createDBConnection()
	defer DB.Close()
	r := gin.Default()
	r.Use(CORSMiddleware())
	setupRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
func setupRoutes(r *gin.Engine) {

	r.POST("/user_link", SaveLongLink)
	r.GET("/user_link/all", GetAllLinks)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// POST
func SaveLongLink(c *gin.Context) {
	reqBody := Link{}
	err := c.Bind(&reqBody)
	if err != nil {
		res := gin.H{
			"error": "invalid request body",
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, res)

		return
	}

	//reqBody.ValidUrl = validurl(reqBody.URL)

	// Data[lastID] = reqBody
	reqBody.TIME_STAMP = time.Now()
	fmt.Println(reqBody)
	res, err := DB.Exec(`INSERT INTO "user_link" ( "link","time_stamp")
	VALUES ( $1, $2)`, reqBody.LINK, reqBody.TIME_STAMP)
	if err != nil {
		fmt.Println("err inserting data: ", err)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	lastInsID, _ := res.LastInsertId()

	fmt.Println("res: ", lastInsID)
	c.JSON(http.StatusOK, reqBody)
	c.Writer.Header().Set("Content-Type", "application/json")
}

// GET

func GetAllLinks(c *gin.Context) {
	rows, err := DB.Query("SELECT link, time_stamp, FROM user_link order by id desc")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	links := []Link{}
	for rows.Next() {
		l := Link{}
		err := rows.Scan(&l.LINK, &l.TIME_STAMP)
		if err != nil {
			panic(err)
		}
		links = append(links, l)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	res := gin.H{
		"data": links,
	}
	c.JSON(http.StatusOK, res)
}
