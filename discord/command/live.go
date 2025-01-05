package command

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"neko-bot/internal/errors"
	"net/http"
	"net/url"
	"strings"
)

type Live struct {
	Detail
}

type Event struct {
	Date   string `json:"date"`
	Title  string `json:"title"`
	Place  string `json:"place"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

func NewLive(name string, prefix *string) *Live {
	live := &Live{
		Detail: Detail{
			name:   name,
			prefix: prefix,
		},
	}

	live.Detail.Command = live

	return live
}

func (l *Live) GetName() string {
	return l.Detail.name
}

func (l *Live) GetCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        l.GetName(),
		Description: "fukuoka live command",
	}
}

func (l *Live) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !l.Detail.isCommand(i) {
		return
	}

	var Url string
	place := ""
	title := ""

	if i.Message != nil && len(i.Message.Content) > len(*l.Detail.prefix)+len(l.Detail.name)+1 {
		remaining := i.Message.Content[len(*l.Detail.prefix)+len(l.Detail.name)+1:]
		parts := strings.SplitN(remaining, " ", 2)

		if len(parts) > 0 {
			place = parts[0]
		}

		if len(parts) > 1 {
			title = parts[1]
		}
	}

	queryParams := url.Values{}
	if title != "" {
		queryParams.Add("title", title)
	}
	if place != "" && place != `""` {
		queryParams.Add("place", place)
	}

	fmt.Printf("QueryParams: %s", queryParams.Encode())

	baseUrl := "url"
	if len(queryParams) > 0 {
		Url = fmt.Sprintf("%s?%s", baseUrl, queryParams.Encode())
	} else {
		Url = baseUrl
	}

	fmt.Printf("URL: %s", Url)

	resp, err := http.Get(Url)
	errors.Catch(err, "cannot get response")

	defer resp.Body.Close()
	body, er := ioutil.ReadAll(resp.Body)
	errors.Catch(er, "cannot read response body")

	var jsonString string
	err = json.Unmarshal([]byte(body), &jsonString)
	errors.Catch(err, "cannot unmarshal json")

	var events []Event
	err = json.Unmarshal([]byte(jsonString), &events)
	errors.Catch(err, "cannot unmarshal json")
	place = "(" + place + ")"
	thread, er := s.MessageThreadStartComplex(i.ChannelID, i.ID, &discordgo.ThreadStart{
		Name: fmt.Sprintf("福岡%sのライブ情報", place),
	})
	errors.Catch(er, "cannot start thread")

	for _, event := range events {
		Image, _ := fetchImage(event.URL)
		errors.Catch(err, "cannot fetch image")
		if Image == "" {
			Image = ""
		}
		_, _ = s.ChannelMessageSendEmbed(thread.ID, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("タイトル\n%s", event.Title),
			Type:        discordgo.EmbedType("rich"),
			Color:       0x00ff00,
			Description: fmt.Sprintf("### 日付\n%s\n### 場所\n%s\n### %s\n### 応募url \n%s\n", event.Date, event.Place, event.Status, event.URL),
			Image: &discordgo.MessageEmbedImage{
				URL: Image,
			},
		})
	}
}

func (l *Live) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !l.Detail.isPrefix(s, m) {
		return
	}

	var Url string
	place := ""
	title := ""

	if len(m.Content) > len(*l.Detail.prefix)+len(l.Detail.name)+1 {
		remaining := m.Content[len(*l.Detail.prefix)+len(l.Detail.name)+1:]
		parts := strings.SplitN(remaining, " ", 2)

		if len(parts) > 0 {
			place = parts[0]
		}

		if len(parts) > 1 {
			title = parts[1]
		}
	}

	queryParams := url.Values{}
	if title != "" {
		queryParams.Add("title", title)
	}
	if place != "" && place != `""` {
		queryParams.Add("place", place)
	}

	fmt.Printf("QueryParams: %s", queryParams.Encode())

	baseUrl := "url"
	if len(queryParams) > 0 {
		Url = fmt.Sprintf("%s?%s", baseUrl, queryParams.Encode())
	} else {
		Url = baseUrl
	}

	fmt.Printf("URL: %s", Url)

	resp, err := http.Get(Url)
	errors.Catch(err, "cannot get response")

	defer resp.Body.Close()
	body, er := ioutil.ReadAll(resp.Body)
	errors.Catch(er, "cannot read response body")

	var jsonString string
	err = json.Unmarshal([]byte(body), &jsonString)
	errors.Catch(err, "cannot unmarshal json")

	var events []Event
	err = json.Unmarshal([]byte(jsonString), &events)
	errors.Catch(err, "cannot unmarshal json")
	place = "(" + place + ")"
	thread, er := s.MessageThreadStartComplex(m.ChannelID, m.ID, &discordgo.ThreadStart{
		Name: fmt.Sprintf("福岡%sのライブ情報", place),
	})
	errors.Catch(er, "cannot start thread")

	for _, event := range events {
		Image, _ := fetchImage(event.URL)
		errors.Catch(err, "cannot fetch image")
		if Image == "" {
			Image = ""
		}
		_, _ = s.ChannelMessageSendEmbed(thread.ID, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("タイトル\n%s", event.Title),
			Type:        discordgo.EmbedType("rich"),
			Color:       0x00ff00,
			Description: fmt.Sprintf("### 日付\n%s\n### 場所\n%s\n### %s\n### 応募url \n%s\n", event.Date, event.Place, event.Status, event.URL),
			Image: &discordgo.MessageEmbedImage{
				URL: Image,
			},
		})
	}
}

func fetchImage(url string) (string, error) {
	resp, err := http.Get(url)
	errors.Catch(err, "cannot get response")
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	Image := doc.Find(`meta[property="og:image"]`).AttrOr("content", "")
	return Image, nil
}
