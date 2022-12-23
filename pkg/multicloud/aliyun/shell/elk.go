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
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type ElkListOptions struct {
		Page int
		Size int
	}
	shellutils.R(&ElkListOptions{}, "elastic-search-list", "List elastic searchs", func(cli *aliyun.SRegion, args *ElkListOptions) error {
		elks, _, err := cli.GetElasticSearchs(args.Size, args.Page)
		if err != nil {
			return err
		}
		printList(elks, 0, 0, 0, nil)
		return nil
	})

	type ElkIdOptions struct {
		ID string
	}

	shellutils.R(&ElkIdOptions{}, "elastic-search-show", "Show elasitc search", func(cli *aliyun.SRegion, args *ElkIdOptions) error {
		elk, err := cli.GetElasitcSearch(args.ID)
		if err != nil {
			return err
		}
		printObject(elk)
		return nil
	})

	shellutils.R(&ElkIdOptions{}, "elastic-search-delete", "Delete elasitc search", func(cli *aliyun.SRegion, args *ElkIdOptions) error {
		return cli.DeleteElasticSearch(args.ID)
	})

}
