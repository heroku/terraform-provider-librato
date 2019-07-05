Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Building The Provider
---------------------

```sh
git clone git@github.com:heroku/terraform-provider-librato
cd terraform-provider-librato
make build
```

Using the provider
----------------------

```sh
export $GOPATH=%{GOPATH:-$HOME/go/bin}
make install
cp $GOPATH/bin/terraform-provider-librato ~/.terraform.d/plugins
```

After recompiling the provider, you will need to re-run `terraform init` init any projects that use it.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is required).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-librato
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```
