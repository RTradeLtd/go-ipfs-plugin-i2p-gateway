go-ipfs-plugin-i2p-gateway
==========================

**WARNING:** This is *only* the gateway part, A.K.A. the easy part. It will make
your IPFS gateway accessible via i2p clients, but it **will not route**
**communication between IPFS nodes over i2p(1)**. This means that it **doesn't**
**make your IPFS instance anonymous**, it just makes it *accssible to clients*
*anonymously(2)*. As such, it's probably not that useful to people in general
yet. It's also emphatically *not* a product of the i2p Project and doesn't carry
a guarantee from them. File an issue here, I'm happy to help.

How it works
------------

First of all, it uses the un-gxed version of IPFS found here for now:
[Un-gxed IPFS](github.com/ipsn/go-ipfs).

Officially speaking, the IPFS project has only documented plugins for
filesystems. This plugin makes use of none of those interfaces. Instead it
simply takes advantage of the plugin system to set up some hidden services
when the IPFS plugin is initialized. It does this by reading the IPFS config
file to find the ports that have been configured by the admin running IPFS, then
using the SAM API to forward those ports to I2P. Once they are forwarded, a
config file is generated containing the i2p configuration and it's base32 and
base64 settings. This is stored in a file called "i2pconfig."

### Using it

Again, this plugin only shows the IPFS gateway, an

### Notes:

(1) I'm going to do that too, but that's the hard part(the plan is to adapt
BiglyBT-style bridging, with a pure-clearnet peers, clearnet-to-i2p peers, and
pure-i2p peers. That may be less straightforward than the simple description
made it sound).

(2) Of course, that leaves the matter of i2p-compatible IPFS applications but
those are almost as simple as the gateway.
