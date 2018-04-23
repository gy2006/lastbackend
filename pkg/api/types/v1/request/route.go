//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package request

// swagger:model request_route_create
type RouteCreateOptions struct {
	Name     string        `json:"name"`
	Security bool          `json:"security"`
	Rules    []RulesOption `json:"rules"`
}

// swagger:model request_route_update
type RouteUpdateOptions struct {
	Security bool          `json:"security"`
	Rules    []RulesOption `json:"rules"`
}

// swagger:ignore
// swagger:model request_route_remove
type RouteRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:model request_route_rules
type RulesOption struct {
	Service string `json:"service"`
	Path    string `json:"path"`
	Port    int    `json:"port"`
}
