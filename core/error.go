package core

import (
	"log"
	"net"
	"runtime/debug"
)

func LoadErrorHandler() func(c net.Conn) {
	return func(c net.Conn) {
		if r := recover(); r != nil {
			log.Printf("%s: %s\n", r, debug.Stack())
		}

		c.Close()
	}
}