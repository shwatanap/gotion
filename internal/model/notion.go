package model

import (
	"context"
	"os"

	"github.com/jomei/notionapi"
	"golang.org/x/oauth2"
)

type NotionClient struct {
	client *notionapi.Client
}

func NewNotionClient(token string) *NotionClient {
	client := notionapi.NewClient(notionapi.Token(token))
	return &NotionClient{client: client}
}

func NewNotionOAuth() *OAuth {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("NOTION_CLIENT_ID"),
		ClientSecret: os.Getenv("NOTION_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("NOTION_AUTH_URL"),
			TokenURL: os.Getenv("NOTION_TOKEN_URL"),
		},
		RedirectURL: os.Getenv("NOTION_REDIRECT_URL"),
	}
	oauth := &OAuth{
		Config: cfg,
	}
	return oauth
}

type CreateDatabaseRequest struct {
	PageID          string
	Title           string
	CalendarOptions []notionapi.Option
}

func (nc *NotionClient) CreateDatabase(ctx context.Context, req CreateDatabaseRequest) (*notionapi.Database, error) {
	request := &notionapi.DatabaseCreateRequest{
		Parent: notionapi.Parent{
			Type:      notionapi.ParentTypePageID,
			PageID:    notionapi.PageID(req.PageID),
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
					Options: req.CalendarOptions,
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
	res, err := nc.client.Database.Create(ctx, request)
	return res, err
}

type AddEventRequest struct {
	DatabaseID       notionapi.DatabaseID
	Title            string
	SelectedCalendar notionapi.Option
	DateStart        *notionapi.Date
	DateEnd          *notionapi.Date
}

func (nc *NotionClient) AddEvent(ctx context.Context, req AddEventRequest) error {
	_, err := nc.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: req.DatabaseID,
			Workspace:  false,
		},
		Properties: map[string]notionapi.Property{
			"name": &notionapi.TitleProperty{
				ID:   "title",
				Type: notionapi.PropertyTypeTitle,
				Title: []notionapi.RichText{
					{
						Type:     notionapi.ObjectTypeText,
						Text:     &notionapi.Text{Content: req.Title},
						Mention:  nil,
						Equation: nil,
					},
				},
			},
			"calendar": &notionapi.SelectProperty{
				ID:     "calendar",
				Type:   notionapi.PropertyTypeSelect,
				Select: req.SelectedCalendar,
			},
			"date": &notionapi.DateProperty{
				ID:   "date",
				Type: notionapi.PropertyTypeDate,
				Date: &notionapi.DateObject{
					Start: req.DateStart,
					End:   req.DateEnd,
				},
			},
		},
	})
	return err
}

// func (nc *NotionClient) GetDatabase(ctx context.Context, databaseID notionapi.DatabaseID) error {
// 	res, err := nc.client.Database.Query(ctx, databaseID, nil)
// 	if err != nil {
// 		log.Println("😡", err.Error())
// 		return err
// 	}
// 	for _, page := range res.Results {
// 		log.Println("=====================================")
// 		// log.Printf("🥺: %v, %v, %v", page.Properties["calendar"], page.Properties["date"], page.Properties["name"])
// 		// log.Printf("🥺:%v", page.Properties["名前"])
// 		// log.Printf("🥺:%v", page.Properties["name"])
// 		// log.Printf("🥺:%v", page.Properties["date"].(*notionapi.DateProperty).Date)
// 		log.Printf("👏: %v", page.Properties["calendar"])
// 		// log.Printf("🥺: %v", page)
// 		// log.Printf("🥺: %v", page)
// 		// log.Printf("🥺: %v", page)
// 	}
// 	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
// 	db, _ := nc.client.Database.Get(ctx, databaseID)
// 	for key, value := range db.Properties {
// 		log.Printf("%s: %v", key, value)
// 	}
// 	return nil
// }

// func (nc *NotionClient) GetPage(ctx context.Context) {
// 	page, err := nc.client.Page.Get(ctx, "34b55e9279764a4d9fc153657446631b")
// 	if err != nil {
// 		log.Println("😡", err.Error())
// 	}
// 	log.Printf("🥺: %v", page.Parent)
// 	log.Printf("🥺: %v", page.Properties)
// 	// log.Printf("🥺: %v", page.Properties["date"])
// 	log.Printf("🥺: %v", page.Properties["date"].(*notionapi.DateProperty).Date)
// 	log.Printf("🥺: %v", page.Properties["calendar"])
// }
