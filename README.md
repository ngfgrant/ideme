# IDEme

IDEme is a quick and easy way to deploy a isolated Theia IDE instance to a
DigitalOcean account.

This is useful for a range of use cases from setting up your own cloud IDE
through to deploying temporary environments to share with others, such as in
interview settings.

# Prerequisites

- I have not built any binaries for this yet so you will need Go installed on
  in order to install this app.

- You will also need a DigitalOcean Account with and API token set in your
  environment as `DO_TOKEN`.

- A domain name that you want to run this application on. Freenom offer some
  free TLDs.

# Install

1. Pull the repository

2. `cd ideme`

3. `go install`

# Usage

1. Ensure you have a writable DigitalOcean API token set in your environment as
   `DO_TOKEN`. E.g. `export DO_TOKEN=token-without-quotes`

2. Check `config.yaml` and adjust the settings appropriately. For example, you
   will want to set the domain to one you own. You may also want to increase
   the power of the droplet used.

3. You can then deploy a Theia IDE like so:

```shell
$ ideme deploy infrastructure

$ ideme deploy app

$ ideme delete app <app-name>
```

The Droplet and IDE containers usually take 3-4 minutes to get up and running
so the application won't be available until then.

# Contributing

Please feel free to open up issues and fork. No specific guidelines for this
project yet.
