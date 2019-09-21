MUD Protocol
---

## Definitions
| **Command**              | **Byte Value** | **Note**                                                                                                   |
| ------------------------ | -------------- | ---------------------------------------------------------------------------------------------------------- |
| **ECHO**                 | 1              |                                                                                                            |
| **SUPPRESSGOAHEAD**      | 3              |                                                                                                            |
| **STATUS**               | 5              |                                                                                                            |
| **TIMINGMARK**           | 6              |                                                                                                            |
| **TERMINALTYPE**         | 24             | [see](https://tintin.sourceforge.io/protocols/mtts/)                                                       |
| **WINDOWSIZE**           | 31             |                                                                                                            |
| **TERMINALSPEED**        | 32             |                                                                                                            |
| **REMOTEFLOWCONTROL**    | 33             |                                                                                                            |
| **LINEMODE**             | 34             |                                                                                                            |
| **ENVIRONMENTVARIABLES** | 36             |                                                                                                            |
| AUTHENTICATION           | 37             |                                                                                                            |
| ENCRYPTION               | 38             |                                                                                                            |
| NEWENVIRONMENT           | 39             |                                                                                                            |
| TN3270E                  | 40             |                                                                                                            |
| XAUTH                    | 41             |                                                                                                            |
| CHARSET                  | 42             |                                                                                                            |
| RSP                      | 43             |                                                                                                            |
| COMPORTCONTROL           | 44             |                                                                                                            |
| TELNETSUPPRESSLOCALECHO  | 45             |                                                                                                            |
| TELNETSTARTTLS           | 46             |                                                                                                            |
| KERMIT                   | 47             |                                                                                                            |
| SENDURL                  | 48             |                                                                                                            |
| FORWARDX                 | 49             |                                                                                                            |
| MSDP                     | 69             | [see](https://tintin.sourceforge.io/protocols/msdp/)                                                       |
| MSSP                     | 70             | [see](https://tintin.sourceforge.io/protocols/mssp/)                                                       |
| MCCP1                    | 85             | [see](https://www.zuggsoft.com/zmud/mcp.htm)                                                               |
| MCCP2                    | 86             | [see](https://tintin.sourceforge.io/protocols/mccp/)                                                       |
| MSP                      | 90             | [see](https://www.zuggsoft.com/zmud/msp.htm)                                                               |
| MXP                      | 91             | [see](https://www.zuggsoft.com/zmud/mxp.htm)                                                               |
| ZMP                      | 93             | [see](http://discworld.starturtle.net/external/protocols/zmp.html)                                         |
| TELOPTPRAGMALOGON        | 138            |                                                                                                            |
| TELOPTSSPILOGON          | 139            |                                                                                                            |
| TELOPTPRAGMAHEARTBEAT    | 140            |                                                                                                            |
| ATCP                     | 200            | [see](http://www.ironrealms.com/rapture/manual/files/FeatATCP-txt.html)                                    |
| GMCP                     | 201            | [see@tintin](https://tintin.sourceforge.io/protocols/gmcp/) & [see@gammon](https://www.gammon.com.au/gmcp) |
| **SE**                   | 240            |                                                                                                            |
| **NOP**                  | 241            |                                                                                                            |
| **DATAMARK**             | 242            |                                                                                                            |
| **BREAK**                | 243            |                                                                                                            |
| **INTERRUPT**            | 244            |                                                                                                            |
| **ABORTOUTPUT**          | 245            |                                                                                                            |
| **AREYOUTHERE**          | 246            |                                                                                                            |
| **ERASECHAR**            | 247            |                                                                                                            |
| **ERASELINE**            | 248            |                                                                                                            |
| **GOAHEAD**              | 249            |                                                                                                            |
| **SB**                   | 250            |                                                                                                            |
| **WILL**                 | 251            |                                                                                                            |
| **WONT**                 | 252            |                                                                                                            |
| **DO**                   | 253            |                                                                                                            |
| **DONT**                 | 254            |                                                                                                            |
| **IAC**                  | 255            |                                                                                                            |

## Examples

#### pkuxkx.net

``` go
255 253 24  // DO TERMINALTYPE
255 253 31  // DO WINDOWSIZE
255 253 91  // DO ??
255 251 70  // WILL MSSP
255 251 93  // WILL ??
255 253 39  // DO NEWENVIRONMENT
255 251 201 // WILL GMCP
255 251 1   // WILL ECHO
```
