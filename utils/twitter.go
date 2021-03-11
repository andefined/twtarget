package utils

import (
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

// Credentials ...
type Credentials struct {
	ConsumerKey       string `yaml:"consumer-key"`
	ConsumerSecret    string `yaml:"consumer-secret"`
	AccessToken       string `yaml:"access-token"`
	AccessTokenSecret string `yaml:"access-token-secret"`
}

// API ...
func API(creds *Credentials) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(creds.ConsumerKey)
	anaconda.SetConsumerSecret(creds.ConsumerSecret)
	api := anaconda.NewTwitterApi(creds.AccessToken, creds.AccessTokenSecret)

	if _, err := api.VerifyCredentials(); err != nil {
		log.Println("Bad Authorization Tokens. Please refer to https://apps.twitter.com/ for your Access Tokens.")
		os.Exit(1)
	}

	return api
}

// StreamHandler ...
type StreamHandler struct {
	All   func(m interface{})
	Tweet func(tweet anaconda.Tweet)
	Other func(m interface{})
}

// NewStreamHandler ...
func NewStreamHandler() *StreamHandler {
	return &StreamHandler{
		All:   func(m interface{}) {},
		Tweet: func(tweet anaconda.Tweet) {},
		Other: func(m interface{}) {},
	}
}

// Handle ...
func (d StreamHandler) Handle(m interface{}) {
	switch t := m.(type) {
	case anaconda.Tweet:
		d.Tweet(t)
	default:
		d.Other(t)
	}
}

// HandleChan ...
func (d StreamHandler) HandleChan(c <-chan interface{}) {
	for m := range c {
		d.Handle(m)
	}
}
