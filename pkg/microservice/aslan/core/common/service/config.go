/*
Copyright 2021 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/koderover/zadig/pkg/microservice/aslan/core/common/dao/models"
	commonrepo "github.com/koderover/zadig/pkg/microservice/aslan/core/common/dao/repo"
)

func GetConfigTemplateByService(serviceName string, log *zap.SugaredLogger) (*models.Config, error) {

	revisions, err := commonrepo.NewConfigColl().ListMaxRevisions("")
	if err != nil {
		errMsg := fmt.Sprintf("[ConfigTmpl.ListMaxRevisions] error: %v", err)
		log.Error(errMsg)
		return nil, err
	}

	for _, rev := range revisions {
		if rev.ServiceName == serviceName {
			return rev, nil
		}
	}
	return &models.Config{}, nil
}
