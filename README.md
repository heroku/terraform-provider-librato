# Librato Terraform Provider (v1)

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

## Building The Provider

```sh
git clone git@github.com:heroku/terraform-provider-librato
cd terraform-provider-librato
make build
```

## Using the provider

```sh
export $GOPATH=%{GOPATH:-$HOME/go/bin}
make build
cp $GOPATH/bin/terraform-provider-librato ~/.terraform.d/plugins
```

<<<<<<< HEAD
Using the provider
----------------------
for local use:
- configure provider in `main.tf`
```
provider "librato" {
  version = "0.1.1"
  email = "${var.username}"
  token = "${var.token}"
}
```
- build provider binary and update plugin
```
mkdir -p ~/.terraform.d/plugins
cd ~/go/src/github.com/notmaxx/terraform-provider-librato
make
mv $GOPATH/bin/terraform-provider-librato ~/.terraform.d/plugins/terraform-provider-librato_v0.1.1
...
cd TERRAFORM_PROJECT_DIR
terraform init    # make sure new provider is initialized
terraform apply
```
After recompiling the provider, you will need to re-run `terraform init` init any projects that use it.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is required).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-librato
...
```

### Running tests

#### Unit tests

```sh
$ make test
```

#### Acceptance tests

Acceptance tests send requests against the Librato API. Hit the [Librato UI](https://metrics.librato.com/tokens) and create a token. The Email should be popiulated for you automatically and the token will be generated. You will want a Full Access token as you'll be creating (and tearing down) Librato resources.

```sh
export LIBRATO_EMAIL=<email>
export LIBRATO_TOKEN=<token>

make testacc
```

If you want to run a single acceptance test:

```sh
export LIBRATO_EMAIL=<email>
export LIBRATO_TOKEN=<token>

# example: TestAccLibratoAlert_Rename
RUN=TestAccLibratoAlert_Rename make testaccone
# or
make RUN=TestAccLibratoAlert_Rename testaccone
```
