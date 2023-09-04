<img alt="pancap logo" src="pancap.png" width="250px" height="250px">

# pancap

## Idea

If you get access to a [PCAP](https://en.wikipedia.org/wiki/Pcap) file, for example during a CTF or captured on your own, you usually have the problem of overlooking all the relevant information to get a basic idea of the capture file. This gets worse if the capture file includes lots of white noise or irrelevant traffic - often included in the capture file to cloak *interesting* packets in a bunch of packets to YouTube, Reddit, Twitter and others.

*pancap* addresses this problem. With multiple submodules, it analyzes the given PCAP file and extracts useful information out of it. In many cases, this saves you a lot of time and can point you into the right direction.

## Features

- Support for different network protocols
	- ARP: collect communication, identify switches and routers, detect [ARP spoofing](https://www.crowdstrike.com/cybersecurity-101/spoofing-attacks/arp-spoofing/).
	- DHCP: analyze requests and responses, get an idea of the network setup
	- DNS: collect hints of user actions and their OS
	- HTTP: dump cleartext communication and embedded files
- Create [GraphViz](https://graphviz.org/) graphs out of network communication flow

## Usage

Simply run

`go install github.com/maride/pancap@latest`

This will build `pancap` and place it into your `GOBIN` directory - means you can directly execute it!
It might be required to install the `pcap` header files, e.g. for Ubuntu with `apt install libpcap-dev`.

In any use case, you need to specify the file you want to analyze, simply handed over to pancap with the `-file` flag.

`pancap -file ~/Schreibtisch/mitschnitt.pcapng`

## Benchmarks

Parsing an `n`GB big pcap takes `y` seconds:

| `n`GB | `y` seconds |
| ----- | ----------- |
| 2     | 30          |

## Contributions

... yes please! There are still a lot of modules missing.
If you are brave enough, you can even implement another Link Type. Pancap currently only supports `Ethernet` (which, to be honest, fits most cases well), but `USB` might be interesting, too. Especially sniffed keyboard and mouse packets are hard to analyze by hand...
