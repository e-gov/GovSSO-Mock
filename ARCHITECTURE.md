# Architecture

GovSSO mock is implemented in Go programming language, which was chosen for the following reasons:

* Go ecosystem has well-maintained libraries and frameworks to implement GovSSO mock functionality.
* Go program’s small memory footprint and fast startup time enables convenient deployment of GovSSO mock in different
  development and test environments, some of which may have significant resource constraints.
* [Go’s increasing popularity](https://survey.stackoverflow.co/2022/#technology-most-popular-technologies) helps to
  attract developers to amend GovSSO mock codebase in the future.

GovSSO mock makes a simplification compared to real GovSSO service: mock stores working state only in memory, which is
lost with mock restart. Therefore, if mock is restarted when a client application flow is in progress, the flow breaks
and needs to be started over. GovSSO mock doesn’t depend on any external services nor doesn’t write to any files on
disk, it logs messages to console by default.

GovSSO mock is built and tested with Go 1.19 (which was latest version as of this writing in 2022-09-06), although it
may work with older versions as newer language features may not yet have been used in mock codebase. GovSSO mock is
built on Gin web framework and uses the following libraries for implementing main functionality:

* jwt-go – for creating and parsing JWT-s;
* uuid – for generating UUID values for sid and jti claims;
* jwx – for producing JWKS endpoint output;
* zerolog – for logging;
* validator - for validating requests;
* Bootstrap – for having a consistent UI style.

GovSSO example client setup in Docker Compose, keypair and TLS certificate generation scripts are based on open-source
[GOVSSO-Session repository](https://github.com/e-gov/GOVSSO-Session). GOVSSO-Session repository’s Docker Compose setup
currently uses *.localhost domain names for default deployment, because these are supported out-of-the-box by major
browsers without needing to modify hosts file on user’s computer. But browsers impose slightly different semantics to
sites served from *.localhost domains, which breaks GovSSO session update flow. Therefore, current project uses *.test
domain names for default deployment, which don’t have special semantics in browsers, but need modification of hosts file
on user’s computer.

GovSSO mock serves only HTTPS and doesn’t have an option to serve plain HTTP, because this corresponds to real GovSSO
service and forces integrating client applications to correctly use TLS from the start.

GovSSO mock serves additional HTTP response headers analogous to real GovSSO service:

* Headers disabling caching in browser, to prevent intermediate or old state being reused during certain navigation
  scenarios.
* Header preventing GovSSO mock page being included from an iframe, to enable client application developers to
  experience similar behavior to real GovSSO service and discover integration mistakes more quickly.

## Considerations for future development

* Simulate error conditions.
    * Issue different HTTP status codes.
    * Issue ID Tokens and Logout Tokens that don’t pass different security checks that are listed in GovSSO protocol
      specification. Token payload could be allowed to be freely edited in GovSSO mock UI.
    * Issue TLS certificates that don’t pass different security checks.
* Support cancelling authentication flow with user_cancel error code when this is added to GovSSO protocol
  specification.
* Support phone scope and phone_number, phone_number_verified claims when this is added to GovSSO protocol
  specification.
* Optimize Docker image (use multi-stage build to reduce resulting image size).
* Enable grouping of preconfigured users and add an additional information field to preconfigured users. Useful when
  client application testing requires preconfiguring tens or hundreds of users which currently becomes a long list in
  mock UI where a specific user is difficult or slow to find).
* Add TARA protocol support to this mock, allowing it to be switched with a configuration option. This would eliminate
  the need to maintain two separate mock codebases.
* Add more logging to simplify debugging of client application interaction with GovSSO mock.
* Strictly enforce ID Token request (OAuth authorization code usage) timeout to better conform to real GovSSO service
  behavior.
* Remember active sessions and allow to select from these in back-channel logout testing screen.
* Enable adding multiple JWK-s to JWKS, to simulate key-pair rollover scenario similar to when it would happen in real
  GovSSO service.
* Set and check a cookie during authentication and session update flows, to better catch client application mistakes
  during session update background JavaScript request, as CORS configuration is difficult to get right.
* Limit TLS protocol versions and cipher suites to those that are possible with real GovSSO service.
* Enable configuring back-channel logout request TLS certificate chain validation.
* Retry back-channel logout requests.
* Consider limiting back-channel logout request testing when mock is deployed in non-local environment (e.g.,
  development or test network) so that requests to arbitrary URL-s could not be maliciously performed.
* Consider replacing custom keypair and TLS certificate generation scripts with some well-known solution (e.g., minica
  or mkcert).
* Consider replacing custom code for loading configuration from JSON file with a third-party library.
* Consider moving inline JavaScript to asset files and using jQuery library.
* Cover code with unit tests.
* Clearer UX design.
