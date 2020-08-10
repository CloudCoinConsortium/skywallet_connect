# raidaGo

RAIDA GO Console program allows you to verify that you have received funds in your Skwyallet and to send fund to another Skywallet account from your Skywallet. You can find both the Linux and Windows version at: https://CloudCoin.global/assets/raida_go.zip


Usage of ./raida_go on Linux Systems:
```console
./raida_go [-debug] <operation> <args>
./raida_go [-help]

<operation> is one of 'view_receipt|transfer'
<args> arguments for operation

  -debug
        Display Debug Information
  -help
        Show Usage

```

Linux Example of how to check how many CloudCoins were sent to the merchant.mydomain.com Skywallet account with a guid in the memo:

## View Receipt

```console
$ ./raida_go view_receipt 080A4CE89126F4F1B93E4745F89F6713 merchant.mydomain.com
{"amount_verified":100,"status":"success","message":"CloudCoins verified"}
```
Same Example in Windows:
```console
C:\xampp\htdocs\cloudcoin\raida_go.exe view_receipt 080A4CE89126F4F1B93E4745F89F6713 merchant.mydomain.com
{"amount_verified":100,"status":"success","message":"CloudCoins verified"}
```
To see additional Debug Info:

```console
$ ./raida_go -debug view_receipt 080A4CE89126F4F1B93E4745F89F6713 merchant.mydomain.com
```


## Transfer:

format: ./raida_go tranfer <amount> <destination_skywallet> <memo>

Example:

```console
$ ./raida_go transfer 2 ax2.skywallet.cc "my memo"
{"amount_sent":2,"Message":"CloudCoins sent","Status":"success"}
```
