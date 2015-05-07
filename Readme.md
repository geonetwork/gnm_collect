h1. Summary

A standalone application/daemon for recording the metrics output from Geonetwork and generating tables and graphs from the recorded data.

The application can be ran on the commandline and interacted with via terminal-like commands or it can be started with a http server and 
interaction can be done via http requests.

h2. Terminal Mode

By default the system is ran in terminal mode.  In this case `?` followed by Enter/Return will show the available commands.  For example:

    q - Write reports to disk and exit the application
    f - Write/Flush reports to disk