# shapeshift-influx

Reads current data from the shapeshift API and loads it into influxDB.

Recent transactions are stored as:
```
transaction,currency_in=<type>,currency_out=<type> amount=<amount> <timestamp>
```

Conversion rates are stored as:
```
market_info,pair=<pair_type> rate=<rate>,limit=<limit>,min=<min>,miner_fee=<minerFee> <timestamp>
```

## Usage

```bash
# shapeshift-influx <influx url> <dbname> <type1> <type2> [typeN...]
```

### Example

The following example will add recent transactions, as well as current rates for `btc_ltc` `btc_doge` `ltc_btc` `ltc_doge` `doge_btc` and `doge_ltc`:

```bash

shapeshift-influx http://foo:bar@localhost:8086 shapeshiftdata btc ltc doge
```



