package terminal

import (
	"os"
	"strings"

	"golang.org/x/term"
)

func terminalWidth() (width int) {

	var err error
	width, _, err = term.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		// Assume normal 80 char width line
		width = 80
	}

	return
}

func complexFormatWidth(stringToFormat string, indent int) string {
	lengthOfStringToFormat := len(stringToFormat)
	returnString := ""

	if stringToFormat[lengthOfStringToFormat-1:] == "\n" {
		// remove newline at end of string before it is formatted
		stringToFormat = stringToFormat[:lengthOfStringToFormat-1]
	}

	stringLines := strings.Split(stringToFormat, "\n")

	for key, line := range stringLines {
		lineLen := len(line)
		if lineLen > terminalWidth() {
			// Format the string
			line = FormatWidth(line[:lineLen], indent) + "\n"
		} else if key != len(stringLines)-1 {
			// add newline except if the line is the last one
			line += "\n"
		}

		returnString += line
	}

	return returnString
}

// FormatWidth formats the width of a string to the max of the terminal width
func FormatWidth(stringToFormat string, indent int) string {

	// if the string is already formatted, handle as a more complex formatting
	// confirm that if a newline character exists that the only one is not at the end of the string
	if strings.Contains(stringToFormat, "\n") && strings.Index(stringToFormat, "\n") < len(stringToFormat)-2 {
		return complexFormatWidth(stringToFormat, indent)
	}

	moreToProcess := true
	startIdx := 0
	width := terminalWidth() - 1
	endIdx := width
	returnString := ""

	runesToFormat := []rune(stringToFormat)

	for moreToProcess {

		if endIdx < len(runesToFormat) {
			if string(runesToFormat[endIdx-1:endIdx]) == " " {
				stringToAdd := string(runesToFormat[startIdx:endIdx])
				if strings.TrimSpace(stringToAdd) == "" {
					// next space is after 80 characters, so just find the next space and append that "word"
					var index int
					for index = startIdx; index < len(runesToFormat); index++ {
						if string(runesToFormat[index:index+1]) == " " {
							break
						}
					}
					stringToAdd = string(runesToFormat[startIdx:index])
					// set endIdx to the index after the index of the space that we found
					endIdx = index + 1
				}
				returnString += stringToAdd + "\n"

				// Indent the next line
				for i := 0; i < indent; i++ {
					returnString += " "
				}

				startIdx = endIdx
				endIdx += width - indent // Decrease number of characters by the indention
			} else {
				endIdx--
				if endIdx < 1 {
					moreToProcess = false
				}
			}
		} else {
			if startIdx < len(runesToFormat) {
				// Add the remainder to the returnString
				returnString += string(runesToFormat[startIdx:])
			}
			moreToProcess = false
		}
	}

	// If there are no spaces in the string to break it up, just return the original string
	if returnString == "" {
		returnString = string(runesToFormat)
	}

	return returnString
}
