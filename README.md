
# UDPM: Messages sent via UDP

This project is mostly meant to teach me about using Go.
The idea is that you can send and receive UDP messages to/from
your friends on the internet.

The core of this app is peer to peer messaging, but there will
be an optional service that power users can host that will aid
user discovery and resolve some networking issues for folks who
cannot set up port forwarding on their local networks.

## Networking

I'm starting with UDP, but it sounds like TCP would be a better choice.
I'd have to come up with a different name if I use TCP.

In order to receive remote connections, you probably have to set up
port forwarding on your local network.
This may not be possible depending on your location.

The workaround here is to host your own chat server.

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
or otherwise generate an RSA key on your own.
This package uses [crypto/rsa.GenerateKey](https://pkg.go.dev/crypto/rsa#GenerateKey)
to generate a private key and write it to your UDPM_HOME dir,
(~/.udpm/udpm_id_rsa by default).

### Signatures

Messages are signed by the senders and verified by the listener upon receipt.
The sender will not receive a success response if any of the following occur:

- the host associated with the message sender does not match the expected host
- the message cannot be decrypted
- the message is not signed as expected

## Webserver

There is an option to host a webserver, so that you can chat with people who
cannot set up port forwarding on their networks.
Currently the config file of the webmaster (the user hosting the webserver)
must have the correct IP address and public RSA key of each user in the group.
With this config the webserver will be able to:

- tell clients what their own (public, if connecting over the internet) IP address is (`/ip`)
- give clients the config for all the users in a group (`/config`)
- authenticate all requests with the public key assigned to the client's IP.

## Protobuf

To compile the protobuf files run the following command:

```shell
protoc -I messages --go_out messages --go_opt=paths=source_relative messages/messages.proto
```

This package uses PB to encode your message before it's sent over to the recipient.

## ToDo

- set up the websocket mechanism.
  Looks like I'll have to maintain a map of subscribers along with the config.
  - messages will be sent to all subscribers who are also in config
  - the server will have to tell the clients which messages were delivered
- the clients will have to be able to send correct headers (with a signed token)

## Reference

- [cobra](https://github.com/spf13/cobra/) for CLI arguments
- [viper](https://github.com/spf13/viper) for config.
- [X.509 standard](https://en.wikipedia.org/wiki/X.509) defines the format of public key certificates