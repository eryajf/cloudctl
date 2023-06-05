package cls

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/eryajf/cloudctl/public"
	jsoniter "github.com/json-iterator/go"
	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"gopkg.in/yaml.v2"
)

type LogTopicCost struct {
	Response struct {
		Analysis        bool          `json:"Analysis"`
		AnalysisRecords []interface{} `json:"AnalysisRecords"`
		AnalysisResults []struct {
			Data []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"Data"`
		} `json:"AnalysisResults"`
		ColNames     []string      `json:"ColNames"`
		Columns      []interface{} `json:"Columns"`
		Context      string        `json:"Context"`
		ListOver     bool          `json:"ListOver"`
		RequestID    string        `json:"RequestId"`
		Results      []interface{} `json:"Results"`
		SamplingRate int           `json:"SamplingRate"`
	} `json:"Response"`
}

func GetRst(region, logid, query string) (rst float64) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cls.tencentcloudapi.com"
	client, _ := cls.NewClient(GetSecret(), region, cpf)

	request := cls.NewSearchLogRequest()

	request.TopicId = common.StringPtr(logid)
	request.From = common.Int64Ptr(public.GetOldTime())
	request.To = common.Int64Ptr(public.GetNowTime())
	request.Query = common.StringPtr(query)
	request.SyntaxRule = common.Uint64Ptr(1)

	response, err := client.SearchLog(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return
	}
	if err != nil {
		fmt.Printf("do cls search failed: %v\n", err)
		return
	}
	var data LogTopicCost
	err = jsoniter.Unmarshal([]byte(response.ToJsonString()), &data)
	if err != nil {
		return
	}
	for _, v := range data.Response.AnalysisResults {
		ret, _ := strconv.ParseFloat(v.Data[0].Value, 64)
		rst = math.Round(ret*100) / 100 // 保留2位小数
	}
	return rst
}

type DetailSetRsp struct {
	Response struct {
		Context    string      `json:"Context"`
		DetailSets []DetailSet `json:"DetailSet"`
		RequestID  string      `json:"RequestId"`
		Total      int         `json:"Total"`
	} `json:"Response"`
}

type DetailSet struct {
	ActionType       string `json:"ActionType"`
	ActionTypeName   string `json:"ActionTypeName"`
	BillID           string `json:"BillId"`
	BusinessCode     string `json:"BusinessCode"`
	BusinessCodeName string `json:"BusinessCodeName"`
	ComponentSet     []struct {
		BlendedDiscount    string      `json:"BlendedDiscount"`
		CashPayAmount      string      `json:"CashPayAmount"`
		ComponentCode      string      `json:"ComponentCode"`
		ComponentCodeName  string      `json:"ComponentCodeName"`
		ContractPrice      string      `json:"ContractPrice"`
		Cost               string      `json:"Cost"`
		Discount           string      `json:"Discount"`
		IncentivePayAmount string      `json:"IncentivePayAmount"`
		InstanceType       string      `json:"InstanceType"`
		ItemCode           string      `json:"ItemCode"`
		ItemCodeName       string      `json:"ItemCodeName"`
		OriginalCostWithRI string      `json:"OriginalCostWithRI"`
		OriginalCostWithSP string      `json:"OriginalCostWithSP"`
		PriceUnit          string      `json:"PriceUnit"`
		RealCost           string      `json:"RealCost"`
		ReduceType         string      `json:"ReduceType"`
		RiTimeSpan         string      `json:"RiTimeSpan"`
		SPDeductionRate    string      `json:"SPDeductionRate"`
		SinglePrice        string      `json:"SinglePrice"`
		TimeSpan           string      `json:"TimeSpan"`
		TimeUnitName       string      `json:"TimeUnitName"`
		TransferPayAmount  interface{} `json:"TransferPayAmount"`
		UsedAmount         string      `json:"UsedAmount"`
		UsedAmountUnit     string      `json:"UsedAmountUnit"`
		VoucherPayAmount   string      `json:"VoucherPayAmount"`
	} `json:"ComponentSet"`
	FeeBeginTime    string        `json:"FeeBeginTime"`
	FeeEndTime      string        `json:"FeeEndTime"`
	OperateUin      string        `json:"OperateUin"`
	OrderID         string        `json:"OrderId"`
	OwnerUin        string        `json:"OwnerUin"`
	PayModeName     string        `json:"PayModeName"`
	PayTime         string        `json:"PayTime"`
	PayerUin        string        `json:"PayerUin"`
	PriceInfo       []interface{} `json:"PriceInfo"`
	ProductCode     string        `json:"ProductCode"`
	ProductCodeName string        `json:"ProductCodeName"`
	ProjectID       int           `json:"ProjectId"`
	ProjectName     string        `json:"ProjectName"`
	RegionID        string        `json:"RegionId"`
	RegionName      string        `json:"RegionName"`
	ResourceID      string        `json:"ResourceId"`
	ResourceName    string        `json:"ResourceName"`
	Tags            []interface{} `json:"Tags"`
	ZoneName        string        `json:"ZoneName"`
}

func GetClsAlarms(region, logid string) ([]DetailSet, error) {
	var (
		offset uint64 = 0
		limit  uint64 = 100
		tmps   []DetailSet
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "billing.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := billing.NewClient(GetSecret(), region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := billing.NewDescribeBillDetailRequest()

	request.Offset = &offset
	request.Limit = &limit

	now := time.Now()
	beforeYesterday := now.Add(-48 * time.Hour)
	beforeYesterdayStart := time.Date(beforeYesterday.Year(), beforeYesterday.Month(), beforeYesterday.Day(), 0, 0, 0, 0, beforeYesterday.Location())
	beforeYesterdayEnd := time.Date(beforeYesterday.Year(), beforeYesterday.Month(), beforeYesterday.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), beforeYesterday.Location())

	request.BeginTime = common.StringPtr(beforeYesterdayStart.Format("2006-01-02 15:04:05"))
	request.EndTime = common.StringPtr(beforeYesterdayEnd.Format("2006-01-02 15:04:05"))
	request.NeedRecordNum = common.Int64Ptr(1)
	request.BusinessCode = common.StringPtr("p_cls")
	request.ResourceId = common.StringPtr(logid)

	for {
		response, err := client.DescribeBillDetail(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s", err)
			return nil, err
		}
		if err != nil {
			fmt.Printf("get list bill err: %v", err)
			return nil, err
		}

		var data DetailSetRsp
		err = jsoniter.Unmarshal([]byte(response.ToJsonString()), &data)
		if err != nil {
			fmt.Printf("json unmarshal err: %v", err)
			return nil, err
		}
		// 输出json格式的字符串回包
		for _, v := range data.Response.DetailSets {
			tmp := DetailSet{
				ActionType:       v.ActionType,
				ActionTypeName:   v.ActionTypeName,
				BillID:           v.BillID,
				BusinessCode:     v.BusinessCode,
				BusinessCodeName: v.BusinessCodeName,
				ComponentSet:     v.ComponentSet,
				FeeBeginTime:     v.FeeBeginTime,
				FeeEndTime:       v.FeeEndTime,
				OperateUin:       v.OperateUin,
				OrderID:          v.OrderID,
				OwnerUin:         v.OwnerUin,
				PayModeName:      v.PayModeName,
				PayTime:          v.PayTime,
				PayerUin:         v.PayerUin,
				PriceInfo:        v.PriceInfo,
				ProductCode:      v.ProductCode,
				ProductCodeName:  v.ProductCodeName,
				ProjectID:        v.ProjectID,
				ProjectName:      v.ProjectName,
				RegionID:         v.RegionID,
				RegionName:       v.RegionName,
				ResourceID:       v.ResourceID,
				ResourceName:     v.ResourceName,
				Tags:             v.Tags,
				ZoneName:         v.ZoneName,
			}
			tmps = append(tmps, tmp)
		}
		if len(data.Response.DetailSets) == 0 {
			break
		}
		offset += limit
	}
	return tmps, nil
}

type ProjectList struct {
	Bot          string        `yaml:"bot"`
	ProjectItems []ProjectItem `yaml:"project_items"`
}

type ProjectItem struct {
	LogID   string   `yaml:"logid"`
	Kind    string   `yaml:"kind"`
	Project []string `yaml:"project"`
}

func LoadProjectList(file string) (rst *ProjectList, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &rst)
	if err != nil {
		return
	}
	return rst, nil
}

func getQuery(kind, project string) (query string) {
	switch kind {
	case "host":
		query = `* | SELECT COUNT(*) * 100.0 / (SELECT COUNT(*)) AS percentage WHERE "__HOSTNAME__" like '` + project + `%'`
	case "k8s":
		query = `* | SELECT ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*)), 2) AS percentage WHERE "__TAG__.pod_label_app" = '` + project + `'`
	}
	return
}

func getLogCost(region, logid string) (cost int) {
	alarms, err := GetClsAlarms(region, logid)
	if err != nil {
		fmt.Println("get alarm err: ", err)
	}
	var tmpint []int
	for _, v := range alarms {
		for _, j := range v.ComponentSet {
			a, _ := strconv.Atoi(strings.Split(j.RealCost, ".")[0])
			tmpint = append(tmpint, a)
		}
	}
	return public.Sum(tmpint)
}
