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

package hcs

import (
	"fmt"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

type SDBInstanceDatabase struct {
	instance *SDBInstance
	multicloud.SDBInstanceDatabaseBase
	huawei.HuaweiTags

	Name         string
	CharacterSet string
}

func (database *SDBInstanceDatabase) GetId() string {
	return database.Name
}

func (database *SDBInstanceDatabase) GetGlobalId() string {
	return database.Name
}

func (database *SDBInstanceDatabase) GetName() string {
	return database.Name
}

func (database *SDBInstanceDatabase) GetStatus() string {
	return api.DBINSTANCE_DATABASE_RUNNING
}

func (database *SDBInstanceDatabase) GetCharacterSet() string {
	return database.CharacterSet
}

func (database *SDBInstanceDatabase) Delete() error {
	return database.instance.region.DeleteDBInstanceDatabase(database.instance.Id, database.Name)
}

func (region *SRegion) DeleteDBInstanceDatabase(instanceId, database string) error {
	resource := fmt.Sprintf("instances/%s/database/%s", instanceId, database)
	err := region.rdsDelete(resource)
	return err
}

// rds查询数据库列表
func (region *SRegion) GetDBInstanceDatabases(instanceId string) ([]SDBInstanceDatabase, error) {
	resource := fmt.Sprintf("%s/database/detail", instanceId)
	databases := []SDBInstanceDatabase{}
	err := region.rdsList(resource, nil, databases)
	if err != nil {
		return nil, err
	}
	return databases, nil
}
