# CentOS 7 下 MySQL 5.7 的安装

## 配置 yum 源

在 [https://dev.mysql.com/downloads/repo/yum/](https://links.jianshu.com/go?to=https%3A%2F%2Fdev.mysql.com%2Fdownloads%2Frepo%2Fyum%2F) 找到 yum 源 rpm 安装包

![img](https://upload-images.jianshu.io/upload_images/1458376-6c3dece1d8bd0650.png?imageMogr2/auto-orient/strip|imageView2/2/w/1193/format/webp)

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

## 安装 MySQL

使用 yum install 命令安装

```undefined
shell> yum install -y mysql-community-server
```

问题：

Failing package is: mysql-community-libs-compat-5.7.37-1.el7.x86_64

原因是Mysql的GPG升级了，需要重新获取
使用以下命令即可

```
rpm --import https://repo.mysql.com/RPM-GPG-KEY-mysql-2022
```



## 启动 MySQL 服务

在 CentOS 7 下，新的启动/关闭服务的命令是 `systemctl start|stop`

```undefined
shell> systemctl start mysqld
```

用 `systemctl status` 查看 MySQL 状态

```undefined
shell> systemctl status mysqld
```



## 设置开机启动

```bash
shell> systemctl enable mysqld
# 重载所有修改过的配置文件
shell> systemctl daemon-reload
```



## 修改 root 本地账户密码

mysql 安装完成之后，生成的默认密码在 `/var/log/mysqld.log` 文件中。使用 grep 命令找到日志中的密码。

```bash
shell> grep 'temporary password' /var/log/mysqld.log
```

![img](https:////upload-images.jianshu.io/upload_images/1458376-6694dca4f9eb39a3.png?imageMogr2/auto-orient/strip|imageView2/2/w/1137/format/webp)

查看临时密码

首次通过初始密码登录后，使用以下命令修改密码

```bash
shell> mysql -uroot -p
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'MyNewPass4!'; 
```

或者

```bash
mysql> set password for 'root'@'localhost'=password('MyNewPass4!'); 
```

以后通过 update set 语句修改密码

```bash
mysql> use mysql;
mysql> update user set password=PASSWORD('MyNewPass5!') where user='root';
mysql> flush privileges;
```

> 注意：mysql 5.7 默认安装了密码安全检查插件（validate_password），默认密码检查策略要求密码必须包含：大小写字母、数字和特殊符号，并且长度不能少于8位。否则会提示 ERROR 1819 (HY000): Your password does not satisfy the current policy requirements 错误。查看 [MySQL官网密码详细策略](https://links.jianshu.com/go?to=https%3A%2F%2Fdev.mysql.com%2Fdoc%2Frefman%2F5.7%2Fen%2Fvalidate-password-options-variables.html%23sysvar_validate_password_policy)