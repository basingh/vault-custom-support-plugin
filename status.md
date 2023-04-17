A custom sercret engine needs to be developed in following stages:

Define Backend -> Define configuration -> Define role -> Implementationg of renew/revoke of secret -> Define credentials


22/02
* Create main.go which is really just a plugin entry point
* The main func serves 2 main functionalities. 1. To setup gRPC communication with Vault 2. Refer to secret engine backend as a form of factory

23/02
* Setup backend.go which contains logic to setup secret engine in vault to store data. Following are the few thing its meant to do:
    * Factory function to setup backend in vault. 
    * Created backend object to secret engine
    * backend function is core of this and essentially  calling different frameworks in Vault SDKs which will used by secret engine. 
    * reset method to lock the backend while target API client object is reset
    * invalidate method which call reset to reset configuration

* For learning purposes we are going to use trimmed down version of https://developer.hashicorp.com/vault/tutorials/custom-secrets-engine/custom-secrets-engine-backend#vaultplugintlsprovider. 

Next up:

* Understanding Factory Design Pattern better
* Implement functions in backend.go to send hardcoded data to Vault. https://github.com/hashicorp/vault/tree/main/sdk/plugin/mock has some good examples.


