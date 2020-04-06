package service

import (
	"Installer/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

//记录ks的header里的mac填充到数据库里
func FillNicInfo(m *models.KSHeader) (err error) {

	if err = db.Model(&models.Machine{}).Where("SerialNumber = ?", m.SerialNumber).UpdateColumn(
		map[string]interface{}{
			"Arch":     m.Arch,
			"System":   m.System,
			"NIC1Name": m.NIC1Name,
			"NIC1MAC":  m.NIC1MAC,
			"NIC2Name": m.NIC2Name,
			"NIC2MAC":  m.NIC2MAC,
			"NIC3Name": m.NIC3Name,
			"NIC3MAC":  m.NIC3MAC,
			"NIC4Name": m.NIC4Name,
			"NIC4MAC":  m.NIC4MAC,
			"NIC5Name": m.NIC5Name,
			"NIC5MAC":  m.NIC5MAC,
			"NIC6Name": m.NIC6Name,
			"NIC6MAC":  m.NIC6MAC,
			"NIC7Name": m.NIC7Name,
			"NIC7MAC":  m.NIC7MAC,
			"NIC8Name": m.NIC8Name,
			"NIC8MAC":  m.NIC8MAC,
			"Count":    gorm.Expr("Count + 1"),
		}).Error; err != nil {
		return fmt.Errorf("FillNicInfo err: %s", err)
	}

	return nil
}

func KsPostStatus(SerialNumber string) error {
	var err error

	if err = db.Model(&models.Machine{}).Where("SerialNumber = ?", SerialNumber).UpdateColumn(
		map[string]interface{}{
			"InstallStatus": 1,
		}).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}
