package jsondb

import (
	"os"
	"sync"
)

type JsonDatabase interface {
	ListTeams() ([]Team, error)
	GetTeam(id uint64) (*Team, error)
	AddTeam(t *Team) error
	UpdateTeam(t *Team) error
	DeleteTeam(id uint64) error

	ListEvents() ([]RaceEvent, error)
	GetEvent(id uint64) (*RaceEvent, error)
	AddEvent(e *RaceEvent) error
	UpdateEvent(e *RaceEvent) error
	DeleteEvent(id uint64) error
}

type fileDatabase struct {
	teamDb   *os.File
	eventsDb *os.File

	teamsReadLocker   sync.Locker
	teamsWriteLocker  sync.Locker
	eventsReadLocker  sync.Locker
	eventsWriteLocker sync.Locker
}

func CreateFileDatabase() JsonDatabase {
	teamFile, err := os.OpenFile("teams.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	eventsFile, err := os.OpenFile("events.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	teamsRwMutex := &sync.RWMutex{}
	eventsRwMutex := &sync.RWMutex{}

	return &fileDatabase{
		teamDb:            teamFile,
		eventsDb:          eventsFile,
		teamsReadLocker:   teamsRwMutex.RLocker(),
		teamsWriteLocker:  teamsRwMutex,
		eventsWriteLocker: eventsRwMutex.RLocker(),
		eventsReadLocker:  eventsRwMutex,
	}
}
