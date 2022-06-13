# Reading exchange rates from an RSS feed

##### This is a basic API written in GO, which reads the latest exchange rates form [Bank of Latvia's RSS feed][url1], stores them to a MySQL database and then displays them through two endpoints.

## Table of Contents
* [Setup](#setup)
* [Commands](#commands)
* [Endpoints](#endpoints)
* [Example usage](#example-usage)

## Setup
Download the full repository and save it on your machine.
This program by default uses ports `8080` for the web app and `3306` for the database. If these ports are taken on your machine, change `DB_PORT=3306` and `WEB_PORT=8080` inside the .env file in the root of the repository.
Open a terminal of your choice and cd into the downloaded folder.
```
$ cd <path-to-folder>/tetTask
```
Run the command 'docker compose build' and wait for the build process to end (This requires Docker being installed on your system)
```
$ docker compose build
```
Run the command 'docker compose up', this will create and start up the docker containers
```
$ docker compose up
```
Now that the containers are running, run the command 'docker exec -it currencies_api sh'
```
$ docker exec -it currencies_api sh
```
This will now launch a new terminal window where you can access the program.

## Commands

After starting and running the docker containers, and launching a command terminal as described in [Setup](#setup). You can now give commands to the program. The program accepts two commands:
```
$ go run main.go -action loadCurrencies
```
This command  will grab the latest currency exchange rates from [Bank of Latvia's RSS feed][url1] and push them to the database. 
```
$ go run main.go -action startEndpoints
```
This command will start the endpoints, so the user can access them. 

## Endpoints
There are two endpoints: /currencies and /currencies/{currency_name}
### /currencies
This endpoint will display the latest exchange rates for each currency. Results are displayed in JSON format.
### /currencies/{currency_name}
This endpoint takes the GET parameter {currency_name} and then returns all of the recorded values of said currency, also in JSON format. For example /currencies/AUD will return all historical values for the Australian dollar.

For a list of available currencies, check [here][url4]

**NOTE:** Before trying to access these endpoints, make sure to run `go run main.go -action startEndpoints`, so that the endpoints are ready. Also, if the endpoints return empty JSON, run `go run main.go -action loadCurrencies`, which will populate the database.

## Example usage
```
$ cd <path-to-folder>/tetTask
$ docker compose build
$ docker compose up
$ docker exec -it currencies_api sh
```
```
$ go run main.go -action loadCurrencies
Successfully added 93 new currencies!
$ go run main.go -action startEndpoints
[GIN-debug] Listening and serving HTTP on :8080
```
Now open [localhost:8080/currencies][url2] in your browser to view all latest currencies or [localhost:8080/currencies/AUD][url3] to view the Australian dollar exchange rates.

   [url1]: <https://www.bank.lv/vk/ecb_rss.xml>
   [url2]: <http:localhost:8080/currencies>
   [url3]: <http:localhost:8080/currencies/AUD>
   [url4]: <https://www.bank.lv/statistika/dati-statistika/valutu-kursi/aktualie>
