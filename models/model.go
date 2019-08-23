package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//信息分两部分导入更新的，所以ks的header里暂时不写unique
type Machines struct {
	//设备标签名
	DeviceLabel string `gorm:"column:DeviceLabel;unique"`
	//系统hostname
	Hostname string `gorm:"size:64;column:Hostname;unique"`
	//带外管理ip/mask
	IPMIIP string `gorm:"size:15;column:IPMIIP;unique"`
	//带外管理ip/mask
	IPMIMask string `gorm:"size:15;column:IPMIMask"`
	//带外管理网关
	IPMIGW string `gorm:"size:15;column:IPMIGW"`
	//机器管理网ip
	MGIP string `gorm:"size:15;column:MGIP;unique"`
	//机器管理网掩码
	MGMask string `gorm:"size:15;column:MGMask"`
	//网关
	MGGW string `gorm:"size:15;column:MGGW"`
	//手动输入机器型号
	InputMachineName string `gorm:"column:InputMachineName"`
	//扫描到的型号
	ScanMachineName string `gorm:"column:ScanMachineName"`

	*ArrayBoard

	//专有网口mac地址
	IPMIMac string `gorm:"size:17;column:IPMIMac;"`

	*KSHeader
	InstallStatus uint8 `gorm:"column:InstallStatus;default:0"` //安装状态，ks里的post阶段发请求更新
	gorm.Model
}

type ArrayBoard struct {
	//阵列卡名字
	ArrayBoardName string `gorm:"column:ArrayBoardName"`
}

//ks信息是安装的时候带header上来，在人为写excel上传导入数据库之后(此时人为写的表格里没有mac地址)，如果mac地址字段加unique则excel的内容第二行就重复为''报错了
type KSHeader struct {
	//序列号,华三序列号20位
	SerialNumber string `gorm:"size:20;column:SerialNumber;unique;not null"`
	Arch         string `gorm:"size:7;column:Arch"`
	System       string `gorm:"size:20;column:System"`
	NIC1Name     string `gorm:"size:10;column:NIC1Name"`
	NIC1MAC      string `gorm:"size:17;column:NIC1MAC"`
	NIC2Name     string `gorm:"size:10;column:NIC2Name"`
	NIC2MAC      string `gorm:"size:17;column:NIC2MAC"`
	NIC3Name     string `gorm:"size:10;column:NIC3Name"`
	NIC3MAC      string `gorm:"size:17;column:NIC3MAC"`
	NIC4Name     string `gorm:"size:10;column:NIC4Name"`
	NIC4MAC      string `gorm:"size:17;column:NIC4MAC"`
	NIC5Name     string `gorm:"size:10;column:NIC5Name"`
	NIC5MAC      string `gorm:"size:17;column:NIC5MAC"`
	NIC6Name     string `gorm:"size:10;column:NIC6Name"`
	NIC6MAC      string `gorm:"size:17;column:NIC6MAC"`
	NIC7Name     string `gorm:"size:10;column:NIC7Name"`
	NIC7MAC      string `gorm:"size:17;column:NIC7MAC"`
	NIC8Name     string `gorm:"size:10;column:NIC8Name"`
	NIC8MAC      string `gorm:"size:17;column:NIC8MAC"`
	Count        uint8  `gorm:"column:Count;default:0"`
}


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
