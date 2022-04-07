# Tars 初体验



## Linux 下开发环境搭建

1. 登陆到开发机上，建立工作目录

   ```sh
   mkdir /root/gocode
   mkdir src bin
   cd src
   ```

2. 安装golang语言环境（tarsgo要求golang版本在1.14.x及以上。）

   - 最新安装包下载 [goland下载](https://golang.org/dl/)，解压缩到/usr/local/。

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

3. 安装TarsGo项目创建脚手架

   ```sh
   # < go 1.17
   go get -u github.com/TarsCloud/TarsGo/tars/tools/tarsgo
   # >= go 1.17
   go install github.com/TarsCloud/TarsGo/tars/tools/tarsgo@latest
   ```

4. 安装编译tars协议转Golang工具

   ```sh
   # < go 1.17
   go get -u github.com/TarsCloud/TarsGo/tars/tools/tars2go
   # >= go 1.17
   go install github.com/TarsCloud/TarsGo/tars/tools/tars2go@latest
   ```

检查下GOPATH路径下tars是否安装成功。



## 服务编写

### 创建服务

运行tarsgo脚手架，自动创建服务必须的文件。

```sh
tarsgo make [App] [Server] [Servant] [GoModuleName]
例如： 
tarsgo make TestApp HelloGo SayHello github.com/Tars/test
```

命令执行后将生成代码至GOPATH中，并以`APP/Server`命名目录，生成代码中也有提示具体路径。

```sh
[root@kaiwen kevin]# tarsgo make TestApp HelloGo SayHello github.com/Tars/test
go install github.com/TarsCloud/TarsGo/tars/tools/tars2go@latest
go: downloading github.com/TarsCloud/TarsGo v1.3.2
go: downloading github.com/TarsCloud/TarsGo/tars/tools/tars2go v0.0.0-20220402065337-f1a1f53940cb
🚀 Creating server TestApp.HelloGo, layout repo is https://github.com/TarsCloud/TarsGo.git, please wait a moment.

From https://github.com/TarsCloud/TarsGo
   02257a2..f1a1f53  master                -> origin/master
 * [new branch]      feature/lbbniu/tarsgo -> origin/feature/lbbniu/tarsgo
 * [new branch]      feature/lbbniu/wrr    -> origin/feature/lbbniu/wrr
 * [new tag]         v1.3.2                -> v1.3.2
Updating 02257a2..f1a1f53
Fast-forward
 CONTRIBUTING.md                                    |  19 +-
 README.md                                          | 976 +--------------------
 README.zh.md                                       | 938 +-------------------
 docs/images/fork_button.png                        | Bin 20890 -> 0 bytes
 docs/images/tars_go_quickstart_bushu1.png          | Bin 32454 -> 0 bytes
 docs/images/tars_go_quickstart_bushu1_en.png       | Bin 38394 -> 0 bytes
 docs/images/tars_go_quickstart_release.png         | Bin 21043 -> 0 bytes
 docs/images/tars_go_quickstart_release_en.png      | Bin 28402 -> 0 bytes
 .../images/tars_go_quickstart_service_inactive.png | Bin 30889 -> 0 bytes
 .../tars_go_quickstart_service_inactive_en.png     | Bin 52413 -> 0 bytes
 docs/images/tars_go_quickstart_service_ok.png      | Bin 25715 -> 0 bytes
 docs/images/tars_go_quickstart_service_ok_en.png   | Bin 39913 -> 0 bytes
 docs/images/tars_web_index.png                     | Bin 45109 -> 0 bytes
 docs/images/tars_web_index_en.png                  | Bin 58615 -> 0 bytes
 docs/tars_go_performance.md                        |  20 -
 docs/tars_go_quickstart.md                         | 308 -------
 docs/tars_go_quickstart_en.md                      | 303 -------
 17 files changed, 35 insertions(+), 2529 deletions(-)
 delete mode 100644 docs/images/fork_button.png
 delete mode 100644 docs/images/tars_go_quickstart_bushu1.png
 delete mode 100644 docs/images/tars_go_quickstart_bushu1_en.png
 delete mode 100644 docs/images/tars_go_quickstart_release.png
 delete mode 100644 docs/images/tars_go_quickstart_release_en.png
 delete mode 100644 docs/images/tars_go_quickstart_service_inactive.png
 delete mode 100644 docs/images/tars_go_quickstart_service_inactive_en.png
 delete mode 100644 docs/images/tars_go_quickstart_service_ok.png
 delete mode 100644 docs/images/tars_go_quickstart_service_ok_en.png
 delete mode 100644 docs/images/tars_web_index.png
 delete mode 100644 docs/images/tars_web_index_en.png
 delete mode 100644 docs/tars_go_performance.md
 delete mode 100644 docs/tars_go_quickstart.md
 delete mode 100644 docs/tars_go_quickstart_en.md

go: creating new go.mod: module github.com/Tars/test
go: to add module requirements and sums:
        go mod tidy

CREATED HelloGo/SayHello.tars (171 bytes)
CREATED HelloGo/SayHello_imp.go (620 bytes)
CREATED HelloGo/client/client.go (444 bytes)
CREATED HelloGo/config.conf (716 bytes)
CREATED HelloGo/debugtool/dumpstack.go (412 bytes)
CREATED HelloGo/go.mod (37 bytes)
CREATED HelloGo/main.go (517 bytes)
CREATED HelloGo/makefile (185 bytes)
CREATED HelloGo/scripts/makefile.tars.gomod (4181 bytes)
CREATED HelloGo/start.sh (56 bytes)

>>> Great！Done! You can jump in HelloGo
>>> Tips: After editing the Tars file, execute the following cmd to automatically generate golang files.
>>>       /root/gocode/bin/tars2go *.tars
$ cd HelloGo
$ ./start.sh
🤝 Thanks for using TarsGo
📚 Tutorial: https://tarscloud.github.io/TarsDocs/
```

### 定义接口文件

接口文件定义请求方法以及参数字段类型等，有关接口定义文件说明参考tars_tup.md

为了测试我们定义一个echoHello的接口，客户端请求参数是短字符串如 "tars"，服务响应"hello tars".

```sh
# cat HelloGo/SayHello.tars 
module TestApp{
    interface SayHello{
        int echoHello(string name, out string greeting); 
    };
};
```

**注意**： 参数中**out**修饰关键字标识输出参数。

### 服务端开发

首先把tars协议文件转化为Golang语言形式

```shell
[root@kaiwen kevin]# cd HelloGo/
[root@kaiwen HelloGo]# tars2go  -outdir=tars-protocol -module=github.com/Tars/test SayHello.tars
```

现在开始实现服务端的逻辑：客户端传来一个名字，服务端回应hello name。

```go
// HelloGo/SayHello_imp.go

package main

import "context"

type SayHelloImp struct {
}

func (imp *SayHelloImp) EchoHello(ctx context.Context, name string, greeting *string) (int32, error) {
    *greeting = "hello " + name
    return 0, nil
}
```

**注意**： 这里函数名要大写，Go语言方法导出规定。

编译main函数，初始代码以及有tars框架实现了。

```go
// HelloGo/main.go

package main

import (
    "github.com/TarsCloud/TarsGo/tars"
    
    "github.com/Tars/test/tars-protocol/TestApp"
)

func main() {
    // Get server config
    cfg := tars.GetServerConfig()
  
    // New servant imp
    imp := new(SayHelloImp)
    // New servant
    app := new(TestApp.SayHello)
    // Register Servant
    app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".SayHelloObj")
  
    // Run application
    tars.Run()
}
```

编译生成可执行文件，并打包发布包。将生成可执行文件HelloGo和发布包HelloGo.tgz

```shell
cd HelloGo && make && make tar
```

运行服务端

```sh
[root@kaiwen HelloGo]# ./HelloGo --config=config.conf
# OR 后台运行
[root@kaiwen HelloGo]# ./HelloGo --config=config.conf &
[1] 11970
```

### 客户端开发

```go
package main

import (
    "fmt"
  
    "github.com/TarsCloud/TarsGo/tars"
  
    "github.com/Tars/test/tars-protocol/TestApp"
)

//只需初始化一次，全局的
var comm *tars.Communicator
func main() {
    comm = tars.NewCommunicator()
    obj := "TestApp.HelloGo.SayHelloObj@tcp -h 127.0.0.1 -p 10015 -t 60000"
    app := new(TestApp.SayHello)
    /*
       // if your service has been registered at tars registry
       obj := "TestApp.HelloGo.SayHelloObj"
       // tarsregistry service at 192.168.1.1:17890
       comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 192.168.1.1 -p 17890")
    */
  
    comm.StringToProxy(obj, app)
    reqStr := "tars"
    var resp string
    ret, err := app.EchoHello(reqStr, &resp)
    if err != nil {
      fmt.Println(err)
      return
    }
    fmt.Println("ret: ", ret, "resp: ", resp)
}
```

- TestApp依赖是tars2go生成的代码。

- obj指定服务端地址端口，如果服务端未在主控注册，则需要知道服务端的地址和端口并在Obj中指定，在例子中，协议为TCP，服务端地址为本地地址，端口为3002。如果有多个服务端，则可以这样写`TestApp.HelloGo.SayHelloObj@tcp -h 127.0.0.1 -p 9985:tcp -h 192.168.1.1 -p 9983`这样请求可以分散到多个节点。

  如果已经在主控注册了服务，则不需要写死服务端地址和端口，但在初始化通信器时需要指定主控的地址。

- com通信器，用于与服务端通信。

编译测试

```shell
[root@kaiwen HelloGo]# go build -o client/client client/client.go
[root@kaiwen HelloGo]# ./client/client
ret:  0 resp:  hello tarss
```



