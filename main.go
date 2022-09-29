package main

import (
	"Day-9/connection"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"

	"github.com/lib/pq"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home)
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/post-project", postProject).Methods("POST")
	route.HandleFunc("/post-contact", postContact).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", getProject).Methods("GET")
	route.HandleFunc("/edit-post/{id}", editProject).Methods("POST")

	fmt.Println("server running on port 8000")
	http.ListenAndServe("localhost:8000", route)
}

type Project struct {
	ID           int
	Title        string
	StartDate    pgtype.Date
	EndDate      pgtype.Date
	Node         string
	React        string
	Next         string
	Typescript   string
	Technologies []string
	Content      string
	Duration     int
	Format_Date1 string
	Format_Date2 string
}

// var dataProject = []Project{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date,technologies, description  FROM tb_project")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.Title, &each.StartDate, &each.EndDate, pq.Array(&each.Technologies), &each.Content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Node = each.Technologies[0]
		each.React = each.Technologies[1]
		each.Next = each.Technologies[2]
		each.Typescript = each.Technologies[3]

		date := "2 Jan 2006"
		each.Format_Date1 = each.StartDate.Time.Format(date)
		each.Format_Date2 = each.EndDate.Time.Format(date)

		each.Duration = int(each.EndDate.Time.Month() - each.StartDate.Time.Month())

		result = append(result, each)
	}

	response := map[string]interface{}{
		"Project": result,
	}

	// fmt.Println(result)
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/Project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func postProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var title = r.PostForm.Get("projectName")
	var startDate = r.PostForm.Get("sd")
	var endDate = r.PostForm.Get("ed")
	var node = r.PostForm.Get("c1")
	var react = r.PostForm.Get("c2")
	var next = r.PostForm.Get("c3")
	var typescript = r.PostForm.Get("c4")
	var content = r.PostForm.Get("Description")

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project(name, start_date, end_date,technologies, description) VALUES ($1, $2, $3, $4, $5)", title, startDate, endDate, pq.Array([]string{node, react, next, typescript}), content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/Contact.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func postContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Nama : " + r.PostForm.Get("name"))
	fmt.Println("Email : " + r.PostForm.Get("email"))
	fmt.Println("Telepon : " + r.PostForm.Get("phone"))
	fmt.Println("Subject : " + r.PostForm.Get("subject"))
	fmt.Println("Pesan : " + r.PostForm.Get("Pesanmu"))

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name,  start_date, end_date, technologies,description FROM tb_project WHERE id=$1", id).Scan(
		&ProjectDetail.ID, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Technologies, &ProjectDetail.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	ProjectDetail.Node = ProjectDetail.Technologies[0]
	ProjectDetail.React = ProjectDetail.Technologies[1]
	ProjectDetail.Next = ProjectDetail.Technologies[2]
	ProjectDetail.Typescript = ProjectDetail.Technologies[3]

	date := "2 Jan 2006"
	ProjectDetail.Format_Date1 = ProjectDetail.StartDate.Time.Format(date)
	ProjectDetail.Format_Date2 = ProjectDetail.EndDate.Time.Format(date)

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func getProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/edit-project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var EditProject = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name,  start_date, end_date, technologies,description FROM tb_project WHERE id=$1", id).Scan(
		&EditProject.ID, &EditProject.Title, &EditProject.StartDate, &EditProject.EndDate, &EditProject.Technologies, &EditProject.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	data := map[string]interface{}{
		"Edit": EditProject,
	}

	tmpl.Execute(w, data)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var title = r.PostForm.Get("projectName")
	var startDate = r.PostForm.Get("sd")
	var endDate = r.PostForm.Get("ed")
	var node = r.PostForm.Get("c1")
	var react = r.PostForm.Get("c2")
	var next = r.PostForm.Get("c3")
	var typescript = r.PostForm.Get("c4")
	var content = r.PostForm.Get("Description")

	Update := `UPDATE tb_project SET name=$2, start_date=$3, end_date=$4, technologies=$5, description=$6 WHERE id=$1`

	_, err = connection.Conn.Exec(context.Background(), Update, id, title, startDate, endDate, pq.Array([]string{node, react, next, typescript}), content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}
