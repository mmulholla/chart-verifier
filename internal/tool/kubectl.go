package tool

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/utils"
)

//go:embed kubeOpenShiftVersionMap.yaml
var content embed.FS

// Based on https://access.redhat.com/solutions/4870701
var kubeOpenShiftVersionMap map[string]string

type versionMap struct {
	Versions []*versionMapping `yaml:"versions"`
}

type versionMapping struct {
	KubeVersion string `yaml:"kube-version"`
	OcpVersion  string `yaml:"ocp-version"`
}

func init() {
	kubeOpenShiftVersionMap = make(map[string]string)

	yamlFile, err := content.ReadFile("kubeOpenShiftVersionMap.yaml")

	if err != nil {
		utils.LogError(fmt.Sprintf("Error reading content of kubeOpenShiftVersionMap.yaml: %v", err))
		return
	}

	versions := versionMap{}
	err = yaml.Unmarshal(yamlFile, &versions)
	if err != nil {
		utils.LogError(fmt.Sprintf("Error reading content of kubeOpenShiftVersionMap.yaml: %v", err))
		return
	}

	for _, versionMap := range versions.Versions {
		kubeOpenShiftVersionMap[versionMap.KubeVersion] = versionMap.OcpVersion
	}

}

type Kubectl struct {
	clientset kubernetes.Interface
}

func NewKubectl(kubeConfig clientcmd.ClientConfig) (*Kubectl, error) {
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	config.APIPath = "/api"
	config.GroupVersion = &schema.GroupVersion{Group: "core", Version: "v1"}
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	kubectl := new(Kubectl)
	kubectl.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubectl, nil
}

func (k Kubectl) WaitForDeployments(context context.Context, namespace string, selector string) error {
	deadline, _ := context.Deadline()
	unavailableDeployments := []string{"None"}
	getDeploymentsError := ""

	utils.LogInfo(fmt.Sprintf("Start wait for deployments. --timeout time left: %s ", deadline.Sub(time.Now()).String()))

	for deadline.After(time.Now()) && len(unavailableDeployments) > 0 {
		unavailableDeployments = []string{}
		deployments, err := k.clientset.AppsV1().Deployments(namespace).List(context, metav1.ListOptions{
			LabelSelector: selector,
		})
		if err != nil {
			unavailableDeployments = []string{"None"}
			getDeploymentsError = fmt.Sprintf("error getting deployments from namespace %s : %v", namespace, err)
			utils.LogWarning(getDeploymentsError)
		} else {
			getDeploymentsError = ""
			for _, deployment := range deployments.Items {
				// Just after rollout, pods from the previous deployment revision may still be in a
				// terminating state.
				unavailableDeployments = append(unavailableDeployments, deployment.Name)
			}
			if len(unavailableDeployments) > 0 {
				utils.LogInfo(fmt.Sprintf("Wait for unavailable deployments: %v", strings.Join(unavailableDeployments, ",")))
				time.Sleep(time.Second)
			} else {
				utils.LogInfo(fmt.Sprintf("Finish wait for deployments, --timeout time left %s", deadline.Sub(time.Now()).String()))
			}
		}
	}

	if len(getDeploymentsError) > 0 {
		errorMsg := fmt.Sprintf("Time out retrying after %s", getDeploymentsError)
		utils.LogError(errorMsg)
		return errors.New(errorMsg)
	}
	if len(unavailableDeployments) > 0 {
		errorMsg := fmt.Sprintf("Time out waiting for unavailable deployments: %v", strings.Join(unavailableDeployments, ","))
		utils.LogError(errorMsg)
		return errors.New(errorMsg)
	}

	return nil
}

func (k Kubectl) DeleteNamespace(context context.Context, namespace string) error {
	if err := k.clientset.CoreV1().Namespaces().Delete(context, namespace, *metav1.NewDeleteOptions(0)); err != nil {
		return err
	}
	return nil
}

func (k Kubectl) GetServerVersion() (*version.Info, error) {
	version, err := k.clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, err
	}
	return version, err
}

func GetKubeOpenShiftVersionMap() map[string]string {
	return kubeOpenShiftVersionMap
}

func GetClientConfig(envSettings *cli.EnvSettings) clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(envSettings.KubeConfig) > 0 {
		loadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: envSettings.KubeConfig}
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{})
}
