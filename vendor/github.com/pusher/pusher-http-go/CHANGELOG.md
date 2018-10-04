1.3.0 / 2018-08-13
==================

* This release adds support for end to end encrypted channels, a new feature for Channels. Read more [in our docs](https://pusher.com/docs/client_api_guide/client_encrypted_channels).

1.2.0 / 2016-05-24
==================

* Add support for batch events

1.1.0 / 2016-02-22
==================

* Introduce a `Cluster` option for the Pusher initializer.

1.0.0 / 2015-05-14
==================

* Users can pass in a `http.Client` instance to the Pusher initializer. They can configure this instance directly to have specific options e.g. timeouts.
* Therefore, the `Timeout` field on `pusher.Client` is deprecated.
* `HttpClient()` function is no longer public. HTTP Client configuration is now done on the `HttpClient` **property** of `pusher.Client`. Read [here](https://github.com/pusher/pusher-http-go#request-timeouts) for more details.
* If no `HttpClient` is specified, the library creates one with a default timeout of 5 seconds.
* The library is now GAE compatible. Read [here](https://github.com/pusher/pusher-http-go#google-app-engine) for more details.

0.2.2 / 2015-05-12
==================

* Socket_ids are now validated upon Trigger*Exclusive and channel authentication.

0.2.1 / 2015-04-30
==================

* Webhook validation uses hmac.Equals to guard against timing attacks.

0.2.0 / 2015-03-30
==================

* A HTTP client is shared between requests to allow configuration. If none is set by the user, the library supplies a default. Allows for pipelining or to change the transport.

0.1.0 / 2015-03-26
==================

* Instantiation of client from credentials, URL or environment variables.
* User can trigger Pusher events on single channels, multiple channels, and exclude recipients
* Authentication of private and presence channels
* Pusher webhook validation
* Querying application state
* Cluster configuration, HTTPS support, timeout configuration.
