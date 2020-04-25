go-garlic-tcp-transport
=======================

A new, more consistent and useful libp2p transport for tcp-like messages over
i2p.

Relationship to sam3
--------------------

This is essentially a shim between libp2p and sam3 which prepares all of the
libp2p-specific parts on top of the sam3 Streaming connection and listener
interfaces.
