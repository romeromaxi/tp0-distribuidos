package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/op/go-logging"
)

const BET_FIELDS_NUMBER = 5

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	FileName      string
	FileDelimiter string
	BatchSize     int
}

// Client Entity that encapsulates how
type Client struct {
	config          ClientConfig
	courier         *Courier
	sigtermReceived chan bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:          config,
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

func (c *Client) SendConnectionMessage() error {
	connection_message := GetConnectionMessage(c.config.ID)
	err := c.courier.SendMessage(connection_message)

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | CONN error: %v",
			c.config.ID,
			err,
		)
		return err
	}

	return nil
}

func (c *Client) SendBetsBatchMessageAndRecv(betsInBatch [][]string) error {
	number_of_bets := len(betsInBatch)
	nbets_message := GetBetsBatchMessage(betsInBatch)

	err := c.courier.SendMessage(nbets_message)

	if err != nil {
		return err
	}

	recv_msg_type, err := c.courier.RecvTypeMessage()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}

	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		recv_msg_type,
	)

	if recv_msg_type != MESSAGE_TYPE_OK_RESPONSE {
		nok_reponse := fmt.Errorf("action: apuesta_enviada | result: fail | client_id: %v | cantidad: %v",
			c.config.ID,
			number_of_bets,
		)
		return nok_reponse
	}

	return nil
}

func (c *Client) SendEndBetsMessage() error {
	end_message := GetEndBetsMessage()
	err := c.courier.SendMessage(end_message)

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | END error: %v",
			c.config.ID,
			err,
		)
		return err
	}

	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	c.handleSigterm()
	defer c.gracefulShutdown()

	err := c.createClientSocket()

	if err != nil {
		return
	}

	err = c.SendConnectionMessage()

	if err != nil {
		return
	}

	file, err := os.Open(c.config.FileName)
	if err != nil {
		log.Errorf("action: open_file | result: fail | error: %v", err)
		return
	}
	defer file.Close()

	log.Debugf("action: open_file | result: success")

	scanner := bufio.NewScanner(file)

	currentBatchSize := 0
	lineNumber := 0
	var betsInCurrentBatch [][]string

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		fields := strings.Split(line, c.config.FileDelimiter)

		if len(fields) < BET_FIELDS_NUMBER {
			log.Errorf("action: read_line | result: fail | error: BAD FORMAT LINE %v", lineNumber)
			continue
		}

		if currentBatchSize < c.config.BatchSize {
			betsInCurrentBatch = append(betsInCurrentBatch, fields)
			currentBatchSize++
		}

		if currentBatchSize == c.config.BatchSize {
			err = c.SendBetsBatchMessageAndRecv(betsInCurrentBatch)

			if err != nil {
				break
			}
			currentBatchSize = 0
			betsInCurrentBatch = nil
		}
	}

	if currentBatchSize > 0 {
		err = c.SendBetsBatchMessageAndRecv(betsInCurrentBatch)
	}

	if err == nil {
		err = scanner.Err()
		if err != nil {
			log.Errorf("action: open_file | result: fail | error: %v", err)
			return
		}

		c.SendEndBetsMessage()
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) handleSigterm() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		c.gracefulShutdown()
		c.sigtermReceived <- true
		close(c.sigtermReceived)
	}()
}

func (c *Client) gracefulShutdown() {
	if c.courier != nil {
		c.courier.Close()
		log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
	}
}
