package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"online-subscription-service/internal/model"
	"online-subscription-service/internal/model/dto"
	"online-subscription-service/internal/usecase"
)

type SubscriptionHandler struct {
	service subscriptionsService
}

func NewSubscriptionHandler(service subscriptionsService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	body := c.Request.Body
	defer func() {
		_ = body.Close()
	}()

	var req *dto.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	parseDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "date is not correct"})
		return
	}
	endDate, err := time.Parse("01-2006", req.EndDate)
	if err != nil {
		endDate = parseDate.AddDate(0, 1, 0)
	}
	subModel := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   parseDate,
		EndDate:     &endDate,
	}
	sub, err := h.service.Create(c.Request.Context(), subModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}

	req = &dto.SubscriptionRequest{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate:     sub.EndDate.Format("01-2006"),
	}
	c.JSON(http.StatusCreated, req)
}

func (h *SubscriptionHandler) GetAll(c *gin.Context) {
	subModel, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: ": "service error"})
		return
	}

	subResp := make([]dto.SubscriptionRequest, 0, len(*subModel))
	for _, v := range *subModel {
		resp := &dto.SubscriptionRequest{
			ServiceName: v.ServiceName,
			Price:       v.Price,
			UserID:      v.UserID,
			StartDate:   v.StartDate.Format("01-2006"),
			EndDate:     v.EndDate.Format("01-2006"),
		}
		subResp = append(subResp, *resp)
	}
	c.JSON(http.StatusOK, subResp)
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + id})
		return
	}

	subModel, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found by id: " + id})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}

	resp := &dto.SubscriptionRequest{
		ServiceName: subModel.ServiceName,
		Price:       subModel.Price,
		UserID:      subModel.UserID,
		StartDate:   subModel.StartDate.Format("01-2006"),
		EndDate:     subModel.EndDate.Format("01-2006"),
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + id})
		return
	}
	body := c.Request.Body
	defer func() {
		_ = body.Close()
	}()

	var req dto.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	newStarDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "date is not correct"})
		return
	}
	modelSub := &model.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   newStarDate,
	}
	endDate, err := time.Parse("01-2006", req.EndDate)
	if err != nil {
		endDate = modelSub.StartDate.AddDate(0, 1, 0)
	}
	modelSub.EndDate = &endDate

	resultSub, err := h.service.Update(c.Request.Context(), modelSub)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found by id: " + id})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}

	resp := &dto.SubscriptionRequest{
		ServiceName: resultSub.ServiceName,
		Price:       resultSub.Price,
		UserID:      resultSub.UserID,
		StartDate:   resultSub.StartDate.Format("01-2006"),
		EndDate:     resultSub.EndDate.Format("01-2006"),
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + id})
		return
	}

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found by id: " + id})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *SubscriptionHandler) GetTotal(c *gin.Context) {
	var req dto.StatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters", "details": err.Error()})
		return
	}

	start, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	end, err := time.Parse("01-2006", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	total, err := h.service.GetTotalCost(c.Request.Context(), req.UserID, req.ServiceName, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate total"})
		return
	}
	c.JSON(http.StatusOK, dto.StatsResponse{Total: total})
}

func (h *SubscriptionHandler) RegisterRoute(r *gin.Engine) {
	v1 := r.Group("/api/v1/subscription")
	{
		v1.POST("", h.Create)
		v1.GET("", h.GetAll)
		v1.GET("/:id", h.GetByID)
		v1.PUT("/:id", h.Update)
		v1.DELETE("/:id", h.Delete)
		v1.GET("/total", h.GetTotal)
	}

}
