package domain

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

// ValidateKeyPress checks whether the given key press is permitted for the player.
// It returns ErrDeleteNotAllowed if a backspace is sent but the player's CanDelete is false,
// and ErrKeyNotAllowed if the key fails printability, AllowedCharacters, or BlockedCharacters checks.
// Both AllowedCharacters and BlockedCharacters are treated as optional: an empty string skips the check.
func ValidateKeyPress(player *Player, key string, isBackspace bool) error {
	if isBackspace {
		if !player.CanDelete {
			return ErrDeleteNotAllowed
		}
		return nil
	}

	r, size := utf8.DecodeRuneInString(key)
	if r == utf8.RuneError || size != len(key) || len(key) == 0 {
		return ErrKeyNotAllowed
	}

	if unicode.IsSpace(r) || !unicode.IsPrint(r) {
		return ErrKeyNotAllowed
	}

	if player.AllowedCharacters != "" {
		re, err := regexp.Compile(player.AllowedCharacters)
		if err != nil || !re.MatchString(key) {
			return ErrKeyNotAllowed
		}
	}

	if player.BlockedCharacters != "" {
		re, err := regexp.Compile(player.BlockedCharacters)
		if err == nil && re.MatchString(key) {
			return ErrKeyNotAllowed
		}
	}

	return nil
}
