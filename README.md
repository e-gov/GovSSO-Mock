<img src="doc/eu_regional_development_fund_horizontal.jpg" width="350" height="200" alt="European Union European Regional Development Fund"/>

# GovSSO Mock

GovSSO mock is an application that serves [GovSSO protocol](https://e-gov.github.io/GOVSSO/TechnicalSpecification) to
clients. Its main use cases are:

* Enable client applications to develop and test integration with GovSSO protocol. Compared to GovSSO service demo
  environment (`govsso-demo.ria.ee`), mock can be used without needing registration with RIA and can also be used
  offline or in closed networks.
* Provide mock authentication data in development and test environments. Compared to GovSSO service demo environment
  (`govsso-demo.ria.ee`), mock can return arbitrary user data to client application and is also simpler to use with
  automated tests.

GovSSO mock returns protocol-compliant responses for all successful flows and allows to simulate many error conditions.
NB! GovSSO mock currently validates most conditions on input data, therefore not all validations in GovSSO mock are
currently as strict as in GovSSO environment.

## Quick start

1. Run `docker compose up`.
2. Open https://client.localhost:11443/ to test authentication, session update and logout flows;
   open https://govsso-mock.localhost:10443/ to test back-channel logout request. Refer to documentation below for more
   detailed instructions.

## Documentation

* [DEPLOYMENT.md](DEPLOYMENT.md) – guide for building, running, and configuring GovSSO mock.
* [USAGE.md](USAGE.md) – guide for testing all usage flows with GovSSO mock.
* [ARCHITECTURE.md](ARCHITECTURE.md) – description of GovSSO mock internals and possible future development.
