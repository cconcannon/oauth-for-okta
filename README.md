# OAuth for Okta

## Intro

Okta encourages our users to *avoid* using tenant-level API keys to interact with the Okta API. All keys issued at the tenant level expire after 30 days.

Okta offers a better way - OAuth for Okta. Applications can be created with the ability to grant scoped access to the Okta API via the *access token* issued via OAuth2. The way this works is that the `scope` param in an authorization request includes the desired Okta API scopes (such as `okta.groups.manage` or `okta.users.manage`). If *both* the *application* and the *user* are authorized for the requested scope, then the access token is issued and can be used to interact with the Okta API.

## Demo

This demo shows how to configure an app in Okta that is a designated Okta API user. The app will not authenticate within the context of an Okta identity, but rather will authenticate against the Okta OAuth2 server as a client using a *client credentials* grant.

Okta allows an app to be created as an *API Services* application. Okta requires that the clients of this app authenticate with a *Public key / Private key* mechanism - this allows organizations to maintain their own keys that can grant access to the Okta API.

## Prerequisites

- **the relevant EA feature flag must be enabled (as of March 31, 2022)**
- Go must be installed on your local machine

## Configure an API Services App

1. configure an API Services application in Okta
2. configure scopes for access by the app
3. configure a JWK in Okta or elsewhere
4. if required for the library implementation, convert the JWK to `.pem` format (for testing purposes I used [https://www.npmjs.com/package/pem-jwk])

**[More details can be found in Okta's developer documentation](https://developer.okta.com/docs/guides/implement-oauth-for-okta-serviceapp/main/#create-and-sign-the-jwt)**

## Run this Example

1. `git clone https://github.com/cconcannon/oauth-for-okta && cd oauth-for-okta`
2. copy `example.env` to `.env` and fill the details
3. `go run main.go`

## Brainstorming

Use Cases: replace tenant API key usage

### Assumptions

- a service actor will need to perform things such as group membership audits and remediation

### Scope

- a service actor will be able to list and update Okta group membership
- the Okta syslog should show granular usage of this service tool, including application context

### Infrastructure

- web app to show client credentials authentication, JWT access_token, then using that JWT access_token to:
1. list group membership for a particular group
2. add a group member
3. remove a group member 

- show sad-path as well...
1. anything outside of group membership operations should be denied