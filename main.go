package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"runtime"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	nameSpace string
	podName   string
)

func main() {
	flag.Parse()
	pods := NewPod()
	fmt.Println("Starting------")
	// 查看指定空间 pod
	pods.GetPod(nameSpace)
	//		 查看pod事件
	pods.ListEvents(nameSpace, podName)
}

type Pod struct {
	ClientSet kubernetes.Interface
	Config    *rest.Config
}

func NewPod() *Pod {
	var err error
	var k8sConfig *rest.Config

	if runtime.GOOS == "windows" {
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", "config")
	} else {
		k8sConfig, err = rest.InClusterConfig()
	}
	//k8sConfig, err := clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		fmt.Println("k8s集群config配置文件加载失败:", err)
	}
	// 初始化客户端
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		fmt.Println("客户端初始化失败:", err)
	}

	return &Pod{
		ClientSet: k8sClient,
		Config:    k8sConfig,
	}
}

func (k *Pod) GetPod(namespace string) {
	//podObj, err := k.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	podObj, err := k.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error getting pods: ", err)
		return
	}

	for _, pod := range podObj.Items {
		fmt.Printf("当前命名空间:%s,pod名称:%s,查询时间:%v\n", namespace, pod.Name, time.Now().Format(time.DateTime))
	}
}

func (k *Pod) ListEvents(namespace, podName string) {
	podEvents, _ := k.ClientSet.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{})
	for _, v := range podEvents.Items {
		if v.Type == "Normal" {
			continue
		}

		if strings.HasPrefix(v.Name, podName) {
			fmt.Println(v.Name)
			fmt.Println(v.Message)
		}
	}
}

func init() {
	flag.StringVar(&nameSpace, "n", "default", "指定命名空间")
	flag.StringVar(&podName, "p", "", "指定pod名称")
}
