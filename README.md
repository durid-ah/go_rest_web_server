# go_rest_web_server
A simple rest web server written in Go and implemented using sockets and a custom HTTP parser
The purpose of this project is to better understand http requests and how they can be handled in a REST API.

## Routes handled:
* GET Requests at `/contacts` : The api returns a list of the contacts available in the database in JSON format
* POST Requests at `/add` : Adds a new contact to the contact table
* The web server can also return static files available in the resources folder of the type `gif`, `html`, `json`
and `jpeg`

## Set up:
* The server can be easily run using the `server_binary` file in a linux machine.
* Another option would be to compile the files using `go build` or `go run`
* When the project is run it uses `PORT: 8080` by default but it can be changed using `--port` flag
* requesting `index.html` will return a page that has a form to submit new contacts as well as a list of all the
existing contacts
