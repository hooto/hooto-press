// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package v1

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/rdb"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/datax"
	"github.com/hooto/hpress/store"
)

var (
	spaceReg                       = regexp.MustCompile(" +")
	term_list_limit          int64 = 15
	term_list_limit_taxonomy int64 = 200
)

type Term struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *Term) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c Term) ListAction() {

	var ls api.TermList

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "editor.list", config.Config.InstanceID) {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	model, err := config.SpecTermModel(c.Params.Get("modname"), c.Params.Get("modelid"))
	if err != nil {
		ls.Error = &types.ErrorMeta{
			Code:    "404",
			Message: "Spec or Model Not Found",
		}
		return
	}

	page, limit := c.Params.Int64("page"), term_list_limit

	dq := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	if model.Type == api.TermTaxonomy {
		limit = term_list_limit_taxonomy
		page = 1
	}

	if page < 1 {
		page = 1
	}

	dq.Limit(limit)
	if page > 1 {
		dq.Offset(int64((page - 1) * limit))
	}

	//
	ls = dq.TermList()

	dqc := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))

	if c.Params.Get("qry_text") != "" {
		dqc.Filter("title.like", "%"+c.Params.Get("qry_text")+"%")
	}

	count, err := dqc.TermCount()
	if err != nil {
		ls.Error = &types.ErrorMeta{api.ErrCodeInternalError, err.Error()}
		return
	}

	ls.Kind = "TermList"
	ls.Meta.TotalResults = uint64(count)
	ls.Meta.StartIndex = uint64((page - 1) * limit)
	ls.Meta.ItemsPerList = uint64(limit)
}

func (c Term) EntryAction() {

	rsp := api.Term{
		TypeMeta: types.TypeMeta{
			APIVersion: api.Version,
		},
	}

	defer c.RenderJson(&rsp)

	if !iamclient.SessionAccessAllowed(c.Session, "editor.read", config.Config.InstanceID) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dq := datax.NewQuery(c.Params.Get("modname"), c.Params.Get("modelid"))
	dq.Limit(100)

	dq.Filter("id", c.Params.Get("id"))

	rsp = dq.TermEntry()
}

func (c Term) SetAction() {

	rsp := api.Term{}

	defer c.RenderJson(&rsp)

	if !iamclient.SessionAccessAllowed(c.Session, "editor.write", config.Config.InstanceID) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	model, err := config.SpecTermModel(c.Params.Get("modname"), c.Params.Get("modelid"))
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    "404",
			Message: "Spec or Model Not Found",
		}
		return
	}

	if err := c.Request.JsonDecode(&rsp); err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    "400",
			Message: "Bad Request " + err.Error(),
		}
		return
	}

	var (
		set   = map[string]interface{}{}
		table = fmt.Sprintf("tx%s_%s", idhash.HashToHexString([]byte(c.Params.Get("modname")), 12), c.Params.Get("modelid"))
	)

	q := rdb.NewQuerySet().From(table).Limit(1)

	switch model.Type {

	case api.TermTag:

		uniTitle := spaceReg.ReplaceAllString(strings.TrimSpace(strings.ToLower(rsp.Title)), " ")

		h := md5.New()
		io.WriteString(h, uniTitle)
		rsp.UID = fmt.Sprintf("%x", h.Sum(nil))[:16]
		rsp.ID = 0

		q.Where.And("uid", rsp.UID)

		rs, err := store.Data.Query(q)
		if err != nil {
			rsp.Error = &types.ErrorMeta{
				Code:    "500",
				Message: "Can not pull database instance",
			}
			return
		}

		if len(rs) == 1 {

			rsp.ID = rs[0].Field("id").Uint32()

			if rs[0].Field("title").String() != rsp.Title {
				set["title"] = rsp.Title
			}

			if rs[0].Field("status").Int16() != rsp.Status {
				set["status"] = rsp.Status
			}

		} else {

			set["uid"] = rsp.UID
			set["title"] = rsp.Title
			set["status"] = rsp.Status
			set["created"] = rdb.TimeNow("datetime")
			set["userid"] = c.us.UserId()
		}

	case api.TermTaxonomy:

		if rsp.ID > 0 {

			q.Where.And("id", rsp.ID)

			rs, err := store.Data.Query(q)
			if err != nil {
				rsp.Error = &types.ErrorMeta{
					Code:    "500",
					Message: "Can not pull database instance",
				}
				return
			}

			if len(rs) < 1 {
				rsp.Error = &types.ErrorMeta{
					Code:    "404",
					Message: "Term Not Found",
				}
				return
			}

			if rs[0].Field("title").String() != rsp.Title {
				set["title"] = rsp.Title
			}

			if rs[0].Field("status").Int16() != rsp.Status {
				set["status"] = rsp.Status
			}

			if rs[0].Field("pid").Uint32() != rsp.PID {
				set["pid"] = rsp.PID
			}

			if rs[0].Field("weight").Int32() != rsp.Weight {
				set["weight"] = rsp.Weight
			}

		} else {

			set["pid"] = rsp.PID
			set["title"] = rsp.Title
			set["status"] = rsp.Status
			set["weight"] = rsp.Weight
			set["created"] = rdb.TimeNow("datetime")
			set["userid"] = c.us.UserId()
		}

		datax.TermTaxonomyCacheClean(c.Params.Get("modname"), c.Params.Get("modelid"))

	default:
		rsp.Error = &types.ErrorMeta{
			Code:    "500",
			Message: "Server Error",
		}
		return
	}

	if len(set) > 0 {

		set["updated"] = rdb.TimeNow("datetime")

		if rsp.ID > 0 {

			ft := rdb.NewFilter()
			ft.And("id", rsp.ID)
			_, err = store.Data.Update(table, set, ft)

		} else {

			fmt.Println("SET", table, "___", set)
			rs, err := store.Data.Insert(table, set)
			if err == nil {
				if incrid, err := rs.LastInsertId(); err == nil && incrid > 0 {
					rsp.ID = uint32(incrid)
				} else {
					err = errors.New("Can Not Get LastInsertId")
				}
			}
		}

		if err != nil {
			rsp.Error = &types.ErrorMeta{
				Code:    "500",
				Message: err.Error(),
			}
			return
		}
	}

	rsp.Model = model

	rsp.Kind = "Term"
}
