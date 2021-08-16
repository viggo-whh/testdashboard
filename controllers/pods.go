package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"net/http"
	"strconv"
	"testdash/k8s"
)

func GetKubeLogs(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("pod")
	container := c.Query("container")
	tailLines, _ := strconv.ParseInt(c.DefaultQuery("tailLines", "500"), 10, 64)
	timestamps, _ := strconv.ParseBool(c.DefaultQuery("timestamps", "true"))
	previous, _ := strconv.ParseBool(c.DefaultQuery("previous", "false"))

	klog.V(2).InfoS("get kube logs request params", "namespace", namespace, "pod", podName, "container", container, "tailLines", tailLines, "timestamps", timestamps, "previous", previous)

	if namespace == "" || podName == "" || container == "" {
		c.String(http.StatusBadRequest, "must spectific namespace,pod and container query params")
		return
	}

	//获取pod日志
	kubelogger, err := NewKubeLogger(c.Writer, c.Request, nil)
	if err != nil {
		klog.V(2).ErrorS(err, "upgrade websocket failed")
		c.String(http.StatusBadRequest, err.Error())
	}

	opts := corev1.PodLogOptions{
		Container:  container,
		Follow:     true,
		Previous:   previous,
		Timestamps: timestamps,
		TailLines:  &tailLines,
	}

	//req := global.K8sClient().CoreV1().Pods(namespace).GetLogs(podName,&opts)
	if err := k8s.Client.Pod.LogsStream(podName, namespace, &opts, kubelogger); err != nil {
		klog.V(2).ErrorS(err, "Get log stream failed")
		_, _ = kubelogger.Write([]byte(err.Error()))
	}
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type KubeLogger struct {
	Conn *websocket.Conn
}

func NewKubeLogger(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*KubeLogger, error) {
	conn, err := upGrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	kubelogger := &KubeLogger{Conn: conn}

	return kubelogger, nil
}

func (kl *KubeLogger) Write(data []byte) (int, error) {
	if err := kl.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (kl *KubeLogger) Close() error {
	return kl.Conn.Close()
}
