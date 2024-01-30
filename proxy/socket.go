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
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			err = humane.Wrap(
				err,
				fmt.Sprintf("Could not accept a new connection on the %s:%d listener.", l.Proto, l.Port),
				"Make sure that your client is still connected to the Tailnet and re-authenticate if necessary.",
			)

			log.Warn().Err(err).Msg("Could not accept connection")
			continue
		}

		go l.handleConnection(ctx, srv, conn)
	}
}

func (l *Listener) handleConnection(ctx context.Context, srv *tsnet.Server, conn net.Conn) {
	defer conn.Close()

	log.Info().Str("remote", conn.RemoteAddr().String()).Msg("New connection")

	remote, err := srv.Dial(ctx, l.Proto, l.Target)
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

	go io.Copy(remote, conn)
	go io.Copy(conn, remote)

	<-ctx.Done()
}
