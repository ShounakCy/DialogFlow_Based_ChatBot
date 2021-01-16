package router

import (
	"github.com/gin-gonic/gin"
	h "github.com/havells/nlp/handler"
	"github.com/havells/nlp/models"
	"github.com/havells/nlp/services"
)

// handler -> app.go -> App struct
var a h.App

//SetupRouter :
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/query/service", a.VoiceReq)
	// r.POST("wa/query",a.VoiceReq)
	return r
}

//InitServices :
func InitServices(projectID, authFile, langCode string, ff models.FFSrvConf) {
	a.Service = &services.ServiceProvider{Lang: langCode, ProjectID: projectID,
		AuthJSONFilePath: authFile, FFConf: ff}
	return
}
