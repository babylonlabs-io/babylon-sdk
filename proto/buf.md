# Protobufs

This is the public protocol buffers API for [Babylon SDK](https://github.com/babylonlabs-io/babylon-sdk).

## Download

The `buf` CLI comes with an export command. Use `buf export -h` for details

#### Examples:

Download Babylon SDK protos for a commit:
```bash
## todo: not published, yet
buf export buf.build/babylonchain/babylon-sdk:${commit} --output ./tmp
```

Download all project protos:
```bash
buf export . --output ./tmp
```