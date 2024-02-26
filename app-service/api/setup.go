package api

import (
	_ "embed"

	"github.com/gin-gonic/gin"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

const Title = "Quible app service"

//go:embed serviceDescription.md
var ServiceDescription string

func Setup(impl libAPI.ServiceAPI, router *gin.Engine, vc libAPI.VersionConfig, withOptions ...libAPI.WithOption) {
	postInit := libAPI.GetPostInit(Title, ServiceDescription)
	postInit(impl, router, vc, withOptions...)
}
