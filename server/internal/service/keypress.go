package service

import (
	"unicode"
	"unicode/utf8"

	commondomain "github.com/sceredi/co-type/common/domain"
)

// validateKey delegates to the shared domain validation so both client and server apply identical rules.
func validateKey(player *commondomain.Player, key string, isBackspace bool) error {
	return commondomain.ValidateKeyPress(player, key, isBackspace)
}

// isDelimiter reports whether r is a whitespace character that should be auto-skipped in the snippet.
func isDelimiter(r rune) bool {
	return unicode.IsSpace(r)
}

// typeableChars returns the non-delimiter runes of the snippet in order.
func typeableChars(snippet string) []rune {
	out := make([]rune, 0, len(snippet))
	for _, r := range snippet {
		if !isDelimiter(r) {
			out = append(out, r)
		}
	}
	return out
}

// applyKeyPress updates game in place and returns true when the game has ended.
func applyKeyPress(game *commondomain.GameInfo, key string, isBackspace bool) (ended bool) {
	typeable := typeableChars(game.Snippet)
	total := int64(len(typeable))

	if isBackspace {
		if game.WrongChars > 0 {
			game.WrongChars--
		}
		game.Revision++
		return false
	}

	r, _ := utf8.DecodeRuneInString(key)
	cursorPos := game.CorrectChars + game.WrongChars

	if game.WrongChars <= 0 {
		if cursorPos < total && typeable[cursorPos] == r {
			game.CorrectChars++
		} else {
			game.WrongChars++
		}
	} else {
		game.WrongChars++
	}

	game.Revision++

	return game.CorrectChars >= total && game.WrongChars <= 0
}
