// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package inputmon

import (
	"strings"

	"github.com/elastic/elastic-agent-libs/monitoring"
)

// NewInputRegistry returns a new monitoring.Registry for metrics related to
// an input instance. The returned registry will be initialized with a static
// string values for the input and id. When the input stops it should invoke
// the returned cancel function to unregister the metrics. For testing purposes
// an optional monitoring.Registry may be provided as an alternative to using
// the global 'dataset' monitoring namespace. The inputType and id must be
// non-empty for the metrics to be published to the global 'dataset' monitoring
// namespace.
func NewInputRegistry(inputType, id string, optionalParent *monitoring.Registry) (reg *monitoring.Registry, cancel func()) {
	// Use the default registry unless one was provided (this would be for testing).
	parentRegistry := optionalParent
	if parentRegistry == nil {
		parentRegistry = globalRegistry()
	}

	// If an ID has not been assigned to an input then metrics cannot be exposed
	// in the global metric registry. The returned registry still behaves the same.
	if (id == "" || inputType == "") && parentRegistry == globalRegistry() {
		// Null route metrics without ID or input type.
		parentRegistry = monitoring.NewRegistry()
	}

	// Sanitize dots from the id because they created nested objects within
	// the monitoring registry, and we want a consistent flat level of nesting
	key := sanitizeID(id)

	reg = parentRegistry.NewRegistry(key)
	monitoring.NewString(reg, "input").Set(inputType)
	monitoring.NewString(reg, "id").Set(id)

	return reg, func() { parentRegistry.Remove(key) }
}

func sanitizeID(id string) string {
	return strings.ReplaceAll(id, ".", "_")
}

func globalRegistry() *monitoring.Registry {
	return monitoring.GetNamespace("dataset").GetRegistry()
}
