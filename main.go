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
		fmt.Fprintf(os.Stderr, "\n<operation> is one of 'view_receipt|transfer'\n")
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
	if operation == "view_receipt" {
		uuid := flag.Arg(1)
		owner := flag.Arg(2)
		if (owner == "") {
			fmt.Printf("%s", core.JsonError("Receive requires two arguments: guid and owner"))
			os.Exit(1)
		}
		r := raida.NewVerifier()
		response, err := r.Receive(uuid, owner)
		if err != nil {
			fmt.Printf("%s", core.JsonError(err.Message))
			os.Exit(1)
		}

		fmt.Println(response)
	} else if operation == "transfer" {
		amount, to, memo := flag.Arg(1), flag.Arg(2), flag.Arg(3)
		cc, err := core.GetIDCoin()
		if err != nil {
			fmt.Printf("%s", core.JsonError("Failed to find ID coin"))
			os.Exit(1)
		}
		if (amount == "" || to == "") {
			fmt.Printf("%s", core.JsonError("Amount and To parameters required: " + os.Args[0] + " transfer 250 destination.skywallet.cc memo"))
			os.Exit(1)
		}

		t := raida.NewTransfer()
		response, err := t.Transfer(cc, amount, to, memo)
		if err != nil {
			fmt.Printf("%s", core.JsonError(err.Message))
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



	os.Exit(0)
}
