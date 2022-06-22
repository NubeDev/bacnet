# bacnet

bacnet client written in go.

# Installation

learn some `golang` if you don't know it first :)

```
go mod tidy
```

For usage run:

```
go mod tidy
cd cli
go run main.go --help
```

## examples

change out the interface and device and so on

### whois

```
go run main.go whois --interface=wlp3s0
```

```
go run main.go whois --interface=wlp3s0
```

### read AO

```
go run main.go read --interface=wlp3s0 --device=202 --objectID=1 --objectType=1 --property=87
```

### write to an AO

```
go run main.go write --interface=wlp3s0 --device=202 --address=192.168.15.20 --network=4 --mstp=1 --objectID=1 --objectType=1 --property=85 --priority=16 --value=21
```

### write null to @16

```
go run main.go write --interface=wlp3s0 --device=202 --objectID=1 --objectType=1 --property=85 --priority=16 --null=true
```

## over a bacnet to ms-tp network

- router ip: 192.168.15.20
- bacnet router network number: 4
- bacnet mstp(rs485) mac address (between 0-255): 1

```
go run main.go read --interface=wlp3s0 --device=202 --address=192.168.15.20 --network=4 --mstp=1 --objectID=1 --objectType=1 --property=85 
```

get device name

```
go run main.go read --interface=wlp3s0 --device=202 --address=192.168.15.20 --network=4 --mstp=1 --objectID=202 --objectType=8 --property=77
```

### Max APDU Length

`Max APDU Length is important on for read/write prop multiple`

In the variable "Max APDU Length Accepted" the following are the values that can be returned:

```
mstp device: 480
ip device: 1476
```

### example same device getting the max APDU

get device MaxApdu length over MSTP will return `480`

```
go run main.go read --interface=wlp3s0 --device=202 --address=192.168.15.20 --network=4 --mstp=1 --objectID=202 --objectType=8 --property=62
```

get device MaxApdu length and on the same device but over IP will return `1476`

```
go run main.go read --interface=wlp3s0 --device=202 --address=192.168.15.202 --network=0 --mstp=0 --objectID=202 --objectType=8 --property=62
```

## Library

- [x] Who Is
- [x] Iam
- [x] Read Property
- [x] Read Multiple Property (beta)
- [ ] Read Range
- [x] Write Property
- [x] Write Property Multiple (beta)
- [ ] Who Has
- [x] What Is Network Number (beta)
- [x] Who Is Router To Network (beta)
- [ ] Change of Value Notification
- [ ] Event Notification
- [ ] Subscribe Change of Value
- [ ] Atomic Read File
- [ ] Atomic Write File

## Command Line Interface

- [x] Who Is
- [x] Iam
- [x] Read Property
- [x] Read Multiple Property
- [ ] Read Range
- [ ] Write Property
- [ ] Write Property Multiple
- [ ] Who Has
- [ ] What Is Network Number
- [x] Who Is Router To Network
- [ ] Atomic Read File
- [ ] Atomic Write File

# testing

## Tested on devices

- [x] Johnson Controls (FEC)
- [x] Easy-IO 30p, tested over IP and ms-tp
- [ ] Delta Controls
- [x] Reliable Controls
- [x] Honeywell Spyder
- [ ] Niagara N4 jace
- [ ] Schneider

## tested with other bacnet-libs

- [x] bacnet-stack
- [ ] bacnet-4j
- [x] bacpypes

This library is heavily based on the BACnet-Stack library originally written by Steve Karg.

- Ported and all credit to alex from https://github.com/alexbeltran/gobacnet
- And ideas from https://github.com/noahtkeller/go-bacnet



example whois
```go
	bytes := []byte{
		0x81, 0x0b, 0x00, 0x08, // BVLC
		0x01, 0x00, // NPDU
		0x10, 0x08, // APDU
	}

	pc, err := net.ListenPacket("udp4", ":47809")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:47808")
	if err != nil {
		panic(err)
	}

	_, err = pc.WriteTo(bytes, addr)
	if err != nil {
		panic(err)
	}
	d := make([]byte, 1)

	a, b, c := pc.ReadFrom(d)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
  ```


