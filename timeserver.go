// Chance O'Day
// CSS 490 - Tactical SWE
// 1/19/2015
// Assignment 2

// This package acts as a primitive web server hosted on the user's local
// machine.  It displays the current time when "/time" is appended to the
// localhost address.  It's default port is 8080, but allows the server
// to be hosted on a non-default port when the flag "-port" is included.
// This package also includes a "-V" flag, which outputs to the console
// the current package version.  If any other address is used except 
// "/time", an error message is displayed in the browser.  If the
// target port is already in use, the package will output an error
// message to the console.  When time is displayed, it is formatted
// as Hour:Minute:Second (AM/PM).

//CITE:: http://stackoverflow.com/questions/12756782/go-http-post-and-use-cookies
package main

import (
	"fmt"
	"time"
	"flag"
	"net/http"
	"os/exec"
)

// Parses flags and starts the server.  Prints an error message
// if the target port is already in use.
func main() {
	port := flag.Int("port", 8080, "Sets the server port")
	version := flag.Bool("V", false, "Program version number")
	flag.Parse()

	cookieJar := make(map[string]string)
	
	portString := fmt.Sprintf(":%v", *port)

	if *version {
		fmt.Println("Current Program Version: 1.0")
		return
	}

	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/login?name=name", formHandler)
    http.HandleFunc("/time", timeHandler)
    error := http.ListenAndServe(portString , nil)

    fmt.Printf("Server already in use.  Error %v", error)
}

// General handler.  Only used when "/time" is appended to the web address
func timeHandler(response http.ResponseWriter, request *http.Request) {
    s := currentTime()
    fmt.Fprintf(response, `
    	<html>
		<head>
		<style>
		p {font-size: xx-large}
		span.time {color: red}
		</style>
		</head>
		<body>
		<p>The time is now <span class="time">%v</span>.</p>
		</body>
		</html>
    	`, s)
}

// Handler for all cases where the web address provided is not "/time"
func formHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, `
	<html>
	<body>
	<form action="login">
	  What is your name, Earthling?
	  <input type="text" name="name" size="50">
	  <input type="submit">
	</form>
	</p>
	</body>
	</html>
		`)
	request.ParseForm()
	formName := request.FormValue("name")
	cookieUUID, _ := exec.Command("uuidgen").Output()
	uuidString := fmt.Sprintf("%v", cookieUUID)

	cookie := &http.Cookie{Name:"COOKIE", Value:uuidString, Expires:time.Now().Add(356*24*time.Hour), HttpOnly:true}
	http.SetCookie(response, cookie)

	cookieJar[uuidString] = formName

	http.Redirect(response, request, "/", http.StatusFound)

}

func homePageHandler(response http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("COOKIE")

	if len(cookieJar) > 0 && cookieJar[cookie.Value] != "" {
		fmt.Fprintf(response, "Greetings, %v", cookieJar[cookie.Value])
	} else {
		http.Redirect(response, request, "/login?name=name", http.StatusFound)
	}
}

// Returns the current time in the format Hour:Minute:Second (AM/PM)
func currentTime() string {
	const layout = "03:04:05 PM"
	t := time.Now()
	s := t.Format(layout)
	return s;
}