# raidaGo

RAIDA GO Console program allows you to verify that you have received funds in your Skwyallet and to send funds to another Skywallet account from your Skywallet. You can find both the Linux and Windows version at: https://CloudCoin.global/assets/raida_go.zip

## Folder Structure
Some of the raida_go commands require that you have a Skywallet ID coin to work. You will need to created a folder called "ID" in the current directory.

## Example Usage

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
  -version
        Display Version
```

You can add -help parameter to any specific operation 

```console
./raida_go -help view_receipt
./raida_go -help transfer
```


Linux Example of how to check how many CloudCoins were sent to the merchant.mydomain.com Skywallet account with a guid in the memo:

## View Receipt
View receipt allows you to see the money that someone sent to your Skywallet. You must provide 
your account name and the GUID the customer has sent you in their memo. Note, the Skywallet point of service (POSJS) software will generate a guid for the customer and they will never see it.  

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

format: ./raida_go transfer <amount> <destination_skywallet> memo

Example:

```console
$ ./raida_go transfer 2 ax2.skywallet.cc "my memo"
{"amount_sent":2,"Message":"CloudCoins sent","Status":"success"}
```

The list of possible errors:

<pre>
{"status":"error", "message":"Failed to find ID coin, please create a folder called ID in the same folder as your raida_go program. Place one ID coins in that folder"}
{"status":"error", "message":"Failed to parse ID Coin"}
{"status":"error", "message":"Failed to generate random string"} // The program failed to generate random hex-string
{"status":"error", "message":"Failed to convert IP octet1"} // The program failed to convert IP address (four octets xxx.xxx.xxx.xxx) to a serial number
{"status":"error", "message":"Failed to convert IP octet2"} 
{"status":"error", "message":"Failed to convert IP octet3"} 
{"status":"error", "message":"Failed to get SN from IP"} 
{"status":"error", "message":"Invalid Destination Address"}  // Input parameters validation
{"status":"error", "message":"Invalid amount"}  // Input parameters validation
{"status":"error", "message":"Stack File is Corrupted"} 
{"status":"error", "message":"Not enough coins"}  
{"status":"error", "message":"Failed to Show Coins"}   // Show service failed
{"status":"error", "message":"Failed to pick coins"}  // The program can't find neither required amount nor a coin for breaking
{"status":"error", "message":"Failed to pick coins after change"}  // If change is needed and it is successfull, but the program still can't find coins
{"status":"error", "message":"Results from the RAIDA are not synchronized"}  
{"status":"error", "message":"Failed to Encode JSON"}  
{"status":"error", "message":"Failed to Break Coin: ..."}  
{"status":"error", "message":"Failed to Transfer: ..."}  
{"status":"error", "message":"Failed to get Change Method"}  // The program doesn't know how to break coin
{"status":"error", "message":"ShowChange results are not synchronized"}  // The program can't receive trustworthy results after ShowChange
</pre>

Example Error if the amount to transfer, the skwyallet of the person to transfer it to or the memo is left out:

```console
{"status":"fail", "message":"Amount, To and Memo parameters required: raida_go transfer 250 destination.skywallet.cc memo "}
```
