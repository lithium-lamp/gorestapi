package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ADMIN_TOKEN := os.Getenv("ADMIN_TOKEN")

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1>Preflight CORS</h1>
		<a href="http://localhost:9000/v1/healthcheck">/v1/healthcheck</a><br>
		<a href="http://localhost:9000/v1/availableitems">/v1/availableitems</a><br>
		<a href="http://localhost:9000/v1/availableitems/1">/v1/availableitems/1</a><br>
		<a href="http://localhost:9000/v1/itemtypes">/v1/itemtypes</a><br>
		<a href="http://localhost:9000/v1/itemtypes/1">/v1/itemtypes/1</a><br>
		<a href="http://localhost:9000/v1/measurements">/v1/measurements</a><br>
		<a href="http://localhost:9000/v1/measurements/1">/v1/measurements/1</a><br>
		<a href="http://localhost:9000/debug/vars">/debug/vars</a><br>
		<pre id="json"></pre>
		<script>
			document.addEventListener('DOMContentLoaded', function() {
				fetch("http://localhost:4000" + window.location.pathname + window.location.search, {
					method: "GET",
					headers: {
						'Content-Type': 'application/json',
						'Authorization': 'Bearer %s'
					}
				}).then(
					function (response) {
						response.text().then(function (text) {
							document.getElementById("json").textContent = text;
						});
					},
					function(err) {
						document.getElementById("output").innerHTML = err;
					}
				);
			});
		</script>
	</body>
	</html>`, ADMIN_TOKEN)

	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()

	log.Printf("starting server on %s", *addr)

	err = http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	log.Fatal(err)
}
