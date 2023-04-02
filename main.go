package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/robfig/cron/v3"

	"github.com/chriszychen/pttpost-line-bot/config"
)

var (
	bot     *linebot.Client
	cronjob *cron.Cron
)

type ApiRes struct {
	Data interface{} `json:"data" description:"API response data"`
}

func main() {
	config.Init()

	gin.ForceConsoleColor()

	var err error
	bot, err = linebot.New(config.Config.ChannelSecret, config.Config.ChannelAccessToken)
	if err != nil {
		fmt.Println("linebot.New error happens, err:", err)
	}

	r := gin.Default()
	r.POST("/callback", handler)

	// r.Run(":3000")
	r.Run() // default localhost:8080
}

func handler(ctx *gin.Context) {

	events, err := bot.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			ctx.JSON(http.StatusBadRequest, ApiRes{
				Data: "wrong signature",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, ApiRes{
				Data: "server error",
			})
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			token := event.ReplyToken

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// print UserID
				userID := event.Source.UserID
				fmt.Println("UserID:", userID)

				msg := message.Text
				switch msg {
				case "turn on":
					startCron()
					if err := replyLineMsg(token, "cron start success!"); err != nil {
						fmt.Println("cron start push msg error happens, err:", err)
					}
				case "turn off":
					stopCron()
					if err := replyLineMsg(token, "cron stop success!"); err != nil {
						fmt.Println("cron stop push msg error happens, err:", err)
					}
				default:
					sameMsg := "UserID: " + userID + ",\nGet TextMessage: " + msg + " , \n OK!"
					if err := replyLineMsg(token, sameMsg); err != nil {
						fmt.Println("push same msg error happens, err:", err)
					}
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, ApiRes{})
}

func crawlPTTPost() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}

func startCron() {
	cronjob = cron.New()
	cronjob.AddFunc("* * * * *", func() { fmt.Println("Every minute") })
	cronjob.Start()
}

func stopCron() {
	cronjob.Stop()
}

func replyLineMsg(replyToken, msg string) error {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(msg),
	).Do(); err != nil {
		return err
	}

	return nil
}

func pushLineMsg(pushMsg string) error {
	if _, err := bot.PushMessage(
		config.Config.SelfLineID,
		linebot.NewTextMessage(pushMsg),
	).Do(); err != nil {
		return err
	}

	return nil
}
