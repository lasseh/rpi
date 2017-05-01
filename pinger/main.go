package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arussellsaw/unicorn-go"
	//"github.com/arussellsaw/unicorn-go/util"
	"math"
	"math/rand"

	"github.com/tatsushid/go-fastping"
)

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

const Brightness = 30

var Green = [3]uint{255, 0, 0}

func main() {
	c := unicorn.Client{Path: unicorn.SocketPath}
	c.Connect()

	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = unicorn.Pixel{255, 0, 0}
	}

	c.SetBrightness(40)
	//c.SetAllPixels(pixels)
	//Startup(&c)
	hostname := "google.com"

	p := fastping.NewPinger()

	netProto := "ip4:icmp"
	ra, err := net.ResolveIPAddr(netProto, hostname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	results := make(map[string]*response)
	results[ra.String()] = nil
	p.AddIPAddr(ra)

	onRecv, onIdle := make(chan *response), make(chan bool)
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}
	p.OnIdle = func() {
		onIdle <- true
	}

	p.MaxRTT = time.Second
	p.RunLoop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
loop:
	for {
		select {
		case <-ch:
			fmt.Println("get interrupted")
			break loop
			// Clear shield
			c.Clear()
			c.Show()
		case res := <-onRecv:
			if _, ok := results[res.addr.String()]; ok {
				results[res.addr.String()] = res
			}
		case <-onIdle:
			for host, r := range results {
				if r == nil {
					//fmt.Printf("%s : unreachable\n", host)
					Warning(&c)
					// Alarm
				} else {
					fmt.Printf("%s : %v\n", host, r.rtt)
					// Green
					// Check rtt
					ms := int(r.rtt.Seconds() * 1000)
					//fmt.Println("MS:", ms)

					// RPI Led
					switch {
					case ms <= 5:
						//fmt.Println("green")
						NoAlarms(&c)
					case ms > 5, ms > 10:
						//fmt.Println("amber")
						Amber(&c)
					case ms > 10:
						//fmt.Println("lol down")
						Warning(&c)
					}
				}
				results[host] = nil
			}
		case <-p.Done():
			if err = p.Err(); err != nil {
				fmt.Println("Ping failed:", err)
				// Alarm
			}
			break loop
			c.Clear()
			c.Show()
		}
	}
	signal.Stop(ch)
	p.Stop()
	// Clear shield
	c.Clear()
	c.Show()

	//Random(&c, 4)
	//Rainbow(&c)

	//WarningBlink(&c, 1)
	//NoAlarms(&c)
	//time.Sleep(10 * time.Second)

	//AlarmBlink(&c, 3)

	//Warning(&c)
	//time.Sleep(10 * time.Second)

	// Reset
	c.Clear()
	c.Show()
}
func Random(c *unicorn.Client, sec int) {
	for i := 1; i <= (sec * 10); i++ {
		for i := 1; i <= 30; i++ {
			c.SetPixel(uint(rand.Intn(8)), uint(rand.Intn(8)), uint(rand.Intn(255)), uint(rand.Intn(255)), uint(rand.Intn(255)))
		}
		c.Show()
		time.Sleep(100 * time.Millisecond)
		c.Clear()
		c.Show()
	}

}
func Rainbow(c *unicorn.Client) {
	i := 0.0
	offset := 3.0
	for {
		i += 0.3
		for y := 1; y <= 8; y++ {
			for x := 1; x <= 8; x++ {
				y := float64(y)
				x := float64(x)

				r := 0.0
				g := 0.0
				b := 0.0

				r = (math.Cos((x+i)/2.0) + math.Cos((y+i)/2.0)*64.0 + 128.0)
				g = (math.Sin((x+i)/1.5) + math.Sin((y+i)/2.0)*64.0 + 128.0)
				b = (math.Sin((x+i)/2.0) + math.Cos((y+i)/1.5)*64.0 + 128.0)
				r = math.Max(10, math.Min(255, r+offset))
				g = math.Max(10, math.Min(255, g+offset))
				b = math.Max(10, math.Min(255, b+offset))
				c.SetPixel(uint(x), uint(y), uint(r), uint(g), uint(b))
			}
		}
		c.Show()
		time.Sleep(10 * time.Millisecond)
	}
}

func Warning(c *unicorn.Client) {
	// Set full frame
	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = unicorn.Pixel{0, 255, 0}
	}
	// Set color
	c.SetAllPixels(pixels)
	c.SetBrightness(uint(Brightness))
	c.Show()
}

func Amber(c *unicorn.Client) {
	// Set full frame
	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = unicorn.Pixel{0, 0, 255}
	}
	// Set color
	c.SetAllPixels(pixels)
	c.SetBrightness(uint(Brightness))
	c.Show()
}

func NoAlarms(c *unicorn.Client) {
	// Set full frame
	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = unicorn.Pixel{255, 0, 0}
	}
	// Set color
	c.SetAllPixels(pixels)
	c.SetBrightness(uint(Brightness))
	c.Show()
}

func WarningBlink(c *unicorn.Client, sec int) {
	// Set full frame
	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = unicorn.Pixel{255, 0, 0}
	}
	// Blink for $sec
	for i := 1; i <= sec; i++ {
		// Up
		c.SetAllPixels(pixels)
		c.Show()
		for i := 1; i < 100; i++ {
			c.SetBrightness(uint(i))
			c.Show()
			time.Sleep(1 * time.Millisecond)
		}

		// Down
		for i := 100; i > 1; i-- {
			c.SetBrightness(uint(i))
			c.Show()
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func AlarmBlink(c *unicorn.Client, sec int) {
	// Set full frame
	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = unicorn.Pixel{0, 255, 0}
	}
	// Blink for $sec
	for i := 1; i <= sec; i++ {
		// Up
		c.SetAllPixels(pixels)
		c.Show()
		for i := 1; i < 100; i++ {
			c.SetBrightness(uint(i))
			c.Show()
			time.Sleep(1 * time.Millisecond)
		}

		// Down
		for i := 100; i > 1; i-- {
			c.SetBrightness(uint(i))
			c.Show()
			time.Sleep(1 * time.Millisecond)
		}
	}
}
func Startup(c *unicorn.Client) {
	// Set full frame
	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = unicorn.Pixel{255, 0, 255}
	}
	c.SetAllPixels(pixels)
	for i := 0; i <= 3; i++ {
		// Red
		for i := 20; i < 100; i++ {
			c.SetBrightness(uint(i))
			c.Show()
			time.Sleep(10 * time.Millisecond)
		}
		for i := 100; i > 20; i-- {
			c.SetBrightness(uint(i))
			c.Show()
			time.Sleep(10 * time.Millisecond)
		}

	}
}
