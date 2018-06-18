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

package types_test

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEqual(t *testing.T) {

	var network = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	var asset = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assert.Equal(t, true, network.Equal(asset), "equal")
}

func TestNotEqual(t *testing.T) {

	var network = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	var assets = make(map[string]*types.NetworkSpec)

	assets["type"] = &types.NetworkSpec{
		Type: "vlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["subnet"] = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/22",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["subnet & type"] = &types.NetworkSpec{
		Type: "vlan",
		CIDR: "10.0.0.0/22",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index"] = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index & type"] = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index & subnet"] = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface name"] = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.0",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface addr"] = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.2",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface haddr"] = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:c8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["addr"] = &types.NetworkSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: types.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.1",
	}

	for attr, asset := range assets {
		assert.Equal(t, false, network.Equal(asset), attr)
	}
}
