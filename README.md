# env
A simple environment variable to struct parsing library

## How to use it
```go
type configs struct {
    Host string `env:"HOST" default:"localhost"`
    Port string `env:"PORT" default:"8080"`
    IsProdReady bool `env:"PROD"`
    MaxConnection int `env:"MAXCONN"`
}

func main() {
    cfgs := configs{}

    if err := env.Parse(&cfgs); err != nil {
        fmt.Println(err)
    }

    fmt.Println(cfgs)
}
```