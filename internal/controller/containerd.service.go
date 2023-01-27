package controller

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
)

var (
	//go:embed template/containerd.service.tmpl
	containerdServiceTemplateStr string
	containerdServiceTemplate    = template.Must(template.New("containerd-service").Parse(containerdServiceTemplateStr))
)

const (
	containerdServiceName     = "containerd"
	containerdServiceFilePath = "/etc/systemd/system/containerd.service"
	containerdServiceFilePERM = 0644
)

// createContainerdService create containerd.service file.
func createContainerdService(cfg config.Instruction) error {
	data, err := createContainerdServiceData(cfg.Spec)
	if err != nil {
		return err
	}

	err = createFile(containerdServiceFilePath, data, containerdServiceFilePERM, "root:root")
	if err != nil {
		return err
	}

	err = startService(containerdServiceName)
	return err
}

func createContainerdServiceData(spec any) ([]byte, error) {
	eargs, err := getMap(spec)
	if err != nil {
		return nil, err
	}

	containerdServiceBuffer := new(bytes.Buffer)
	err = containerdServiceTemplate.Execute(containerdServiceBuffer, eargs)
	return containerdServiceBuffer.Bytes(), err
}
