package server

import (
	"e-commerce/internal/e-commerce/handler"
	"e-commerce/internal/e-commerce/repository"
	"e-commerce/internal/e-commerce/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutesEcommerce(router *gin.Engine, cleanerMap map[string]string) {

	r := repository.NewRepositoryAdapter(nil)
	s := service.NewServiceAdapter(r, cleanerMap)
	h := handler.NewHandlerAdapter(s)

	router.POST("/api/NormalizeOrder", h.NormalizeOrderHandlers)
}
