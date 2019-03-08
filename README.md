go-ipfs-plugin-i2p-gateway
==========================

Plugin for presenting an IPFS gateway over i2p.

**WARNING:** This is *only* the gateway part, A.K.A. the easy part. It will make
your IPFS gateway accessible via i2p clients, but it **will not route**
**communication between IPFS nodes over i2p(1)**. This means that it **doesn't**
**make your IPFS instance anonymous**, it just makes it *accssible to clients*
*anonymously(2)*. As such, it's probably not that useful to people in general
yet. It's also emphatically *not* a product of the i2p Project and doesn't carry
a guarantee from them. File an issue here, I'm happy to help.

Also, even if your IP address and/or location are obfuscated, your IPFS identity
is unique. You can serve information to anonymous users, but you have an
identity which corresponds to the IPFS node that you are using(No different to
the key fingerprint of an SSH server). If it isn't possible to associate it to
a real-life identity or organization, and can't be correlated with a real
physical location, then it could be regarded as "pseudonymous."

*Todo:* For right now the plugin exposes the HTTP and RPC API's only. It's
possible to talk directly to the coreAPI now, and soon, it will do so.

How it works
------------

It simply takes advantage of the plugin system to set up some hidden services
when the IPFS plugin is initialized. It does this by reading the IPFS config
file to find the ports that have been configured by the admin running IPFS, then
using the SAM API to forward those ports to I2P. Once they are forwarded, a
config file is generated containing the i2p configuration and it's base32 and
base64 addresses. This is stored in a file called "i2pconfig." Finally, it
simply forwards the IPFS gateway to I2P.

Installation And Setup
---------

The gotcha to the IPFS plugin system is that you'll need to run an IPFS node
that is using the same version of IPFS which was used to build the plugin.
Due to `gx` being extensively used by the IPFS project, and it being hard to
work with, you will not be able to use it with any IPFS daemon that wasn't built
from the code within the `vendor/github.com/ipfs/go-ipfs` folder of this
repository. This is somewhat undesirable, but it is the only way to
reliably run this plugin. As of IPFS 0.4.19, we use the latest tagged release of
IPFS. Even so, the easiest way to use this plugin is to build against the
dependencies vendored in this repository. To re-generate the vendor directory,
run `make vendor`.

To build the IPFS version needed to use this plugin run `make ipfs` which will
copy the built `ipfs` binary to `$GOPATH/bin`. After that, the usual IPFS setup
is needed which won't be covered here.

Once you have all the dependencies vendored in, you can run:

        make build

To build the plugin.

Installing and Using it
-----------------------

Assuming you have IPFS_PATH set, you can simply:

        make install

Again, this plugin only shows the IPFS gateway, it does not make you anonymous.
Once you have it in your plugin directory, restart the IPFS daemon and the
plugin should load. Retrieve the base32 address from the config file and visit
the page in ~5-30 minutes.

### Notes:

(1) I'm going to do that too, but that's the hard part(the plan is to adapt
BiglyBT-style bridging, with a pure-clearnet peers, clearnet-to-i2p peers, and
pure-i2p peers. That may be less straightforward than the simple description
made it sound).

(2) Of course, that leaves the matter of i2p-compatible IPFS applications but
those are almost as simple as the gateway.

### Other plugins:

There are some other i2p plugins that are companions to this plugin. For now,
if you install one of them, you end up getting them all. So if you've followed
the install procedure for go-ipfs-plugin-i2p-gateway(this plugin), you've also
got go-ipfs-plugin-i2p-bootstrap.
