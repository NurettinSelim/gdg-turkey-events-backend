package database

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"github.com/NurettinSelim/gdg-turkey-events-backend/api"
	"google.golang.org/api/option"
	"strconv"
	"time"
)

type QueryType string

const (
	ALL      QueryType = "all"
	LATEST   QueryType = "latest"
	UPCOMING QueryType = "upcoming"
	OLD      QueryType = "old"
)

var ValidQueries = map[QueryType]bool{
	ALL:      true,
	LATEST:   true,
	UPCOMING: true,
	OLD:      true,
}

type FsDatabase struct {
	client *firestore.Client
}

func (d *FsDatabase) Init() error {
	sa := option.WithCredentialsFile("ServiceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)

	d.client, err = app.Firestore(context.Background())
	return err

}
func (d *FsDatabase) Close() error {
	err := d.client.Close()
	return err
}

func (d *FsDatabase) SaveEvent(event api.Event) (string, error) {
	colRef, _, err := d.client.Collection("events").Add(context.Background(), event)
	if err != nil {
		return "", err
	}
	return colRef.ID, nil
}

func (d *FsDatabase) SaveEvents(events []api.Event) error {
	batch := d.client.Batch()
	eventCollection := d.client.Collection("events")
	for _, event := range events {
		batch.Set(eventCollection.Doc(strconv.Itoa(event.Id)), event)
	}
	_, err := batch.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (d *FsDatabase) GetEventIds() ([]int, error) {
	eventCollection := d.client.Collection("events")
	documentSnapshots, err := eventCollection.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}
	var eventIds []int
	for _, snapshot := range documentSnapshots {
		event := api.Event{}
		err := snapshot.DataTo(&event)
		if err != nil {
			return nil, err
		}
		eventIds = append(eventIds, event.Id)
	}
	return eventIds, nil
}

func (d *FsDatabase) GetEvents(queryType QueryType, page int, pageSize int) ([]api.Event, error) {
	eventCollection := d.client.Collection("events").Offset((page-1)*pageSize).Limit(pageSize).OrderBy("start_date", firestore.Desc)

	var events []api.Event

	switch queryType {
	case ALL:
		documentSnapshots, err := eventCollection.Documents(context.Background()).GetAll()

		if err != nil {
			return nil, err
		}

		for _, snapshot := range documentSnapshots {
			event := api.Event{}

			err := snapshot.DataTo(&event)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	case LATEST:
		documentSnapshots, err := eventCollection.Documents(context.Background()).GetAll()

		if err != nil {
			return nil, err
		}

		for _, snapshot := range documentSnapshots {
			event := api.Event{}

			err := snapshot.DataTo(&event)
			if err != nil {
				return nil, err
			}
			daysAgo := time.Now().AddDate(0, 0, -3)
			if daysAgo.Before(snapshot.CreateTime) {
				events = append(events, event)
			}
		}
	case OLD:
		query := eventCollection.Where("start_date", "<", time.Now())
		documentSnapshots, err := query.Documents(context.Background()).GetAll()

		if err != nil {
			return nil, err
		}

		for _, snapshot := range documentSnapshots {
			event := api.Event{}

			err := snapshot.DataTo(&event)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	case UPCOMING:
		query := eventCollection.Where("start_date", ">", time.Now())
		documentSnapshots, err := query.Documents(context.Background()).GetAll()

		if err != nil {
			return nil, err
		}

		for _, snapshot := range documentSnapshots {
			event := api.Event{}

			err := snapshot.DataTo(&event)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	}

	return events, nil
}
