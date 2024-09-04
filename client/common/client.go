package common

import (
	"bufio"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	FileName      string
	FileDelimiter string
	BatchSize     int
}

// Client Entity that encapsulates how
type Client struct {
	config          ClientConfig
	messageHandler  *MessageHandler
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
		return err
	}
	c.messageHandler = NewMessageHandler(c.config.ID, conn)
	return nil
}

func (c *Client) getWinnersLoop() error {
	has_result := false

	for msgID := 1; msgID <= c.config.LoopAmount && !has_result; msgID++ {
		select {
		case <-c.sigtermReceived:
			return nil

		default:

			response_winners, err := c.messageHandler.SendGetWinnersMessageAndRecv()

			if err != nil {
				log.Errorf("action: consulta_ganadores | result: fail | error: %v", err)
				return err
			}

			if response_winners.HasResult {
				log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", response_winners.NumberOfWinners)
				has_result = true
				break
			}

			if c.messageHandler != nil {
				c.messageHandler.Close()
			}

			time.Sleep(c.config.LoopPeriod)

			err = c.createClientSocket()

			if err != nil {
				return err
			}

			err = c.messageHandler.SendConnectionMessage()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	c.handleSigterm()
	defer func() {
		c.gracefulShutdown()
	}()

	err := c.createClientSocket()

	if err != nil {
		return
	}

	err = c.messageHandler.SendConnectionMessage()

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
			err = c.messageHandler.SendBetsBatchMessageAndRecv(betsInCurrentBatch)

			if err != nil {
				break
			}
			currentBatchSize = 0
			betsInCurrentBatch = nil
		}
	}

	if currentBatchSize > 0 {
		err = c.messageHandler.SendBetsBatchMessageAndRecv(betsInCurrentBatch)
	}

	if err == nil {
		err = scanner.Err()
		if err != nil {
			log.Errorf("action: open_file | result: fail | error: %v", err)
			return
		}
	}

	err = c.messageHandler.SendEndBetsMessage()
	if err != nil {
		log.Errorf("action: fin_apuestas | result: fail | error: %v", err)
		return
	}

	err = c.getWinnersLoop()

	if err != nil {
		log.Errorf("action: winners_loop | result: fail | error: %v", err)
		return
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
	if c.messageHandler != nil {
		c.messageHandler.Close()
		log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
	}
}
