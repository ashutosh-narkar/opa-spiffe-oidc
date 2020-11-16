package envoy.authz

import input.attributes.request.http as http_request

default allow = false

# helper to get the token payload
token = {"payload": payload} { io.jwt.decode(http_request.headers.token, [_, payload, _]) }

# allow GET access to the "/claims" API to users with the "BILLING_MANAGER" role
allow {
    http_request.method == "GET"
    http_request.path == "/claims"
    token.payload.userRole == "BILLING_MANAGER"
}
