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

package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	deliveryservice "github.com/koderover/zadig/pkg/microservice/aslan/core/delivery/service"
	internalhandler "github.com/koderover/zadig/pkg/microservice/aslan/internal/handler"
	e "github.com/koderover/zadig/pkg/tool/errors"
)

func ListDeliveryProduct(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	//params validate
	orgIDStr := c.Query("orgId")
	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc("orgId can't be empty!")
		return
	}
	ctx.Resp, ctx.Err = deliveryservice.FindDeliveryProduct(orgID, ctx.Logger)
}
