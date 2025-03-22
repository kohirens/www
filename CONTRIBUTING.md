# How To Contribute

This is a library to share functions in Go across packages. Function that
are generic enough should fit right in here. If you're not sure what that is,
take a look around, hopefully it will become apparent.

## When To Make A Sub-Package

What is meant by sub-package here is putting related functionality into
subdirectories and give it a package new name; which can then be called on its
own.

We do not want to fracture the code into too many packages, nor make one big
monolith. Try and group related functionality in the same file. When you start
using external packages then it may be time to decide if it belongs in a
subpackage or even own its on.

These things tend to work out over time. Functionality can always be moved
later. Yes, moving them into another package later on will force a semantic
major version increase. But don't worry about that, just get it done, and we'll
figure it out.

See [What to put into a package] for further assistance.

## Setup Local Development.

You can use a Docker environment to get going if you have Docker on your
computer. In fact, there is no documentation for any other way in this reading.

### Run Docker

1. Clone this repository.
2. Execute a command such as `go test -v ./...`
   ```output
   ~/src/github.com/kohirens/stdlib $ go test -v ./...
   PASS
   ok      github.com/kohirens/www/session      0.004s
   ```

---

[What to put into a package]: https://go.dev/blog/organizing-go-code#what-to-put-into-a-package