package handler

import (
	"e-commerce/internal/e-commerce/model"
	"e-commerce/internal/e-commerce/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	NormalizeOrderHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s service.ServicePort
}

func NewHandlerAdapter(s service.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

// Uncomment this version if client sends JSON in the format:
// {
//     "items": [
//         { "no": 1, "platformProductId": "...", ... },
//         ...
//     ]
// }
// Suitable for scalable systems and easier to extend in the future.

// func (h *handlerAdapter) NormalizeOrderHandlers(c *gin.Context) {
// 	var input model.InputOrder
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
// 		return
// 	}
// 	output, err := h.s.NormalizeOrderService(input.Items)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"status": "OK", "message": "Order normalized successfully.", "data": output})
// }

func (h *handlerAdapter) NormalizeOrderHandlers(c *gin.Context) {
	var items []model.InputOrderItem

	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	output, err := h.s.NormalizeOrderService(items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Order normalized successfully.", "data": output})
}
