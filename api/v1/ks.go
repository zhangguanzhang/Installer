package v1

import (
	"Installer/api"
	"Installer/models"
	"Installer/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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
func GetKsFile(KSTemplate string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sn := c.GetHeader(SerialNumber)
		//vmware的虚机序列号会超过21，华为华三的序列号都是21位
		//最少有一个网卡的header
		if sn != "" && len(sn) <= 21 && c.GetHeader(ProMac1) != "" {
			ks := &models.KSHeader{
				SerialNumber: sn,
				Arch:         c.GetHeader(Architecture),
				System:       c.GetHeader(SystemRelease),
				NIC1Name:     strings.Split(c.GetHeader(ProMac1), " ")[0],
				NIC1MAC:      strings.Split(c.GetHeader(ProMac1), " ")[1],
				//NIC2Name:     strings.Split(c.GetHeader(ProMac2), " ")[0],
				//NIC2MAC:      strings.Split(c.GetHeader(ProMac2), " ")[1],
				//NIC3Name:     strings.Split(c.GetHeader(ProMac3), " ")[0],
				//NIC3MAC:      strings.Split(c.GetHeader(ProMac3), " ")[1],
				//NIC4Name:     strings.Split(c.GetHeader(ProMac4), " ")[0],
				//NIC4MAC:      strings.Split(c.GetHeader(ProMac4), " ")[1],
				//NIC5Name:     strings.Split(c.GetHeader(ProMac5), " ")[0],
				//NIC5MAC:      strings.Split(c.GetHeader(ProMac5), " ")[1],
				//NIC6Name:     strings.Split(c.GetHeader(ProMac6), " ")[0],
				//NIC6MAC:      strings.Split(c.GetHeader(ProMac6), " ")[1],
				//NIC7Name:     strings.Split(c.GetHeader(ProMac7), " ")[0],
				//NIC7MAC:      strings.Split(c.GetHeader(ProMac7), " ")[1],
				//NIC8Name:     strings.Split(c.GetHeader(ProMac8), " ")[0],
				//NIC8MAC:      strings.Split(c.GetHeader(ProMac8), " ")[1],
			}
			err := service.FillNicInfo(ks)
			if err != nil {
				log.Errorf("service.FillNicInfo: %s", err)
				api.Error(c, http.StatusInternalServerError, "something wrong")
				return
			}
			m, err := service.GetMachine(sn)
			if err != nil {
				api.Error(c, http.StatusInternalServerError, "something wrong")
				return
			}
			if m.SerialNumber == "" {
				//此处记录未录入库发起ks请求的机器的序列号
				log.Warnf("record not found %v", sn)
				//ks错误的话不能返回200状态码
				api.Error(c, http.StatusBadRequest, "record not found in database")
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

			t1, _ := template.ParseFiles(KSTemplate)
			if err = t1.Execute(c.Writer, info); err != nil { // 返回渲染的ks文件
				log.Errorf("template err:%v", err)
			}
			return

		}
		log.Println(c.Request.Header)
		api.Error(c, http.StatusBadRequest, "Not a kickstart request")
	}
}

//ks的%post阶段请求，更改数据库字段count+1表明安装完成
func UpdateStatusFromKs(c *gin.Context) {
	sn := c.GetHeader(SerialNumber)
	if sn == "" {
		api.Error(c, http.StatusBadRequest, "Not a kickstart %post request")
		return
	}
	if err := service.KsPostStatus(sn); err != nil {
		log.Error(err)
		api.Error(c, http.StatusForbidden, "Not an internal machine")
		return
	}
	api.Success(c, nil, "ok")
}
