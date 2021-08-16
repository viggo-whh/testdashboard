package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"testdash/k8s"
	"testdash/routers"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()
	_ = flag.Set("logtostderr", "false")
	_ = flag.Set("alsologtostderr", "false")
	flag.Parse()

	if err := initialize(); err != nil {
		klog.V(2).ErrorS(err, "init Global failed")
		return
	}

	serv := gin.Default()

	//todo:注册路由
	routers.InitApi(serv)

	if err := serv.Run(":8888"); err != nil {
		klog.V(2).ErrorS(err, "Server Run Error")
		panic(err)
	}
}

var err error

func initialize() error {
	if err := k8s.NewKubeClient(); err != nil {
		return err
	}
	return nil
}
