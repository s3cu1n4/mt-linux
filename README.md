# mt-linux

一款基于 Elkeid(AgentSmith-HIDS) Driver 的实时 Linux 信息采集软件，使用AgentSmith-HIDS 的预编译内核模块，可快速收集Linux系统内核级进程执行，特权升级监控，网络审计等信息，在没有部署Linux HIDS产品的情况下，可使用该软件进行应急响应排查。

## 使用方式

本软件提供了编译后的可执行文件，可直接下载本软件使用。

若安装AgentSmith-HIDS Driver或下载失败，请自行编译内核模块,
编译方式可参考：

[About Elkeid(AgentSmith-HIDS) Driver](https://github.com/bytedance/Elkeid/blob/main/driver/README-zh_CN.md)

编译完成后执行：
```shell script
insmod hids_driver.ko
```
加载完内核模块，再运行本软件即可




## 其他
注意：使用 ` insmod ` 方式手动加载的内核模块，本软件不会主动卸载，需要自己使用 

```shell script
rmmod hids_driver
```
 卸载内核模块

