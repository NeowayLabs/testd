# testd

It is a library that makes easier to write integration tests that depends on runnings daemons/services.

If you are thinking about unit tests, this library will not help you with that.
They are a good idea on a lot of cases, and it is not on scope of this document to
argue about when it is better to write integration tests, instead of unit tests.

When you found yourself needing to write integration tests, then you will need some stuff:

* A way to start services
* A way to stop them
* A way to get the logs from the services

Well, if you value deterministic tests, isolation will be your friend. One good way to get isolation is
to start/stop a instance of the service on setup/teardown.

Another thing that is important is the logs, if a test fails, why it failed ? On integration tests this
is a little harder to get only by the last performed change (it is actually one of the drawbacks of integration
tests, in comparison with unit tests).

**testd** provides an easy way to do that, and just that :-).


## What would be the logs ?

The logs are actually everything that the daemon sends to stdout and stderr. This fits
well the same model used by docker, where services just send logs to stdout and the
docker daemon (or some other daemon, like journald) captures it.

With **testd** you will be able to chose where to save the logs.


## How to use it ?

[Checkout the godocs documentation](TODO).
