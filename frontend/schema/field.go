/*
Copyright © 2022 ConsenSys Software Inc.

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

package schema

// Field represent a schema Field and is analogous to reflect.StructField (but simplified)
type Field struct {
	Name       string
	Tag        string
	FullName   string
	Visibility Visibility
	Type       FieldType

	// if this field represent a struct, SubFields contains the struct members
	//
	// if it is an array or a slice of consistent elements, len(SubFields) == 1, the type of the
	// consistent element to repeat ArraySize times.
	//
	// SubFields are consistents if they have the exact same type, and if they contain slices,
	// exact same length.
	SubFields  []Field
	NbElements int
}

// FieldType represents the type a field is allowed to have in a gnark Schema
type FieldType uint8

const (
	Leaf FieldType = iota
	Array
	Slice
	Struct
)

// Visibility encodes a Variable (or wire) visibility
// Possible values are Unset, Internal, Secret or Public
type Visibility uint8

const (
	Unset Visibility = iota
	Internal
	Secret
	Public
	Virtual
)

func (v Visibility) String() string {
	switch v {
	case Internal:
		return "internal"
	case Secret:
		return "secret"
	case Public:
		return "public"
	case Virtual:
		return "virtual"
	}

	return "unset"
}
