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

package datax

import (
	"fmt"
	"strings"

	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/store"
)

func (q *QuerySet) NodeCount() (int64, error) {

	table := fmt.Sprintf("hpn_%s_%s", utils.StringEncode16(q.ModName, 12), q.Table)

	return store.Data.Count(table, q.filter)
}

func (q *QuerySet) NodeList(fields, terms []string) api.NodeList {

	rsp := api.NodeList{}

	model, err := config.SpecNodeModel(q.ModName, q.Table)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeBadArgument,
			Message: "Spec Not Found",
		}
		return rsp
	}

	table := fmt.Sprintf("hpn_%s_%s", utils.StringEncode16(q.ModName, 12), q.Table)

	qs := store.Data.NewQueryer().
		Select(q.cols).
		From(table).
		Limit(q.limit).
		Offset(q.offset)

	if q.order != "" {
		qs.Order(q.order)
	} else {
		qs.Order("created desc")
	}

	qs.SetFilter(q.filter)

	rs, err := store.Data.Query(qs)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: "Can not pull database instance",
		}
		return rsp
	}

	var (
		termBufs     = map[string][]string{}
		termTaxonomy = map[string]api.Term{}
		ar_fields    = types.ArrayString(fields)
		ar_terms     = types.ArrayString(terms)
	)

	if len(rs) > 0 {

		for _, v := range rs {

			item := api.Node{
				ID:      v.Field("id").String(),
				PID:     v.Field("pid").String(),
				Status:  v.Field("status").Int16(),
				UserID:  v.Field("userid").String(),
				Created: v.Field("created").Uint32(),
				Updated: v.Field("updated").Uint32(),
			}

			if model.Extensions.AccessCounter {
				item.ExtAccessCounter = v.Field("ext_access_counter").Uint32()
			}

			if model.Extensions.CommentEnable {
				if model.Extensions.CommentPerEntry && v.Field("ext_comment_perentry").Bool() == false {
					item.ExtCommentEnable = false
					item.ExtCommentPerEntry = false
				} else {
					item.ExtCommentEnable = true
					item.ExtCommentPerEntry = true
				}
			}

			if model.Extensions.Permalink != "" {
				if pn := v.Field("ext_permalink_name").String(); pn != "" {
					item.ExtPermalinkName = pn
				} else {
					item.ExtPermalinkName = item.ID
				}
				item.SelfLink = fmt.Sprintf("%s", item.ExtPermalinkName)
			} else {
				item.SelfLink = fmt.Sprintf("%s.html", item.ID)
			}

			for _, field := range model.Fields {

				if field.Name != "title" && len(ar_fields) > 0 && !ar_fields.Has(field.Name) {
					continue
				}

				nodeField := api.NodeField{
					Name:  field.Name,
					Value: v.Field("field_" + field.Name).String(),
				}

				if field.Type == "text" &&
					len(v.Field("field_"+field.Name+"_attrs").String()) > 10 {

					var attrs types.KvPairs
					if err := v.Field("field_" + field.Name + "_attrs").JsonDecode(&attrs); err == nil {
						nodeField.Attrs = attrs
					}
				}

				if l := field.Attrs.Get("langs"); len(l) > 3 {

					if len(v.Field("field_"+field.Name+"_langs").String()) > 5 {
						var node_langs api.NodeFieldLangs
						if err := v.Field("field_" + field.Name + "_langs").JsonDecode(&node_langs); err == nil {
							nodeField.Langs = &node_langs
						}
					}
				}

				if field.Name == "title" {
					item.Title = nodeField.Value
				}

				item.Fields = append(item.Fields, &nodeField)
			}

			for _, term := range model.Terms {

				if len(ar_terms) > 0 && !ar_terms.Has(term.Meta.Name) {
					continue
				}

				termItem := api.NodeTerm{
					Name:  term.Meta.Name,
					Value: v.Field("term_" + term.Meta.Name).String(),
					Type:  term.Type,
				}

				item.Terms = append(item.Terms, termItem)
				if term.Type == api.TermTaxonomy {

					if te := TermTaxonomyCacheEntry(q.ModName, term.Meta.Name, v.Field("term_"+term.Meta.Name).Uint32()); te != nil {

						termTaxonomy[v.Field("term_"+term.Meta.Name).String()] = api.Term{
							ID:    te.ID,
							Title: te.Title,
						}

					} else if !utilx.ArrayContain(termItem.Value, termBufs[termItem.Name]) {
						termBufs[termItem.Name] = append(termBufs[termItem.Name], termItem.Value)
					}
				}
			}

			rsp.Items = append(rsp.Items, item)
		}
	}

	// Fetch Terms
	for _, term := range model.Terms {

		termids, ok := termBufs[term.Meta.Name]
		if !ok || len(termids) < 1 {
			continue
		}
		ids := []interface{}{}
		for _, tv := range termids {
			ids = append(ids, tv)
		}

		switch term.Type {

		case api.TermTaxonomy:

			table := fmt.Sprintf("hpt_%s_%s", utils.StringEncode16(q.ModName, 12), term.Meta.Name)
			qs := store.Data.NewQueryer().From(table).Limit(1000)
			qs.Where().And("id.in", ids...)

			if rs, err := store.Data.Query(qs); err == nil && len(rs) > 0 {

				for _, v := range rs {
					termTaxonomy[v.Field("id").String()] = api.Term{
						ID:    v.Field("id").Uint32(),
						Title: v.Field("title").String(),
					}
				}
			}
		}
	}

	//
	for k, v := range rsp.Items {

		for tk, tv := range v.Terms {

			if tv.Value == "" {
				continue
			}

			switch tv.Type {

			case api.TermTaxonomy:

				if tvs, ok := termTaxonomy[tv.Value]; ok {
					rsp.Items[k].Terms[tk].Items = append(rsp.Items[k].Terms[tk].Items, tvs)
				}

			case api.TermTag:

				tags := strings.Split(tv.Value, ",")

				for _, vtag := range tags {

					rsp.Items[k].Terms[tk].Items = append(rsp.Items[k].Terms[tk].Items, api.Term{
						Title: vtag,
					})
				}
			}
		}
	}

	rsp.Model = model

	rsp.Kind = "NodeList"

	if q.Pager {
		num, _ := store.Data.Count(table, q.filter)
		rsp.Meta.TotalResults = uint64(num)
		rsp.Meta.StartIndex = uint64(q.offset)
		rsp.Meta.ItemsPerList = uint64(q.limit)
	}

	// qryhash := q.Hash()

	// if rsp.Model.CacheTTL > 0 && entry.Title != "" {
	// 	store.CacheSetJson(qryhash, rsp, rsp.Model.CacheTTL)
	// }

	return rsp
}

func (q *QuerySet) NodeEntry() api.Node {

	var (
		rsp = api.Node{}
		err error
	)

	rsp.Model, err = config.SpecNodeModel(q.ModName, q.Table)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeBadArgument,
			Message: "Node Not Found",
		}
		return rsp
	}

	table := fmt.Sprintf("hpn_%s_%s", utils.StringEncode16(q.ModName, 12), q.Table)

	qs := store.Data.NewQueryer().
		Select(q.cols).
		From(table).
		Order(q.order).
		Limit(1).
		Offset(q.offset)

	qs.SetFilter(q.filter)
	// qs.Where().And("id", c.Params.Get("id"))

	rs, err := store.Data.Fetch(qs)
	if err != nil {
		rsp.Error = &types.ErrorMeta{
			Code:    api.ErrCodeInternalError,
			Message: err.Error(),
		}
		return rsp
	}

	for _, field := range rsp.Model.Fields {

		nodeField := api.NodeField{
			Name:  field.Name,
			Value: rs.Field("field_" + field.Name).String(),
		}

		if field.Type == "text" &&
			len(rs.Field("field_"+field.Name+"_attrs").String()) > 10 {

			var attrs types.KvPairs
			if err := rs.Field("field_" + field.Name + "_attrs").JsonDecode(&attrs); err == nil {
				nodeField.Attrs = attrs
			}
		}

		if l := field.Attrs.Get("langs"); len(l) > 3 {

			if len(rs.Field("field_"+field.Name+"_langs").String()) > 5 {
				var node_langs api.NodeFieldLangs
				if err := rs.Field("field_" + field.Name + "_langs").JsonDecode(&node_langs); err == nil {
					nodeField.Langs = &node_langs
				}
			}
		}

		if field.Name == "title" {
			rsp.Title = nodeField.Value
		}

		rsp.Fields = append(rsp.Fields, &nodeField)
	}

	for _, term := range rsp.Model.Terms {

		rsp.Terms = append(rsp.Terms, api.NodeTerm{
			Name:  term.Meta.Name,
			Value: rs.Field("term_" + term.Meta.Name).String(),
			Type:  term.Type,
		})
	}

	rsp.Terms = NodeTermQuery(q.ModName, rsp.Model, rsp.Terms)

	rsp.ID = rs.Field("id").String()
	rsp.Status = rs.Field("status").Int16()
	rsp.UserID = rs.Field("userid").String()
	rsp.Created = rs.Field("created").Uint32()
	rsp.Updated = rs.Field("updated").Uint32()

	if rsp.Model.Extensions.AccessCounter {
		rsp.ExtAccessCounter = rs.Field("ext_access_counter").Uint32()
	}

	if rsp.Model.Extensions.CommentEnable {
		if rsp.Model.Extensions.CommentPerEntry && rs.Field("ext_comment_perentry").Bool() == false {
			rsp.ExtCommentEnable = false
			rsp.ExtCommentPerEntry = false
		} else {
			rsp.ExtCommentEnable = true
			rsp.ExtCommentPerEntry = true
		}
	}

	if rsp.Model.Extensions.Permalink != "" {
		if pn := rs.Field("ext_permalink_name").String(); pn != "" {
			rsp.ExtPermalinkName = pn
		} else {
			rsp.ExtPermalinkName = rsp.ID
		}
	}

	if rsp.Model.Extensions.NodeRefer != "" {
		rsp.ExtNodeRefer = rs.Field("ext_node_refer").String()
	}

	rsp.Kind = "Node"

	// qryhash := q.Hash()

	// if rsp.Model.CacheTTL > 0 && entry.Title != "" {
	// 	store.CacheSetJson(qryhash, rsp, rsp.Model.CacheTTL)
	// }

	return rsp
}
