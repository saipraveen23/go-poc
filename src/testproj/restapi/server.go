package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var IndexTpl, FormTpl *template.Template

type Details struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Company  string `json:"company"`
	Location string `json:"location"`
}

type detailsHandler struct {
	store map[string]Details
}

func (d *detailsHandler) methods(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			d.getDetails(w, r)
			return
		}
	case "POST":
		{
			d.insertDetails(w, r)
			return
		}
	case "DELETE":
		{
			d.deleteDetails(w, r)
			return
		}
	}

}

func (d *detailsHandler) getDetails(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method != "GET" {
		errorMsg := "unsuported method "
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(errorMsg))
		return
	}
	details := make([]Details, len(d.store))
	i := 0
	for _, e := range d.store {
		details[i] = e
		i++
	}

	data := map[string]interface{}{
		"title":   "SampleData",
		"details": details,
	}
	fmt.Println("detauls", details)

	var w bytes.Buffer
	err := IndexTpl.ExecuteTemplate(&w, "index.gohtml", data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %s\n", w.String())
	rw.Write(w.Bytes())
}

func (d *detailsHandler) insertDetails(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorMsg := "unsuported method "
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(errorMsg))
		return
	}
	fmt.Println("insert function ")

	ageValue, error := strconv.ParseInt(r.FormValue("age")[0:], 10, 64)
	if error != nil {
		log.Println(error)
	}
	fmt.Println("a", ageValue)
	employee := Details{
		Id:       r.FormValue("id"),
		Name:     r.FormValue("name"),
		Age:      int(ageValue),
		Company:  r.FormValue("company"),
		Location: r.FormValue("location"),
	}
	fmt.Println("employee", employee)
	d.store[employee.Id] = employee
	fmt.Println("employee", d.store)
	http.Redirect(rw, r, "/getemployees", 301)

}

func (d *detailsHandler) displayInsertForm(rw http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"title": "Employee form",
	}
	var w bytes.Buffer
	err := FormTpl.ExecuteTemplate(&w, "form.gohtml", data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %s\n", w.String())
	rw.Write(w.Bytes())
}

func (d *detailsHandler) deleteDetails(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorMsg := "unsuported method "
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(errorMsg))
		return
	}
	query := r.URL.Query()
	fmt.Println("query", query)
	id, present := query["id"]
	if !present || len(id) == 0 {
		fmt.Println("filters not present")
	}
	fmt.Println("ids", id)

	delete(d.store, id[0])
	fmt.Println("store in delete", d.store)

	http.Redirect(rw, r, "/getemployees", 301)

}

func newDetailsHandler() *detailsHandler {
	return &detailsHandler{
		store: map[string]Details{},
	}
}

func initialiseTemplate() {
	FormTpl = template.Must(template.ParseGlob("form.gohtml"))
	IndexTpl = template.Must(template.ParseGlob("index.gohtml"))
}

func main() {
	initialiseTemplate()
	detailsHandler := newDetailsHandler()
	http.HandleFunc("/getemployees", detailsHandler.getDetails)
	http.HandleFunc("/insertform", detailsHandler.displayInsertForm)
	http.HandleFunc("/insertemployee", detailsHandler.insertDetails)
	http.HandleFunc("/deleteemployee", detailsHandler.deleteDetails)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
