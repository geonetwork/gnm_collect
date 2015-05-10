Summary
=======

A standalone application/daemon for recording the metrics output from Geonetwork and generating tables and graphs from the recorded data.

The application have 2 APIs for interacting with it. CLI and Web

If you have the binary an example of running the application is as follows:

    gnm_collect -target http://localhost:8989/gn -logging=true -user=monitor -pass=monpas -out=gnm-reports

Command-line Interface
----------------------

By default the system is ran in terminal mode.  In this case `?` followed by Enter/Return will show the available commands.  For example:

    q - Write reports to disk and exit the application
    f/s - Write/Flush reports to disk
    
Web API
-------

The -port=... parameter can be used to set the web server port.  By default the server is started on port 10100.  

The main page is:  http://localhost:10100/index.html (or just http://localhost:10100/)

Build from Source
=================
You can use the go get commands to build from sources.  To do this you need to install:

* Go
* Git
* Hg

If this is done then you just need to do the following:

    go get github.com/geonetwork/gnm_collect
    go build github.com/geonetwork/gnm_collect
    ./bin/gnm_collect -user admin -pass admin 
 
You can run directly from the source (rather than compiled binary with:

    go run src/github.com/geonetwork/gnm_collect/gnm_collect.go
 
All dependencies are cloned and you can build from your GOPATH directory

Distribution
============

To distribute this application, clone and build the application as described in _Build from Source_.  For the distribution bundle copy
`github.com/gonum/plot/vg/fonts` along with the bin/gnm_collect binary

    