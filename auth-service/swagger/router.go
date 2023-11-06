package swagger

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed *
var FS embed.FS

const templatesDir = "templates"

//	@Summary		Swagger spec
//	@Description	Swagger spec in YAML format
//	@Tags			docs
//	@Produce		application/yaml
//	@Success		200	{string}	string
//	@Router			/docs/spec [get]
func swaggerSpec(c *gin.Context) {
	c.Writer.Header().Add("content-type", "application/yaml")
	data, _ := FS.ReadFile("swagger.yaml")
	c.Writer.Write(data)
}

//	@Summary		Swagger UI
//	@Description	Render Swagger UI page
//	@Tags			docs
//	@Produce		text/html
//	@Param			ui	query		string	false	"UI template"	Enums(swagger,rapidoc,redoc),	default(swagger)
//	@Success		200	{string}	string
//	@Router			/docs [get]
func swaggerPage(c *gin.Context) {
	dir, _ := FS.ReadDir(templatesDir)
	names := make(map[string]struct{}, len(dir))
	for idx := range dir {
		names[dir[idx].Name()] = struct{}{}
	}
	template := "swagger.html"
	ui := c.Query("ui") + ".html"
	if _, ok := names[ui]; ok {
		template = ui
	}
	data, err := FS.ReadFile(templatesDir + "/" + template)
	fmt.Println(err)
	c.Writer.Header().Add("content-type", "text/html")
	c.Writer.Write(data)
}

func Register(r *gin.RouterGroup, path string) {
	if _, err := FS.ReadFile("swagger.yaml"); err != nil {
		log.Fatal("unable to read Swagger spec file")
	}
	path = strings.TrimRight(path, "/")
	r.GET(path, swaggerPage)
	r.GET(path+"/spec", swaggerSpec)
}
