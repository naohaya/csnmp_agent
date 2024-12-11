# snmp_agent
This is an implementation of a SNMP agent using GoSNMP (twsnmp version).

## To test the implementation.
On the agent side.
`go run main.go`

On the client side at localhost.
`snmpget -v2c -c public localhost:161 .1.3.6.1.2.1.1.1.0`

