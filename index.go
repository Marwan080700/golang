package main

import (
	"b47-s1/connection"
	"net/http"
	"strconv"
	"text/template"
	"fmt"
	"time"
	"context"
	
	"github.com/labstack/echo/v4"
)

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
}

var dataProject = []Project{
	{
		Title: "beda nih",
		StartDate: "24/04/2023",
		EndDate: "25/04/2023",
		Desc: "ini adalah ayam semesta",
		Duration: "1 Bulan",
		CheckJs: true,
		CheckNodejs: true,
		CheckPhp: true,
		CheckPython: true,
		File: "hih",
	},
	{
		Title: "beda nih 2",
		StartDate: "24/04/2022",
		EndDate: "25/04/2023",
		Desc: "ini adalah ayam semesta",
		Duration: "1 tahun",
		CheckJs: true,
		CheckNodejs: true,
		CheckPhp: true,
		CheckPython: true,
		File: "hih",
	},
}



func main() {
	connection.DataBaseConnect()

	e := echo.New()
	e.Static("/public", "public")


	// Routing
	// GET
	e.GET("/", home)
	e.GET("/form-add-project", formAddProject)
	e.GET("/form-contact", formContact)
	e.GET("/testimonial-page", testimonPage)
	e.GET("/detail-project-page/:id", detailPage)
	e.GET("/form-project-update/:id", updateProject)

	// POST
	e.POST("/add-project", addProject)
	e.POST("/project-delete/:id", deleteProject)
	e.POST("/form-project-update/:id", resUpdate)


	e.Logger.Fatal(e.Start("localhost:5100"))
}

  
func home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, duration, content, js, java, php, python, image FROM tb_project")

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
		)
		
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"massage": err.Error()})
		}

		result = append(result, each)
	}


	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message":err.Error()})
	}

	projects := map[string]interface{}{
		"Projects": result,
	}


	return tmpl.Execute(c.Response(), projects)
}

func formAddProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/form-add-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message":err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func formContact(c echo.Context) error{
	var tmpl, err = template.ParseFiles("views/form-contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func testimonPage(c echo.Context) error{
	var tmpl, err = template.ParseFiles("views/testimonial-page.html")

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(),nil)
}

func detailPage(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))
	// data := map[string]interface{}{
	// 	"Id" : id,
	// 	"Content": " Lorem ipsum dolor, sit amet consectetur adipisicing elit. Accusamus,neque repellat nam odio officia magni velit amet rerum voluptatum ut",
	// }

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, duration, content, js, java, php, python, image FROM tb_project WHERE id=$1",id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Desc, &ProjectDetail.CheckJs, &ProjectDetail.CheckNodejs, &ProjectDetail.CheckPhp, &ProjectDetail.CheckPython, &ProjectDetail.File)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	StartTime, _ := time.Parse("2006-01-02", ProjectDetail.StartDate)
	EndTime, _ := time.Parse("2006-01-02", ProjectDetail.EndDate)
	ProjectDetail.StartDate = StartTime.Format("2 January 2006")
	ProjectDetail.EndDate = EndTime.Format("2 January 2006")
	fmt.Println(ProjectDetail.StartDate, ProjectDetail.EndDate)

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errOut = template.ParseFiles("views/detail-project-page.html")

	if errOut != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errOut.Error()})
	}

	return tmpl.Execute(c.Response(), data)
	
}


// function date

func diffDate(startDate string, endDate string)string{
	var duration string

	dateLayout := "2006-01-02"
	startDateParse, err := time.Parse(dateLayout, startDate )
	if err != nil {
		fmt.Println("Parsing Date Error", err)
	}
	endDateParse, err := time.Parse(dateLayout, endDate)
	if err != nil {
		fmt.Println("Parsing Date Error", err)
	}

	difference := endDateParse.Sub(startDateParse).Hours()

	day := int(difference / 24)
	week := day / 7
	month := week / 4
	year := month / 12

	if day >= 0 {
		duration = strconv.Itoa(day) + " days"
	}
	if week > 0 {
		duration = strconv.Itoa(week) + " weeks"
	}
	if month > 0 {
		duration = strconv.Itoa(month) + " Months"
	}
	if year > 0 {
		duration = strconv.Itoa(year) + " years"
	}

	return duration
}

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
	file := c.FormValue("input-file") 

	_, err:= connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (title, start_date, end_date, duration, content, js, java, php, python, image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", title, startDate, endDate, duration, desc , js, nodeJs, php, python, file)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error add project": err.Error()})
	}


	return c.Redirect(http.StatusMovedPermanently, "/")
}

func updateProject(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, duration, content, js, java, php, python, image FROM tb_project WHERE id=$1",id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Desc, &ProjectDetail.CheckJs, &ProjectDetail.CheckNodejs, &ProjectDetail.CheckPhp, &ProjectDetail.CheckPython, &ProjectDetail.File)

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errOut = template.ParseFiles("views/form-edit-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errOut.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

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
	file := c.FormValue("input-file") 

	// konversi cekbox string to boolean
	js := checkJs !=""
	nodeJs := checkNodejs !=""
	php := checkPhp !=""
	python := checkPython !=""

	_,err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title=$1, start_date=$2, end_date=$3, duration=$4, content=$5, js=$6, java=$7, php=$8, python=$9, image=$10", title, startDate, endDate, duration, desc, js, nodeJs, php, python, file)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Id :", id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}