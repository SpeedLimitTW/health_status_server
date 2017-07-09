package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"

	"fmt"

	"bufio"
	"os"
	"time"
)

type deathboard struct {
	dashboard int64
	warring   string
	war_msg   string
}

type status struct {
	phone_using bool
}

func main() {
	state := status{
		phone_using: false,
	}
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)

		go func() {
			time.Sleep(200 * time.Millisecond)
			m.Broadcast([]byte("Wattap? Bitch!"))
		}()
	})

	r.GET("/status/phone", func(c *gin.Context) {
		now_using := c.Query("now_using")

		if now_using == "1" {
			state.phone_using = true
		} else {
			state.phone_using = false
		}

		c.String(http.StatusOK, "Setting now changed in server.")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
		fmt.Println(string(msg))
	})

	go func() {
		for {
			typing_message(m)
		}
	}()

	r.Run(":5000")
}

func typing_message(m *melody.Melody) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	m.Broadcast([]byte(text))
	fmt.Println("You just send the message.")

}