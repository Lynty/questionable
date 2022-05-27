// assumes 1. already authenticated with Google Cloud SDK, kubectl binary installed, using ssh protocol for config sync repository

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Cluster struct {
	Zone string
	Name string
}

func getClusters(project string) []Cluster {
	fmt.Println("Getting clusters...")
	out, err := exec.Command("gcloud", "container", "clusters", "list", "--project", project, "--format", "json(name,zone)").Output()
	if err != nil {
		log.Fatalf("Failed to retrieve clusters: %v", err)
	}
	clusterJson := string(out)

	var clusters []Cluster
	var clusterMap []map[string]interface{}
	err2 := json.Unmarshal([]byte(clusterJson), &clusterMap)
	if err2 != nil {
		panic(err2)
	}

	for _, clusterData := range clusterMap {
		// convert map to array of clusters struct
		var c Cluster
		c.Name = fmt.Sprintf("%s", clusterData["name"])
		c.Zone = fmt.Sprintf("%s", clusterData["zone"])
		clusters = append(clusters, c)
	}
	return clusters
}

func enableACM(project string) {
	fmt.Println("Enabling Config Management Feature...")
	out, err := exec.Command("gcloud", "beta", "container", "hub", "config-management", "enable", "--project", project).Output()
	if err != nil {
		log.Fatalf("Failed to enable config management: %v", err)
	}
	fmt.Println(string(out))
}

func setupKubeconfig(cluster, zone, project string) {

	fmt.Println("Setting up Kubeconfig...")
	fmt.Println("Cluster: ", cluster)
	fmt.Println("Zone: ", zone)
	fmt.Println("Project: ", project)
	out, err := exec.Command("gcloud", "container", "clusters", "get-credentials", cluster, "--zone", zone, "--project", project).Output()
	if err != nil {
		log.Fatalf("Failed to get credentials: %v", err)
	}
	fmt.Println(string(out))
}

func createNamespace(namespace string) {
	fmt.Println("Creating", namespace, "namespace...")
	out, err := exec.Command("kubectl", "create", "namespace", namespace).Output()
	if err != nil {
		println("Namespace already exists")
	}
	fmt.Println(string(out))
}
func createGitCreds(namespace, key string) {
	fmt.Println("Creating git-creds secret")
	out, err := exec.Command("kubectl", "create", "secret", "generic", "git-creds", "--namespace", namespace, "--from-file=ssh="+key).Output()
	if err != nil {
		println("Secret already exists")
	}
	fmt.Println(string(out))
}

func installACM(cluster, config, project string) {
	fmt.Println("Setting up Config Sync...")
	out, err := exec.Command("gcloud", "beta", "container", "hub", "config-management", "apply", "--membership", cluster, "--config", config, "--project", project).Output()
	if err != nil {
		log.Fatalf("Config sync install failed: %v", err)
	}
	fmt.Println(string(out))
}

// for CICD cluster output that need ACM installed, run all the things
func main() {
	namespace := "config-management-system"
	key := os.Args[1]
	config := os.Args[2]
	project := os.Args[3]
	//var project string
	//flag.StringVar(&project, "project", "", "gcp project id")
	//flag.Parse()
	//println("PROJECT: ", project)

	clusters := getClusters(project)
	for _, c := range clusters {
		println("Installing ACM on: ", c.Name)
		enableACM(project)
		setupKubeconfig(c.Name, c.Zone, project)
		createNamespace(namespace)
		createGitCreds(namespace, key)
		installACM(c.Name, config, project)
	}
}
