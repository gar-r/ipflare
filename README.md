# ipflare

![banner](etc/banner.png)

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

You can pull the latest docker image from docker hub...

```
docker pull garric/ipflare
```

...or use the included Dockerfile to build the image yourself.

```
cd <source code root>
docker build -t garric/ipflare .
```

When running the app with docker, the command line arguments can be specified in the usual fashion.

Example:

```
docker run ipflare \
   -t "..." \
   -e "mysite/*.example.com"
```

```
ipflare starting with the following parameters:
auth token: [...]
 frequency: 30
   entries: [mysite/*.example.com]
```

## running as a system service

It is recommended to run `ipflare` as a system service, allowing it to run always without interruption. The steps vary from OS to OS. On GNU/Linux operating systems one way to achieve this is by creating a systemd service, preferably using [systemd-docker](https://github.com/ibuildthecloud/systemd-docker) to wrap the systemd unit. The following example shows a systemd unit which uses the app's docker image:

```
[Unit]
Description=ipflare service
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=120
Restart=always
ExecStartPre=/usr/bin/docker pull garric/ipflare
ExecStart=/usr/bin/docker run --rm --name %n garric/ipflare \
    -t "token" \
    -f 10 \
    -e "my-site/*.example.com"

[Install]
WantedBy=multi-user.target

```

To start and enable the above systemd service:

```
sudo cp etc/ipflare.service /etc/systemd/system/
sudo systemctl enable ipflare.service
sudo systemctl start ipflare.service
```

The application prints to STDOUT, so you can easily check the output as well:

```
sudo systemctl status ipflare.service
```

