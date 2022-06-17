[![Go Report Card](https://goreportcard.com/badge/github.com/emanuelhristea/lime)](https://goreportcard.com/report/github.com/emanuelhristea/lime)
(https://www.codefactor.io/repository/github/emanuelhristea/lime/badge)](https://www.codefactor.io/repository/github/emanuelhristea/lime) ![Docker](https://github.com/emanuelhristea/lime/workflows/Docker/badge.svg) 

<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/Icon.svg" height="70" />


## Installation 
```
$ git clone https://github.com/emanuelhristea/lime.git
```


## Setup
1. Modify config for DB in `config/config.go`
2. Update parameters for privateKey, publicKey in file `license/license.go` 
To generate new key pair use command ```go run main.go pkey```

## Run server
```
$ go run main.go server 
```

## Available Commands:
- `healthcheck` : Check healthcheck
- `help` : Help about any command
- `server` : Start license server
- `pkey` : Generating key pair


## Admin console
Link for admin console http://localhost:8080/admin/
default login - admin, password - admin
<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/admin/login.png" />
<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/admin/pricings.png" />
<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/admin/createpricing.png" />
<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/admin/customers.png" />
<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/admin/createcustomer.png" />
<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/admin/subscriptions.png" />
<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/admin/createsubscription.png" />
<img src="https://raw.githubusercontent.com/emanuelhristea/lime/master/.github/assets/admin/createlicense.png" />


## API list
* `GET      /ping ` : Health server
* `POST     /api/key` : Generate new license
* `GET      /api/subscriptions ` : Get subscriptions and licenses by email and MAC address
* `PATCH    /api/key/:customer_id` : Update license
* `POST     /api/verify` : Check status license
* `DELETE   /api/key` : Release license key.

## TODO
- [x] Generating license
- [x] Verification license
- [X] Create and install license on the client
- [x] Command-line utility for generating key pair 
- [ ] Integration with Stripe
- [x] Example client
- [x] Admin console
- [x] Support MAC address check
- [ ] Support country check
