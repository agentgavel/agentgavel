package security

import (
	"strings"
	"unicode"
)

// NormalizeForMatch applies ADR 009 CI normalization: casefold and whitespace
// collapse. Used for exact-nonce compare after normalize.
func NormalizeForMatch(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	b.Grow(len(s))
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace && b.Len() > 0 {
				b.WriteByte(' ')
				prevSpace = true
			}
			continue
		}
		prevSpace = false
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

// FoldLeetspeak maps a small fixed set of leetspeak / punctuation variants
// to letters (ADR 009 string-variant matcher).
func FoldLeetspeak(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '4', '@':
			b.WriteByte('a')
		case '3':
			b.WriteByte('e')
		case '1', '!':
			b.WriteByte('i')
		case '0':
			b.WriteByte('o')
		case '5', '$':
			b.WriteByte('s')
		case '7':
			b.WriteByte('t')
		case '.', ',', ';', ':', '-', '_', '\'', '"':
			// drop punctuation noise
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// StringVariantContains reports whether haystack contains needle under ADR 009
// string-variant rules (normalize + leetspeak fold on both sides).
func StringVariantContains(haystack, needle string) bool {
	if needle == "" {
		return false
	}
	h := FoldLeetspeak(NormalizeForMatch(haystack))
	n := FoldLeetspeak(NormalizeForMatch(needle))
	if n == "" {
		return false
	}
	return strings.Contains(h, n)
}

// ExactNonceContains reports whether haystack contains the exact nonce after
// ADR 009 normalize (no leetspeak — nonce must survive as the issued token).
func ExactNonceContains(haystack, nonce string) bool {
	if nonce == "" {
		return false
	}
	return strings.Contains(NormalizeForMatch(haystack), NormalizeForMatch(nonce))
}
