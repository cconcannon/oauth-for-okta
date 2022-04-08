# OAuth for Okta

## Intro

Okta encourages using [OAuth for Okta](https://developer.okta.com/docs/guides/implement-oauth-for-okta-serviceapp/main/) when building services that will interact with the Okta API. 

### Benefits of OAuth for Okta

- bring-your-own PKI infrastructure (or use Okta's)
- application context is captured in syslogs
- scope management at the application level

OAuth for Okta works by using the `scope` param in an authorization request. The param must include the desired Okta API scopes (such as `okta.groups.manage` or `okta.users.manage`). If the application is authorized for the requested scope(s), then an access token is issued via the *client credentials* grant. The token can then be used to interact with the Okta API.

## Demo

This demo shows how to configure an API-only application for use with the Okta API (machine-to-machine, no user context).

Okta allows an app to be created as an *API Services* application. Okta requires that the clients of this app authenticate with a *Public key / Private key* mechanism. If desired, organizations can maintain their own PKI infrastructure, and Okta will fetch and cache the public keys. Okta can also issue and manage keys.

**Note: Okta issues keys in `.json` format. If your app library requires a `.pem` format to sign requests (such as the Okta Golang SDK used here), then you can use a tool such as [pem-jwk](https://github.com/dannycoates/pem-jwk) to perform the conversions. Do not use external libraries for production keys without first performing a thorough audit of the library.**

## Prerequisites to run this demo

- **the relevant EA feature flag must be enabled for your Okta tenant (as of March 31, 2022)**
- [Go v1.17+](https://go.dev/doc/install) must be installed on your local machine

## Configure an API Services App in the Okta Admin Console

1. configure an API Services application in Okta
2. configure scopes for access by the app (this demo uses `okta.groups.manage`)
3. configure public/private keys in Okta (in the app settings) or BYO
4. copy the private key, clientId, and org url for use in the next steps

## Configure and Run the App

1. convert the private key to a file named `private.pem` [according to the certificate spec](https://www.rfc-editor.org/rfc/rfc1422) (for testing purposes I used [https://github.com/dannycoates/pem-jwk])
2. create a `secrets` folder in the project root directory and copy `private.pem` to this folder
3. copy `example.env` to `.env` and add your org/app details
4. `go run main.go`

**More details can be found in [OAuth for Okta guide](https://developer.okta.com/docs/guides/implement-oauth-for-okta-serviceapp/main/) and in the [Okta Go SDK documentation](https://github.com/okta/okta-sdk-golang/).**
