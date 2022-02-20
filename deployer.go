package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"strings"
)

func DeployNewContainer(yamlString string, container string, service string) (string, error) {
	newYamlString := strings.Replace(yamlString, "{{CONTAINER_IMAGE}}", container, 1)
	port, err := evokeKubectl(newYamlString, service)
	if err != nil {
		return "", err
	}
	return port, nil
}

func evokeKubectl(yamlString string, service string) (string, error) {
	//here save temporary file and run kubectl
	fileName := "/root/" + uuid.NewString() + ".yaml"
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	_, err = file.WriteString(yamlString)
	if err != nil {
		return "", err
	}
	err = file.Close()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("kctl apply -f", fileName)
	stdout, err := cmd.Output()

	if err != nil {
		return "", err
	}
	fmt.Println(string(stdout))
	returnMessage := string(stdout)
	if !strings.Contains(returnMessage, "created") {
		return "", errors.New("couldn't create a new pod")
	}

	cmd = exec.Command(fmt.Sprintf("kctl expose pod interceptor --type=NodePort --name=%s", service))
	stdout, err = cmd.Output()
	if err != nil {
		return "", err
	}

	cmd = exec.Command(fmt.Sprintf("kubectl get svc %s --all-namespaces -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}}{{\"\\n\"}}{{end}}{{end}}{{end}}'", service))
	stdout, err = cmd.Output()
	if err != nil {
		return "", nil
	}

	return string(stdout), nil
}
