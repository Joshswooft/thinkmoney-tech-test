# Checkout Kata

Implement the code for a checkout system that handles pricing schemes such as "pineapples cost 50, three pineapples cost 130."

Implement the code for a supermarket checkout that calculates the total price of a number of items. In a normal supermarket, things are identified using Stock Keeping Units, or SKUs. In our store, we’ll use individual letters of the alphabet (A, B, C, and so on). Our goods are priced individually. In addition, some items are multi-priced: buy n of them, and they’ll cost you y pence. For example, item A might cost 50 individually, but this week we have a special offer: buy three As and they’ll cost you 130. In fact the prices are:

| SKU  | Unit Price | Special Price |
| ---- | ---------- | ------------- |
| A    | 50         | 3 for 130     |
| B    | 30         | 2 for 45      |
| C    | 20         |               |
| D    | 15         |               |

The checkout accepts items in any order, so that if we scan a B, an A, and another B, we’ll recognize the two Bs and price them at 45 (for a total price so far of 95). **The pricing changes frequently, so pricing should be independent of the checkout.**

The interface to the checkout could look like:

```cs
interface ICheckout
{
    void Scan(string item);
    int GetTotalPrice();
}
```

## Prerequisites

A version of golang will be needed (version: 1.21).

## Running the application

```sh
go run main.go
```

## Building the application

This command will create a binary which can then be run from your shell.

```sh
go build -o bin/checkout
```

To run:

```sh
./bin/checkout
```


## Running tests

```sh
go test ./...
```

## Improvements

Here is a list of improvements which could be made:

### Adding contexts

By adding contexts we can free up resources e.g. a database connection or cancel long running tasks i.e. the ScanItems().

These would typically be used on the Basket and the Checkout.

### Better scanner

Currently the scanner can read from anything that implements `io.Reader` this includes strings, files, buffers, network streams, http requests etc. 

However for streaming data we don't have any control over the stream. Cancelling a long running stream could be done with the help of a context but being able to resume a stream in this implementation is not possible. It could be nice to add this option by simply following the implementation of an `io.ReaderSeeker` or using something like `bufio` so you can read the input from the last token you read.

If the product data got more complicated then you could easily use the Scanner interface to create a new scanner implementation e.g. JSON.

### Pricing rules

The pricing rules fits the task at hand but would need to be modified to do anything more advanced, it also has the drawback of not being go-routine safe (out of time).

In future you could modify it to accept pricing rules from different file formats: e.g. json, yaml, csv etc.


### Misc

- CI/CD (Github Actions)
- Pre-commit hooks
- Static analysis
- Goroutines
- Creating a CLI script with flags to input pricing rules and skus e.g. using cobra