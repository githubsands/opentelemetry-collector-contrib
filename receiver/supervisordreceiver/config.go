// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package supervisordreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/supervisorddreceiver"

import (
	"errors"
	"fmt"
	"os"

	"github.com/gosimple/conf"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"gopkg.in/yaml.v2"
)

// Predefined error responses for configuration validation failures
var (
	errMissingUsername = errors.New(`"username" not specified in config`)
	errMissingPassword = errors.New(`"password" not specified in config`)
	errInvalidEndpoint = errors.New(`"endpoint" must be in the form of unix://<hostname>:<port>`)
	errNoScraperConfig = errors.New(`supervisord config file must be defined`)
)

const defaultEndpoint = "http://localhost:8098"

type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	confighttp.HTTPClientSettings           `mapstructure:",squash"`
	MetricsSettings                         MetricsSettings `mapstructure:"metrics"`

	SvConfigLocation string `mapstructure:"supervisordd_config"`
	SvPassword       string `mapstructure:"supervisordd_password"`
	SvUsername       string `mapstructure:"supervisordd_user"`
	SvUnixSocket     string `mapstructure:"supervisordd_unix_socket"`

	MetricsBuilder *MetricsBuilder
}

type SuperVisordConfig struct {
	SvUnixSocket string
	SvPassword   string
	SvUsername   string
}

const (
	supervisorUser                  = "user"
	supervisorPassword              = "password"
	supervisorctlUnixSocketLocation = "serverurl"
	supervisordConfigKey            = "supervisord"
)

func (c *Config) parseSVConfig() error {
	cfg, err := conf.ReadFile(c.SvConfigLocation)
	if err != nil {
		return err
	}
	cfg.String("default", supervisorUser)
	cfg.String("default", supervisorPassword)
	cfg.String("default", supervisorctlUnixSocketLocation)
	return nil
}

func checkUnixSocket(fn string) error {
	if fn == "" {
		return errors.New("no unix socket given")
	}
	_, err := os.Stat(fn)
	return err
}

/*
func checkSVFile(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	content, err := io.ReadAll(fd)
	if err != nil {
		return err
	}

	var targetGroups []*targetgroup.Group

	switch ext := filepath.Ext(filename); strings.ToLower(ext) {
	case ".json":
		if err := json.Unmarshal(content, &targetGroups); err != nil {
			return fmt.Errorf("error in unmarshaling json file extension: %w", err)
		}
	case ".yml", ".yaml":
		if err := yaml.UnmarshalStrict(content, &targetGroups); err != nil {
			return fmt.Errorf("error in unmarshaling yaml file extension: %w", err)
		}
	default:
		return fmt.Errorf("invalid file extension: %q", ext)
	}

	for i, tg := range targetGroups {
		if tg == nil {
			return fmt.Errorf("nil target group item found (index %d)", i)
		}
	}
	return nil
}
*/

/*
// Validate checks the receiver configuration is valid.
func (cfg *Config) Validate() error {
	supervisorddconfig := cfg.supervisorddConfig
	if supervisorddconfig != nil {
		err := cfg.validatesupervisorddconfig(supervisorddconfig)
		if err != nil {
			return err
		}
	}
	if cfg.TargetAllocator != nil {
		err := cfg.validateTargetAllocatorConfig()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cfg *Config) validateTargetAllocatorConfig() error {
	// validate targetAllocator
	targetAllocatorConfig := cfg.TargetAllocator
	if targetAllocatorConfig == nil {
		return nil
	}
	// ensure valid endpoint
	if _, err := url.ParseRequestURI(targetAllocatorConfig.Endpoint); err != nil {
		return fmt.Errorf("TargetAllocator endpoint is not valid: %s", targetAllocatorConfig.Endpoint)
	}
	// ensure valid collectorID without variables
	if targetAllocatorConfig.CollectorID == "" || strings.Contains(targetAllocatorConfig.CollectorID, "${") {
		return fmt.Errorf("CollectorID is not a valid ID")
	}

	return nil
}
*/
func newConfigError(s string) string {
	return "local supervisord.conf" + s + "must equal the collector's" + s
}

func compareConfigs(svLocalCfg, hvCollectorCfg *SuperVisordConfig) error {
	var errStrings []string = make([]string, 0)
	if svLocalCfg.SvPassword != hvCollectorCfg.SvPassword {
		errStrings = append(errStrings, newConfigError("password"))
	}
	if svLocalCfg.SvUsername != hvCollectorCfg.SvUsername {
		errStrings = append(errStrings, newConfigError("user"))
	}
	if svLocalCfg.SvUnixSocket != hvCollectorCfg.SvUnixSocket {
		errStrings = append(errStrings, newConfigError("unixSocket"))
	}
	var errString string
	for _, v := range errStrings {
		errString = errString + v + "\n"
	}
	return errors.New(errString)
}

// func (cfg *Config) ID() cfg.ComponentID {}

// Unmarshal a config.Parser into the config struct.
func (cfg *Config) Unmarshal(componentParser *confmap.Conf) error {
	if componentParser == nil {
		return nil
	}
	// We need custom unmarshaling because supervisordd "config" subkey defines its own
	// YAML unmarshaling routines so we need to do it explicitly.
	err := componentParser.UnmarshalExact(cfg)
	if err != nil {
		return fmt.Errorf("supervisord receiver failed to parse config: %w", err)
	}

	// Unmarshal supervisordd's config values. Since supervisordd uses `yaml` tags, so use `yaml`.
	svCfg, err := componentParser.Sub(supervisordConfigKey)
	if err != nil || len(svCfg.ToStringMap()) == 0 {
		return err
	}
	out, err := yaml.Marshal(svCfg.ToStringMap())
	if err != nil {
		return fmt.Errorf("supervisord receiver failed to marshal config to yaml: %w", err)
	}

	err = yaml.UnmarshalStrict(out, &svCfg)
	if err != nil {
		return fmt.Errorf("supervisord receiver failed to unmarshal yaml to supervisord config: %w", err)
	}

	return nil
}
