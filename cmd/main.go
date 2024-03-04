package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var client *route53domains.Client

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func configureAWS() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client = route53domains.NewFromConfig(cfg, func(o *route53domains.Options) {
		o.Region = "us-east-1"
	})
}

func checkAvailability(domain *string) string {
	output, err := client.CheckDomainAvailability(context.TODO(), &route53domains.CheckDomainAvailabilityInput{
		DomainName: domain,
	})
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", output.Availability)
}

func loadTemplates() *template.Template {
	tmpls, err := template.New("").ParseGlob("web/template/*.html")
	if err != nil {
		log.Fatalf("couldn't initialize templates: %v", err)
	}
	return tmpls
}

func main() {
	loadEnv()
	configureAWS()

	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: loadTemplates(),
	}

	e.Use(middleware.Logger())
	e.Static("/web/static/dist", "dist")
	e.Static("/web/static/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", nil)
	})

	e.POST("/available", func(c echo.Context) error {
		domain := c.FormValue("domain")
		available := checkAvailability(&domain)

		return c.HTML(200, fmt.Sprintf("<div>%s</div>", available))
	})

	e.Logger.Fatal(e.Start(":3000"))
}
