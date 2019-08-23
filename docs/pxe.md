
# PXE

## 什么是PXE？
- PXE(Pre-boot Execution Environment，预启动执行环境)是由Intel公司开发的最新技术，工作于Client/Server的网络模式，支持工作站通过网络从远端服务器下载映像，并由此支持通过网络启动操作系统。
- 通过网络接口启动计算机，不依赖本地存储设备（如硬盘）或本地已安装的操作系统
- 在启动过程中，PXE客户端会调用网际协议(IP)、用户数据报协议(UDP)、动态主机设定协议(DHCP)、小型文件传输协议(TFTP)等网络协议。
- 终端要求服务器分配IP地址，再用TFTP（trivial file transfer protocol）或MTFTP(multicast trivial file transfer protocol)协议下载一个启动软件包到本机内存中执行，由这个启动软件包完成终端基本软件设置，从而引导预先安装在服务器中的终端操作系统。
- 严格来说，PXE 并不是一种安装方式，而是一种引导方式。进行 PXE 安装的必要条件是在要安装的计算机中必须包含一个 PXE 支持的网卡（NIC），即网卡中必须要有 PXE Client。PXE 协议可以使计算机通过网络启动。此协议分为 Client端 和 Server 端，而PXE Client则在网卡的 ROM 中。当计算机引导时，BIOS 把 PXE Client 调入内存中执行，然后由 PXE Client 将放置在远端的文件通过网络下载到本地运行。运行 PXE 协议需要设置 DHCP 服务器和 TFTP 服务器。DHCP 服务器会给 PXE Client（将要安装系统的主机）分配一个 IP 地址，由于是给 PXE Client 分配 IP 地址，所以在配置 DHCP 服务器时需要增加相应的 PXE 设置。此外，在 PXE Client 的 ROM 中，已经存在了 TFTP Client，那么它就可以通过 TFTP 协议到 TFTP Server 上下载所需的文件了。
- PXE技术与RPL技术不同之处为RPL是静态路由，PXE是动态路由。RPL是根据网卡上的ID号加上其他记录组成的一个Frame（帧）向服务器发出请求。而服务器中已有这个ID数据，匹配成功则进行远程启动。PXE则是根据服务器端收到的工件站MAC地址，使用DHCP服务为这个MAC地址指定个IP地址。每次启动可能同一台工作站有与上次启动有不同的IP，即动态分配地址。

简单讲就是计算机可以通过网络的形式去安装系统

## PXE工作过程

服务器是多张网卡，所有网卡都支持pxe，买来后开机后默认是pxe启动，会每张网卡去尝试，尝试失败了换下一张，然后ipv6的方式来一遍，如此循环

### legacy PXE启动过程

![pxe](https://raw.githubusercontent.com/zhangguanzhang/Image-Hosting/master/machineInstaller/030.jpg)

- 1、PXE Client向DHCP发送请求

> PXE Client开机后，PXE BootROM（自启动芯片）获得控制权之前执行自我测试，然后以UDP发送一个广播请求（FIND帧） ，向本网络中的DHCP服务器索取IP；

- 2、DHCP服务器提供信息 

> DHCP 收到客户端的请求，送回DHCP响应应，包括用户端的IP地址、预设通信通道、PXE文件的放置位置(一般放在一台TFTP服务器上) 以及开机映像文件；

- 3、 PXE客户端请求下载启动文件

> - 客户端收到服务器发回的响应后则会回应一个帧，以请求传送启动所需文件，并把自己的MAC地址写到服务器端的Netnames.db文件中。
> - 启动所需文件包含：pxelinux.0、pxelinux.cfg / default、vmlinuz、initrd.img等文件。


- 4、Boot Server响应客户端请求并传送文件 

> - 当服务器收到客户端的请求后，他们之间之后将有更多的信息在客户端与服务器之间作应答, 用以决定启动参数。
> - BootROM基于TFTP通讯协议从Boot Server下载启动安装程序所必须的文件(pxelinux.0、pxelinux.cfg / default)。default文件下载完成后，会根据该文件中定义的引导顺序，启动Linux安装程序的引导内核。

- 5、请求下载自动应答文件 

> 客户端通过pxelinux.cfg / default文件成功的引导Linux安装内核后，安装程序首先必须确定你通过什么安装介质来安装linux，如果是通过网络安装(NFS, FTP, HTTP)，则会在这个时候初始化网络，并定位安装源位置。接着会读取default文件中指定的自动应答文件ks.cfg所在位置，根据该位置请求下载该文件。

> **注意：** 在第2步和第5步初始化2次网络，这是由于PXE获取的是安装用的内核以及安装程序等，而安装程序要获取的是安装系统所需的二进制包以及配置文件。因此PXE模块和安装程序是相对独立的，PXE的网络配置并不能传递给安装程序，从而进行两次获取IP地址过程，但IP地址在DHCP的租期内是一样的。

- 6、客户端安装操作系统 

> 正常是无人值守安装，所以都是用kickstart无人值守。将ks.cfg文件下载回来后，通过该文件找到OS Server，并按照该文件的配置请求下载安装过程需要的软件包。 
OS Server和客户端建立连接后，将开始传输软件包，客户端将开始安装操作系统。安装完成后，将提示重新引导计算机。

### UEFI PXE启动过程

网上的文档应该都是legacy的启动过程，实际上我用物理服务器测试都不能启动，现在使用的也都是UEFI启动，最终发现大体上是一样的

- 1: PXE Client向DHCP发送请求

- 2: DHCP服务器提供信息 

- 3: PXE客户端请求下载启动文件
> - UEFI启动依赖ISO里的EFI/BOOT下的`BOOTX64.EFI`文件，所以boot-filename得设置成`BOOTX64.EFI`

- 4: 64位是加载运行`BOOTX64.EFI`，然后加载运行`grubx64.efi`,然后向tftp服务器请求grub.cfg文件(也就是启动的菜单)

- 5: 按照grub.cfg里走

## pxe请求的cfg文件规律

无论是pxelinux.cfg的default还是grub.cfg实际上tftp日志可以看到有个请求规律，查询了相关文档是这样介绍的
假设pxe一个客户端的uuid为`b8945908-d6a6-41a9-611d-74a6ab80b83d`，mac地址为`88:99:AA:BB:CC:DD`，获取到的DHCP ip为`192.168.2.91`，请求的文件顺序为:
```
mybootdir/pxelinux.cfg/b8945908-d6a6-41a9-611d-74a6ab80b83d
mybootdir/pxelinux.cfg/01-88-99-aa-bb-cc-dd
mybootdir/pxelinux.cfg/C0A8025B
mybootdir/pxelinux.cfg/C0A8025
mybootdir/pxelinux.cfg/C0A802
mybootdir/pxelinux.cfg/C0A80
mybootdir/pxelinux.cfg/C0A8
mybootdir/pxelinux.cfg/C0A
mybootdir/pxelinux.cfg/C0
mybootdir/pxelinux.cfg/C
mybootdir/pxelinux.cfg/default
```
先请求uuid的文件名，然后是网卡mac地址的文件，然后是ip转换成十六进制，找不到就减少最后一个字符，最后是default
而uefi的文件规律也差不多:
```
Aug 14 03:07:43 tftp daemon.notice in.tftpd[2124]: RRQ from 10.1.1.30 filename BOOTX64.EFI
Aug 14 03:07:44 tftp daemon.notice in.tftpd[2125]: RRQ from 10.1.1.30 filename grubx64.efi
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2126]: RRQ from 10.1.1.30 filename /grub.cfg-01-50-98-b8-1a-34-8d
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2127]: RRQ from 10.1.1.30 filename /grub.cfg-0A01011E
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2128]: RRQ from 10.1.1.30 filename /grub.cfg-0A01011
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2129]: RRQ from 10.1.1.30 filename /grub.cfg-0A0101
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2130]: RRQ from 10.1.1.30 filename /grub.cfg-0A010
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2131]: RRQ from 10.1.1.30 filename /grub.cfg-0A01
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2132]: RRQ from 10.1.1.30 filename /grub.cfg-0A0
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2133]: RRQ from 10.1.1.30 filename /grub.cfg-0A
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2134]: RRQ from 10.1.1.30 filename /grub.cfg-0
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2135]: RRQ from 10.1.1.30 filename /grub.cfg
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2136]: RRQ from 10.1.1.30 filename /grub.cfg
```

最后放下运行grub后请求相关文件的log
```
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2137]: RRQ from 10.1.1.30 filename /EFI/BOOT/x86_64-efi/command.lst
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2138]: RRQ from 10.1.1.30 filename /EFI/BOOT/x86_64-efi/fs.lst
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2139]: RRQ from 10.1.1.30 filename /EFI/BOOT/x86_64-efi/crypto.lst
Aug 14 03:07:45 tftp daemon.notice in.tftpd[2140]: RRQ from 10.1.1.30 filename /EFI/BOOT/x86_64-efi/terminal.lst
Aug 14 03:08:02 tftp daemon.notice in.tftpd[2145]: RRQ from 10.1.1.30 filename vmlinuz
Aug 14 03:08:05 tftp daemon.notice in.tftpd[2146]: RRQ from 10.1.1.30 filename initrd.img
```

## 要注意的坑

我们看到默认的grub.cfg文件是这样的
```
menuentry 'Install CentOS 7' --class fedora --class gnu-linux --class gnu --class os {
	linuxefi /images/pxeboot/vmlinuz inst.stage2=hd:LABEL=CentOS\x207\x20x86_64 quiet
	initrdefi /images/pxeboot/initrd.img
}
```
所以我们得把`/images/pxeboot/`下的`vmlinuz`和`initrd.img`复制到tftp的根目录下，然后菜单改成:
```
menuentry 'Auto Install CentOS 7' --class fedora --class gnu-linux --class gnu --class os {
	linuxefi vmlinuz inst.repo=http://10.1.0.2/centos quiet
	initrdefi initrd.img
}
```
这里要注意的是`inst.stage2=hd:LABEL=CentOS\x207\x20x86_64`这个必须得改，这个意思是寻找带LABEL:`CentOS X86_64`的介质，也就是iso，而我们用的是pxe启动，ISO的文件一般是解压放到一个http server上，我使用的时候是放nginx上，所以改成`inst.repo`

不用自己起dhcp服务的话，我是在交换机上dhcp应该加俩属性
```
[H3C-S5130-dhcp-pool-1] bootfile-name BOOTX64.EFI
[H3C-S5130-dhcp-pool-1] next-server 10.1.0.2
```
最终tftp我是改了个镜像增加了日志输出的docker镜像部署的，tftp目录文件为
```
tftp
├── BOOTX64.EFI
├── fonts
│   ├── TRANS.TBL
│   └── unicode.pf2
├── grub.cfg
├── grubx64.efi
├── initrd.img
├── MokManager.efi
├── TRANS.TBL
└── vmlinuz

1 directory, 9 files
```

参考文档：

- https://linuxgeeks.github.io/2018/01/22/162310-Kickstart%E6%97%A0%E4%BA%BA%E5%80%BC%E5%AE%88%E5%AE%89%E8%A3%85CentOS6.8/
- https://wiki.syslinux.org/wiki/index.php?title=PXELINUX