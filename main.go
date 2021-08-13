package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"testdash/global"
	"testdash/routers"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()
	flag.Set("logtostderr","false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()


	if err := global.Init(); err != nil {
		klog.V(2).ErrorS(err, "init Global failed")
		return
	}

	serv := gin.Default()

	//todo:注册路由
	routers.InitApi(serv)

	if err := serv.Run(":8888"); err != nil {
		panic(err)
		klog.V(2).ErrorS(err, "Server Run Error")
	}
}
