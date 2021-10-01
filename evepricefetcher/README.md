# evepricefetcher

evepricefetcher is a short golang script made for quickly fetching market prices from the eve esi

## Usage
### Input File
a list of typeids to pull the prices for
```bash
626
627
3756
...
```
### Running the script
```bash
go run pricefetcher.go -c 5 -i items.txt -o itemprices.json -t sell
```

### Output file
```bash
{
    "626":8707000,
    "627": 9311000,
    "3756": 36160000,
    ...
}
```
## License
[MIT](https://choosealicense.com/licenses/mit/)