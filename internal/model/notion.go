package model

import (
	"context"

	"github.com/jomei/notionapi"
)

type NotionClient struct {
	Client *notionapi.Client
}

func NewNotionClient(apikey string) *NotionClient {
	client := notionapi.NewClient(notionapi.Token(apikey))
	return &NotionClient{Client: client}
}

type CreateDatabaseRequest struct {
	PageId notionapi.PageID
	Title  string
}

func (nc *NotionClient) CreateDatabase(ctx context.Context, req CreateDatabaseRequest) error {
	request := &notionapi.DatabaseCreateRequest{
		Parent: notionapi.Parent{
			Type:      notionapi.ParentTypePageID,
			PageID:    req.PageId,
			Workspace: false,
		},
		Title: []notionapi.RichText{
			{
				Type:     notionapi.ObjectTypeText,
				Text:     &notionapi.Text{Content: req.Title},
				Mention:  nil,
				Equation: nil,
			},
		},
		Properties: map[string]notionapi.PropertyConfig{
			"name": &notionapi.TitlePropertyConfig{
				ID:    "title",
				Type:  notionapi.PropertyConfigTypeTitle,
				Title: struct{}{},
			},
			"calendar": &notionapi.SelectPropertyConfig{
				ID:   "calendar",
				Type: notionapi.PropertyConfigTypeSelect,
				Select: notionapi.Select{
					Options: []notionapi.Option{},
					// Options: []notionapi.Option{
					// 	{
					// 		ID:    "0",
					// 		Name:  "イベント",
					// 		Color: "blue",
					// 	},
					// 	{
					// 		ID:    "1",
					// 		Name:  "遊び",
					// 		Color: "red",
					// 	},
					// },
				},
			},
			"date": &notionapi.DatePropertyConfig{
				ID:   "date",
				Type: notionapi.PropertyConfigTypeDate,
				Date: struct{}{},
			},
		},
		IsInline: true,
	}
	_, err := nc.Client.Database.Create(ctx, request)
	return err
}
