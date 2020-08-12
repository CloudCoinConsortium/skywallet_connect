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

const VERSION = "0.0.1"

func Usage() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-debug] <operation> <args>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-help]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-version]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n<operation> is one of 'view_receipt|transfer'\n")
		fmt.Fprintf(os.Stderr, "<args> arguments for operation\n\n")
		flag.PrintDefaults()
}

func Version() {
		fmt.Fprintf(os.Stderr, "%s\n", VERSION)
}

func main() {
	//flag.StringVar(&config.CmdCommand, "", "", "Operation")
	flag.BoolVar(&config.CmdDebug, "debug", false, "Display Debug Information")
	flag.BoolVar(&config.CmdHelp, "help", false, "Show Usage")
	flag.BoolVar(&config.CmdVersion, "version", false, "Display version")
	flag.Usage = Usage
	flag.Parse()

	if config.CmdVersion {
		Version()
		os.Exit(0)
	}


	if flag.NArg() == 0 {
		Usage()
		os.Exit(1)
	}

	core.MkDirs()

	operation := flag.Arg(0)
	if operation == "view_receipt" {
		if (config.CmdHelp) {
			fmt.Fprintf(os.Stderr, "view_receipt checks if the receipt with a given uuid exists and shows amount of coins in it\n\n")
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fmt.Fprintf(os.Stderr, "%s [-debug] view_receipt <uuid> <skywallet>\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "<uuid> - uuid of the receipt\n")
			fmt.Fprintf(os.Stderr, "<skywallet> - serial number, ip address, or skywallet address\n\n")
			fmt.Fprintf(os.Stderr, "Example:\n")
			fmt.Fprintf(os.Stderr, "%s view_receipt 080A4CE89126F4F1B93E4745F89F6713 demo.skywallet.cc\n", os.Args[0])
			os.Exit(0)
		}

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
		if (config.CmdHelp) {
			fmt.Fprintf(os.Stderr, "transfer sends coins from your Sky Wallet to another Sky Wallet. You need to create an 'ID' folder in the current directory and put your Sky Wallet ID Coin there before you can use transfer\n\n")
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fmt.Fprintf(os.Stderr, "%s [-debug] transfer <amount> <destination skywallet> <memo>\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "<amount> - amount to transfer\n")
			fmt.Fprintf(os.Stderr, "<destination skywallet> - serial number, ip address, or skywallet address of the receiver\n")
			fmt.Fprintf(os.Stderr, "<memo> - memo\n\n")
			fmt.Fprintf(os.Stderr, "Example:\n")
			fmt.Fprintf(os.Stderr, "%s transfer 10 ax2.skywallet.cc \"my memo\"\n", os.Args[0])
			os.Exit(0)
		}
		amount, to, memo := flag.Arg(1), flag.Arg(2), flag.Arg(3)
		cc, err := core.GetIDCoin()
		if err != nil {
			fmt.Printf("%s", core.JsonError(err.Message))
			os.Exit(1)
		}
		if (amount == "" || to == "" || memo == "") {
			fmt.Printf("%s", core.JsonError("Amount, To and Memo parameters required: " + os.Args[0] + " transfer 250 destination.skywallet.cc memo"))
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
		if config.CmdHelp {
			Usage()
			os.Exit(0)
		}
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
