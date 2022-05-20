# Tars Docker 部署

## 1 介绍

本节主要介绍采用 docker 来完成框架的部署:

- framework: Tars 框架 Docker 制作脚本, 制作的 docker 包含了框架核心服务和 web 管理平台
- tars-node: tarsnode 节点的镜像，包含各语言的运行时环境，可以将服务发布到 tars-node 容器中, tarsnode 每台机器都存在, 它连接到 framework
- 部署完成后, 打开 framework 主节点上 web 管理平台, 你可以通过 web 管理平台部署和发布服务, 将这些服务发布到 tarsnode 所在的机器上

Docker 开发环境部署可以很方便的在本地拉起服务开始服务的部署、开发和测试。开发环境部署采用单机多容器的部署方式模拟生产环境的服务部署结构。Docker 生产环境部署为生产主机部署 Tars 服务提供参考，相关参数需要根据具体环境变更调整。

开始操作之前，请确保你的服务上已经安装了 docker 环境, 如果没有, 可以参考[docker install]()



## 2 Docker 部署服务开发环境

**如果你想源码自己编译 docker, 请参见** [**Install**]()

当然你可以在多台机器上搭建一套环境, 将你的业务服务发布到这套环境上测试.

如果你没有多台机器, 你又需要完整的搭建环境, 你可以通过以下方式来完成:

- 使用 docker 将 framework/tarsnode 部署在同一台机器上;
- tarsnode 你可以启动多个 docker, 每个 ip 都不同, 表示多台节点机器
- 由于使用 docker, 为了保证 framework/tarsnode 这些 docker 的网络互通, 你需要创建虚拟网络, 以连接这些 docker(这里是 docker 的知识, 有需要自己百度)
- 如果使用--net=host 方式, 你无法在同一台机器上部署 framework/tarsnode, 因为他们使用了相同的端口, 会带来端口冲突

### 2.1 创建 docker 虚拟网络

为了方便虚拟机、Mac、Linux 主机等各种环境下的 docker 部署，在本示例中先创建虚拟网络，模拟现实中的局域网内网环境(注意 docker 都还是在同一台机器, 只是 docker 的虚拟 ip 不同, 模拟多机)

```sh
# 创建一个名为tars的桥接(bridge)虚拟网络，网关172.25.0.1，网段为172.25.0.0
docker network create -d bridge --subnet=172.25.0.0/16 --gateway=172.25.0.1 tars
```

### 2.2 启动 MySQL

- 为框架运行提供 MySQL 服务，若使用宿主机或现有的 MySQL 可以跳过此步骤，建议框架和应用使用不同的 MySQL 服务。
- 注意 MySQL 的 IP 和 root 密码，后续构建中需要使用

```sh
docker run -d -p 3306:3306 \
    --net=tars \
    -e MYSQL_ROOT_PASSWORD="123456" \
    --ip="172.25.0.2" \
    -v /data/framework-mysql:/var/lib/mysql \
    -v /etc/localtime:/etc/localtime \
    --name=tars-mysql \
    mysql:5.7
```

为了验证 MySQL 是否正常启动且能正常连接，可通过 host 中的 mysql 客户端进行登录验证

```
mysql -h 172.25.0.2 -u root -p
```

也可以使用后面已经下载启动的 tars-framework docker 节点进行验证，可以等下再回来操作；

执行 tars-framework 中的 mysql-tool

```sh
docker exec -it tars-framework /bin/bash

cd /usr/local/tars/cpp/deploy/

./mysql-tool --host=172.25.0.2 --user="root" --pass="123456" --port=3306 --check
```

### 2.3 使用 tarscloud/framework 部署框架

#### 2.3.1 拉取镜像

最新版本

```
docker pull tarscloud/framework:latest
```

指定版本:

```
docker pull tarscloud/framework:v{x.y.z}
```

如：

```
docker pull tarscloud/framework:v3.0.4
```

说明:

- 使用指定版本，如：`v2.4.17`，便于开发和生产环境的部署，后期需要升级时可选择更新的版本 tag，升级之前请先查看 GitHub 的 changelog，避免升级到不兼容的版本造成损失。
- 注意这里的 Tag 不是源码 TarsFramework 的 tag 号, 而是 Tars 这个 GIT 仓库的 tag 号, 因为 tarscloud/framework 还包含了 TarsWeb
- 当部署好以后, TarsWeb 页面的右上角显示了内置的 TarsFramework&TarsWeb 的版本号, 以便你对应源码.

#### 2.3.2 启动镜像(目前只考虑了 linux 上, 时间和本机同步)

```sh
# 挂载的/etc/localtime是用来设置容器时区的，若没有可以去掉
# 3000端口为web程序端口
# 3001端口为web授权相关服务端口(docker>=v2.4.7可以不暴露该端口)
docker run -d \
    --name=tars-framework \
    --net=tars \
    -e MYSQL_HOST="172.25.0.2" \
    -e MYSQL_ROOT_PASSWORD="123456" \
    -e MYSQL_USER=root \
    -e MYSQL_PORT=3306 \
    -e REBUILD=false \
    -e INET=eth0 \
    -e SLAVE=false \
    --ip="172.25.0.3" \
    -v /data/framework:/data/tars \
    -v /etc/localtime:/etc/localtime \
    -p 3000:3000 \
    -p 3001:3001 \
    tarscloud/framework:v3.0.4
```

安装完毕后, 访问 `http://${your_machine_ip}:3000` 打开 web 管理平台

### 2.4 Docker 部署 Tars 应用节点

**Tars 应用节点镜像默认为集合环境(Java+GoLang+NodeJs+PHP)的镜像，如果需要可登陆 Docker Hub 查看各语言相关 tag**

#### 2.4.1 拉取镜像

最新版本:

```sh
docker pull tarscloud/tars-node:latest
```

#### 2.4.2 启动 Node(目前只考虑了 linux 上, 时间和本机同步)

最新版本:

```sh
docker run -d \
    --name=tars-node \
    --net=tars \
    -e INET=eth0 \
    -e WEB_HOST="http://172.25.0.3:3000" \
    --ip="172.25.0.5" \
    -v /data/tars:/data/tars \
    -v /etc/localtime:/etc/localtime \
    -p 6600-6610:6600-6610 \
    tarscloud/tars-node:latest
```

- 初始开放了 6600~6610 端口供应用使用，若不够可自行添加
- Node 启动之后会自动向框架 172.25.0.3 进行注册，部署完成之后在框架的 运维管理-》节点管理 中可以看到 IP 为 `172.25.0.5` 的节点启动

**注意, 如果在同一台机器上采用--net=host, 同时启动 framework 和 tars-node 镜像, 是不行的, 因为 framework 中也包含了一个 tars-node, 会导致端口冲突, 启动不了**











## 参考

[TarsDocker 部署](https://tarscloud.gitbook.io/tarsdocs/kuang-jia-bu-shu/docker)