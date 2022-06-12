package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "strings"
	_ "time"

	_ "github.com/antchfx/xmlquery"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// doc, err := xmlquery.LoadURL("https://www.bank.lv/vk/ecb_rss.xml")
	// if err != nil {
	// 	panic(err)
	// }

	// for i, n := range xmlquery.Find(doc, "//item/description") {
	// 	values := strings.Fields(n.InnerText())
	// 	dates := xmlquery.Find(doc, "//item/pubDate")

	// 	date, _ := time.Parse(time.RFC1123Z, dates[i].InnerText())
	// 	sqldate := date.Format("2006-01-02")

	// 	for j := 0; j < len(values); j += 2 {
	// 		albID, err := addCurrencies(Currency{
	// 			Name:  values[j],
	// 			Value: values[j+1],
	// 			Date:  sqldate,
	// 		})
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		fmt.Printf("ID of added currency: %v\n", albID)
	// 		fmt.Println(values[j], values[j+1], sqldate)
	// 	}
	// }

	/*currencies, err := getAllCurrencies("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Currencies found: %v\n", currencies)*/

	fmt.Println(getCurrency("ISK"))

}

type Currency struct {
	ID    int8
	Name  string
	Value string
	Date  string
}

func addCurrencies(curr Currency) (int64, error) {
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:8080)/currencies")

	result, err := db.Exec("INSERT INTO currency (name, value, date) SELECT ?, ?, ? WHERE NOT EXISTS(SELECT id FROM currency WHERE name = ? AND value = ? AND date = ?)", curr.Name, curr.Value, curr.Date, curr.Name, curr.Value, curr.Date)
	if err != nil {
		return 0, fmt.Errorf("addCurrency: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addCurrency: %v", err)
	}
	return id, nil
}

func getCurrency(name string) ([]Currency, error) {
	var currencies []Currency
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:8080)/currencies")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM currency WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("getCurrency %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cur Currency
		if err := rows.Scan(&cur.ID, &cur.Name, &cur.Value, &cur.Date); err != nil {
			return nil, fmt.Errorf("getCurrency %q: %v", name, err)
		}
		currencies = append(currencies, cur)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCurrency %q: %v", name, err)
	}
	return currencies, nil
}

// albumsByArtist queries for albums that have the specified artist name.
func getAllCurrencies(name string) ([]Currency, error) {
	// An albums slice to hold data from returned rows.
	var currencies []Currency

	db, err := sql.Open("mysql", "root:admin@tcp(localhost:8080)/currencies")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM currency")
	if err != nil {
		return nil, fmt.Errorf("getAllCurrencies %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var cur Currency
		if err := rows.Scan(&cur.ID, &cur.Name, &cur.Value, &cur.Date); err != nil {
			return nil, fmt.Errorf("getAllCurrencies %q: %v", name, err)
		}
		currencies = append(currencies, cur)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllCurrencies %q: %v", name, err)
	}
	return currencies, nil
}
