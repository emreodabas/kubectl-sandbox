package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var k3sLink = "https://github.com/rancher/k3s/releases/download/v0.2.0-rc5/k3s"
var kubeHome = "/home/emreo/.kube/"
var kubeConfigCmd = "--kubeconfig " + kubeHome + "k3s.yaml"

var demoList = []string{
	"Basic Commands (Beginner) ",
	"Basic Commands (Intermediate) ",
	"Deploy Commands ",
	"Cluster Management Commands ",
	"Troubleshooting and Debugging Commands ",
	"Advanced Commands ",
	"Settings Commands ",
	"Other Commands ",
}

func main() {
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
		fmt.Println("Server start error. You could check details in ~/.kube/k3s.logs ")
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
	showDemo(item)
	//commandRun("kubectl get ns " + kubeConfigCmd)
}
func showDemo(i int) {
}

func loadLessonData(i int) {

}

func continueDemo() int {
	return 1
}
func isDemoStarted() bool {
	return false
}

func showDemoList() int {
	return Choose("", demoList)
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

	return fileExists(kubeHome + "k3s")
}

// Confirm continues prompting until the input is boolean-ish.
func Confirm(prompt string, args ...interface{}) bool {
	for {
		switch String(prompt, args...) {
		case "Yes", "yes", "y", "Y":
			return true
		case "No", "no", "n", "N":
			return false
		}
	}
}

// String prompt.
func String(prompt string, args ...interface{}) string {
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

func Choose(prompt string, list []string) int {
	fmt.Println()
	for i, val := range list {
		fmt.Printf("  %d) %s\n", i+1, val)
	}

	fmt.Println()
	i := -1

	for {
		s := String(prompt)

		// index
		n, err := strconv.Atoi(s)
		if err == nil {
			if n > 0 && n <= len(list) {
				i = n - 1
				break
			} else {
				continue
			}
		}

		// value
		i = indexOf(s, list)
		if i != -1 {
			break
		}
	}

	return i
}

func indexOf(s string, list []string) int {
	for i, val := range list {
		if val == s {
			return i
		}
	}
	return -1
}
