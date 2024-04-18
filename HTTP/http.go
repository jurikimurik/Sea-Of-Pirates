package http

// - TO IMPLEMENT ---------------------------------------------------------------------

//func InitGame() error
//func Board() ([]string, error)
//func Status() (*StatusResponse, error)
//func Fire(coord string) (string, error)

// ----- GLOBAL  ----------------------------------------------------------------------

var serverURL string = "https://go-pjatk-server.fly.dev"

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
