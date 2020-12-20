# go-infra-kit

This repository consists of multiple go packages to help with the infrastructure of your application.

## config

---

config is a lightweight Golang package for integrating both YAML and Environment variable into one config object.

It is primarily used to read config from YAML file during development and use Environment variables during production.

### Install

```shell
go get github.com/aamuz/go-infra-kit/config
```
### Usage

Define a configuration structure:
```go
type Config struct {
	Server struct {
		Port string `yaml: "port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Server   string `yaml:"server" envconfig:"DB_SERVER"`
		Port     int    `yaml:"port" envconfig:"DB_PORT" default:"1433"`
		User     string `yaml:"user" envconfig:"DB_USER"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
		Database string `yaml:"database" envconfig:"DB_DATABASE"`
	} `yaml:"database"`
}
```

Create a YAML file according to the configuration:
```yaml
server:
  port: 8080

database:
  server: "localhost"
  port: 1433
  user: "user"
  password: "password"
  database: "database"

```

Import the package
```go
import "github.com/aamuz/go-infra-kit/config"

cfg := infra.Config{}
if err := infra.Read("config.yml", &cfg); err != nil {
    logger.Fatalf("Error reading config: %v", err)
}
```

Now we can use the merged config from yaml and environment file as one.

```go
logger.Printf("server listening at port %s", cfg.Server.Port)
logger.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
```


## retry

---

retry provides a simple retry function to add retry policies in your application. 

This is taken from the blog post by Nick Stogner [Simple Golang Retry Function](https://upgear.io/blog/simple-golang-retry-function/).

### Install

```shell
go get -u github.com/aamuz/go-infra-kit/retry
```

### Usage

In the following 
```go

import "github.com/aamuz/go-infra-kit/retry"

...

err := retry.Retry(5, time.Second, func() error {
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        // This error will result in a retry
        return err
    }
    defer resp.Body.Close()

	s := resp.StatusCode
        switch {
        case s >= 500:
            // Retry
            return fmt.Errorf("server error: %v", s)
        case s >= 400:
            // Don't retry, it was client's fault
            return stop{fmt.Errorf("client error: %v", s)}
        default:
            // All is good no need to retry
            return nil
    }
})
```