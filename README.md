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

## Technical Design

The main design decision for this task was to decouple the pricing rules from the checkout (as hinted by the challenge). This was accomplished by using a seperate `pricing` package which could then be used within the checkout object.
The checkout accepts a `PricingRules` interface which can be used to accept different types of pricing. In this application you can see that I have used a simple pricing (just multiplies the unit price by the product quantity) along
with a more advanced pricing system which handles the special pricing for the challenge.

Adding in these interfaces made unit testing a breeze as I could slide in my mocked implementations and control error states etc.


### Basket

Checkouts normally have a basket i.e. where you store your scanned items. For this challenge it would have been enough to simply use a `map[sku]quantity` on the `checkout` object and call it a day
but I wanted to show off and decouple the basket storage from the checkout by introducing another interface which allows us to have different kinds of storage i.e. in memory, from a database, even a file.

```go
type Basket interface {
	// Adds a new item or updates the existing item's quantity by its sku
	AddItem(sku sku.SKU, quantity quantity.Quantity) error

	// Gets an item by it's product SKU
	// if the item is not found then it returns a checkout.ErrItemNotFound error
	GetItem(sku sku.SKU) (qty quantity.Quantity, err error)
	// runs the iterator func over every item in the basket
	Range(iterator func(id itemID, quantity quantity.Quantity))
}
```

The interface could be improved upon i.e accepting a context so we have control over cancelling the function call. E.g. your DB might be down which means the connection could hang indefinitely, using a context with a cancel timeout
would solve this issue.


### Scanner

The scanner is definitely an add on to this challenge but it definitely expands the system in an interesting way. 

The idea behind the scanner is much like in real life where some supermarkets have scanners which you can take around the store with you, use from a website or even from an app on your phone.

The scanner I've built can amazingly accomplish that with very few lines of code by leveraging the super powers of the `io.Reader` interface.

```go
type skuScanner struct {
	reader io.Reader
}
```

By injecting the `io.Reader` we can read in anything that uses this interface which includes: buffers, strings, files, network calls (API's) and more.
But whats more interesting is that by having my own `Read(p []byte) (int, error)` function within the scanner means that the scanner itself is an `io.Reader`!.

This means you can combine this scanner with other readers e.g.`gzip.NewReader()` to read in compressed data.

Or more advanced use cases:  We can use a `io.TeeReader` to read in data from the scanner and write the data out as a backup to AWS S3.

### Skus and Quantity

For this challenge the sku is simply a single character or in golangs world a `rune`. I could have written my sku like this: `type SKU rune` however this doesn't stop me
from making a bad sku e.g. `SKU(4)`.

A common workaround for this problem can be seen in golangs std library e.g. `time.Time`.

Instead we have a `sku.New(r rune) (SKU, error)` to initialize the skus properly and catch any invalid skus. 

So why go through all this effort?

1. It simplifies error handling a LOT because anytime the `SKU` type is used we can guarantee that object is correct meaning we don't need to run validation logic in every
function.

Note: Quantity is also setup in a similar manner.

### Currency

For currency I decided to use an `int` which represents a `Penny`. This is fine for the task but in real life we would want a better data type. 

Why?

Because some grocery items are priced by weight e.g. bananas = 17.9p per 100g. We can't do this type of pricing when we are confined by an `int`.

Why not floats?

`float` is a big no no as its prone to rounding errors.


Preference would be to use a library that deals with currency e.g. £, $ etc. 

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
- Logging
- Metrics
- Mocks could be generated from library e.g. gomock, mockery etc. This way they are kept up to date with implementation
