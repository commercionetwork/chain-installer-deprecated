# Commercio.network chain installer
The Commercio.network chains installer allows you to easily download, install and start a Commercio.network 
full node in matter of seconds. 

Current version: ![GitHub release](https://img.shields.io/github/release/commercionetwork/chain-installer.svg)


## Usage
### Setup
If you are operating on a Linux or MAC system, in order to run the script you must give it the execution permissions. 
To do so, run 

```bash
chmod +x commercio-network-chain-installer-{...}
``` 

### Steps
0. On UNIX-based systems, run `./commercio-network-chain-installer-{...}`; on Windows-based system start the `.exe` file.
1. Select the chain version you wish to install, by choosing from the list that is shown you and 
   fetched from the  [chains repo](https://github.com/Commercionetwork/Chains).  
   In order to select the option, use your up and down arrow key and press ENTER to select the option you prefer. 
2. Once selected the chain version, type the installation directory. 
3. Wait for the magic to happen.
4. Select whenever you wish to start your full node or not.     

## Under the hood
In order to work, the following steps are performed: 

1. The list of chains is fetched from the [chains repo](https://github.com/Commercionetwork/Chains).
2. Once selected the chain id, the proper binaries are downloaded from the associated tag from 
   the [binaries repo](https://github.com/Commercionetwork/Commercionetwork).
3. The `cnd init` command is executed to set everything up properly. 
4. The `genesis.json` file and the seed nodes are fetched from the chains repo. 
5. Everything is properly saved, unzipped and installed inside the installation directory.    
6. The `cnd start` command is run. 


## Building
Inside the [release page](https://github.com/Commercionetwork/Chain-installer/releases) you will find all the 
binaries ready to be downloaded and used. 