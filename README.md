# raidaGo

RAIDA GO Console Program

```console
Usage of ./raida_go:
./raida_go [-debug] <operation> <args>
./raida_go [-help]

<operation> is one of 'receive|send'
<args> arguments for operation

  -debug
        Display Debug Information
  -help
        Show Usage

```

Example:


```console
$ ./raida_go receive 080A4CE89126F4F1B93E4745F89F6713
{"amount_verified":150}
```

Debug Info:

```console
$ ./raida_go -debug receive 080A4CE89126F4F1B93E4745F89F6713
```
