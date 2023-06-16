package main

import (
	// package
	"b47-s1/connection"
	"b47-s1/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// different between struct and interface
// struct get all data at database
// interface get single data at database

//different between Query and QueryRow
//Query get all data at table in database
//QueryRow get one table row at table in database

// struct ini sbagai membangun object/properties
type Project struct{
	Id int
	Title string
	StartDate string
	EndDate string
	StartTime time.Time
	EndTime time.Time
	Desc string
	Duration string
	CheckJs bool
	CheckNodejs bool
	CheckPhp bool
	CheckPython bool
	File string
	Author string
}

type User struct{
	Id int
	Name string
	Email string
	Password string
}

// session if login will appear data
type SessionData struct{
	IsLogin bool
	Name string
}

var userData = SessionData{}

// var dataProject = []Project{
	// {
	// 	Title: "beda nih",
	// 	StartDate: "24/04/2023",
	// 	EndDate: "25/04/2023",
	// 	Desc: "ini adalah ayam semesta",
	// 	Duration: "1 Bulan",
	// 	CheckJs: true,
	// 	CheckNodejs: true,
	// 	CheckPhp: true,
	// 	CheckPython: true,
	// 	File: "hih",
	// },
	// {
	// 	Title: "beda nih 2",
	// 	StartDate: "24/04/2022",
	// 	EndDate: "25/04/2023",
	// 	Desc: "ini adalah ayam semesta",
	// 	Duration: "1 tahun",
	// 	CheckJs: true,
	// 	CheckNodejs: true,
	// 	CheckPhp: true,
	// 	CheckPython: true,
	// 	File: "hih",
	// },
// }



func main() {

	// connect to database
	connection.DataBaseConnect()

	e := echo.New()
	e.Static("/public", "public")
	e.Static("/uploads", "uploads")

	//to use sessions using echo
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))


	// Routing
	// GET
	e.GET("/", home)
	e.GET("/form-add-project", formAddProject)
	e.GET("/form-contact", formContact)
	e.GET("/testimonial-page", testimonPage)
	e.GET("/detail-project-page/:id", detailPage)
	e.GET("/form-project-update/:id", updateProject)

	// LOGIN
	e.GET("/form-login", loginPage)
	e.POST("/login", login)

	// SIGNUP
	e.GET("/form-signup", signUpPage)
	e.POST("/signup", signup)

	//SIGN OUT
	e.POST("/signout", signout)

	// POST
	e.POST("/add-project", middleware.UploadFile(addProject))
	e.POST("/project-delete/:id", deleteProject)
	e.POST("/form-project-update/:id", middleware.UploadFile(resUpdate))


	e.Logger.Fatal(e.Start("localhost:5100"))
}

//  home /landing page
func home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, title, start_date, end_date, duration, content, js, java, php, python, image, tb_user.name AS author FROM tb_project JOIN tb_user ON tb_project.author_id = tb_user.id ORDER BY tb_project.id DESC")

	var result []Project
	for data.Next(){
		var each = Project{}

		err := data.Scan(&each.Id, 
						&each.Title, 
						&each.StartDate, 
						&each.EndDate,
						&each.Duration, 
						&each.Desc, 
						&each.CheckJs, 
						&each.CheckNodejs, 
						&each.CheckPhp, 
						&each.CheckPython,
						&each.File,
						&each.Author,
		)
		
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"massage": err.Error()})
		}

		result = append(result, each)
	}

	sess,_ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else{
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	datas := map[string]interface{}{
		"Projects": result,
		"FlashStatus": sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"DataSession": userData,
	}
	
	delete(sess.Values, "status")
	delete(sess.Values, "message")
	sess.Save(c.Request(),c.Response())


	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message":err.Error()})
	}


	return tmpl.Execute(c.Response(), datas)
}

// Form add project
func formAddProject(c echo.Context) error {
	sess,_ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else{
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"DataSession": userData,
	}
	var tmpl, err = template.ParseFiles("views/form-add-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message":err.Error()})
	}

	return tmpl.Execute(c.Response(),projects)
}

// Form Contact
func formContact(c echo.Context) error{
	sess,_ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else{
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"DataSession": userData,
	}
	var tmpl, err = template.ParseFiles("views/form-contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

// Testimonial page
func testimonPage(c echo.Context) error{
	sess,_ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else{
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"DataSession": userData,
	}


	var tmpl, err = template.ParseFiles("views/testimonial-page.html")

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), projects)
}

// Detail page
func detailPage(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))
	sess,_ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else{
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	// data := map[string]interface{}{
	// 	"Id" : id,
	// 	"Content": " Lorem ipsum dolor, sit amet consectetur adipisicing elit. Accusamus,neque repellat nam odio officia magni velit amet rerum voluptatum ut",
	// }

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT tb_project.id, title, start_date, end_date, duration, content, js, java, php, python, image, tb_user.name AS author FROM tb_project JOIN tb_user ON tb_project.author_id = tb_user.id WHERE tb_project.id=$1",id).Scan(
		&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Desc, &ProjectDetail.CheckJs, &ProjectDetail.CheckNodejs, &ProjectDetail.CheckPhp, &ProjectDetail.CheckPython, &ProjectDetail.File, &ProjectDetail.Author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}


	data := map[string]interface{}{
		"Project": ProjectDetail,
		"DataSession": userData,
	}

	var tmpl, errOut = template.ParseFiles("views/detail-project-page.html")

	if errOut != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errOut.Error()})
	}

	return tmpl.Execute(c.Response(), data)
	
}

// Duration Date Function
func diffDate(startDate string, endDate string)string{
	var duration string

	layoutDate := "2006-01-02"
	startDateParse, err := time.Parse(layoutDate, startDate )
	if err != nil {
		fmt.Println("Parsing Date Error", err)
	}
	endDateParse, err := time.Parse(layoutDate, endDate)
	if err != nil {
		fmt.Println("Parsing Date Error", err)
	}

	diff := endDateParse.Sub(startDateParse).Hours()

	days := int(diff / 24)
	weeks := days / 7
	months := weeks / 4
	years := months / 12

	if days >= 0 {
		duration = strconv.Itoa(days) + " days"
	}
	if weeks > 0 {
		duration = strconv.Itoa(weeks) + " weeks"
	}
	if months > 0 {
		duration = strconv.Itoa(months) + " Months"
	}
	if years > 0 {
		duration = strconv.Itoa(years) + " years"
	}

	return duration
}

// Add Project Page
func addProject(c echo.Context)error{
	title := c.FormValue("input-title")
	startDate := c.FormValue("input-start-date")
	endDate := c.FormValue("input-end-date")
	duration := diffDate(startDate, endDate)
	desc := c.FormValue("input-desc")
	checkJs := c.FormValue("checkJs")
	checkNodejs := c.FormValue("checkNodejs")
	checkPhp := c.FormValue("checkPhp")
	checkPython := c.FormValue("checkPython")
	js := checkJs !=""
	nodeJs := checkNodejs !=""
	php := checkPhp !=""
	python := checkPython !=""
	file := c.Get("dataFile").(string) 
	sess,_ := session.Get("session", c)
	author := sess.Values["id"].(int)

	_, err:= connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (title, start_date, end_date, duration, content, js, java, php, python, image, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,$11)", title, startDate, endDate, duration, desc , js, nodeJs, php, python, file, author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error add project": err.Error()})
	}


	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Add Project Page GET
func updateProject(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, duration, content, js, java, php, python, image FROM tb_project WHERE id=$1",id).Scan(
		&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Desc, &ProjectDetail.CheckJs, &ProjectDetail.CheckNodejs, &ProjectDetail.CheckPhp, &ProjectDetail.CheckPython, &ProjectDetail.File)

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errOut = template.ParseFiles("views/form-edit-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errOut.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

// Add Project Page POST
func resUpdate(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Index :", id)

	title := c.FormValue("input-title")
	startDate := c.FormValue("input-start-date")
	endDate := c.FormValue("input-end-date")
	duration := diffDate(startDate, endDate)
	desc := c.FormValue("input-desc")
	checkJs := c.FormValue("checkJs")
	checkNodejs := c.FormValue("checkNodejs")
	checkPhp := c.FormValue("checkPhp")
	checkPython := c.FormValue("checkPython")
	file :=  c.Get("dataFile").(string) 

	sess,_ := session.Get("session", c)
	author := sess.Values["id"].(int)

	// konversi cekbox string to boolean
	js := checkJs !=""
	nodeJs := checkNodejs !=""
	php := checkPhp !=""
	python := checkPython !=""

	_,err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title=$1, start_date=$2, end_date=$3, duration=$4, content=$5, js=$6, java=$7, php=$8, python=$9, image=$10, author_id=$11 WHERE id=$12", title, startDate, endDate, duration, desc, js, nodeJs, php, python, file, author, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Delete Project Function
func deleteProject(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Id :", id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}


// Login Page
func loginPage(c echo.Context)error{
	sess,_ := session.Get("session", c)


	flash := map[string]interface{}{
		"FlashStatus": sess.Values["status"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "status")
	delete(sess.Values, "message")
	sess.Save(c.Request(), c.Response())


	var tmpl, err = template.ParseFiles("views/form-login.html")
	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

// Sign Up Page
func signUpPage(c echo.Context)error{
	var tmpl,err = template.ParseFiles("views/form-signup.html")
	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"massage":err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

// POST data Sign up to database
func signup(c echo.Context)error{
	//make sure request data is form data format not from any format like json, xml ,etc
	err := c.Request().ParseForm()
	if err != nil{
		log.Fatal(err)
	}
	name := c.FormValue("input-name")
	email := c.FormValue("input-email")
	password := c.FormValue("input-password")
	
	// to hashing password 
	hash, _ := bcrypt.GenerateFromPassword([]byte(password),10)

	_,err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1,$2,$3)", name, email, hash)

	if err != nil{
		redirectWithMessage(c, "Signup Failed, Please Try Again", false, "/form-signup")
	}

	return redirectWithMessage(c, "Signup Success", true, "/form-login")
}

func login(c echo.Context)error{
	err := c.Request().ParseForm()
	if err != nil{
		log.Fatal(err)
	}
	email := c.FormValue("input-email")
	password := c.FormValue("input-password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(),"SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.Id,&user.Name,&user.Email,&user.Password)
	if err != nil{
		return redirectWithMessage(c, "Incorrect Email, Please Try Again", false, "/form-login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
	if err != nil{
		return redirectWithMessage(c, "Incorrect Password, Please Try Again", false, "form-login")
	}

	sess, _ := session.Get("session", c)
	// for duration per session ex 3 hour
	sess.Options.MaxAge = 10800
	sess.Values["message"] = "Sign In Succes!"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["id"] =  user.Id
	sess.Values["isLogin"] = true
	sess.Save(c.Request(),c.Response())
	
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func signout(c echo.Context)error{
	sess,_:= session.Get("session",c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}


func redirectWithMessage(c echo.Context, message string, status bool,  path string) error {
	sess,_:= session.Get("session",c)
	sess.Values["message"]= message
	sess.Values["status"] = status
	sess.Save(c.Request(),c.Response())
	return c.Redirect(http.StatusMovedPermanently, path)
}