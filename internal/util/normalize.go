package util

import (
	"net/url"
	"strings"
)

func CanonURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil { return raw }
	u.Fragment = ""
	q := u.Query()
	// strip common trackers
	for _, k := range []string{"utm_source","utm_medium","utm_campaign","utm_term","utm_content","utm_id"} {
		q.Del(k)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func EscapeTelegram(s string) string {
	// minimal escaping for MarkdownV2
	replacer := strings.NewReplacer(
		"_","\\_","*","\\*","[","\\[","]","\\]", "(", "\\(", ")", "\\)",
		"~","\\~","`","\\`",">","\\>", "#","\\#", "+","\\+","-","\\-",
		"=","\\=","|","\\|","{","\\{","}","\\}",".","\\.","!","\\!",
	)
	return replacer.Replace(s)
}
