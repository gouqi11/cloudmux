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
	"yunion.io/x/cloudmux/pkg/multicloud/apsara"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

func init() {
	type VSwitchListOptions struct {
		Limit  int `help:"page size"`
		Offset int `help:"page offset"`
	}
	shellutils.R(&VSwitchListOptions{}, "vswitch-list", "List vswitches", func(cli *apsara.SRegion, args *VSwitchListOptions) error {
		vswitches, total, e := cli.GetVSwitches(nil, "", args.Offset, args.Limit)
		if e != nil {
			return e
		}
		printList(vswitches, total, args.Offset, args.Limit, []string{})
		return nil
	})

	type VSwitchShowOptions struct {
		ID string `help:"show vswitch details"`
	}
	shellutils.R(&VSwitchShowOptions{}, "vswitch-show", "Show vswitch details", func(cli *apsara.SRegion, args *VSwitchShowOptions) error {
		vswitch, e := cli.GetVSwitchAttributes(args.ID)
		if e != nil {
			return e
		}
		printObject(vswitch)
		return nil
	})

	shellutils.R(&VSwitchShowOptions{}, "vswitch-delete", "Show vswitch details", func(cli *apsara.SRegion, args *VSwitchShowOptions) error {
		e := cli.DeleteVSwitch(args.ID)
		if e != nil {
			return e
		}
		return nil
	})
}
