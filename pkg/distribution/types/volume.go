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

package types

import "fmt"

const VOLUMETYPELOCAL = "local"

// swagger:ignore
// swagger:model types_volume
type Volume struct {
	Runtime
	// Volume meta
	Meta VolumeMeta `json:"meta" yaml:"meta"`
	// Volume spec
	Spec VolumeSpec `json:"spec" yaml:"spec"`
	// Volume status
	Status VolumeStatus `json:"status" yaml:"status"`
}

// swagger:ignore
// swagger:model types_volume_map
type VolumeMap struct {
	Runtime
	Items map[string]*Volume
}

// swagger:ignore
// swagger:model types_volume_list
type VolumeList struct {
	Runtime
	Items []*Volume
}

// swagger:ignore
// swagger:model types_volume_meta
type VolumeMeta struct {
	Meta
	Namespace string `json:"namespace"`
}

// swagger:model types_volume_spec
type VolumeSpec struct {
	State     VolumeSpecState     `json:"state"`
	Type      string              `json:"type"`
	Path      string              `json:"path"`
	Mode      string              `json:"mode"`
	Resources SpecVolumeResources `json:"resources"`
}

// swagger:model types_volume_spec_state
type VolumeSpecState struct {
	Destroy bool `json:"destroy"`
}

// swagger:ignore
// swagger:model types_volume_create
type VolumeCreateOptions struct {
}

// swagger:ignore
// swagger:model types_volume_update
type VolumeUpdateOptions struct {
}

type VolumeStatus struct {
	// volume state
	State string `json:"state" yaml:"state"`
	// volume status
	Status VolumeState `json:"status" yaml:"status"`
	// volume status message
	Message string `json:"message" yaml:"message"`
}

// swagger:ignore
// swagger:model types_volume_status
type VolumeState struct {
	Type string `json:"type" yaml:"type"`
	// Volume root path
	Path string `json:"path" yaml:"path"`
	// Volume state ready
	Ready bool `json:"ready" yaml:"ready"`
}

func (v *Volume) SelfLink() string {
	if v.Meta.SelfLink == "" {
		v.Meta.SelfLink = v.CreateSelfLink(v.Meta.Namespace, v.Meta.Name)
	}
	return v.Meta.SelfLink
}

func (v *Volume) CreateSelfLink(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

func NewVolumeList() *VolumeList {
	dm := new(VolumeList)
	dm.Items = make([]*Volume, 0)
	return dm
}

func NewVolumeMap() *VolumeMap {
	dm := new(VolumeMap)
	dm.Items = make(map[string]*Volume)
	return dm
}
