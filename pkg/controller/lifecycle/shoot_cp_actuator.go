// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package lifecycle

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	"github.com/Kristian-ZH/gardener-extension-logging/pkg/imagevector"
	gardenerkubernetes "github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/utils/chart"
	gardeneriv "github.com/gardener/gardener/pkg/utils/imagevector"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	shootChart = &chart.Chart{
		Name: "hello-world",
		Path: filepath.Join("charts/shoot-bootstrap", "hello-world"),
		Objects: []*chart.Object{
			{Type: &appsv1.Deployment{}, Name: "hello-world"},
		},
	}
)

type shootActuator struct {
	logger       logr.Logger // logger
	chart        chart.Interface
	chartApplier gardenerkubernetes.ChartApplier
	imageVector  gardeneriv.ImageVector

	client            client.Client
	clientset         kubernetes.Interface
	gardenerClientset gardenerkubernetes.Interface
}

func (a *shootActuator) InjectConfig(config *rest.Config) error {
	var err error

	a.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("could not create Kubernetes client: %w", err)
	}

	a.gardenerClientset, err = gardenerkubernetes.NewWithConfig(gardenerkubernetes.WithRESTConfig(config))
	if err != nil {
		return fmt.Errorf("could not create Gardener client: %w", err)
	}

	a.chartApplier = a.gardenerClientset.ChartApplier()

	return nil
}

func (a *shootActuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

// NewActuator returns an actuator responsible for Extension resources.
func NewShootActuator() Actuator {
	return &shootActuator{
		logger:      log.Log.WithName("logging actuator"),
		chart:       shootChart,
		imageVector: imagevector.ImageVector(),
	}
}

// Reconcile the Extension resource.
func (a *shootActuator) Reconcile(ctx context.Context, _ logr.Logger, ex *extensionsv1alpha1.Logging) error {
	a.logger.Info("Hello World, I just entered the Reconcile method")
	fmt.Println("SHOOT Hello World, I just entered in the ACTUATOR")
	if err := a.chart.Apply(ctx, a.chartApplier, ex.Namespace, a.imageVector, "", "", map[string]interface{}{}); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

// Delete the Extension resource.
func (a *shootActuator) Delete(ctx context.Context, _ logr.Logger, ex *extensionsv1alpha1.Logging) error {
	a.logger.Info("SHOOT Hello World, I just entered the Delete method")
	if err := a.chart.Delete(ctx, a.client, ex.Namespace); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

// Restore the Extension resource.
func (a *shootActuator) Restore(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Logging) error {
	return a.Reconcile(ctx, log, ex)
}

// Migrate the Extension resource.
func (a *shootActuator) Migrate(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Logging) error {
	return a.Delete(ctx, log, ex)
}
