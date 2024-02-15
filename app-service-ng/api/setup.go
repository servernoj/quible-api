package api

import (
	_ "embed"

	"github.com/gin-gonic/gin"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

const Title = "Quible app service"

//go:embed serviceDescription.md
var ServiceDescription string

func Setup[Impl libAPI.ErrorReporter](router *gin.Engine, vc libAPI.VersionConfig, withOptions ...libAPI.WithOption) {
	libSetup := libAPI.SetupFactory[Impl](Title, ServiceDescription)
	libSetup(router, vc, withOptions...)
}
