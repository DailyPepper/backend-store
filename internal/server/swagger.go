package server

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupSwagger(router *gin.Engine) {
	router.GET("/openapi.yaml", serveOpenAPISpec)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/openapi.yaml"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
}

func serveOpenAPISpec(c *gin.Context) {
	openAPIPath := filepath.Join("docs", "openapi.yaml")
	content, err := ioutil.ReadFile(openAPIPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "OpenAPI spec not found"})
		return
	}
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", content)
}
