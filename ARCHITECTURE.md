# Architecture

GovSSO mock is implemented in Go programming language, which was chosen for the following reasons:

* Go ecosystem has well-maintained libraries and frameworks to implement GovSSO mock functionality.
* Go program's small memory footprint and fast startup time enable convenient deployment of GovSSO mock in different
  development and test environments, some of which may have significant resource constraints.
* [Go's increasing popularity](https://survey.stackoverflow.co/2022/#technology-most-popular-technologies) helps to
  attract developers to amend GovSSO mock codebase in the future.

GovSSO mock makes a simplification compared to real GovSSO service: mock stores working state only in memory, which is
lost with a mock restart. Therefore, if the mock is restarted when a client application flow is in progress, the flow
breaks and needs to be started over. GovSSO mock doesn't depend on any external services nor doesn't write to any files
on disk, it logs messages to the console by default.

GovSSO mock is built and tested with Go 1.19 (which was the latest version as of this writing in 2022-09-06), although
it may work with older versions as newer language features may not yet have been used in the mock codebase. GovSSO mock
is built on the Gin web framework and uses the following libraries for implementing the main functionality:

* go-jose – for creating and parsing JWT-s;
* uuid – for generating UUID values for `sid` and `jti` claims;
* jwx – for producing JWKS endpoint output;
* zerolog – for logging;
* validator - for validating requests;
* Bootstrap – for having a consistent UI style.

GovSSO example client setup in Docker Compose, key pair and TLS certificate generation scripts are based on open-source
[GovSSO-Session repository](https://github.com/e-gov/GovSSO-Session). Docker Compose setup currently uses `*.localhost`
domain names for default deployment, because these are supported out-of-the-box by major browsers without needing to
modify the `hosts` file on the user's computer.

GovSSO mock serves only HTTPS and doesn't have an option to serve plain HTTP, because this corresponds to real GovSSO
service and forces integrating client applications to correctly use TLS from the start.

GovSSO mock serves additional HTTP response headers analogous to real GovSSO service:

* Headers disabling caching in the browser, to prevent intermediate or old states from being reused during certain
  navigation scenarios.
* Header preventing GovSSO mock page being included from an iframe, to enable client application developers to
  experience similar behavior to real GovSSO service and discover integration mistakes more quickly.

## Considerations for future development

* Validate more conditions on input data.
* Simulate more error conditions.
    * Issue different HTTP status codes.
    * Issue TLS certificates that don't pass different security checks.
* Enable grouping of preconfigured users and add an additional information field to preconfigured users. Useful when
  client application testing requires preconfiguring tens or hundreds of users which currently becomes a long list in
  mock UI where a specific user is difficult or slow to find.
* Add TARA protocol support to this mock, allowing it to be switched with a configuration option. This would eliminate
  the need to maintain two separate mock codebases.
* Add more logging to simplify debugging of client application interaction with GovSSO mock.
* Strictly enforce ID Token request (OAuth authorization code usage) timeout to better conform to real GovSSO service
  behavior.
* Remember active sessions and allow to select from these in the back-channel logout testing screen.
* Enable adding multiple JWK-s to JWKS, to simulate a key pair rollover scenario similar to when it would happen in real
  GovSSO service.
* Limit TLS protocol versions and cipher suites to those that are possible with real GovSSO service.
* Enable configuring back-channel logout request TLS certificate chain validation.
* Retry back-channel logout requests.
* Consider limiting back-channel logout request testing when the mock is deployed in a non-local environment (e.g.,
  development or test network) so that requests to arbitrary URLs could not be maliciously performed.
* Consider replacing custom key pair and TLS certificate generation scripts with some well-known solution (e.g., minica
  or mkcert).
* Consider replacing custom code for loading configuration from JSON files with a third-party library.
* Consider moving inline JavaScript to asset files and using the jQuery library.
* Cover code with unit tests.
* Clearer UX design.
