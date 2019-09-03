package service

import (
	"Installer/models"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jinzhu/gorm"
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

//传入文件名和SheetName导入数据库，返回已经存在的机器
func LoadToDB(f *excelize.File, SheetName string) ([]string, error) {
	var instance = make([]string, 0)

	rows, err := f.GetRows(SheetName)
	if err != nil {
		return nil, err
	}

	rows = append(rows[:0], rows[1:]...) //删掉标题栏

	for _, row := range rows {
		if err = AddMachine(&models.Machine{
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
			// 错误不为空则表明机器已经存在
			instance = append(instance, row[HOSTNAME])
		}
	}
	return instance, nil
}

func GetMachinesTotal(data map[string]interface{}) (int, error) {
	var (
		count int
		err error
	)

	if len(data) == 0 {
		err = db.Model(&models.Machine{}).Count(&count).Error
	} else {
		err = db.Model(&models.Machine{}).Where(data).Count(&count).Error
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetMachines(data map[string]interface{}) ([]models.Machine, error) {
	var err error
	var instance []models.Machine
	err = db.Where(data).Find(&instance).Error
	if err != nil {
		return nil, err
	} else {
		return instance, nil
	}
}

func GetMachine(SN string) (*models.Machine, error) {
	var err error
	var instance models.Machine
	err = db.Model(&models.Machine{}).Where("SerialNumber = ?", SN).First(&instance).Error
	if err != nil {
		return nil, err
	} else {
		return &instance, nil
	}
}


func AddMachine(m *models.Machine) error {
	return db.Create(m).Error
}


func UpdateMachine(m *models.Machine) error {
	var (
		count int
		err error
	)
	err = db.Model(&models.Machine{}).Where("SerialNumber = ?", m.SerialNumber).Count(&count).Updates(m).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}


func DeleteMachine(SerialNumber string) error {
	var (
		count int
		err error
	)
	err = db.Where("SerialNumber = ?", SerialNumber).Unscoped().Count(&count).Delete(models.Machine{}).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
