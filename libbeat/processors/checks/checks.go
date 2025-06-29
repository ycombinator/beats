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

package checks

import (
	"fmt"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/processors"
	"github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
)

// ConfigChecked returns a wrapper that will validate the configuration using
// the passed checks before invoking the original constructor.
func ConfigChecked(
	constr processors.Constructor,
	checks ...func(*config.C) error,
) processors.Constructor {
	validator := checkAll(checks...)
	return func(cfg *config.C, log *logp.Logger) (beat.Processor, error) {
		err := validator(cfg)
		if err != nil {
			return nil, fmt.Errorf("%w in %v", err, cfg.Path())
		}

		return constr(cfg, log)
	}
}

func checkAll(checks ...func(*config.C) error) func(*config.C) error {
	return func(c *config.C) error {
		for _, check := range checks {
			if err := check(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// RequireFields checks that the required fields are present in the configuration.
func RequireFields(fields ...string) func(*config.C) error {
	return func(cfg *config.C) error {
		for _, field := range fields {
			if !cfg.HasField(field) {
				return fmt.Errorf("missing %v option", field)
			}
		}
		return nil
	}
}

// AllowedFields checks that only allowed fields are used in the configuration.
func AllowedFields(fields ...string) func(*config.C) error {
	return func(cfg *config.C) error {
		for _, field := range cfg.GetFields() {
			found := false
			for _, allowed := range fields {
				if field == allowed {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("unexpected %v option", field)
			}
		}
		return nil
	}
}

// MutuallyExclusiveRequiredFields checks that only one of the given
// fields is used at the same time. It is an error for none of the fields to be
// present.
func MutuallyExclusiveRequiredFields(fields ...string) func(*config.C) error {
	return func(cfg *config.C) error {
		var foundField string
		for _, field := range cfg.GetFields() {
			for _, f := range fields {
				if field == f {
					if len(foundField) == 0 {
						foundField = field
					} else {
						return fmt.Errorf("field %s and %s are mutually exclusive", foundField, field)
					}
				}
			}
		}

		if len(foundField) == 0 {
			return fmt.Errorf("missing option, select one from %v", fields)
		}
		return nil
	}
}
