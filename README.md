<div align="center">
<h1>cloudctl</h1>

[![Auth](https://img.shields.io/badge/Auth-eryajf-ff69b4)](https://github.com/eryajf)
[![GitHub Pull Requests](https://img.shields.io/github/stars/eryajf/cloudctl)](https://github.com/eryajf/cloudctl/stargazers)
[![HitCount](https://views.whatilearened.today/views/github/eryajf/cloudctl.svg)](https://github.com/eryajf/cloudctl)
[![GitHub license](https://img.shields.io/github/license/eryajf/cloudctl)](https://github.com/eryajf/cloudctl/blob/main/LICENSE)
[![](https://img.shields.io/badge/Awesome-MyStarList-c780fa?logo=Awesome-Lists)](https://github.com/eryajf/awesome-stars-eryajf#readme)

<p> 🌉 与云产品的一些交互 🌉</p>

<img src="https://camo.githubusercontent.com/82291b0fe831bfc6781e07fc5090cbd0a8b912bb8b8d4fec0696c881834f81ac/68747470733a2f2f70726f626f742e6d656469612f394575424971676170492e676966" width="800"  height="3">

</div>

运维也可以如此优雅！

## 如何使用

先拷贝配置文件，然后正确配置里边的内容：

```sh
cp config.example.yml config.yml
```

执行如下指令，查询某日志主题前天单日费用：

```sh
# 你也可以直接在release下载二进制来执行
$ cloudctl cls getlogcost -l "xxxxxxxxxxxxxxxxxx" -r "ap-shanghai"

前天的花费为: 174元
日志主题为: xxxxxxxxxxxxxxxxxx
```

查询某个应用在一个大日志主题中花费的费用。

```sh
$ cloudctl cls getprojectcost -b "bot_webhook" -k k8s -l xxxxxxxxxxxxxxxxxxxxxxx -p testeryajf -r ap-shanghai

> 应用名称：`testeryajf`
> 应用在主题日志占比： `10%`
> 应用单月费用预计： `3000元`
>该主题前天单日的费用为: `1000元`
该主题的ID为: `xxxxxxxxxxxxxxxxxxxxxxx`
```

通知到群里的消息大概长这样：

![](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20230605_184855.png)

> 我们的应用场景是，生产环境中，有两个日志主题，一个采集所有 CVM 主机应用的日志，一个采集所有部署在 tke 之中的日志。所以在一个日志主题中，会有许多个应用，每个应用对应了一组服务器。当然在 k8s 中对应的名字就是 deployment 的名字。

上边所有的参数都比较容易理解，示例当中也给出了明确的解析。这里单独把 project 的参数拿出来说一说。传递应用的名字，程序会首先拼接出如下语句，查询应用在整个日志主题中的日志量占比：

- CVM：`* | SELECT ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*)), 2) AS percentage WHERE "__HOSTNAME__" like 'testeryajf%'`
- K8S: `* | SELECT ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*)), 2) AS percentage WHERE "__TAG__.pod_label_app"  'testeryajf'`

所以你在运行项目之前，应该先确认如上语句，在你的日志主题当中，能够正常拿到该应用在主题中的日志量百分比。这是算出应用一个月成本评估的基础。

如果你有多个项目想要共同关注，则可以参照 project_list.example.yml 文件，对内容进行配置，运行命令如下：

```sh
$ cloudctl cls getprojectcost -f project_list.example.yml
```

你也可以使用 docker 镜像来运行服务。

```sh
$ docker run -it -e TC_SERCRET_ID="xxxxx" -e TC_SERCRET_KEY="xxxxxxxxxx" dockerproxy.com/eryajf/cloudctl /app/cloudctl cls getlogcost  -l "222345432454425"

前天的花费为: 1000元
日志主题为: 222345432454425
```

## 感谢开源

- [eryajfctl](https://github.com/eryajf/eryajfctl)

如果觉得项目不错，请别忘了一键三连，给个 star。
