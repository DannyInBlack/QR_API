package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	godotenv "github.com/joho/godotenv"
	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := gin.Default()
	router.GET("/qrcode", genQRcode)
	router.POST("/qrcode", genQRcodes)
	router.Run("localhost:8080")

	uri := os.Getenv("MONGODB_URI")
	println(uri)
}

func genQRcode(c *gin.Context) {
	link := c.Query("link")
	if link == "" {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Link is invalid"})
		return
	}

	err := qrcode.WriteFile(link, qrcode.Medium, 256, link+".png")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server couldn't generate QRcode"})
		return
	}

	c.File(link + ".png")
}

func generateZIP(images []string) error {

	if err := exec.Command("mkdir", "temp").Run(); err != nil {
		// log.Fatal(err)
		return err
	}

	args := append(images, "temp")

	_, err := exec.Command("mv", args...).Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("zip", "-r", "temp.zip", "temp").Output()
	if err != nil {
		return err
	}

	return exec.Command("rm", "-R", "temp").Run()
}

func genQRcodes(c *gin.Context) {

	if err := c.Request.ParseMultipartForm(5000); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	links := c.Request.Form["links"]
	for i, link := range links {
		link += ".png"
		err := qrcode.WriteFile(link, qrcode.Medium, 256, link)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server couldn't generate QRcode"})
			return
		}

		links[i] = link
	}

	err := generateZIP(links)

	if err != nil {
		log.Fatal(err)
	}

	c.File("temp.zip")
}
