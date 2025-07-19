<h1 align="center">
  <img src="Meta.png" alt="Meta Kennel" width="200">
  <br>Meta Kernel<br>
</h1>

<h3 align="center">Another Mihomo Kernel.</h3>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/MetaCubeX/mihomo">
    <img src="https://goreportcard.com/badge/github.com/MetaCubeX/mihomo?style=flat-square">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/MetaCubeX/mihomo/Alpha?style=flat-square">
  <a href="https://github.com/MetaCubeX/mihomo/releases">
    <img src="https://img.shields.io/github/release/MetaCubeX/mihomo/all.svg?style=flat-square">
  </a>
  <a href="https://github.com/MetaCubeX/mihomo">
    <img src="https://img.shields.io/badge/release-Meta-00b4f0?style=flat-square">
  </a>
</p>

## Features

- Local HTTP/HTTPS/SOCKS server with authentication support
- VMess, VLESS, Shadowsocks, Trojan, Snell, TUIC, Hysteria protocol support
- Built-in DNS server that aims to minimize DNS pollution attack impact, supports DoH/DoT upstream and fake IP.
- Rules based off domains, GEOIP, IPCIDR or Process to forward packets to different nodes
- Remote groups allow users to implement powerful rules. Supports automatic fallback, load balancing or auto select node
  based off latency
- Remote providers, allowing users to get node lists remotely instead of hard-coding in config
- Netfilter TCP redirecting. Deploy Mihomo on your Internet gateway with `iptables`.
- Comprehensive HTTP RESTful API controller

## Dashboard

A web dashboard with first-class support for this project has been created; it can be checked out at [metacubexd](https://github.com/MetaCubeX/metacubexd).

## Configration example

Configuration example is located at [/docs/config.yaml](https://github.com/MetaCubeX/mihomo/blob/Alpha/docs/config.yaml).

## Docs

Documentation can be found in [mihomo Docs](https://wiki.metacubex.one/).

## For development

Requirements:
[Go 1.20 or newer](https://go.dev/dl/)

Build mihomo:

```shell
git clone https://github.com/MetaCubeX/mihomo.git
cd mihomo && go mod download
go build
```

Set go proxy if a connection to GitHub is not possible:

```shell
go env -w GOPROXY=https://goproxy.io,direct
```

Build with gvisor tun stack:

```shell
go build -tags with_gvisor
```

### IPTABLES configuration

Work on Linux OS which supported `iptables`

```yaml
# Enable the TPROXY listener
tproxy-port: 9898

iptables:
  enable: true # default is false
  inbound-interface: eth0 # detect the inbound interface, default is 'lo'
```

## Debugging

Check [wiki](https://wiki.metacubex.one/api/#debug) to get an instruction on using debug
API.

## API: Content Check

A new API endpoint `/content_check` has been added to check the content of a URL through all available proxies.

**Endpoint:** `GET /content_check`

**Query Parameters:**

*   `url` (required): The URL to check. e.g., `https://checkip.amazonaws.com/`
*   `timeout` (optional): Timeout for each request in duration format. Defaults to `5s`.
*   `retries` (optional): The number of retries for a failed request. Defaults to `0`.

**Example Usage:**

```bash
curl "http://127.0.0.1:9093/content_check?url=https%3A%2F%2Fcheckip.amazonaws.com%2F&timeout=3s"
```

**Response Body:**

The API returns a JSON object containing a list of results, one for each proxy.

```json
{
  "results": [
    {
      "proxy": "ProxyName1",
      "url": "https://checkip.amazonaws.com/",
      "content": "123.45.67.89",
      "error": ""
    },
    {
      "proxy": "ProxyName2",
      "url": "https://checkip.amazonaws.com/",
      "content": "",
      "error": "request timed out"
    }
  ]
}
```

**Field Description:**

*   `proxy`: The name of the proxy being tested.
*   `url`: The target URL that was requested.
*   `content`: 
    *   **On success:** The response body from the URL as a string. Any leading/trailing whitespace (including newlines) is removed. For `https://checkip.amazonaws.com/`, this will be the clean IP address.
    *   **On failure:** An empty string `""`.
*   `error`:
    *   **On success:** An empty string `""`.
    *   **On failure:** A string containing the error message, e.g., `"request timed out"`, `"connection refused"`, etc.

## Credits

- [Dreamacro/clash](https://github.com/Dreamacro/clash)
- [SagerNet/sing-box](https://github.com/SagerNet/sing-box)
- [riobard/go-shadowsocks2](https://github.com/riobard/go-shadowsocks2)
- [v2ray/v2ray-core](https://github.com/v2ray/v2ray-core)
- [WireGuard/wireguard-go](https://github.com/WireGuard/wireguard-go)
- [yaling888/clash-plus-pro](https://github.com/yaling888/clash)

## License

This software is released under the GPL-3.0 license.

**In addition, any downstream projects not affiliated with `MetaCubeX` shall not contain the word `mihomo` in their names.**