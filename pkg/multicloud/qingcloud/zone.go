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

package qingcloud

import (
	"strings"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type SZone struct {
	region *SRegion

	host     *SHost
	Status   string
	ZoneId   string
	RegionId string
}

func (zone *SRegion) GetZones() ([]SZone, error) {
	params := map[string]string{}
	resp, err := zone.ec2Request("DescribeZones", params)
	if err != nil {
		return nil, err
	}
	ret := []SZone{}
	err = resp.Unmarshal(&ret, "zone_set")
	if err != nil {
		return nil, err
	}
	result := []SZone{}
	for i := range ret {
		if strings.HasPrefix(ret[i].ZoneId, zone.Region) {
			result = append(result, ret[i])
		}
	}
	return result, nil
}

func (zone *SZone) GetIHosts() ([]cloudprovider.ICloudHost, error) {
	return []cloudprovider.ICloudHost{zone.getHost()}, nil
}

func (zone *SZone) getHost() *SHost {
	if zone.host == nil {
		zone.host = &SHost{zone: zone}
	}
	return zone.host
}
