package myleaf

import (
	"os"
	"os/signal"

	"github.com/sharehdm/myleaf/cluster"

	"github.com/sharehdm/myleaf/log"
	"github.com/sharehdm/myleaf/module"
)

func Run(mods ...module.Module) {
	log.Release("Leaf %v starting up", version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()
	// cluster
	cluster.Init()
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("Leaf closing down (signal: %v)", sig)
	cluster.Destroy()
	module.Destroy()
}
