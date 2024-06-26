/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package ionhash

import "github.com/amazon-ion/ion-go/ion"

type scalarSerializer struct {
	baseSerializer
}

func newScalarSerializer(hashFunction IonHasher, depth int) serializer {
	return &scalarSerializer{baseSerializer{hashFunction: hashFunction, depth: depth}}
}

func (ss *scalarSerializer) scalar(ionValue hashValue) error {
	err := ss.handleAnnotationsBegin(ionValue)
	if err != nil {
		return err
	}

	err = ss.beginMarker()
	if err != nil {
		return err
	}

	var ionVal interface{}
	var ionType ion.Type
	if ionValue.IsNull() {
		ionVal = nil
		ionType = ion.NoType
	} else {
		ionVal, err = ionValue.value()
		if err != nil {
			return err
		}
		ionType = ionValue.Type()
	}

	scalarBytes, err := ss.getBytes(ionValue.Type(), ionVal, ionValue.IsNull())
	if err != nil {
		return err
	}

	var symbol *ion.SymbolToken = nil
	if ionValue.Type() == ion.SymbolType {
		if token, ok := ionVal.(ion.SymbolToken); ok {
			symbol = &token
		} else if token, ok := ionVal.(*ion.SymbolToken); ok {
			symbol = token
		}
	}

	tq, representation, err := ss.scalarOrNullSplitParts(ionType, symbol, ionValue.IsNull(), scalarBytes)
	if err != nil {
		return err
	}

	err = ss.write([]byte{tq})
	if err != nil {
		return err
	}

	if len(representation) > 0 {
		err = ss.write(escape(representation))
		if err != nil {
			return err
		}
	}

	err = ss.endMarker()
	if err != nil {
		return err
	}

	err = ss.handleAnnotationsEnd(ionValue, false)
	if err != nil {
		return err
	}

	return nil
}

func (ss *scalarSerializer) handleAnnotationsBegin(ionValue hashValue) error {
	return ss.baseSerializer.handleAnnotationsBegin(ionValue, false)
}
