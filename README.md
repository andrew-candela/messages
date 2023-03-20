
# A not so good chat app

I never had AIM when I was a kid.
I think that's what I have in mind with this project.
Ultimately I want to be able to chat with my friends in our terminals
without sending data to Apple or Google.
I'll also get to learn about a bunch of CS fundamentals on the way:

## Networking

I'm starting with UDP, but it sounds like TCP would be a better choice.
I'll probably try both.

In order to receive remote connections, you will probably have to set up
port forwarding on your local network.

## Cryptography

I haven't quite settled on a design here, but I'm imagining a web service
where each user can put their public keys in a "group", and then
when you send a message, your machine uses the public key of the recipient to
encrypt the message.
I haven't yet figured out how to make IP addresses available to senders.

## Protobuf

To compile the protobuf files run the following command:

```shell
protoc -I messages --go_out messages --go_opt=paths=source_relative messages/messages.proto
```

## Reference

- [example TCP and UDP servers in Go](https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/)
