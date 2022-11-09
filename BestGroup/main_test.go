// Get /players/{name}/wins
//Post /players/{name}/wins

package main

import "testing"

// // use json_handler package
// func TestRequesting(t *testing.T) {
// 	t.Run("call request with only /players/", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodGet, "/players/", nil)
// 		response := httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		got := response.Body.String()

// 		want := "No player name called"
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})
// 	t.Run("call request with only /players", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodGet, "/players", nil)
// 		response := httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		got := response.Body.String()

// 		want := "No player name called"
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})
// }

// func TestRecordingWinsAndRetrievingThem(t *testing.T) {

// 	t.Run("returns Pepper's score", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
// 		response := httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		got := response.Body.String()

// 		want := "300"
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})
// 	t.Run("returns Salt's score", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodGet, "/players/Salt", nil)
// 		response := httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		got := response.Body.String()

// 		want := "400"
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})
// 	t.Run("update Salt's wins", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodPost, "/players/Salt", nil)
// 		response := httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		request, _ = http.NewRequest(http.MethodGet, "/players/Salt", nil)
// 		response = httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		want := "401"
// 		got := response.Body.String()
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})
// 	t.Run("update Pepper's wins twice", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
// 		response := httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		PlayerServer(response, request)
// 		request, _ = http.NewRequest(http.MethodGet, "/players/Pepper", nil)
// 		response = httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		want := "302"
// 		got := response.Body.String()
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})

// 	t.Run("call request with non-existing player", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodGet, "/players/Jacques", nil)
// 		response := httptest.NewRecorder()
// 		PlayerServer(response, request)
// 		got := response.Body.String()

// 		want := "Player Jacques doesn't exist"
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})

// }

// test base64tohex
func TestBase64ToHex(t *testing.T) {
	t.Run("convert base64 to hex", func(t *testing.T) {
		got := base64Toutf8("qqo=")
		want := "aaaa"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
