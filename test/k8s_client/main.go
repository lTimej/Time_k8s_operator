package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	config *rest.Config
	err    error
)

func init() {
	config, err = clientcmd.BuildConfigFromFlags("", "./conf/config")
	if err != nil {
		panic(err)
	}

}
func get_pods() {
	// 1. 构造访问config的配置，从文件中加载，将 home目录下的 .kube/config拷贝到当前./conf/下
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.APIPath = "/api"
	// 2. 创建rest client
	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	// 3. 查找命名空间dev下的pod
	var podList v1.PodList
	err = client.Get().Namespace("dev").Resource("pods").Do(context.Background()).Into(&podList)
	if err != nil {
		log.Printf("get pods error:%v\n", err)
		return
	}

	fmt.Println("dev pod count:", len(podList.Items))
	for _, pod := range podList.Items {
		fmt.Printf("name: %s\n", pod.Name)
	}
}
func create_pods() {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	mypod := v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mynginx",
			Namespace: "dev",
			Labels: map[string]string{
				"run": "nginx",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image: "nginx:1.14-alpine",
					Name:  "mynginx",
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}
	_, err = clientset.CoreV1().Pods("dev").Create(context.Background(), &mypod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("pod create failure,err:", err)
		return
	}
	fmt.Println("pod create success")
}
func main() {
	get_pods()
	// create_pods()
}
