// Copyright 2020 ConsenSys Software Inc.
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

// Code generated by gnark DO NOT EDIT

package witness

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/schema"
	"github.com/consensys/gnark/internal/backend/compiled"

	"github.com/consensys/gnark-crypto/ecc/bw6-633/fr"

	curve "github.com/consensys/gnark-crypto/ecc/bw6-633"
)

type Witness []fr.Element

// that's public and modifiable. Fine since it's internal for now
// TODO @gbotrel get rid of that through generics
var T = reflect.TypeOf(fr.Element{})

// WriteTo encodes witness to writer (implements io.WriterTo)
func (witness *Witness) WriteTo(w io.Writer) (int64, error) {
	// encode slice length
	if err := binary.Write(w, binary.BigEndian, uint32(len(*witness))); err != nil {
		return 0, err
	}

	enc := curve.NewEncoder(w)
	for i := 0; i < len(*witness); i++ {
		if err := enc.Encode(&(*witness)[i]); err != nil {
			return enc.BytesWritten() + 4, err
		}
	}
	return enc.BytesWritten() + 4, nil
}

func (witness *Witness) ReadFrom(r io.Reader) (int64, error) {

	var buf [4]byte
	if read, err := io.ReadFull(r, buf[:4]); err != nil {
		return int64(read), err
	}
	sliceLen := binary.BigEndian.Uint32(buf[:4])

	if len(*witness) != int(sliceLen) {
		*witness = make([]fr.Element, sliceLen)
	}

	dec := curve.NewDecoder(r)

	for i := 0; i < int(sliceLen); i++ {
		if err := dec.Decode(&(*witness)[i]); err != nil {
			return dec.BytesRead() + 4, err
		}
	}

	return dec.BytesRead() + 4, nil
}

// FromAssignment extracts the witness and its schema
func (witness *Witness) FromAssignment(w interface{}, publicOnly bool) (*schema.Schema, error) {
	s, err := schema.Parse(w, tVariable, nil)
	if err != nil {
		return nil, err
	}
	nbSecret, nbPublic := s.NbSecret, s.NbPublic

	if publicOnly {
		nbSecret = 0
	}

	if len(*witness) < (nbPublic + nbSecret) {
		(*witness) = make(Witness, nbPublic+nbSecret)
	} else {
		(*witness) = (*witness)[:nbPublic+nbSecret]
	}

	var i, j int // indexes for secret / public variables
	i = nbPublic // offset

	var collectHandler schema.LeafHandler = func(visibility compiled.Visibility, name string, tInput reflect.Value) error {
		if publicOnly && visibility != compiled.Public {
			return nil
		}
		if tInput.IsNil() {
			return fmt.Errorf("when parsing variable %s: missing assignment", name)
		}
		v := tInput.Interface().(frontend.Variable)

		if v == nil {
			return fmt.Errorf("when parsing variable %s: missing assignment", name)
		}

		if !publicOnly && visibility == compiled.Secret {
			if _, err := (*witness)[i].SetInterface(v); err != nil {
				return fmt.Errorf("when parsing variable %s: %v", name, err)
			}
			i++
		} else if visibility == compiled.Public {
			if _, err := (*witness)[j].SetInterface(v); err != nil {
				return fmt.Errorf("when parsing variable %s: %v", name, err)
			}
			j++
		}
		return nil
	}
	return schema.Parse(w, tVariable, collectHandler)
}

// VectorToAssignment sets to leaf values to witness underlying vector element values (in order)
// see witness.MarshalBinary protocol description
func (witness *Witness) VectorToAssignment(to interface{}, toLeafType reflect.Type, publicOnly bool) {
	i := 0
	setAddr := toLeafType.Kind() == reflect.Ptr
	setHandler := func(v compiled.Visibility) schema.LeafHandler {
		return func(visibility compiled.Visibility, name string, tInput reflect.Value) error {
			if visibility == v {
				if setAddr {
					tInput.Set(reflect.ValueOf((&(*witness)[i])))
				} else {
					tInput.Set(reflect.ValueOf(((*witness)[i])))
				}

				i++
			}
			return nil
		}
	}
	_, _ = schema.Parse(to, toLeafType, setHandler(compiled.Public))
	if publicOnly {
		return
	}
	_, _ = schema.Parse(to, toLeafType, setHandler(compiled.Secret))

}

func (witness *Witness) String() string {
	var sbb strings.Builder
	sbb.WriteByte('[')
	for i := 0; i < len(*witness); i++ {
		sbb.WriteString((*witness)[i].String())
		sbb.WriteByte(',')
	}
	sbb.WriteByte(']')
	return sbb.String()
}

var tVariable reflect.Type

func init() {
	tVariable = reflect.ValueOf(struct{ A frontend.Variable }{}).FieldByName("A").Type()
}
