package k8s

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/lza730/k8s-test/tools/errorHelper"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Client *kubernetes.Clientset
}

const timeOut = 300

func CreateK8sClient(kubeConfig string, inCluster bool) *Client {
	glog.Info("Creating Kubernetes client")
	inClusterConf := ""
	if inCluster {
		inClusterConf = "in"
	} else {
		inClusterConf = "out of"
	}
	glog.Infof("Creating %s cluster config", inClusterConf)

	var config *rest.Config
	var err error
	if inCluster {
		config, err = rest.InClusterConfig()
		errorHelper.PanicOnError(err, fmt.Sprintf("InClusterConfig failed: %v", err))
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		errorHelper.PanicOnError(err, fmt.Sprintf("BuildConfigFromFlags kubeConfig=%s failed: %v", kubeConfig, err))
	}
	client, err := kubernetes.NewForConfig(config)
	errorHelper.PanicOnError(err, fmt.Sprintf("NewForConfig kubeConfig=%s failed: %v", kubeConfig, err))

	v, err := client.Discovery().ServerVersion()
	errorHelper.PanicOnError(err, "k8sClient discovery server version failed")
	glog.Infof("Running %s Kubernetes Cluster - version v%v.%v (%v) - platform %v",
		inClusterConf, v.Major, v.Minor, v.GitVersion, v.Platform)

	return &Client{
		Client: client,
	}
}

func (c *Client) GetPodList(namespace string) (*coreV1.PodList, error) {
	return c.Client.CoreV1().Pods(namespace).List(metaV1.ListOptions{})
}

func (c *Client) GetPodListBySelector(namespace string, selector string, podName string) (*coreV1.PodList, error) {
	return c.Client.CoreV1().Pods(namespace).List(metaV1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", selector, podName),
	})
}

func (c *Client) WaitForPodToRunning(namespace string, selector string, podName string, replicas int) bool {
	waitTime := time.Second
	for int(waitTime) < timeOut * replicas {
		time.Sleep(waitTime)
		waitTime *= 2

		pods, err := c.GetPodListBySelector(namespace, selector, podName)
		errorHelper.InfoOnError(err, fmt.Sprintf("Error getting list of loadbots: %s", err))
		if len(pods.Items) != replicas {
			glog.Infof("Replicas of pods is inaccurate. Waiting %v then checking again.", waitTime)
			continue
		}
		flag := true
		for _, pod := range pods.Items {
			if pod.Status.Phase != coreV1.PodRunning {
				flag = false
				break
			}
		}
		if flag == true {
			return true
		}
	}
	return false
}

func (c *Client)ReplicationControllers(namespace string, name string, replicas int32) {
	glog.Infof("Scaling %s to %d replicas", name, replicas)
	rc, err := c.Client.CoreV1().ReplicationControllers(namespace).Get(name, metaV1.GetOptions{})
	errorHelper.PanicOnError(err, fmt.Sprintf("Error scaling %s to %d replicas: %s", name, replicas, err))
	rc.Spec.Replicas = &replicas
	_, err = c.Client.CoreV1().ReplicationControllers(namespace).Update(rc)
	errorHelper.PanicOnError(err, fmt.Sprintf("Error scaling %s to %d replicas: %s", name, replicas, err))
}

func (c *Client)GetPodsBasic(pods *coreV1.PodList) []*PodBasic {
	var podsBasic []*PodBasic
	for _, pod := range pods.Items {
		podsBasic = append(podsBasic, &PodBasic{
			PodName: pod.Name,
			PodIp:   pod.Status.PodIP,
		})
	}
	return podsBasic
}

func (c *Client)RestGetAbsPathDoRaw(url string) ([]byte, error) {
	return c.Client.Discovery().RESTClient().Get().AbsPath(url).DoRaw()
}
