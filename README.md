# raidaGo


Note: Before you can run the raida_go program on Linux, you must first make it an executable by running: 
```bash
chmod +x raida_go
```
RAIDA GO Console program allows you to verify that you have received funds in your Skwyallet and to send fund to another Skywallet account from your Skywallet. You can find both the Linux and Windows version at: https://CloudCoin.global/assets/raida_go.zip

## Folder Structure
Some of the raida_go commands require that you have a Skywallet ID coin to work. You will need to created a folder called "ID" in the current directory.

## Example Usage

Usage of ./raida_go on Linux Systems:
```console
./raida_go [-debug] [-logfile logfile] <operation> <args>
./raida_go [-help]
./raida_go [-version]

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

Error Codes:
Note: These error codes are the same for all commands
```bash
1: Incorrect Usage
2: Could not get serial number from IP Address
3: ID, Could not open file
4: ID, Could not read file
5: ID, Corupted PNG
6. ID, CloudCoin not found in PNG ID
7: ID, PNG CRC32 Incorrect
8: ID, Failed to parse CloudCoin in PNG
9: ID, Invalid CloudCoin format
10: Random Number generation failed
11: The ID Coin not found
12: The ID directory could not not be read from
13: Change Method not found
14: Show Change failed
15: Break-in-bank change making failed
16: Insufficient funds to make the transfer
17: Results from RAIDA were out of sync
18: Invalid amount specified 
19: Invalide Skywallet Address
20: Show Coins failed
21: Could not pick coins after showing
22: Cloud not pick coins after change
23: Failed to encode JSON 
24: Transfering Coins Failed
25: Invalid Receipt ID
26: Invalid Skywallet Owner

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
