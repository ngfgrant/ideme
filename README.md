# IDEme

IDEme is a quick and easy way to deploy a isolated Theia IDE instance to a
DigitalOcean account.

This is useful for a range of use cases from setting up your own cloud IDE
through to deploying temporary environments to share with others, such as in
interview settings.

The application actually just runs any docker image and exposes this over an
Nginx proxy with Lets Encrypt certificates add on.

# Prerequisites

- I have not built any binaries for this yet so you will need Go installed on
  in order to install this app.

- You will also need a DigitalOcean Account with and API token set in your
  environment as `DO_TOKEN`.

- You will need a SSH public key set as an environment variable
  `IDEME_PUB_KEY`. This public key will be automatically added to your DigitalOcean
  account and can be used to login to your droplets.

- A domain name that you want to run this application on. Freenom offer some
  free TLDs.

# Install

1. Pull the repository

2. `cd ideme`

3. `go install`

# Usage

1. Ensure you have a writable DigitalOcean API token set in your environment as
   `DO_TOKEN`. E.g. `export DO_TOKEN=token-without-quotes`

2. Ensure you have a SSH Public key your environment as `IDEME_PUB_KEY`.

3. Check `config.yaml` and adjust the settings appropriately. For example, you
   will want to set the domain to one you own. You may also want to increase
   the power of the droplet used.

4. You can then deploy an app such as Theia IDE like so:

```shell
$ ideme deploy infrastructure # deploys vpc, project, ssh keys, firewall, domain, tags

$ ideme deploy app

$ ideme deploy app --image theiaide/theia-go

$ ideme delete app <app-name>
```

The original idea behind this app was to quickly spin up cloud IDEs for use in
an interview setting.

Theia IDE DockerHub gives a list of images that you can use and of course you
can build your own. By default, the app uses whichever image is specified in
`config.yaml` at:

```yaml
application:
  image: username/image:tag
```

This means that if you deploy an application with a `--image` flag the app will
fall back to the above default.

The Droplet and IDE containers usually take 3-4 minutes to get fully up and
running so the application won't be available until then.

# Contributing

Please feel free to open up issues and fork. No specific guidelines for this
project yet.
