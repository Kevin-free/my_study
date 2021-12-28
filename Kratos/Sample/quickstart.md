# 快速开始

### Requirements 

- Go version>=1.13
- 设置环境变量：
- 开启go mod：`GO111MODULE=on`
- 设置包下载代理：`GOPROXY=https://goproxy.cn,direct`
- 代理忽略公司git地址：`GONOPROXY=git.huoys.com`
- 关闭校验：`GOSUMDB=off`



### Installation

1. **安装工具：**

   - 通过go get安装：

     ```
     go get -u git.huoys.com/middle-end/tool/kratos
     ```

   - 安装全部工具

     ```
     kratos tool install all
     ```

2. **生成基于kratos库的脚手架工程：**

   - 一键生成服务工程项目：

     ```
     kratos new [servcieName] --service
     kratos new [servcieName] （实测用此命令，not defined: -service）
     ```

   - 一键生成游戏工程项目：

     ```
     kratos new [gameName] --game
     ```



### Build & Run

```
cd [servcieName]/cmd
go build
./cmd -conf ../configs  （Linux 中）
cmd.exe -conf ../configs（Windows 中）
```

打开浏览器访问：http://localhost:8000/[servcieName]/start，你会看到输出了`Golang 大法好 ！！！`



### Protobuf 协议

目前kratos tool可基于protobuf协议生成http、grpc、tcp接口代码。

1. **服务开发**：生成http、grpc代码，提供http、grp协议接口

   ```
   kratos tool protoc --service api.proto
   ```

   服务协议生成成功后，需将协议文件上传到git：https://git.huoys.com/middle-end/proto，方便服务调用方引入使用，服务调用方import proto库，直接使用proto库中的client.go初始化grpc client即可。

2. **游戏开发**：生成http、tcp代码，提供http、tcp协议接口

   ```
   kratos tool protoc --game api.proto
   ```

   关于tcp协议定义：

   - command的enum名称必须为GameCommand，错误码的enum名称必须为ErrCode，不可修改！！！
   - 每个command的名称与定义rpc的函数名称一致！！！

```
   service Demo {
       rpc Ping (.google.protobuf.Empty) returns (.google.protobuf.Empty);
   	rpc SayHello(HelloReq) returns (HelloResp) {
           option (google.api.http) = {
               get:"/Demo/SayHello"
           };
       };
   }
   
   enum GameCommand {
       Ping = 0;
       SayHello = 1;
   }
   
   enum ErrCode {
       HelloSuccess = 0;
       HelloFailed = 1500;
   }
```

**多个proto文件支持：**

proto协议文件可以新建多个，其中必须有一个main文件，其他文件为子模块

例：game_message_main.proto、game_message_ping.proto、game_message_login.proto

game_message_main.proto引入子文件：

```
syntax = "proto2";
package GameProto;

import "github.com/gogo/protobuf/gogoproto/gogo.proto";
import "google/api/annotations.proto";
import "game_message_ping.proto";
import "game_message_login.proto";

service DjsTf {
	rpc Req_LogoutGame(CSLogoutGame) returns (SCLogoutGame) {
        option (google.api.http) = {
            post:"/DjsTf/Req_LogoutGame"
        };
    };
}
enum GameCommand {
	Req_RequestAttack 				= 1004; 
}
enum ErrCode {
	QPBaseUserInfoGetFailed 	= 10002;
}
```

game_message_login.proto:

```
package GameProto;

import "github.com/gogo/protobuf/gogoproto/gogo.proto";
import "google/api/annotations.proto";
import "game_message_main.proto";

service DjsTf {
	rpc Req_LoginGameWithToken(CSLoginGameWithToken) returns (SCLoginGameWithToken) {
		option (google.api.http) = {
            post:"/DjsTf/Req_LoginGameWithToken"
        };
	};
}
enum GameCommand {
	Req_LoginGameWithToken = 1001; //登录
	Req_LogoutGame = 1002; //退出
	Push_LogoutGame = 1003; //退出结果
}

enum ErrCode {
	QPEnterRoomFailed = 10001; //请求房间失败，你在其他房间的游戏还没结束
}

// 登录
message CSLoginGameWithToken {
	required int64 PlayerId = 1;
	required string Token = 2;
	optional string ClientVersion = 3;	//游戏版本号
}

message SCLoginGameWithToken {
	required int64 Money = 1;			//玩家登录时的金币
	repeated HeroInfo Info = 2;			//英雄信息
	required int64 PassNo = 3;			//玩家当前关卡等级
	required int32 IsFirstLogin = 4;	//第一次登陆推送(0:未推送,1:已推送)
	required int32 IsForceGuide = 5;	//是否已完成强制引导(0:未完成,1:已完成)
	optional int32 UseRatio = 6;		//当前使用的倍率
}
```

协议接口生成命令：

```
kratos tool protoc --game game_message_main.proto
```

**注：业务逻辑实现在internal/service目录下，不要写在service.go文件里，另新建一个.go文件即可。**



### 项目目录结构

```
├── CHANGELOG.md 
├── OWNERS
├── README.md
├── api                     # api目录为对外保留的proto文件及生成的pb.go文件
│   ├── api.bm.go
│   ├── api.pb.go           # 通过go generate生成的pb.go文件
│   ├── api.proto
│   └── client.go
├── cmd
│   └── main.go             # cmd目录为main所在
├── configs                 # configs为配置文件目录
│   ├── application.toml    # 应用的自定义配置文件，可能是一些业务开关如：useABtest = true
│   ├── db.toml             # db相关配置
│   ├── grpc.toml           # grpc相关配置
│   ├── http.toml           # http相关配置
│   ├── memcache.toml       # memcache相关配置
│   └── redis.toml          # redis相关配置
├── go.mod
├── go.sum
└── internal                # internal为项目内部包，包括以下目录：
│   ├── dao                 # dao层，用于数据库、cache、MQ、依赖某业务grpc|http等资源访问
│   │   ├── dao.bts.go
│   │   ├── dao.go
│   │   ├── db.go
│   │   ├── mc.cache.go
│   │   ├── mc.go
│   │   └── redis.go
│   ├── di                  # 依赖注入层 采用wire静态分析依赖
│   │   ├── app.go
│   │   ├── wire.go         # wire 声明
│   │   └── wire_gen.go     # go generate 生成的代码
│   ├── model               # model层，用于声明业务结构体
│   │   └── model.go
│   ├── server              # server层，用于初始化grpc和http server
│   │   ├── grpc            # grpc层，用于初始化grpc server和定义method
│   │   │   └── server.go
│   │   └── http            # http层，用于初始化http server和声明handler
│   │       └── server.go
│   └── service             # service层，用于业务逻辑处理，且为方便http和grpc共用方法，建议入参和出参保持grpc风格，且使用pb文件生成代码
│       └── service.go
└── test                    # 测试资源层 用于存放测试相关资源数据 如docker-compose配置 数据库初始化语句等
    └── docker-compose.yaml
```





-------------

[文档目录树](https://git.huoys.com/middle-end/kratos/blob/master/doc/wiki-cn/summary.md)