package config

import (
	"io"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func v2ParseConfig(rd io.Reader) (Config, error) {
	dec := yaml.NewDecoder(rd)
	dec.KnownFields(true)
	var conf Config
	if err := dec.Decode(&conf); err != nil {
		return conf, err
	}
	if conf.Version == "" {
		return conf, ErrMissingVersion
	}
	if conf.Version != "2" {
		return conf, ErrUnknownVersion
	}
	if len(conf.SQL) == 0 {
		return conf, ErrNoPackages
	}
	// if err := conf.validateGlobalOverrides(); err != nil {
	// 	return conf, err
	// }

	for j := range conf.SQL {

		if conf.SQL[j].Gen.Go != nil {
			if conf.SQL[j].Gen.Go.Out == "" {
				return conf, ErrNoPackagePath
			}
			if conf.SQL[j].Gen.Go.Package == "" {
				conf.SQL[j].Gen.Go.Package = filepath.Base(conf.SQL[j].Gen.Go.Out)
			}

		}

		for _, cg := range conf.SQL[j].Codegen {
			if cg.Plugin == "" {
				return conf, ErrPluginNoName
			}
			if cg.Out == "" {
				return conf, ErrNoOutPath
			}

		}
	}
	return conf, nil
}
