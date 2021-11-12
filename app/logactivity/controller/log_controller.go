package controller

import (
	"bytes"
	"encoding/csv"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"streetbox.id/app/logactivity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// LogController ..
type LogController struct {
	Svc logactivity.ServiceInterface
}

// GetAll godoc
// @Summary Get All Log Activity Pagination (permission = superadmin)
// @Id GetAllLog
// @Tags Log Activity
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Success 200 {object} model.ResponseSuccess "model.Pagination"
// @Router /log [get]
func (r *LogController) GetAll(c *gin.Context) {
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sorted := util.SortedBy(c.QueryArray("sort"))
	data := r.Svc.GetAll(limit, page, sorted)
	model.ResponsePagination(c, data)
	return
}

// GenerateCSV ...
func (r *LogController) GenerateCSV(c *gin.Context) {
	var data = [][]string{{"Log Date", "Activity"}}

	logList := r.Svc.GetList()
	var listData []string
	for _, value := range *logList {
		listData = append(listData, value.LogTime.Format("2006-01-02 15:04:05"), value.Activity)
		data = append(data, listData)
		listData = nil
	}
	b := &bytes.Buffer{}
	writer := csv.NewWriter(b)
	writer.WriteAll(data)
	if err := writer.Error(); err != nil {
		log.Fatalln("Error Writing csv:", err)
	}
	today := time.Now().Format("Mon 02 Jan 2006 15:04:05")
	fileName := "logactivity " + today + ".csv"

	writer.Flush()
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment;filename="+fileName)
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}
