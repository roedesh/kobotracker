# CryptoKobo

Crypto tracker app for Kobo devices.

![MIT License](https://img.shields.io/badge/license-MIT-green "MIT License")
![Version 0.0.1](https://img.shields.io/badge/version-0.0.1-blue "Version 0.0.1")

## Installation

Download the latest release and extract the contents of `kobotracker.tar.gz` into the root of your Kobo device.
This will also add a NickelMenu entry for you.

After installing the app, you need to create the following file `.adds/kobotracker/config.ini`. It should at least have the `id` option set, which should be one or multiple CoinGecko ids separated by a space. Example:

```
ids=bitcoin ethereum cardano
```

Finding the CoinGecko ID is easy. Just go to the CoinGecko page for your crypto and copy the last part of the URL.

After you have created the configuration file, you may start up the app by choosing the "Kobotracker" entry in NickelMenu.

## Configuration
- `ids`: list of CoinGecko ids separated by a space
- `update_price_interval`: how long the app should wait before updating prices again (in seconds)
- `show_next_interval`: how long each crypto will be shown (in seconds)
- `fiat`: in which Fiat currency the price should be shown (defaults to 'eur')

**NOTE:** keep in mind that CoinGecko has a rate limit of 50 requests per minute, which should be enough most of the time. Just dont set the `update_price_interval` too short, otherwise you might hit their rate limiting. Every price update does 1 request per id you have set.


## Tested devices

- Kobo Libra H2O

## Building from scratch

**NOTE:** To build this app, you will need a Linux machine. The koxtoolchain fails to build on other platforms. You might be able to get it to build on MacOS, but it would probably be best to use a virtual machine instead.

Install [koxtoolchain]() and [Go](https://golang.org/), then run `make` to build the app.

To create a release tar file that can be used for installing the app on the Kobo, which includes assets and a NickelMenu configuration, run `make release`.