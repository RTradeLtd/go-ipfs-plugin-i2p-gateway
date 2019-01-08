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



Figuring out what the IPFS plugin system actually allows me to do
-----------------------------------------------------------------

The IPFS plugin system is relatively new and hasn't been used alot. It hasn't
been used for something like this at all as far as I can tell. Hopefully I can
help rectify this issue by documenting this process very methodically.

### Notes:

(1) I'm going to do that too, but that's the hard part(the plan is to adapt
BiglyBT-style bridging, with a pure-clearnet peers, clearnet-to-i2p peers, and
pure-i2p peers.

(2) Of course, that leaves the matter of i2p-compatible IPFS applications but
those are almost as simple as the gateway.
