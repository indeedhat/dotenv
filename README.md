# Dotenv
Runtime dependency free library for loading .env files into your environment

## Installation
```console
go get github.com/indeedhat/dotenv
```

## Usage
### Autoload
```go
// this will automatically load the local .env file from the current working directory into your
// environment.
// existing evnvars will be overwritten by any in the .env file (should there be any crossover)
import _ "github.com/indeedhat/dotenv/autoload"
```

### Manual loading
```go
import "github.com/indeedhat/dotevn"

func main() {
    // Load in envars while maintaining existing envars should there be any conflict
    err := dotenv.Load()
    ...

    // Optionally pass in your own list of files
    err := dotenv.Load(".env", ".env.local", ".env.test")
    ...

    // Load in envars while maintaining existing envars should there be any conflict
    // This operation will fail on any file that contains invalid syntax
    err := dotenv.LoadStrict()
    ...
}
```

### Manual overloading
```go
import "github.com/indeedhat/dotevn"

func main() {
    // The Overload variants work the same way as the Load variants above except they will override
    // any existing envars with versions in the .env files
    err := dotenv.Load()
    ...
    err := dotenv.Load(".env", ".env.local", ".env.test")
    ...
    err := dotenv.LoadStrict()
    ...
}
```

### Manual parsing
```go
import "github.com/indeedhat/dotevn"

func main() {
    // You can manually parse a file returning a list of the envars found in the file for you to 
    // manually handle, in this case nothing will actually be loaded into the environment
    parser, err := dotenv.ParseFile()

    list := parser.Parse()
    // or
    list, err := parser.ParseStrict()
}
```

### Helper Types
There are a number of helper types for handling environment variables along with type conversions
within your application. They are provided for String, Int, Float and Bool values:

```.env
MY_STRING_ENVAR="my string"
MY_INT_ENVAR=1234 # or Hex 0xFF or Octal 0o775 or Binary 0b00110
MY_FLOAT_ENVAR_=5.2 # or Exponent 0e6
MY_BOOL_ENVAR=true # or 0/1 or y/n
MY_EMPTY_ENVAR=
```

```go
import "github.com/indeedhat/dotevn"

const (
    envString dotenv.String = "MY_STRING_ENVAR"
    envInt dotenv.Int = "MY_INT_ENVAR"
    envFloat dotenv.Float = "MY_FLOAT_ENVAR"
    envBool dotenv.Bool = "MY_BOOL_ENVAR"
    envEmpty dotenv.String = "MY_EMPTY_ENVAR"
    envMissing dotenv.String = "MY_MISSING_ENVAR"
)

func main() {
    // Get methods on the types will return the value of the variable defaulting back to the
    // optional fallback value if the variable is not found or contains an empty string
    envString.Get("fallback") // "my string"
    envInt.Get(4321) // 1234
    envFloat.Get(12.3) // 5.2
    envBool.Get(false) // true
    envEmpty.Get("fallback") // "fallback"
    envMissing.Get("fallback") // "fallback"

    // On the other hand the Lookup methods will only defalt to the fallback value if the envar
    // does not exist
    envString.Lookup("fallback") // "my string"
    envInt.Lookup(4321) // 1234
    envFloat.Lookup(12.3) // 5.2
    envBool.Lookup(false) // true
    envEmpty.Lookup("fallback") // ""
    envMissing.Lookup("fallback") // "fallback"
}
```

## Is it fast?
I haven't done any benchmarking against other similar libraries because i dont feel that speed is 
all that important when it comes to a library like this that will likely only be ran once at startup.

The test suite runs on my mid range laptop in 0.003 seconds though so i would say its fast enough
for its use case

## Its not really dependency free tho is it?!
Well no i do use [stretchr/testify](github.com/stretchr/testify) in the tests but the runtime library
does not have any external dependencies

## Is it spec complient?
I'm not sure if there is actually an official spec for .env files, looking at various different
libraries across languages there seems to be a number of defferences it is however capable of parsing
what seems to be a common syntax between the libraries i have looked at.

I may consider adding a full syntax guide but alternatively you can just look at the fixtures files
to see the syntax that is being tested against
