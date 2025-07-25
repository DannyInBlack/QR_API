package main

import (
	"fmt"
	// "log"
	"net/http"
	"os"
	// "os/exec"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"  
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

var running = false
var statusText = widget.NewLabel("Status: Stopped")
var generatedQRImage = widget.NewCard("", "Waiting for QR Codes", nil)
var enterText = widget.NewEntry()

func serve() {
	router := gin.Default()
	router.GET("/qrcode", genQRcode)
	// router.POST("/qrcode", genQRcodes)
	router.Run("localhost:8080")
}

func genQRcode(c *gin.Context) {

	if !running {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server is stopped"})
		return
	}

	link := c.Query("link")

	if link == "" { // link argument is missing
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid link argument"})
		return
	}

	// Try making the links directory if it doesn't exist
	if err := os.Mkdir("links", 0750); err != nil && !os.IsExist(err) {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server couldn't access links directory"})
		return
	}

	filepath := "links/" + link + ".png"

	println("New link received: " + filepath)

	if _, err := os.ReadFile(filepath); err != nil { // File doesn't exist - generate new code
		err := qrcode.WriteFile(link, qrcode.Medium, 256, filepath)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server couldn't generate QRcode" + err.Error()})
			return
		}
	} else { // File exists
		println("Link already found: " + filepath)
	}

	img := canvas.NewImageFromFile(filepath)
	img.FillMode = canvas.ImageFillOriginal
	generatedQRImage.SetImage(img)
	generatedQRImage.SetSubTitle(link)

	c.File(filepath)
}

// func generateZIP(images []string) error {

// 	if err := exec.Command("mkdir", "temp").Run(); err != nil {
// 		// log.Fatal(err)
// 		return err
// 	}

// 	args := append(images, "temp")

// 	_, err := exec.Command("mv", args...).Output()
// 	if err != nil {
// 		return err
// 	}

// 	_, err = exec.Command("zip", "-r", "temp.zip", "temp").Output()
// 	if err != nil {
// 		return err
// 	}

// 	return exec.Command("rm", "-R", "temp").Run()
// }

// func genQRcodes(c *gin.Context) {

// 	if err := c.Request.ParseMultipartForm(5000); err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	links := c.Request.Form["links"]
// 	for i, link := range links {
// 		link += ".png"
// 		err := qrcode.WriteFile(link, qrcode.Medium, 256, link)
// 		if err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server couldn't generate QRcode"})
// 			return
// 		}

// 		links[i] = link
// 	}

// 	err := generateZIP(links)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	c.File("temp.zip")
// }

func setStatusOk() {
	running = true
	statusText.SetText("Status: Running")
}

func setStatusNotOk() {
	running = false
	statusText.SetText("Status: Stopped")
}

func main() {
	a := app.New()
	w := a.NewWindow("QR Code Generator")

	startButton := widget.NewButton("Start Server", setStatusOk)
	stopButton := widget.NewButton("Stop Server", setStatusNotOk)

	buttons := container.New(layout.NewVBoxLayout(), startButton, stopButton)
	status := container.New(layout.NewVBoxLayout(), statusText, generatedQRImage)
	startPart := container.New(layout.NewHBoxLayout(), status, layout.NewSpacer(), buttons)
	go serve()

	w.SetContent(startPart)
	w.ShowAndRun()

	fmt.Println("Exited")
}
