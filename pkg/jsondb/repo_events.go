package jsondb

import (
	"encoding/json"
	"fmt"
	"io"
)

func (db *fileDatabase) ListEvents() ([]RaceEvent, error) {
	db.eventsReadLocker.Lock()
	defer db.eventsReadLocker.Unlock()

	schema, err := db.readEvents()
	if err != nil {
		return nil, err
	}

	return schema.Events, nil
}

func (db *fileDatabase) AddEvent(e *RaceEvent) error {
	db.eventsWriteLocker.Lock()
	defer db.eventsWriteLocker.Unlock()

	schema, err := db.readEvents()
	if err != nil {
		return err
	}

	e.ID = schema.NextEventID
	schema.NextEventID++

	schema.Events = append(schema.Events, *e)
	if err := db.writeEvents(schema); err != nil {
		return err
	}
	return nil
}

func (db *fileDatabase) GetEvent(id uint64) (*RaceEvent, error) {
	db.eventsWriteLocker.Lock()
	defer db.eventsWriteLocker.Unlock()

	schema, err := db.readEvents()
	if err != nil {
		return nil, fmt.Errorf("unable to read events: %w", err)
	}

	for _, existingEvent := range schema.Events {
		if existingEvent.ID == id {
			return &existingEvent, nil
		}
	}

	return nil, fmt.Errorf("missing event %d", id)
}

func (db *fileDatabase) UpdateEvent(e *RaceEvent) error {
	db.eventsWriteLocker.Lock()
	defer db.eventsWriteLocker.Unlock()

	schema, err := db.readEvents()
	if err != nil {
		return err
	}

	found := false
	for index, existingEvent := range schema.Events {
		if existingEvent.ID == e.ID {
			schema.Events[index] = *e
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("missing event %d", e.ID)
	}

	if err := db.writeEvents(schema); err != nil {
		return err
	}
	return nil
}

func (db *fileDatabase) DeleteEvent(id uint64) error {
	db.eventsWriteLocker.Lock()
	defer db.eventsWriteLocker.Unlock()

	schema, err := db.readEvents()
	if err != nil {
		return err
	}

	filteredEvents := make([]RaceEvent, 0, len(schema.Events))
	for _, existingEvent := range schema.Events {
		if existingEvent.ID != id {
			filteredEvents = append(filteredEvents, existingEvent)
		}
	}

	schema.Events = filteredEvents

	if err := db.writeEvents(schema); err != nil {
		return err
	}
	return nil
}

func (db *fileDatabase) readEvents() (*EventSchema, error) {
	if _, err := db.eventsDb.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error resetting file cursor for event file: %w", err)
	}

	eventsBuf, err := io.ReadAll(db.eventsDb)
	if err != nil {
		return nil, fmt.Errorf("error reading events from file: %w", err)
	}

	if len(eventsBuf) <= 0 {
		return &EventSchema{}, nil
	}

	schema := &EventSchema{}
	if err := json.Unmarshal(eventsBuf, schema); err != nil {
		return nil, fmt.Errorf("erro unmarshaling event json: %w", err)
	}

	return schema, nil
}

func (db *fileDatabase) writeEvents(eventSchema *EventSchema) error {
	if _, err := db.eventsDb.Seek(0, 0); err != nil {
		return fmt.Errorf("error resetting file cursor for event file: %w", err)
	}

	eventBuf, err := json.Marshal(eventSchema)
	if err != nil {
		return fmt.Errorf("unable to unmarshal events to json: %w", err)
	}

	if err := db.eventsDb.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate event file: %w", err)
	}
	_, err = db.eventsDb.Write(eventBuf)
	if err != nil {
		return fmt.Errorf("unable to write events to file: %w", err)
	}

	return nil
}
