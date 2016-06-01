# Additional information

## Decision making

Initially I thought about using `libchan` or `gRPC` as framework to communicate
server/client offering the best performance. However I used in the past `libchan`
(https://github.com/hectorj2f/search_networking) and I learn of their problems
when handling multiple concurrent clients. I also thought about using `gRPC` but
I realized we just want to use plain text to send the messages. That could help to
allow heterogeneous types of clients. Therefore I decided to build my own socket
server/client implementation.

I combined goroutines/channels to be able to compute concurrent requests without
having to use many system resources when handling the incoming connections.

I used cobra framework to build the command line interface because it makes things
to like nice and easy out of the box.

I hesitated to use two different data structures. Initially I just wanted to be as
much efficient as possible, so I used the data model defined in `package.go`.
But I thought we also want to provide a more tree-based data model, that would allow us
to add additional features for rendering. This new data model was used in `packageV2.go`.

## Future ideas

I could add more tests and better documentation for each part of the application.

I plan to use `gRPC` for the client/server that would improve the performance, as
well as, helps us to transfer the information via JSON.
