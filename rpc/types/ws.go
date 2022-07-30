package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/libs/log"
	mrand "math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rcrowley/go-metrics"
)

var (
	_ Client = &WSClient{}
)

type WSOptions struct {
	MaxReconnectAttempts uint          // maximum attempts to reconnect
	ReadWait             time.Duration // deadline for any read op
	WriteWait            time.Duration // deadline for any write op
	PingPeriod           time.Duration // frequency with which pings are sent
}

func DefaultWSOptions() WSOptions {
	return WSOptions{
		MaxReconnectAttempts: 10, // first: 2 sec, last: 17 min.
		WriteWait:            10 * time.Second,
		ReadWait:             0,
		PingPeriod:           0,
	}
}

type WSClient struct { // nolint: maligned
	*RunState
	conn *websocket.Conn

	Address  string // IP:PORT or /path/to/socket
	Endpoint string // /websocket/url/endpoint
	Dialer   func(string, string) (net.Conn, error)

	ResponsesCh chan RPCResponse

	onReconnect func()

	send            chan RPCRequest // user requests
	backlog         chan RPCRequest // stores a single user request received during a conn failure
	reconnectAfter  chan error      // reconnect requests
	readRoutineQuit chan struct{}   // a way for readRoutine to close writeRoutine

	maxReconnectAttempts uint

	protocol string

	wg sync.WaitGroup

	mtx            sync.RWMutex
	sentLastPingAt time.Time
	reconnecting   bool
	nextReqID      int

	writeWait time.Duration

	readWait time.Duration

	pingPeriod time.Duration

	PingPongLatencyTimer metrics.Timer
}

func NewWS(remoteAddr, endpoint string) (*WSClient, error) {
	return NewWSWithOptions(remoteAddr, endpoint, DefaultWSOptions())
}

func NewWSWithOptions(remoteAddr, endpoint string, opts WSOptions) (*WSClient, error) {
	parsedURL, err := newParsedURL(remoteAddr)
	if err != nil {
		return nil, err
	}
	// default to ws protocol, unless wss is explicitly specified
	if parsedURL.Scheme != protoWSS {
		parsedURL.Scheme = protoWS
	}

	dialFn, err := makeHTTPDialer(remoteAddr)
	if err != nil {
		return nil, err
	}

	c := &WSClient{
		RunState:             NewRunState("WSClient", nil),
		Address:              parsedURL.GetTrimmedHostWithPath(),
		Dialer:               dialFn,
		Endpoint:             endpoint,
		PingPongLatencyTimer: metrics.NewTimer(),

		maxReconnectAttempts: opts.MaxReconnectAttempts,
		readWait:             opts.ReadWait,
		writeWait:            opts.WriteWait,
		pingPeriod:           opts.PingPeriod,
		protocol:             parsedURL.Scheme,

		// sentIDs: make(map[types.JSONRPCIntID]bool),
	}
	return c, nil
}

func (c *WSClient) OnReconnect(cb func()) {
	c.onReconnect = cb
}

func (c *WSClient) String() string {
	return fmt.Sprintf("WSClient{%s (%s)}", c.Address, c.Endpoint)
}

func (c *WSClient) Start() error {
	if err := c.RunState.Start(); err != nil {
		return err
	}
	err := c.dial()
	if err != nil {
		return err
	}

	c.ResponsesCh = make(chan RPCResponse)
	c.send = make(chan RPCRequest)
	c.reconnectAfter = make(chan error, 1)
	c.backlog = make(chan RPCRequest, 1)

	c.startReadWriteRoutines()
	go c.reconnectRoutine()

	return nil
}

func (c *WSClient) Stop() error {
	if err := c.RunState.Stop(); err != nil {
		return err
	}
	c.wg.Wait()
	close(c.ResponsesCh)

	return nil
}

func (c *WSClient) IsReconnecting() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.reconnecting
}

func (c *WSClient) IsActive() bool {
	return c.IsRunning() && !c.IsReconnecting()
}

func (c *WSClient) Send(ctx context.Context, request RPCRequest) error {
	select {
	case c.send <- request:
		c.Logger.Info("sent a request", "req", request)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *WSClient) Call(ctx context.Context, method string, params map[string]interface{}) error {
	request, err := MapToRequest(c.nextRequestID(), method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

func (c *WSClient) CallWithArrayParams(ctx context.Context, method string, params []interface{}) error {
	request, err := ArrayToRequest(c.nextRequestID(), method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

func (c *WSClient) nextRequestID() JSONRPCIntID {
	c.mtx.Lock()
	id := c.nextReqID
	c.nextReqID++
	c.mtx.Unlock()
	return JSONRPCIntID(id)
}

func (c *WSClient) dial() error {
	dialer := &websocket.Dialer{
		NetDial: c.Dialer,
		Proxy:   http.ProxyFromEnvironment,
	}
	rHeader := http.Header{}
	urlStr := c.protocol + "://" + c.Address + c.Endpoint
	conn, _, err := dialer.Dial(urlStr, rHeader)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *WSClient) reconnect() error {
	attempt := uint(0)

	c.mtx.Lock()
	c.reconnecting = true
	c.mtx.Unlock()
	defer func() {
		c.mtx.Lock()
		c.reconnecting = false
		c.mtx.Unlock()
	}()

	for {
		jitter := time.Duration(mrand.Float64() * float64(time.Second)) // 1s == (1e9 ns)
		backoffDuration := jitter + ((1 << attempt) * time.Second)

		c.Logger.Info("reconnecting", "attempt", attempt+1, "backoff_duration", backoffDuration)
		time.Sleep(backoffDuration)

		err := c.dial()
		if err != nil {
			c.Logger.Error("failed to redial", "err", err)
		} else {
			c.Logger.Info("reconnected")
			if c.onReconnect != nil {
				go c.onReconnect()
			}
			return nil
		}

		attempt++

		if attempt > c.maxReconnectAttempts {
			return fmt.Errorf("reached maximum reconnect attempts: %w", err)
		}
	}
}

func (c *WSClient) startReadWriteRoutines() {
	c.wg.Add(2)
	c.readRoutineQuit = make(chan struct{})
	go c.readRoutine()
	go c.writeRoutine()
}

func (c *WSClient) processBacklog() error {
	select {
	case request := <-c.backlog:
		if c.writeWait > 0 {
			if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
				c.Logger.Error("failed to set write deadline", "err", err)
			}
		}
		if err := c.conn.WriteJSON(request); err != nil {
			c.Logger.Error("failed to resend request", "err", err)
			c.reconnectAfter <- err
			// requeue request
			c.backlog <- request
			return err
		}
		c.Logger.Info("resend a request", "req", request)
	default:
	}
	return nil
}

func (c *WSClient) reconnectRoutine() {
	for {
		select {
		case originalError := <-c.reconnectAfter:
			// wait until writeRoutine and readRoutine finish
			c.wg.Wait()
			if err := c.reconnect(); err != nil {
				c.Logger.Error("failed to reconnect", "err", err, "original_err", originalError)
				if err = c.Stop(); err != nil {
					c.Logger.Error("failed to stop conn", "error", err)
				}

				return
			}
			// drain reconnectAfter
		LOOP:
			for {
				select {
				case <-c.reconnectAfter:
				default:
					break LOOP
				}
			}
			err := c.processBacklog()
			if err == nil {
				c.startReadWriteRoutines()
			}

		case <-c.Quit():
			return
		}
	}
}

func (c *WSClient) writeRoutine() {
	var ticker *time.Ticker
	if c.pingPeriod > 0 {
		// ticker with a predefined period
		ticker = time.NewTicker(c.pingPeriod)
	} else {
		// ticker that never fires
		ticker = &time.Ticker{C: make(<-chan time.Time)}
	}

	defer func() {
		ticker.Stop()
		c.conn.Close()
		c.wg.Done()
	}()

	for {
		select {
		case request := <-c.send:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteJSON(request); err != nil {
				c.Logger.Error("failed to send request", "err", err)
				c.reconnectAfter <- err
				// add request to the backlog, so we don't lose it
				c.backlog <- request
				return
			}
		case <-ticker.C:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Logger.Error("failed to write ping", "err", err)
				c.reconnectAfter <- err
				return
			}
			c.mtx.Lock()
			c.sentLastPingAt = time.Now()
			c.mtx.Unlock()
			c.Logger.Debug("sent ping")
		case <-c.readRoutineQuit:
			return
		case <-c.Quit():
			if err := c.conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			); err != nil {
				c.Logger.Error("failed to write message", "err", err)
			}
			return
		}
	}
}

func (c *WSClient) readRoutine() {
	defer func() {
		c.conn.Close()
		c.wg.Done()
	}()

	c.conn.SetPongHandler(func(string) error {
		// gather latency stats
		c.mtx.RLock()
		t := c.sentLastPingAt
		c.mtx.RUnlock()
		c.PingPongLatencyTimer.UpdateSince(t)

		c.Logger.Debug("got pong")
		return nil
	})

	for {
		// reset deadline for every message type (control or data)
		if c.readWait > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.readWait)); err != nil {
				c.Logger.Error("failed to set read deadline", "err", err)
			}
		}
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			c.Logger.Error("failed to read response", "err", err)
			close(c.readRoutineQuit)
			c.reconnectAfter <- err
			return
		}

		var response RPCResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			c.Logger.Error("failed to parse response", "err", err, "data", string(data))
			continue
		}

		if err = validateResponseID(response.ID); err != nil {
			c.Logger.Error("error in response ID", "id", response.ID, "err", err)
			continue
		}

		c.Logger.Info("got response", "id", response.ID, "result", response.Result)

		select {
		case <-c.Quit():
		case c.ResponsesCh <- response:
		}
	}
}

// Subscribe to a query. Note the server must have a "subscribe" route
// defined.
func (c *WSClient) Subscribe(ctx context.Context, query string) error {
	params := map[string]interface{}{"query": query}
	return c.Call(ctx, "subscribe", params)
}

// Unsubscribe from a query. Note the server must have a "unsubscribe" route
// defined.
func (c *WSClient) Unsubscribe(ctx context.Context, query string) error {
	params := map[string]interface{}{"query": query}
	return c.Call(ctx, "unsubscribe", params)
}

// UnsubscribeAll from all. Note the server must have a "unsubscribe_all" route
// defined.
func (c *WSClient) UnsubscribeAll(ctx context.Context) error {
	params := map[string]interface{}{}
	return c.Call(ctx, "unsubscribe_all", params)
}

func validateResponseID(id interface{}) error {
	if id == nil {
		return errors.New("no ID")
	}
	_, ok := id.(JSONRPCIntID)
	if !ok {
		return fmt.Errorf("expected JSONRPCIntID, but got: %T", id)
	}
	return nil
}

type RunState struct {
	Logger log.Logger

	mu        sync.Mutex
	name      string
	isRunning bool
	quit      chan struct{}
}

// NewRunState returns a new unstarted run state tracker with the given logging
// label and log sink. If logger == nil, a no-op logger is provided by default.
func NewRunState(name string, logger log.Logger) *RunState {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &RunState{
		name:   name,
		Logger: logger,
	}
}

// Start sets the state to running, or reports an error.
func (r *RunState) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.isRunning {
		r.Logger.Error("not starting client, it is already started", "client", r.name)
		return ErrClientRunning
	}
	r.Logger.Info("starting client", "client", r.name)
	r.isRunning = true
	r.quit = make(chan struct{})
	return nil
}

// Stop sets the state to not running, or reports an error.
func (r *RunState) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.isRunning {
		r.Logger.Error("not stopping client; it is already stopped", "client", r.name)
		return ErrClientNotRunning
	}
	r.Logger.Info("stopping client", "client", r.name)
	r.isRunning = false
	close(r.quit)
	return nil
}

// SetLogger updates the log sink.
func (r *RunState) SetLogger(logger log.Logger) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Logger = logger
}

// IsRunning reports whether the state is running.
func (r *RunState) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.isRunning
}

// Quit returns a channel that is closed when a call to Stop succeeds.
func (r *RunState) Quit() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.quit
}
