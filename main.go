package main

import (
	"fmt"
	"raida"
	"logger"
	"flag"
	"config"
	"os"
)

func Usage() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-debug] <operation> <args>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-help]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n<operation> is one of 'receive|send'\n")
		fmt.Fprintf(os.Stderr, "<args> arguments for operation\n\n")
		flag.PrintDefaults()
}


func main() {
	//flag.StringVar(&config.CmdCommand, "", "", "Operation")
	flag.BoolVar(&config.CmdDebug, "debug", false, "Display Debug Information")
	flag.BoolVar(&config.CmdHelp, "help", false, "Show Usage")
	flag.Usage = Usage
	flag.Parse()

	if config.CmdHelp {
		Usage()
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		Usage()
		os.Exit(1)
	}

	operation := flag.Arg(0)
	if operation != "receive" {
		Usage()
		os.Exit(1)
	}

	//fmt.Printf("cmd=%d\n",flag.NArg(), flag.Arg(0), flag.Arg(1))
	/*args := os.Args[1:]

	fmt.Println("cmd=", config.CmdCommand)
	for _, e := range args {
		fmt.Println("arg " + e)
	}
*/

	r := raida.NewVerifier()

	response, err := r.Receive("222")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to Receive Coins")
		os.Exit(1)
	}

	fmt.Println(response)

	os.Exit(1)

	/*
	var i interface{}

	i = "xxx"
	i = 2
	i = 3.2

	r := i.(float64)
	y := i.(string)

	fmt.Println("x=", r)
	fmt.Println("x=", y)

	*/

	//client := httpclient.New("https://raida14.cloudcoin.global")

/*
	da := raida.NewDetectionAgent(142)
	if resp, err := da.DoRequest("/service/echo"); err != nil {
		fmt.Println("error")
	} else {

		fmt.Println("r=", resp)
		}
	*/

	
	raida := raida.New()
	logger.Info("xxx")

	raida.SendRequest("/service/echo", nil)







//	client.Send("/service/echo")
	//httpclient.Send("ggg")
}
