## Transit Vault: A simple script to sign/verify documents using HashiCorp Vault transit engine:

### This is only a POC.

Prerequisites:

* Have a transit key already configured to be used to sign/verify
* Have a HashiCorp Vault instance with transit engines enabled on `/transit` path
* Have a vault token to be used to auth and perform operations.

## To build

```
# go build transit_vault.go
```
## Usage

```
./transit_vault -cmd sign -key test123bacd -input main.go -signature main.go.signature -token s.0wtcFidsdcdscscsm -vaultaddress https://vault1:8200
```

```
./transit_vault -cmd verify -key test123bacd -input main.go -signature main.go.signature -token s.0wtcFidsdcdscscsm -vaultaddress https://vault1:8200
```

* -cmd is the operations that we need (sign/verify)
* -key is the key used to sign/verify documents
* -input is the file we want to sign/verify
* -signature is the file where signature is saved when sign OR is the file readed to verify file against the signature
* -token is the Vault token used to authenticate on vault.
* -vaultaddress is the Vault address

The code is not well written. There's some duplicate code that should be deduplicated using functions. 

