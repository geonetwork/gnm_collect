h1. Summary

A standalone application/daemon for recording the metrics output from Geonetwork and generating tables and graphs from the recorded data.

The application have 2 APIs for interacting with it. CLI and Web
h2. Terminal Mode

By default the system is ran in terminal mode.  In this case `?` followed by Enter/Return will show the available commands.  For example:

    q - Write reports to disk and exit the application
    f/s - Write/Flush reports to disk
    
h2. Web API

The -port=... parameter can be used to set the web server port.  By default the server is started on port 10100.  

The main page is:  http://localhost:10100/index.html (or just http://localhost:10100/)