package main

import (
	"flag"
	"strings"

	"go.uber.org/zap"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/controller"
)

var (
	Version = "undefined"
)

func main() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	var configFile string
	var skipKindList string
	flag.StringVar(&configFile, "config", "", "path to dir with configs")
	flag.StringVar(&skipKindList, "slip-kinds", "", `list of skipped kind
	supported kind:
		KubeletService
		KubeletConfiguration
		ContainerdService
		ContainerdConfiguration
		SysctlNetworkConfiguration
		ModProbConfiguration`)
	flag.Parse()

	if configFile == "" {
		zap.L().Fatal("the path to the config file is not set")
	}

	skippingKind := make(map[string]struct{})
	kindList := strings.Split(skipKindList, " ")
	for _, k := range kindList {
		skippingKind[k] = struct{}{}
	}

	zap.L().Debug("configuration", zap.String("version", Version))

	files, err := config.ReadConfig(configFile)
	if err != nil {
		zap.L().Fatal("read config", zap.Error(err))
	}

	zap.L().Info("started")

	err = controller.Generation(files, skippingKind)
	if err != nil {
		zap.L().Fatal("generation", zap.Error(err))
	}

	zap.L().Info("goodbye")
}
