package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testdash/models"
	"time"
)

func GetMetrics(c *gin.Context) {
	//post
	var metricsQuery models.MetricsQuery
	if err := c.ShouldBindJSON(&metricsQuery); err != nil {
		klog.V(2).ErrorS(err, "bind models.MetricsQuery to json failed")
		//WriteOK(c, gin.H{})
		WriteError(c,err.Error())
		return
	}

	//todo,从数据库中获取Prometheus的服务地址
	//判断服务是否可用（配置超时时间）
	ctx, cancel := context.WithTimeout(context.Background(),time.Second * 10)
	defer cancel()
	readyReq, err := http.NewRequest("GET", "http://172.21.204.62:30250/-/ready",nil)
	if err != nil {
		klog.V(2).ErrorS(err, "check prometheus ready failed")
		//WriteOK(c,gin.H{})
		WriteError(c,err.Error())
		return
	}
	readyResp, err := http.DefaultClient.Do(readyReq.WithContext(ctx))
	if err != nil {
		klog.V(2).ErrorS(err,"check prometheus response failed")
		//WriteOK(c,gin.H{})
		WriteError(c,err.Error())
		return
	}
	if readyResp.StatusCode != http.StatusOK {
		//WriteOK(c,gin.H{})
		WriteError(c,err.Error())
		return
	}

	wg := sync.WaitGroup{}

	step := 60
	end := time.Now().Unix()
	start := end - 3600

	tracker := models.NewPrometheusTracker()

	//执行查询任务
	e := reflect.ValueOf(&metricsQuery).Elem()
	for i := 0; i < e.NumField(); i++ {
		wg.Add(i)
		go func(i int) {
			//执行promeql查询
			defer wg.Done()
			fName := e.Type().Field(i).Name
			fValue := e.Field(i).Interface().(*models.MetricsCategory)
			fTag := e.Type().Field(i).Tag
			if fValue == nil {
				return
			}
            klog.V(2).Info("start query prometheus data", "field",fName)

			prometheusQuery := fValue.GenerateQuery()
			if prometheusQuery == nil {
				klog.V(2).Info("no promql", "field",fName)
				return
			}

			promql := prometheusQuery.GetValueByField(fName)

            resp, err := http.Get(fmt.Sprintf("http://172.21.204.62:30250/api/v1/query_range?query=%s&start=%d&end=%dstep=%d", promql, start, end, step))
            if err != nil {
            	klog.V(2).ErrorS(err, "request prometheus failed","field",fName, "promql", promql)
				return
			}

			body ,err := ioutil.ReadAll(resp.Body)
			defer func(Body io.ReadCloser) {
				_ =Body.Close()
			}(resp.Body)
			if err != nil {
				klog.V(2).ErrorS(err, "read resp body error","field",fName, "promql", promql)
			}

			var prometheusResp models.PrometheusQueryResp
			if err := json.Unmarshal(body,&prometheusResp); err != nil {
				klog.V(2).ErrorS(err,"unmarshal prometheus body error","field",fName, "promql", promql)
			}

			//把数据组装到要返回的json中
			tag := fTag.Get("json")
			tracker.Set(tag[:strings.Index(tag,",omitempty")], &prometheusResp)

		}(i)
	}

	//等待所有的查询完成
	wg.Wait()

	WriteOK(c, gin.H{
		"metrics": tracker.Metrics,
	})

}
