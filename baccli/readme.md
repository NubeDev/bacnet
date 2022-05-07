## whois

```
go run main.go discover --interface=wlp3s0
```

## whois

```
go run main.go whois --interface=wlp3s0
```

## read

```
go run main.go read --interface=wlp3s0 --device=202 --objectID=1 --objectType=1 --property=87
```

## write to an AO

```
go run main.go write --interface=wlp3s0 --device=202 --objectID=1 --objectType=1 --property=85 --priority=1 --null=true
```

## write null

```
go run main.go write --interface=wlp3s0 --device=202 --objectID=1 --objectType=1 --property=85 --priority=1 --null=true
```
