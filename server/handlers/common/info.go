package common

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func RenderInfoPage(w http.ResponseWriter, code int, message string) {
	infoTemplateName := "info"
	type Info struct {
		Code    int
		Title   string
		Message string
	}
	/*
		info:
			{http.StatusInternalServerError, text{Code},"internal server error"}
			{http.StatusNotFound, text{Code}, "page not found"}
			{http.StatusMethodNotAllowed, text{Code}, "method not allowed"}
	*/

	info := Info{
		Code:    code,
		Title:   http.StatusText(code),
		Message: message,
	}
	var htmlString strings.Builder
	err := tmpl.ExecuteTemplate(&htmlString, infoTemplateName, info)
	if err != nil {
		http.Error(w, message, code)
		log.Printf("Failed to execute template %s: %v", infoTemplateName, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprint(w, htmlString.String())
}
