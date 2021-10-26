### station_trading
my first attempt at using shared c objects to execute golang code in python.
gathers market data relevant for station trading (volume, margins, sell2buy etc.)

### Compile golang code
```
go build -buildmode=c-shared -o Libraries/esiRequests.dll esiRequests.go
```
### Run the script
```
> python station_trading.py
```