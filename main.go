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
	Content      string
	StartDate    pgtype.Date
	EndDate      pgtype.Date
	Node         string
	React        string
	Next         string
	Typescript   string
	Technologies []string
	Duration     int
}

var dataProject = []Project{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date,technologies, description FROM tb_project")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.Title, &each.StartDate, &each.EndDate, pq.Array(&each.Technologies), &each.Content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	response := map[string]interface{}{
		"Project": result,
	}

	fmt.Println(result)

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

	// var title = r.PostForm.Get("projectName")
	// var content = r.PostForm.Get("Description")
	// var startDate = r.PostForm.Get("sd")
	// var endDate = r.PostForm.Get("ed")
	// var node = r.PostForm.Get("c1")
	// var react = r.PostForm.Get("c2")
	// var next = r.PostForm.Get("c3")
	// var typescript = r.PostForm.Get("c4")

	// format := "2006-01-02"
	// startDateD, _ := time.Parse(format, startDate)
	// endDateD, _ := time.Parse(format, endDate)

	// hour := endDateD.Sub(startDateD).Hours()
	// days := hour / 24
	// week := days / 7
	// month := week / 4

	// var getMonth, _ float64 = math.Modf(month)

	// var postProject = Project{
	// 	Title:      title,
	// 	Content:    content,
	// 	StartDate:  startDate,
	// 	EndDate:    endDate,
	// 	Node:       node,
	// 	React:      react,
	// 	Next:       next,
	// 	Typescript: typescript,
	// 	Duration:   int(getMonth),
	// }

	// dataProject = append(dataProject, postProject)

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

	for i, data := range dataProject {
		if id == i {
			ProjectDetail = Project{
				Title:      data.Title,
				Content:    data.Content,
				StartDate:  data.StartDate,
				EndDate:    data.EndDate,
				Node:       data.Node,
				React:      data.React,
				Next:       data.Next,
				Typescript: data.Typescript,
				Duration:   data.Duration,
			}
		}
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	dataProject = append(dataProject[:id], dataProject[id+1:]...)

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

	for i, data := range dataProject {
		if id == i {
			EditProject = Project{
				Title:      data.Title,
				Content:    data.Content,
				StartDate:  data.StartDate,
				EndDate:    data.EndDate,
				Node:       data.Node,
				React:      data.React,
				Next:       data.Next,
				Typescript: data.Typescript,
				Duration:   data.Duration,
			}
		}
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

	// id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// dataProject[id].Title = r.PostForm.Get("projectName")
	// dataProject[id].Content = r.PostForm.Get("Description")
	// dataProject[id].StartDate = r.PostForm.Get("sd")
	// dataProject[id].EndDate = r.PostForm.Get("ed")
	// dataProject[id].Node = r.PostForm.Get("c1")
	// dataProject[id].React = r.PostForm.Get("c2")
	// dataProject[id].Next = r.PostForm.Get("c3")
	// dataProject[id].Typescript = r.PostForm.Get("c4")

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}
