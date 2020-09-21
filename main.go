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

const VERSION = "0.0.6"

func Usage() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-debug] [-log logfile] <operation> <args>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-help]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-version]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n<operation> is one of 'view_receipt|transfer|send|inventory'\n")
		fmt.Fprintf(os.Stderr, "<args> arguments for operation\n\n")
		flag.PrintDefaults()
}

func Version() {
		fmt.Printf("%s\n", VERSION)
}

func main() {
	//flag.StringVar(&config.CmdCommand, "", "", "Operation")
	flag.BoolVar(&config.CmdDebug, "debug", false, "Display Debug Information")
	flag.StringVar(&config.CmdLogfile, "logfile", "", "Logfile path")
	flag.BoolVar(&config.CmdHelp, "help", false, "Show Usage")
	flag.BoolVar(&config.CmdVersion, "version", false, "Display version")
	flag.Usage = Usage
	flag.Parse()

	if config.CmdVersion {
		Version()
		os.Exit(0)
	}

	if config.CmdLogfile != "" {
		stat, _ := os.Stat(config.CmdLogfile)
	  if stat != nil {
			if (stat.Size() > config.MAX_LOG_SIZE) {
				core.RotateLog(config.CmdLogfile)
			}
		}

	  file, err0 := os.OpenFile(config.CmdLogfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644);
	  if err0 != nil {
			core.ShowError(config.ERROR_INCORRECT_USAGE, "Failed to open logfile")
		}

		config.LogDesc = file
	}

	if flag.NArg() == 0 {
		if config.CmdHelp {
			Usage()
			os.Exit(0)
		}
		core.ShowError(config.ERROR_INCORRECT_USAGE, "Operation is not specified")
	}

	core.CreateFolders()

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
			core.ShowError(config.ERROR_INCORRECT_USAGE, "Receive requires two arguments: guid and owner")
		}
		r := raida.NewVerifier()
		response, err := r.Receive(uuid, owner)
		if err != nil {
			core.ShowError(err.Code, err.Message)
		}

		fmt.Println(response)
	} else if operation == "send" {
		if (config.CmdHelp) {
				fmt.Fprintf(os.Stderr, "send command transfers coins from your local wallet to a remote Sky Wallet\n\n")
				fmt.Fprintf(os.Stderr, "Usage:\n")
				fmt.Fprintf(os.Stderr, "%s [-debug] send <amount> <destination skywallet> <memo>\n\n", os.Args[0])
				fmt.Fprintf(os.Stderr, "<amount> - amount to transfer\n")
				fmt.Fprintf(os.Stderr, "<destination skywallet> - serial number, ip address, or skywallet address of the receiver\n")
				fmt.Fprintf(os.Stderr, "<memo> - memo\n\n")
				fmt.Fprintf(os.Stderr, "Example:\n")
				fmt.Fprintf(os.Stderr, "%s send 10 ax2.skywallet.cc \"my memo\"\n", os.Args[0])
				os.Exit(0)
		}

		amount, to, memo := flag.Arg(1), flag.Arg(2), flag.Arg(3)
		if (amount == "" || to == "" || memo == "") {
			core.ShowError(config.ERROR_INCORRECT_USAGE, "Amount, To, Memo parameters required: " + os.Args[0] + " send 250 destination.skywallet.cc memo")
		}

		s := raida.NewSender()
		response, err := s.Send(amount, to, memo)
		if err != nil {
			core.ShowError(err.Code, err.Message)
		}
		fmt.Println(response)
	} else if operation == "transfer" {
		if (config.CmdHelp) {
			fmt.Fprintf(os.Stderr, "transfer sends coins from your Sky Wallet to another Sky Wallet. You need to create an 'ID' folder in the current directory and put your Sky Wallet ID Coin there before you can use transfer\n\n")
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fmt.Fprintf(os.Stderr, "%s [-debug] transfer <amount> <destination skywallet> <memo> <idcoin>\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "<amount> - amount to transfer\n")
			fmt.Fprintf(os.Stderr, "<destination skywallet> - serial number, ip address, or skywallet address of the receiver\n")
			fmt.Fprintf(os.Stderr, "<memo> - memo\n\n")
			fmt.Fprintf(os.Stderr, "<idcoin> - full path to the ID coin\n\n")
			fmt.Fprintf(os.Stderr, "Example:\n")
			fmt.Fprintf(os.Stderr, "%s transfer 10 ax2.skywallet.cc \"my memo\" /home/user/my.skywallet.cc.stack\n", os.Args[0])
			os.Exit(0)
		}
		amount, to, memo, idcoin := flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4)
		if (amount == "" || to == "" || memo == "" || idcoin == "") {
			core.ShowError(config.ERROR_INCORRECT_USAGE, "Amount, To, Memo and IDCoin parameters required: " + os.Args[0] + " transfer 250 destination.skywallet.cc memo")
		}

		cc, err := core.GetIDCoinFromPath(idcoin)
		if err != nil {
			core.ShowError(err.Code, err.Message)
		}

		t := raida.NewTransfer()
		response, err := t.Transfer(cc, amount, to, memo)
		if err != nil {
			core.ShowError(err.Code, err.Message)
		}
		fmt.Println(response)

	} else if operation == "show" {
		if (config.CmdHelp) {
				fmt.Fprintf(os.Stderr, "show command shows coins int your local wallet\n\n")
				fmt.Fprintf(os.Stderr, "Usage:\n")
				fmt.Fprintf(os.Stderr, "%s [-debug] show <localwallet>\n\n", os.Args[0])
				fmt.Fprintf(os.Stderr, "<localwallet> local wallet (optional)\n")
				fmt.Fprintf(os.Stderr, "Example:\n")
				fmt.Fprintf(os.Stderr, "%s mywallet\n", os.Args[0])
				os.Exit(0)
		}

		s := raida.NewShowCoins()
		response, err := s.ShowCoins()
		if err != nil {
			core.ShowError(err.Code, err.Message)
		}
		fmt.Println(response)
	} else {

		core.ShowError(config.ERROR_INCORRECT_USAGE, "Invalid command")
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
