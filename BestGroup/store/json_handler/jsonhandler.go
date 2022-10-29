package jsonhandler

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/tuxago/go/BestGroup/store"
)

type PlayerStorage struct {
	muPlayers sync.RWMutex
	Players   []store.Player
}

func (ps *PlayerStorage) InitJSON(jsonI string) error {
	jsonFile, err := os.Open(jsonI)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	json.Unmarshal(byteValue, &ps.Players)
	return nil
}

func NewPlayerStorage() *PlayerStorage {
	return &PlayerStorage{}
}

func (ps *PlayerStorage) SaveJSON(jsonI string) error {
	ps.muPlayers.Lock()
	defer ps.muPlayers.Unlock()
	// open output file
	jsonFile, err := os.OpenFile(jsonI, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// convert to json
	jsonData, err := json.Marshal(ps.Players)
	if err != nil {
		return err
	}

	//write to file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PlayerStorage) GetPlayer(name string) (store.Player, error) {
	ps.muPlayers.RLock()
	defer ps.muPlayers.RUnlock()
	for _, player := range ps.Players {
		if player.Name == name {
			return player, nil
		}
	}
	return store.Player{}, errors.New("player not found in getplayer")
}

func (ps *PlayerStorage) IncWins(name string) (int, error) {
	ps.muPlayers.Lock()
	defer ps.muPlayers.Unlock()
	for i, player := range ps.Players {
		if player.Name == name {
			wins := ps.Players[i].Wins + 1
			ps.Players[i].Wins = wins
			return wins, nil
		}
	}
	return -1, errors.New("player not found in setplayer")
}

func (ps *PlayerStorage) RemovePlayer(name string) error {
	ps.muPlayers.Lock()
	defer ps.muPlayers.Unlock()
	for i, player := range ps.Players {
		if player.Name == name {
			ps.Players = append(ps.Players[:i], ps.Players[i+1:]...)
			return nil
		}
	}
	return errors.New("player not found in removeplayer")
}

func (ps *PlayerStorage) AddPlayer(name string) error {
	ps.muPlayers.Lock()
	defer ps.muPlayers.Unlock()
	for _, player := range ps.Players {
		if player.Name == name {
			return errors.New("player already exists")
		}
	}
	ps.Players = append(ps.Players, store.Player{Name: name, Wins: 0})
	return nil
}

func (ps *PlayerStorage) GetAllPlayers() ([]store.Player, error) {
	ps.muPlayers.RLock()
	defer ps.muPlayers.RUnlock()

	var players []store.Player
	// we create a full copy of the slice, to avoid access outside of the lock
	copy(players, ps.Players)

	return players, nil
}

func (ps *PlayerStorage) FormatPlayers_old(format string) (string, error) {
	ps.muPlayers.RLock()
	defer ps.muPlayers.RUnlock()
	// sort players by name
	sort.SliceStable(ps.Players, func(i, j int) bool {
		return ps.Players[i].Name < ps.Players[j].Name
	})

	switch format {
	case "string":
		// Players to string in var str
		var str string
		for _, player := range ps.Players {
			str += player.Name + " " + strconv.Itoa(player.Wins) + "\n"
		}
		return str, nil
	case "csv":
		var str string
		str += "Name,Wins\n"
		// Players to string in var str
		for _, player := range ps.Players {
			str += player.Name + "," + strconv.Itoa(player.Wins) + "\n"
		}
		return str, nil
	case "html":
		var str string
		str += "<table>\n"
		str += "<tr><th>Name</th><th>Wins</th></tr>\n"
		// Players to string in var str
		for _, player := range ps.Players {
			str += "<tr><td>" + player.Name + "</td><td>" + strconv.Itoa(player.Wins) + "</td></tr>\n"
		}
		str += "</table>\n"
		return str, nil
	case "xml":
		var str string
		str += "<players>\n"
		// Players to string in var str
		for _, player := range ps.Players {
			str += "<player>\n"
			str += "<name>" + player.Name + "</name>\n"
			str += "<wins>" + strconv.Itoa(player.Wins) + "</wins>\n"
			str += "</player>\n"
		}
		str += "</players>\n"
		return str, nil
	default:
		// convert to json
		jsonData, err := json.Marshal(ps.Players)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}
}

// Make backups of the json file in 3 other files
func (ps *PlayerStorage) Backup(timeMult int, jsonI2 string, jsonI3 string, jsonI4 string) {
	// make a goroutine backup of the original file every 30 minute in jsonI2, jsonI3, jsonI4
	go func() {
		nbrTrack := 0
		for range time.Tick(time.Duration(timeMult) * time.Minute) {
			func() {
				ps.muPlayers.Lock()
				defer ps.muPlayers.Unlock()
				var jsonFile *os.File
				var err error
				switch nbrTrack {
				case 0:
					// open output files
					jsonFile, err = os.OpenFile(jsonI2, os.O_WRONLY|os.O_TRUNC, 0644)
					if err != nil {
						return
					}
				case 1:
					jsonFile, err = os.OpenFile(jsonI2, os.O_WRONLY|os.O_TRUNC, 0644)
					if err != nil {
						return
					}
				default:
					jsonFile, err = os.OpenFile(jsonI2, os.O_WRONLY|os.O_TRUNC, 0644)
					if err != nil {
						return
					}
				}

				defer jsonFile.Close()

				// convert to json
				jsonData, err := json.Marshal(ps.Players)
				if err != nil {
					return
				}

				//write to file
				_, err = jsonFile.Write(jsonData)
				if err != nil {
					return
				}
				nbrTrack++
				if nbrTrack > 2 {
					nbrTrack = 0
				}
			}()
		}
	}()
}
