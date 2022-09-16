// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package v2

import low "github.com/pb33f/libopenapi/datamodel/low/v2"

type SecurityDefinitions struct {
	Definitions map[string]*SecurityScheme
	low         *low.SecurityDefinitions
}

func NewSecurityDefinitions(definitions *low.SecurityDefinitions) *SecurityDefinitions {
	sd := new(SecurityDefinitions)
	sd.low = definitions
	schemes := make(map[string]*SecurityScheme)
	for k := range definitions.Definitions {
		schemes[k.Value] = NewSecurityScheme(definitions.Definitions[k].Value)
	}
	sd.Definitions = schemes
	return sd
}

func (sd *SecurityDefinitions) GoLow() *low.SecurityDefinitions {
	return sd.low
}