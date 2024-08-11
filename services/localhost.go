package services

import (
	"02.08.2024-L0/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	pageTemplate = `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Simple Form</title>
	</head>
	<body>
		<h1>Insert Order UID</h1>
		<form method="POST" action="/">
			<input type="text" name="inputText" value="{{.Input}}" />
			<button type="submit">Get</button>
		</form>
		<p>{{.Output}}</p>
	</body>
	</html>`
)

type PageData struct {
	Input  string
	Output string
}

type CacheHandler struct {
	Cache map[string]*models.Order
}

func StartupServ(cache map[string]*models.Order) {
	handler := CacheHandler{Cache: cache}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}

func (cacheHandler CacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		input := r.FormValue("inputText")
		var output string
		if val, ok := cacheHandler.Cache[input]; ok {
			jsonString, err := json.Marshal(val)
			if err != nil {
				log.Fatal(err)
				return
			} else {
				output = fmt.Sprintf(string(jsonString))
			}
		} else {
			output = fmt.Sprintf("Заказа с таким идентификатором не существует")
		}

		data := PageData{
			Input:  input,
			Output: output,
		}
		renderTemplate(w, data)
		return
	}

	renderTemplate(w, PageData{})
}

func renderTemplate(w http.ResponseWriter, data PageData) {
	t, err := template.New("form").Parse(pageTemplate)
	if err != nil {
		http.Error(w, "Ошибка при создании шаблона", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
	}
}
