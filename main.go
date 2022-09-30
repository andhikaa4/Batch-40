package main

import (
	"Day-9/connection"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgtype"
	"golang.org/x/crypto/bcrypt"

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
	route.HandleFunc("/register", Register).Methods("POST")
	route.HandleFunc("/form-register", formRegister).Methods("GET")
	route.HandleFunc("/form-login", formLogin).Methods("GET")
	route.HandleFunc("/login", Login).Methods("POST")
	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("server running on port 8000")
	http.ListenAndServe("localhost:8000", route)
}

type SessionData struct {
	IsLogin   bool
	UserName  string
	FlashData string
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
	Author       string
	Format_Date1 string
	Format_Date2 string
	IsLogin      bool
}

var Data = SessionData{}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// var dataProject = []Project{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			// meamasukan flash message
			flashes = append(flashes, f1.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, tb_project.name, start_date, end_date, technologies, description,  tb_users.name as author FROM tb_project LEFT JOIN tb_users ON tb_project.user_id = tb_users.id ORDER BY id DESC")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ID, &each.Title, &each.StartDate, &each.EndDate, pq.Array(&each.Technologies), &each.Content, &each.Author)
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
		"DataSession": Data,
		"Project":     result,
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

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	author := session.Values["ID"].(int)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project(name, start_date, end_date,technologies, description, user_id ) VALUES ($1, $2, $3, $4, $5, $6)", title, startDate, endDate, pq.Array([]string{node, react, next, typescript}), content, author)
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

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name,  start_date, end_date, technologies,description, tb_users.name as author FROM tb_project LEFT JOIN tb_users ON tb_project.user_id = tb_users.id WHERE tb_project.id=$1", id).Scan(
		&ProjectDetail.ID, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Technologies, &ProjectDetail.Content, &ProjectDetail.Author)
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

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/Register.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("pass")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_users(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/login.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			// meamasukan flash message
			flashes = append(flashes, f1.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")
	tmpl.Execute(w, Data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("pass")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(),
		"SELECT * FROM tb_users WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {

		var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Email belum terdaftar!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {

		var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Password Salah!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)

		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["ID"] = user.ID
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 10800

	session.AddFlash("succesfull login", "message")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func logout(w http.ResponseWriter, r *http.Request) {

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/form-login", http.StatusSeeOther)
}
