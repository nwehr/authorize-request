# Authorize Request

Example:

```
package main

import (
	"net/http"
	auth "github.com/nwehr/authorize-request"
)


func main() {
	http.HandleFunc("/hello", auth.Require(hello))
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		println(err.Error())
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
```