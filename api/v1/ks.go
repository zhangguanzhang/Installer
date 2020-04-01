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
		if sn != "" && len(sn) <= 21 { //vmware的虚机序列号会超过21，华为华三的序列号都是21位
			ks := &models.KSHeader{
				SerialNumber: sn,
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
				api.NewResponse(c, http.StatusInternalServerError, "something wrong")
				return
			}
			m, err := service.SearchKSInfo(c.GetHeader(SerialNumber))
			if err != nil {
				log.Printf("service.SearchKSInfo: %s|SN: %s", err, ks.SerialNumber)
				api.NewResponse(c, http.StatusInternalServerError, "something wrong")
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
			if err = t1.Execute(c.Writer, info); err != nil {
				log.Errorf("template err:%v", err)
			}
			return

		} else {
			log.Println(c.Request.Header)
			api.NewResponse(c, http.StatusBadRequest, "not a kickstart request!")
		}

	}
}


//ks的%post阶段请求，更改数据库字段表明安装完成
func UpdateStatusFromKs(c *gin.Context) {

	if err := service.KsPostStatus(c.GetHeader(SerialNumber)); err != nil {
		log.Error(err)
		api.NoResponse(c)
	}
	api.NewResponse(c, http.StatusOK)
}
