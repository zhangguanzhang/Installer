package models

import (
	"reflect"
	"time"
)

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

//信息分两部分导入更新的，所以ks的header里暂时不写unique
type Machine struct {
	//设备标签名
	DeviceLabel string `gorm:"column:DeviceLabel;unique" json:"DeviceLabel,omitempty" binding:"required"`
	//系统hostname
	Hostname string `gorm:"size:64;column:Hostname;unique" json:"Hostname" binding:"required"`
	//带外管理ip/mask
	IPMIIP string `gorm:"size:15;column:IPMIIP;unique" json:"IPMIIP" binding:"required"`
	//带外管理ip/mask
	IPMIMask string `gorm:"size:15;column:IPMIMask" json:"IPMIMask" binding:"required"`
	//带外管理网关
	IPMIGW string `gorm:"size:15;column:IPMIGW" json:"IPMIGW" binding:"required"`
	//机器管理网ip
	MGIP string `gorm:"size:15;column:MGIP;unique" json:"MGIP" binding:"required"`
	//机器管理网掩码
	MGMask string `gorm:"size:15;column:MGMask" json:"MGMask" binding:"required"`
	//网关
	MGGW string `gorm:"size:15;column:MGGW" json:"MGGW" binding:"required"`
	//手动输入机器型号
	InputMachineName string `gorm:"column:InputMachineName;" json:"InputMachineName,omitempty"`
	//扫描到的型号
	ScanMachineName string `gorm:"column:ScanMachineName;" json:"ScanMachineName,omitempty"`

	*ArrayBoard

	//专有网口mac地址
	IPMIMac string `gorm:"size:17;column:IPMIMac;" json:"IPMIMac,omitempty"`

	*KSHeader
	InstallStatus uint8 `gorm:"column:InstallStatus;default:0" json:"InstallStatus,omitempty"` //安装状态，ks里的post阶段发请求更新
	Model
}

type ArrayBoard struct {
	//阵列卡名字
	ArrayBoardName string `gorm:"column:ArrayBoardName" json:"ArrayBoardName,omitempty"`
}

//ks信息是安装的时候带header上来，在人为写excel上传导入数据库之后(此时人为写的表格里没有mac地址)
// 避免零值重复，mac地址字段不能用unique属性
type KSHeader struct {
	//序列号,华三序列号20位
	SerialNumber string `gorm:"size:20;column:SerialNumber;unique;not null;primary_key" json:"SerialNumber" header:"X-System-Serial-Number"`
	Arch         string `gorm:"size:7;column:Arch" json:"Arch,omitempty" header:"X-Anaconda-Architecture"`
	System       string `gorm:"size:20;column:System" json:"System,omitempty" header:"X-Anaconda-System-Release"`
	NIC1Name     string `gorm:"size:10;column:NIC1Name"  json:"NIC1Name,omitempty"`
	NIC1MAC      string `gorm:"size:17;column:NIC1MAC"  json:"NIC1MAC,omitempty"`
	NIC2Name     string `gorm:"size:10;column:NIC2Name"  json:"NIC2Name,omitempty"`
	NIC2MAC      string `gorm:"size:17;column:NIC2MAC"  json:"NIC2MAC,omitempty"`
	NIC3Name     string `gorm:"size:10;column:NIC3Name"  json:"NIC3Name,omitempty"`
	NIC3MAC      string `gorm:"size:17;column:NIC3MAC"  json:"NIC3MAC,omitempty"`
	NIC4Name     string `gorm:"size:10;column:NIC4Name"  json:"NIC4Name,omitempty"`
	NIC4MAC      string `gorm:"size:17;column:NIC4MAC"  json:"NIC4MAC,omitempty"`
	NIC5Name     string `gorm:"size:10;column:NIC5Name"  json:"NIC5Name,omitempty"`
	NIC5MAC      string `gorm:"size:17;column:NIC5MAC"  json:"NIC5MAC,omitempty"`
	NIC6Name     string `gorm:"size:10;column:NIC6Name"  json:"NIC6Name,omitempty"`
	NIC6MAC      string `gorm:"size:17;column:NIC6MAC"  json:"NIC6MAC,omitempty"`
	NIC7Name     string `gorm:"size:10;column:NIC7Name"  json:"NIC7Name,omitempty"`
	NIC7MAC      string `gorm:"size:17;column:NIC7MAC"  json:"NIC7MAC,omitempty"`
	NIC8Name     string `gorm:"size:10;column:NIC8Name"  json:"NIC8Name,omitempty"`
	NIC8MAC      string `gorm:"size:17;column:NIC8MAC"  json:"NIC8MAC,omitempty"`
	//不在header里，主要是kickstart的请求计数
	Count uint8 `gorm:"column:Count;default:0"  json:"Count,omitempty"`
}

func (m *Machine) IsEmpty() bool {
	return reflect.DeepEqual(m, &Machine{})
}

//ks模板
type KSTemplateInfo struct {
	Hostname     string
	MGIP         string
	MGMask       string
	MGGw         string
	IPMIIP       string
	IPMIMask     string
	IPMIGw       string
	SerialNumber string //如果ip和mask为空则在模板里作为hostname
	RequsetHost  string
}

type Status struct {
	SerialNumber string `gorm:"column:SerialNumber" json:"SerialNumber"`
	Arch         string `gorm:"column:Arch" json:"Arch"`
	System       string `gorm:"column:System" json:"System"`
	Count        uint8  `gorm:"column:Count" json:"Count"`
	//安装状态，ks里的post阶段发请求更新
	InstallStatus uint8 `gorm:"column:InstallStatus" json:"InstallStatus"`
}
