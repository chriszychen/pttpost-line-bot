package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/chriszychen/pttpost-line-bot/config"
)

var bot *linebot.Client

type ApiRes struct {
	Data interface{} `json:"data" description:"API response data"`
}

func main() {
	config.Init()

	var err error
	bot, err = linebot.New(config.Config.ChannelSecret, config.Config.ChannelAccessToken)
	if err != nil {
		panic(err)
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
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				// message.ID: Msg unique ID
				// message.Text: Msg text
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+":"+"Get:"+message.Text+" , \n OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, ApiRes{})
}
