package main

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt" // for debug

	"github.com/twsnmp/gosnmp" //	"github.com/gosnmp/gosnmp"
	// "ecrypto" // Ego's crypto package
)

var startTime time.Time
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var key = []byte("astaxie12798akljzmknm.ahkjkljl;k")

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
	plaintext := []byte("test")
	// create a new AES cipher block
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// create a new CFB encrypter
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	fmt.Printf("%s => %x\n", plaintext, ciphertext)

	// create a new CFB decrypter
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(plaintext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	fmt.Printf("%x => %s\n", ciphertext, plaintextCopy)

	//	return "test"
	return string(ciphertext)
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
