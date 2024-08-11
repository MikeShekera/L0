package services

import (
	"encoding/json"
	"fmt"
	"github.com/MikeShekera/L0/transport"
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
	Stash *transport.AppStash
}

func StartupServ(stash *transport.AppStash) {
	handler := CacheHandler{Stash: stash}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}

func (cacheHandler CacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		input := r.FormValue("inputText")
		var output string
		if val, ok := cacheHandler.Stash.OrdersCache[input]; ok {
			jsonString, err := json.Marshal(val)
			if err != nil {
				log.Println(err)
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
