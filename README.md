# lproto

This package implements length prefix framing to send and receive data over a reader/writer.
Each message is prefixed with a 4 byte header to signify its length.

A reader and writer have been provided:

- FrameWriter
- FrameReader
