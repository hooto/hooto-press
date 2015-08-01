// Copyright 2015 lessOS.com, All rights reserved.
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

package modset

const (
	dsTplNodeModels = `
{
    "columns": [
        {
            "name": "id",
            "type": "string",
            "length": "16"
        },
        {
            "name": "pid",
            "type": "string",
            "length": "16"
        },
        {
            "name": "state",
            "type": "int16"
        },
        {
            "name": "userid",
            "type": "string",
            "length": "10"
        },
        {
            "name": "title",
            "type": "string",
            "length": "100"
        },
        {
            "name": "created",
            "type": "datetime"
        },
        {
            "name": "updated",
            "type": "datetime"
        }
    ],
    "indexes": [
        {
            "name": "PRIMARY",
            "type": 3,
            "cols": ["id"]
        },
        {
            "name": "pid",
            "type": 1,
            "cols": ["pid"]
        },
        {
            "name": "state",
            "type": 1,
            "cols": ["state"]
        },
        {
            "name": "userid",
            "type": 1,
            "cols": ["userid"]
        },
        {
            "name": "created",
            "type": 1,
            "cols": ["created"]
        },
        {
            "name": "updated",
            "type": 1,
            "cols": ["updated"]
        }
    ]
}
`
	dsTplTermModels = `
{
    "name": "template",
    "columns": [
        {
            "name": "id",
            "type": "uint32",
            "IncrAble": true
        },
        {
            "name": "state",
            "type": "int16"
        },
        {
            "name": "userid",
            "type": "string",
            "length": "10"
        },
        {
            "name": "title",
            "type": "string",
            "length": "100"
        },
        {
            "name": "created",
            "type": "datetime"
        },
        {
            "name": "updated",
            "type": "datetime"
        }
    ],
    "indexes": [
        {
            "name": "PRIMARY",
            "type": 3,
            "cols": ["id"]
        },
        {
            "name": "state",
            "type": 1,
            "cols": ["state"]
        },
        {
            "name": "userid",
            "type": 1,
            "cols": ["userid"]
        },
        {
            "name": "created",
            "type": 1,
            "cols": ["created"]
        },
        {
            "name": "updated",
            "type": 1,
            "cols": ["updated"]
        }
    ]
}
`
)
