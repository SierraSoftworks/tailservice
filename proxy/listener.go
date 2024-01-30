package proxy

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	humane "github.com/sierrasoftworks/humane-errors-go"
	"tailscale.com/tsnet"
)

type Listener struct {
	Proto  string
	Port   int
	Secure bool

	Target string
}

// Takes in a string describing a listening port and the target service to forward
// to, parsing it into a Listener object.
//
// The format of the input string can be one of the following:
// - "port:target_port" - A listener on the specified port, forwarding to the target service on the local machine.
// - "port:target_host:target_port" - A listener on the specified port, forwarding to the target service on a remote machine.
func ParseListener(s, proto string, secure bool) (*Listener, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, humane.New(
			"The listener you provided did not match the expected format.",
			"To forward port 80 to a local port 8080, you can specify '80:8080'.",
			"To forward port 80 to a remote port 8080, you can specify '80:example.com:8080'.",
		)
	}

	port, err := strconv.ParseInt(parts[0], 10, 16)
	if err != nil {
		return nil, humane.Wrap(
			err,
			"The port you provided was not a valid number.",
			"Make sure you've provided the port number you want to listen on as the first part of the listener string (i.e. 80:8080).",
			"Make sure that the port number you specify is between 1-65384.",
		)
	}

	target := parts[1]

	if !strings.Contains(target, ":") {
		target = fmt.Sprintf("127.0.0.1:%s", target)
	}

	return &Listener{
		Proto:  proto,
		Port:   int(port),
		Secure: secure,
		Target: target,
	}, nil
}

func (l *Listener) Start(ctx context.Context, srv *tsnet.Server) error {
	var listener net.Listener
	var err error

	if l.Secure {
		listener, err = srv.ListenTLS(l.Proto, fmt.Sprintf(":%d", l.Port))
	} else {
		listener, err = srv.Listen(l.Proto, fmt.Sprintf(":%d", l.Port))
	}

	if err != nil {
		err = humane.Wrap(
			err,
			"Could not start listening on the provided port.",
			"Make sure that you only specify one listener for each port.",
			"Make sure that, if you're using the --tls flag, you've enabled HTTPS certificates on your Tailnet.",
		)
		log.Error().Err(err).Int("port", l.Port).Str("proto", l.Proto).Bool("secure", l.Secure).Msg("Could not start listener")
		return err
	}

	if strings.HasPrefix(l.Target, "http://") || strings.HasPrefix(l.Target, "https://") {
		log.Info().Int("port", l.Port).Str("proto", l.Proto).Bool("secure", l.Secure).Str("target", l.Target).Msgf("Forwarding traffic from %s:%d to %s:%s (using http proxy mode)", l.Proto, l.Port, l.Proto, l.Target)
		go l.listenHttp(ctx, srv, listener)
	} else {
		log.Info().Int("port", l.Port).Str("proto", l.Proto).Bool("secure", l.Secure).Str("target", l.Target).Msgf("Forwarding traffic from %s:%d to %s:%s", l.Proto, l.Port, l.Proto, l.Target)
		go l.listenSocket(ctx, srv, listener)
	}

	return nil
}
