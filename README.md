# 概述

> 基于ssh的能够在远程主机组上批量执行命令或上传下载文件的命令行工具（跨平台、无依赖、免安装）。

> 在线帮助：[https://xied5531.github.io/xsh/](https://xied5531.github.io/xsh/)

# 原理

```
                     ----------                          ----------    
                     | config |  <--------|        --->  | group1 |    
                     ----------           |        |     ----------    
                                          |        |                    ---> command
   ----------        ----------        -------     |     ----------     |
   | auth   |  <---  | host   |  <---  | xsh |  ------>  | group2 |  ---|
   ----------        ----------        -------     |     ----------     |
                                                   |                    ---> copy
                                                   |     ----------        
                                                   --->  | group3 |
                                                         ----------           
```

- config，控制xsh行为
- auth，定义认证信息
- host，定义主机组信息并关联auth信息
- xsh，根据指定的group主机组，建立ssh连接，执行命令或传输文件

# 核心特性

- 运行模式：支持单命令行、任务编排、交互式常驻操作
- 执行命令：支持单命令，多命令；支持普通用户执行、切换其他用户执行
- 传输文件：支持文件和文件夹；支持上传和下载；协议：scp
- 任务编排：支持行为组合；支持自定义变量
- 结果输出：支持text、json、yaml
- 安全性：支持配置命令黑名单、密码加密保存、密码和私钥认证、超时控制

# 依赖

| 第三方包 | License |
| ------ | ------ |
| github.com/patrickmn/go-cache | MIT |
| github.com/peterh/liner | X11 |
| github.com/hnakamur/go-scp | MIT |
| golang.org/x/crypto | BSD |
| gopkg.in/yaml.v2 | Apache |
