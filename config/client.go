package config

import (
	"context"
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type kubeconfigclient struct {
	namespace string
	cmi       v1.ConfigMapInterface
	Mutex     *sync.Mutex
}

type Interface interface {
	ReadConfig(name string) map[string]string
	WatchConfig(name string)
}

func GetClient(namespace string) Interface {
	config, err := rest.InClusterConfig()

	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	cmi := clientset.CoreV1().ConfigMaps(namespace) //ConfigMapsInterface
	return &kubeconfigclient{
		namespace: namespace,
		cmi:       cmi,
		Mutex:     &sync.Mutex{},
	}
}

func (client *kubeconfigclient) ReadConfig(name string) map[string]string {

	configmap, err := client.cmi.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	return configmap.Data
}

func (client *kubeconfigclient) WatchConfig(name string) {
	watchif, err := client.cmi.Watch(context.TODO(), metav1.SingleObject(metav1.ObjectMeta{
		Namespace: client.namespace,
		Name: name,
	}))
	if err != nil {
		panic(err.Error())
	}

	eventchan := watchif.ResultChan()

	for {
		event, open := <-eventchan
		if open {
			switch event.Type {
			case watch.Added:
				fallthrough
			case watch.Modified:
				client.Mutex.Lock()
				// handle the update
				if updatedMap, ok := event.Object.(*corev1.ConfigMap); ok {
					fmt.Println("Received updates on ConfigMap")
					fmt.Println(updatedMap.Data)
				client.Mutex.Unlock()
				}
			
			default:
				// do nothing here
			}
		} else {
			// chan closed -> stop watching
			return			
		}

	}

}
