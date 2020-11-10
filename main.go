package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	db_filename = "db.txt"
)

func main()  {

	var handler Myhandler

	http.ListenAndServe(":8008", &handler)
}

type Myhandler struct {
	//empty struct
}

func (receiver Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request){

	switch r.Method {
	case "GET":
			tools, err := ioutil.ReadFile(db_filename)
			if nil != err {
				http.Error(w, "Internal Server Error", 444)
			}
			var buffer strings.Builder
	    	buffer.WriteString("<body><html>")
			buffer.WriteString("<pre>")
			_, err = (&buffer).Write(tools)
			if nil != err {
				fmt.Printf("Error With File" )
				return
			}
	 		buffer.WriteString("</pre>")
			buffer.WriteString(
				`<form method="POST" action="/">
 					<input type="text" name="tool" />
					<input type="submit" value="add tool" />
					</form>
			`)
			buffer.WriteString("</body></html>")
			_, err = io.WriteString(w, buffer.String())
			if nil != err {
				fmt.Printf("Problem Writing File" )
				return
			}
	case "POST":
			fmt.Fprintln(w, "HELLO FROM POST")
			body, err := ioutil.ReadAll(r.Body)
			href, err := url.Parse("http://kompir.com/?"+string(body))
			fmt.Print(href)
			tool := href.Query().Get("tool")
			db, err := os.OpenFile(db_filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if nil != err {
				http.Error(w,"Error", 404)
			}
			defer db.Close()
			_, err = io.WriteString(db, tool+"\n")
			fmt.Fprintf(w, "Add %q as new tool", tool)
			w.Write(body)
	default:
			http.Error(w, "Method Not Allowed", 405)
			return
	}

}


