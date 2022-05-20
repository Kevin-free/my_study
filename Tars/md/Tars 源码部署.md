# Tars 源码部署

## 1. 环境准备

| 软件            | 软件要求                                                     |
| --------------- | ------------------------------------------------------------ |
| linux 内核版本: | 2.6.18 及以上版本（操作系统依赖）                            |
| gcc 版本:       | 4.8.2 及以上版本、glibc-devel（c++语言框架依赖）             |
| bison 工具版本: | 2.5 及以上版本（c++语言框架依赖）                            |
| flex 工具版本:  | 2.5 及以上版本（c++语言框架依赖）                            |
| cmake 版本：    | 3.2 及以上版本（c++语言框架依赖）                            |
| mysql 版本:     | 5.6 及以上版本（框架运行依赖）                               |
| nvm 版本：      | 0.35.1 及以上版本（web 管理系统依赖, 脚本安装过程中自动安装） |
| node 版本：     | 12.13.0 及以上版本（web 管理系统依赖, 脚本安装过程中自动安装） |

运行服务器要求：安装 linux 系统的机器 or mac 机器（笔者使用的是 CentOS 7）

### 1.1. 编译包依赖下载安装介绍

源码编译过程需要安装:gcc, glibc, bison, flex, cmake, ncurses-devel zlib-devel 等

例如，在 Centos7 下，执行：

```sh
yum install glibc-devel gcc gcc-c++ bison flex cmake which psmisc ncurses-devel zlib-devel yum-utils psmisc telnet net-tools wget unzip
```

在 ubuntu 下执行:

```sh
sudo apt-get install build-essential bison flex cmake psmisc libncurses5-dev zlib1g-dev psmisc telnet net-tools wget unzip
```

在 mac 安装, 请先安装 brew(如何在 mac 上安装 brew, 请自行搜索)

```sh
brew install bison flex cmake
```

### 1.2. MySQL 安装

正式部署时, 如果你的 mysql 可以安装在其他机器上.

Tars 框架安装需要在 mysql 中读写数据, 因此需要安装 mysql, 如果你已经存在 mysql, 可以忽略该步骤.

#### 配置 yum 源

在 [https://dev.mysql.com/downloads/repo/yum/](https://links.jianshu.com/go?to=https%3A%2F%2Fdev.mysql.com%2Fdownloads%2Frepo%2Fyum%2F) 找到 yum 源 rpm 安装包

安装 mysql 源

```csharp
# 下载
shell> wget https://dev.mysql.com/get/mysql57-community-release-el7-11.noarch.rpm
# 安装 mysql 源
shell> yum localinstall mysql57-community-release-el7-11.noarch.rpm
```

用下面的命令检查 mysql 源是否安装成功

```bash
shell> yum repolist enabled | grep "mysql.*-community.*"
```

#### 安装 MySQL

使用 yum install 命令安装

```undefined
shell> yum install -y mysql-community-server
```

#### 启动 MySQL 服务

在 CentOS 7 下，新的启动/关闭服务的命令是 `systemctl start|stop`

```undefined
shell> systemctl start mysqld
```

用 `systemctl status` 查看 MySQL 状态

```undefined
shell> systemctl status mysqld
```

#### 设置开机启动

```bash
shell> systemctl enable mysqld
# 重载所有修改过的配置文件
shell> systemctl daemon-reload
```

#### 修改 root 本地账户密码

mysql 安装完成之后，生成的默认密码在 `/var/log/mysqld.log` 文件中。使用 grep 命令找到日志中的密码。

```bash
shell> grep 'temporary password' /var/log/mysqld.log
```

首次通过初始密码登录后，使用以下命令修改密码

```bash
shell> mysql -uroot -p
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'MyNewPassword'; 
```

或者

```bash
mysql> set password for 'root'@'localhost'=password('MyNewPassword'); 
```

以后通过 update set 语句修改密码

```bash
mysql> use mysql;
mysql> update user set password=PASSWORD('MyNewPassword') where user='root';
mysql> flush privileges;
```

> 注意：mysql 5.7 默认安装了密码安全检查插件（validate_password），默认密码检查策略要求密码必须包含：大小写字母、数字和特殊符号，并且长度不能少于8位。否则会提示 ERROR 1819 (HY000): Your password does not satisfy the current policy requirements 错误。查看 [MySQL官网密码详细策略](https://links.jianshu.com/go?to=https%3A%2F%2Fdev.mysql.com%2Fdoc%2Frefman%2F5.7%2Fen%2Fvalidate-password-options-variables.html%23sysvar_validate_password_policy)



#### 问题&解决

##### 问题

启动 TarsWeb 时，出现 mysql_real_connect: Access denied for user 'root'@'kaiwen.fsbm.cc' 

![image-20220408143540579](C:\Users\win10\AppData\Roaming\Typora\typora-user-images\image-20220408143540579.png)

##### 解决

```sh
mysql> GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'mypassword' WITH GRANT OPTION;  

mysql> FLUSH   PRIVILEGES; 
```

### 1.3. 关闭防火墙

CentOS 7.0 默认使用的是 firewall 作为防火墙

查看防火墙状态

```
firewall-cmd --state
```

停止firewall

```
systemctl stop firewalld.service
```

禁止firewall开机启动

```
systemctl disable firewalld.service 
```



## 2. Tars C++开发环境(源码安装框架必备)

**源码安装框架才需要做这一步, 如果只是用 c++写服务, 只需要下载 tarscpp 代码即可**

下载 TarsFramework 源码

```sh
cd ${source_folder}
git clone https://github.com/TarsCloud/TarsFramework.git --recursive
```

然后进入 build 源码目录

```sh
cd TarsFramework
git submodule update --remote --recursive
cd build
cmake ..
make -j4
```

切换至 root 用户，创建安装目录

```sh
cd /usr/local
mkdir tars
mkdir app
```

如果在非 root 下执行, 则需要先创建目录, 并赋予权限, 比如以: ubuntu(当前用户为 ubuntu)为例

```
sudo mkdir -p /usr/local/tars
sudo mkdir -p /usr/local/app

sudo chown -R ubuntu:ubuntu /usr/local/tars
sudo chown -R ubuntu:ubuntu /usr/local/web
```

安装

```sh
cd build
make install
```



## 3. Tars 框架安装

### 3.1. 框架安装模式

**框架有两种模式:**

- centos/ubuntu/mac 一键部署, 安装过程中需要网络从外部下载资源
- tars-framework>=2.1.0 支持 mac 部署
- 制作成 docker 镜像来完成安装, 制作 docker 过程需要网络下载资源, 但是启动 docker 镜像不需要外网

**框架安装注意事项:**

- 安装过程中, 由于 tars-web 依赖 nodejs, 所以会自动下载 nodejs, npm, pm2 以及相关的依赖, 并设置好环境变量, 保证 nodejs 生效.
- nodejs 的版本目前默认下载的 v12.13.0
- 如果你本机装了低版本 nodejs, 最好提前卸载掉

**注意:需要完成 TarsFramework 的编译和安装** **注意:框架依赖 mysql, 如果你使用 mysql8, 注意需要关闭 ssl 以及启用 mysql_native_password**

下载 tarsweb 并 copy 到/usr/local/tars/cpp/deploy 目录下(注意目录名是 web, 不要搞错!):

```sh
git clone https://github.com/TarsCloud/TarsWeb.git
mv TarsWeb web
cp -rf web /usr/local/tars/cpp/deploy/
```

例如, 这是/usr/local/tars/cpp/deploy 下的文件:

```sh
ubuntu@VM-0-14-ubuntu:/usr/local/tars/cpp/deploy$ ls -l
total 1030h
-rw-r--r--  1 root root  443392 Apr  3 17:22 busybox.exe
-rw-r--r--  1 root root    1922 Apr  3 17:22 centos7_base.repo
-rw-r--r--  1 root root    1395 Apr  3 17:22 Dockerfile
-rwxr-xr-x  1 root root    3260 Apr  4 11:31 docker-init.sh
-rwxr-xr-x  1 root root     319 Apr  3 22:13 docker.sh
drwxr-xr-x  7 root root    4096 Apr  3 17:57 framework
-rwxr-xr-x  1 root root    4537 Apr  4 11:31 linux-install.sh
-rwxr-xr-x  1 root root 9820288 Apr  3 22:16 mysql-tool
-rwxr-xr-x  1 root root     811 Apr  4 11:31 tar-server.sh
-rwxr-xr-x  1 root root   16449 Apr  3 17:22 tars-install.sh
-rwxr-xr-x  1 root root     320 Apr  4 11:31 tars-stop.sh
drwxr-xr-x  2 root root    4096 Apr  3 17:57 tools
drwxr-xr-x 12 root root    4096 Apr  3 21:07 web
-rwxr-xr-x  1 root root    3590 Apr  3 17:22 web-install.sh
-rwxr-xr-x  1 root root    1476 Apr  3 17:22 windows-install.sh
```

### 3.2. 框架部署说明

框架可以部署在单机或者多机上, 多机是一主多从模式, 通常一主一从足够了:

- 主节点只能有一台, 从节点可以多台
- 主节点默认会安装:tarsAdminRegistry, tarspatch, tarsweb, tarslog, tarsstat, tarsproperty, 这几个服务在从节点上不会安装
- tarslog 用于收集所有服务的远程日志, 建议单节点, 否则日志会分散在多机上
- 原则上 tarspatch, tarsweb 可以是多点, 如果部署成多点, 需要把/usr/local/app/patchs 目录做成多机间共享(可以通过 NFS), 否则无法正常发布服务
- 虽然 tarsAdminRegistry 上记录了正在发布服务的状态, 但是原则上也可以可以多节点, tarsweb 调用 tarsAdminRegistry 是 hash 调用
- 后续强烈建议把 tarslog 部署到大硬盘服务器上
- 实际使用中, 即使主从节点都挂了, 也不会影响框架上服务的正常运行, 只会影响发布
- 一键部署会自动安装好 web(自动下载 nodejs, npm, pm2 等相关依赖), 同时开启 web 权限

部署完成后会创建 5 个数据库，分别是 db_tars、db_tars_web、db_user_system、 tars_stat、tars_property。

其中 db_tars 是框架运行依赖的核心数据库，里面包括了服务部署信息、服务模版信息、服务配置信息等等；

db_tars_web 是 web 管理平台用到数据库

db_user_system 是 web 管理平台用到的权限管理数据库

tars_stat 是服务监控数据存储的数据库；

tars_property 是服务属性监控数据存储的数据库；

无论哪种安装方式, 如果成功安装, 都会看到类似如下输出:

```sh
 2019-10-31 11:06:13 INSTALL TARS SUCC: http://xxx.xxx.xxx.xxx:3000/ to open the tars web.
 2019-10-31 11:06:13 If in Docker, please check you host ip and port.
 2019-10-31 11:06:13 You can start tars web manual: cd /usr/local/app/web; npm run prd
```

打开你的浏览器输入: http://xxx.xxx.xxx.xxx:3000/ 如果顺利, 可以看到 web 管理平台

请参考[检查 web 的问题]()中的检查 web 问题章节, 如果没有问题, 请检查机器防火墙

### 3.3. (centos/ubuntu/mac)一键部署

进入/usr/local/tars/cpp/deploy, 执行:

```sh
chmod a+x linux-install.sh
./linux-install.sh MYSQL_HOST MYSQL_PASSWORD INET REBUILD(false[default]/true) SLAVE(false[default]/true) MYSQL_USER MYSQL_PORT
```

说明：

MYSQL_HOST: mysql 数据库的 ip 地址

MYSQL_PASSWORD: mysql 数据库的 MYSQL_USER 的密码(注意密码不要有太特殊的字符, 例如!, 否则 shell 脚本识别有问题, 因为是特殊字符)

INET: 网卡的名称(ifconfig 可以看到, 比如 eth0), 表示框架绑定的本机 IP, 注意不能是 127.0.0.1

REBUILD: 是否重建数据库,通常为 false, 如果中间装出错, 希望重置数据库, 可以设置为 true

SLAVE: 是否是从节点

MYSQL_USER: mysql 用户, 默认是 root

MYSQL_PORT: mysql 端口

举例：

安装两台节点, 一台数据库(假设: 主[192.168.7.151], 从[192.168.7.152], mysql:[192.168.7.153])

主节点上执行(192.168.7.151)

```sh
chmod a+x linux-install.sh
./linux-install.sh 192.168.7.153 tars2015 eth0 false false root 3306
```

主节点执行完毕后, 从节点执行:

```sh
chmod a+x linux-install.sh
./linux-install.sh 192.168.7.153 tars2015 eth0 false true root 3306
```

执行过程中的错误参见屏幕输出, 如果出错可以重复执行(一般是下载资源出错)

注意:

- 脚本会自动根据传入的 MYSQL_USER 和 MYSQL_PASSWORD 来登录数据库，创建 TarsAdmin 账号和授权 Tars 相关数据库供框架使用
- 如果是 ubuntu, 需要 sudo linux-install.sh ...来执行
- 注意: 执行完毕以后, 可以检查 nodejs 环境变量是否生效: node --version
- 安装完成以后, 会在/etc/profile 下写入 nodejs 相关的环境变量
- 如果没生效, 手动执行: source /etc/profile, 如果是 ubuntu 请注意权限的问题

### 3.4. 制作成 docker

目标: 将框架制作成一个 docker, 部署时启动 docker 即可.

首先将 TarsWeb clone 到 TarsFramework 源码根目录, 然后直接在 TarsFramework 源码目录下, 执行:

```sh
git clone https://github.com/TarsCloud/TarsWeb.git web
#x64
sudo ./deploy/docker.sh v1 amd64
#arm64
sudo ./deploy/docker.sh v1 arm64
```

查看系统的架构：arch

```sh
[root@kaiwen TarsFramework]# arch
x86_64
```

问题：

> docker pull i/o timeout

解决：

由于daemon.json没有配置造成的，需修改`daemon.json`

```sh
vim /etc/docker/daemon.json
```

增加拉镜像的地址

```json
{
"registry-mirrors":["https://hub-mirror.c.163.com","https://registry.aliyuncs.com","https://registry.docker-cn.com","https://docker.mirrors.ustc.edu.cn"]
}
```

修改后重启[docker](https://so.csdn.net/so/search?q=docker&spm=1001.2101.3001.7020)服务

```sh
service docker restart
```

可以将 docker 发布到你的机器, 然后执行

```sh
docker run -d --net=host -e MYSQL_HOST=xxxxx -e MYSQL_ROOT_PASSWORD=xxxxx \
        -e MYSQL_USER=root -e MYSQL_PORT=3306 \
        -eREBUILD=false -eINET=enp3s0 -eSLAVE=false \
        -v/data/tars:/data/tars \
        -v/etc/localtime:/etc/localtime \
        tarscloud/framework:v1
```

例如：

```sh
docker run -d --net=host -e MYSQL_HOST=10.4.87.87 -e MYSQL_ROOT_PASSWORD=tars2015 \
        -e MYSQL_USER=root -e MYSQL_PORT=3306 \
        -eREBUILD=false -eINET=eth0 -eSLAVE=false \
        -v/data/tars:/data/tars \
        -v/etc/localtime:/etc/localtime \
        tarscloud/framework:v1
```

说明：

MYSQL_IP: mysql 数据库的 ip 地址

MYSQL_ROOT_PASSWORD: mysql 数据库的 root 密码

INET: 网卡的名称(ifconfig 可以看到, 比如 eth0), 表示框架绑定本机 IP, 注意不能是 127.0.0.1

REBUILD: 是否重建数据库,通常为 false, 如果中间装出错, 希望重置数据库, 可以设置为 true

SLAVE: 是否是从节点

MYSQL_USER: mysql 用户, 默认是 root

MYSQL_PORT: mysql 端口

映射三个目录到宿主机

- 

  -v/data/tars:/data/tars

  > - 
  >
  >   包含了 tars 应用日志, tarsnode/data 目录(业务服务的运行包, 保证 docker 重启, 发布到 docker 内部的服务不会丢失)
  >
  > - 
  >
  >   如果是主机则还包含: web 日志, 发布包目录

**如果希望多节点部署, 则在不同机器上执行 docker run ...即可, 注意参数设置!**

**这里必须使用 --net=host, 表示 docker 和宿主机在相同网络**

docker 制作完毕: tarscloud/framework:v1

```sh
docker ps
```

```sh
[root@kaiwen TarsFramework]# docker ps
CONTAINER ID   IMAGE                           COMMAND                  CREATED          STATUS          PORTS     NAMES
ed9843ea6e7e   tarscloud/framework:v1          "/usr/local/tars/cpp…"   2 seconds ago    Up 1 second               silly_fermi
4d6b2e791f46   moby/buildkit:buildx-stable-1   "buildkitd"              26 minutes ago   Up 17 minutes             buildx_buildkit_tars-builder0
```



## 6. 建议部署方案

虽然本章节是介绍的源码部署 tars 框架的方案, 但是实际使用上, 不建议源码部署, 会导致升级维护比较麻烦, 下面介绍实际使用过程重点部署方案:

- 采用 docker 化部署框架(参考相关文档), 部署两台机器, 主从模式, 并启用--net=host 模式
- 节点机采用物理化部署方案, tarsnode 连接到框架即可
- 升级框架的情况下, 可以升级 docker 即可, 停止老的 docker, 启动新 docker, 注意参与不要选择 rebuild db!!!
- 在 web 上可以远程升级 tarsnode, 正常情况下, tarsnode 一般可以不用升级, 除非有重大功能优化(这里未来不排除会做成自动升级!)
- tarslog 需要大硬盘机器, 部署完框架以后, 建议将框架的 tarslog 服务扩容到其他大硬盘的节点机上
- 如果你使用 windows, 你可以考虑框架部署使用 docker, 节点机使用 windows 部署 tarsnode

如果你对 k8s 比较熟悉, 也可以将 Tars 部署在 k8s 上, [请参见 k8s 部署文档]()



## 参考

[Linux/Mac 源码部署](https://tarscloud.gitbook.io/tarsdocs/kuang-jia-bu-shu/source)

[CentOS 7 下 MySQL 5.7 的安装与配置](https://www.jianshu.com/p/1dab9a4d0d5f)





