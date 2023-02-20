package generator

import (
	"bytes"
	_ "embed"

	// containerd "github.com/containerd/containerd/services/server/config"
	// "github.com/irbgeo/go-structure"
	"github.com/pelletier/go-toml/v2"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/utils"
)

const (
	containerdConfigurationFilePath = "/etc/kubernetes/containerd/config.toml"
	containerdConfigurationFilePERM = 0644
)

var testCfg = &config.Config{
	ExtraArgs: map[any]any{
		"version": 2,
		"plugins": map[any]any{
			"io.containerd.grpc.v1.cri": map[any]any{
				"containerd": map[any]any{
					"runtimes": map[any]any{
						"runc": map[any]any{
							"runtime_type": "io.containerd.runc.v2",
							"options": map[any]any{
								"SystemdCgroup": true,
							},
						},
					},
				},
			},
		},
	},
}

// createContainerdConfiguration create containerd.service file.
func createContainerdConfiguration(i config.Instruction) error {
	data, err := createContainerdConfigurationData(i.Spec.Configuration)
	if err != nil {
		return err
	}

	return utils.CreateFile(containerdConfigurationFilePath, data, containerdConfigurationFilePERM, "root:root")
}

func createContainerdConfigurationData(cfg *config.Config) ([]byte, error) {
	eargs, err := getMap(cfg.ExtraArgs)
	if err != nil {
		return nil, err
	}

	// tomlData, err := toml.Marshal(eargs)
	// if err != nil {
	// 	return nil, err
	// }

	// cfgTest := new(containerd.Config)

	// cc, err := structure.New(new(containerd.Config))
	// if err != nil {
	// 	return nil, err
	// }

	// cc.ChangeTags(getContainerdTag)

	// err = toml.Unmarshal(tomlData, cfgTest)
	// if err != nil {
	// 	return nil, err
	// }

	buf := new(bytes.Buffer)
	e := toml.NewEncoder(buf)
	e.SetIndentTables(true)
	err = e.Encode(eargs)

	return buf.Bytes(), err
}

// var regexpContainerdTag = regexp.MustCompile(`"$`)

// func getContainerdTag(fieldName, fieldTag, _ string) string {
// 	if fieldTag != "" {
// 		fieldTag = regexpContainerdTag.ReplaceAllString(fieldTag, `,omitempty"`)
// 	}
// 	return fieldTag
// }
