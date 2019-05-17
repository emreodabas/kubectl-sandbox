package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

var (
	usr, _        = user.Current()
	kubeHome      = usr.HomeDir + "/.kube/"
	path          = usr.HomeDir + "/.kube/k3s"
	kubeConfigCmd = "--kubeconfig /etc/rancher/k3s/k3s.yaml"
	rancherPath   = usr.HomeDir + "/.rancher/k3s/"
)

const (
	env                     = "DE2V"
	k3sLink                 = " https://github.com/rancher/k3s/releases/download/v0.5.0/k3s"
	loadDemoDataPromptValue = " You could load sample demo data to your Kubernetes Instance. Do you want to install/reset demo data [y/N] ?"
	downloadPromptValue     = " Kubectl Demo will download k3s (lighweight kubernetes) <<https://github.com/rancher/k3s/releases>> \n" +
		" and sample files from github repo <<https://github.com/emreodabas/samples>>. " +
		" Approximately 40mb file will download." +
		"Do you agree with this  [y/N] ? "
	resetK3sPromtValue = "kubectl demo will remove your k3s instance. Do you want to continue [y/N] ? "
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		initK3s()
	} else {
		switch argsWithoutProg[0] {

		case "reset":
			if Confirm(resetK3sPromtValue) {
				uninstallK3s()
			}
			initK3s()

		}
	}
}
func uninstallK3s() {

	stopServer()
	commandRun("rm " + path)
	commandRun("rm " + kubeHome + "k3s.* ")
	commandRun("rm -R " + rancherPath)
}

func initK3s() {

	if !isInstalled() {
		fmt.Println("Starting K3s installation")
		installK3s()
	}
	if startK3sServer() {
		fmt.Println("Server succesfully started")
	} else {
		fmt.Println("Server start error. You could check details in ~/.kube/k3s.log ")
		return
	}

	loadDemoData()

	createTerminal()

	defer stopServer()
}
func createTerminal() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\033[1;36m%s\033[0m", "kubectl-demo$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if cmdString != "" {
			err = runCommand(cmdString)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	if strings.Contains(commandStr, "kubectl") {
		commandStr = commandStr + " " + kubeConfigCmd
	}
	arrCommandStr := strings.Fields(commandStr)

	if len(arrCommandStr) != 0 {
		switch arrCommandStr[0] {
		case "exit", "quit", "q":
			os.Exit(0)
			// add another case here for custom commands.
		}
	}
	if len(arrCommandStr) > 1 {
		cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	} else if len(arrCommandStr) == 1 {
		cmd := exec.Command(arrCommandStr[0])
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
	return nil
}

func stopServer() {
	fmt.Println("Stopping K3s server")
	err := commandRun("systemctl stop k3s")
	//err := commandRun("pkill k3s")
	if err != nil {
		fmt.Printf(" Server start failed %v\n", err)
		return
	}
}

func startK3sServer() bool {
	fmt.Println("Starting K3s server")
	err := commandRun("sudo systemctl start k3s ")
	//err := commandRun("sudo ~/.kube/k3s server > ~/.kube/k3s.log 2>&1 &")
	//err = commandRun("nohup ~/.kube/k3s agent --server localhost  > ~/.kube/k3s.log 2>&1 &")
	if err != nil {
		fmt.Printf(" Server start failed %v\n", err)
		return false
	}
	return serverHealth()

}
func serverHealth() bool {

	for start := time.Now(); time.Since(start) < time.Second*10; {
		cmd := exec.Command("systemctl", "check", " k3s")
		bytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("System service check error %v\n", err)
		}

		if strings.Contains(string(bytes), "active") {
			return true
		} else {
			time.Sleep(time.Second)
		}
	}
	return false

	//if !fileExists(rancherPath) {
	//	fmt.Println("waiting for first configuration of Rancher")
	//	time.Sleep(time.Second * 5)
	//}
	//
	//for start := time.Now(); time.Since(start) < time.Second*10; {
	//	bytes, _ := commandRunAndReturn("kubectl get ns " + kubeConfigCmd + " | grep default")
	//	output := string(bytes)
	//	if strings.Contains(output, "default") {
	//		return true
	//	}
	//}
	//return false

}

// TODO Demo yamls need to define
func loadDemoData() {
	if Confirm(loadDemoDataPromptValue) {
		commandRun("kubectl apply " + kubeConfigCmd + " -f https://raw.githubusercontent.com/kubernetes/examples/master/guestbook/all-in-one/guestbook-all-in-one.yaml ")
	}
}

func installK3s() {
	if Confirm(downloadPromptValue) {
		var err error
		if env == "DEV" {
			err = commandRun("cp " + usr.HomeDir + "/Documents/k3s/k3s05 " + kubeHome + "/k3s")
		} else {
			err = commandRun("cd " + usr.HomeDir + " && curl -sfL https://get.k3s.io | sh -")
			//err = commandRun("wget" + k3sLink + "  && chmod +x k3s && mv k3s " + kubeHome)
		}
		if err != nil {
			fmt.Printf("Download k3s failed please try again %v\n", err)
		}
	} else {
		fmt.Println("Kubectl Demo is just useless without k3s. \n " +
			"=============================================== \n  " +
			"====> May the Kubernetes be with you :) <====== \n " +
			"=============================================== \n ")
		os.Exit(0)
	}
}

func isInstalled() bool {
	fmt.Println("IsInstalled")
	cmd := exec.Command("systemctl", "status", "k3s")
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("System service status error %v\n", err)
	}
	if strings.Contains(string(bytes), "could not be found") {
		return false
	} else if strings.Contains(string(bytes), "active") {
		return true
	} else {
		return false
	}
	//return fileExists(path) || fileExists(kubeHome+"k3s")
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
