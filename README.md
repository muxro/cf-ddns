# cf-ddns
Dynamic DNS Client for [CloudFlare](https://cloudflare.com).

## Getting started

### Instalation
You will need a [Go](http://golang.org) development environment to build this program. 

After the environment is properly set up, you can build and install the programin `$GOPATH/bin`
```sh
go get github.com/muxro/cf-ddns
```

### Configuration
Configuration is done trough a [YAML](http://yaml.org/) file, which should have the following format:

```yaml
checkInterval: 60

zones:
    example.com:
        email: 'you@example.com'
        token: 'CLOUDFLARE TOKEN'
        domains:
            - 'example.com'
            - 'www.example.com'
            - '*.dev.example.com'
```
`checkInterval` is the interval, in seconds, at which to check for IP address changes.

`zones` define the zones that you need to update, each zone has to have an email, Cloudflare API token, and the domains that need to be updated.

### Running
Assuming that the binary is in your `$PATH`, run the following command:
```sh
cf-ddns -config /path/to/config
```
This will start the program in the foreground, with logging to stdout. To run this program as a daemon, use the following command.
```sh
nohup cf-ddns -config /path/to/config >cf-ddns.log &
```

### License
Copyright (C) 2014 [Mux Systems SRL-D](http://mux.ro)

Shared under the GPLv3 license. See the LICENSE file for more details.
