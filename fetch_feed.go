package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	req, err := http.NewRequestWithContext(ctx,"GET", feedURL, nil)
	if err != nil{
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp,err := client.Do(req)
	if err != nil{
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rssfeed := RSSFeed{}
	err = xml.Unmarshal(body, &rssfeed)
	if err != nil{
		return nil, err
	}
	rssfeed.Channel.Title = html.UnescapeString(rssfeed.Channel.Title)
	rssfeed.Channel.Description= html.UnescapeString(rssfeed.Channel.Description)

	for i := range rssfeed.Channel.Item{
		rssfeed.Channel.Item[i].Title = html.UnescapeString(rssfeed.Channel.Item[i].Title)
		rssfeed.Channel.Item[i].Description= html.UnescapeString(rssfeed.Channel.Item[i].Description)
	}
	return &rssfeed, nil
}
