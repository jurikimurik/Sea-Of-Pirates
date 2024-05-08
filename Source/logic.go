package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	http "sea-of-pirates/HTTP"
	"time"

	util "sea-of-pirates/util"

	gui "github.com/grupawp/warships-gui/v2"
)

// ----- GLOBAL  ----------------------------------------------------------------------
var playerStates [10][10]gui.State
var opponentStates [10][10]gui.State
var ui *gui.GUI

// ----- GUI     ----------------------------------------------------------------------

// DrawGUIText immediately draws text on the screen.
//
//	Arguments:
//
// x - Integer x coordinate of text.
//
// y - Integer y coordinate of text.
//
// text - String text to show up.
//
// cfg - Configuration for text label.
//
//	Returns:
//
// *gui.Text - Pointer at previously created text for future use.
func DrawGUIText(x int, y int, text string, cfg *gui.TextConfig) *gui.Text {
	newText := gui.NewText(x, y, text, cfg)
	ui.Draw(newText)
	return newText
}

// DrawGUITextFor immediately draws text on the screen for specific amount of time!
//
//	Arguments:
//
// x - Integer x coordinate of text.
//
// y - Integer y coordinate of text.
//
// text - String text to show up.
//
// cfg - Configuration for text label.
//
// time - Time in seconds for showing the text up.
func DrawGUITextFor(x int, y int, text string, cfg *gui.TextConfig, time int) {
	timerText := DrawGUIText(x, y, text, cfg)
	for waiting := 0; waiting < time; waiting++ {
		WaitSecond()
	}
	if timerText != nil {
		ui.Remove(timerText)
	}

}

// SetupFillBoard is filling up the entire board with default states.
//
//	Arguments:
//
// Board - UI board that needs to be filled with states.
//
//	Returns:
//
// [10][10]gui.State - Array of states for future use that represents all the states
// inside of the board.
func SetupFillBoard(Board *gui.Board) [10][10]gui.State {

	//Fill board with default state
	states := [10][10]gui.State{}
	for i := range states {
		states[i] = [10]gui.State{}
		for y := range states[i] {
			states[i][y] = gui.Empty
		}
	}
	Board.SetStates(states)

	return states
}

// GetLogicStateChange helps to logically change the state of one field.
// Based on old state of the field, function will return the new logical
// state that needs to be used. Only if useLogic is true, otherwise it will
// always return the new state of the field.
//
//	Arguments:
//
// oldState - The previous state of the field
//
// newState - The new state that will affect the old one
//
// useLogic - Is logic needs to be used. If false - just forcefully return
// the new state of the field.
//
//	Returns:
//
// gui.State - The new state of the current field
func GetLogicStateChange(oldState gui.State, newState gui.State, useLogic bool) gui.State {
	if oldState == newState || !useLogic {
		return newState
	}

	if oldState == gui.Empty && newState == gui.Hit {
		return gui.Miss
	}

	if oldState == gui.Miss && newState == gui.Hit {
		return gui.Miss
	}

	if oldState == gui.Ship && newState == gui.Hit {
		return gui.Hit
	}

	return oldState
}

// FillStatesWith is filling provided state of the board with new states
// depending on recievied places that needs to be updated. It will only
// affect the provided places and can or cannot use the logic.
//
//	Arguments:
//
// board - Pointer on board where states needs to be updated.
//
// states - Pointer on states that are connected with previous board.
//
// places - Array interface with places as string (example of the inside: {"A2", "B5", "I10"}).
//
// newState - The new state that should be applied for places.
//
// useFireLogic - Should it use the fire logic or just force to bring a new.
// state for places.
func FillStatesWith(board *gui.Board, states *[10][10]gui.State, places []interface{}, newState gui.State, useFireLogic bool) {
	for _, value := range places {
		//Retrieving information
		first, second, err := util.CoordToIntegers(value.(string))
		if err != nil {
			errorOccured(err)
			return
		}

		//Placing ship in array
		states[first-1][second-1] = GetLogicStateChange(states[first-1][second-1], newState, useFireLogic)
	}

	board.SetStates(*states)
	//ui.Draw(board)
}

// CreateBoard is creating a board and instatly draws it after the configuration of states.
// Function can also get a setup with were the ships are placed on the map.
//
//	Arguments:
//
// x - Integer for the x coordinate on the screen where board should start
//
// y - Integer for the y coordinate on the screen where board should start
//
// cfg - Board configuration for display
//
// shipPlaces - Additional argument (can be nil) that is used for setting up
// the locations of the ships on the map.
//
//	Returns:
//
// *gui.Board - Newly created pointer on board
//
// [10][10]gui.State - Array of states connected with the board
func CreateBoard(x int, y int, cfg *gui.BoardConfig, shipPlaces []interface{}) (*gui.Board, [10][10]gui.State) {
	//Creating the new board
	Board := gui.NewBoard(x, y, cfg)

	//Default configuration (empty spaces)
	states := SetupFillBoard(Board)

	//If there is some configuration, use it to draw ships on it
	if shipPlaces != nil {
		FillStatesWith(Board, &states, shipPlaces, gui.Ship, false)
		return Board, states
	}

	ui.Draw(Board)

	return Board, states
}

// ----- GAME    ----------------------------------------------------------------------

// BeginGame is a function that start the whole game process.
func BeginGame() {
	//Prepare screen
	ui = gui.NewGUI(true)
	prepareText := DrawGUIText(1, 1, "Game is loading...", nil)

	//Send HTTP Request to begin the game
	response := http.StartGame(JSONGetDummy())
	if response.Err != nil {
		errorOccured(response.Err)
	}

	//Draw screen
	go ui.Start(context.TODO(), nil)

	for {
		status := prepareGame()
		if jsonGetParam(status, "game_status") == "game_in_progress" {
			break
		}
		//jsonPrint(status)
		WaitSecond()
	}

	//Clear screen and enter game flow
	ui.Remove(prepareText)
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

	//Retrieving board configuration
	setupShipsData := jsonGetParam(jsonToMap(setupBoard.Body), "board").([]interface{})

	//Creating Player board
	var playerBoard *gui.Board
	playerBoard, playerStates = CreateBoard(1, 5, nil, setupShipsData)
	ui.Draw(playerBoard)

	//Creating Enemy board
	var enemyBoard *gui.Board
	enemyBoard, opponentStates = CreateBoard(50, 5, nil, nil)
	ui.Draw(enemyBoard)

	//Real game flow (loop)
	for {

		//Checking status
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

		//Filling board of player with shots from opponent
		FillStatesWith(playerBoard, &playerStates, enemyShots, gui.Hit, true)

		//Showing up text indicating turn of the player
		turnText := DrawGUIText(15, 0, "Your turn!", nil)
		char := enemyBoard.Listen(context.TODO())
		ui.Remove(turnText)

		// Send Fire as HTTP request
		response := http.Fire(char)
		if response.Err != nil {
			errorOccured(response.Err)
		}
		// If shot were accepted by server, proceed
		if response.StatusCode == 200 {
			//Retrieve hit result
			result := jsonGetParam(jsonToMap(response.Body), "result")

			//Set up go routine for text with result that shows up for 2 seconds and then dissapears
			go DrawGUITextFor(40, 0, result.(string), nil, 2)

			//Checking the effect of player's shot
			var effect gui.State
			if result == "hit" || result == "sunk" {
				effect = gui.Hit
			} else {
				effect = gui.Miss
			}

			//Creating an array interface with one shot for updating the enemy board
			data := []string{char}
			fire := make([]interface{}, len(data))
			fire[0] = char

			//Updating enemy board with player's shot effect
			FillStatesWith(enemyBoard, &opponentStates, fire, effect, false)
		}
		//Repeat until the end of the game
	}

	EndOfGame()
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

	gameResultTest := DrawGUIText(25, 25, dataMap["last_game_status"].(string), nil)
	WaitSeconds(5)
	ui.Remove(gameResultTest)
}

// WaitSecond is function that forcing thread to get some sleep for 1 second.
func WaitSecond() {
	time.Sleep(1 * time.Second)
}

// WaitSeconds is an additional function that will wait for specific amount of time
//
//	Arguments:
//
// time - Amount of time in seconds as integer
func WaitSeconds(time int) {
	for waiting := 0; waiting < time; waiting++ {
		WaitSecond()
	}
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
