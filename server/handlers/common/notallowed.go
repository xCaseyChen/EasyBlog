package common

import (
	"fmt"
	"net/http"
)

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "自定义 405 页面\n")
}
