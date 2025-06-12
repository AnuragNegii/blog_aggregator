package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed,
	error){

	req, err := http.NewRequestWithContext(ctx,
		"GET", feedURL, nil )
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while doing the request: %v\n", err) 
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading: %v\n", err)
	}

	var rsFeed RSSFeed

	if err := xml.Unmarshal(body, &rsFeed); err != nil{
		return nil, fmt.Errorf("error while unmarshaling xml data: %v\n", err)
	}
	rsFeed.Channel.Title = html.UnescapeString(rsFeed.Channel.Title)
	rsFeed.Channel.Description = html.UnescapeString(rsFeed.Channel.Description)

	for _, item := range rsFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)	
		item.Description = html.UnescapeString(item.Description)
	}

	return &rsFeed, nil
}


