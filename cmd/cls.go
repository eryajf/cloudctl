/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/eryajf/cloudctl/api/cls"
	"github.com/spf13/cobra"
)

// exCmd represents the jenkins command
var clsCmd = &cobra.Command{
	Use:  "cls",
	Long: `与CLS产品相关的一些交互。`,
}

func init() {
	rootCmd.AddCommand(clsCmd)
	// 获取一个应用在cls中的花费预估
	clsCmd.AddCommand(cls.GetProjectCostCmd)
	cset := cls.GetProjectCostCmd.Flags()
	cset.StringP("logid", "l", "", "应用所在日志主题的ID")
	cset.StringP("project", "p", "", "应用主机名，如果你检索的日志主题为CVM，其中一组服务器名为 prod-admin-node-1..4,那么此处应该写为 prod-admin-node")
	cset.StringP("bot", "b", "", "结果同步的群机器人，目前仅支持企业微信群")
	cset.StringP("kind", "k", "host", "传入host/k8s，默认为host，不同的值会使用不同的查询语句")
	cset.StringP("region", "r", "ap-shanghai", "传入日志所在的区域，默认为上海")
	cset.StringP("file", "f", "", "指定要关注的应用配置文件，当指定了文件时，其他参数将不需要指定。")
	clsCmd.AddCommand(cls.GetLogCostCmd)
	dset := cls.GetLogCostCmd.Flags()
	dset.StringP("logid", "l", "", "应用所在日志主题的ID")
	dset.StringP("region", "r", "ap-shanghai", "传入日志所在的区域，默认为上海")
	_ = cls.GetLogCostCmd.MarkFlagRequired("logid")
}
