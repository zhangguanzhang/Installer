
# PXE

## sn

- `GET` `/api/v1/sns` 返回所有序列号和total
```json
{
    "code": 200,
    "data": {
        "results": [
            "210200A00QH19700xxx1",
            "210200A00QH19700xxx2",
            "210200A00QH19700xxx3",
            "210200A00QH19700xxx4",
            "210200A00QH19700xx11"
        ],
        "total": 5
    }
}
```

## ks

- `GET` `/api/v1/ks` header里带上信息获取kickstart文件，因为接收是文件，所以这里不是标准的接口，可以用脚本里去模拟
```bash
curl -s -H 'Accept: [*/*]' \
  -H 'User-Agent: curl/7.29.0' \
  -H 'X-Anaconda-Architecture: x86_64' \
  -H 'X-Anaconda-System-Release: CentOS Linux' \
  -H 'X-Rhn-Provisioning-Mac-0: enp61s0f0 9c:e8:95:d8:3c:cc' \
  -H 'X-Rhn-Provisioning-Mac-1: enp61s0f1 9c:e8:95:d8:3c:cd' \
  -H 'X-Rhn-Provisioning-Mac-2: enp61s0f2 9c:e8:95:d8:3c:ce' \
  -H 'X-Rhn-Provisioning-Mac-3: enp61s0f3 9c:e8:95:d8:3c:cf' \
  -H 'X-Rhn-Provisioning-Mac-4: ens1f0 3c:f5:cc:91:1f:68' \
  -H 'X-Rhn-Provisioning-Mac-5: ens1f1 3c:f5:cc:91:1f:6a' \
  -H 'X-Rhn-Provisioning-Mac-6: ens2f0 3c:f5:cc:91:1e:48' \
  -H 'X-Rhn-Provisioning-Mac-7: ens2f1 3c:f5:cc:91:1e:4a' \
  -H 'X-System-Serial-Number: 210200A00QH185002000' localhost:8080/api/v1/ks
```
每请求一次会记录一次`Count`

- `POST` `/api/v1/ks` ks的post阶段发送请求标明安装结束，下面为ks的模板里`%post`阶段上报状态写法
```bash
if [ -r /sys/class/dmi/id/product_serial ];then
 curl -X POST -H "X-System-Serial-Number:$(cat /sys/class/dmi/id/product_serial)" http://10.1.0.2:8080/api/v1/ks
fi
```

## status

`GET /api/v1/status`返回当前库里机器的状态信息, 例如：
```json
{
    "code": 200,
    "data": {
        "results": [
            {
                "SerialNumber": "210200A00QH19700xxx1",
                "Arch": "x86_64",
                "System": "CentOS Linux",
                "Count": 1,
                "InstallStatus": 1
            },
            {
                "SerialNumber": "210200A00QH19700xxx2",
                "Arch": "",
                "System": "",
                "Count": 0,
                "InstallStatus": 0
            }
        ],
        "total": 2
    }
}
```

## machine

### 查询
- `GET` `/api/v1/machines` 返回json
  - 不带url参数默认返回所有机器
  - `/api/v1/machines?list=ColumnNames`返回表的列名
```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "ColumnNames": [
      "DeviceLabel",
      "Hostname",
      "IPMIIP",
      "IPMIMask",
      "IPMIGW",
      "MGIP",
      "MGMask",
      "MGGW",
      ...
    ]
  }
}
```
  - 也可以用列名+值来查询一个或者多个 `curl -sX GET "http://localhost:8080/api/v1/machines?MGGW=10.124.24.254" -H "accept: application/json"  | jq .`

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "result": [
      {
        "DeviceLabel": "HN1-GZ2-DYDC-C05-YW-SecPool-SRV-R4900G3-001",
        "Hostname": "HN1-GZ2-SECCVK-001",
        "IPMIIP": "10.124.24.181",
        "IPMIMask": "255.255.255.0",
        "IPMIGW": "10.124.24.254",
        "MGIP": "10.124.24.129",
        "MGMask": "255.255.255.0",
        "MGGW": "10.124.24.254",
        "SerialNumber": "210200A00QH197000769",
        "CreatedAt": "2019-09-01T19:30:54+08:00",
        "UpdatedAt": "2019-09-01T19:30:54+08:00",
        "DeletedAt": null
      },
      {
        "DeviceLabel": "HN1-GZ2-DYDC-C06-YW-SecPool-SRV-R4900G3-002",
        "Hostname": "HN1-GZ2-SECCVK-002",
        "IPMIIP": "10.124.24.182",
        "IPMIMask": "255.255.255.0",
        "IPMIGW": "10.124.24.254",
        "MGIP": "10.124.24.130",
        "MGMask": "255.255.255.0",
        "MGGW": "10.124.24.254",
        "SerialNumber": "210200A00QH197000765",
        "CreatedAt": "2019-09-01T19:30:54+08:00",
        "UpdatedAt": "2019-09-01T19:30:54+08:00",
        "DeletedAt": null
      },
      ...
    ],
    "total": 4
  }
}
```

- `GET` `/api/v1/machine/:sn` 序列号查找机器，返回json,，例如`/api/v1/machine/210200A00QH197000769`
```json
{
    "code": 200,
    "message": "ok",
    "data": {
        "DeviceLabel": "HN1-GZ2-DYDC-C05-YW-SecPool-SRV-R4900G3-001",
        "Hostname": "HN1-GZ2-SECCVK-001",
        "IPMIIP": "10.124.24.181",
        "IPMIMask": "255.255.255.0",
        "IPMIGW": "10.124.24.254",
        "MGIP": "10.124.24.129",
        "MGMask": "255.255.255.0",
        "MGGW": "10.124.24.254",
        "SerialNumber": "210200A00QH197000769",
        "CreatedAt": "2019-09-01T19:30:54+08:00",
        "UpdatedAt": "2019-09-01T19:30:54+08:00",
        "DeletedAt": null
    }
}
```

### 新加

- `POST` `/api/v1/machine` 接收和返回json,json要求为下面
```
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
```

### 更新

- `PUT` `/api/v1/machine/:sn` model的Machine里的属性名就是key名，可以发一部分更新，例如更新`210200A00QH197000769`的Arch为`x86_64`
```json
{
	"Arch": "x86_64"
}
```


### 删除

- `DELETE` `/api/v1/machine/:sn` 删除指定sn的机器

