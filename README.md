# ipflare

![banner](etc/logo.svg)

## about

This small utility can help with operating a web service from behind a dynamic IP address - similar to dynamic DNS.

Internet service providers often dynamically allocate customer IP address and either charge extra for, or outright not allow a static IP. This is a problem if you want to host a service from your home, as the IP address can change any time resulting in the exposed endpoints no longer being available.

With `ipflare` you can quickly setup a small utility program that can run on the server host (or any host in the same network) and continuously check if the public IP address has changed. If it detects a change, it will update the DNS records associated with your domain to the new IP address.

## requirements

   * Your own domain
   * DNS routed through CloudFlare for your domain
   * CloudFlare API key

## usage

You can configure the utility through command line arguments. To get help for the usage you can use the `-h` flag:

```
ipflare -h
```

The command line flags are listed in the following table:

| Flag | Description                                            | Mandatory | Default |
| ---- | ------------------------------------------------------ | --------- | ------- |
| -f   | frequency of ip change detection in seconds            | no        | 30s     |
| -t   | your CloudFlare api token                              | yes       |         |
| -e   | DNS entry to keep up-to-date in the zone/record format | no        |         |


Notes:

   * The `-e` flag can be added multiple times to update multiple records. Regardless of the order in which you specify the entries, requests to CloudFlare are automatically grouped by the zones to reduce the amount of requests required.
   * It is recommended to adjust the defult update frequency (`-f`) while keeping the CloudFlare api rate limiting and the amount of entries to update in mind.


## example usage

```
ipflare "-t", "token" \
        "-f", "15" \
        "-e", "website1/mail.example.com" \
        "-e", "website1/*.example.com" \
        "-e", "website2/conference.chat.com"
```

The above command starts the application with a `15` second checking frequency, and the token `token` (for each request to CloudFlare this will be added to the Authentication request headers).

The command adds the following DNS entries:

   * zone: `website1`, records: `mail.example.com`, `*.example.com`
   * zone: `website2`, records: `conference.chat.com`


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

When running the app with docker, the command line arguments can be specified in the usual fashion.

Example:

```
docker run ipflare \
   -t "..." \
   -e "mysite/*.example.com"
```

Example output:

```
ipflare starting with the following parameters:
auth token: [...]
 frequency: 30
   entries: [mysite/*.example.com]
```

## running as a service

### using docker restart policy

If you are using the docker image, the recommended way is to use a docker restart policy (`--restart`):

```
docker run -d --restart always --name ipflare git.okki.hu/garric/ipflare \
   -t "..." \
   -e "mysite/*.example.com"
```

To checking the logs, use docker to list the container id, then use the container id to display the logs:

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
ExecStart=/path/to/ipflare \
    -t "token" \
    -f 10 \
    -e "my-site/*.example.com"

[Install]
WantedBy=multi-user.target

```

Copy the service file to `/etc/systemd/system`, then start and enable the systemd service:

```
sudo systemctl enable ipflare.service
sudo systemctl start ipflare.service
```

The application prints to STDOUT, so you can easily check the output as well:

```
sudo systemctl status ipflare.service
```

