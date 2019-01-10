package i2pgateconfig

import (
	//"errors"
    "io"
    "os"
    "path/filepath"

    "github.com/ipfs/go-ipfs-util"
    "github.com/mitchellh/go-homedir"

    //config "gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config"
	serialize "gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config/serialize"
)

type Config struct {
	tunname string
	hops    int
	tunnels int
}

const (
	// DefaultPathName is the default config dir name
	DefaultPathName = ".ipfs"
	// DefaultPathRoot is the path to the default config dir location.
	DefaultPathRoot = "~/" + DefaultPathName
	// DefaultConfigFile is the filename of the configuration file
	DefaultConfigFile = "i2pconfig"
	// EnvDir is the environment variable used to change the path root.
	EnvDir = "IPFS_PATH"
)

// PathRoot returns the default configuration root directory
func PathRoot() (string, error) {
	dir := os.Getenv(EnvDir)
	var err error
	if len(dir) == 0 {
		dir, err = homedir.Expand(DefaultPathRoot)
	}
	return dir, err
}

// ConfigAt loads an i2p gateway plugin from the IPFS_PATH directory. It's a
// file intended to be as similar to the IPFS config as possible.
func ConfigAt(ipfs_path string) (*Config, error) {
	var final error
	if filename, final := Filename(ipfs_path); final == nil {
		return Load(filename)
	}
	return nil, final
}

func Path(configroot, extension string) (string, error) {
	if len(configroot) == 0 {
		dir, err := PathRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, extension), nil

	}
	return filepath.Join(configroot, extension), nil
}

func Filename(ipfs_path string) (string, error) {
	return Path(ipfs_path, "i2pconfig")
}

func Init(out io.Writer) (*Config, error) {
    cfg := &Config {
        tunname: "ipfs",
        hops: 5,
        tunnels: 5,
    }
    return cfg, nil
}

func ReadConfig(filename string, cfg interface{}) (*Config, error) {
    if err := serialize.ReadConfigFile(filename, cfg); err != nil {
        return nil, err
    }
    return cfg.(*Config), nil
}

func WriteConfig(filename string, cfg interface{}) (*Config, error) {
    if err := serialize.WriteConfigFile(filename, cfg); err != nil {
        return nil, err
    }
    return cfg.(*Config), nil
}

func Load(filename string) (*Config, error) {
	// if nothing is there, generate a 'safe(paranoid)' default config and
	// inform the user thusly
	if !util.FileExists(filename) {
        f, err := os.Create(filename)
        if err != nil {
            return nil, err
        }
		return Init(f)
	}

	var cfg Config
	return ReadConfig(filename, &cfg)
}
