# SSS

Amazon Simple Storage Service (S3) for HTTP Session handling.

It uses RAM to store/retrieve data, only when you call Load or Save does
it send data across the network. This should provide good performance for the
average HTTP session use case.

However, if your storing large amounts of data then this may not be performant
for your use case.


Example:
```go
package main

import (
	"context"
	"fmt"
	"github.com/kohirens/www/session"
	"github.com/kohirens/www/session/sss"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//  to store session data.
sessionHandler := sss.NewStorageClient(bucket, context.Background())
// set where to store the session in the bucket.
sessionHandler.Prefix(sessionPrefix)
// HTTP Session handler using RAM and then saving to  Amazon S3 for longer-term.
sessionManager := session.NewManager(sessionHandler, sessionTimeout)
```