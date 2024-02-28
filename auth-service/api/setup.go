package api

import (
	_ "embed"

	"github.com/gin-gonic/gin"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

const Title = "Quible auth service"

//go:embed serviceDescription.md
var ServiceDescription string

func Setup(serviceAPI libAPI.ServiceAPI, router *gin.Engine, vc libAPI.VersionConfig, withOptions ...libAPI.WithOption) {
	postInit := libAPI.GetPostInit(Title, ServiceDescription)
	postInit(serviceAPI, router, vc, withOptions...)
}
