# gogen

Monorepo and service generator for Golang projects.

## Install

### Linux / Mac

```
curl -fsSL https://raw.githubusercontent.com/fresh8gaming/gogen/trunk/install.sh | sh
```

Or

```
wget -q https://raw.githubusercontent.com/fresh8gaming/gogen/trunk/install.sh -O- | sh
```

## Usage

### Generate Monorepo

```sh
gogen repo /path/to/repo --team dmp --domain example
```

### Generate gRPC/HTTP Service

`gogen` does not distinguish between gRPC and HTTP services, but provides both as entry points to give you a false sense
of control. Both commands create a service with HPP and gRPC routing, utilising `grpc-gateway`. This way, HTTP services
can be documented with Swagger nicely.

```sh
gogen service grpc /path/to/repo --name neat-service
```

Or

```sh
gogen service http /path/to/repo --name neat-service
```

### Generate cron
`gogen` can also generate code a basic cron service with its deploment with:

```sh
gogen service cron /path/to/repo --name neat-service
```

### Generate crawler cron
`gogen` can also generate a feed crawler specific cron, where you can also specify weather this is an inplay crawler or not.

```sh
gogen service crawler-cron /path/to/repo --name neat-crawler
```

Or, for inplay.

```sh
gogen service crawler-cron /path/to/repo --name neat-crawler --inplay true
```
