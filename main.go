package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

//The general currency structure used by the microservice, formatted for JSON output
type Currency struct {
	ID    int8   `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Date  string `json:"date"`
}

//Global DB variable
var DB *sql.DB

//Establishes a global connection to the DB
func ConnectDB() {
	db, err := sql.Open("mysql", "user:admin@tcp(db:3306)/currencies")
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

//Based on the input flags, loads fresh data or starts the microservice
func main() {
	//Opens a DB connection
	ConnectDB()

	action := flag.String("action", "none", "Use the flag -action loadCurrencies to load fresh data or\n -action startEndpoints to start the microservice")
	flag.Parse()

	switch {
	case *action == "loadCurrencies":
		loadCurrencies()
	case *action == "startEndpoints":
		startEndpoints()
	default:
		fmt.Println("Unrecognised command, please use -action loadCurrencies or -action startEndpoints")
	}
}

//Routing logic and starting the HTTP server
func startEndpoints() {
	router := gin.Default()

	router.GET("/currencies", getLatestCurrencies)
	router.GET("/currencies/:name", getCurrency)

	router.Run(":8080")
}

//Gets the current exchange rates form an RSS feed and calls addCurrency() to submit all read values to the database
func loadCurrencies() {
	doc, err := xmlquery.LoadURL("https://www.bank.lv/vk/ecb_rss.xml")
	if err != nil {
		log.Fatal(err)
	}

	//Counter for displaying info to the user
	var count int

	//Loops through the XML data and finds all currency items
	for i, n := range xmlquery.Find(doc, "//item/description") {
		values := strings.Fields(n.InnerText())
		dates := xmlquery.Find(doc, "//item/pubDate")

		date, _ := time.Parse(time.RFC1123Z, dates[i].InnerText())
		sqldate := date.Format("2006-01-02")

		//Processes values for each item and and calls addCurrency()
		for j := 0; j < len(values); j += 2 {
			currID, err := addCurrency(Currency{
				Name:  values[j],
				Value: values[j+1],
				Date:  sqldate,
			})
			if err != nil {
				log.Fatal(err)
			}

			if currID != 0 {
				count++
			}
		}
	}

	if count > 0 {
		fmt.Printf("Successfully added %v new currencies!\n", count)
	} else {
		fmt.Println("No new currencies added! (Most likely because there are no new values in the RSS feed)")
	}

}

//Recieves a currency structure and inserts it to the database IF it is unique.
func addCurrency(curr Currency) (int64, error) {
	result, err := DB.Exec("INSERT INTO currency (name, value, date) SELECT ?, ?, ? WHERE NOT EXISTS(SELECT id FROM currency WHERE name = ? AND value = ? AND date = ?)", curr.Name, curr.Value, curr.Date, curr.Name, curr.Value, curr.Date)
	if err != nil {
		return 0, fmt.Errorf("addCurrency: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addCurrency: %v", err)
	}
	return id, nil
}

// Queries for the latest exchange rates and displays them in the endpoint
func getLatestCurrencies(c *gin.Context) {
	var currencies []Currency

	//Selects the LATEST value by date for each currency
	rows, err := DB.Query("SELECT c.* FROM currency c INNER JOIN (SELECT name, value, MAX(date) AS maxdate FROM currency GROUP BY name) grouped ON c.name = grouped.name AND c.date = grouped.maxdate")
	if err != nil {
		return
	}
	defer rows.Close()

	//For each record found, writes it to a Currency structure and adds that to the "currencies" variable
	for rows.Next() {
		var cur Currency
		if err := rows.Scan(&cur.ID, &cur.Name, &cur.Value, &cur.Date); err != nil {
			return
		}
		currencies = append(currencies, cur)
	}
	if err := rows.Err(); err != nil {
		return
	}

	//Displays the list as an JSON or a message that no records are found
	if currencies != nil {
		c.IndentedJSON(http.StatusOK, currencies)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no records found"})
	}
}

//Reads the GET parameter "name" from the URL and queries the database for all records of that currency
func getCurrency(c *gin.Context) {
	name := c.Param("name")

	var currencies []Currency

	rows, err := DB.Query("SELECT * FROM currency WHERE name = ?", name)
	if err != nil {
		return
	}
	defer rows.Close()

	//For each record found, writes it to a Currency structure and adds that to the "currencies" variable
	for rows.Next() {
		var cur Currency
		if err := rows.Scan(&cur.ID, &cur.Name, &cur.Value, &cur.Date); err != nil {
			return
		}
		currencies = append(currencies, cur)
	}
	if err := rows.Err(); err != nil {
		return
	}

	//Displays the list as an JSON or a message that no records are found
	if currencies != nil {
		c.IndentedJSON(http.StatusOK, currencies)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "currency not found"})
	}
}
