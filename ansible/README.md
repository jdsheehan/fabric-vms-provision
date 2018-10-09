## About

Most Hyperledger Fabric examples focus on provisioning it with `docker-compose`, which means if the machine on which docker-compose is running goes down, the whole fabric network goes with it. This example shows how to provision fabric on multiple virtual machines (vms), with each fabric component (each peer, orderer, fabric-ca) residing on its own virtual machine (vm).

**Note:**
- **This is not a production ready example.**
- **It will also incur the financial costs of running softlayer vms.**

The example leverages `ansible` and `softlayer` to provision a network.

<p align="center">
  <img src="../orgs.svg">
</p>


## Quickstart

### Prerequisites:
- softlayer account
- softlayer api key
- an ssh key registered with softlayer
- vagrant
- virtualbox

### Initial steps

```
vagrant up
vagrant ssh
```

On the vagrant vm, populate `~/.softlayer`. See the softlayer python api read-the-docs `Configuration File` section for details.
```
[softlayer]
username = 
api_key = 
endpoint_url = 
timeout = 
```

On the vagrant vm, create an ssh public/private keypair (called `~/.ssh/fabric`) and copy the public key to softlayer. See `Devices->Manage->SSH Keys` in softlayer.

### In general

Open a terminal
```
vagrant up
vagrant ssh
cd /vagrant/ansible

slcli sshkey list
export SLKEYID='REPLACE WITH SSH KEY ID'

# create password for vms
# http://docs.ansible.com/ansible/latest/reference_appendices/faq.html#how-do-i-generate-crypted-passwords-for-the-user-module
mkpasswd --method=sha-512
export SLPASS='REPLACE ME'

eval `ssh-agent`

./provision.sh <number of peers in org0> <number of peers in org1> ... <number of peers in org N>
# for example to create to create 3 orgs, org0 with 2 peers, org1 with 1 peer, and org2 with 3 peers
# ./provision.sh 2 1 3

# list vms
slcli vs list --columns id,hostname,primary_ip,backend_ip,datacenter,action,power_state

# there is also a utility to list ssh connection details for each vm
python3 utils/env_hosts.py
```

### Notes
It might be worthwhile adding 
```
export SLKEYID='REPLACE WITH SSH KEY ID'
export SLPASS='REPLACE ME'
eval `ssh-agent` 
```
to the vagrant vm `/home/vagrant/.bashrc`

**This creates a lot of vms**

For example, with `./provision.sh 2 1 3`, each org has a seperate vm for zookeeper, kafka, fabric-ca, orderer, cli, and a seperate vm for each peer. In the example `./provision.sh 2 1 3` this equates to **21 vms** _[1(zk) + 1(kafka) + 1(fca) + 1(orderer) + 1(cli) + 2(peer)] + [1(zk) + 1(kafka) + 1(fca) + 1(orderer) + 1(cli) + 1(peer)] + [1(zk) + 1(kafka) + 1(fca) + 1(orderer) + 1(cli) + 3(peer)]_

|           |  org0 |  org1 |  org2 |  total |
| :-------- | ----: | ----: | ----: | -----: |
| zk        |     1 |     1 |     1 |        |
| kafka     |     1 |     1 |     1 |        |
| fabric-ca |     1 |     1 |     1 |        |
| orderer   |     1 |     1 |     1 |        |
| cli       |     1 |     1 |     1 |        |
| peer      |     2 |     1 |     3 |        |
|           |     7 |     6 |     8 | **21** |

## Next steps

Once the network is up, a channel can be added and chaincode installed. A simple key/value chaincode is included.

Open a new terminal and connect to `cli0`

```
vagrant ssh
ssh -i ~/.ssh/fabric fabric@$CLI0
cd pkg-cli0

# this contains a set of scripts (which leverage the peer binary) to create, join, and update a channel as well as install and instantiate a chaincode on the channel
./01-create
./02-join
./03-update
./04-install
./05-instantiate
```

Open a new terminal and connect to `cli1`
```
vagrant ssh
ssh -i ~/.ssh/fabric fabric@$CLI1
cd pkg-cli1

# this contains a set of scripts (which leverage the peer binary) to fetch, join, and update a channel as well as install a chaincode on the channel
./01-fetch
./02-join
./03-update
./04-install
```

Open a new terminal and connect to `cli2`
```
vagrant ssh
ssh -i ~/.ssh/fabric fabric@$CLI2
cd pkg-cli2

# this contains a set of scripts (which leverage the peer binary) to fetch, join, and update a channel as well as install a chaincode on the channel
./01-fetch
./02-join
./03-update
./04-install
```

Now that the channel is created and chaincode instantiated, the `./invoke` and `./query` scripts can be used to append data to the ledger.

### Notes
- connecting to the `peerXorgY` vm and running `docker logs -f <chaincode-container>` will show chaincode being executed.
- there is also a sample program on each cli, in `~/pkg-invoke` which uses the unofficial `gohfc` library to perform invokes
- on the build vm `~/build/fabric/orgX` shows the files required for an org


## Tear down

Open a terminal
```
vargrant ssh
cd /vagrant/ansible
ansible-playbook --key-file "~/.ssh/fabric"  cancel.yml
```
