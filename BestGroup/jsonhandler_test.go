package main

import (
	"testing"
)

var Json string = "test.json"

func TestInitJSON(t *testing.T) {
	t.Run("TestInitJSON", func(t *testing.T) {
		err := InitJSON(Json, 0)
		//if Players is empty, then the test fails
		if err != nil {
			t.Errorf("Error in InitJSON, %v", err)
		}
	})
}

// func TestSaveJSON(t *testing.T) {
// 	t.Run("TestSaveJSON", func(t *testing.T) {
// 		err := InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in InitJSON %v", err)
// 		}
// 		//add a new player to the list
// 		Players.Players = append(Players.Players, {Company: "Test", Address: "AAA"})
// 		err = SaveJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in SaveJSON %v", err)
// 		}
// 		err = InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in second InitJSON %v", err)
// 		}
// 		//if Players does not contain Name : Test with Wins :0, then the test fails
// 		if Players.Players[len(Players.Players)-1].Address != "Test" || Players.Players[len(Players.Players)-1].Id != 0 {
// 			t.Errorf("Players.JPlayers does not contain Name : Test with Wins :0")
// 		}
// 	})

// }

// func TestAddPlayer(t *testing.T) {
// 	t.Run("TestAddPlayer", func(t *testing.T) {
// 		err := InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in InitJSON")
// 		}
// 		err = AddPlayer("Test1")
// 		if err != nil {
// 			t.Errorf("Error in AddPlayer")
// 		}
// 		err = SaveJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in SaveJSON")
// 		}
// 		err = InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in second InitJSON")
// 		}
// 		player, err := GetPlayer("Test1")
// 		if player.Name != "Test1" || err != nil {
// 			t.Errorf("Players.JPlayers does not contain Name : Test with Wins :0, or an error : %v", err)
// 		}
// 	})
// }

// func TestSetPlayer(t *testing.T) {
// 	t.Run("TestSetPlayer 'Test' wins", func(t *testing.T) {
// 		err := InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in InitJSON")
// 		}
// 		wins, err := SetPlayer("Test")
// 		if err != nil {
// 			t.Errorf("Error in SetPlayer")
// 		}
// 		if wins == -1 || err != nil {
// 			t.Errorf("Players.JPlayers does not contain Name : Test with Wins :4, or an error : %v", err)
// 		}
// 	})
// 	t.Run("TestSetPlayer 'Test123' wins and get error", func(t *testing.T) {
// 		err := InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in InitJSON")
// 		}
// 		_, err = SetPlayer("Test123")
// 		if err == nil {
// 			t.Errorf("no error in SetPlayer")
// 		}
// 	})
// }

// func TestGetPlayer(t *testing.T) {
// 	InitJSON(Json)
// 	//if GetPlayer returns an empty JPlayer, then the test fails
// 	player, err := GetPlayer("Test")
// 	if player == (JPlayer{}) || err != nil {
// 		t.Errorf("GetPlayer(\"Test\") returned an empty JPlayer, or an error : %v", err)
// 	}
// }

// func TestRemovePlayer(t *testing.T) {
// 	t.Run("TestRemoveJSON", func(t *testing.T) {
// 		err := InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in InitJSON")
// 		}
// 		err = RemovePlayer("Test")
// 		if err != nil {
// 			t.Errorf("Error in RemovePlayer")
// 		}
// 		err = RemovePlayer("Test1")
// 		if err != nil {
// 			t.Errorf("Error in RemovePlayer")
// 		}
// 		err = SaveJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in SaveJSON")
// 		}
// 		err = InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in second InitJSON")
// 		}
// 		player, err := GetPlayer("Test")
// 		if player.Name == "Test" || err.Error() != "player not found in getplayer" {
// 			t.Errorf("Players.JPlayers does contain Name : Test with Wins :0, or an error : %v", err)
// 		}
// 	})
// }

// func TestFormatPlayers(t *testing.T) {
// 	t.Run("TestFormatPlayers", func(t *testing.T) {
// 		err := InitJSON(Json)
// 		if err != nil {
// 			t.Errorf("Error in InitJSON")
// 		}
// 		players, err := FormatPlayers("")
// 		if err != nil {
// 			t.Errorf("Error in GetPlayers")
// 		}
// 		if len(players) == 0 {
// 			t.Errorf("GetPlayers returned an empty list")
// 		}
// 	})
// }
