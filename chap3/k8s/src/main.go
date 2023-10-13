package main

import (
	"context"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/labstack/echo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type User struct {
	Name  string
	Email string
}

type AppConfig struct {
	Environment string `yaml:"environment"`
}

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	e := echo.New()
	e.GET("/", healthCheck)
	e.GET("/configmap", showConfigMap)
	e.GET("/secret", showSecret)
	e.GET("/user", showUser)

	e.Logger.Fatal(e.Start(":8080"))
}

func showConfigMap(c echo.Context) error {
	namespace := "default"
	configMapName := "myconfig"

	clientset, err := createKubernetesClient()
	if err != nil {
		fmt.Printf("Error creating Kubernetes client:%s\n", err)
		return fmt.Errorf("Error creating Kubernetes client:%s\n", err)
	}

	configMapData, err := getConfigMap(clientset, namespace, configMapName)
	if err != nil {
		fmt.Printf("Error getting ConfigMap data:%s\n", err)
		return fmt.Errorf("Error getting ConfigMap data:%s\n", err)
	}

	// cofigMapを利用する
	appData := configMapData["app"]
	var appConfig AppConfig
	err = yaml.Unmarshal([]byte(appData), &appConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML:%s\n", err)
		return fmt.Errorf("Error parsing YAML:%s\n", err)
	}
	return c.JSON(http.StatusOK, appConfig.Environment)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func showUser(c echo.Context) error {
	name := c.QueryParam("name")
	email := c.QueryParam("email")

	u := new(User)
	u.Name = name
	u.Email = email

	return c.JSON(http.StatusOK, u)
}

func createKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getConfigMap(clientset *kubernetes.Clientset, namespace, name string) (map[string]string, error) {
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return configMap.Data, nil
}

func showSecret(c echo.Context) error {
	namespace := "default"         // Secretが存在するネームスペース
	secretName := "my-credentials" // 読み取るSecretの名前

	clientset, err := createKubernetesClient()
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		return fmt.Errorf("Error creating Kubernetes client:%s", err)
	}

	secretData, err := getSecret(clientset, namespace, secretName)
	if err != nil {
		fmt.Println("Error getting Secret data:", err)
		return fmt.Errorf("Error getting Secret data:%s", err)
	}

	configJSON := secretData["config.json"]

	var config Config
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return fmt.Errorf("Error parsing JSON:%s", err)
	}
	// keyJSON を利用して必要な処理を行う
	return c.JSON(http.StatusOK, config)
}

func getSecret(clientset *kubernetes.Clientset, namespace, name string) (map[string][]byte, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret.Data, nil
}
