package javbus

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type NameLinkPair struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type Cast struct {
	NameLinkPair
	Avatar string `json:"avatar"`
}

type Detail struct {
	Title    string         `json:"title"`
	Thumb    string         `json:"thumb"`
	Cover    string         `json:"cover"`
	Date     string         `json:"date"`
	Duration string         `json:"duration"`
	Director NameLinkPair   `json:"director"`
	Studio   NameLinkPair   `json:"studio"`
	Label    NameLinkPair   `json:"label"`
	Stars    []Cast         `json:"stars"`
	Genres   []NameLinkPair `json:"genres"`
}

type Config struct {
	API   string
	Token string
}

type JavBusClient struct {
	client *http.Client
	config *Config
}

func New(config *Config) (client *JavBusClient) {
	if config.API == "" {
		config.API = "https://www.javbus.com"
	}
	client = &JavBusClient{
		config: config,
		client: http.DefaultClient,
	}
	return
}

func (javbus *JavBusClient) request(path string) (doc *html.Node, err error) {
	url := javbus.config.API + path
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Cookie", fmt.Sprintf("bus_auth=%s", javbus.config.Token))
	res, err := javbus.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	doc, err = htmlquery.Parse(res.Body)
	return
}

func (javbus *JavBusClient) Search(keyword string) {
}

func getThumbnail(pic string) (url string) {
	name := filepath.Base(pic)
	id := strings.Split(name, "_")[0]
	url = fmt.Sprintf("https://pics.javbus.com/thumb/%s.jpg", id)
	return
}

func (javbus *JavBusClient) GetDetail(id string) (detail Detail, err error) {
	doc, err := javbus.request(fmt.Sprintf("/%s", id))
	if err != nil {
		return
	}
	pic := htmlquery.SelectAttr(htmlquery.FindOne(doc, "/html/body/div[5]/div[1]/div[1]/a"), "href")
	detail.Title = htmlquery.InnerText(htmlquery.FindOne(doc, "/html/body/div[5]/h3"))
	detail.Cover = javbus.config.API + pic
	detail.Thumb = getThumbnail(pic)
	detail.Date = strings.TrimSpace(htmlquery.InnerText(htmlquery.FindOne(doc, "/html/body/div[5]/div[1]/div[2]/p[2]/text()")))
	detail.Duration = strings.TrimSpace(htmlquery.InnerText(htmlquery.FindOne(doc, "/html/body/div[5]/div[1]/div[2]/p[3]/text()")))
	items := htmlquery.Find(doc, "/html/body/div[5]/div[1]/div[2]//a")
	for _, link := range items {
		text := htmlquery.InnerText(link)
		href := htmlquery.SelectAttr(link, "href")
		if strings.Contains(href, "director") {
			detail.Director = NameLinkPair{
				Name: text,
				Link: href,
			}
		} else if strings.Contains(href, "studio") {
			detail.Studio = NameLinkPair{
				Name: text,
				Link: href,
			}
		} else if strings.Contains(href, "label") {
			detail.Label = NameLinkPair{
				Name: text,
				Link: href,
			}
		} else if strings.Contains(href, "genre") {
			detail.Genres = append(detail.Genres, NameLinkPair{
				Name: text,
				Link: href,
			})
		} else if strings.Contains(href, "star") {
			img := htmlquery.FindOne(link, "img")
			if img == nil {
				continue
			}
			star := Cast{
				NameLinkPair: NameLinkPair{
					Name: htmlquery.SelectAttr(img, "title"),
					Link: htmlquery.SelectAttr(link, "href"),
				},
				Avatar: javbus.config.API + htmlquery.SelectAttr(img, "src"),
			}
			detail.Stars = append(detail.Stars, star)
		}
	}
	return
}
