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

package middleware

import (
	"github.com/gin-gonic/gin"

	commonservice "github.com/koderover/zadig/pkg/microservice/aslan/core/common/service"
	"github.com/koderover/zadig/pkg/util/ginzap"
)

// 更新操作日志状态
func UpdateOperationLogStatus(c *gin.Context) {
	log := ginzap.WithContext(c).Sugar()
	c.Next()
	if c.GetString("operationLogID") == "" {
		return
	}
	err := commonservice.UpdateOperation(c.GetString("operationLogID"), c.Writer.Status(), log)
	if err != nil {
		log.Errorf("UpdateOperation err :%v", err)
	}
}
