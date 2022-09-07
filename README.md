# GovSSO Mock

GovSSO mock is an application that serves [GovSSO protocol](https://e-gov.github.io/GOVSSO/TechnicalSpecification) to
clients. Its main use cases are:

* Enable client applications to develop and test integration with GovSSO protocol. Compared to GovSSO service demo
  environment, mock can be used without needing registration and can also be used offline or in closed networks.
* Provide mock authentication data in development and test environments. Compared to GovSSO service demo environment,
  mock can return arbitrary user data to client and is also simpler to use with automated tests.

Currently GovSSO mock returns protocol-compliant responses for all successful flows. Validating input data and
simulating error conditions should be considered for future development.

## Documentation

* [DEPLOYMENT.md](DEPLOYMENT.md) – guide for building, running, and configuring GovSSO mock.
* [USAGE.md](USAGE.md) – guide for testing all usage flows with GovSSO mock.
* [ARCHITECTURE.md](ARCHITECTURE.md) – description of GovSSO mock internals and possible future development.
