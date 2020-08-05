package main

import (
	"fmt"
	"raida"
	//"logger"
	"flag"
	"config"
	"os"
	"core"
)

func Usage() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-debug] <operation> <args>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-help]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n<operation> is one of 'receive|transfer'\n")
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

	core.MkDirs()

	operation := flag.Arg(0)
	if operation == "receive" {
		uuid := flag.Arg(1)
		owner := flag.Arg(2)
		if (owner == "") {
			fmt.Fprintf(os.Stderr, "Receive requires two arguments: guid and owner\n")
			os.Exit(1)
		}
		r := raida.NewVerifier()
		response, err := r.Receive(uuid, owner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Message)
			os.Exit(1)
		}

		fmt.Println(response)
	} else if operation == "transfer" {
		amount, to, memo := flag.Arg(1), flag.Arg(2), flag.Arg(3)
		if (amount == "" || to == "") {
			fmt.Fprintf(os.Stderr, "Amount and To parameters required: %s transfer 250 destination.skywallet.cc memo\n", os.Args[0])
			os.Exit(1)
		}

		t := raida.NewTransfer()
		response, err := t.Transfer(amount, to, memo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Message)
			os.Exit(1)
		}
		fmt.Println(response)

	} else {
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



	os.Exit(1)

}
