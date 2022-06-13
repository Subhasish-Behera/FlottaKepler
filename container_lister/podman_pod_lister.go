package container_lister

import (
	"encoding/json"
	"fmt"
	libpod "github.com/containers/libpod"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"log"
	"net/http"
	"os"
)

type PodmanContainerLister struct{}

const (
	//	saPath         = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	EdgeDeviceEnv  = "NODE_NAME"
	kubeletPortEnv = "KUBELET_PORT"
)

var (
	containerUrl, metricsUrl string

	EdgeDeviceCpuUsageMetricName = "node_cpu_usage_seconds_total"
	EdgeDeviceMemUsageMetricName = "node_memory_working_set_bytes"
	containerCpuUsageMetricName  = "container_cpu_usage_seconds_total"
	containerMemUsageMetricName  = "container_memory_working_set_bytes"
	containerStartTimeMetricName = "container_start_time_seconds"

	containerNameTag = "container"
	//	podNameTag       = "pod"
	//	namespaceTag     = "namespace"
)

func init() {
	nodeName := os.Getenv(EdgeDeviceEnv)
	if len(nodeName) == 0 {
		nodeName = "localhost"
	}
	port := os.Getenv(kubeletPortEnv)
	if len(port) == 0 {
		port = "10250"
	}
	containerUrl = "https://" + nodeName + ":" + port + "/pods"
	metricsUrl = "https://" + nodeName + ":" + port + "/metrics/resource"
}

func httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from %q: %v", url, err)
	}
	return resp, err
}
func (k *PodmanContainerLister) ListContainers() (*[]libpod.Container, error) {
	resp, err := httpGet(containerUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	podList := corev1.PodList{}
	err = json.Unmarshal(body, &podList)
	if err != nil {
		log.Fatalf("failed to parse response body: %v", err)
	}

	pods := &podList.Items

	return pods, nil
}

func main() {

}
