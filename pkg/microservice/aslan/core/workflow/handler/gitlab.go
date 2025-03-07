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
	"github.com/gin-gonic/gin"

	"github.com/koderover/zadig/pkg/microservice/aslan/core/workflow/service/webhook"
	internalhandler "github.com/koderover/zadig/pkg/microservice/aslan/internal/handler"
	e "github.com/koderover/zadig/pkg/tool/errors"
)

// @Router /workflow/webhook/gitlabhook [POST]
// @Summary Process webhook for giblab
// @Accept  json
// @Produce json
// @Success 200 {object} map[string]string "map[string]string - {message: success}"
func ProcessGitlabHook(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	log := ctx.Logger
	err := webhook.ProcessGitlabHook(c.Request, ctx.RequestID, log)
	if err != nil {
		c.JSON(e.ErrorMessage(err))
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
