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

func Register(r *gin.RouterGroup, spec, path string) {
	if len(spec) == 0 {
		log.Fatal("error reading Swagger spec file")
	}
	path = strings.TrimRight(path, "/")
	r.GET(path, swaggerPage)
	r.GET(path+"/spec", swaggerSpec(spec))
}

// @Summary		Swagger spec
// @Description	Swagger spec in YAML format
// @Tags			docs
// @Produce		application/yaml
// @Success		200	{string}	string
// @Router			/docs/spec [get]
func swaggerSpec(spec string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("content-type", "application/yaml")
		_, _ = c.Writer.WriteString(spec)
	}
}

// @Summary		Swagger UI
// @Description	Render Swagger UI page
// @Tags			docs
// @Produce		text/html
// @Param			ui	query		string	false	"UI template"	Enums(swagger,rapidoc,redoc),	default(swagger)
// @Success		200	{string}	string
// @Router			/docs [get]
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
	_, _ = c.Writer.Write(data)
}
