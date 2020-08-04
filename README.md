# raidaGo

RAIDA GO Console program allows you to verify that you have received funds in your Skwyallet and to send fund to another Skywallet account from your Skywallet. You can find both the Linux and Windows version at: https://CloudCoin.global/assets/raida_go.zip


Usage of ./raida_go on Linux Systems:
```console
./raida_go [-debug] <operation> <args>
./raida_go [-help]

<operation> is one of 'receive|send'
<args> arguments for operation

  -debug
        Display Debug Information
  -help
        Show Usage

```

Linux Example of how to check how many CloudCoins were sent to the merchant.mydomain.com Skywallet account with a guid in the memo:
```console
$ ./raida_go receive 080A4CE89126F4F1B93E4745F89F6713 merchant.mydomain.com
{"amount_verified":150}
```
Same Example in Windows:
```console
C:/xampp/htdocs/cloudcoin/raida_go receive 080A4CE89126F4F1B93E4745F89F6713 merchant.mydomain.com
{"amount_verified":150}
```
To see additional Debug Info:

```console
$ ./raida_go -debug receive 080A4CE89126F4F1B93E4745F89F6713 merchant.mydomain.com
```
