package rss

import (
	"encoding/xml"
	"io"
	"time"
)

type Feed struct {
	XMLName  xml.Name  `xml:"rss"`
	Channels []Channel `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	PubDate     ItemDate `xml:"pubDate"`
	Description string   `xml:"description"`
	Enclosure   Enclosure
}

type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  int64    `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

type ItemDate struct {
	time.Time
}

func (date *ItemDate) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var s string
	var err error

	err = decoder.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	date.Time, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", s)

	return err
}

func Parse(source io.Reader) (*[]Channel, error) {
	feed := Feed{}

	if err := xml.NewDecoder(source).Decode(&feed); err != nil {
		return nil, err
	}

	return &feed.Channels, nil
}
