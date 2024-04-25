package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// ----- GLOBAL  ----------------------------------------------------------------------

var token string
var serverURL string = "https://go-pjatk-server.fly.dev/api/"

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

// Response is a all-in-one structure that have all neccessary information of HTTP request:
//
// http.Header - Header of the HTTP request
//
// int - Status code after the HTTP request
//
// []byte - Readed body of the HTTP request
//
// error - Error that might occured while calling HTTP request!
type Response struct {
	Header     http.Header
	Body       []byte
	StatusCode int
	Err        error
}

// ----- SERVER -----------------------------------------------------------------------

// --- BASIC -------------------

// GameStatus is a http function for handling game status HTTP request
//
//	Returns:
//
// Response - All in one structure that have neccessary info of HTTP Request
func GameStatus() Response {
	return call(GET, "game", nil, nil, true)
}

// StartGame is a http function for sending HTTP request of beginning the game.
// It also retrives authorization token!
//
//	Arguments:
//
// jsonParameters - Map of JSON data parameter of request
//
//	Returns:
//
// Response - All in one structure that have neccessary info of HTTP Request
func StartGame(jsonParameters map[string]any) Response {
	resp := call(POST, "game", nil, jsonParameters, false)
	//Check if some error occured
	if resp.Err != nil {
		return resp
	}

	//Retrieving authorization token
	token = resp.Header.Get("X-Auth-Token")
	if token == "" {
		resp.Err = errors.New("can't retrieve authorization token from request")
	}

	//All good
	return resp
}

// GetMyGameBoard is a http function for handling getting the game board in HTTP request
//
//	Returns:
//
// Response - All in one structure that have neccessary info of HTTP Request
func GetMyGameBoard() Response {
	return call(GET, "game/board", nil, nil, true)
}

// Fire is a http function for sending HTTP request of firing to specific coordinate
//
//	Returns:
//
// Response - All in one structure that have neccessary info of HTTP Request
func Fire(coord string) Response {
	jParam := map[string]any{"coord": coord}
	return call(POST, "game/fire", nil, jParam, true)
}

// --- OPTIONAL ----------------
//func GiveUp() error
//func GetMyAndOpponentDesc() error
//func RefreshSession() error
//func Lobby() error
//func Stats() error
//func StatsOfPlayer(nick string) error

// ----- NETWORK ----------------------------------------------------------------------

// Call is a all-in-one function for GET, POST and DELETE HTTP Requests
//
//	Arguments:
//
// TYPE - Type of HTTP Request (GET / POST / DELETE)
//
// addURL - URL to add to serverURL
//
// parameters - Map of parameters to be inside of URL after "?"
//
// jsonParameters - Map of JSON data parameter of request
//
// includeToken - Should token be includen into HTTP request
//
//	Returns:
//
// Response - All in one structure that have neccessary info of HTTP Request
func call(TYPE string, addURL string, parameters map[string]string, jsonParameters map[string]any, includeToken bool) Response {

	// Creating URL with parameters
	finalUrl := urlWithParameters(parameters, addURL)

	// Creating JSON data to send
	json_data, err := json.Marshal(jsonParameters)
	if err != nil {
		return Response{nil, []byte{}, -1, err}
	}

	// Preparing variables
	var resp *http.Response

	var header http.Header
	var body []byte
	var statusCode int
	var errHttp error

	var req *http.Request

	switch TYPE {
	case GET:

		// Making the GET request
		req, err = http.NewRequest(GET, finalUrl, bytes.NewBuffer(json_data))
		if err != nil {
			errHttp = err
		}

	case POST:

		// Making the POST request
		req, err = http.NewRequest(POST, finalUrl, bytes.NewBuffer(json_data))
		if err != nil {
			errHttp = err
		}

	case DELETE:

		// Making the DELETE request
		req, err = http.NewRequest(DELETE, finalUrl, bytes.NewBuffer(json_data))
		if err != nil {
			errHttp = err
		}
	}

	// Adding information to header
	req.Header.Add("Content-Type", "application/json")
	if includeToken {
		req.Header.Add("X-AUTH-TOKEN", token)
	}

	//Making an HTTP request
	resp, errHttp = http.DefaultClient.Do(req)
	if errHttp != nil {
		return Response{nil, []byte{}, -1, err}
	}

	// Reading the header and body
	body, errHttp = io.ReadAll(resp.Body)
	header = resp.Header
	statusCode = resp.StatusCode

	resp.Body.Close()

	//Packing all information into one single response
	packagedResponse := Response{header, body, statusCode, errHttp}
	return packagedResponse
}

// Function for creating final URL with parameters after "?"
//
//	Arguments:
//
// parameters - parameters to add after "?".
//
// addURL - text to add after serverURL.
//
//	Returns:
//
// finalURL - modified URL string.
func urlWithParameters(parameters map[string]string, addURL string) string {
	URLParameters := url.Values{}
	for key, value := range parameters {
		URLParameters.Add(key, value)
	}
	finalUrl := serverURL + addURL + "?" + URLParameters.Encode()
	return finalUrl
}

// ----- TESTS   ----------------------------------------------------------------------

func TestPackage() string {
	return "HTTP package is working just fine!"
}

// ----- GETTERS ----------------------------------------------------------------------

func GetServerURL() string {
	return serverURL
}

func GetToken() string {
	return token
}

// ----- SETTERS ----------------------------------------------------------------------

func SetServerURL(URL string) {
	serverURL = URL
}
