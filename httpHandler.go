package main

import(
	"strings"
	"fmt"
	"./repository"
	"encoding/json"
	"io/ioutil"
)

const (
	// CRLF : New line in response header
	CRLF = "\r\n"
	
	emptyString = ""

	// badRequest : Bad request status response code
	badRequest = "HTTP/1.1 400 Bad Request"

	// ok : Success status response code
	ok = "HTTP/1.1 200 OK"

	okNoContent = "HTTP/1.1 204 No Content"

	// notFound : Not found status response code
	notFound = "HTTP/1.1 404 Not Found"

	// contentType : The content type header line
	contentType = "Content-Type: "
)

// ParseRequest : takes in the client's request and parses it 
func ParseRequest(request string, requsetLength int) string {
	
	trimmedRequest := request[:requsetLength]

	fmt.Println()
	fmt.Println("##########################################################")
	fmt.Println("Incoming Request:")
	fmt.Println("##########################################################")
	fmt.Println(trimmedRequest)
	fmt.Println("##########################################################")

	requestArray := strings.Split(trimmedRequest, "\r\n")
	var routingArray []string

	if requestArray != nil && len(requestArray) > 0 {
		routingArray = strings.Split(requestArray[0], " ")
	}

	if len(routingArray) < 3 { return badRequest + CRLF }

	requestMethod := routingArray[0] 

	return handleRouting(
		&requestMethod,
		strings.Split(routingArray[1], "/"),
		&requestArray,
	)
}

// takes care of the routing section in the http request
func handleRouting(method *string, route []string, requestArray *[]string) string {
	
	routeLength := len(route)
	var repo *(repository.ContactRepo) = repository.GetInstance()
	
	// handle GET requests
	if (*method == "GET") {
		file := route[len(route) - 1] // Get the last item in the route array to prevent any attempt of accessing anything else on the system
		cType := getContentTypeFromRoute(file)

		// if the request is routed to /contacts it will return a json list
		// of all the contacts
		if (routeLength == 2 && route[1] == "contacts") {
			contacts := repo.GetAllContacts()
			data, _ := json.Marshal(contacts)
			return ok + CRLF + 
				"Content-Type: " + getContentTypeFromRoute(".json") +
				CRLF + CRLF +
				string(data) + CRLF
		
		// else handle any requested files
		} else if (routeLength == 2) {

			if len(cType) == 0 { return notFound + CRLF }

			dat, err := ioutil.ReadFile("./resources/" + file)
			if (err != nil) {
				fmt.Println(err)
				return notFound + CRLF
			}
			
			return ok + CRLF + 
			"Content-Type: " + cType +
			CRLF + CRLF +
			string(dat) + CRLF
		}


	
	} else if (*method == "POST") {
		
		if routeLength == 2 && route[1] == "add" && route[0] == "" {
			contentType, body := getContentTypeAndBody(requestArray)
			if contentType == "Content-Type: application/json" && len(body) > 0 {
				var newContact repository.Contact
				json.Unmarshal([]byte(body), &newContact)
				repo.InsertContact(&newContact)
				return ok + CRLF + CRLF + emptyString + CRLF
			}

			return badRequest + CRLF
		}
	}

	return notFound + CRLF
}

func getContentTypeFromRoute(fileName string) string {
	var fileType string
	for i := len(fileName) - 1; i >= 0; i-- {
		if (fileName[i] == '.') {
			break
		}

		fileType = string(fileName[i]) + fileType
	}

	fmt.Println("FileName " + fileType)

	if fileType == "html" || fileType == "htm" {
		return "text/html"
	} else if fileType == "json" { 
		return "application/json"
	} else if fileType == "gif" {
		return "image/gif"
	} else if fileType == "txt" {
		return "text/plain"
	} else if fileType == "jpg" ||
		fileType == "jpeg" ||
		fileType == "jfif" ||
		fileType == "pjpeg" ||
		fileType == "pjp" {
			return "image/jpeg"
		}

	return ""
}

func getContentTypeAndBody(requestArray *[]string) (string, string) {
	
	var body string = ""
	var bodyStarted bool = false
	var contentType string
	for _, line := range *requestArray {
		if  len(line) >= 12 && line[:12] == "Content-Type" && !bodyStarted { 
			contentType = line 
		} else if len (line) == 0 { 
			bodyStarted = true 
		} else if bodyStarted {
			body += line
		}
	}

	return contentType, body
}
