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

package node

import (
	"net/http"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"strings"
)

const logLevel = 2

func NodeGetH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: list node")

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("Handler: Node: node `%s` not found", nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	response, err := v1.View().Node().New(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeGetSpecH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: list node")

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		cid = utils.Vars(r)["cluster"]
		nid = utils.Vars(r)["node"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("Handler: Node: node `%s` not found", cid)
		errors.New("node").NotFound().Http(w)
		return
	}

	spec, err  := nm.GetSpec(n)
	if err != nil {
		log.V(logLevel).Warnf("Handler: Node: node `%s` not found", cid)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Node().NewSpec(spec).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeListH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: list node")

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	nodes, err := nm.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Node().NewList(nodes).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["node"]

	log.V(logLevel).Debugf("Handler: Node: update node `%s`", nid)

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := new(request.NodeUpdateOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Node: validation incoming data", err)
		err.Http(w)
		return
	}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	err = nm.SetMeta(n, opts.Meta)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: update node `%s` err: %s", nid, err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Node().New(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetInfoH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set info")

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	// request body struct
	opts := new(request.NodeInfoOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Node: validation incoming data", err)
		err.Http(w)
		return
	}

	node, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := nm.SetInfo(node, types.NodeInfo{
		Hostname: opts.Hostname,
		Architecture: opts.Architecture,
		OSName: opts.OSName,
		OSType: opts.OSType,
		ExternalIP: opts.ExternalIP,
		InternalIP: opts.InternalIP,
	}); err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetStateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set state")

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	// request body struct
	opts := new(request.NodeStateOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Node: validation incoming data", err)
		err.Http(w)
		return
	}

	node, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := nm.SetState(node, types.NodeState{
		Capacity: opts.Capacity,
		Allocated: opts.Allocated,
	}); err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetPodStatusH(w http.ResponseWriter, r *http.Request) {

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		pm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
		pid = utils.Vars(r)["pod"]
	)

	// request body struct
	opts := new(request.NodePodStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Node: validation incoming data", err)
		err.Http(w)
		return
	}

	_, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.NotFound(w)
		return
	}

	keys:=strings.Split(pid, ":")
	if len(keys) != 5 {
		log.V(logLevel).Errorf("Handler: Node: invalid pod selflink err: %s", pid)
		errors.HTTP.BadRequest(w)
		return
	}

	pod, err := pm.Get(keys[1],keys[2],keys[3],keys[4])
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: pod not found selflink err: %s", pid)
		errors.HTTP.NotFound(w)
		return
	}

	if err := pm.SetStatus(pod, &types.PodStatus{
		Stage: opts.Stage,
		Message: opts.Message,
		Steps: opts.Steps,
		Network: opts.Network,
		Containers: opts.Containers,
	}); err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetVolumeStatusH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set volume state")

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		vm = distribution.NewVolumeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
		vid = utils.Vars(r)["volume"]
	)

	// request body struct
	opts := new(request.NodeVolumeStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Node: validation incoming data", err)
		err.Http(w)
		return
	}

	_, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	keys:=strings.Split(vid, ":")
	if len(keys) != 3 {
		log.V(logLevel).Errorf("Handler: Node: invalid volume selflink err: %s", vid)
		errors.HTTP.BadRequest(w)
		return
	}

	volume, err := vm.Get(keys[1],keys[2])
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: pod not found selflink err: %s", vid)
		errors.HTTP.NotFound(w)
		return
	}

	if err := vm.SetStatus(volume, &types.VolumeStatus{
		Stage: opts.Stage,
		Message: opts.Message,
	}); err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetRouteStatusH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set route state")

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
		vid = utils.Vars(r)["route"]
	)

	// request body struct
	opts := new(request.NodeRouteStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Node: validation incoming data", err)
		err.Http(w)
		return
	}

	_, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	keys:=strings.Split(vid, ":")
	if len(keys) != 3 {
		log.V(logLevel).Errorf("Handler: Node: invalid route selflink err: %s", vid)
		errors.HTTP.BadRequest(w)
		return
	}

	route, err := rm.Get(keys[1],keys[2])
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: pod not found selflink err: %s", vid)
		errors.HTTP.NotFound(w)
		return
	}

	if err := rm.SetStatus(route, &types.RouteStatus{
		Stage: opts.Stage,
		Message: opts.Message,
	}); err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeRemoveH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: create node")

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	// request body struct
	opts := v1.Request().Node().RemoveOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Node: validation incoming data err: %s", err)
		err.Http(w)
		return
	}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: remove node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	if n == nil {
		log.V(logLevel).Warnf("Handler: Node: remove node `%s` not found", nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	if err := nm.Remove(n); err != nil {
		log.V(logLevel).Errorf("Handler: Node: remove node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}
