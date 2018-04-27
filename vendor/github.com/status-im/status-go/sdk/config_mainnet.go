package sdk

import (
	"fmt"
	"log"

	"github.com/status-im/status-go/geth/params"
)

// DefaultConfig : prepares a default config
func MainnetConfig() *Config {
	configFormat := `{
    "APIModules": "db,eth,net,web3,shh,personal,admin",
    "ClusterConfig": {
        "BootNodes": [ ],
        "Enabled": true,
        "RootHash": "77eedcf6f940940b3615da49109c1ba57b95c3fff8bcf16f20ac579c3ae24e58",
        "RootNumber": 478
    },
    "DataDir": "%s/data/ethereum/mainnet_rpc",
    "DevMode": true,
    "HTTPHost": "localhost",
    "HTTPPort": 8545,
    "IPCEnabled": false,
    "IPCFile": "geth.ipc",
    "KeyStoreDir": "%s/data/keystore",
    "LightEthConfig": {},
    "ListenAddr": ":0",
    "LogEnabled": true,
    "LogFile": "",
    "LogLevel": "DEBUG",
    "LogToStderr": true,
    "MaxPeers": 25,
    "MaxPendingPeers": 0,
    "Name": "StatusIM",
    "NetworkId": 1,
    "NodeKeyFile": "",
    "RPCEnabled": false,
    "SwarmConfig": {
        "Enabled": false
    },
    "TLSEnabled": false,
    "UpstreamConfig": {
        "Enabled": true,
        "URL": "https://mainnet.infura.io/z6GCTmjdP3FETEJmMBI4"
    },
    "Version": "0.9.9-unstable",
    "WSEnabled": false,
    "WSHost": "localhost",
    "WSPort": 8546,
    "WhisperConfig": {
        "DataDir": "%s/data/ethereum/mainnet_rpc/wnode",
        "EnableMailServer": false,
        "EnablePushNotification": false,
        "Enabled": true,
        "FirebaseConfig": {
            "AuthorizationKeyFile": "",
            "NotificationTriggerURL": "https://fcm.googleapis.com/fcm/send"
        },
        "IdentityFile": "",
        "MinimumPoW": 0.001,
        "Password": "",
        "PasswordFile": "",
        "TTL": 120
    }
}`

	cwd := "/data/sg_bots/node"

	config := fmt.Sprintf(configFormat, cwd, cwd, cwd, cwd)
	cfg, err := params.LoadNodeConfig(config)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &Config{NodeConfig: cfg}
}
