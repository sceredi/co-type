package service

import (
	"regexp"
	"unicode"
	"unicode/utf8"

	commondomain "github.com/sceredi/co-type/common/domain"
)

// isDelimiter reports whether r is a whitespace character that should be auto-skipped in the snippet.
func isDelimiter(r rune) bool {
	return unicode.IsSpace(r)
}

// typeableChars returns a slice of the non-delimiter runes in the snippet, in order.
func typeableChars(snippet string) []rune {
	out := make([]rune, 0, len(snippet))
	for _, r := range snippet {
		if !isDelimiter(r) {
			out = append(out, r)
		}
	}
	return out
}

// validateKey returns ErrKeyNotAllowed if the rune does not pass the player's
// AllowedCharacters / BlockedCharacters regex filters, or if it is not a single
// printable non-whitespace character. It returns ErrDeleteNotAllowed when the
// player is not permitted to delete.
func validateKey(player *commondomain.Player, key string, isBackspace bool) error {
	if isBackspace {
		if !player.CanDelete {
			return commondomain.ErrDeleteNotAllowed
		}
		return nil
	}

	// Must be exactly one rune.
	r, size := utf8.DecodeRuneInString(key)
	if r == utf8.RuneError || size != len(key) || len(key) == 0 {
		return commondomain.ErrKeyNotAllowed
	}

	// No whitespace or control characters.
	if unicode.IsSpace(r) || !unicode.IsPrint(r) {
		return commondomain.ErrKeyNotAllowed
	}

	// AllowedCharacters regex must match (if set).
	if player.AllowedCharacters != "" {
		re, err := regexp.Compile(player.AllowedCharacters)
		if err != nil || !re.MatchString(key) {
			return commondomain.ErrKeyNotAllowed
		}
	}

	// BlockedCharacters regex must NOT match (if set).
	if player.BlockedCharacters != "" {
		re, err := regexp.Compile(player.BlockedCharacters)
		if err == nil && re.MatchString(key) {
			return commondomain.ErrKeyNotAllowed
		}
	}

	return nil
}

// applyKeyPress updates game in place according to the key press and returns
// whether the game has ended (all typeable characters typed correctly).
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
	cursorPos := game.CorrectChars + game.WrongChars // index into typeable chars

	if game.WrongChars <= 0 {
		// In the "correct zone" — check if the pressed key matches the snippet.
		if cursorPos < total && typeable[cursorPos] == r {
			game.CorrectChars++
		} else {
			game.WrongChars++
		}
	} else {
		// Already in the "wrong zone" — every key press digs deeper.
		game.WrongChars++
	}

	game.Revision++

	ended = game.CorrectChars >= total && game.WrongChars <= 0
	return ended
}
