package v1

import (
	"Installer/api"
	"Installer/models"
	"Installer/service"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)


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
		api.NewResponse(c, http.StatusOK, nil, data)
	} else {

		paramMap := make(map[string]interface{}, 0)
		for k, v := range c.Request.URL.Query() {
			if len(v) == 1 && len(v[0]) != 0 {
				paramMap[k] = v[0]
			} else {
				api.NewResponse(c, http.StatusBadRequest)
				return
			}
		}
		total, _ := service.GetMachinesTotal(paramMap)

		result, err := service.GetMachines(paramMap)
		if err != nil {
			api.NewResponse(c, http.StatusInternalServerError, err, data)
			return
		}
		data["total"] = total
		data["result"] = result
		api.NewResponse(c, http.StatusOK, nil, data)
	}
}

func GetMachine(c *gin.Context) {

	sn := c.Param("sn")

	data, err := service.GetMachine(sn)
	if err != nil {
		api.NewResponse(c, http.StatusBadRequest, err)
		return
	}
	api.NewResponse(c, http.StatusOK, nil, data)
}


func AddMachine(c *gin.Context) {

	var data models.Machine

	if err := c.BindJSON(data); err != nil {
		api.NewResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := service.AddMachine(&data); err != nil {
		api.NewResponse(c, http.StatusInternalServerError, err)
		log.Printf("AddWorkRecord error:%v, %v", err, data)
		return
	}
	api.NewResponse(c, http.StatusOK)
}



//更新字段
func UpdateMachine(c *gin.Context) {

	var data models.Machine

	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		api.NewResponse(c, http.StatusBadRequest, err)
		return
	}

	if data.IsEmpty() {
		api.NewResponse(c, http.StatusBadRequest, "At least one correct field name is required")
		return
	}

	data.KSHeader = &models.KSHeader{SerialNumber:c.Param("sn")}


	err := service.UpdateMachine(&data)
	if err != nil {
		api.NewResponse(c, http.StatusBadRequest, err)
		return
	}

	api.NewResponse(c, http.StatusOK)
}


func DeleteMachine(c *gin.Context) {

	sn := c.Param("sn")
	if err := service.DeleteMachine(sn); err != nil {
		api.NewResponse(c, http.StatusBadRequest, err)
		return
	}
	api.NewResponse(c, http.StatusOK)
}


//机器信息的excel表格上传导入到mysql
func UploadExcel(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 上传文件至指定目录
	if err := c.SaveUploadedFile(file, file.Filename); err != nil {
		log.Error(err)
		api.NewResponse(c, http.StatusInternalServerError, err)
	} else {

		f, err := excelize.OpenFile(file.Filename)
		if err != nil {
			log.Errorf("excelize.OpenFile: %v", err)
		}

		instance, err := service.LoadToDB(f, c.DefaultQuery(`SheetName`, service.SheetName))
		if err != nil {
			log.Error(err)
			api.NewResponse(c, http.StatusInternalServerError, err)
		}
		_ = os.Remove(file.Filename)
		data := make(map[string]interface{})
		data["alreadyExisted"] = instance
		api.NewResponse(c, http.StatusCreated, fmt.Sprintf("'%s' uploaded success!", file.Filename), data)
	}

}
