
# Messages sent via UDP

This project is mostly meant to teach me about using Go.
The idea is that you can send and receive UDP messages to/from
your friends on the internet.

This project is quickly growing out of control...

## Networking

I'm starting with UDP, but it sounds like TCP would be a better choice.
I'll probably try both.

In order to receive remote connections, you probably have to set up
port forwarding on your local network.
This is really going to put a damper on using this with non-technical folks.
Jon gave me a really good idea of how to get around this.

## Cryptography

We use asymetric encryption.
The producer has public keys for each target in the group,
and the message is encrypted for each target using that target's key.

Messages received are decrypted with the client's private key.

### Key management

I will not support all types of keys.
You must either create a new private key using this library,
or otherwise generate an RSA key that conforms with
the [X.509 standard](https://en.wikipedia.org/wiki/X.509).
This package uses [crypto/rsa.GenerateKey](https://pkg.go.dev/crypto/rsa#GenerateKey)
to generate a private key and write it to a known location ( # todo: what's the location?).

### ToDo: Signatures
Messages are signed by the senders and verified by the listener upon receipt.
The sender will not receive a success response if any of the following occurs:

- the host associated with the message sender does not match the expected host
- the message cannot be decrypted
- the message is not signed as expected

I can include the signature in the message metadata.


## Protobuf

To compile the protobuf files run the following command:

```shell
protoc -I messages --go_out messages --go_opt=paths=source_relative messages/messages.proto
```

This package uses PB to encode your message before it's sent over to the recipient.
For now this is overkill.
If I ever start including more metadata in the messages, then this will be useful.

## ToDo

- Check max message length.
How do I chunk messages up if they exceed the buffer size?
- read public keys from a file/memory
- Think about command line interface.
I'll use [cobra](https://github.com/spf13/cobra/) and [viper](https://github.com/spf13/viper).
I'll have a manual mode and and then a mode that grabs address/username/public key data from a service
- write up details about key formats in README
- make directory where stuff is written and config lives configurable


## Reference

- [example TCP and UDP servers in Go](https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/)
- stack overflow referencing
[loading keys](https://stackoverflow.com/questions/13555085/save-and-load-crypto-rsa-privatekey-to-and-from-the-disk)
- Looks like I'll have to chunk up the payloads in order to encrypt them..
[this so answer](https://stackoverflow.com/a/67035019/14223687) seems promising.
