package supervisordreceiver

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.uber.org/multierr"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		cfg         *Config
		expectedErr error
	}{
		{
			desc: "missing username, password, and invalid endpoint",
			cfg: &Config{
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingSvUsername,
				errMissingSvPassword,
				fmt.Errorf("%s: %w", errInvalidEndpoint, errors.New(`parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`)),
			),
		},
		{
			desc: "missing password and invalid endpoint",
			cfg: &Config{
				SvUsername: "otelu",
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingSvPassword,
				fmt.Errorf("%s: %w", errInvalidEndpoint, errors.New(`parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`)),
			),
		},
		{
			desc: "missing username and invalid endpoint",
			cfg: &Config{
				SvPassword: "otelp",
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingUsername,
				fmt.Errorf("%s: %w", errInvalidEndpoint, errors.New(`parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`)),
			),
		},
		{
			desc: "invalid endpoint",
			cfg: &Config{
				SvUsername: "otelu",
				SvPassword: "otelp",
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				fmt.Errorf("%s: %w", errInvalidEndpoint, errors.New(`parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`)),
			),
		},
		{
			desc: "valid config",
			cfg: &Config{
				SvUsername: "otelu",
				SvPassword: "otelp",
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: defaultEndpoint,
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualErr := tc.cfg.Validate()
			if tc.expectedErr != nil {
				require.EqualError(t, actualErr, tc.expectedErr.Error())
			} else {
				require.NoError(t, actualErr)
			}

		})
	}
}
