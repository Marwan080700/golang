package main

import (
	"net/http"
	"strconv"
	"text/template"
	"fmt"
	"time"
	
	"github.com/labstack/echo/v4"
)

type Project struct{
	Id int
	Title string
	StartDate string
	EndDate string
	Desc string
	Duration string
	CheckJs bool
	CheckNodejs bool
	CheckPhp bool
	CheckPython bool
	// File string
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
		// File: "hih",
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
		// CheckBox: "java",
		// File: "hih",
	},
}



func main() {
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
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message":err.Error()})
	}

	projects := map[string]interface{}{
		"Projects": dataProject,
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

	for i, data := range dataProject{
		if id == i{
			ProjectDetail = Project{
				Title: data.Title,
				StartDate: data.StartDate,
				EndDate: data.EndDate,
				Desc: data.Desc,
				CheckJs: data.CheckJs,
				CheckNodejs: data.CheckNodejs,
				CheckPhp: data.CheckPhp,
				CheckPython: data.CheckPython,
				// CheckBox: checkBox,
				// File: data.file,
			}
		}
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, err = template.ParseFiles("views/detail-project-page.html")

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
	
}


// function date

func diffDate(startDate, endDate string)string{
	startTime, _ := time.Parse("2005-02-02", startDate)
	endTime, _ := time.Parse("2005-02-02", endDate)


	duraTime := int(endTime.Sub(startTime).Hours())
	duraDays := duraTime / 24
	duraWeeks := duraDays / 7
	duraMonths := duraWeeks / 4
	duraYears := duraMonths / 12

	var duration string

	if duraYears > 1 {
		duration = strconv.Itoa(duraYears) + " years"
	} else if duraYears > 0 {
		duration = strconv.Itoa(duraYears) + " year"
	} else {
		if duraMonths > 1 {
			duration = strconv.Itoa(duraMonths) + " months"
		} else if duraMonths > 0 {
			duration = strconv.Itoa(duraMonths) + " month"
		} else {
			if duraWeeks > 1 {
				duration = strconv.Itoa(duraWeeks) + " weeks"
			} else if duraWeeks > 0 {
				duration = strconv.Itoa(duraWeeks) + " week"
			} else {
				if duraDays > 1 {
					duration = strconv.Itoa(duraDays) + " days"
				} else {
					duration = strconv.Itoa(duraDays) + " day"
				}
			}
		}
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
	// file := c.FormValue("input-file") 



	var newProject = Project{
		Title: title,
		Desc: desc,
		StartDate: startDate,
		EndDate: endDate,
		Duration: duration,
		CheckJs: js,
		CheckNodejs: nodeJs,
		CheckPhp: php,
		CheckPython: python,
		// File: file,
	}

	dataProject = append(dataProject, newProject)

	fmt.Println(dataProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func updateProject(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	for i, data := range dataProject {
		if id == i {
			ProjectDetail = Project{
				Id:          id,
				Title: data.Title,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Duration:    data.Duration,
				Desc: data.Desc,
				CheckJs: data.CheckJs,
				CheckNodejs: data.CheckNodejs,
				CheckPhp: data.CheckPhp,
				CheckPython: data.CheckPython,
			}
		}
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, err = template.ParseFiles("views/form-edit-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func resUpdate(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Index :", id)

	title := c.FormValue("input-title")
	startDate := c.FormValue("input-start-date")
	endDate := c.FormValue("input-end-date")
	desc := c.FormValue("input-desc")
	checkJs := c.FormValue("checkJs")
	checkNodejs := c.FormValue("checkNodejs")
	checkPhp := c.FormValue("checkPhp")
	checkPython := c.FormValue("checkPython")

	// konversi cekbox string to boolean
	js := checkJs !=""
	nodeJs := checkNodejs !=""
	php := checkPhp !=""
	python := checkPython !=""

	var resUpdate = Project{
		Title: title,
		Desc: desc,
		StartDate: startDate,
		EndDate: endDate,
		Duration: diffDate(startDate,endDate),
		CheckJs: js,
		CheckNodejs: nodeJs,
		CheckPhp: php,
		CheckPython: python,
	}

	dataProject[id] = resUpdate

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context)error{
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Index :", id)

	dataProject = append(dataProject[:id], dataProject[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}