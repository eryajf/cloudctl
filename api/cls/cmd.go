package cls

import (
	"fmt"
	"os"
	"time"

	"github.com/eryajf/cloudctl/public"
	"github.com/spf13/cobra"
)

var GetProjectCostCmd = &cobra.Command{
	Use:   "getprojectcost",
	Short: "获取一个应用所花的费用",
	Long: `通过命令行获取一个应用在cls中所消耗的费用，需要传入的参数有：
  -l: 应用所在日志主题的ID
  -p: 应用主机名，如果你检索的日志主题为CVM，其中一组服务器名为 prod-admin-node-1..4,那么此处应该写为 prod-admin-node
  -b: 指定机器人的webhook
  -r: 指定日志所在的region
  -k: 指定应用所在日志主题的类型，host/k8s`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		// 判断是否有传入标志
		if flags.NFlag() == 0 {
			fmt.Println("Please specify a flag")
			os.Exit(1)
		}
		var (
			msg string
			bot string
		)
		file, _ := cmd.Flags().GetString("file")
		region, _ := cmd.Flags().GetString("region")
		if file != "" {
			rst, err := LoadProjectList(file)
			if err != nil {
				fmt.Println("err:", err)
			}
			bot = rst.Bot
			for _, v := range rst.ProjectItems {
				logid := v.LogID
				sumcost := getLogCost(region, logid)
				msg += fmt.Sprintf("\n\n> 主题前日花费：`%d元`\n> 主题ID：`%s`\n", sumcost, logid)
				for _, j := range v.Project {
					percent := GetRst(region, logid, getQuery(v.Kind, j))
					project_cost := fmt.Sprintf("%.2f", float64(sumcost)*percent/100*30)
					msg += fmt.Sprintf("> **应用：** `%s`\n> 日志量占比： `%v%%` \n> 单月花费预计： `%v元`", j, percent, project_cost) + "\n"
				}
			}
		} else {
			logid, _ := cmd.Flags().GetString("logid")
			project, _ := cmd.Flags().GetString("project")
			bot, _ = cmd.Flags().GetString("bot")
			kind, _ := cmd.Flags().GetString("kind")
			sumcost := getLogCost(region, logid)
			percent := GetRst(region, logid, getQuery(kind, project))
			project_cost := fmt.Sprintf("%.2f", float64(sumcost)*percent/100*30)
			msg = fmt.Sprintf("> 应用名称：`%s`\n> 应用在主题日志占比： `%v%%` \n> 应用单月费用预计： `%v元`\n>该主题前天单日的费用为: `%v元`\n该主题的ID为: `%s`", project, percent, project_cost, sumcost, logid)
		}
		fmt.Println(msg)
		err := public.SendMsg(bot, fmt.Sprintf("=========`%v`汇总数据如下=========\n\n\n\n%s\n\n\n\n 说明：目前取样该应用的数据为最近三个小时占该主题的占比，金额则是取的该主题前天整天的费用，计算方式为 `前天单日费用 x 近三小时百分比 x 30天` 得出该应用单月的费用，是一个相对粗略的参考，但排除日志主题新增或减少应用的变更因素，这个计算数值是相对准确的。", time.Now().Format("2006-01-02 15:04"), msg))
		if err != nil {
			fmt.Printf("发送机器人消息失败: %v\n", err)
		}
	},
}

var GetLogCostCmd = &cobra.Command{
	Use:   "getlogcost",
	Short: "获取一个日志主题在前天的花费总额",
	Long: `通过命令行获取一个日志主题在前天所消耗的费用，需要传入的参数有：
  -l: 应用所在日志主题的ID`,
	Run: func(cmd *cobra.Command, args []string) {
		logid, _ := cmd.Flags().GetString("logid")
		region, _ := cmd.Flags().GetString("region")
		fmt.Printf("前天的花费为: %d元\n日志主题为: %s\n", getLogCost(region, logid), logid)
	},
}
