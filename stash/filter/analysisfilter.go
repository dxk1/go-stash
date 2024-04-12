/**
 * Package filter
 * @Author: tbb
 * @Date: 2024/3/5 11:55
 */
package filter

import (
	"fmt"
)

// DataInfo ...
type DataInfo struct {
	Region       string
	Cluster      string
	FunctionId   string
	AppId        string
	Namespace    string
	FunctionName string
	ErrorCode    *map[string]int64
}

var data = make(map[string]*DataInfo)

// AnalysisFilter 函数异常查询
func AnalysisFilter() FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		Analysis(m)
		return m
	}
}

func Analysis(m map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Errorf("analysis panic:%v", r))
		}
	}()

	//if m["Module"].(string) != "RegionInvoke" {
	//	return
	//}
	//statusCode, ok := m["StatusCode"].(int64)
	//if !ok {
	//	fmt.Println(fmt.Sprintf("StatusCode err: %+v", m["StatusCode"]))
	//	return
	//}
	//
	//if statusCode == 500 {
	//	reqCntVal, _ := m["ReqCnt"].(int64)
	//	functionIdVal, _ := m["FunctionId"].(string)
	//	appIdVal, _ := m["AppId"].(string)
	//	functionNameVal, _ := m["FunctionName"].(string)
	//	namespaceVal, _ := m["Namespace"].(string)
	//	clusterVal, _ := m["CalleeCluster"].(string)
	//	timestampVal, _ := m["Timestamp"].(float64)
	//	regionVal, _ := m["Region"].(string)
	//	functionTypeVal, _ := m["FunctionType"].(string)
	//
	//	//分钟级别聚合
	//	min := time.Unix(int64(timestampVal), 0).Minute()
	//	key := fmt.Sprintf("%s:%s:%s:%d", regionVal, clusterVal, functionIdVal, min)
	//	if d, ok := data[key]; ok {
	//		if errReqNum, ok := (*d.ErrorCode)[functionTypeVal]; ok {
	//			(*d.ErrorCode)[functionTypeVal] = errReqNum + reqCntVal
	//		} else {
	//			(*d.ErrorCode)[functionTypeVal] = reqCntVal
	//		}
	//	} else {
	//		errorCode := make(map[string]int64)
	//		errorCode[functionTypeVal] = reqCntVal
	//		data[key] = &DataInfo{
	//			Region:       regionVal,
	//			Cluster:      clusterVal,
	//			FunctionId:   functionIdVal,
	//			Namespace:    namespaceVal,
	//			AppId:        appIdVal,
	//			FunctionName: functionNameVal,
	//			ErrorCode:    &errorCode,
	//		}
	//	}
	//}
}
