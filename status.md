A custom sercret engine needs to be developed in following stages:

Define Backend -> Define configuration -> Define role -> Implementationg of renew/revoke of secret -> Define credentials


22/02
* Create main.go which is really just a plugin entry point
* The main func serves 2 main functionalities. 1. To setup gRPC communication with Vault 2. Refer to secret engine backend as a form of factory

Next:

* Setup backend.go which contains logic to setup secret engine in vault to store data. 
