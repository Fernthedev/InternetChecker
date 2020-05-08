package main

import (
	"InternetChecker/internal/config"
	"InternetChecker/pkg/network"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	succeedPing := false

	var silent bool

	for _,i := range os.Args {
		if i == "silent" || strings.Contains(i, "silent") {
			silent = true
			break
		}
	}

	ifaces := network.GetInterfaces()

	configuration, err := config.LoadConfiguration("./config.json")

	logIfNoneSilent(silent, "Loaded configuration")

	if err != nil {
		if errVal := err.(*os.PathError); errVal != nil {
			configuration, err = config.SaveDefaultConfig("./config.json")

			if err != nil {
				log.Fatal("Unable to load error: " + err.Error())
			}
		} else {
			log.Fatal("Unable to load error: " + err.Error())
		}
	}

	p := fastping.NewPinger()


	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {

		logIfNoneSilent(silent, "IP Addr: %s receive, RTT: %v", addr.String(), rtt)
		succeedPing = true

		os.Exit(0)
	}

	p.OnIdle = func() {
		logIfNoneSilent(silent, "finish")
	}


	var found bool
	for _, i := range ifaces {
		if _, ok := configuration.CheckNetworkInterfaces[i.Name]; ok && i.StatusBool {
			found = true

				for _, ipSetting := range configuration.CheckNetworkInterfaces[i.Name].PingAddressSetting {
					ra, err := net.ResolveIPAddr("ip4:icmp", ipSetting.PingAddress)
					if err != nil {
						fmt.Println(err)
					}

					logIfNoneSilent(silent, "Queuing to ping %v", ipSetting.PingAddress)
					p.AddIPAddr(ra)
				}
		} else {
			logIfNoneSilent(silent, "Could not find %v in config", i.Name)
		}
	}

	if found {
		for _, i := range configuration.GlobalPingAddresses {
			ra, err := net.ResolveIPAddr("ip4:icmp", i.PingAddress)
			if err != nil {
				fmt.Println(err)
			}

			p.AddIPAddr(ra)
		}
	}


	logIfNoneSilent(silent, "Pinging sites")
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}

	if !succeedPing {
		os.Exit(2)
	}
}


func logIfNoneSilent(silent bool, format string, a ...interface{}) (n int, err error) {
	if !silent {


		if len(a) == 0 {
			return fmt.Printf(format + "\n")
		} else {
			return fmt.Printf(format + "\n", a)
		}


	}

	return 0, nil
}
