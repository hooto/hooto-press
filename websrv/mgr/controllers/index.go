// Copyright 2015~2017 hooto Author, All rights reserved.
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

package controllers

import (
	"github.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/httpsrv"

	"code.hooto.com/hooto/hootopress/config"
	"code.hooto.com/hooto/hootopress/status"
)

type Index struct {
	*httpsrv.Controller
}

func (c Index) IndexAction() {

	status.Locker.RLock()
	defer status.Locker.RUnlock()

	if c.Params.Get("_iam_out") != "" {
		c.Redirect(c.UrlBase(""))
		return
	}

	if !iamclient.SessionIsLogin(c.Session) {
		c.Redirect(iamclient.AuthServiceUrl(
			config.Config.InstanceID,
			c.UrlBase("auth/cb"),
			c.Request.RawAbsUrl(),
		))
		return
	}

	if status.IamServiceStatus == status.IamServiceUnRegistered {
		c.Redirect("mgr/setup/index")
		return
	}

	c.Render("index.tpl")
}
