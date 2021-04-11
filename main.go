package main

import (
	"fmt"
	"raida"
	"logger"
	"flag"
	"config"
	"os"
	"core"
	"strings"
	"cloudcoin"
	"error"
)

const VERSION = "0.0.11"

func Usage() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-debug] <operation> <args>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-help]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [-version]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n<operation> is one of 'view_receipt|verify_payment|transfer|send|inventory|balance'\n")
		fmt.Fprintf(os.Stderr, "<args> arguments for operation\n\n")
		flag.PrintDefaults()
}

func Version() {
		fmt.Printf("%s\n", VERSION)
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
		if config.CmdHelp {
			Usage()
			os.Exit(0)
		}
		core.ShowError(config.ERROR_INCORRECT_USAGE, "Operation is not specified")
	}

	core.CreateFolders()
	core.InitLog()
	core.ReadConfig()

	logger.Debug("Program started")
	logger.Debug(strings.Join(os.Args, " "))

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
	} else if operation == "verify_payment" {
		if (config.CmdHelp) {
			fmt.Fprintf(os.Stderr, "verify_payment allows the sender and receiver of a payment to verify the reciept\n\n")
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fmt.Fprintf(os.Stderr, "%s [-debug] verify_payment <uuid> [<idcoin]\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "<uuid> - uuid of the receipt\n")
			fmt.Fprintf(os.Stderr, "<idcoin> - ID coin path\n\n")
			fmt.Fprintf(os.Stderr, "Example:\n")
			fmt.Fprintf(os.Stderr, "%s verify_payment 080A4CE89126F4F1B93E4745F89F6713\n", os.Args[0])
			os.Exit(0)
		}

		memo := flag.Arg(1)
		if (memo == "") {
			core.ShowError(config.ERROR_INCORRECT_USAGE, "Memo parameter required: " + os.Args[0])
		}

		var cc *cloudcoin.CloudCoin
		var err *error.Error
		if (flag.NArg() == 2) {
			cc, err = core.GetIDCoin()
		} else if (flag.NArg() == 3) {
			idcoin := flag.Arg(2)
			cc, err = core.GetIDCoinFromPath(idcoin)
		} else {
			core.ShowError(config.ERROR_INCORRECT_USAGE, "Memo parameter required: " + os.Args[0])
		}
		if err != nil {
			core.ShowError(err.Code, err.Message)
		}

		r := raida.NewPaymentVerifier()
		response, err := r.Verify(memo, cc)
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
			fmt.Fprintf(os.Stderr, "%s [-debug] transfer <amount> <destination skywallet> <memo> [<idcoin>]\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "<amount> - amount to transfer\n")
			fmt.Fprintf(os.Stderr, "<destination skywallet> - serial number, ip address, or skywallet address of the receiver\n")
			fmt.Fprintf(os.Stderr, "<memo> - memo\n\n")
			fmt.Fprintf(os.Stderr, "<idcoin> - full path to the ID coin. If not defined it will be taken from the ID folder\n\n")
			fmt.Fprintf(os.Stderr, "Example:\n")
			fmt.Fprintf(os.Stderr, "%s transfer 10 ax2.skywallet.cc \"my memo\" /home/user/my.skywallet.cc.stack\n", os.Args[0])
			os.Exit(0)
		}
		amount, to, memo := flag.Arg(1), flag.Arg(2), flag.Arg(3)
		if (amount == "" || to == "" || memo == "") {
			core.ShowError(config.ERROR_INCORRECT_USAGE, "Amount, To, Memo parameters required: " + os.Args[0] + " transfer 250 destination.skywallet.cc memo")
		}

		var cc *cloudcoin.CloudCoin
		var err *error.Error
		if (flag.NArg() == 4) {
			cc, err = core.GetIDCoin()
		} else {
			idcoin := flag.Arg(4)
			cc, err = core.GetIDCoinFromPath(idcoin)
		}
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
	} else if operation == "balance" {
		if (config.CmdHelp) {
				fmt.Fprintf(os.Stderr, "balance command shows your skywallet balance\n\n")
				fmt.Fprintf(os.Stderr, "Usage:\n")
				fmt.Fprintf(os.Stderr, "<idcoin> - full path to the ID coin. If not defined it will be taken from the ID folder\n\n")
				os.Exit(0)
		}
		var cc *cloudcoin.CloudCoin
		var err *error.Error
		if (flag.NArg() == 1) {
			cc, err = core.GetIDCoin()
		} else {
			idcoin := flag.Arg(1)
			cc, err = core.GetIDCoinFromPath(idcoin)
		}
		if err != nil {
			core.ShowError(err.Code, err.Message)
		}

		b := raida.NewShowTransferBalance()
		response, err := b.ShowTransferBalance(cc)
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
