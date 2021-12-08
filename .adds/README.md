Create the following file `config.ini` in this folder. It should at least have the `id` option set, which should be one or multiple CoinGecko ids separated by a space. Example:

```
ids=bitcoin ethereum cardano
```

Finding the CoinGecko ID is easy. Just go to the CoinGecko page for your crypto and copy the last part of the URL.

After you have created the configuration file, you may start up the app by choosing the "Kobotracker" entry in NickelMenu.

## Configuration

- `ids`: list of CoinGecko ids separated by a space
- `update_price_interval`: how long the app should wait before updating prices again (in seconds, defaults to 30)
- `show_next_interval`: how long each crypto will be shown (in seconds, defaults to 10)
- `fiat`: in which Fiat currency the price should be shown (defaults to 'eur')

**NOTE:** keep in mind that CoinGecko has a rate limit of 50 requests per minute, which should be enough most of the time. Just dont set the `update_price_interval` too short, otherwise you might hit their rate limiting. Every price update does 1 request per id you have set.