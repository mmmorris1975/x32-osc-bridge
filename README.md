# X32 OSC Bridge

This program provides a bridge/proxy to the OSC interface of the [Behringer X32](https://www.behringer.com/product.html?modelCode=P0ASF)
family of digital audio mixer, including the Compact, Producer, and Rack models.  It may also work for the XAir family
of digital mixers as well.

As a proxy, it can enable tools like X32-Edit or the excellent [X32 Mixing Station](https://play.google.com/store/apps/details?id=com.davidgiga1993.mixingstation)
mobile app to access mixers which may reside on another network.

As a bridge, it can allow tools which speak OSC, but don't implement the Behringer extensions, to access the mixer and
receive updates using the "xremote" OSC extension used by these mixers.

### Installation

Download the compiled program for your operating system and processor architecture from the releases section and run.
There are no external libraries or language runtime dependencies.

Raspberry Pi uses should download the Linux Arm file for a given version.

**TODO** provide installation packages for the various platforms to simplify installation

### Usage

The program is designed to be simple to use, with little to no configuration needed.

```text
Usage of ./x32-osc-bridge:
  -p int client port
  -v	 enable verbose output
```

The `-p` option can be used when a connecting device requires a specific port for traffic returned from this program.
It may be necessary to use this flag with certain OSC-enabled control devices, refer to your product's documentation to
find the value needed for this option.

#### Proxy

When used as a proxy, the only thing required is to execute the program (optionally, with the `-v` option), and the
program will create the necessary network plumbing to pass traffic to the mixer.  The proxy supports the auto-discovery
process used by the X32-Edit program, so as soon as this program is running, you can use your tools the same as you always
have.

#### Bridge

Similar to the proxy, the only thing required is to execute the program.  In your client program you will need to specify
the network information needed to communicate with this program.  In your client program, configure the IP address of the
system you run this program on, and use port 10023.

To enable bi-directional communication with the client, you may need to use the `-p` option to configure the port the
x32-osc-bridge needs to use to talk back with the client.  If you use TouchOSC or QLC+ (and likely others!), you will
need to use this option to configure the necessary port so changes on the mixer are communicated back to the client.

In the future it may be possible to specify configuration which allows you to remap an OSC path on the input side to a
path which is recognized by the mixer.  This could be useful for client tools which output pre-defined OSC messages.

### Examples

Some examples for running the tool with bi-directional communication enabled for various client applications.  The port
numbers used appears to be default values in the client app, and should hopefully make integration as easy as copy & paste.

#### Bridge Mode for X32-Edit, X32 Mixing Station, and possibly X-Air
```shell
x32-osc-bridge
```

#### Touch OSC
```shell
x32-osc-bridge -p 9000
```

#### QLC+
```shell
x32-osc-bridge -p 7700
```

### Reference

X32 OSC protocol documentation https://behringerwiki.musictribe.com/index.php?title=OSC_Remote_Protocol  
Open Sound Control (OSC) documentation http://opensoundcontrol.org/introduction-osc