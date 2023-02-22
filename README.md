# vault-custom-support-plugin :construction:

Learn path to developing custom plugin for vault using Go

This an attempt to learn Go programming language by building custom plugin to interact with Vault and perform operations. 
Following are phases identified so far with the help of my mentor Glenn Sarti:

1. Build a simple go plugin to create a secret engine path in Vault and pass hardcoded values as part of configuration. 
2. For starters an api call to this secret engine will just return a random gernerated number as a secret. 
3. Extend above functionality in [1] above to read values from a json file and pass it to configure secret engine. 
4. Next up, send an API request to dummy service e.g. (https://github.com/hashicorp/tf-testing-mocks), read response and use data in response to pass it to Vault to configure secret engine. 


#### Further ideas

* Create CI/CD build pipelines around this custom plugin and release. 
* Create different functionalities within this plugin to hack or break vault usable to reproduce tough support issues. 
* Explore topics around Go-routines and GRPC etc.?


#### Implementation:

Hashicorp has a good tutorial around how to use Vault SDK to create a customer secret engine here: https://developer.hashicorp.com/vault/tutorials/custom-secrets-engine. We would use similar approach to build this plugin.

The current state of implementation can be find in `status.md`  