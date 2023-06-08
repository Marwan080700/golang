package main

import (
	"net/http"
	"strconv"
	"text/template"
	
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/public", "public")


	e.GET("/", home)
	e.GET("/form-add-project", formAddProject)
	e.GET("/form-contact", formContact)
	e.GET("/testimonial-page", testimonPage)
	e.GET("/detail-project-page/:id", detailPage)
	e.POST("/add-project", addProject)


	e.Logger.Fatal(e.Start("localhost:5100"))
}

  
func home(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message":err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
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
	data := map[string]interface{}{
		"Id" : id,
		"Content": " Lorem ipsum dolor, sit amet consectetur adipisicing elit. Accusamus,neque repellat nam odio officia magni velit amet rerum voluptatum ut",
	}

	var tmpl, err = template.ParseFiles("views/detail-project-page.html")

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
	
}

func addProject(c echo.Context)error{
	title := c.FormValue("input-title")
	content := c.FormValue("input-desc")
	startDate := c.FormValue("input-start-date")
	endDate := c.FormValue("input-end-date")
	desc := c.FormValue("input-desc")
	checkBox := c.FormValue("input-check")
	file := c.FormValue("input-file")

	println("Title :" + title)
	println("Content :" + content)
	println("Start :" + startDate)
	println("End :" + endDate)
	println("Desc :" + desc)
	println("Check :" + checkBox)
	println("File :" + file)

	return c.Redirect(http.StatusMovedPermanently, "/")
}