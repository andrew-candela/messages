package messages

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

const (
	WEBSOCKET_PROTOCOL = "UDPMWS"
	SIG_TOKEN_HEADER   = "RSA_SIG_TOKEN"
	SIG_HEADER         = "RSA_SIG"
	IP_TARGET_HEADER   = "IP_TARGET"
)

type MessageServer struct {
	logf          func(f string, v ...interface{})
	serveMux      http.ServeMux
	subscribersMu sync.Mutex
	subscribers   map[string]*subscriber
	config        map[string]DeliverConfig
}

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func newMessageServer(config map[string]DeliverConfig) *MessageServer {
	ms := &MessageServer{
		logf:        log.Printf,
		subscribers: make(map[string]*subscriber),
		config:      config,
	}
	ms.serveMux.HandleFunc("/config", ms.authenticateRequest(ms.handleConfig))
	// Dont need to authenticate this endpoint.
	// It's cheap to compute and if it is behind my auth system, you
	// can't use it if you don't already know your IP.
	ms.serveMux.HandleFunc("/ip", ms.handleIP)
	ms.serveMux.HandleFunc("/publish", ms.authenticateRequest(ms.handlePublish))
	ms.serveMux.HandleFunc("/subscribe", ms.authenticateRequest(ms.handleSubscribe))
	return ms
}

func (s *MessageServer) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Config:\n%v\n", s.config)
}

func (s *MessageServer) authenticateRequest(endpoint func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := s.config[StripPort(r.RemoteAddr)]
		// Couldn't find the requesting IP address in the config. Exiting.
		if !ok {
			fmt.Printf("Couldn't find IP: %v in config\n", r.RemoteAddr)
			http.Error(w, "Origination address mismatch", http.StatusForbidden)
			return
		}
		rsa_token_str := r.Header.Get(SIG_TOKEN_HEADER)
		rsa_token, err := StringToBytes(rsa_token_str)
		if err != nil {
			fmt.Printf("Couldn't unmarshal signature token passed in header from: %v\n", r.RemoteAddr)
			http.Error(w, "unverifiable signature", http.StatusForbidden)
			return
		}
		rsa_sig_str := r.Header.Get(SIG_HEADER)
		rsa_sig, err := StringToBytes(rsa_sig_str)
		if err != nil {
			fmt.Printf("Couldn't unmarshal signature string passed in header from: %v\n", r.RemoteAddr)
			http.Error(w, "unverifiable signature", http.StatusForbidden)
			return
		}
		pub_key := ParsePublicKey(user.Key)
		verified := RSAVerify(pub_key, []byte(rsa_token), []byte(rsa_sig))
		if !verified {
			fmt.Printf("Unable to verify request from IP: %v\n", r.RemoteAddr)
			http.Error(w, "signature mismatch", http.StatusForbidden)
			return
		}
		endpoint(w, r)
	}
}

func (s *MessageServer) handleIP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IP: %v\n", r.RemoteAddr)
}

func (s *MessageServer) handlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	ip_target := r.Header.Get(IP_TARGET_HEADER)
	if ip_target == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 1024)
	msg, err := io.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}
	err = s.publish(msg, StripPort(ip_target))
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s *MessageServer) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "Internal Error")
	err = s.subscribe(r.Context(), c, StripPort(r.RemoteAddr))
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		s.logf("%v", err)
		return
	}
}

func (s *MessageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

// Creates a subscriber and registers the websocket connection to it
// Returns any error encountered, otherwise loops forever until the connection is
// closed.
func (s *MessageServer) subscribe(ctx context.Context, c *websocket.Conn, ip string) error {
	// We will write to this websocket connection and never read from it.
	// We use CloseRead so that we can pass the cancellable context on to the
	// write loop below.
	new_context := c.CloseRead(ctx)
	new_sub := &subscriber{
		msgs: make(chan []byte, 2),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "Connection too slow to keep up")
		},
	}
	s.addSubscriber(ip, new_sub)
	defer s.deleteSubscriber(ip)
	for {
		select {
		case msg := <-new_sub.msgs:
			err := writeTimeout(new_context, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-new_context.Done():
			return new_context.Err()
		}
	}
}

// Adds the given subscriber to the subscribers map.
// Must lock the mutex when editing the global map.
func (s *MessageServer) addSubscriber(ip string, sub *subscriber) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()
	s.subscribers[ip] = sub

}

// Removes a subscriber from the map
func (s *MessageServer) deleteSubscriber(ip string) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()
	delete(s.subscribers, ip)
}

// Parses the msg body and passes the message to the proper channel
func (s *MessageServer) publish(msg []byte, ip_target string) error {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()
	sub, ok := s.subscribers[ip_target]
	if !ok {
		return errors.New(fmt.Sprintf("%v is not in the map of subscribers", ip_target))
	}
	sub.msgs <- msg
	return nil
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Write(ctx, websocket.MessageBinary, msg)
}

func StartServer(address string, users map[string]DeliverConfig) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("listening on %v\n", l.Addr())
	ms := newMessageServer(users)
	s := &http.Server{
		Handler:      ms,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("Failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("Terminating: %v", sig)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return s.Shutdown(ctx)

}
