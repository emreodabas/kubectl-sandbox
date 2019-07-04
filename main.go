package main

import (
	"fmt"
	"github.com/peterh/liner"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	kubeConfigCmd = "--kubeconfig /etc/rancher/k3s/k3s.yaml"
	history       = filepath.Join(os.TempDir(), ".liner_example_history")
	keywords      = []string{"kubectl", "create", "update", "delete", "deployment"}
)

const (
	downloadPromptValue = " Kubectl Sandbox will download k3s (lighweight kubernetes detail in https://k3s.io/ ) \n" +
		"and create systemctl service approximately 40mb file will download." +
		"Do you agree with this  [y/N] ? "
	loadDemoDataPromptValue = " You could load sample demo data to your Kubernetes Instance. Do you want to install/reset demo data [y/N] ?"
	deletek3sPromptValue    = "kubectl sandbox will delete your k3s instance. Do you want to continue [y/N] ? "
	resetk3sPromptValue     = "kubectl sandbox will delete and install your k3s instance. Do you want to continue [y/N] ? "
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		initK3s(false)
	} else {
		switch argsWithoutProg[0] {

		case "uninstall", "remove", "delete":
			if Confirm(deletek3sPromptValue) {
				err := uninstallK3s()
				if err != nil {
					fmt.Println("Error while uninstall K3s")
				}
			}

		case "reset":
			if Confirm(resetk3sPromptValue) {
				err := uninstallK3s()
				if err != nil {
					fmt.Println("Error while uninstall K3s")
				}
				installK3s()
			}
			initK3s(false)
		case "load":
			initK3s(true)
		default:
			initK3s(false)
		}
	}
}
func uninstallK3s() error {

	stopServer()
	return commandSudoRun("/usr/local/bin/k3s-uninstall.sh")

}

func initK3s(loadData bool) {

	if !isInstalled() {
		fmt.Println("Starting K3s installation")
		installK3s()
		loadData = true
	}
	if startK3sServer() {
		fmt.Println("Server succesfully started")
	} else {
		fmt.Println("Server start error. You could check details in k3s service log ")
		return
	}

	if loadData {
		loadDemoData()
	}

	createTerminal()
	defer stopServer()
}
func createTerminal() {
	fmt.Println("You could exit terminal with exit|quit commands or Ctrl+C|Ctrl+D ")
	line := liner.NewLiner()
	defer func() {
		e := line.Close()
		if e != nil {
			fmt.Println("Error while history file close with defer")
		}
	}()

	line.SetCtrlCAborts(true)

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range keywords {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(history); err == nil {
		_, err := line.ReadHistory(f)
		if err != nil {
			fmt.Println("Error while history read")
		}
		er := f.Close()
		if er != nil {
			fmt.Println("Error while history file close")
		}
	}
	for {
		if cmdString, err := line.Prompt("kubectl-sandbox$ "); err == nil {
			line.AppendHistory(cmdString)
			if strings.Contains(cmdString, "exit") || strings.Contains(cmdString, "quit") {
				return
			}
			signal.Ignore(syscall.SIGINT)
			err = runCommand(cmdString)
			if err != nil {
				fmt.Println(err)
			}
		} else if err == liner.ErrPromptAborted {
			return
		} else {
			//CTRL+D
			return
		}
		if f, err := os.Create(history); err != nil {
			log.Print("Error writing history file: ", err)
		} else {
			_, err := line.WriteHistory(f)
			if err != nil {
				fmt.Println("Error while history write")
			}
			er := f.Close()
			if er != nil {
				fmt.Println("Error while history file close")
			}
		}
	}
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	if strings.Contains(commandStr, "kubectl sandbox") {
		commandStr = " echo \"kubectl sandbox in kubectl sandbox could be deadlock. YOU SHALL NOT PASS :) \" "
	} else if strings.Contains(commandStr, "kubectl") {
		commandStr = "sudo " + commandStr + " " + kubeConfigCmd
	}
	arrCommandStr := strings.Fields(commandStr)

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
	err := commandSudoRun("systemctl stop k3s")
	if err != nil {
		fmt.Printf(" Server start failed %v\n", err)
		return
	}
}

func startK3sServer() bool {
	fmt.Println("Starting K3s server")
	err := commandSudoRun("systemctl start k3s ")
	if err != nil {
		fmt.Printf(" Server start failed %v\n", err)
		return false
	}
	return serverHealth()

}
func serverHealth() bool {

	for start := time.Now(); time.Since(start) < time.Second*10; {
		cmd := exec.Command("systemctl", "check", " k3s")
		bytes, _ := cmd.CombinedOutput()

		if strings.Contains(string(bytes), "active") {
			return true
		} else {
			time.Sleep(time.Second)
		}
	}
	return false

}

func installK3s() {

	if Confirm(downloadPromptValue) {
		var command = ""
		if isKubectlAvailable() {
			command = "cd " + os.TempDir() + " && curl -sfL https://get.k3s.io | INSTALL_K3S_BIN_DIR_READ_ONLY=\"false\" sh -s - "
		} else {
			command = "cd " + os.TempDir() + " && curl -sfL https://get.k3s.io | sh -s - "
		}
		fmt.Println(command)
		err := commandRun(command)
		if err != nil {
			fmt.Printf("Download k3s failed please try again %v\n", err)
		}

	} else {
		fmt.Println("Kubectl Sandbox is just useless without k3s. \n " +
			"=============================================== \n  " +
			"====> May the Kubernetes be with you :) <====== \n " +
			"=============================================== \n ")
		os.Exit(0)
	}
}

func isInstalled() bool {
	cmd := exec.Command("systemctl", "status", "k3s")
	bytes, _ := cmd.CombinedOutput()
	if strings.Contains(string(bytes), "could not be found") {
		fmt.Println("installation could not be found")
		return false
	} else if strings.Contains(string(bytes), "active") || strings.Contains(string(bytes), "k3s.io") {
		fmt.Println("installation and service active")
		return true
	} else {
		fmt.Println("installation error" + string(bytes))
		return false
	}
}

func loadDemoData() {
	if Confirm(loadDemoDataPromptValue) {
		err := commandSudoRun("kubectl apply " + kubeConfigCmd + " -f https://raw.githubusercontent.com/kubernetes/examples/master/guestbook/all-in-one/guestbook-all-in-one.yaml ")
		if err != nil {
			fmt.Println("Loading sample data error")
		}
	}
}

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
	_, err := fmt.Scanln(&s)
	if err != nil {
		fmt.Println("Prompt value could not read")
	}
	return s
}

func isKubectlAvailable() bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v kubectl")
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func commandRun(command string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func commandSudoRun(command string) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo "+command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
