package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	//	"github.com/gosnmp/gosnmp"
	"github.com/twsnmp/gosnmp"
)

func main() {
	// SNMP エージェント設定
	g := &gosnmp.GoSNMP{}
	g.Port = 161
	g.Community = "public"
	g.Version = gosnmp.Version2c
	//	g.Logger := log.New(os.Stdout, "", log.LstdFlags)
	agent := &gosnmp.GoSNMPAgent{
		Port:   161,
		IPAddr: "0.0.0.0",
		Logger: g.Logger,
		Snmp:   g,
	}

	// MIB の初期化
	initMib(agent)

	// サーバーを起動
	go func() {
		log.Println("Starting SNMP Agent...")
		if err := agent.Start(); err != nil {
			log.Fatalf("Error starting SNMP Agent: %v", err)
		}
	}()
	//	defer a.Snmp.Conn.Close()

	// シグナルでの終了処理
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down SNMP Agent...")
}

func initMib(a *GoSNMPAgent) {
	a.AddMibList(".1.3.6.1.2.1.1.1.0", gosnmp.OctetString, getSysDescr)
}

func getSysDescr(oid string) interface{} {
	return "test"
}
