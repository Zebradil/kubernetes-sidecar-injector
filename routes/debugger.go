package routes

import (
	"net/http"
)

/*Debug prints request payload*/
func Debug(writer http.ResponseWriter, request *http.Request) {
	readRequestBody(request)
}
