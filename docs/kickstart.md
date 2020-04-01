## kickstart

我是写了个根据请求ks的header来返回渲染的ks文件的http server，但是下面的只是必须清楚

### boot option
grub.cfg被我们修改成这样了
```
menuentry 'Auto Install CentOS 7' --class fedora --class gnu-linux --class gnu --class os {
	linuxefi vmlinuz inst.repo=http://10.1.0.2/centos quiet
	initrdefi initrd.img
}
```
实际上我们后面的操作得依赖ks的某些属性，需要修改boot option，boot option是一些`key=value`写在linuxefi后quiet前，相关的文档地址为 https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/installation_guide/chap-anaconda-boot-options
boot option可以在安装启动后在终端里`cat /proc/cmdline`看到当前的boot option有哪些，`squashfs.img`里很多脚本都是从这个路径读取启动参数决定了行为

### 从没有人用过的`inst.ks.sendmac`和`inst.ks.sendsn`
最开始考虑的是起一个web，每个机器请求ks的时候渲染ip和设置命令，这样保证了预期目标。用go写了个http server，详细的打印了下一些信息，默认的请求ks的header为:
```
"Accept": "[*/*]"
"User-Agent": "curl/7.29.0"
"X-Anaconda-Architecture: x86_64 
"X-Anaconda-System-Release: "CentOS Linux"  #这个貌似有些版本是Centos
```
一开始是考虑扫描交换机获取机器的mac地址，假设提前录入了mac地址下，请求ks的时候根据客户端的ip扫描交换机的mac地址就找到了是哪台上来了(当时也写出了)。
但是后面查看资料发现实际上在boot option加俩选项后，请求ks的时候会带上mac地址和更详细的信息。我在查看`squashfs.img`里的一些function脚本的时候，在`/images/pxeboot/initrd.img/initrd/lib/dracut/hooks/initqueue/online/11-fetch-kickstart-net.sh` 里找到了一个逻辑
```bash
# If we're doing sendmac, we need to run after anaconda-ks-sendheaders.sh
if getargbool 0 inst.ks.sendmac kssendmac; then
    newjob=$hookdir/initqueue/settled/fetch-ks-${netif}.sh
else
    newjob=$hookdir/initqueue/fetch-ks-${netif}.sh
fi
```
我用grep查找这个`getargbool`函数是在`usr/lib/dracut/hooks/initqueue/settled/00-anaconda-ks-sendheaders.sh`里，内容为
```bash
#/bin/sh
# anaconda-ks-sendheaders.sh - set various HTTP headers for kickstarting

[ -f /tmp/.ks_sendheaders ] && return
command -v set_http_header >/dev/null || . /lib/url-lib.sh

# inst.ks.sendmac: send MAC addresses in HTTP headers
if getargbool 0 kssendmac inst.ks.sendmac; then
    ifnum=0
    for ifname in /sys/class/net/*; do
        [ -e "$ifname/address" ] || continue
        mac=$(cat $ifname/address)
        ifname=${ifname#/sys/class/net/}
        # TODO: might need to choose devices better
        if [ "$ifname" != "lo" ] && [ -n "$mac" ]; then
            # set_http_header is from url-lib.sh, sourced earlier
            set_http_header "X-RHN-Provisioning-MAC-$ifnum" "$ifname $mac"
            ifnum=$(($ifnum+1))
        fi
    done
fi

# inst.ks.sendsn: send system serial number in HTTP headers
if getargbool 0 kssendsn inst.ks.sendsn; then
    system_serial=$(cat /sys/class/dmi/id/product_serial 2>/dev/null)
    if [ -z "$system_serial" ]; then
        warn "inst.ks.sendsn: can't find system serial number"
    else
        set_http_header "X-System-Serial-Number" "$system_serial"
    fi
fi
> /tmp/.ks_sendheaders
```
可以看出boot的选项加了`inst.ks.sendmac`和`inst.ks.sendsn`后在curl请求ks文件的时候会附加header，所有网卡的mac地址和序列号。这个选项我搜了一圈，发现国内的和国外的基本很少用这个。http server打印了下加了这俩选项后的header信息

```
"Accept": "[*/*]"
"User-Agent": "curl/7.29.0"
"X-Anaconda-Architecture: x86_64 
"X-Anaconda-System-Release: "CentOS Linux" 
"X-Rhn-Provisioning-Mac-0: "enp61s0f0 9c:e8:95:d8:3c:cc"
"X-Rhn-Provisioning-Mac-1: "enp61s0f1 9c:e8:95:d8:3c:cd" 
"X-Rhn-Provisioning-Mac-2: "enp61s0f2 9c:e8:95:d8:3c:ce" 
"X-Rhn-Provisioning-Mac-3: "enp61s0f3 9c:e8:95:d8:3c:cf" 
"X-Rhn-Provisioning-Mac-4: "ens1f0 3c:f5:cc:91:1f:68" 
"X-Rhn-Provisioning-Mac-5: "ens1f1 3c:f5:cc:91:1f:6a" 
"X-Rhn-Provisioning-Mac-6: "ens2f0 3c:f5:cc:91:1e:48" 
"X-Rhn-Provisioning-Mac-7: "ens2f1 3c:f5:cc:91:1e:4a" 
"X-System-Serial-Number: "210200A00QH18500xxxx"
```
序列号我对比了下和物理机的一样，应该是个国际标准，而安装完进系统里也可以通过`cat /sys/class/dmi/id/product_serial`看序列号。改完的最终grub.cfg为:
```
menuentry 'Auto Install CentOS 7' --class fedora --class gnu-linux --class gnu --class os {
	linuxefi vmlinuz inst.repo=http://10.1.0.2/centos ks=http://10.1.0.2:8080/api/v1/ks inst.ks.sendmac  inst.ks.sendsn  quiet
	initrdefi initrd.img
}
```

ks=http://10.1.0.2:8080/api/v1/ks是指定从http获取ks文件，这个后面我再讲

### ks的`%pre`和`%post`阶段
ks的内容里`%pre`到`%end`里的命令都是在安装输入界面之前执行，而`%post`则是安装完成后执行的

#### 利用%pre设置带外网络信息
基本上不买厂商的服务费服务器出厂的带外全部是同一个ip，这里我们用ks的%pre阶段里用ipitool去设置下所有服务器的带外网络，可以提前找台物理服务器测试看看生效否
```
%pre
#!/bin/bash

# 获取eth1的信息,该字段为8
ipmitool raw 0x0c 0x02   0x08   0x04 0x00 0x00 | grep -qw '02'   # dhcp
if [ "$?" -eq 0 ];then
    ipmitool lan set 8 ipsrc static #如果带外是dhcp就设置成静态ip
fi

ipmitool lan set 8 ipaddr {{.IPMIIP}}
ipmitool lan set 8 netmask {{.IPMIMask}}
ipmitool lan set 8 defgw ipaddr {{.IPMIGw}}
ipmitool chassis bootdev disk options=efiboot,persistent # force boot from desik, and set the mode to uefi
#  ipmitool raw 0x00 0x08 0x05  0xE0 0x08  0x00 0x00 0x00 #uefi启动,强制boot from disk，测试这步没必要

%end
```

#### %pre阶段做raid
修改`squashfs.img`就是配合%pre阶段做raid，下面逻辑是我根据cli的输出和查看命令帮助写的，系统盘都是最后俩槽位，所以我是这样写的
```
arcconf getconfig 1 ld | grep -qw 'No logical'
if [ "$?" -eq 0 ];then # 不为0表示没有创建阵列

    #获取实际硬盘的Channel id存成数组:0,1  0,0  取最后两个槽位做系统盘的raid1
    hardDiskChannelIDArray=( $(arcconf getconfig 1 pd | grep Channel | tac | awk -F'[: ()]+' 'NR!=1&&NR<4{print $6}') )
    arcconf task start 1 device all initialize noprompt  #初始化所有硬盘
    #              控制器id               容量  raid级别
    arcconf create    1     logicaldrive  max    1   ${hardDiskChannelIDArray[@]//,/ }  noprompt

fi
```
当然也可以不修改iso，可以在pre阶段下载阵列的cli包安装了再做raid，例如
```
%pre
wget http://10.1.0.2/soft/xxx
chmod u+x xxx
xxx 创建阵列
...
...
%end
```


#### %pre阶段自定义交互
默认的安装进程都是占据了tty1，我们可以利用exec和chvt命令达到交互的方式，例如我们写一个需求，10秒内没输入就设置带外
```
%pre
#!/bin/bash

exec < /dev/tty3 > /dev/tty3 2>&1
chvt 3
read -n1 -t 10 -p 'not need to setup the ipmi network?' select
if [ "$select" != 'n' ];then
  # 此处设置ipmi的网络
fi
chvt 1
```
ks下载是在/tmp里一个随机名字的文件，ks的运行日志是名字.log，如果有语法错误，例如if和fi没闭合会卡住，所以建议ks的%pre和%post阶段每一步操作都要echo信息当作打印日志

#### %post阶段
%post阶段根据自己的场景去写，建议是写一个curl去带上序列号访问后端，这样能表明这台机器安装完成了，单单的请求了ks并不保证了安装完成
```
%post
if [ -r /sys/class/dmi/id/product_serial ];then
 curl -X POST -H "X-System-Serial-Number:$(cat /sys/class/dmi/id/product_serial)" http://10.1.0.2:8080/api/v1/ks
fi
%end
```

#### 一些建议
如果机器的网卡都是一样的可以在boot option那加选项`ksdevice=exxx`指定获取ks的网卡，不然每个一张张网卡去retry比较耗时间
因为是pxe，iso解压放在web上了，所以ks里得指定下从哪里获取包`url --url="http://10.1.0.2/centos"`

最后是ks详细的中文文档链接 https://fedoraproject.org/wiki/Anaconda/Kickstart/zh-cn