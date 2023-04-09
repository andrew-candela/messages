
# Messages sent via UDP

This project is mostly meant to teach me about using Go.
The idea is that you can send and receive UDP messages to/from
your friends on the internet.

## Networking

I'm starting with UDP, but it sounds like TCP would be a better choice.
I'll probably try both.

In order to receive remote connections, you probably have to set up
port forwarding on your local network.
This may not be possible depending on your location.
I may implement a workaround.

## Encryption

This project uses hybrid encryption.
The producer has public keys for each target in the group,
Every message generates a random AES encryption key.
The message content is encrypted using the AES key.
The AES key is encrypted using the recipient's RSA public key
and shared with the encrypted message in the Packet fields.

### Key management

I will not support all types of keys.
You must either create a new private key using this library,
or otherwise generate an RSA key that conforms with
the [X.509 standard](https://en.wikipedia.org/wiki/X.509).
This package uses [crypto/rsa.GenerateKey](https://pkg.go.dev/crypto/rsa#GenerateKey)
to generate a private key and write it to a known location ( # todo: what's the location?).

### Signatures

Messages are signed by the senders and verified by the listener upon receipt.
The sender will not receive a success response if any of the following occurs:

- the host associated with the message sender does not match the expected host
- the message cannot be decrypted
- the message is not signed as expected

The signature is included as a field in the Packet object.

## Protobuf

To compile the protobuf files run the following command:

```shell
protoc -I messages --go_out messages --go_opt=paths=source_relative messages/messages.proto
```

This package uses PB to encode your message before it's sent over to the recipient.

## ToDo

- Check max message length.
Making some progress on this one.
I still need to add support in the consumer for messages composed of multiple datagrams.
- Think about command line interface.
I'll use [cobra](https://github.com/spf13/cobra/) and [viper](https://github.com/spf13/viper).
I'll have a manual mode and and then a mode that grabs address/username/public key data from a service
- make directory where stuff is written and config lives configurable

## Reference
