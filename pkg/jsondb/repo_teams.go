package jsondb

import (
	"encoding/json"
	"fmt"
	"io"
)

func (db *fileDatabase) ListTeams() ([]Team, error) {
	db.teamsReadLocker.Lock()
	defer db.teamsReadLocker.Unlock()

	schema, err := db.readTeams()
	if err != nil {
		return nil, err
	}

	return schema.Teams, nil
}

func (db *fileDatabase) GetTeam(id uint64) (*Team, error) {
	db.teamsReadLocker.Lock()
	defer db.teamsReadLocker.Unlock()

	teams, err := db.ListTeams()
	if err != nil {
		return nil, err
	}

	for _, t := range teams {
		if t.ID == id {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("no team found matching ID %d", id)
}

func (db *fileDatabase) AddTeam(t *Team) error {
	db.teamsWriteLocker.Lock()
	defer db.teamsWriteLocker.Unlock()

	schema, err := db.readTeams()
	if err != nil {
		return err
	}

	t.ID = schema.NextTeamID
	schema.NextTeamID++

	for index := range t.Drivers {
		t.Drivers[index].ID = schema.NextDriverID
		schema.NextDriverID++
	}

	schema.Teams = append(schema.Teams, *t)
	if err := db.writeTeams(schema); err != nil {
		return err
	}
	return nil
}

func (db *fileDatabase) UpdateTeam(t *Team) error {
	db.teamsWriteLocker.Lock()
	defer db.teamsWriteLocker.Unlock()

	schema, err := db.readTeams()
	if err != nil {
		return err
	}

	found := false
	for index, extsingTeam := range schema.Teams {
		if extsingTeam.ID == t.ID {
			schema.Teams[index] = *t
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("cant update missing team %s", t.Name)
	}

	if err := db.writeTeams(schema); err != nil {
		return fmt.Errorf("unable to write update teams: %w", err)
	}

	return nil
}

func (db *fileDatabase) DeleteTeam(id uint64) error {
	db.teamsWriteLocker.Lock()
	defer db.teamsWriteLocker.Unlock()

	schema, err := db.readTeams()
	if err != nil {
		return err
	}

	newTeamList := make([]Team, 0, len(schema.Teams))
	for _, extsingTeam := range schema.Teams {
		if extsingTeam.ID != id {
			newTeamList = append(newTeamList, extsingTeam)
		}
	}

	schema.Teams = newTeamList

	if err := db.writeTeams(schema); err != nil {
		return fmt.Errorf("unable to write update teams: %w", err)
	}

	return nil
}

func (db *fileDatabase) readTeams() (*TeamSchema, error) {
	if _, err := db.teamDb.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error resetting file cursor for teams file: %w", err)
	}

	teamBuf, err := io.ReadAll(db.teamDb)
	if err != nil {
		return nil, fmt.Errorf("error reading teams from file: %w", err)
	}

	if len(teamBuf) <= 0 {
		return &TeamSchema{}, nil
	}

	schema := &TeamSchema{}
	if err := json.Unmarshal(teamBuf, schema); err != nil {
		return nil, fmt.Errorf("erro unmarshaling team json: %w", err)
	}

	return schema, nil
}

func (db *fileDatabase) writeTeams(schema *TeamSchema) error {
	if _, err := db.teamDb.Seek(0, 0); err != nil {
		return fmt.Errorf("error resetting file cursor for teams file: %w", err)
	}

	teamBuf, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("unable to unmarshal teams to json: %w", err)
	}

	if err := db.teamDb.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate teams file: %w", err)
	}
	_, err = db.teamDb.Write(teamBuf)
	if err != nil {
		return fmt.Errorf("unable to write teams to file: %w", err)
	}

	return nil
}
