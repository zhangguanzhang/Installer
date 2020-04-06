package service

import (
	"Installer/models"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"strings"
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
			DeviceLabel: strings.Trim(row[MachineTagName], " "),
			KSHeader: &models.KSHeader{
				SerialNumber: strings.Trim(row[SeriaNum], " "),
			},
			Hostname: strings.Trim(row[HOSTNAME], " "),
			IPMIIP:   strings.Trim(row[IPMIIP], " "),
			IPMIMask: strings.Trim(row[IPMIMask], " "),
			IPMIGW:   strings.Trim(row[IPMIGateWay], " "),
			MGIP:     strings.Trim(row[MGIP], " "),
			MGMask:   strings.Trim(row[MGMask], " "),
			MGGW:     strings.Trim(row[MGGateway], " "),
		}); err != nil {
			// 错误不为空则表明机器已经存在
			instance = append(instance, row[SeriaNum])
		}
	}
	return instance, nil
}

func GetStatus() ([]models.Status, error) {
	var (
		result []models.Status
		err    error
	)

	err = db.Table("machines").Select([]string{
		"SerialNumber",
		"Arch",
		"System",
		"Count",
		"InstallStatus",
	}).Find(&result).Error
	if err != nil {
		log.Errorf("GetAllSN err: %s", err.Error())
		return nil, err
	}

	return result, nil
}

func GetAllSN() ([]string, error) {
	var result []string

	if err := db.Model(&models.Machine{}).Pluck("SerialNumber", &result).Error; err != nil {
		log.Errorf("GetAllSN err: %s", err.Error())
		return nil, err
	}

	return result, nil
}

func GetMachines(data map[string]interface{}) ([]models.Machine, int, error) {
	var (
		count int
		err   error
	)
	var instances []models.Machine
	err = db.Model(&models.Machine{}).Where(data).Find(&instances).Count(&count).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("GetMachines %v err: %s", data, err.Error())
		return nil, 0, err
	}

	return instances, count, nil
}

func GetMachine(SN string) (*models.Machine, error) {
	var err error
	var instance models.Machine
	err = db.Model(&models.Machine{}).Where("SerialNumber = ?", SN).First(&instance).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("sn %s search err: %s", SN, err.Error())
		return nil, err
	}

	return &instance, nil
}

func AddMachine(m *models.Machine) error {
	return db.Create(m).Error
}

func UpdateMachine(m *models.Machine) error {
	var err error

	err = db.Model(&models.Machine{}).Where("SerialNumber = ?", m.SerialNumber).Updates(m).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("instance %v update err: %s", m, err.Error())
		return err
	}
	return nil
}

func DeleteMachine(sn string) (int64, error) {

	DB := db.Where("SerialNumber = ?", sn).Unscoped().Delete(&models.Machine{})

	if DB.Error != nil && !errors.Is(DB.Error, gorm.ErrRecordNotFound) {
		log.Errorf("sn %v dalete err: %s", sn, DB.Error)
		return 0, DB.Error
	}
	log.Println(DB.RowsAffected)
	return DB.RowsAffected, nil
}
