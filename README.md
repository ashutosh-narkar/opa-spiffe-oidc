# opa-spiffe-oidc

OPA-SPIFFE OIDC Demo.

## Overview

This repository contains the code for the OPA-SPIFFE OIDC Demo.

## Running the Demo

### Step 1: Install Docker

Ensure that you have recent versions of `docker` and `docker-compose` installed.

### Step 2: Build

Build the binaries for the `backend` and `invoice` service.

```bash
$ ./build.sh
```

### Step 3: Start containers

```bash
$ docker-compose up --build -d
$ docker-compose ps

          Name                         Command               State                 Ports
------------------------------------------------------------------------------------------------------
opa-spiffe-oidc_backend_1   /bin/sh -c /usr/local/bin/ ...   Up      10000/tcp, 0.0.0.0:8001->8001/tcp
opa-spiffe-oidc_invoice_1   /bin/sh -c /usr/local/bin/ ...   Up      0.0.0.0:5000->5000/tcp
opa-spiffe-oidc_opa_1       ./opa_envoy_linux_amd64 ru ...   Up      0.0.0.0:9191->9191/tcp
```

### Step 4: Start the `frontend` service

Open a new terminal window to start the `frontend` service.

```bash
$ cd src/frontend-app
$ FLASK_APP=main.py flask run -p 8000
```

### Step 4: Exercise the demo

In a browser open [http://localhost:8000](http://localhost:8000) and login as `alice@hooli.com` as the username and `Opastyra5%` as the password.