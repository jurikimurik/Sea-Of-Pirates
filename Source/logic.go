package source

import (
	"encoding/json"
	"errors"
	"fmt"
)

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
	DrawText(jsonGetParam(dummy, "desk"), true)
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
func jsonGetParam(jmap map[string]any, param string) string {
	if !jsonCheckParam(jmap, param) {
		errorOccured(errorCreate("No such parameter (" + param + ") in JSON!"))
		return ""
	}

	return jmap[param].(string)
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
