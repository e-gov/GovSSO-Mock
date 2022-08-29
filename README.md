# GovSSO Mock

## Prerequisites for development

* Golang version 1.19+ (might also work with older versions)

## Building Dependencies

1. Follow [GOVSSO-Client/README.md](https://github.com/e-gov/GOVSSO-Client#running-in-docker) to build its Docker image
2. Generate required resources (TLS certificates, id-token signing keys)
   NB! By default, generates TLS certificates for '*.test' subdomains. Can be modified in './config/tls/generate-tls-resources.sh' script.
   ```shell
   cd ./config
   ./generate-resources.sh
   ```

## Running in Docker Compose

1. Run in docker compose
   ```shell
   docker compose up
   ```
2. Rebuild GovSSO Mock image and run in docker compose
   ```shell
   docker compose up --build
   ```

## Endpoints

* GovSSO Mock
    * https://govsso-mock.test:10443/
    * https://govsso-mock.test:10443/.well-known/openid-configuration
    * https://govsso-mock.test:10443/.well-known/jwks.json
    * https://govsso-mock.test:10443/oauth2/auth
    * https://govsso-mock.test:10443/oauth2/token
    * https://govsso-mock.test:10443/oauth2/sessions/logout
    * https://govsso-mock.test:10443/backchannel/sessions/logout
* Example Client
    * https://client.test:11443/ - UI
    * https://client.test:11443/actuator - maintenance endpoints

## Configuration

[//]: # (TODO: clients.json)
[//]: # (TODO: config.json)
[//]: # (TODO: users.json)