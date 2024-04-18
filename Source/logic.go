package source

import "fmt"

// ----- TEXTS   ----------------------------------------------------------------------

func DrawText(str string, newLine bool) {
	if newLine {
		fmt.Println(str)
	} else {
		fmt.Print(str)
	}
}

// ----- ERRORS -----------------------------------------------------------------------

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
