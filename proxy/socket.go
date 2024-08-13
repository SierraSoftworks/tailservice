package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/rs/zerolog/log"
	humane "github.com/sierrasoftworks/humane-errors-go"
	"tailscale.com/tsnet"
)

func (l *Listener) listenSocket(ctx context.Context, srv *tsnet.Server, listener net.Listener) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			herr := humane.Wrap(
				err,
				fmt.Sprintf("Could not accept a new connection on the %s:%d listener.", l.Proto, l.Port),
				"Make sure that your client is still connected to the Tailnet and re-authenticate if necessary.",
			)

			log.Warn().Err(err).Msg(herr.Display())
			continue
		}

		go l.handleConnection(ctx, srv, conn)
	}
}

func (l *Listener) handleConnection(ctx context.Context, srv *tsnet.Server, conn net.Conn) {
	defer conn.Close()

	log.Debug().Str("remote", conn.RemoteAddr().String()).Msg("New connection")

	var remote net.Conn
	var err error

	if isTailscaleHost(l.Target) {
		remote, err = srv.Dial(ctx, l.Proto, l.Target)
	} else {
		remote, err = net.Dial(l.Proto, l.Target)
	}

	if err != nil {
		err = humane.Wrap(
			err,
			fmt.Sprintf("Could not establish a connection to the target service %s:%s.", l.Proto, l.Target),
			"Make sure that the target service is still running and reachable.",
			"Make sure that your client is still connected to the Tailnet and re-authenticate if necessary.",
			"Make sure that you've specified the correct target service in your listener configuration.",
		)

		return
	}

	defer remote.Close()

	close := make(chan struct{})

	go func() {
		io.Copy(remote, conn)
		close <- struct{}{}
	}()
	go func() {
		io.Copy(conn, remote)
		close <- struct{}{}
	}()

	select {
	case <-ctx.Done():
	case <-close:
		if addr := remote.RemoteAddr(); addr != nil {
			log.Debug().Str("remote", addr.String()).Msg("Connection closed")
		} else {
			log.Debug().Msg("Connection closed")
		}
	}
}
