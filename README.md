# ipflare

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
        "-f", "15",
        "-e", "website1/mail.example.com",
        "-e", "website1/*.example.com",
        "-e", "website2/conference.chat.com",
```

The above command starts the application with a `15` second checking frequency, and the token `token` (for each request to CloudFlare this will be added to the Authentication request headers).

The command adds the following DNS entries:

   * zone: `website1`, records: `mail.example.com`, `*.example.com`
   * zone: `website2`, records: `conference.chat.com`


## docker image

TODO