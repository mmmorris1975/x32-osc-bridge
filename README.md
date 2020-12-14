# X32 OSC Bridge

This program provides a bridge/proxy to the OSC interface of the [Behringer X32](https://www.behringer.com/product.html?modelCode=P0ASF)
family of digital audio mixer, including the Compact, Producer, and Rack models.

As a proxy, it can enable tools like X32-Edit or the excellent [X32 Mixing Station](https://play.google.com/store/apps/details?id=com.davidgiga1993.mixingstation)
mobile app to access mixers which may reside on another network.

As a bridge, it can allow tools which speak OSC, but don't implement the Behringer extensions, to access the mixer and
receive updates using the "xremote" OSC extension used by these mixers.

### Usage

The program is designed to be simple to use, with little to no configuration needed.  A `-v` command option is provided
to enable verbose output, if needed for troubleshooting, or if you are interested in seeing what the program is doing.

#### Proxy

When used as a proxy, the only thing required is to execute the program (optionally, with the `-v` option), and the
program will create the necessary network plumbing to pass traffic to the mixer.  The proxy supports the auto-discovery
process used by the X32-Edit program, so as soon as this program is running, you can use your tools the same as you always
have.

#### Bridge

Similar to the proxy, the only thing required is to execute the program.  In your client program you will need to specify
the network information needed to communicate with this program.  In your client program, configure the IP address of the
system you run this program on, and use port 10023.

In the future it may be possible to specify configuration which allows you to remap an OSC path on the input side to a
path which is recognized by the mixer.  This could be useful for client tools which output pre-defined OSC messages.

### Reference

Details about the X32 OSC protocol can be found at https://behringerwiki.musictribe.com/index.php?title=OSC_Remote_Protocol