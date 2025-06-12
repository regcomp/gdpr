# Authentication and Authorization

Upon successful authentication, the auth provider will give the service data in a currently unknown shape, likely unique to each provider. The immediately known relevant details are the refresh token, access token, and permissions of the user. There will be a model conversion implementation to a model shape that the service uses. Likewise, requesting authentication may require some model conversion for each provider.

The users of this service will need to add role and permission elements to their providers flow for this service to use. They should all be configurable with the exception of admin. For example:

Roles:
- admin
- moderator
- user
- auditor

Resources:
- secrets (api keys, webhooks, ect)
- config (token expiry, run frequency, ect)
- database schema (what tables and relevant key to find all the data)
- logs
- notifications

Controls (will we need anything beyond CRUD?):
- edit
- create
- delete
- view



### In-Memory Tokens

Pros:
- Will integrate into all auth pipelines

Cons:
- Vulnerable to XSS attacks
- Requires a lot of client-side JS for security!!

### Cookies

Pros:
- automatically sent in requests, easy to set in responses
- cannot be accessed via JS with HTTPOnly
- requires minimal client-side JS for security
- can be set to vanish or persist on the browser

Cons:
- Does not fit into the OAuth pipline
- vulerable to CSRF attacks


# Security

### CSP

A way to control at the browser level how certain aspects of code are handled. Necessary to implement
[Reference](https://developer.mozilla.org/en-US/curriculum/extensions/security-and-privacy/)

### CSRF
An actor is able to send certs from a different client/uri than they were sent to to gain access. 

### XSS
A malicious actor is able to input .js into the page and it executes giving them access to the internals of the DOM
