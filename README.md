# bacnet

bacnet is a client for bacnet written in go.

# Installation

For usage run:

```
go mod tidy
cd baccli
go run baccli --help
```

## examples

change out the interface and device and so on

### whois

```
go run main.go whois --interface=wlp3s0
```

### read AO

```
go run main.go read --interface=wlp3s0 --device=202 --objectID=1 --objectType=1 --property=87
```

### write to an AO

```
go run main.go write --interface=wlp3s0 --device=202 --objectID=1 --objectType=1 --property=85 --priority=1 --null=true
```

### write null to @16

```
go run main.go write --interface=wlp3s0 --device=202 --objectID=1 --objectType=1 --property=85 --priority=1 --null=true
```

## over a bacnet to ms-tp network

- router ip: 192.168.15.20
- bacnet router network number: 4
- bacnet mstp(rs485) mac address (between 0-255): 1

```
go run main.go read --interface=wlp3s0 --device=202 --address=192.168.15.20 --network=4 --mstp=1 --objectID=1 --objectType=1 --property=85 
```

## Library

- [x] Who Is
- [x] Read Property
- [x] Read Multiple Property (beta)
- [ ] Read Range
- [x] Write Property
- [x] Write Property Multiple (beta)
- [ ] Who Has
- [ ] Change of Value Notification
- [ ] Event Notification
- [ ] Subscribe Change of Value
- [ ] Atomic Read File
- [ ] Atomic Write File

## Command Line Interface

- [x] Who Is
- [x] Read Property
- [x] Read Multiple Property
- [ ] Read Range
- [ ] Write Property
- [ ] Write Property Multiple
- [ ] Who Has
- [ ] Atomic Read File
- [ ] Atomic Write File

# testing

## Tested on devices

- [x] Johnson Controls (FEC)
- [x] Easy-IO 30p, tested over IP and mst-tp
- [] Delta
- [] Reliable Controls
- [x] Honeywell Spyder
- [] Schneider

## tested with other bacnet-libs

- [x] bacnet-stack
- [] bacnet-4j
- [x] bacpypes

This library is heavily based on the BACnet-Stack library originally written by Steve Karg.

- Ported and all credit to alex from https://github.com/alexbeltran/gobacnet
- And ideas from https://github.com/noahtkeller/go-bacnet
