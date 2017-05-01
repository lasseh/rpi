package main

import (
	"net/http"
	"time"

	unicorn "github.com/arussellsaw/unicorn-go"
	"github.com/gin-gonic/gin"
)

//const RED = (255, 0, 0)
var (
	Red    = unicorn.Pixel{0, 255, 0}
	Green  = unicorn.Pixel{255, 0, 0}
	Orange = unicorn.Pixel{230, 150, 0}
	Blue   = unicorn.Pixel{0, 0, 255}
	Cyan   = unicorn.Pixel{0, 255, 255}
	Pink   = unicorn.Pixel{255, 0, 127}
)

func main() {
	// Unicorn
	u := unicorn.Client{Path: unicorn.SocketPath}
	u.Connect()
	u.SetBrightness(20)

	// Gin Mode
	gin.SetMode(gin.DebugMode)

	// Creates a new router
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/fap/:status", func(c *gin.Context) {
		status := c.Param("status")
		go blinkFap(&u, status)
		c.JSON(http.StatusOK, gin.H{"blink": "success"})
	})

	r.GET("/dhcp", func(c *gin.Context) {
		go blinkDhcp(&u)
		c.JSON(http.StatusOK, gin.H{"blink": "success"})
	})

	r.GET("/color/:color", func(c *gin.Context) {
		color := c.Param("color")
		go blinkColor(&u, color)
		c.JSON(http.StatusOK, gin.H{"blink": "success"})
	})

	r.GET("/", indexHandler)

	// Listen and serv
	r.Run(":80")

}

func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"blink.tg17.gathering.org": "blink the light at noc"})
}

func blinkColor(u *unicorn.Client, color string) {
	// Set full frame
	pixels := [64]unicorn.Pixel{}

	switch color {
	case "red":
		for i := range pixels {
			pixels[i] = unicorn.Pixel(Red)
		}
	case "green":
		for i := range pixels {
			pixels[i] = unicorn.Pixel(Green)
		}
	case "orange":
		for i := range pixels {
			pixels[i] = unicorn.Pixel(Orange)
		}
	case "blue":
		for i := range pixels {
			pixels[i] = unicorn.Pixel(Blue)
		}
	case "cyan":
		for i := range pixels {
			pixels[i] = unicorn.Pixel(Cyan)
		}
	case "pink":
		for i := range pixels {
			pixels[i] = unicorn.Pixel(Pink)
		}
	}
	u.SetAllPixels(pixels)
	u.Show()
	for i := 1; i < 100; i++ {
		u.SetBrightness(uint(i))
		u.Show()
		time.Sleep(1 * time.Millisecond)
	}

	// Down
	for i := 100; i > 1; i-- {
		u.SetBrightness(uint(i))
		u.Show()
		time.Sleep(1 * time.Millisecond)
	}
}

func blinkFap(u *unicorn.Client, status string) {
	sec := 10
	pixels := [64]unicorn.Pixel{}

	switch status {
	case "success":
		// Set leds to green
		for i := range pixels {
			pixels[i] = unicorn.Pixel(Green)
		}
	case "fail":
		// Set leds to red
		for i := range pixels {
			pixels[i] = unicorn.Pixel(Red)
		}
	}

	// Blink for $sec
	for i := 1; i <= sec; i++ {
		// Up
		u.SetAllPixels(pixels)
		u.Show()
		for i := 1; i < 100; i++ {
			u.SetBrightness(uint(i))
			u.Show()
			time.Sleep(1 * time.Millisecond)
		}

		// Down
		for i := 100; i > 1; i-- {
			u.SetBrightness(uint(i))
			u.Show()
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func blinkDhcp(u *unicorn.Client) {
	// Set full frame
	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = unicorn.Pixel{255, 0, 255}
	}
	// Up
	u.SetAllPixels(pixels)
	u.Show()
	for i := 1; i < 100; i++ {
		u.SetBrightness(uint(i))
		u.Show()
		time.Sleep(1 * time.Millisecond)
	}

	// Down
	for i := 100; i > 1; i-- {
		u.SetBrightness(uint(i))
		u.Show()
		time.Sleep(1 * time.Millisecond)
	}
}
