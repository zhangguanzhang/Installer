## 项目由来
设想一下: 你公司某个节点要建设，进货了200台服务器(不同规格),都是安装同样的系统，区别在于hostname和业务相关的IP以及IPMI(带外的)网络信息。接好所有的线后机柜和网络规划为
```
                                            +----------------+
                                            |  Corenet-Spine |
                                            |                |
                                            +-------+--------+
                                                    ^
                                                    |
         +---------------------+-------------------------------------------+---------------------+
         |                     ^                    |                      ^                     |
         |                     |                    |                      |                     |
         |                     |                    |                      |                     |
+--------+-------+    +--------+-------+    +-------+--------+    +--------+-------+    +--------+-------+
|    M/C/S-ACC   |    |    M/C/S-ACC   |    |    M/C/S-ACC   |    |    M/C/S-ACC   |    |    M/C/S-ACC   |
|                |    |                |    |                |    |                |    |                |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+-----------------    ------------------    ------------------    ------------------    -----------------+
```
每个机柜都有各个网络接入的交换机，考虑到最坏最省钱的情况，由于每台服务器都没有做raid(出场做raid得加服务费)。你得到现场去用笔记本一台台直连服务器的带外口网线设置好带外(IPMI)的ip后(出场配置好指定ip或者配置成DHCP得加服务费)，才能在办公室通过带外远程继续操作。重启服务器进bios或者阵列卡里去配置阵列(raid)(服务器厂商没有cli tool来远程配置raid)，最后在pxe安装完后所有系统后人为一台台上去改hostname和相关的网络设置。
但是如果所有服务器只需要开机通电，就能达到最终需求呢

## 目前做到的
上架布好线后，统计服务器的序列号和标签名，整理成一份excel表格(现场的人员统计前两列，后面可以让办公室的人填写)，例如

| 设备名 | 序列号 | hostname | IPMI IP | MASK | GW | MG IP | MASK | GW 
| :--- | --- | --- | --- | --- | --- | --- | --- | --- |
|HN1-xxx-R4900G3-001 | 210200A00QH19700xxx1 | HN1-XX-ComCVK-001 | 10.101.0.101 | 255.255.255.0 | 10.101.0.254 | 10.102.0.1 | 255.255.255.0 | 10.102.0.254 |
|HN1-xxx-R4900G3-002 | 210200A00QH19700xxx2 | HN1-XX-ComCVK-002 | 10.101.0.102 | 255.255.255.0 | 10.101.0.254 | 10.102.0.2 | 255.255.255.0 | 10.102.0.254 |
|HN1-xxx-R4900G3-003 | 210200A00QH19700xxx3 | HN1-XX-ComCVK-003 | 10.101.0.103 | 255.255.255.0 | 10.101.0.254 | 10.102.0.3 | 255.255.255.0 | 10.102.0.254 |
|HN1-xxx-R4900G3-004 | 210200A00QH19700xxx4 | HN1-XX-ComCVK-004 | 10.101.0.104 | 255.255.255.0 | 10.101.0.254 | 10.102.0.4 | 255.255.255.0 | 10.102.0.254 |

然后(通过curl命令或者后期开发个前端页面)上传excel到后端，autoInstaller相关进程启动后，交换机配置好pxe需要的DHCP的boot-filename，服务器开机后最终就会按照表格里的设置好，全程不需要人为干预

## 技术前提需知

### 前提
- [PXE](./docs/pxe.md)
- [Linux安装启动过程(Centos为例)以及ISO的修改](./docs/linux-start.md)
- [详细讲解kickstart](./docs/kickstart.md)
- [Installer接口说明](./docs/api.md)


### 总结

- 修改grub.cfg适应pxe场景以及ks依赖的俩boot option
- 修改iso的`squashfs.img`增加阵列卡的cli tool命令，这步不一定，也可以kickstart的%pre阶段下载cli的包安装然后做阵列
- 编写后端，支持上传excel和根据来源序列号返回渲染是ks文件

## 实际部署步骤

因为pxe和请求ks的时候都用到了dhcp，所以定义了一个段`10.1.0.0/16`供安装阶段使用。最终的实际ip和安装段不重合，也就是下面的组网，核心上设置路由都能通，我自己笔记本接核心上设置静态ip:
```

                              10.1.0.2      +----------------+
                                +---+       |  Corenet-Spine |
                                |pc +------>+                |
                                +---+       +-------+--------+
                                                    ^
                                                    |
         +---------------------+-------------------------------------------+---------------------+
         |                     ^                    |                      ^                     |
         |                     |                    |                      |                     |
         |                     |                    |                      |                     |
+--------+-------+    +--------+-------+    +-------+--------+    +--------+-------+    +--------+-------+
|    M/C/S-ACC   |    |    M/C/S-ACC   |    |    M/C/S-ACC   |    |    M/C/S-ACC   |    |    M/C/S-ACC   |
|                |    |                |    |                |    |                |    |                |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+----------------+    +----------------+    +----------------+    +----------------+    +----------------+
|     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |    |     Machine    |
+-----------------    ------------------    ------------------    ------------------    -----------------+
    10.1.1.0/24           10.1.2.0/24           10.1.3.0/24           10.1.4.0/24           10.1.5.0/24
```
一个机柜是一个最小单位，所以一个机柜一个段
### 实操部署步骤:

#### 准备docker环境

准备一台centos7.6+系统，有条件单独机器也行，我是笔记本虚机桥接网口，笔记本网口接核心上

- [安装docker](https://github.com/zhangguanzhang/docker-need-to-know/blob/master/1.container-and-vm/1.2.install-docker.md)
- [安装docker-compose](https://docs.docker.com/compose/install/)

#### 编译Installer
##### 准备go
安装go，建议1.13以上版本，从官网获取下载直链，当然懂容器的话可以容器编译，这里是在Linux上非容器编译
```bash
VERSION=1.14.2
OS=linux
ARCH=amd64
wget https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz #下载多半需要梯子
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
# 配置环境变量
cat <<'EOF'>> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct
export CGO_ENABLED=0
EOF
. ~/.bashrc
```
执行`go version`有输出则往下走

##### 编译

```bash
git clone https://github.com/zhangguanzhang/Installer.git  #下载文件
```

我们这标准都是八张网卡，所以代码里也是八张的逻辑，为了避免少于8张网卡会触发我代码的panic。所以我代码目前是只留一个网卡的逻辑，自行按照自己实际去取消代码`api/v1/ks.go里的41-54行`的注释，我代码里录入mac地址是为了以后的cmdb，所以mac这部分录入不是那么重要。例如你两批机器分别是6张和8张网卡，你可以留到nic6的部分就行了。

```bash
cd Installer
go build -o docker/Installer main.go # 编译可执行文件到docker目录下
```

#### 准备相关文件
##### 准备镜像相关文件
准备目录，先进入docker目录
```bash
cd docker
mkdir -p http/centos tftp mysql http/soft
```
启动mysql，tftp，nginx
```bash
docker-compose up -d
```
准备pxe启动文件，系统先把安装的iso挂载了

- `cp -a`复制centos iso的`EFI/BOOT/`下文件到上面的`tftp/`目录
- 把iso解压到`http/centos/`目录下，复制iso里`/images/pxeboot/`的`vmlinuz`、`initrd.img`到tftp目录里
- 如果你不想改iso的话可以在`http/soft/`下放入安装过程中下载和安装运行做硬raid的命令的安装包啥的
- 更改grub.cfg(tftp里的`grub.cfg.tmpl`文件可参考)，里面的ip和你实际的pc部署ip要一致，因为安装阶段下载文件是nginx提供的，`http/centos/`是放centos iso解开的文件，8080是installer，所以是
```
inst.repo=http://10.1.0.2/centos ks=http://10.1.0.2:8080/api/v1/ks
```
不要加参数`ksdevice=xxx`指定ks请求指定网卡，除非你明确现场连线正确，centos7是会去尝试每张网卡的

##### kickstart部分
kickstart部分自己按照自己的实际写，分区，密码(ks的密码可以用`cd scripts`里运行`python2.7 grub-crypt.py --sha-512`生成)啥的。可以参考文档 https://fedoraproject.org/wiki/Anaconda/Kickstart/zh-cn
特别是ipmi(也就是bmc)的网络和阵列要自己去测试下，如果阵列的cli没做到iso里这里可以wget cli的包安装然后做阵列，例如把阵列的cli放`http/soft`下，然后`%pre`阶段从nginx下载安装
```
%pre
wget http://10.1.0.2/soft/xxx
chmod u+x xxx
xxx 创建阵列
...
...
%end
```
`ks-mini.tmpl`文件是大概的列出了下模板的变量和基础设置，两个文件都参考下
编辑kickstart文件记得`dos2unix`它，可能windows下编辑有回车

##### 启动Installer
```bash
./Installer -ks=template/ks.tmpl # 指定你自己的ks文件
```
然后再打开个终端进入docker目录往下继续,Installer提供的http接口看文件`docs/api.md`文件

##### 上传excel文件导入数据库

准备填写好excel文件后，scripts目录里有脚本，参照`machines.xlsx`的列要求写好这个excel(列的作用是在代码里的常量固定的)，然后用脚本`upload.sh`去curl模拟http上传excel到我后端(url的path为localhost:8080/api/v1/ks POST请求)，后端会把excel信息导入到mysql里，例如上传excel文件`conf.xlsx`
```bash
curl -X POST localhost:8080/api/v1/upload   \
  -H "Content-Type: multipart/form-data"  \
  -F "file=@conf.xlsx"  
```
然后进入mysql查看表
```
# 密码是docker-compose里，默认zhangguanzhang
docker exec -ti pxe_mysql mysql -u root -p
> use pxe;
> select * from machines;
```
表信息都有了的话看下Installer运行界面有错误信息没，没就往后继续走，有的话多半是值的长度(例如序列号规定是21，你22位了)或者例如ip或者唯一值的列有值重复了。可以把mysql的machine表内容清空了改好excel后再上传

##### 利用curl模拟机器请求kickstart文件

`scripts`目录里有个`header.sh`脚本，自己改下序列号模拟下真实的kickstart请求，看下返回的ks内容是否渲染了机器唯一的信息



##### 配置dhcp
这个看自己的组网情况，如果你是全部都在一个设备下面，可以自己起一个软dhcp server，我这里是交换机起dhcp server，现场每个机柜交换机都配置好后dhcp和pxe部分，bootfile和next-server要写对
```
[H3C-S5130-dhcp-pool-1] bootfile-name BOOTX64.EFI
[H3C-S5130-dhcp-pool-1] next-server 10.1.0.2
[H3C-S5130-dhcp-pool-1] tftp-server ip-address 10.1.0.2
```
因为我pc是`10.1.0.0/24`，得写个安装段16的路由 
- `Installer`机器配置好ip为`10.1.0.2`, 并配置`10.1.0.0/16`的路由`route add -net 10.1.0.0/16 gw 10.1.0.254 dev ens37`，我pc是虚拟机桥接，根据自己情况去配置
- 所有物理机开机即可，不熟悉的话可以先一台物理机测试下，可行了后续基本就都固化了

### 可扩展的设想

- 目前的dhcp是在交换机上，可以开发dhcp server，因为ipmi出场是固定ip的，有些服务器厂商的mac地址贴在机器表面了，可以根据mac地址返回不同的boot-filename
- 同样的，可以用tftp server根据请求grub.cfg时候第一次请求的是自身mac地址，可以提前导入mac地址去渲染返回不同的grub.cfg从而安装不同的操作系统
