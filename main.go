package main

import (
	"./drawing"
	"./lessons"
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

var k3sLink = " https://github.com/rancher/k3s/releases/download/v0.2.0-rc5/k3s"
var usr, _ = user.Current()
var kubeHome = usr.HomeDir + "/.kube/"
var path = kubeHome + "k3s"
var kubeConfigCmd = "--kubeconfig " + kubeHome + "k3s.yaml"

var selections = []selection{
	{"intro.json", "1-> Introduction to kubectl", "kubectl controls the Kubernetes cluster manager.   \n" +
		"  \n" +
		"Find more information at: https://kubernetes.io/docs/reference/kubectl/overview/  \n" +
		"  \n" +
		"Usage:  \n" +
		"  kubectl [flags] [options]  \n" +
		"Use \"kubectl <command> --help\" for more information about a given command.  \n" +
		"Use \"kubectl options\" for a list of global command-line options (applies to all commands). "},
	{"lessons/basic-commands-beginner.json", "2-> Basic Commands (Beginner)", "\n \t create         Create a resource from a file or from stdin.  \n " +
		"\t expose         Take a replication controller, service, deployment or pod and  expose it as a new Kubernetes Service  \n " +
		"\t run            Run a particular image on the cluster  \n " +
		"\t set            Set specific features on objects  \n " +
		"   \n "},
	{"basic-commands-intermediate.json", "3-> Basic Commands (Intermediate) ", "\n \t explain        Documentation of resources  \n " +
		"\t get            Display one or many resources  \n " +
		"\t edit           Edit a resource on the server  \n " +
		"\t delete         Delete resources by filenames, stdin, resources and names, or by resources and label selector  \n " +
		"   \n "},
	{"deploy-commands.json", "4-> Deploy Commands ", "\n \t rollout        Manage the rollout of a resource  \n " +
		"\t scale          Set a new size for a Deployment, ReplicaSet, Replication  Controller, or Job  \n " +
		"\t autoscale      Auto-scale a Deployment, ReplicaSet, or ReplicationController  \n " +
		"   \n "},
	{"cluster-managements-commands.json", "5-> Cluster Management Commands ", "\n \t certificate    Modify certificate resources.  \n " +
		"\t cluster-info   Display cluster info  \n " +
		"\t top            Display Resource (CPU/Memory/Storage) usage.  \n " +
		"\t cordon         Mark node as unschedulable  \n " +
		"\t uncordon       Mark node as schedulable  \n " +
		"\t drain          Drain node in preparation for maintenance  \n " +
		"\t taint          Update the taints on one or more nodes  \n " +
		"   \n "},
	{"troubleshooting-debugging-commands.json", "6-> Troubleshooting and Debugging Commands ", "\n \t describe       Show details of a specific resource or group of resources  \n " +
		"\t logs           Print the logs for a container in a pod  \n " +
		"\t attach         Attach to a running container  \n " +
		"\t exec           Execute a command in a container  \n " +
		"\t port-forward   Forward one or more local ports to a pod  \n " +
		"\t proxy          Run a proxy to the Kubernetes API server  \n " +
		"\t cp             Copy files and directories to and from containers.  \n " +
		"\t auth           Inspect authorization  \n " +
		"   \n "},
	{"advanced-commands.json", "7-> Advanced Commands ", "\n \t diff           Diff live version against would-be applied version  \n " +
		"\t apply          Apply a configuration to a resource by filename or stdin  \n " +
		"\t patch          Update field(s) of a resource using strategic merge patch  \n " +
		"\t replace        Replace a resource by filename or stdin  \n " +
		"\t wait           Experimental: Wait for a specific condition on one or many resources.  \n " +
		"\t convert        Convert config files between different API versions  \n " +
		"   \n "},
	{"settings-commands.json", "8-> Settings Commands ", "\n \t label          Update the labels on a resource  \n " +
		"\t annotate       Update the annotations on a resource  \n " +
		"\t completion     Output shell completion code for the specified shell (bash or zsh)  \n " +
		"   \n "},
	{"other-commands.json", "9-> Other Commands ", "\n \t api-resources  Print the supported API resources on the server  \n " +
		"\t api-versions   Print the supported API versions on the server, in the form of \"group/version\"  \n " +
		"\t config         Modify kubeconfig files  \n " +
		"\t plugin         Provides utilities for interacting with plugins.  \n " +
		"\t version        Print the client and server version information  \n "},
}

type selection struct {
	filePath    string
	name        string
	description string
}

func main() {
	fmt.Println(kubeHome)
	initK3s()
}

func initK3s() {

	if !isInstalled() {
		fmt.Println("Starting K3s installation")
		installK3s()
		//loadDemoData()
	}
	if startK3sServer() {
		fmt.Println("Server succesfully started")
	} else {
		fmt.Println("Server start error. You could check details in ~/.kube/k3s.log ")
		return
	}

	startDemo()
	defer stopServer()
}

func stopServer() {
	fmt.Println("Stopping K3s server")
	err := commandRun("pkill k3s")
	if err != nil {
		fmt.Printf(" Server start failed %v\n", err)
		return
	}
}

func startDemo() {
	var item int
	if isDemoStarted() {
		item = continueDemo()
		Confirm("Do you want to continue ? [Y/N] ")
	} else {
		item = showDemoList()
	}
	showDemo(selections[item])
	//commandRun("kubectl get ns " + kubeConfigCmd)
}
func showDemo(s selection) {

	lesson, _ := lessons.Init(s.filePath)
	fmt.Println("Showing Descriptions")
	drawing.ShowLesson(lesson, lessons.Desc, 1)
}

//TODO
func continueDemo() int {
	return 1
}

//TODO
func isDemoStarted() bool {
	return false
}

func showDemoList() int {
	idx, err := fuzzyfinder.FindMulti(
		selections,
		func(i int) string {
			return selections[i].name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s \nDescription: %s",
				strings.SplitAfter(selections[i].name, "> ")[1],
				selections[i].description)
		}))
	if err != nil {
		log.Fatal(err)
	}
	return idx[0]
}

func startK3sServer() bool {
	fmt.Println("Starting K3s server")
	err := commandRun("nohup ~/.kube/k3s server --disable-agent  > ~/.kube/k3s.log 2>&1 &")
	if err != nil {
		fmt.Printf(" Server start failed %v\n", err)
		return false
	}
	return serverHealth()

}
func serverHealth() bool {
	for start := time.Now(); time.Since(start) < time.Second; {

		bytes, _ := commandRunAndReturn("kubectl get ns " + kubeConfigCmd + " | grep default")
		output := string(bytes)
		if strings.Contains(output, "default") {
			return true
		}
	}
	return false

}

// TODO Demo yamls need to define
func loadDemoData() {

}

func installK3s() {
	if Confirm(" Kubectl Demo will download k3s (lighweight kubernetes) <<https://github.com/rancher/k3s/releases>> \n" +
		" and sample files from github repo <<https://github.com/emreodabas/samples>>. " +
		" Approximately 40mb file will download." +
		"Do you agree with this  [y/N] ?? }") {
		err := commandRun("wget" + k3sLink + "  && chmod +x k3s && mv k3s " + kubeHome)
		if err != nil {
			fmt.Printf("Download k3s failed please try again %v\n", err)
			return
		}
	}
}

func isInstalled() bool {
	fmt.Println(path)
	return fileExists(path)
}

// Confirm continues prompting until the input is boolean-ish.
func Confirm(promptValue string, args ...interface{}) bool {
	for {
		switch Prompt(promptValue, args...) {
		case "Yes", "yes", "y", "Y":
			return true
		case "No", "no", "n", "N":
			return false
		}
	}
}

func Prompt(prompt string, args ...interface{}) string {
	var s string
	fmt.Printf(prompt+": ", args...)
	fmt.Scanln(&s)
	return s
}

func commandRun(command string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func commandRunAndReturn(command string) ([]byte, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	return cmd.CombinedOutput()
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//func Choose(promptValue string, list map[string]string) string {
//	fmt.Println()
//	for i, val := range list {
//		fmt.Printf("  %d) %s\n", i+1, val)
//	}
//
//	fmt.Println()
//	i := ""
//
//	for {
//		s := Prompt(promptValue)
//
//		n, err := strconv.Atoi(s)
//		if err == nil {
//			if n > 0 && n <= len(list) {
//				i = n - 1
//				break
//			} else {
//				continue
//			}
//		}
//
//		// value
//		i = indexOf(s, list)
//		if i != "" {
//			break
//		}
//	}
//
//	return i
//}

func indexOf(s string, list map[string]string) string {
	for i, val := range list {
		if val == s {
			return i
		}
	}
	return ""
}
