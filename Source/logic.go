package source

import (
	"encoding/json"
	"errors"
	"fmt"
	http "sea-of-pirates/HTTP"
	"time"
)

// ----- GLOBAL  ----------------------------------------------------------------------

var shots []string

// ----- GAME    ----------------------------------------------------------------------

// BeginGame is a function that start the whole game process.
func BeginGame() {
	response := http.StartGame(JSONGetDummy())
	if response.Err != nil {
		errorOccured(response.Err)
	}

	for {
		status := prepareGame()
		if jsonGetParam(status, "game_status") == "game_in_progress" {
			break
		}
		jsonPrint(status)
		WaitSecond()
	}

	enterGameFlow()
}

// prepareGame is a function that is responsible for pre-game preparations.
//
// It is also responsible for showing status screen and can be called many times!
//
// Returns:
//
//	map[string]any - Body of HTTP request as map
func prepareGame() map[string]any {
	response := http.GameStatus()
	if response.Err != nil {
		errorOccured(response.Err)
		return map[string]any{}
	}

	return jsonToMap(response.Body)
}

// enterGameFlow is a function that is responsible for in-game flow.
//
// It waits, consumes input and is resposible for displaying the screen
func enterGameFlow() {

	//Battleship area setup
	setupBoard := http.GetMyGameBoard()
	if setupBoard.Err != nil {
		errorOccured(setupBoard.Err)
	}

	jsonPrint(jsonToMap(setupBoard.Body))

	//Real game flow (loop)
	for {
		status := http.GameStatus()
		if status.Err != nil {
			errorOccured(status.Err)
			return
		}

		//Checks for game end
		dataMap := jsonToMap(status.Body)
		if jsonGetParam(dataMap, "game_status").(string) == "ended" {
			break
		}

		//If there is no "should_fire" param, wait for your turn
		if !jsonCheckParam(dataMap, "should_fire") {
			WaitSecond()
			continue
		}

		// Get opponents shots coordinates
		enemyShots, assert := dataMap["opp_shots"].([]interface{})
		if !assert {
			errorOccured(errors.New("caution: assertion of enemyShots is not successful"))
		}
		// Show up opponents shots coordinates
		DrawText("Opponents shots: ", true)
		for _, value := range enemyShots {
			DrawText(value.(string)+", ", false)
		}
		DrawText("\n", true)

		DrawText("Your shots: ", true)
		for _, value := range shots {
			DrawText(value+", ", false)
		}
		DrawText("\n", true)
		// Show text that it is your turn
		DrawText("Your turn!", true)
		// Take input of coordinates
		input := ReadInput()
		// Send Fire as HTTP request
		response := http.Fire(input)
		if response.Err != nil {
			errorOccured(response.Err)
		}
		// If shot were accepted by server, add it to the shots
		if response.StatusCode == 200 {
			shots = append(shots, input)
		}

		jsonPrint(jsonToMap(response.Body))
	}
}

// EndOfGame is responsible for ending battleship game.
//
// It also prints if you won or lose.
func EndOfGame() {
	status := http.GameStatus()
	if status.Err != nil {
		errorOccured(status.Err)
	}
	dataMap := jsonToMap(status.Body)

	DrawText(dataMap["last_game_status"].(string), true)

	clear(shots)
}

// WaitSecond is function that forcing thread to get some sleep for 1 second.
func WaitSecond() {
	time.Sleep(1 * time.Second)
}

// ----- TEXTS   ----------------------------------------------------------------------

func DrawText(str string, newLine bool) {
	if newLine {
		fmt.Println(str)
	} else {
		fmt.Print(str)
	}
}

// ----- ERRORS -----------------------------------------------------------------------

func errorCreate(str string) error {
	return errors.New(str)
}

func errorOccured(err error) {
	// TODO: In future, let error print to LOG
	DrawText("ERROR: "+err.Error(), true)
}

// ----- INPUT   ----------------------------------------------------------------------

func ReadInput() string {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		errorOccured(err)
		return ""
	}
	return input
}

// -----  JSON   -----------------------------------------------------------------------
// Funciton for converting map[string]any to JSON
func mapToJSON(jmap map[string]any) []byte {
	j, err := json.Marshal(jmap)
	if err != nil {
		errorOccured(err)
	}

	return j
}

// Function for convecting JSON to map[string]any
func jsonToMap(j []byte) map[string]any {
	retrieved := make(map[string]any)
	if err := json.Unmarshal(j, &retrieved); err != nil {
		errorOccured(err)
		return retrieved
	}

	return retrieved
}

// -----  JSON (MAP) -------------------------------------------------------------------

// Funciton for JSON testing
func JSONTest() {
	dummy := JSONGetDummy()
	DrawText(jsonGetParam(dummy, "desk").(string), true)
	jsonAddParam(dummy, "age", "21")
	jsonModParam(dummy, "age", "121")
	jsonPrint(dummy)
	jsonDeleteParam(dummy, "age")
	jsonDeleteParam(dummy, "no such param")
	jsonPrint(dummy)
}

// Funciton for getting empty JSON map
func jsonGetEmptyMap() map[string]any {
	return map[string]any{}
}

// Function for printing JSON
func jsonPrint(jmap map[string]any) {
	DrawText(string(mapToJSON(jmap)), true)
}

// Function for checking existing of the given parameter in given JSON
func jsonCheckParam(jmap map[string]any, param string) bool {
	_, ok := jmap[param]
	return ok
}

// Function for getting value of requested param in JSON
func jsonGetParam(jmap map[string]any, param string) any {
	if !jsonCheckParam(jmap, param) {
		errorOccured(errorCreate("No such parameter (" + param + ") in JSON!"))
		return ""
	}

	return jmap[param]
}

// Funciton for modifing parameter in given JSON
func jsonModParam(jmap map[string]any, param string, value string) {

	if !jsonCheckParam(jmap, param) {
		errorOccured(errorCreate("No such parameter (" + param + ") in JSON!"))
		return
	}

	//Setting new value for given parameter
	jmap[param] = value
}

func jsonDeleteParam(jmap map[string]any, param string) {

	//Checks the existing of the parameter
	if !jsonCheckParam(jmap, param) {
		errorOccured(errorCreate("No such parameter (" + param + ") in JSON!"))
		return
	}

	//Deleting the given parameter
	delete(jmap, param)
}

// Funciton for modifing parameter in given JSON
func jsonAddParam(jmap map[string]any, param string, value string) {

	//Checks the existing of the parameter
	if jsonCheckParam(jmap, param) {
		errorOccured(errorCreate("Parameter (" + param + ") is already in JSON!"))
		return
	}

	jmap[param] = value
}

// Function for getting the dummy JSON text
func JSONGetDummy() map[string]any {
	var coords [20]string = [20]string{"A1", "A3", "B9", "C7", "D1", "D2", "D3", "D4", "D7", "E7", "F1", "F2", "F3", "F5", "G5", "G8", "G9", "I4", "J4", "J8"}
	desc := "My first name"
	nick := "John_Doe_YM"
	targetNick := ""
	wpbot := true
	return map[string]any{"coords": coords, "desc": desc, "nick": nick, "target_nick": targetNick, "wpbot": wpbot}
}
