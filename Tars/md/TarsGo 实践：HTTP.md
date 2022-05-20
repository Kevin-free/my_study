# TarsGo 实践：HTTP

## 准备工作

### Go 环境

TarsGo 要求 Golang 版本在1.14.x 及以上。

1. - 最新安装包下载 [goland下载](https://golang.org/dl/)，解压缩到/usr/local/。

     ```sh
     tar -C /usr/local -xzf go.XXX.tar.gz
     ```

   - 配置环境变量

     ```sh
     vim ~/.bashrc
     ```

   - 配置内容并退出

     ```sh
     export GOROOT=/usr/local/go
     export GOPATH=$HOME/gocode
     export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
     ```

   - 导入系统环境

     ```sh
     source ~/.bashrc
     ```

   - 查看版本

     ```sh
     go version
     ```

     如果在国内, 可以设置go代理:  

     ```
     go env -w GOPROXY=https://goproxy.cn   
     ```

     另外请设置go模式为:

     ```
     go env -w GO111MODULE=auto
     ```

2. 安装TarsGo项目创建脚手架

   ```sh
   # < go 1.17
   go get -u github.com/TarsCloud/TarsGo/tars/tools/tarsgo
   # >= go 1.17
   go install github.com/TarsCloud/TarsGo/tars/tools/tarsgo@latest
   ```

3. 安装编译tars协议转Golang工具

   ```sh
   # < go 1.17
   go get -u github.com/TarsCloud/TarsGo/tars/tools/tars2go
   # >= go 1.17
   go install github.com/TarsCloud/TarsGo/tars/tools/tars2go@latest
   ```

检查下GOPATH路径下tars是否安装成功。



## 代码设计

TarsGo 的官方 [Quick Start 文档](https://github.com/TarsCloud/TarsGo/blob/master/docs/tars_go_quickstart.md) 的第一个例子，就是使用 tars 协议进行 server-client 的通信。不过我个人觉得，要说后台服务程序的 hello world 的话，第一个应该是 http 服务嘛，毕竟程序一运行就可以看到效果，这才是 hello world 嘛。

### 给服务命名

Tars 实例的名称，有三个层级，分别是 App（应用）、Server（服务）、Servant（服务者，有时也称 Object）三级。在[前文](https://cloud.tencent.com/developer/article/1372998?from=10910)我们已经初步接触到了：比如 Tars 基础框架中的 `tarsstat`，其服务的完整名称即为：`tars.tarsstat.StatObj`。

Tars 实例的名称其中一个非常重要的作用就是用于服务间名字服务寻址。而对于 HTTP 这样的直接对外提供服务的实例而言，其实这块相对不是很重要，我们更多的是以描述服务功能的角度去命名。这里我把我的 HTTP 服务命名为 `kevin.GoWebServer.GoWebObj`

### 创建基础框架

和 TarsCpp 一样，TarsGo 也提供了一个 `create_tars_server.sh` 脚本用于生成 tars 服务，但却没有提供 `create_http_server.sh` 生成 HTTP 服务。所以这里我们就直接用它就行了：

```js
$ cd $GOPATH/src/github.com/TarsCloud/TarsGo/tars/tools
$ chmod +x create_tars_server.sh
$ ./create_tars_server.sh kevin GoWebServer GoWeb
```

执行后我们可以查看生成的文件，清除不需要的：

```js
$ cd $GOPATH/src/kevin/GoWebServer
$ rm -rf GoWeb.tars client debugtool
$ chmod +x start.sh
$ ls -l
total 44
-rw-rw-r-- 1 centos centos  964 Jan  5 22:09 config.conf
-rw-rw-r-- 1 centos centos  303 Jan  5 22:09 goweb_imp.go
-rw-rw-r-- 1 centos centos  422 Jan  5 22:09 main.go
-rw-rw-r-- 1 centos centos  252 Jan  5 22:09 makefile
-rw-rw-r-- 1 centos centos   59 Jan  5 22:09 start.sh
drwxrwxr-x 2 centos centos 4096 Jan  5 22:09 vendor
```

其实留下的，各文件里的内容，实际上我们都要完全替换掉的……首先是修改 makefile，自动生成的 makefile 内容是这样的：

```js
$ cat makefile 
APP       := kevin
TARGET    := GoWebServer
MFLAGS    :=
DFLAGS    :=
CONFIG    := client
STRIP_FLAG:= N
J2GO_FLAG:= 

libpath=${subst :, ,$(GOPATH)}
$(foreach path,$(libpath),$(eval -include $(path)/src/github.com/TarsCloud/TarsGo/tars/makefile.tars))
```

我们把 “`CONFIG := client`” 行去掉就行了。

### 代码修改

#### main.go

接着是修改代码了。首先是 `main.go`，这里参照官方 Guide 的写法就好了，TarsGo 的 HTTP 实现用的是 Go 原生的组件。我稍微调整了一下，把回调函数放在 `goweb_imp.go` 中，将 `main.go` 简化为：

```js
//kevin/GoWebServer/main.go

package main

import (
	"github.com/TarsCloud/TarsGo/tars"
)

func main() {
	mux := &tars.TarsHttpMux{}
	mux.HandleFunc("/", HttpRootHandler)
	cfg := tars.GetServerConfig()
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".GoWebObj") //Register http server
	tars.Run()
}
```

代码还是比较简单的，无需多言。

#### goweb_imp.go

`main.go` 中的 `HTTPRootHandler` 回调函数定义在业务的主要实现逻辑 `goweb_imp.go` 文件中：

```js
// kevin/GoWebServer/goweb_imp.go

package main

import (
	"fmt"
    "time"
	"net/http"
)

func HttpRootHandler(w http.ResponseWriter, r *http.Request) {
    time_fmt := "2006-01-02 15:04:05"
    local_time := time.Now().Local()
    time_str = local_time.Format(time_fmt)
    ret_str = fmt.Sprintf("{\"msg\":\"Hello, Tars-Go!\", \"time\":\"%s\"}\n", time_str)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(ret_str))
	return
}
```



## 部署发布

如果开发过程中, 每次都需要手工发布到 web 平台调试, 调试效率是非常低, 因此 Tars 平台提供了一个方式, 能够一键发布服务到 Tars 框架上.

使用方式如下:

- 这需要 web >= 2.0.0, tarscpp>=2.1.0 的版本才能支持.
- 完成框架安装后，登录 TarsWeb 后，在【用户中心-Token 管理】中，新增一个 token
- linux 上使用 curl 命令即可完成服务的上传和发布,以 Test/HelloServer 为例, [参考 cmake 管理规范]()

```sh
curl http://${your-web-host}/api/upload_and_publish?ticket=${token} -Fsuse=@HelloServer.tgz -Fapplication=Test -Fmodule_name=HelloServer -Fcomment=dev
```

### 新增 Token

登录 TarsWeb 后，在【用户中心-Token 管理】中，新增一个 token

![image-20220411152458413](C:\Users\win10\AppData\Roaming\Typora\typora-user-images\image-20220411152458413.png)

### 新建 CMakeLists.txt

c++版本的 cmake 已经内嵌了命令行在服务的 CMakeLists.txt 中, 比如

```sh
# kevin/GoWebServer/CMakeLists.txt

execute_process(COMMAND go env GOPATH OUTPUT_VARIABLE GOPATH)

string(REGEX REPLACE "\n$" "" GOPATH "${GOPATH}")

include(${GOPATH}/src/github.com/TarsCloud/TarsGo/cmake/tars-tools.cmake)

cmake_minimum_required(VERSION 2.8)

# 注意匹配：GoWebServer 为服务名（Server）
project(GoWebServer Go) # select GO compile

# 注意匹配：Kevin 为应用名（APP），GoWebServer 为服务名（Server）
gen_server(kevin GoWebServer)

# go mod init
# mkdir build
# cd build
# cmake ..
# make
```

新建 CMakeLists.txt 后，执行如下命令：

```sh
go mod init
mkdir build
cd build
cmake .. -DTARS_WEB_HOST=${TARS_WEB_HOST} -DTARS_TOKEN=${TARS_TOKEN}
make
```

注意:

- 替换 TARS_WEB_HOST 和 TARS_TOKEN
- HelloServer.tgz 是 c++的发布包, java 对应是 war 包, 其他语言类似, 对应你上传到 web 平台的发布包

### 打包和上传

你可以使用以下命令一键打包和上传服务:

```sh
# build 目录下
make GoWebServer-tar
make GoWebServer-upload
```

































