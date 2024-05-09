package source

import (
	"context"
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
var errorGUIConfig *gui.TextConfig
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
		errorCheck(err)

		//Placing ship in array
		states[first-1][second-1] = GetLogicStateChange(states[first-1][second-1], newState, useFireLogic)
	}

	board.SetStates(*states)
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

	//Do not forget to draw board on exit!
	defer ui.Draw(Board)

	//Default configuration (empty spaces)
	states := SetupFillBoard(Board)

	//If there is some configuration, use it to draw ships on it
	if shipPlaces != nil {
		FillStatesWith(Board, &states, shipPlaces, gui.Ship, false)
		return Board, states
	}

	return Board, states
}

// ----- GAME    ----------------------------------------------------------------------

// BeginGame is a function that start the whole game process.
func BeginGame() {
	//Prepare screen
	ui = gui.NewGUI(true)
	prepareText := DrawGUIText(1, 1, "Game is loading...", nil)

	//Send HTTP Request to begin the game
	response := http.StartGame(util.JSONGetDummy())
	errorCheck(response.Err)

	//Draw screen
	go ui.Start(context.TODO(), nil)

	for {
		status := prepareGame()
		param, err := util.JSONGetParam(status, "game_status")
		errorCheck(err)
		if param == "game_in_progress" {
			break
		}
		//util.JSONPrint(status)
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
	if errorCheck(response.Err) {
		return map[string]any{}
	}

	info, err := util.JSONToMap(response.Body)
	errorCheck(err)

	return info
}

// enterGameFlow is a function that is responsible for in-game flow.
//
// It waits, consumes input and is resposible for displaying the screen
func enterGameFlow() {

	//Battleship area setup
	setupBoard := http.GetMyGameBoard()
	errorCheck(setupBoard.Err)

	//Retrieving board configuration
	setupShipsDataRaw, err := util.JSONGetParamFromJSON(setupBoard.Body, "board")
	errorCheck(err)
	setupShipsData := setupShipsDataRaw.([]interface{})

	//Creating Player board
	var playerBoard *gui.Board
	playerBoard, playerStates = CreateBoard(1, 5, nil, setupShipsData)

	//Creating Enemy board
	var enemyBoard *gui.Board
	enemyBoard, opponentStates = CreateBoard(50, 5, nil, nil)

	//Real game flow (loop)
	for {

		//Checking status
		status := http.GameStatus()
		errorCheck(status.Err)

		//Checks for game end
		dataMap, err := util.JSONToMap(status.Body)
		errorCheck(err)

		param, err2 := util.JSONGetParam(dataMap, "game_status")
		if errorCheck(err2) {
			WaitSecond()
			continue
		}
		if param.(string) == "ended" {
			break
		}

		//If there is no "should_fire" param, wait for your turn
		if !util.JSONCheckParam(dataMap, "should_fire") {
			WaitSecond()
			continue
		}

		// Get opponents shots coordinates
		enemyShots, assert := dataMap["opp_shots"].([]interface{})
		if !assert {
			errorCheck(errors.New("caution: assertion of enemyShots is not successful"))
		}

		//Filling board of player with shots from opponent
		FillStatesWith(playerBoard, &playerStates, enemyShots, gui.Hit, true)

		//Showing up text indicating turn of the player
		turnText := DrawGUIText(15, 0, "Your turn!", nil)
		char := enemyBoard.Listen(context.TODO())
		ui.Remove(turnText)

		// Send Fire as HTTP request
		response := http.Fire(char)
		errorCheck(response.Err)

		// If shot were accepted by server, proceed
		if response.StatusCode == 200 {
			//Retrieve hit result
			result, err := util.JSONGetParamFromJSON(response.Body, "result")
			errorCheck(err)

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

	//Cleaning up the boards adn nicks
	ui.Remove(playerBoard)
	ui.Remove(enemyBoard)

	EndOfGame()
}

// EndOfGame is responsible for ending battleship game.
//
// It also prints if you won or lose.
func EndOfGame() {
	status := http.GameStatus()
	errorCheck(status.Err)
	dataMap, err := util.JSONToMap(status.Body)
	errorCheck(err)

	gameResultTest := DrawGUIText(1, 1, dataMap["last_game_status"].(string), nil)
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

// ----- ERRORS -----------------------------------------------------------------------
func errorCheck(err error) bool {
	if err != nil {
		errorOccured(err)
		return true
	}
	return false
}

func errorOccured(err error) {
	// TODO: In future, let error print to LOG

	//If there is no config, create one.
	if errorGUIConfig == nil {
		errorGUIConfig = gui.NewTextConfig()
		errorGUIConfig.BgColor = gui.Red
		errorGUIConfig.FgColor = gui.Black
	}

	//Draw warning
	DrawGUITextFor(50, 0, err.Error(), errorGUIConfig, 5)
}

// ----- INPUT   ----------------------------------------------------------------------

func ReadInput() string {
	var input string
	_, err := fmt.Scanln(&input)
	if errorCheck(err) {
		return ""
	}
	return input
}
