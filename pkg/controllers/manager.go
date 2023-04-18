package controllers

import (
	"context"
	"os"
	"time"

	"github.com/spf13/pflag"

	"github.com/morvencao/xcm-connector/pkg/controllers/cluster"
	"github.com/morvencao/xcm-connector/pkg/helpers"
	"github.com/openshift/library-go/pkg/controller/controllercmd"

	clusterclient "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clusterinformers "open-cluster-management.io/api/client/cluster/informers/externalversions"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ManagerOptions defines the flags for kcp-ocm integration controller manager
type ManagerOptions struct {
	ControlPlaneKubeConfigFile string
	XCMServer                  string
	AccessTokenFile            string
	RefreshTokenFile           string
}

// NewManagerOptions returns the flags with default value set
func NewManagerOptions() *ManagerOptions {
	return &ManagerOptions{}
}

// AddFlags register and binds the default flags
func (o *ManagerOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(
		&o.ControlPlaneKubeConfigFile,
		"control-plane-kubeconfig",
		o.ControlPlaneKubeConfigFile,
		"Location of control plane kubeconfig file to connect to control plane cluster.",
	)

	flags.StringVar(
		&o.XCMServer,
		"xcm-server",
		o.XCMServer,
		"The host url of the xCM server, including scheme.",
	)

	// flags.StringVar(
	// 	&o.AccessTokenFile,
	// 	"xcm-pull-secret-token-file",
	// 	o.AccessTokenFile,
	// 	"The access token file for xCM server.",
	// )

	// flags.StringVar(
	// 	&o.AccessTokenFile,
	// 	"xcm-access-token-file",
	// 	o.AccessTokenFile,
	// 	"The access token file for xCM server.",
	// )

	// flags.StringVar(
	// 	&o.RefreshTokenFile,
	// 	"xcm-refresh-token-file",
	// 	o.RefreshTokenFile,
	// 	"The refresh token file.",
	// )
}

// Run starts all of controllers for kcp-ocm integration
func (o *ManagerOptions) Run(ctx context.Context, controllerContext *controllercmd.ControllerContext) error {
	wait.Poll(5*time.Second, 60*time.Second, func() (bool, error) {
		if _, err := os.Stat(o.ControlPlaneKubeConfigFile); err != nil {
			return false, nil
		}

		return true, nil
	})

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	controlPlaneKubeConfig, err := clientcmd.BuildConfigFromFlags("", o.ControlPlaneKubeConfigFile)
	if err != nil {
		return err
	}

	clusterClient, err := clusterclient.NewForConfig(controlPlaneKubeConfig)
	if err != nil {
		return err
	}

	// retrieve access token
	accessToken, err := helpers.RetrieveAccessToken(ctx, kubeClient)

	// kubeInfomer := informers.NewSharedInformerFactory(kubeClient, 10*time.Minute)
	clusterInformers := clusterinformers.NewSharedInformerFactory(clusterClient, 10*time.Minute)

	clusterController := cluster.NewClusterController(
		clusterClient,
		clusterInformers.Cluster().V1().ManagedClusters(),
		o.XCMServer,
		accessToken,
		controllerContext.EventRecorder,
	)

	go clusterInformers.Start(ctx.Done())

	go clusterController.Run(ctx, 1)

	<-ctx.Done()
	return nil
}