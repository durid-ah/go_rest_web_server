# go_rest_web_server
A simple rest web server written in Go and implemented using sockets and a custom HTTP parser
The purpose of this project is to better understand http requests and how they can be handled in a REST API.

## Routes handled:
* GET Requests at `/contacts` : The api returns a list of the contacts available in the database in JSON format
* POST Requests at `/add` : Adds a new contact to the contact table
