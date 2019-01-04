# webpack

* * *

a `golang` package designed to support both serving webpack resources and handling api requests. can be used standalone or 
as an imported module. not much fluff here, very basic stuff.

why?

for fun, because i really enjoy writing `golang` backend services but also really
enjoy using `reactjs` for the frontend ui. the most obvious and straight-forward
approach, if one needed a backend for a `react` frontend, is to use `nodejs` but that's no fun.

so i built this.

## usage

simply create an instance of an `EntrypointRouter` and a `ResourceRouter` and register these to any
other router that takes the type `http.HandlerFunc`. if you don't have a preference use the included `Router`
as your mux.

for example:

```
entrypoints := webpack.NewEntrypointRouter()

entrypoints.RegisterFile("/", "../some/index.html")
entrypoints.RegisterFile("/dashboard", "../some/dashboard.html")

resources := webpack.NewResourceRouter()

resources.RegisterDirectory("../some/other/dist/")
resources.RegisterDirectory("../random/images/")

router := webpack.NewRouter()

router.GET("/$", entrypoints.ServeHTTP)
router.GET("(\\.js|\\.json|\\.css|\\.png|\\.gif|\\.jpe?g|\\.ico)$", resources.ServeHTTP)

server := &http.Server{
    Addr:    "127.0.0.1:1339",
    Handler: router,
}

log.Fatal(server.ListenAndServe())
```

see `router_test.go` for a working example (that is, if you change the constants).

from your terminal:

```
~$ git clone https://github.com/msawangwan/webpack.git
~$ cd webpack
~$ vim router_test.go
```

update the two constants, `entrypointfilepath` and `resourcedirpath` to actual filepaths on your local filesystem, save
the changes, exit and then:

```
~$ go run test -v
```

finally in your browser navigate to `localhost:1337` to see it all in action.