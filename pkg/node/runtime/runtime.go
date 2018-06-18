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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/network"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/endpoint"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/pod"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/volume"
	"time"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logNodeRuntimePrefix = "%s"
)

type Runtime struct {
	ctx  context.Context
	spec chan *types.NodeManifest
}

func (r *Runtime) Restore() {
	log.Debugf("%s:restore:> restore init", logNodeRuntimePrefix)
	network.Restore(r.ctx)
	volume.Restore(r.ctx)
	pod.Restore(r.ctx)
	endpoint.Restore(r.ctx)
}

func (r *Runtime) Provision(ctx context.Context, spec *types.NodeManifest) error {

	log.Debugf("%s> provision init", logNodeRuntimePrefix)

	log.Debugf("%s> provision networks", logNodeRuntimePrefix)
	for cidr, n := range spec.Network {
		log.Debugf("network: %v", n)
		if err := network.Manage(ctx, cidr, &n); err != nil {
			log.Errorf("Network [%s] create err: %s", n.CIDR, err.Error())
		}
	}

	log.Debugf("%s> provision pods", logNodeRuntimePrefix)
	for p, spec := range spec.Pods {
		log.Debugf("pod: %v", p)
		if err := pod.Manage(ctx, p, &spec); err != nil {
			log.Errorf("Pod [%s] manage err: %s", p, err.Error())
		}
	}

	log.Debugf("%s> provision endpoints", logNodeRuntimePrefix)
	for e, spec := range spec.Endpoints {
		log.Debugf("endpoint: %v", e)
		if err := endpoint.Manage(ctx, e, &spec); err != nil {
			log.Errorf("Endpoint [%s] manage err: %s", e, err.Error())
		}
	}

	log.Debugf("%s> provision volumes", logNodeRuntimePrefix)
	for _, v := range spec.Volumes {
		log.Debugf("volume: %v", v)
	}

	return nil
}

func (r *Runtime) Subscribe() {

	log.Debugf("%s:subscribe:> subscribe init", logNodeRuntimePrefix)
	pc := make(chan string)

	go func() {

		for {
			select {
			case p := <-pc:
				log.Debugf("%s:subscribe:> new pod state event: %s", logNodeRuntimePrefix, p)
				events.NewPodStatusEvent(r.ctx, p)
			}
		}
	}()

	envs.Get().GetCRI().Subscribe(r.ctx, envs.Get().GetState().Pods(), pc)
}

func (r *Runtime) Connect(ctx context.Context) error {

	log.Debugf("%s:connect:> connect init", logNodeRuntimePrefix)
	if err := events.NewConnectEvent(ctx); err != nil {
		log.Errorf("%s:connect:> connect err: %s", logNodeRuntimePrefix, err.Error())
		return err
	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			if err := events.NewStatusEvent(ctx); err != nil {
				log.Errorf("%s:connect:> send status err: %s", logNodeRuntimePrefix, err.Error())
			}
		}
	}(ctx)

	return nil
}

func (r *Runtime) GetSpec(ctx context.Context) error {

	log.Debugf("%s:getspec:> getspec request init", logNodeRuntimePrefix)

	var (
		c = envs.Get().GetClient()
	)

	spec, err := c.GetSpec(ctx)
	if err != nil {
		log.Errorf("%s:getspec:> request err: %s", logNodeRuntimePrefix, err.Error())
		return err
	}

	if spec == nil {
		log.Warnf("%s:getspec:> new spec is nil", logNodeRuntimePrefix)
		return nil
	}

	r.spec <- spec.Decode()
	return nil
}

func (r *Runtime) Clean(ctx context.Context, manifest *types.NodeManifest) error {

	log.Debugf("%s> clean up endpoints", logNodeRuntimePrefix)
	endpoints := envs.Get().GetState().Endpoints().GetEndpoints()
	for e := range endpoints {
		if _, ok := manifest.Endpoints[e]; !ok {
			endpoint.Destroy(context.Background(), e, endpoints[e])
		}
	}

	log.Debugf("%s> clean up pods", logNodeRuntimePrefix)
	pods := envs.Get().GetState().Pods().GetPods()

	for k := range pods {
		if _, ok := manifest.Pods[k]; !ok {
			pod.Destroy(context.Background(), k, pods[k])
		}
	}

	log.Debugf("%s> clean up networks", logNodeRuntimePrefix)
	nets := envs.Get().GetState().Networks().GetSubnets()

	for cidr := range nets {
		if _, ok := manifest.Network[cidr]; !ok {
			network.Destroy(ctx, cidr)
		}
	}

	return nil
}

func (r *Runtime) Loop() {
	log.Debugf("%s:loop:> start runtime loop", logNodeRuntimePrefix)

	var clean = true

	go func(ctx context.Context) {
		for {
			select {
			case spec := <-r.spec:
				log.Debugf("%s:loop:> provision new spec", logNodeRuntimePrefix)

				if clean {
					if err := r.Clean(ctx, spec); err != nil {
						log.Errorf("%s:loop:> clean err: %s", err.Error())
						continue
					}
					clean = false
				}

				if err := r.Provision(ctx, spec); err != nil {
					log.Errorf("%s:loop:> provision new spec err: %s", logNodeRuntimePrefix, err.Error())
				}
			}
		}
	}(r.ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			err := r.GetSpec(r.ctx)
			if err != nil {
				log.Debugf("%s:loop:> new spec request err: %s", logNodeRuntimePrefix, err.Error())
			}
		}
	}(context.Background())

	err := r.GetSpec(r.ctx)
	if err != nil {
		log.Debugf("%s:loop:> new spec request err: %s", logNodeRuntimePrefix, err.Error())
	}
}

func NewRuntime(ctx context.Context) *Runtime {
	r := Runtime{
		ctx:  ctx,
		spec: make(chan *types.NodeManifest),
	}

	return &r
}
