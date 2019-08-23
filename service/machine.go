package service

import (
	"Installer/models"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

//对应excel里列的字段
const (
	MachineTagName = iota //设备机房标签名
	SeriaNum              //序列号
	HOSTNAME              //最终系统hostname
	IPMIIP                //ipmi ip
	IPMIMask
	IPMIGateWay
	MGIP //系统 ip/mask
	MGMask
	MGGateway
	InputMachineName //人为写入的型号
	RaidCapacity     // 系统盘raid后容量
	IPMIMac
	NIC1Mac
	NIC2Mac
	NIC3Mac
	NIC4Mac

	SheetName = "Sheet1"
)

func LoadToDB(filename string) error {

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}

	rows, err := f.GetRows(SheetName)
	if err != nil {
		return err
	}

	rows = append(rows[:0], rows[1:]...) //删掉标题栏

	for _, row := range rows {
		if err = machinesAdd(&models.Machines{
			DeviceLabel: row[MachineTagName],
			KSHeader: &models.KSHeader{
				SerialNumber: row[SeriaNum],
			},
			Hostname: row[HOSTNAME],
			IPMIIP:   row[IPMIIP],
			IPMIMask: row[IPMIMask],
			IPMIGW:   row[IPMIGateWay],
			MGIP:     row[MGIP],
			MGMask:   row[MGMask],
			MGGW:     row[MGGateway],
		}); err != nil {
			log.Error(err)
		}
	}
	return nil
}

func machinesAdd(m *models.Machines) error {

	var temp models.Machines
	err := db.Where("SerialNumber = ?", m.SerialNumber).First(&temp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			db.Create(&m)
			return nil
		}
		log.Println(err)
		return err
	}

	//if err = db.Model(&Machines{}).Where("SerialNumber = ?", m.SerialNumber).UpdateColumn(
	//	map[string]interface{}{
	//		"DeviceLabel": m.Arch,
	//		"IPMIIPMask":  m.System,
	//		"IPMIGW":      m.NIC1Name,
	//		"Hostname":    m.NIC1MAC,
	//		"MGIPMask":    m.NIC2Name,
	//		"MGGW":        m.NIC2MAC,
	//	}).Error; err != nil {
	//	return err
	//}

	return nil
}

//根据序列号返回整列
func SearchKSInfo(SerialNumber string) (*models.Machines, error) {

	var err error
	instance := models.Machines{}

	if err = db.Where("SerialNumber = ?", SerialNumber).First(&instance).Error; err != nil {
		return nil, err
	}
	if instance.ID == 0 {
		return nil, errors.New("not found")
	}
	return &instance, nil
}

//记录ks的header里的mac填充到数据库里
func FillNicInfo(m *models.KSHeader) (err error) {

	if err = db.Model(&models.Machines{}).Where("SerialNumber = ?", m.SerialNumber).UpdateColumn(
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
		}).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}

func KsPostStatus(SerialNumber string) error {
	var err error

	if err = db.Model(&models.Machines{}).Where("SerialNumber = ?", SerialNumber).UpdateColumn(
		map[string]interface{}{
			"InstallStatus": 1,
		}).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}

func Delete(SerialNumber string) error {
	var err error
	err = db.Where("SerialNumber = ?", SerialNumber).Delete(models.Machines{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}
