package main

//import "fmt"
import (
	"datx/ListenServer/chainlib"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getTransfer(c *gin.Context) {
	c.HTML(http.StatusOK, "transfer/index.html", gin.H{
		"title": "Main website",
	})
}

func postTransfer(c *gin.Context) {
	var trans chainlib.TransferInfo
	err := c.Bind(&trans) // c.BindJSON(&form)
	if err != nil {
		c.JSON(404, gin.H{"JSON=== status": "binding JSON error!"})
		return
	}

	quantity, err := strconv.ParseFloat(trans.Quantity, 64)
	trans.Quantity = strconv.FormatFloat(quantity, 'f', 4, 64) + " D" + "EOS"

	_, err = chainlib.ClWalletUnlock("PW5JHPpaGrS7bKhmQJ5Rb7rNSXhp3S3sXN2fGWaqQNzQufQaWrkUJ")
	transID, err := chainlib.ClPushTransfer("user", "transfer", trans)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("push transfer succeeded\t", transID)

	c.JSON(http.StatusOK, gin.H{
		"from":   trans.From,
		"to":     trans.To,
		"amount": trans.Quantity,
	})
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/**/*")

	router.GET("/transfer", getTransfer)

	router.POST("/transfer", postTransfer)

	router.GET("/user/:type/:name", getRouteStr)
	router.GET("/user", getQueryStrs)
	router.POST("/", getFormStr)

	router.Run(":8081")
}

func postHome(c *gin.Context) {
	uName := c.PostForm("name")
	c.JSON(200, gin.H{
		"say": "Hello " + uName,
	})
}

func getRouteStr(c *gin.Context) {
	ctype := c.Param("type")
	cname := c.Param("name")
	c.JSON(200, gin.H{
		"typeName": ctype,
		"username": cname,
	})
}

func getQueryStrs(c *gin.Context) {
	name := c.Query("name")           //如果没有相应值，默认为空字符串
	age := c.DefaultQuery("age", "0") //可设置默认值,string类型
	c.JSON(200, gin.H{
		"name": name,
		"age":  age,
	})
}

func getFormStr(c *gin.Context) {
	title := c.PostForm("title")
	cont := c.DefaultPostForm("cont", "没有内容")
	c.JSON(200, gin.H{
		"title": title,
		"cont":  cont,
	})
}
