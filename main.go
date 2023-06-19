package main

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	pods := NewPod()
	fmt.Println("Starting------")
	for {
		pods.GetPod("dev")
		time.Sleep(time.Second * 2)
	}

}

type Pod struct {
	ClientSet kubernetes.Interface
	Config    *rest.Config
}

func NewPod() *Pod {
	k8sConfig, err := rest.InClusterConfig()
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
