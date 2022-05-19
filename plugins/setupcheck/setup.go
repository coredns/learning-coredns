package setupcheck

import (
	"log"

	"github.com/coredns/coredns/plugin"

	"github.com/coredns/caddy"
)

func init() {
	caddy.RegisterPlugin("setupcheck", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	log.Printf("START : setup")
	err := parse(c)
	if err != nil {
		return plugin.Error("setupcheck", err)
	}
	log.Printf("FINISH: setup")
	return nil
}

func parse(c *caddy.Controller) error {
	log.Printf("START : parse")
	for c.Next() {
		args := c.RemainingArgs()
		log.Printf("START : parse/setupcheck args %v", args)
		for c.NextBlock() {
			log.Printf("          parse/setupcheck block %s, args %v", c.Val(), c.RemainingArgs())
		}
		log.Printf("FINISH: parse/setupcheck args %v", args)
	}
	log.Printf("FINISH: parse")
	return nil
}
