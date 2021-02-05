package main

import (
	"github.com/kintohub/utils-go/klog"
	utilsGrpc "github.com/kintohub/utils-go/server/grpc"
	"github.com/kintoproj/kinto-core/internal/build"
	"github.com/kintoproj/kinto-core/internal/build/api"
	"github.com/kintoproj/kinto-core/internal/config"
	"github.com/kintoproj/kinto-core/internal/controller"
	"github.com/kintoproj/kinto-core/internal/server"
	"github.com/kintoproj/kinto-core/internal/store"
	"github.com/kintoproj/kinto-core/internal/store/kube"
	pkgTypes "github.com/kintoproj/kinto-core/pkg/types"
)

// container method for all singletons in the project
type container struct {
	store       store.StoreInterface
	buildClient build.BuildInterface
	controller  controller.ControllerInterface

	kintoCoreService *server.KintoCoreService
}

func main() {

	klog.InitLogger()
	config.InitConfig()

	container := initContainer()

	utilsGrpc.RunServer(config.GrpcPort, config.GrpcWebPort, config.CORSAllowedHost,
		container.kintoCoreService.RegisterToServer,
	)
}

// Do not change the order of initialization due to dependencies needing to be instantiated!
func initContainer() *container {
	container := &container{}

	container.store = kube.NewKubeStore(config.KubeConfigPath)

	container.buildClient =
		api.NewBuildAPI(pkgTypes.NewWorkflowAPIServiceClient(
			utilsGrpc.CreateConnectionOrDie(config.BuildApiHost, false)))

	container.controller = controller.NewController(container.store, container.buildClient)
	container.kintoCoreService = server.NewKintoCoreService(container.controller)

	return container
}
