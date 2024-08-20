# Tail Service
**Easily expose services on your [Tailscale](https://tailscale.com/) network.**

This project provides an extremely easy to use proxy for exposing services on your
[Tailscale](https://tailscale.com/) network under their own MagicDNS name. It is
particularly useful if you're hosting several services on the same machine and
want to access them by name instead of needing to remember port numbers.

## Installation

```bash
go install github.com/sierrasoftworks/tailservice@latest
```

## Usage
At its simplest, you can expose a service on your local machine by running the
`tailservice` command with a `--name` and one or more `--tcp`, `--udp`, or `--tls`
arguments specifying the ports to expose.

```bash
# Expose port 80 on the local machine as my-service on your tailnet,
# listening on ports 80 and 443 (port 443 will get a TLS certificate
# automatically).
tailservice --name my-service --tls 443:80 --tcp 80:80
```

### Exposing Ports
The `tailservice` command supports exposing ports using TCP and UDP
network protocols. It also supports automatically generating TLS certificates
for your services using [Let's Encrypt](https://letsencrypt.org/) (if you
have configured your Tailscale account to support HTTPS certificates).

When specifying a port to expose, you first indicate the type of protocol
you'd like to receive traffic on (e.g. `--tcp`, `--udp` or the special `--tls`
variant), followed by the listener specification.

```bash
# Forwards raw TCP traffic from port 80 on the Tailnet service
# to port 8080 on the local machine.
tailservice --name my-service --tcp 80:8080

# Forwards raw UDP traffic from port 53 on the Tailnet service
# to port 53 on a remote machine.
tailservice --name my-service --udp 53:8.8.4.4:53

# Forwards TLS traffic from port 443 on the Tailnet service
# to port 8080 on another Tailnet node.
tailservice --name my-service --tls 443:example-node.tails-scales.ts.net:8080
```

### Exposing services using Funnel
Tailscale's Funnel functionality allows you to expose tailnet services to the public
internet, allowing clients without Tailscale installed to access the service. To enable
this functionality, you can use the `--funnel` flag when starting `tailservice` and
configure a listener on port `443`, `8443` or `10000`.

```bash
# Forwards TLS traffic from port 443 on your funnel endpoint to
# port 8080 on the local machine.
tailservice --name my-service --tls 443:8080 --funnel
```

### Running in Ephemeral Mode
By default, `tailservice` will save its configuration to disk so that it can
be restarted without the need to re-authenticate. Running in this manner retains
the IP address of the service on your Tailnet, allowing you to use the same
DNS name to access it regardless of how fresh your DNS cache is.

If you'd prefer that the service is removed from your Tailnet when it is
stopped, you can use the `--ephemeral` flag to run in ephemeral mode. This
mode is particularly useful if you're running `tailservice` in a container
or for test purposes.

```bash
tailservice --name my-service --tcp 80:80 --ephemeral
```

### Specifying a Tailscale Authkey
If you're running `tailservice` in a container or on a headless machine,
you may find it useful to specify the Tailscale Authkey using an environment
variable. Doing so is only necessary on the first run, as the resulting config
will be saved to disk (note that this does not apply if `--ephemeral` is used
or if the config file is deleted).

```bash
# Configure your Tailscale authentication key
export TS_AUTHKEY="tskey-1234567890abcdef"
```

### Enabling Tailscale Debug Logging
If you're having trouble getting `tailservice` to work, you can enable debug
logging by passing the `--ts-debug` flag. This will cause `tailservice` to
print out the raw Tailscale logs to the console, which may help you to
diagnose the problem.

```bash
tailservice --name my-service --tcp 80:80 --ts-debug
```
