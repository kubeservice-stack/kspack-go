/*
Copyright 2023 The KubeService-Stack Authors.

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

package pack

const (
	KSPACK_INVALID      = 0x00
	KSPACK_OBJECT       = 0x10
	KSPACK_ARRAY        = 0x20
	KSPACK_STRING       = 0x50
	KSPACK_BINARY       = 0x60
	KSPACK_INT8         = 0x11
	KSPACK_INT16        = 0x12
	KSPACK_INT32        = 0x14
	KSPACK_INT64        = 0x18
	KSPACK_UINT8        = 0x21
	KSPACK_UINT16       = 0x22
	KSPACK_UINT32       = 0x24
	KSPACK_UINT64       = 0x28
	KSPACK_BOOL         = 0x31
	KSPACK_FLOAT        = 0x44
	KSPACK_DOUBLE       = 0x48
	KSPACK_DATE         = 0x58
	KSPACK_NULL         = 0x61
	KSPACK_SHORT_ITEM   = 0x80
	KSPACK_FIXED_ITEM   = 0xf0
	KSPACK_DELETED_ITEM = 0x70

	KSPACK_SHORT_STRING = KSPACK_STRING | KSPACK_SHORT_ITEM
	KSPACK_SHORT_BINARY = KSPACK_BINARY | KSPACK_SHORT_ITEM

	KSPACK_KEY_MAX_LEN = 254

	MAX_SHORT_VITEM_LEN = 255
)
