package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	//	"github.com/gosnmp/gosnmp"
	"github.com/twsnmp/gosnmp"
)

var startTime time.Time

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

// function for responding requests
func getSysDescr(oid string) interface{} {
	return "test"
}

func getSysObjectID(oid string) interface{} {
	return ".1.3.6.1.2.1.1.1.0"
}

func getSysUpTime(oid string) interface{} {
	return uint32((time.Now().UnixNano() - startTime.UnixNano()) / (1000 * 1000 * 10))
}

func getSysContact(oid string) interface{} {
	return "test sysContact"
}

func getSysName(oid string) interface{} {
	return "test sysName"
}

func getSysLocation(oid string) interface{} {
	return "test sysLocation"
}

func getSysServices(oid string) interface{} {
	return 72
}

// function for initializing MIB
func initMib(a *gosnmp.GoSNMPAgent) {
	a.AddMibList(".1.3.6.1.2.1.1.1.0", gosnmp.OctetString, getSysDescr)
	a.AddMibList(".1.3.6.1.2.1.1.2.0", gosnmp.ObjectIdentifier, getSysObjectID)
	a.AddMibList(".1.3.6.1.2.1.1.3.0", gosnmp.TimeTicks, getSysUpTime)
	a.AddMibList(".1.3.6.1.2.1.1.7.0", gosnmp.Integer, getSysServices)
	a.AddMibList(".1.3.6.1.2.1.1.4.0", gosnmp.OctetString, getSysContact)
	a.AddMibList(".1.3.6.1.2.1.1.5.0", gosnmp.OctetString, getSysName)
	a.AddMibList(".1.3.6.1.2.1.1.6.0", gosnmp.OctetString, getSysLocation)
}
