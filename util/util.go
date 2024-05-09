package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Substring function to get a substring from the string
//
//	Arguments:
//
// str - Basic string where to cut text from.
//
// start - First index from where to start cutting.
//
// end - End index where to stop cutting.
func substring(str string, start, end int) string {
	return strings.TrimSpace(str[start:end])
}

// Function that translates coords (like 'B10') to two separate numbers (like 2 and 10)
//
//	Arguments:
//
// coords - String coordinate (eg. "B10")
//
//	Returns:
//
// int - Translated letter to number as integer
//
// int - Translated string number to integer
//
// error - If error occurs, it will return -1, -1 and occured error
func CoordToIntegers(coords string) (int, int, error) {
	//Checking size of coords
	if len(coords) < 2 {
		return -1, -1, errors.New("coords can't be less than size of 2")
	}

	//Taking out letter from coords and make it small
	letter := substring(coords, 0, 1)
	letter = strings.ToLower(letter)

	// Taking out numbers from coords
	numbers := substring(coords, 1, len(coords))

	//String numbers to integer numbers
	newNumbers, err := strconv.Atoi(numbers)
	if err != nil {
		return -1, -1, err
	}
	return int(rune(letter[0]) - 96), newNumbers, nil
}

// ----- TEXTS   ----------------------------------------------------------------------

// DrawText is a function that prints out text using standard FMT library.
//
//	Arguments:
//
// str - String to be printed
//
// newLine - Should function use println (true) or print.
func DrawText(str string, newLine bool) {
	if newLine {
		fmt.Println(str)
	} else {
		fmt.Print(str)
	}
}

// ----- ERRORS -----------------------------------------------------------------------

// errorCreate creates simple error with given text inside of it.
//
//	Arguments:
//
// str - Message that needs to be inside of the created error.
//
//	Returns:
//
// error - Recently created error with given string inside of it.
func errorCreate(str string) error {
	return errors.New(str)
}

// -----  JSON   -----------------------------------------------------------------------

// mapToJSON is resposible for converting map[string]any to JSON
//
//	Arguments:
//
// jmap - map that needs to be translated to JSON.
//
//	Returns:
//
// []byte - JSON data that was translated from map.
//
// error - If some error occurs, it will return that error and nil as JSON.
func mapToJSON(jmap map[string]any) ([]byte, error) {
	j, err := json.Marshal(jmap)
	if err != nil {
		return nil, err
	}

	return j, nil
}

// JSONToMap is resposible for converting JSON to map[string]any
//
//	Arguments:
//
// j - JSON that needs to be translated to map.
//
//	Returns:
//
// map[string]any - map that was translated from JSON data.
//
// error - If some error occurs, it will return that error and nil as map.
func JSONToMap(j []byte) (map[string]any, error) {
	retrieved := make(map[string]any)
	if err := json.Unmarshal(j, &retrieved); err != nil {
		return nil, err
	}

	return retrieved, nil
}

// -----  JSON (MAP) -------------------------------------------------------------------

// JSONTest is a simple JSON test for functions related to JSON.
func JSONTest() {
	dummy := JSONGetDummy()
	param, _ := JSONGetParam(dummy, "desk")
	DrawText(param.(string), true)
	JSONAddParam(dummy, "age", "21")
	JSONModParam(dummy, "age", "121")
	JSONPrint(dummy)
	JSONDeleteParam(dummy, "age")
	JSONDeleteParam(dummy, "no such param")
	JSONPrint(dummy)
}

// JSONGetEmptyMap return the empty map[string]any.
//
//	Returns:
//
// map[string]any - Recently created empty map.
func JSONGetEmptyMap() map[string]any {
	return map[string]any{}
}

// JSONPrint print out text on screen using DrawText function.
//
//	Arguments:
//
// jmap - map that looks like JSON that needs to be printed.
//
//	Returns:
//
// error - If some error occurs, it will return that error.
func JSONPrint(jmap map[string]any) error {
	info, err := mapToJSON(jmap)
	if err != nil {
		return err
	}

	DrawText(string(info), true)
	return nil
}

// JSONCheckParam checks if there is an requested param in map that looks like JSON.
//
//	Arguments:
//
// jmap - map that looks like JSON where requested parameter might be.
//
// param - string parameter that need to be found.
//
//	Returns:
//
// bool - Result of searching for given parameter.
func JSONCheckParam(jmap map[string]any, param string) bool {
	_, ok := jmap[param]
	return ok
}

// JSONGetParam checks and retrieving the given parameter in map. If there is no
// such parameter, returns nil and error.
//
//	Arguments:
//
// jmap - map where parameter is.
//
// param - string as name of the required parameter.
//
//	Returns:
//
// any - Retrieved parameter as any (assertion RECOMMENDED right after retrieving).
//
// error - If error occurs, returns nil and error with it.
func JSONGetParam(jmap map[string]any, param string) (any, error) {
	if !JSONCheckParam(jmap, param) {
		return nil, errorCreate("No such parameter (" + param + ") in JSON!")
	}

	return jmap[param], nil
}

// JSONGetParamFromJSON retrieving requested parameter after translation from
// JSON to map[string]any.
//
//	Arguments:
//
// jmap - JSON data where parameter might be.
//
// param - The name of the required parameter.
//
//	Returns:
//
// any - Retrieved parameter as any (assertion RECOMMENDED right after retrieving).
//
// error - If error occurs, returns nil and error with it.
func JSONGetParamFromJSON(jmap []byte, param string) (any, error) {
	data, err1 := JSONToMap(jmap)
	if err1 != nil {
		return nil, err1
	}

	info, err2 := JSONGetParam(data, param)
	if err2 != nil {
		return nil, err2
	}

	return info, nil
}

// JSONModParam checks and modify specific parameter in map that looks like JSON.
//
//	Arguments:
//
// jmap - map that looks like JSON where requested parameter might be.
//
// param - string as name of the required parameter.
//
// value - New value for parameter as string.
//
//	Returns:
//
// error - If some error occurs, it will return that error.
func JSONModParam(jmap map[string]any, param string, value string) error {

	if !JSONCheckParam(jmap, param) {
		return errorCreate("No such parameter (" + param + ") in JSON!")
	}

	//Setting new value for given parameter
	jmap[param] = value
	return nil
}

// JSONDeleteParam checks and deletes specific parameter in map that looks like JSON.
//
//	Arguments:
//
// jmap - map that looks like JSON where requested parameter might be.
//
// param - string as name of the required parameter that needs to be deleted.
//
//	Returns:
//
// error - If some error occurs, it will return that error.
func JSONDeleteParam(jmap map[string]any, param string) error {

	//Checks the existing of the parameter
	if !JSONCheckParam(jmap, param) {
		return errorCreate("No such parameter (" + param + ") in JSON!")
	}

	//Deleting the given parameter
	delete(jmap, param)
	return nil
}

// JSONAddParam checks and adds specific parameter into map that looks like JSON.
//
//	Arguments:
//
// jmap - map that looks like JSON where parameter needs to be added.
//
// param - string as name of the new parameter.
//
// value - New value for added parameter as string.
//
//	Returns:
//
// error - If some error occurs, it will return that error.
func JSONAddParam(jmap map[string]any, param string, value string) error {

	//Checks the existing of the parameter
	if JSONCheckParam(jmap, param) {
		return errorCreate("Parameter (" + param + ") is already in JSON!")
	}

	jmap[param] = value
	return nil
}

// Function for getting the dummy JSON text
// JSONGetDummy returns pre-configurated map that looks like JSON. It includes:
//   - coords - Array of strings with locations
//   - desc - Description of the Player.
//   - nick - The nick of the Player.
//   - targetNick - Nick of required opponent.
//   - wpbot - Should it use WP bot as AI opponent.
func JSONGetDummy() map[string]any {
	var coords [20]string = [20]string{"A1", "A3", "B9", "C7", "D1", "D2", "D3", "D4", "D7", "E7", "F1", "F2", "F3", "F5", "G5", "G8", "G9", "I4", "J4", "J8"}
	desc := "My first name"
	nick := "John_Doe_YM"
	targetNick := ""
	wpbot := true
	return map[string]any{"coords": coords, "desc": desc, "nick": nick, "target_nick": targetNick, "wpbot": wpbot}
}
