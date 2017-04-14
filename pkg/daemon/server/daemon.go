//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package server

import (
	c "context"
	"github.com/jawher/mow.cli"
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/storage"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/lastbackend/lastbackend/pkg/wss"
	"os"
	"os/signal"
	"syscall"
)

func Daemon(cmd *cli.Cmd) {

	var (
		ctx = context.Get()
		cfg = config.Get()
	)

	cmd.Spec = "[-d]"

	var debug = cmd.Bool(cli.BoolOpt{
		Name: "d debug", Desc: "Enable debug mode",
		EnvVar: "DEBUG", Value: false, HideValue: true,
	})
	var secretToken = cmd.String(cli.StringOpt{
		Name: "secret-token", Desc: "Secret token for signature",
		EnvVar: "SECRET-TOKEN", Value: "b8tX!ae4", HideValue: true,
	})
	var templateHost = cmd.String(cli.StringOpt{
		Name: "template-host", Desc: "Address for template registry",
		EnvVar: "TEMPLATE-HOST", Value: "http://localhost:3003", HideValue: true,
	})
	var proxyServerPort = cmd.Int(cli.IntOpt{
		Name: "proxy-server-port", Desc: "Proxy server port",
		EnvVar: "PROXY-SERVER-PORT", Value: 2968, HideValue: true,
	})
	var httpServerHost = cmd.String(cli.StringOpt{
		Name: "http-server-host", Desc: "Http server host",
		EnvVar: "HTTP-SERVER-HOST", Value: "", HideValue: true,
	})
	var httpServerPort = cmd.Int(cli.IntOpt{
		Name: "http-server-port", Desc: "Http server port",
		EnvVar: "HTTP-SERVER-PORT", Value: 2967, HideValue: true,
	})
	var registryServer = cmd.String(cli.StringOpt{
		Name: "registry-server", Desc: "Http server port",
		EnvVar: "REGISTRY-SERVER", Value: "hub.registry.net", HideValue: true,
	})
	var registryUsername = cmd.String(cli.StringOpt{
		Name: "registry-username", Desc: "Http server port",
		EnvVar: "REGISTRY-USERNAME", Value: "demo", HideValue: true,
	})
	var registryPassword = cmd.String(cli.StringOpt{
		Name: "registry-password", Desc: "Http server port",
		EnvVar: "REGISTRY-PASSWORD", Value: "IU1yxkTD", HideValue: true,
	})
	var etcdEndpoints = cmd.Strings(cli.StringsOpt{
		Name: "etcd-endpoints", Desc: "Set etcd endpoints list",
		EnvVar: "ETCD-ENDPOINTS", Value: []string{"localhost:2379"}, HideValue: true,
	})
	var etcdTlsKey = cmd.String(cli.StringOpt{
		Name: "etcd-tls-key", Desc: "Etcd tls key",
		EnvVar: "ETCD-TLS-KEY", Value: "", HideValue: true,
	})
	var etcdTlsSert = cmd.String(cli.StringOpt{
		Name: "etcd-tls-cert", Desc: "Etcd tls cert",
		EnvVar: "ETCD-TLS-CERT", Value: "", HideValue: true,
	})
	var etcdTlsCA = cmd.String(cli.StringOpt{
		Name: "etcd-tls-ca", Desc: "Etcd tls ca",
		EnvVar: "ETCD-TLS-CA", Value: "", HideValue: true,
	})

	cmd.Before = func() {

		cfg.Debug = *debug
		cfg.SecretToken = *secretToken
		cfg.TemplateRegistry.Host = *templateHost
		cfg.ProxyServer.Port = *proxyServerPort
		cfg.HttpServer.Host = *httpServerHost
		cfg.HttpServer.Port = *httpServerPort
		cfg.Registry.Server = *registryServer
		cfg.Registry.Username = *registryUsername
		cfg.Registry.Password = *registryPassword
		cfg.Etcd.Endpoints = *etcdEndpoints
		cfg.Etcd.TLS.Key = *etcdTlsKey
		cfg.Etcd.TLS.Cert = *etcdTlsSert
		cfg.Etcd.TLS.CA = *etcdTlsCA

		ctx.SetConfig(cfg)
		ctx.SetHttpTemplateRegistry(http.New(cfg.TemplateRegistry.Host))
		ctx.SetLogger(logger.New(cfg.Debug, 9))
		ctx.SetWssHub(new(wss.Hub))
		strg, err := storage.Get(cfg.GetEtcdDB())
		if err != nil {
			panic(err)
		}
		ctx.SetStorage(strg)

		ns, err := ctx.GetStorage().Namespace().GetByName(c.Background(), "demo")
		if err != nil {
			ctx.GetLogger().Error(err)
			return
		}

		ctx.GetLogger().Debug(ns)
	}

	cmd.Action = func() {

		var (
			log  = ctx.GetLogger()
			sigs = make(chan os.Signal)
			done = make(chan bool, 1)
		)

		go func() {
			if err := Listen(cfg.HttpServer.Host, cfg.HttpServer.Port); err != nil {
				log.Warnf("Http server start error: %s", err.Error())
			}
		}()

		// Handle SIGINT and SIGTERM.
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for {
				select {
				case <-sigs:
					done <- true
					return
				}
			}
		}()

		<-done

		log.Info("Handle SIGINT and SIGTERM.")
	}
}
