package main

import (
	"flag"
	"log"
	"net/http"
)

const html = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1>Preflight CORS</h1>
		<a href="http://localhost:9000/v1/healthcheck">/v1/healthcheck</a> 
		<a href="http://localhost:9000/v1/availableitems">/v1/availableitems</a> 
		<a href="http://localhost:9000/v1/availableitems/1">/v1/availableitems/1</a> 
		<a href="http://localhost:9000/v1/itemtypes">/v1/itemtypes</a> 
		<a href="http://localhost:9000/v1/itemtypes/1">/v1/itemtypes/1</a> 
		<a href="http://localhost:9000/debug/vars">/debug/vars</a> 
		<pre id="json"></pre>
		<script>
			document.addEventListener('DOMContentLoaded', function() {
				fetch("http://localhost:4000" + window.location.pathname + window.location.search, {
					method: "GET",
					headers: {
						'Content-Type': 'application/json',
						'Authorization': 'Bearer PLACEHOLDER'
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
	</html>`

func main() {
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()

	log.Printf("starting server on %s", *addr)

	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	log.Fatal(err)
}
