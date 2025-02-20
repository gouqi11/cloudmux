// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shell

import (
	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

func init() {
	type KubeClusterListOptions struct {
		IDs    []string `help:"Kube Cluster ids"`
		Offset int      `help:"List offset"`
		Limit  int      `help:"List limit"`
	}
	shellutils.R(&KubeClusterListOptions{}, "kube-cluster-list", "List kube cluster", func(cli *qcloud.SRegion, args *KubeClusterListOptions) error {
		clusters, total, err := cli.GetKubeClusters(args.IDs, args.Offset, args.Limit)
		if err != nil {
			return err
		}
		printList(clusters, total, args.Offset, args.Limit, []string{})
		return nil
	})

	type KubeClusterIdOptions struct {
		ID string
	}

	shellutils.R(&KubeClusterIdOptions{}, "kube-cluster-delete", "Delete kube cluster", func(cli *qcloud.SRegion, args *KubeClusterIdOptions) error {
		return cli.DeleteKubeCluster(args.ID, false)
	})

	type KubeClusterKubeconfigOptions struct {
		ID      string
		Private bool
	}

	shellutils.R(&KubeClusterKubeconfigOptions{}, "kube-cluster-kubeconfig", "Get kube cluster kubeconfig", func(cli *qcloud.SRegion, args *KubeClusterKubeconfigOptions) error {
		config, err := cli.GetKubeConfig(args.ID, args.Private)
		if err != nil {
			return err
		}
		printObject(config)
		return nil
	})

}
