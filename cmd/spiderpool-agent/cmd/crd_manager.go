// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	spiderpoolv1 "github.com/spidernet-io/spiderpool/pkg/k8s/apis/spiderpool.spidernet.io/v1"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(spiderpoolv1.AddToScheme(scheme))
}

func newCRDManager() (ctrl.Manager, error) {
	config := ctrl.GetConfigOrDie()
	config.Burst = 100
	config.QPS = 50

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     "0",
		HealthProbeBindAddress: "0",
	})
	if err != nil {
		return nil, err
	}

	if err := mgr.GetFieldIndexer().IndexField(agentContext.InnerCtx, &spiderpoolv1.SpiderSubnet{}, "spec.default", func(raw client.Object) []string {
		subnet := raw.(*spiderpoolv1.SpiderSubnet)
		return []string{strconv.FormatBool(*subnet.Spec.Default)}
	}); err != nil {
		return nil, err
	}

	if err := mgr.GetFieldIndexer().IndexField(agentContext.InnerCtx, &spiderpoolv1.SpiderIPPool{}, "spec.default", func(raw client.Object) []string {
		ipPool := raw.(*spiderpoolv1.SpiderIPPool)
		return []string{strconv.FormatBool(*ipPool.Spec.Default)}
	}); err != nil {
		return nil, err
	}

	if err := mgr.GetFieldIndexer().IndexField(agentContext.InnerCtx, &spiderpoolv1.SpiderReservedIP{}, "spec.ipVersion", func(raw client.Object) []string {
		reservedIP := raw.(*spiderpoolv1.SpiderReservedIP)
		return []string{strconv.FormatInt(*reservedIP.Spec.IPVersion, 10)}
	}); err != nil {
		return nil, err
	}

	return mgr, nil
}
