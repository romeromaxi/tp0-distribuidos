package common

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
}

type ClientBetConfig struct {
	Name    string
	Surname string
	DNI     string
	Birth   string
	Number  string
}

// Client Entity that encapsulates how
type Client struct {
	config          ClientConfig
	bet             ClientBetConfig
	courier         *Courier
	sigtermReceived chan bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet ClientBetConfig) *Client {
	client := &Client{
		config:          config,
		bet:             bet,
		sigtermReceived: make(chan bool, 1),
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.courier = NewCourier(c.config.ID, conn)
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	defer c.gracefulShutdown()

	c.handleSigterm()

	c.createClientSocket()

	if c.courier == nil {
		return
	}

	err := c.courier.SendMessage("CONN", c.config.ID)

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | CONN error: %v",
			c.config.ID,
			err,
		)
		return
	}

	betMessage := fmt.Sprintf("%s|%s|%s|%s|%s", c.bet.Name, c.bet.Surname, c.bet.DNI, c.bet.Birth, c.bet.Number)
	err = c.courier.SendMessage("BET", betMessage)

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | BET error: %v",
			c.config.ID,
			err,
		)
		return
	}

	recv_msg_type, err := c.courier.RecvTypeMessage()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	if recv_msg_type != "OK" {
		log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v",
			c.bet.DNI,
			c.bet.Number,
		)
		return
	}

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
		c.bet.DNI,
		c.bet.Number,
	)
}

func (c *Client) handleSigterm() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		c.gracefulShutdown()
		//c.sigtermReceived <- true
		//close(c.sigtermReceived)
	}()
}

func (c *Client) gracefulShutdown() {
	if c.courier != nil {
		c.courier.Close()
		log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
	}
}
