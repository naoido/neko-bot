package ipa

import (
	"encoding/xml"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"neko-bot/discord/bot"
	"neko-bot/redis"
	"net/http"
	"slices"
	"time"
)

type SpecialAlert struct {
	XMLName xml.Name `xml:"RDF"`
	Channel Channel  `xml:"channel"`
	Items   []Item   `xml:"item"`
}

type Channel struct {
	About       string `xml:"about,attr"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       struct {
		Seq struct {
			Items []struct {
				Resource string `xml:"resource,attr"`
			} `xml:"resource"`
		} `xml:"Seq"`
	} `xml:"items"`
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Creator string `xml:"creator"`
	Date    string `xml:"date"`
}

var isError = false

func init() {
	_, err := fetch()
	if err != nil {
		panic(err)
	}
}

func StartWatch() {
	go func() {
		timer := time.NewTicker(10 * time.Minute)
		for {
			<-timer.C
			noticeChannel := redis.Client().Get(redis.Context(), redis.IpaNoticeChannel).Val()
			items, err := fetch()
			if err != nil && !isError {
				bot.SendMessage(noticeChannel, "IPA情報の取得時にエラーが発生しました。")
				isError = true
				continue
			}
			for _, item := range items {
				//date, _ := time.Parse("2006-01-02T15:01:05-07:00", item.Date)
				embed := &discordgo.MessageEmbed{
					Title:       item.Title,
					URL:         item.Link,
					Color:       0xff0000,
					Description: fmt.Sprintf("by %s", item.Creator),
					Timestamp:   item.Date,
				}
				bot.SendMessageEmbed(redis.Client().Get(redis.Context(), redis.IpaNoticeChannel).Val(), embed)
			}
		}
	}()
}

func fetch() ([]Item, error) {
	res, err := http.Get("https://www.ipa.go.jp/security/alert-rss.rdf")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	alert := SpecialAlert{}
	if err = xml.Unmarshal(body, &alert); err != nil {
		return nil, err
	}

	c := redis.Client()
	ctx := redis.Context()
	cached := c.SMembers(ctx, redis.IpaSecurityAlert).Val()

	var newItems []Item
	for _, i := range alert.Items {
		if !slices.Contains(cached, i.Title) {
			newItems = append(newItems, i)
			c.SAdd(ctx, redis.IpaSecurityAlert, i.Title)
		}
	}

	return newItems, nil
}
