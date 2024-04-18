package http

import (
	"bytes"
	"encoding/json"
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

// ----- SERVER -----------------------------------------------------------------------

// --- BASIC -------------------
//func GameStatus() error
//func StartGame() error
//func GetMyGameBoard() error
//func Fire(coord string) error

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
//	Returns:
//
// http.Header - Header of the HTTP request
//
// int - Status code after the HTTP request
//
// []byte - Readed body of the HTTP request
//
// error - Error that might occured while calling this method!
func Call(TYPE string, addURL string, parameters map[string]string, jsonParameters map[string]any) (http.Header, []byte, int, error) {

	// Creating URL with parameters
	finalUrl := urlWithParameters(parameters, addURL)

	// Creating JSON data to send
	json_data, err := json.Marshal(jsonParameters)
	if err != nil {
		return nil, []byte{}, -1, err
	}

	// Preparing variables
	var resp *http.Response

	var header http.Header
	var body []byte
	var statusCode int
	var errHttp error

	switch TYPE {
	case GET:

		// Making the GET request
		resp, err = http.Get(finalUrl)
		if err != nil {
			errHttp = err
		}

	case POST:

		// Making the POST request
		resp, err = http.Post(finalUrl, "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			errHttp = err
		}

	case DELETE:
		// CAUTION!! MAY BE UNSTABLE BECAUSE OF REQUESTION INSTEAD OF RESPONSE!
		// Making the DELETE request
		req, err := http.NewRequest(DELETE, finalUrl, bytes.NewBuffer(json_data))
		if err != nil {
			errHttp = err
		}

		resp = req.Response
	}

	if errHttp != nil {
		return nil, []byte{}, -1, err
	}

	// Reading the header and body
	body, errHttp = io.ReadAll(resp.Body)
	header = resp.Header
	statusCode = resp.StatusCode

	resp.Body.Close()
	return header, body, statusCode, errHttp
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

// ----- SETTERS ----------------------------------------------------------------------

func SetServerURL(URL string) {
	serverURL = URL
}
