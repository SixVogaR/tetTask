package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
)

func main() {

	doc, err := xmlquery.LoadURL("https://www.bank.lv/vk/ecb_rss.xml")
	if err != nil {
		panic(err)
	}

	for i, n := range xmlquery.Find(doc, "//item/description") {
		values := strings.Fields(n.InnerText())
		dates := xmlquery.Find(doc, "//item/pubDate")

		date, _ := time.Parse(time.RFC1123Z, dates[i].InnerText())
		sqldate := date.Format("2006-01-02")

		for j := 0; j < len(values); j += 2 {
			fmt.Println(values[j], values[j+1], sqldate)
		}
	}

}
