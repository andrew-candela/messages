package messages

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	WEBSOCKET_PROTOCOL = "UDPMWS"
	SIG_TOKEN_HEADER   = "RSA_SIG_TOKEN"
	SIG_HEADER         = "RSA_SIG"
)

// Sets up a websocket connection to the server and
// writes a message
func ExampleClient(webserverAddress string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://"+webserverAddress, &websocket.DialOptions{
		Subprotocols: []string{WEBSOCKET_PROTOCOL},
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "There was an error")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message_bytes := scanner.Bytes()
		err = wsjson.Write(ctx, c, string(message_bytes))
		if err != nil {
			log.Panicln(err)
			return
		}

	}
	c.Close(websocket.StatusNormalClosure, "Bye")
}

type MessageServer struct {
	logf          func(f string, v ...interface{})
	serveMux      http.ServeMux
	subscribersMu sync.Mutex
	subscribers   map[*subscriber]struct{}
	config        map[string]DeliverConfig
}

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func newMessageServer(config map[string]DeliverConfig) *MessageServer {
	ms := &MessageServer{
		logf:        log.Printf,
		subscribers: make(map[*subscriber]struct{}),
		config:      config,
	}
	ms.serveMux.HandleFunc("/config", ms.authenticateRequest(ms.handleConfig))
	ms.serveMux.HandleFunc("/ip", ms.authenticateRequest(ms.handleIP))
	ms.serveMux.HandleFunc("/publish", ms.authenticateRequest(ms.handlePublish))
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

}

func (s *MessageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func StartServer(address string, users map[string]DeliverConfig) error {
	fmt.Println("Trying to start server")
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("listening on http://%v\n", l.Addr())
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
