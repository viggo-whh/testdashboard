package controllers

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"testdash/global"
	"testdash/k8s"
)

func GetNodeList(ctx *gin.Context) {
	nodes, err :=k8s.GetNodes(global.K8sClient(),"")
	if err != nil {
		klog.V(2).ErrorS(err,"Get Nodes Failed", "controllers","GetNodeList")
		WriteError(ctx,"Get Nodes Failed")
		return
	}
	WriteOK(ctx,gin.H{"nodes":nodes})
}
