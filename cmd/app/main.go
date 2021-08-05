package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		hystrix.Do("global", func() error {
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				return err
			}
			for _, addr := range addrs {
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					return err
				}

				if ip.IsLoopback() {
					continue
				}

				_, err = fmt.Fprintln(rw, ip.String())
				return err
			}
			return nil
		}, func(err error) error {
			return nil
		})
	})
	s := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
	s.ListenAndServe()
}

var (
	hystrixConfig *hystrix.CommandConfig
)

func init() {
	hystrixConfig = &hystrix.CommandConfig{
		Timeout:               10,
		MaxConcurrentRequests: 3,
		ErrorPercentThreshold: 25,
	}
	hystrix.ConfigureCommand("global", *hystrixConfig)
	hystrix.SetLogger(log.Default())
}

func Limit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hystrix.Do("global", func() error {
			fmt.Println("...")
			h.ServeHTTP(w, r)
			return nil
		}, func(err error) error {
			fmt.Println("... ", err)
			return nil
		})
	})
}
