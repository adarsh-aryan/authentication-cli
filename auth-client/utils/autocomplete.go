package utils

import "strings"

// define all available cli commands for autocompletion
var availableCommands = []string{"login", "logout", "register", "whoami", "quit", "help"}

// create the autocomplete callback function
func AutoCompleteHook(line string, pos int, key rune) (newLine string, newPos int, ok bool) {

	// we only care about the tab key (which triggers autocomplete)
	if key != '\t' {
		return "", 0, false
	}

	// find commands that start with whatever the user typed
	var matches []string
	for _, cmd := range availableCommands {
		if strings.HasPrefix(cmd, line) {
			matches = append(matches, cmd)
		}
	}

	// if there is no matches , do nothing
	if len(matches) == 0 {
		return "", 0, false
	}

	// if there exactly one match, autocomplete the whole and add a space
	if len(matches) == 1 {
		completed := matches[0] + " "
		return completed, len(completed), true
	}

	// if there are multiple matches, autocomplete it with first match
	completed := matches[0] + " "
	return completed, len(completed), true
}
