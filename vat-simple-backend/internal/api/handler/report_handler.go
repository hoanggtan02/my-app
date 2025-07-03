package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/service"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: s}
}

func (h *ReportHandler) GetSalesOverTime(c *gin.Context) {
	companyID, _ := c.Get("companyID")

	// Lấy query params, nếu không có thì mặc định 30 ngày gần nhất
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -29) // 30 ngày bao gồm cả hôm nay

	// Ghi đè nếu có query params
	if val, ok := c.GetQuery("start_date"); ok {
		parsedTime, err := time.Parse("2006-01-02", val)
		if err == nil {
			startDate = parsedTime
		}
	}
	if val, ok := c.GetQuery("end_date"); ok {
		parsedTime, err := time.Parse("2006-01-02", val)
		if err == nil {
			endDate = parsedTime
		}
	}

	dateFormat := "2006-01-02"
	salesData, appErr := h.reportService.GetSalesOverTime(companyID.(string), startDate.Format(dateFormat), endDate.Format(dateFormat))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, salesData)
}
