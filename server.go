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

var state = map[string]bool{
	"phone_using":    false,
	"rotation":       false,
	"speed_up":       false,
	"steering_wheel": false,
	"drive_at_night": false,
	"door_open":      false,
}

var state_msg = map[string]string{
	"phone_using":    "請關閉手機",
	"rotation":       "左右轉請注意後方來車",
	"speed_up":       "您的速度已加快，請注意生命安全",
	"steering_wheel": "請打起精神來專心開車",
	"drive_at_night": "夜間行駛請放慢速度",
	"door_open":      "車門已開啟，請注意後方來車",
}

func main() {
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

	//
	r.GET("/vibrate", func(c *gin.Context) {
		var healthDangerValue int8 = 0
		checkoutAllStatus(&healthDangerValue)

		if healthDangerValue > 80 {
			c.String(http.StatusOK, "yes")
		} else {
			c.String(http.StatusOK, "no")
		}
	})

	//http://localhost:5000/status?now_using=phone_using&sw=0 or 1
	r.GET("/status", func(c *gin.Context) {
		now_using := c.Query("now_using")
		sw := c.Query("sw")

		state[now_using] = (sw == "1")

		if state[now_using] {
			m.Broadcast([]byte(fmt.Sprintf(`{ "warring_msg": "%s" }`, state_msg[now_using])))
		} else {
			m.Broadcast([]byte(fmt.Sprintf(`{ "warring_msg": "請注意，請排除所有障礙，降低生命風險" }`)))
		}

		c.String(http.StatusOK, "Setting now changed in server.")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if string(msg) == "ping" {
			m.Broadcast([]byte("pong"))
			return
		}

		m.Broadcast(msg)
		fmt.Println(string(msg))
	})

	go func() {
		for {
			typing_message(m)
		}
	}()

	go computeHealth(m)

	r.Run(":5000")
}

func typing_message(m *melody.Melody) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	m.Broadcast([]byte(text))
	fmt.Println("You just send the message.")

}

func computeHealth(m *melody.Melody) {
	var healthDangerValue int8 = 0

	for {
		//initialization
		healthDangerValue = 0

		//wait for sleep
		time.Sleep(500 * time.Millisecond)

		//compute add danger value of percent.
		checkoutAllStatus(&healthDangerValue)

		//send back to browser
		json :=
			`{
	"health" : %d
	
}`
		m.Broadcast([]byte(fmt.Sprintf(json, healthDangerValue)))

		if healthDangerValue == 0 {
			m.Broadcast([]byte(fmt.Sprintf(`{ "warring_msg": "您目前很安全" }`)))
		}
	}
}

func checkoutAllStatus(healthDangerValue *int8) {
	if state["phone_using"] {
		*healthDangerValue += 40
	}

	if state["door_open"] {
		*healthDangerValue += 40
	}

	if state["drive_at_night"] {
		*healthDangerValue += 30
	}

	if state["steering_wheel"] {
		*healthDangerValue += 20
	}

	if state["rotation"] {
		*healthDangerValue += 20
	}

	if state["speed_up"] {
		*healthDangerValue += 40
	}

	*healthDangerValue = checkPercent(*healthDangerValue)
}

func checkPercent(v int8) int8 {

	if v > 99 || v < 0 {
		v = 99
	}

	return v
}
