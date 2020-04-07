  这里以ISO安装centos在UEFI模式下的安装举例，pxe启动安装和iso的步骤差距非常小

## 启动过程

- 1: 上电BIOS POST自检然后按照设备顺序启动，加载ISO里的`EFI/BOOT/BOOTX64.EFI`

- 2: `EFI/BOOT/BOOTX64.EFI`会运行`grubx64.efi`也就是我们的grub，然后会加载grub的配置文件`grub.cfg`，接着会出现我们看到安装菜单

![install-menu](https://raw.githubusercontent.com/zhangguanzhang/Image-Hosting/master/machineInstaller/031.png)


- 3: 我们选择`Install CentOS 7`而`grub.cfg`里对应内容为下，意思是寻找带LABEL:`CentOS X86_64`的介质，也就是iso里加载vmlinuz和inittd.img

```
menuentry 'Install CentOS 7' --class fedora --class gnu-linux --class gnu --class os {
	linuxefi /images/pxeboot/vmlinuz inst.stage2=hd:LABEL=CentOS\x207\x20x86_64 quiet
	initrdefi /images/pxeboot/initrd.img
}
```

- 4: 加载`vmlinuz`到内存里，它是个压缩的内核，然后加载解压`initrd.img`当作rootfs，组成一个小型的Linux
- 5: 执行Linux里的init，也就是讲会运行[dracut](https://dracut.wiki.kernel.org/index.php/Main_Page)，这期间我们会看到机器核心数量的企鹅画面，然后到下面图
  
![install-menu](https://raw.githubusercontent.com/zhangguanzhang/Image-Hosting/master/machineInstaller/007.png)

- 6: 图片我是截图的带有ks启动选项的截图，没有指定ks的话就没有用网卡retry获取ks文件的信息，启动到`dracut initquenue hook`的时候进入`stage2`阶段，查看安装介质的根目录下有没有文件`.treeinfo`，正常是下面的内容，也就是没有这个文件从而加载`/LiveOS/squashfs.img`这个rootfs

```
[general]
name = CentOS-7
family = CentOS
timestamp = 1543162874.22
variant = 
version = 7
packagedir = 
arch = x86_64

[stage2]
mainimage = LiveOS/squashfs.img

[images-x86_64]
kernel = images/pxeboot/vmlinuz
initrd = images/pxeboot/initrd.img
boot.iso = images/boot.iso

[images-xen]
kernel = images/pxeboot/vmlinuz
initrd = images/pxeboot/initrd.img
```
- 7: 然后就是我们看到的什么选择时区和分区root密码的那个图形界面了，ks是在这个阶段执行的，如果需要看log我们可以`ctrl+alt+F2-F7`进tty去看日志
> 我们也可以`ip addr add`给此时的机器配置ip然后配置路由，可以scp其他机器上的文件过来，记住这点，可以用来测试出阵列命令配置raid的shell写法

## 修改iso要注意的坑

### 修改initrd.img(有需要的话)

修改initrd的步骤
参考https://wiki.gentoo.org/wiki/Custom_Initramfs
虚机上挂载了centos的iso
```
$ cd 
$ mkdir /root/test
$ mount -t auto /dev/cdrom /mnt
$ cp -ar /mnt /root/iso   # 原来root下没有iso目录，拷贝过来重命名为iso
$ umount /mnt
$ cp iso/images/pxeboot/initrd.img test/
$ cd  test
$ xz -dc ../initrd.img | cpio -id
271820 blocks
```
修改一些东西后下面是打包
```
$ find . -print0 | cpio --null --create --verbose --format=newc | gzip --best > ../myinitrd
$ yum install mkisofs isomd5sum -y
$ cd iso
  #执行封装镜像的命令：
$ mkisofs -U \
  -A 'CentOS 7 x86_64' \
  -V 'CentOS 7 x86_64' \
  -volset 'CentOS 7 x86_64' \
  -J -joliet-long -r -v -T \
   -o /root/CentOS-MY.iso \
   -b isolinux/isolinux.bin \
   -c isolinux/boot.cat \
   -no-emul-boot -boot-load-size 4 \
    -boot-info-table  \
   -eltorito-alt-boot \
  -e images/efiboot.img -no-emul-boot .
```
注意`CentOS 7 x86_64`这部分不能乱写，我看网上的基本都是乱写或者没详细说明的，导致了uefi启动压根就没识别到iso，这个字段必须得是grub.cfg里的字段。也就是因为那个`inst.stage2=`和search以及这里制作iso的要一致
```
$ grep search EFI/BOOT/grub.cfg 
search --no-floppy --set=root -l 'CentOS 7 x86_64'
```
这个字段对上了才能识别到iso，参考 https://access.redhat.com/discussions/762253

校验并写入 md5 值(可选)：
```
$ implantisomd5 /root/CentOS-MY.iso
```


### 修改squashfs.img
使用的厂商服务器没有cli工具远程不交互做raid，然后搜到了个pxe自动做raid的文章: https://blog.csdn.net/lizhihua0925/article/details/53198483
安装界面的时候就是一个在跑着的系统，文章里是把阵列命令加到initrd里，可以在kickstart的`%pre`阶段创建出raid。然后我测试的物理机安装了系统后把storcli传进去获取不到，询问了存储同事说可以用arcconf，最后自己摸索出了一些信息。

#### 判断自己的阵列卡应该使用哪个cli
一般服务器都是下面几大阵列卡厂商，都有对应的操作cli tool
- IBM(MegaCli) 已经被storcli整合了，LSI控制器的阵列卡均可以使用storcli，部分dell的R系列也是megaraid啥的，应该也可以用strocli
- HP(hpacucli)
- Adaptec(arcconf)
- 还有其他的cli，可以去这篇文章看看https://blog.51cto.com/1130739/1771506

如果是要下载storcli可以去博通官网https://www.broadcom.com搜索sotrcli下载`Latest MegaRAID Storcli`

其他的厂商就不知道啥命令了，大家自行摸索，我这里查到是`Adaptec`的
```
$ dmesg | grep -i raid  
[    4.058107] Adaptec aacraid driver 1.2.1[50877]-custom
```
arcconf安装在yum源里是没有的，Adaptec现在改名为Microsemi，下载arcconf的话去官网 https://storage.microsemi.com/en-us/downloads/ 进所有产品->`Adaptec SCSI RAID` -> `Adaptec Series 8 SAS/SATA 12 Gb`  RAID `Storage Manager Downloads` -> ` Downloads  Microsemi Adaptec ARCCONF Command Line Utility v3.01.23531`

#### 修改squashfs.img
自己的虚机上安装arcconf的rpm，然后服务器进入安装界面后ctrl+alt+f2进入终端，配置网络信息后把命令scp过来
```
$ ip addr add 10.0.23.41/24 dev enp61s0f2
$ ip route add default via 10.0.23.1 # 或者route add -net 0.0.0.0/0 gw 10.0.23.1 dev enp61s0f2
$ scp 10.0.23.79:/usr/sbin/arcconf /usr/sbin/
$ arcconf getconfig 1 ld
Controllers found: 1
--------------------------------------------------------------
Logical device infomation
--------------------------------------------------------------
   No logical devices configured



Command completed successfully.
```
命令测试可用，我们后面直接修改`squashfs.img`添加/usr/sbin/arcconf后打包即可。

##### cli报错缺少so的话

squashfs.img这个rootfs大概800多M，里面有很多so，依赖是基本都有的，但是假设添加的命令不可用，我们得把命令相关的动态链接库给拷贝过去
在正常系统上通过ldd看so的信息
```
$ cd /lib64
$ ldd /usr/sbin/arcconf 
	linux-vdso.so.1 =>  (0x00007ffc8a548000)
	libdl.so.2 => /lib64/libdl.so.2 (0x00007fb1fdeee000)
	libpthread.so.0 => /lib64/libpthread.so.0 (0x00007fb1fdcd2000)
	libstdc++.so.6 => /lib64/libstdc++.so.6 (0x00007fb1fd9cb000)
	libm.so.6 => /lib64/libm.so.6 (0x00007fb1fd6c9000)
	libgcc_s.so.1 => /lib64/libgcc_s.so.1 (0x00007fb1fd4b3000)
	libc.so.6 => /lib64/libc.so.6 (0x00007fb1fd0e6000)
	/lib64/ld-linux-x86-64.so.2 (0x00007fb1fe0f2000)
$ echo md5sum $(ldd /usr/sbin/arcconf | grep -Po '/lib64/\K\S+')
md5sum libdl.so.2 libpthread.so.0 libstdc++.so.6 libm.so.6 libgcc_s.so.1 libc.so.6 ld-linux-x86-64.so.2
$ md5sum libdl.so.2 libpthread.so.0 libstdc++.so.6 libm.so.6 libgcc_s.so.1 libc.so.6 ld-linux-x86-64.so.2
df69ee999269a70ee78fed4d39b6ab0a  libdl.so.2
390e1ad4fd8b47508a0b80799acf83bf  libpthread.so.0
bde5a21296b3dc19b0c374b324f97d4a  libstdc++.so.6
0dc5febf77645a7d2e0b8aabbb85995d  libm.so.6
ed3dac8e74ed913de13ee3dd7093e83e  libgcc_s.so.1
03ce524a3e41c8f70daff7ca7a82f9ba  libc.so.6
545bc0249fd1bee457dcec4bbda603b7  ld-linux-x86-64.so.2
```
然后在安装界面的tty里去自行看有没有对应的so和md5值一样，其中的`linux-vdso.so.1`不用管 https://blog.csdn.net/wang_xya/article/details/43985241

##### 给squashfs.img添加arcconf
虚机挂载centos的iso，然后把内容拷贝到`/root/iso`目录，安装相关工具解开squashfs.img
```
$ yum install livecd-tools -y
$ cd /root/iso/LiveOS
$ unsquashfs  squashfs.img 
Parallel unsquashfs: Using 2 processors
1 inodes (16384 blocks) to write

[=========================================================================================================================================================|] 16384/16384 100%

created 1 files
created 2 directories
created 0 symlinks
created 0 devices
created 0 fifos
```
添加命令并打包
```
$ mkdir /root/sq
$ mount -t loop,rw squashfs-root/LiveOS/rootfs.img /root/sq
$ cp /usr/sbin/arcconf /root/sq/usr/sbin/
$ umount /root/sq
$ mv /root/iso/LiveOS/squashfs-root /root/sq/
$ cd /root/sq
$ mksquashfs squashfs-root squashfs.img
$ rm -rf squashfs-root
$ mv squashfs.img /root/iso/LiveOS/
$ mv: overwrite ‘/root/iso/LiveOS/squashfs.img’? y
$ cd /root/iso
$ mkisofs -U \
  -A 'CentOS 7 x86_64' \
  -V 'CentOS 7 x86_64' \
  -volset 'CentOS 7 x86_64' \
  -J -joliet-long -r -v -T \
   -o /root/CentOS-MY.iso \
   -b isolinux/isolinux.bin \
   -c isolinux/boot.cat \
   -no-emul-boot -boot-load-size 4 \
    -boot-info-table  \
   -eltorito-alt-boot \
  -e images/efiboot.img -no-emul-boot .
```
squashfs.img里自带了ipmitool，后面ks的部分会用到它
