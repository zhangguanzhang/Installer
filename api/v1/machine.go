package v1

import (
	"Installer/api"
	"Installer/models"
	"Installer/service"
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func GetStatus(c *gin.Context) {
	data := make(map[string]interface{})
	results, err := service.GetStatus()
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	data["total"] = len(results)
	data["results"] = results
	api.Success(c, data, "")
}

// @Summary Get machines
// @Produce  json
// @Param list query string false "string enums" Enums(ColumnNames)
// @Param anyKey query string false "key name could use `?list=ColumnNames` result"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /machine [get]
func GetMachines(c *gin.Context) {

	data := make(map[string]interface{})

	listTarget := c.Query("list")
	if listTarget == `ColumnNames` {
		//返回表的所有字段
		data["ColumnNames"] = service.ReturnColumnNames()
		api.Success(c, data, "")
		return
	}

	paramMap := make(map[string]interface{}, 0)
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 && len(v[0]) != 0 { //key=value
			paramMap[k] = v[0]
		} else {
			api.Error(c, http.StatusBadRequest, "Bad Request")
			return
		}
	}

	result, total, err := service.GetMachines(paramMap)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	data["total"] = total
	data["results"] = result
	api.Success(c, data, "")

}

func GetSNS(c *gin.Context) {
	data := make(map[string]interface{})
	sns, err := service.GetAllSN()
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "server error")
		return
	}
	data["total"] = len(sns)
	data["results"] = sns
	api.Success(c, data, "")
}

//use the sn to search the machine
func GetMachine(c *gin.Context) {

	sn := c.Param("sn")
	data, err := service.GetMachine(sn)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "server error")
		return
	}
	if data.SerialNumber == "" { //没找到
		api.Fail(c, "machine not found")
		return
	}
	api.Success(c, data, "")
}

func AddMachine(c *gin.Context) {

	var data models.Machine

	if err := c.BindJSON(&data); err != nil {
		api.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := service.AddMachine(&data); err != nil {
		api.Error(c, http.StatusInternalServerError, "server error")
		log.Printf("AddWorkRecord error:%v, %v", err, data)
		return
	}

	api.Success(c, nil, "")
}

//更新字段
func UpdateMachine(c *gin.Context) {

	var data models.Machine

	sn := c.Param("sn")
	m, err := service.GetMachine(sn)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, "server error")
		return
	}
	if m.SerialNumber == "" { //没找到
		api.Fail(c, "machine not found")
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		api.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if data.IsEmpty() {
		api.Error(c, http.StatusBadRequest, "At least one correct field name is required")
		return
	}

	data.KSHeader = &models.KSHeader{SerialNumber: m.SerialNumber}

	err = service.UpdateMachine(&data)
	if err != nil {
		api.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	api.Success(c, nil, "")
}

func DeleteMachine(c *gin.Context) {

	sn := c.Param("sn")
	count, err := service.DeleteMachine(sn)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if count == 0 {
		api.Fail(c, "machine not found")
		return
	}
	api.Success(c, map[string]interface{}{
		"count": count,
	}, "")
}

//机器信息的excel表格上传导入到mysql
func UploadExcel(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 上传文件至指定目录
	err = c.SaveUploadedFile(file, file.Filename)
	if err != nil {
		log.Error(err)
		api.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	f, err := excelize.OpenFile(file.Filename)
	if err != nil {
		log.Errorf("excelize.OpenFile: %v", err)
	}

	instance, err := service.LoadToDB(f, c.DefaultQuery(`SheetName`, service.SheetName))
	if err != nil {
		log.Error(err)
		api.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	_ = os.Remove(file.Filename)
	data := make(map[string]interface{})
	if len(instance) >0 {
		data["alreadyExisted"] = instance
	}
	api.Success(c, data, "success")

}
