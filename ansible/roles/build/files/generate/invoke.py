INVOKE_INTRO = '''#!/bin/bash

set -x
set -e

if [ -z "$BDIR" ]; then
  echo "env var BDIR not set. This variable contains the directory where all build artifacts are placed"
  exit 1
fi

mkdir -p $BDIR/pkg-invoke
cd $BDIR/pkg-invoke

cp $GOPATH/src/localhost/invoke/invoke .
cp $GOPATH/src/localhost/invoke/main.go .
cp -r $BDIR/artifacts/crypto-config .
'''

INVOKE_CA_INTRO = '''
cat << EOF > caOrg{org_id}.yaml
'''

INVOKE_CA_YAML = '''---
url: http://fabric-ca{org_id}.example.com:7054
skipTLSValidation: true
mspId: Org{org_id}MSP
crypto:
  family: ecdsa
  algorithm: P256-SHA256
  hash: SHA2-256
'''

INVOKE_CA_FINISH = '''
EOF
'''

INVOKE_CLIENT_INTRO = '''
cat << EOF > clientOrg{org_id}.yaml
'''

INVOKE_CLIENT = '''---
crypto:
  family: ecdsa
  algorithm: P256-SHA256
  hash: SHA2-256
'''

INVOKE_CLIENT_ORDERERS = '''
orderers:
'''

INVOKE_CLIENT_ORDERER = '''
  orderer:
    host: orderer{org_id}.example.com:7050
    useTLS: true
    tlsPath: /home/fabric/pkg-invoke/crypto-config/ordererOrganizations/example.com/orderers/orderer{org_id}.example.com/tls/server.crt
'''

INVOKE_CLIENT_PEERS = '''
peers:
'''

INVOKE_CLIENT_PEER = '''
  peer{peer_id}:
    host: peer{peer_id}.org{org_id}.example.com:7051
    useTLS: true
    tlsPath: /home/fabric/pkg-invoke/crypto-config/peerOrganizations/org{org_id}.example.com/peers/peer{peer_id}.org{org_id}.example.com/tls/server.crt
'''

INVOKE_CLIENT_EVENTS = '''
eventPeers:
'''

INVOKE_CLIENT_EVENT = '''
  peer0:
    host: peer0.org{org_id}.example.com:7053
    useTLS: true
    tlsPath: /home/fabric/pkg-invoke/crypto-config/peerOrganizations/org{org_id}.example.com/peers/peer0.org{org_id}.example.com/tls/server.crt
'''


INVOKE_CLIENT_FINISH = '''
EOF
'''

SAMPLE_INTRO = '''
cat << EOF > sample.yaml
'''

SAMPLE = '''---
routines:
'''
SAMPLE_ROUTINE = '''  - name: "peer{peer_id}Org{org_id}Routine{routine_id}"
    yamlClient: "./clientOrg{org_id}.yaml"
    yamlCA: "./caOrg{org_id}.yaml"
    peer: "peer{peer_id}"
    orderer: "orderer"
    username: "admin"
    password: "passwd"
    channel: "mychannel"
    mspid: "Org{org_id}MSP"
    prefix: "p{peer_id}o{org_id}r{routine_id}"
    msgSizeBytes: 512
    duration: 10
    chaincode: "mycc"
    version: "1.0"
'''

SAMPLE_FINISH = '''
EOF
'''
