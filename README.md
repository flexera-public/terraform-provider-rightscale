Release 1.0
==================

This is version 1.0 of the Terraform provider for RightScale.  Acceptance into the terraform community and as an official provider is now in progress, but as of this version all tests, docs, and sufficient resources for full operational coverage is complete and tested. [ Note that instructions below are forward-looking for where this repo will move in the near future. ]

Markdown (Documentation) is available here:
- [Resources](https://github.com/rightscale/terraform-provider-rightscale/tree/master/website/docs/r)
- [Datasources](https://github.com/rightscale/terraform-provider-rightscale/tree/master/website/docs/d)

Please [open an issue](https://github.com/rightscale/terraform-provider-rightscale/issues/new) if you find a bug or otherwise are interested in contributing to this open source effort.  PRs accepted! 

Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 0.10.8
- [Go](https://golang.org/doc/install) >= 1.8 (to build the provider plugin)
- [Dep](https://github.com/golang/dep) >= 0.3.2 (Handle dependencies)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-rightscale`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-rightscale
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-rightscale
$ make build
```

Using the provider
----------------------

See the [RightScale Provider documentation](https://www.terraform.io/docs/providers/rightscale/index.html) to get started using the RightScale provider.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-rightscale
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
