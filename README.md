# CryptoKobo

Crypto tracker app for Kobo devices.

![MIT License](https://img.shields.io/badge/license-MIT-green "MIT License")
![Version 0.0.1](https://img.shields.io/badge/version-0.0.1-blue "Version 0.0.1")

## Installation

Download the latest release. Create a directory called `cryptokobo` inside your `.adds` directory. Extract the tar file to this directory.

Then create a `cryptokobo` file inside `/.adds/nm` with the following content:

```
menu_item:main:CryptoKobo:cmd_spawn:quiet:exec /mnt/onboard/.adds/cryptokobo/run.sh
```

## Tested devices

- Kobo Libra H2O

## Building from scratch

**NOTE:** To build this app, you will need a Linux machine. The koxtoolchain fails to build on other platforms. You might be able to get it to build on MacOS, but it would probably be best to use a virtual machine instead.

Install [koxtoolchain]() and [Go](https://golang.org/), then run `make` to build the app.