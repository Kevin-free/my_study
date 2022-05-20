# TarsGo 实践：RPC 通信



## 准备工作

| 当前所用环境  | 说明                                    |
| ------------- | --------------------------------------- |
| CentOS 7.9    | 操作系统依赖 linux 2.6.18 及以上版本    |
| go 1.18       | tarsgo 要求 golang 版本在1.14.x及以上。 |
| MySQL 5.7     | 框架运行依赖                            |
| tarsgo v1.2.0 | 生成代码脚手架                          |
| tars2go       | tars协议转Golang工具                    |

有关环境准备，请参看前文。





## 设计目标

这次我来模拟服务间的 RPC 调用，设计一个服务 A 提供 HTTP 服务给客户端调用，设计一个服务 B 提供 RPC 服务给服务 A 调用。

这里通过一个简单的需求来实现：HTTP 服务根据传入的时间格式，调用 RPC 服务返回相应的服务器时间。

![](http://kevinpub.ifree258.top/mystudy/Tars/001.png)



## 代码设计

### 服务命名

这里我需要创建两个服务：一个服务 A（GoWebServer）提供 HTTP 服务给客户端调用，一个服务 B（GoRpcServer）提供 RPC 服务给服务 A 调用。

Tars 实例的名称其中一个非常重要的作用就是用于服务间名字服务寻址。对于 HTTP 这样的直接对外提供服务的实例而言，其实这块相对不是很重要。但是供内部 RPC 调用的服务，其名称就很重要了，它是其他服务进行寻址的重要依据。

这里我把两个服务分别命名为 `kevin.GoWebServer.GoWebObj` 和 `kevin.GoRpcServer.GoRpcObj` 。



### 生成代码

运行tarsgo脚手架，自动创建服务必须的文件。

```sh
tarsgo make [App] [Server] [Servant] [GoModuleName]
例如： 
tarsgo make TestApp HelloGo SayHello github.com/Tars/test
```

在 `kevin/tars-go-demo/` 目录下，分别生成两个服务的基础框架代码。

```sh
[root@kaiwen tars-go-demo]# tarsgo make kevin GoWebServer GoWeb kevin/tars-go-demo/GoWebServer
go install github.com/TarsCloud/TarsGo/tars/tools/tars2go@latest
🚀 Creating server kevin.GoWebServer, please wait a moment.

go: creating new go.mod: module kevin/tars-go-demo/GoWebServer
go: to add module requirements and sums:
        go mod tidy

CREATED GoWebServer/GoWeb.tars (166 bytes)
CREATED GoWebServer/GoWeb_imp.go (602 bytes)
CREATED GoWebServer/client/client.go (446 bytes)
CREATED GoWebServer/config/config.conf (714 bytes)
CREATED GoWebServer/debugtool/dumpstack.go (411 bytes)
CREATED GoWebServer/go.mod (47 bytes)
CREATED GoWebServer/main.go (511 bytes)
CREATED GoWebServer/makefile (156 bytes)
CREATED GoWebServer/scripts/makefile.tars.gomod (4181 bytes)
CREATED GoWebServer/start.sh (67 bytes)

>>> Great！Done! You can jump in GoWebServer
>>> Tips: After editing the Tars file, execute the following cmd to automatically generate golang files.
>>>       /root/gocode/bin/tars2go *.tars
$ cd GoWebServer
$ ./start.sh
🤝 Thanks for using TarsGo
📚 Tutorial: https://doc.tarsyun.com/


[root@kaiwen tars-go-demo]# tarsgo make kevin GoRpcServer GoRpc kevin/tars-go-demo/GoRpcServer
go install github.com/TarsCloud/TarsGo/tars/tools/tars2go@latest
🚀 Creating server kevin.GoRpcServer, please wait a moment.

go: creating new go.mod: module kevin/tars-go-demo/GoRpcServer
go: to add module requirements and sums:
        go mod tidy

CREATED GoRpcServer/GoRpc.tars (166 bytes)
CREATED GoRpcServer/GoRpc_imp.go (602 bytes)
CREATED GoRpcServer/client/client.go (446 bytes)
CREATED GoRpcServer/config/config.conf (714 bytes)
CREATED GoRpcServer/debugtool/dumpstack.go (411 bytes)
CREATED GoRpcServer/go.mod (47 bytes)
CREATED GoRpcServer/main.go (511 bytes)
CREATED GoRpcServer/makefile (156 bytes)
CREATED GoRpcServer/scripts/makefile.tars.gomod (4181 bytes)
CREATED GoRpcServer/start.sh (67 bytes)

>>> Great！Done! You can jump in GoRpcServer
>>> Tips: After editing the Tars file, execute the following cmd to automatically generate golang files.
>>>       /root/gocode/bin/tars2go *.tars
$ cd GoRpcServer
$ ./start.sh
🤝 Thanks for using TarsGo
📚 Tutorial: https://doc.tarsyun.com/
```



### 修改代码

#### 服务提供方：GoRpcServer

我先修改 GoRpcServer 服务的代码，主要有如下几步：

- 设计协议
- 实现协议
- 

##### 设计协议：GoRpc.tars

修改 `GoRpc.tars` ：

```sh
// kevin/tars-go-demo/GoRpcServer/GoRpc.tars

module kevin
{
	struct GetTimeReq
    {
        0 optional  string  timeFmt;
    };

    struct GetTimeRsp
    {
        0 require   long    utcTimestamp;   // UTC UNIX timestamp
        1 require   long    localTimestamp;
        2 require   string  localTimeStr;
    };

    interface DateTime
    {
        int GetTime(GetTimeReq req, out GetTimeRsp rsp);
    };
};
```

接着，使用 TarsGo 的工具，将协议文件转换为源文件：

```sh
[root@kaiwen GoRpcServer]# tars2go GoRpc.tars 
GoRpc.tars [GoRpc.tars]
```

执行后，`tars2go` 会在当前目录下，根据 `.tars` 文件中指定的 `module` 字段，生成一个新的目录。比如上面的协议文件，module 是 “`kevin`”，那么 tars2go 就生成 `kevin` 目录。读者可以自行查看目录下的文件，如果 `.tats` 文件更新的话，需要再次执行 `tats2go` 命令刷新相应的文件——当然，我觉得完全可以调整 makefile 的逻辑来自动实现这一点。

##### 实现协议：GoRpc_imp.go

在 `GoRpc_imp.go` 文件中实现协议：

```go
// kevin/tars-go-demo/GoRpcServer/GoRpc_imp.go

package main

import (
	"fmt"
	"kevin/tars-go-demo/GoRpcServer/kevin" // Note 1
	"strings"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
)

// GoRpcImp servant implementation
type GoRpcImp struct{}            // Note 2
var log = tars.GetLogger("logic") // Note 3

func (imp *GoRpcImp) GetTime(req *kevin.GetTimeReq, rsp *kevin.GetTimeRsp) (int32, error) { // Note 4
	log.Debug("Enter GetTime ")

	// get timestamp
	utc_time := time.Now()
	local_time := utc_time.Local()

	// convert time string
	var time_str string
	if "" == (*req).TimeFmt {
		log.Debug("Use default time format")
		time_str = local_time.Format("01/02 15:04:05 2006")
	} else {
		/**
		 * reference:
		 * - [go 时间格式风格详解](https://my.oschina.net/achun/blog/142315)
		 * - [Go 时间格式化和解析](https://www.kancloud.cn/itfanr/go-by-example/81698)
		 */
		log.Info(fmt.Sprintf("Got format string: %s", (*req).TimeFmt))
		time_str = (*req).TimeFmt
		time_str = strings.Replace(time_str, "YYYY", "2006", -1)
		time_str = strings.Replace(time_str, "yyyy", "2006", -1)
		time_str = strings.Replace(time_str, "YY", "06", -1)
		time_str = strings.Replace(time_str, "yy", "06", -1)
		time_str = strings.Replace(time_str, "MM", "01", -1)
		time_str = strings.Replace(time_str, "dd", "02", -1)
		time_str = strings.Replace(time_str, "DD", "02", -1)
		time_str = strings.Replace(time_str, "hh", "15", -1)
		time_str = strings.Replace(time_str, "mm", "04", -1)
		time_str = strings.Replace(time_str, "ss", "05", -1)
		log.Info("Convert as golang format: ", time_str)
		time_str = local_time.Format(time_str)
	}

	// construct response
	(*rsp).UtcTimestamp = utc_time.Unix()
	(*rsp).LocalTimestamp = local_time.Unix()
	(*rsp).LocalTimeStr = time_str
	return 0, nil
}
```

针对代码里的几个 Note 注意说明如下：

1. 这里导入的包，就是前文 `tars2go` 所生成的 `kevin` 目录下的 go 文件。通过导入该包，我们就可以获取到我们在前面的 `.tars` 文件中所定义的结构体和方法。这里其实是写了一个基于 `$GOPATH` 的绝对路径来存取该包。
2. 定义了该 servant 的对象，供 server 调用——这个后文讲到 server 时会再提到。
3. 使用 tars 自带的服务器本地日志模块。该模块需要传入一个文件名参数，模块会根据该文件名，在 `/usr/local/app/tars/app_log/amc/GoTarsServer/` 目录下生成日志文件。比如我用的 log 文件名就是：`kevin.GoTarsServer_logic.log`。
4. 这是 `.tars` 文件中 `GetTime` 的实现，它作为 `GoTarsImp` 对象的一个方法来实现。从返回值的角度，TarsGo rpc 方法的返回值除了协议中定义的（本例中是 `int`，对应于 Go 的 `int32`）之外，还有一个 `error`，如果需要的话，读者可以利用。

**Tips**

> 细心的读者可能会发现，在上面的实现中，数据变量名和协议中定义的并不相同。是的，这就是刚转 Go 的开发者很容易遇到的坑之一：Go 语言是使用变量 / 方法 / 常量的命名方式来决定其可见性的，只有在首字母为大写的时候，该元素才能供外部访问。
>
> 笔者特意在 `.tars` 文件中，变量名采用了首字母小写的驼峰式命名法。读者可以看到，`tars2go`会自动将变量名和方法名的首字母改为大写，以保证其可见性。请开发者注意，否则会在编译时遇到未定义错误。

##### 服务入口：main.go

```go
// kevin/tars-go-demo/GoRpcServer/main.go

package main

import (
	"kevin/tars-go-demo/GoRpcServer/kevin"

	"github.com/TarsCloud/TarsGo/tars"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()
	// New servant imp
	imp := new(GoRpcImp)
	// New servant
	app := new(kevin.DateTime)
	// Register Servant
	app.AddServant(imp, cfg.App+"."+cfg.Server+".GoTarsObj") //Register Servant
	// Run application
	tars.Run()
}

```



#### 服务调用方：GoWebServer

##### 服务入口：main.go

首先是 `main.go`，这里参照官方 Guide 的写法就好了，TarsGo 的 HTTP 实现用的是 Go 原生的组件。

```go
// kevin/tars-go-demo/GoWebServer/main.go

package main

import (
	"github.com/TarsCloud/TarsGo/tars"
)

func main() {
	mux := &tars.TarsHttpMux{}
	mux.HandleFunc("/", RpcRootHandler)
	cfg := tars.GetServerConfig()
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".GoWebObj") //Register http server
	tars.Run()
}
```

##### 调用 RPC：GoWeb_imp.go

```go
// kevin/tars-go-demo/GoWebServer/GoWeb_imp.go

package main

import (
	"fmt"
	"kevin/tars-go-demo/GoRpcServer/kevin"
	"net/http"

	"github.com/TarsCloud/TarsGo/tars"
)

var log = tars.GetLogger("logic")

// RPC
func RpcRootHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug("Enter RPCRootHandler")

	comm := tars.NewCommunicator()
	app := new(kevin.DateTime)
	// obj := "kevin.GoRpcServer.GoRpcObj@tcp -h 172.25.0.5 -p 9001 -t 60000" // Note 1-1 
	// Fixed: tarsregistry is Inactive!
	obj := "kevin.GoRpcServer.GoRpcObj"
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 172.25.0.3 -p 17890 -t 60000") // Note 1-2

	log.Debug("RpcRootHandler obj:", obj)

	req := kevin.GetTimeReq{} // Node 2
	rsp := kevin.GetTimeRsp{} // Node 2
	req.TimeFmt = "YYYY-MM-DD hh:mm:ss"

	comm.StringToProxy(obj, app)        // Note 3
	ret, err := app.GetTime(&req, &rsp) // Node 3
	if err != nil {
		// ...... 系统错误处理
		log.Error("GetTime error: ", err)
	} else {
		// ...... 从 rsp 中取出
		log.Debug("GetTime ret:", ret)
	}
	
	ret_str := fmt.Sprintf("{\"msg\":\"Hello, TarsGo! RPC!\", \"time\":\"%s\"}\n", rsp.LocalTimeStr)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(ret_str)) // 写入返回数据
}

```

主要逻辑的说明如下：

1. 选择路由：
   - 1-1 方式为直连服务，其中  `10.4.87.87` 和 `9001` 是  服务`GoRpcServer` 的地址
   - 1-2 方式为主控，其中 `172.25.0.3` 和 `17890` 是 Tars 主控 `tarsregistry` 的地址
2. 准备用于承载参数和返回值的结构体
3. 这两行就是实际的 rpc 调用

##### 模块依赖：go.mod

`GoWebServer `包需要依赖 `GoRpcServer` 包中的 `GetTime()` 方法。

但是编译时却报错：

```
package kevin/tars-go-demo/GoRpcServer/kevin is not in GOROOT (/usr/local/go/src/kevin/tars-go-demo/GoRpcServer/kevin) (compile)
```

这是因为这两个包不在同一个项目路径下，你想要导入本地包，并且这些包也没有发布到远程的github或其他代码仓库地址。这个时候我们就需要在`go.mod`文件中使用`replace`指令。

手动添加，并使用 `replace` **相对路径**来寻找`GoRpcServer`这个包。

```go
// kevin/tars-go-demo/GoWebServer/go.mod

module kevin/tars-go-demo/GoWebServer

go 1.18

require (
	kevin/tars-go-demo/GoRpcServer v0.0.0
	... 忽略其他
)

replace kevin/tars-go-demo/GoRpcServer => ../GoRpcServer

```





## 打包部署

如果开发过程中, 每次都需要手工发布到 web 平台调试, 调试效率是非常低, 因此 Tars 平台提供了一个方式, 能够一键发布服务到 Tars 框架上.

使用方式如下:

- 这需要 web >= 2.0.0, tarscpp>=2.1.0 的版本才能支持.
- 完成框架安装后，登录 TarsWeb 后，在【用户中心-Token 管理】中，新增一个 token

### 新增 Token

登录 TarsWeb 后，在【用户中心-Token 管理】中，新增一个 token

![image-20220411152458413](http://kevinpub.ifree258.top/mystudy/Tars/002.png)

### 新增部署

在【运维管理-部署申请】中，填写相关信息，新部署一个服务。

![image-20220414154210984](http://kevinpub.ifree258.top/mystudy/Tars/003.png)

所以我们新增 “`kevin.GoWebServer.GoWebObj`”，就是在各项中如下填写：

- 应用：`kevin`
- 服务名称：`GoWebServer`
- 服务类型：`tars_go`
- 模板：`tars.default`
- OBJ：`GoWebObj`
- 节点：填写你打算部署的 IP 地址`172.25.0.5`
- 端口类型：`TCP`
- 协议：`非TARS`
- 端口：填写你打算部署的端口`6600`，也可以填好信息后点 “获取端口” 来生成。

各项填写完毕后，点 “确定”，然后刷新界面，重新进入 Tars 管理平台主页，可以看到界面左边的列表就多了上面的配置：

![image-20220414154344928](http://kevinpub.ifree258.top/mystudy/Tars/004.png)

`GoRpcServer` 服务，注意这里选择的是【TARS协议】

![image-20220414160009297](http://kevinpub.ifree258.top/mystudy/Tars/005.png)



### 新建 CMakeLists.txt

c++版本的 cmake 已经内嵌了命令行在服务的 CMakeLists.txt 中, 比如

```sh
# kevin/tars-go-demo/GoWebServer/CMakeLists.txt

execute_process(COMMAND go env GOPATH OUTPUT_VARIABLE GOPATH)

string(REGEX REPLACE "\n$" "" GOPATH "${GOPATH}")

include(${GOPATH}/src/github.com/TarsCloud/TarsGo/cmake/tars-tools.cmake)

cmake_minimum_required(VERSION 2.8)

project(GoWebServer Go) # select GO compile

gen_server(kevin GoWebServer)

# go mod init
# mkdir build
# cd build
# cmake .. -DTARS_WEB_HOST=${TARS_WEB_HOST} -DTARS_TOKEN=${TARS_TOKEN}
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

##### 注意：

第一次执行  `make GoWebServer-upload` 会显示 `0no active server, please start server first!`

![image-20220414154808430](http://kevinpub.ifree258.top/mystudy/Tars/006.png)

第一次需要先手动发布！

![image-20220414154924099](http://kevinpub.ifree258.top/mystudy/Tars/007.png)

之后执行 `make GoWebServer-upload` 才会成功。

![image-20220414160930390](http://kevinpub.ifree258.top/mystudy/Tars/008.png)





## 服务测试

因为服务是部署在宿主机（10.4.87.87）上 Docker 虚拟出的 IP 节点（172.25.0.5）中，所以访问两个 IP 节点都可以。

![image-20220414162434915](http://kevinpub.ifree258.top/mystudy/Tars/009.png)





如果是直接部署在宿主机（10.4.87.87）中，则只能通过宿主机节点访问。

![image-20220414162544435](http://kevinpub.ifree258.top/mystudy/Tars/010.png)





## 参考

[腾讯 Tars-Go 服务 Hello World——RPC 通信](https://cloud.tencent.com/developer/inventory/2481/article/1382458)

[使用go module导入本地包](https://zhuanlan.zhihu.com/p/109828249)



