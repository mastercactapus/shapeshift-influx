package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatalln("Usage: shapeshift-influx <influx_url> <dbname> <2 or more coin types ...>")
	}

	influxURL := strings.TrimSuffix(os.Args[1], "/") + "/write?db=" + os.Args[2]
	types := os.Args[3:]

	n := len(types)
	pairs := make([]string, 0, n*n-n)
	for _, a := range types {
		for _, b := range types {
			if a == b {
				continue
			}
			pairs = append(pairs, a+"_"+b)
		}
	}

	txns := make(chan []transaction, 1)
	ps := make(chan *marketInfo, len(pairs))
	go func() {
		tx, err := getTx()
		if err != nil {
			log.Errorln("Could not get transaction data:", err)
		}
		txns <- tx
	}()
	for _, pair := range pairs {
		go func(p string) {
			m, err := getPair(p)
			if err != nil {
				log.Errorf("Could not get pair '%s': %s", p, err.Error())
			}
			ps <- m
		}(pair)
	}

	transactions := <-txns
	infos := make([]marketInfo, 0, len(pairs))
	for range pairs {
		p := <-ps
		if p != nil {
			infos = append(infos, *p)
		}
	}

	buf := new(bytes.Buffer)
	if transactions != nil {
		for _, t := range transactions {
			buf.WriteString(t.Entry() + "\n")
		}
	}
	for _, p := range infos {
		buf.WriteString(p.Entry() + "\n")
	}
	resp, err := http.Post(influxURL, "", buf)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != 204 {
		data, _ := ioutil.ReadAll(resp.Body)
		log.Errorln(string(data))
		log.Fatalln("Unexpected status code:", resp.Status+"")
	}
}
