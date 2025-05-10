# net-scan

**net-scan** is a modular CLI tool for performing network diagnostics using a centralized host file. It supports port scanning, ICMP pings, HTTP availability checks, DNS queries, and host list management. Designed for developers, sysadmins, and automation workflows.

## Features

- Port scanning with support for custom ports, ranges, and protocols (TCP, UDP, etc.)
- ICMP ping with support for raw sockets, TTL, traffic class, and custom intervals
- HTTP health checks with configurable frequency, timeout, and HTTPS support
- DNS record lookup (CNAME, A, AAAA, NS, MX, TXT)
- Host management commands for storing frequently used targets
- Configuration via YAML or CLI flags using [Viper](https://github.com/spf13/viper)
