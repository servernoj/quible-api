package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

const Title = "Quible auth service"

func Setup(serviceAPI libAPI.ServiceAPI, router *gin.Engine, vc libAPI.VersionConfig, withOptions ...libAPI.WithOption) huma.API {
	postInit := libAPI.GetPostInit(Title)
	return postInit(serviceAPI, router, vc, withOptions...)
}
