package server

import (
	"fmt"
	"html/template"
	"time"
)

func formatAsDate(t time.Time) string {
	return fmt.Sprintf("%d.%02d.%02d %02d:%02d:%02d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second())
}

func formatAsPrice(price int) string {
	realPrice := float64(price) / 100
	return fmt.Sprintf("%.2f", realPrice)
}

func columnStatus(status bool) template.HTML {
	result := ""
	if status {
		result = "<span class=\"column-green\">&nbsp;Active&nbsp;</span>"
	} else {
		result = "<span class=\"column-red\">Inactive</span>"
	}
	return template.HTML(result)
}

func keyBytesToString(data []byte) template.HTML {
	return template.HTML(string(data[:]))
}
