package core

import (
	"strings"
	"unicode"
)

// Summarize extracts up to maxSentences sentences from raw HTML/text
func Summarize(raw string, maxSentences int) string {
	raw = strings.TrimSpace(stripHTML(raw))
	if raw == "" {
		return ""
	}

	sents := splitSentences(raw)
	if len(sents) > maxSentences {
		sents = sents[:maxSentences]
	}

	out := strings.Join(sents, " ")
	if len(out) > 650 { // keep Telegram-friendly
		out = out[:650]
		// Trim partial word at the end
		for len(out) > 0 && !unicode.IsLetter(rune(out[len(out)-1])) {
			out = out[:len(out)-1]
		}
		out += "…"
	}

	return out
}

// splitSentences splits text by '.', '!', '?' safely
func splitSentences(text string) []string {
	var sents []string
	start := 0
	for i, r := range text {
		if r == '.' || r == '!' || r == '?' {
			// include punctuation
			sent := strings.TrimSpace(text[start : i+1])
			if sent != "" {
				sents = append(sents, sent)
			}
			start = i + 1
		}
	}
	// add remaining text as last sentence
	remaining := strings.TrimSpace(text[start:])
	if remaining != "" {
		sents = append(sents, remaining)
	}
	return sents
}

// stripHTML removes script/style tags and other HTML tags
func stripHTML(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	// Remove script and style content
	for _, tag := range []string{"script", "style"} {
		for {
			start := strings.Index(strings.ToLower(s), "<"+tag)
			if start == -1 {
				break
			}
			end := strings.Index(strings.ToLower(s[start:]), "</"+tag+">")
			if end == -1 {
				break
			}
			s = s[:start] + s[start+end+len(tag)+3:]
		}
	}
	// Remove all remaining tags
	var out []rune
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			out = append(out, r)
		}
	}
	return string(out)
}
