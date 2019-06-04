# RPaaS v2

[![Build Status](https://travis-ci.org/tsuru/rpaas-operator.svg?branch=master)](https://travis-ci.org/tsuru/rpaas-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/tsuru/rpaas-operator)](https://goreportcard.com/report/github.com/tsuru/rpaas-operator)

NOTE: This project is the replacement of the [RPaaS][rpaas-v1-repository],
so we'll refer to it only as RPaaS v2 (although there are no breaking
changes between them).

---

## Description

RPaaS, which stands for Reverse Proxy as a Service, provides a easy and fast
way to manage [NGINX][nginx-site]-based reverse proxy into a cloud
infrastructure. Such reverse proxy (aka RPaaS instance) handles incoming HTTP
request and forward to the configured destination application (backend).

Furthermore, it supports adding TLS terminating, cache (according to HTTP cache
headers from backend's response), purge cached objects, scale up/down the
instances and so forth.

RPaaS v2 is broken into two parts: Operator and API.

### Operator

An Kubernetes application, built following the [Operator framework][kubernetes-operator],
which transform the high-level RPaaS Custom Resources into more basic Kubernetes
objects (such as Secret, ConfigMap and so on) and Nginx Custom Resources
(provided by [nginx-operator][nginx-operator-repository] project).

### API

Just a web API which manages the high-level RPaaS Custom Resources inside the
Kubernetes cluster. Unlike the Operator, it does not need run inside the
Kubernetes cluster but needs the credentials to manipulate basic Kubernetes
object as well as RPaaS Custom Resources.

## Installing

Prior any installation step, you need to have a Kubernetes cluster and the `kubectl`
command-line tool must be configured to communicate with your cluster.

First, you need creating the Rpaas CRDs.

```bash
$ for crd in $(ls ./deploy/crds/*_crd.yaml); do kubectl apply -f ${crd}; done
```

Create a new namespace to group the RPaaS core components. For instance, we'll use
the namespace `rpaas-system`.

```bash
$ kubectl create namespace rpaas-system
```

Run the NGINX operator, Rpaas Operator and Rpaas API using the YAML manifests.

```bash
$ kubectl -n rpaas-system -f vendor/github.com/tsuru/nginx-operator/deploy/
$ sed -E 's/(namespace): .+/\1: rpaas-system/' vendor/github.com/tsuru/nginx-operator/deploy/role_binding.yaml | \
  kubectl -n rpaas-system -f -

$ kubectl -n rpaas-system -f deploy/
$ sed -E 's/(namespace): .+/\1: rpaas-system/' deploy/role_binding.yaml | \
  kubectl -n rpaas-system -f -
```

Now, you can access the RPaaS API using the NodePort address of the service 
`rpaas-api`.

```bash
$ kubectl -n rpaas-system services rpaas-api

$ curl http://<node-ip>:<svc-port>/healthcheck
> OK
```

## Contributing

TODO

## License

RPaaS v2 is an open source project authored by [Globo.com][opensource-globocom]
and released under the BS3 3-Clause license.

[rpaas-v1-repository]: https://github.com/tsuru/rpaas.git
[opensource-globocom]: https://opensource.globo.com
[nginx-site]: https://nginx.org/
[kubernetes-operator]: https://coreos.com/operators/
[nginx-operator-repository]: https://github.com/tsuru/nginx-operator.git
