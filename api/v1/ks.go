package v1

import (
	"Installer/api"
	"Installer/models"
	"Installer/service"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const (
	Architecture  = "X-Anaconda-Architecture"
	SystemRelease = "X-Anaconda-System-Release"
	ProMac1       = "X-Rhn-Provisioning-Mac-0"
	ProMac2       = "X-Rhn-Provisioning-Mac-1"
	ProMac3       = "X-Rhn-Provisioning-Mac-2"
	ProMac4       = "X-Rhn-Provisioning-Mac-3"
	ProMac5       = "X-Rhn-Provisioning-Mac-4"
	ProMac6       = "X-Rhn-Provisioning-Mac-5"
	ProMac7       = "X-Rhn-Provisioning-Mac-6"
	ProMac8       = "X-Rhn-Provisioning-Mac-7"
	SerialNumber  = "X-System-Serial-Number"
)


//物理机请求ks文件，返回渲染好的ks文件
func GetKsFile(c *gin.Context) {

	if c.GetHeader(SerialNumber) != "" {
		ks := &models.KSHeader{
			SerialNumber: c.GetHeader(SerialNumber),
			Arch:         c.GetHeader(Architecture),
			System:       c.GetHeader(SystemRelease),
			NIC1Name:     strings.Split(c.GetHeader(ProMac1), " ")[0],
			NIC1MAC:      strings.Split(c.GetHeader(ProMac1), " ")[1],
			NIC2Name:     strings.Split(c.GetHeader(ProMac2), " ")[0],
			NIC2MAC:      strings.Split(c.GetHeader(ProMac2), " ")[1],
			NIC3Name:     strings.Split(c.GetHeader(ProMac3), " ")[0],
			NIC3MAC:      strings.Split(c.GetHeader(ProMac3), " ")[1],
			NIC4Name:     strings.Split(c.GetHeader(ProMac4), " ")[0],
			NIC4MAC:      strings.Split(c.GetHeader(ProMac4), " ")[1],
			NIC5Name:     strings.Split(c.GetHeader(ProMac5), " ")[0],
			NIC5MAC:      strings.Split(c.GetHeader(ProMac5), " ")[1],
			NIC6Name:     strings.Split(c.GetHeader(ProMac6), " ")[0],
			NIC6MAC:      strings.Split(c.GetHeader(ProMac6), " ")[1],
			NIC7Name:     strings.Split(c.GetHeader(ProMac7), " ")[0],
			NIC7MAC:      strings.Split(c.GetHeader(ProMac7), " ")[1],
			NIC8Name:     strings.Split(c.GetHeader(ProMac8), " ")[0],
			NIC8MAC:      strings.Split(c.GetHeader(ProMac8), " ")[1],
		}
		err := service.FillNicInfo(ks)
		if err != nil {
			log.Printf("service.FillNicInfo: %s", err)
			api.NoResponse(c)
			return
		}
		m, err := service.SearchKSInfo(c.GetHeader(SerialNumber))
		if err != nil {
			log.Printf("service.SearchKSInfo: %s|SN: %s", err, ks.SerialNumber)
			api.NoResponse(c)
			return
		}
		info := &models.KSTemplateInfo{
			Hostname: m.Hostname,
			IPMIGw:   m.IPMIGW,
			IPMIIP:   m.IPMIIP,
			IPMIMask: m.IPMIMask,
			MGIP:     m.MGIP,
			MGMask:   m.MGMask,
			MGGw:     m.MGGW,
		}

		t1, _ := template.ParseFiles("templates/ks.tmpl")
		if err = t1.Execute(c.Writer, info); err != nil {
			log.Error("template err")
		}
		return

	} else {
		fmt.Println(c.Request.Header)
		api.NoResponse(c)
	}

}

//ks的%post阶段请求，更改数据库字段表明安装完成
func KsUpdate(c *gin.Context) {

	if err := service.KsPostStatus(c.GetHeader(SerialNumber)); err != nil {
		log.Error(err)
		api.NoResponse(c)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
	})
}

//机器信息的excel表格上传导入到mysql
func ExcelUpload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 上传文件至指定目录
	if err := c.SaveUploadedFile(file, file.Filename); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    fmt.Sprintf("'%s' uploaded err!", file.Filename),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	}

	if err := service.LoadToDB(file.Filename); err != nil {
		log.Error(err)
	}
	_ = os.Remove(file.Filename)
}
