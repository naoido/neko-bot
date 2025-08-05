package ipa

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"neko-bot/internal/infra/redis"
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
	_, err := fetch(context.Background())
	if err != nil {
		panic(err)
	}
}

func StartWatch(s *discordgo.Session) {
	go func() {
		timer := time.NewTicker(10 * time.Minute)
		for {
			<-timer.C
			ctx := context.Background()
			noticeChannel := redis.Client().Get(ctx, redis.IpaNoticeChannel).Val()
			items, err := fetch(ctx)
			if err != nil && !isError {
				s.ChannelMessageSend(noticeChannel, "IPA情報の取得時にエラーが発生しました。")
				isError = true
				continue
			}
			for _, item := range items {
				embed := &discordgo.MessageEmbed{
					Title:       item.Title,
					URL:         item.Link,
					Color:       0xff0000,
					Description: fmt.Sprintf("by %s", item.Creator),
					Timestamp:   item.Date,
				}
				s.ChannelMessageSendEmbed(redis.Client().Get(ctx, redis.IpaNoticeChannel).Val(), embed)
			}
		}
	}()
}

func fetch(ctx context.Context) ([]Item, error) {
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
