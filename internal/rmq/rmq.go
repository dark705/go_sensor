package rmq

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/streadway/amqp"
)

type Config struct {
	Login    string `yaml:"user"`
	Pass     string `yaml:"password"`
	Host     string `yaml:"ip"`
	Port     string `yaml:"port"`
	Timeout  int    `yaml:"timeout"`
	Queue    string `yaml:"queue"`
	Exchange string `yaml:"exchange"`
}

type RMQ struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
}

func NewRMQ(conf Config) (r *RMQ, err error) {
	log.Println("Start connect to RMQ")
	r = &RMQ{}
	r.conn, err = amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s/", conf.Login, conf.Pass, net.JoinHostPort(conf.Host, conf.Port)),
		amqp.Config{Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, time.Second*time.Duration(conf.Timeout))
		}})
	if err != nil {
		return r, err
	}

	r.ch, err = r.conn.Channel()
	if err != nil {
		return r, err
	}

	err = r.ch.ExchangeDeclare(conf.Exchange, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return r, err
	}

	r.exchange = conf.Exchange

	log.Println("Success connected to RMQ")
	return r, nil
}

func (r *RMQ) Send(message []byte) error {
	return r.ch.Publish(r.exchange, "", false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Body:         message,
		})

}

func (r *RMQ) CloseConnect() (err error) {
	err = r.ch.Close()
	if err != nil {
		_ = r.conn.Close()
		return err
	}
	return r.conn.Close()
}
