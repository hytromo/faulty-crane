# How do I run this?

Add the GOPATH's bin to PATH so you are able to execute things directly through there
```bash
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

Then you need to run `go install` at the specific path (or pointing to the specific file) to build and move the app inside the bin directory. See `tasks.json` for more information.
