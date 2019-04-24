// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type installCredentials struct {
	ClusterName               string `json:"clusterName"`
	ClusterZone               string `json:"clusterZone"`
	JenkinsUser               string `json:"jenkinsUser"`
	JenkinsPassword           string `json:"jenkinsPassword"`
	GithubPersonalAccessToken string `json:"githubPersonalAccessToken"`
	GithubUserEmail           string `json:"githubUserEmail"`
	GithubOrg                 string `json:"githubOrg"`
}

/*
type gCloudClusterDescription struct {
	AddonsConfig struct {
		HTTPLoadBalancing struct {
		} `yaml:"httpLoadBalancing,omitempty"`
		KubernetesDashboard struct {
			Disabled bool `yaml:"disabled"`
		} `yaml:"kubernetesDashboard,omitempty"`
		NetworkPolicyConfig struct {
			Disabled bool `yaml:"disabled,omitempty"`
		} `yaml:"networkPolicyConfig,omitempty"`
	} `yaml:"addonsConfig,omitempty"`
	ClusterIpv4Cidr          string    `yaml:"clusterIpv4Cidr"`
	CreateTime               time.Time `yaml:"createTime,omitempty"`
	CurrentMasterVersion     string    `yaml:"currentMasterVersion"`
	CurrentNodeCount         int       `yaml:"currentNodeCount"`
	CurrentNodeVersion       string    `yaml:"currentNodeVersion"`
	DefaultMaxPodsConstraint struct {
		MaxPodsPerNode string `yaml:"maxPodsPerNode"`
	} `yaml:"defaultMaxPodsConstraint,omitempty"`
	Endpoint              string   `yaml:"endpoint,omitempty"`
	InitialClusterVersion string   `yaml:"initialClusterVersion,omitempty"`
	InstanceGroupUrls     []string `yaml:"instanceGroupUrls,omitempty"`
	IPAllocationPolicy    struct {
	} `yaml:"ipAllocationPolicy,omitempty"`
	LabelFingerprint string `yaml:"labelFingerprint,omitempty"`
	LegacyAbac       struct {
	} `yaml:"legacyAbac,omitempty"`
	Location       string   `yaml:"location,omitempty"`
	Locations      []string `yaml:"locations,omitempty"`
	LoggingService string   `yaml:"loggingService,omitempty"`
	MasterAuth     struct {
		ClientCertificate    string `yaml:"clientCertificate,omitempty"`
		ClientKey            string `yaml:"clientKey,omitempty"`
		ClusterCaCertificate string `yaml:"clusterCaCertificate,omitempty"`
		Password             string `yaml:"password,omitempty"`
		Username             string `yaml:"username,omitempty"`
	} `yaml:"masterAuth,omitempty"`
	MasterAuthorizedNetworksConfig struct {
	} `yaml:"masterAuthorizedNetworksConfig,omitempty"`
	MonitoringService string `yaml:"monitoringService,omitempty"`
	Name              string `yaml:"name,omitempty"`
	Network           string `yaml:"network,omitempty"`
	NetworkConfig     struct {
		Network    string `yaml:"network,omitempty"`
		Subnetwork string `yaml:"subnetwork,omitempty"`
	} `yaml:"networkConfig,omitempty"`
	NetworkPolicy struct {
	} `yaml:"networkPolicy,omitempty"`
	NodeConfig struct {
		DiskSizeGb     int      `yaml:"diskSizeGb,omitempty"`
		DiskType       string   `yaml:"diskType,omitempty"`
		ImageType      string   `yaml:"imageType,omitempty"`
		MachineType    string   `yaml:"machineType,omitempty"`
		OauthScopes    []string `yaml:"oauthScopes,omitempty"`
		ServiceAccount string   `yaml:"serviceAccount,omitempty"`
	} `yaml:"nodeConfig,omitempty"`
	NodeIpv4CidrSize int `yaml:"nodeIpv4CidrSize,omitempty"`
	NodePools        []struct {
		Autoscaling struct {
		} `yaml:"autoscaling,omitempty"`
		Config struct {
			DiskSizeGb     int      `yaml:"diskSizeGb,omitempty"`
			DiskType       string   `yaml:"diskType,omitempty"`
			ImageType      string   `yaml:"imageType,omitempty"`
			MachineType    string   `yaml:"machineType,omitempty"`
			OauthScopes    []string `yaml:"oauthScopes,omitempty"`
			ServiceAccount string   `yaml:"serviceAccount,omitempty"`
		} `yaml:"config,omitempty"`
		InitialNodeCount  int      `yaml:"initialNodeCount,omitempty"`
		InstanceGroupUrls []string `yaml:"instanceGroupUrls,omitempty"`
		Management        struct {
		} `yaml:"management,omitempty"`
		Name     string `yaml:"name,omitempty"`
		SelfLink string `yaml:"selfLink,omitempty"`
		Status   string `yaml:"status,omitempty"`
		Version  string `yaml:"version,omitempty"`
	} `yaml:"nodePools,omitempty"`
	ResourceLabels struct {
		Owner string `yaml:"owner,omitempty"`
	} `yaml:"resourceLabels,omitempty"`
	SelfLink         string `yaml:"selfLink,omitempty"`
	ServicesIpv4Cidr string `yaml:"servicesIpv4Cidr,omitempty"`
	Status           string `yaml:"status,omitempty"`
	Subnetwork       string `yaml:"subnetwork,omitempty"`
	Zone             string `yaml:"zone,omitempty"`
}
*/
// installCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs keptn on your Kubernetes cluster",
	Long: `Installs keptn on your Kubernetes cluster

Example:
	keptn install`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing keptn...")
		// var creds *installCredentials

		clusterName, clusterZone, _ := connectToCluster()
		fmt.Printf("Cluster: %s", clusterName)

		var err error
		var creds *installCredentials
		creds, err = getInstallCredentials()

		// execCmd = exec.Command("kubectl", "logs", "control-lxxt6-deployment-569499f6cb-gbzgp", "-n", "keptn", "-c", "user-container", "-f")

		// get the YAML for the installer pod
		fileURL := "https://raw.githubusercontent.com/keptn/keptn/install-keptn/install/manifests/installer/installer.yaml"

		if err := DownloadFile("installer.yaml", fileURL); err != nil {
			panic(err)
		}

		gcloudUser := getGcloudUser()
		clusterIPCIDR, servicesIPCIDR := getGcloudClusterIPCIDR(clusterName, clusterZone)
		setDeploymentFileKey("JENKINS_USER", "admin")
		setDeploymentFileKey("JENKINS_PASSWORD", "AiTx4u8VyUV8tCKk")
		setDeploymentFileKey("GITHUB_PERSONAL_ACCESS_TOKEN", creds.GithubPersonalAccessToken)
		setDeploymentFileKey("GITHUB_USER_EMAIL", creds.GithubUserEmail)
		setDeploymentFileKey("GITHUB_ORGANIZATION", creds.GithubOrg)
		setDeploymentFileKey("GCLOUD_USER", gcloudUser)
		setDeploymentFileKey("CLUSTER_IPV4_CIDR", clusterIPCIDR)
		setDeploymentFileKey("SERVICES_IPV4_CIDR", servicesIPCIDR)

		execCmd := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/keptn/keptn/install-keptn/install/manifests/installer/rbac.yaml")

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := execCmd.StdoutPipe()
		stderrIn, _ := execCmd.StderrPipe()
		err = execCmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			copyAndCapture(os.Stdout, stdoutIn)
			wg.Done()
		}()

		copyAndCapture(os.Stderr, stderrIn)

		wg.Wait()

		err = execCmd.Wait()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}
		outStr, errStr := string(stdout), string(stderr)
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func getInstallCredentials() (*installCredentials, error) {
	var err error
	var credsStr string
	var creds *installCredentials
	credsStr, err = credentialmanager.GetInstallCreds()
	if err != nil {
		credsStr = ""
	}
	err = json.Unmarshal([]byte(credsStr), &creds)
	if err != nil {
		creds = new(installCredentials)
	}

	fmt.Print("Please enter the following information (press enter to keep the old value):\n")

	reader := bufio.NewReader(os.Stdin)

	// default creds for jenkins
	creds.JenkinsUser = "admin"
	creds.JenkinsPassword = "AiTx4u8VyUV8tCKk"

	fmt.Printf("GitHub User Email [%s]: ", creds.GithubUserEmail)
	githubUserEmail, _ := reader.ReadString('\n')
	if githubUserEmail != "\n" {
		creds.GithubUserEmail = strings.TrimSuffix(githubUserEmail, "\n")
	}

	fmt.Printf("GitHub Organization [%s]: ", creds.GithubOrg)
	githubOrg, _ := reader.ReadString('\n')
	if githubOrg != "\n" {
		creds.GithubOrg = strings.TrimSuffix(githubOrg, "\n")
	}

	fmt.Printf("GitHub Personal Access Token [%s]: ", creds.GithubPersonalAccessToken)
	githubPersonalAccessToken, _ := reader.ReadString('\n')
	if githubPersonalAccessToken != "\n" {
		creds.GithubPersonalAccessToken = strings.TrimSuffix(githubPersonalAccessToken, "\n")
	}

	newCreds, _ := json.Marshal(creds)
	newCredsStr := strings.Replace(string(newCreds), "\n", "", -1)
	fmt.Printf("new creds file content: %s\n", newCredsStr)
	credentialmanager.SetInstallCreds(newCredsStr)
	return creds, err
}

func copyAndCapture(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// text := strings.Replace(scanner.Text(), "found", "", -1)
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func setDeploymentFileKey(key string, value string) {
	input, err := ioutil.ReadFile("installer.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte("value: "+key), []byte("value: "+value), -1)

	if err = ioutil.WriteFile("installer.yaml", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func connectToCluster() (string, string, string) {
	reader := bufio.NewReader(os.Stdin)
	setupNewCluster := false
	clusterName, clusterZone, gkeProject := getClusterInfo()
	if clusterName != "" && clusterZone != "" && gkeProject != "" {
		fmt.Printf("You are currently connected to the following cluster:\n")
		fmt.Printf("Project: %s\n", gkeProject)
		fmt.Printf("Compute Zone: %s\n", clusterZone)
		fmt.Printf("Cluster Name: %s\n", clusterName)
		fmt.Printf("Would you like to use this cluster to set up keptn? [Y/n]\n")
		setupNewClusterPrompt, _ := reader.ReadString('\n')

		if setupNewClusterPrompt != "\n" {
			setupNewCluster = true
		}
	} else {
		setupNewCluster = true
	}

	if setupNewCluster {
		connectionSuccessful := false
		for tryConnection := true; tryConnection; tryConnection = !connectionSuccessful {
			for ok := true; ok; ok = (gkeProject == "") {
				fmt.Printf("Please enter the GKE project [%s]:", gkeProject)
				newGkeProject, _ := reader.ReadString('\n')
				if newGkeProject != "\n" {
					gkeProject = strings.TrimSuffix(newGkeProject, "\n")
				}
			}

			for ok := true; ok; ok = (clusterName == "") {
				fmt.Printf("Please enter the GKE cluster name [%s]:", clusterName)
				newClusterName, _ := reader.ReadString('\n')
				if newClusterName != "\n" {
					clusterName = strings.TrimSuffix(newClusterName, "\n")
				}
			}

			for ok := true; ok; ok = (clusterZone == "") {
				fmt.Printf("Please enter the GKE cluster zone [%s]:", clusterZone)
				newClusterZone, _ := reader.ReadString('\n')
				if newClusterZone != "\n" {
					clusterZone = strings.TrimSuffix(newClusterZone, "\n")
				}
			}

			cmd := exec.Command(
				"gcloud",
				"container",
				"clusters",
				"get-credentials",
				clusterName,
				"--zone",
				clusterZone,
				"--project",
				gkeProject,
			)

			_, err := cmd.Output()
			if err != nil {
				fmt.Println("Could not connect to cluster. Please verify that you have entered the correct information:")
			} else {
				connectionSuccessful = true
				fmt.Println("Connection to cluster successful:")
			}
		}
	}

	return clusterName, clusterZone, gkeProject
}

func getClusterInfo() (string, string, string) {
	var clusterName string
	var clusterZone string
	var project string
	fmt.Println("Trying to get GKE cluster info...")
	// try to get current cluster from gcloud config
	cmd := exec.Command("gcloud", "config", "get-value", "container/cluster")
	out, err := cmd.Output()
	if err != nil {
		clusterName = ""
	}
	clusterName = strings.Replace(string(out), "\n", "", -1)

	cmd = exec.Command("gcloud", "config", "get-value", "compute/zone")
	out, err = cmd.Output()
	if err != nil {
		clusterZone = ""
	}
	clusterZone = strings.Replace(string(out), "\n", "", -1)

	cmd = exec.Command("gcloud", "config", "get-value", "core/project")
	out, err = cmd.Output()
	if err != nil {
		project = ""
	}
	project = strings.Replace(string(out), "\n", "", -1)

	return clusterName, clusterZone, project
}

func getGcloudUser() string {
	var err error

	cmd := exec.Command("gcloud", "config", "get-value", "account")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return strings.Replace(string(out), "\n", "", -1)
}

func getGcloudClusterIPCIDR(clusterName string, clusterZone string) (string, string) {
	var err error
	var clusterDescription map[string]interface{}
	cmd := exec.Command("gcloud", "container", "clusters", "describe", clusterName, "--zone="+clusterZone)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Could not get cluster Info\n")
	}

	err = yaml.Unmarshal([]byte(out), &clusterDescription)
	clusterIPCIDR := clusterDescription["clusterIpv4Cidr"].(string)
	servicesIPCIDR := clusterDescription["servicesIpv4Cidr"].(string)

	return clusterIPCIDR, servicesIPCIDR
}
