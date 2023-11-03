# ipflare

![banner](etc/logo.svg)

## about

This small utility can help with operating a web service from behind a dynamic IP address - similar to dynamic DNS.

Internet service providers often dynamically allocate customer IP address and either charge extra for, or outright not allow a static IP. This is a problem if you want to host a service from your home, as the IP address can change any time resulting in the exposed endpoints no longer being available.

With `ipflare` you can quickly setup a small utility program that can run on the server host (or any host in the same network) and continuously check if the public IP address has changed. If it detects a change, it will update the DNS records associated with your domain to the new IP address.

## requirements

   * Your own domain
   * DNS routed through Cloudflare for your domain
   * Cloudflare API key token that can update DNS records for the configured zones

## usage

You can configure the utility through a configuration file, and environment variables.

### environment variables

| Environment variable | Description                   |
| -------------------- | ----------------------------- |
| CONFIG_PATH          | override the config file path |
| CLOUDFLARE_API_TOKEN | set/override the api key      |


### configuration file

The configuration file is loaded from `/etc/ipflare/config.yaml` by default.
An example configuration can be found in the `etc/config.yaml` file under the source root.

Example configuration:

```yaml
api_token: "xyz"          # the cloudflare api key token
frequency: 30             # check frequency in seconds
entries:
  example.com:            # cloudflare zone name
    - "foo.example.com"   # dns record to update (type "A" only)
    - "bar.example.com"
    - "baz.example.com"
  domain.org:
    - "domain.org"        # update root record
```

Notes:

   * the `CLOUDFLARE_API_TOKEN` environment variable overrides the one defined in the config file
   * the tool only supports type A DNS records
   * the root record can also be updated



## docker image

You can pull the latest docker image:

```
docker pull git.okki.hu/garric/ipflare
```

...or use the included Dockerfile to build the image yourself.

```
cd <source code root>
docker build -t ipflare .
```

In order to configure the when using the docker image, use volumes/mounts and environment variables:

```
docker run -e CLOUDFLARE_API_TOKEN="foobar" -v "$(pwd)/config.yaml:/etc/ipflare/config.yaml" git.okki.hu/ipflare
```

## running as a service

### using docker restart policy

If you are using the docker image, the recommended way is to use a docker restart policy (`--restart`):

```
docker run -d --restart always ...
```

To check the logs, use docker to list the container id, then use the container id to display the logs:

```
docker ps -a
docker logs <id>
```

### running the binary as a systemd service

In case you want to run `ipflare` directly as a systemd service, the following example shows how to create a systemd unit:

```
[Unit]
Description=ipflare service
After=network-online.target
Wants=network-online.target

[Service]
Restart=always
ExecStart=/path/to/ipflare

[Install]
WantedBy=multi-user.target

```
Make sure the config file is added to `/etc/ipflare/config.yaml`.
Copy the service file to `/etc/systemd/system`, then start and enable the systemd service:

```
sudo systemctl enable ipflare.service
sudo systemctl start ipflare.service
```

The application prints to STDOUT, so you can easily check the output as well:

```
sudo systemctl status ipflare.service
```

