package envoy.authz

import input.attributes.request.http as http_request

default allow = false

# helper to get the token payload
token = {"payload": payload} { io.jwt.decode(input.sm_token, [_, payload, _]) }

# allow GET access to the "/claims" API for all services except the "invoice_service"
allow {
    http_request.method == "GET"
    http_request.path == "/claims"
    http_request.headers.serviceid != "invoice_service"
}

# allow GET access to the "/invoices" API to users with the "BILLING_MANAGER" role
allow {
    input.method == "GET"
    input.path == ["invoices", "opa"]
    token.payload.userRole == "BILLING_MANAGER"
}