# fsconvert

fsconvert is a collecion of functions to convert from and
to Go's filesystem interface, defined in io/fs.

## FromZigbee

FromZigbee creates a filesystem with the content of a Zigbee network.

## FromKNX

FromKNX creates a filesystem with the values sent to KNX Group Addresses.

## ToFUSE

ToFUSE exposes a filesystem via FUSE, to be able to mount it.

## FromMQTT

FromMQTT creates a filesystem with the values published to a MQTT server.

## ToMQTT

ToMQTT publishes the content of a filesystem to a MQTT server.

## KNXtoMQTT

KNXtoMQTT listens to datagrams from a KNX network and publishes them
to a MQTT server.

## PrintTree

PrintTree displays a representation of a filesystem.
