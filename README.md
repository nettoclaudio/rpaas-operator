# RPaaS v2

[![Build Status](https://travis-ci.org/tsuru/rpaas-operator.svg?branch=master)](https://travis-ci.org/tsuru/rpaas-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/tsuru/rpaas-operator)](https://goreportcard.com/report/github.com/tsuru/rpaas-operator)

NOTE: This project is the replacement of the [RPaaS][rpaas-v1-repository].
Hence, we'll refer to it only as RPaaS v2 (although there are no breaking
changes between them).

---

## About

RPaaS, which stands for Reverse Proxy as a Service, provides a easy and fast
way to manage [NGINX][nginx-site]-based reverse proxy into a cloud
infrastructure. Such reverse proxy (aka RPaaS instance) handles incoming HTTP
request and forward to the configured destination application (backend).

Futhermore, it supports adding TLS terminating, cache (according to HTTP cache
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
