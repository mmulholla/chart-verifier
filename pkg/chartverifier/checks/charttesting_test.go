package checks

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/cli"
)

func TestChartTesting(t *testing.T) {
	type testCase struct {
		config      map[string]interface{}
		description string
		uri         string
	}

	testCases := []testCase{
		{
			config:      map[string]interface{}{},
			description: "with chart-testing defaults",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
		{
			config: map[string]interface{}{
				"upgrade": true,
			},
			description: "override chart-testing upgrade",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
		{
			config: map[string]interface{}{
				"skipMissingValues": true,
			},
			description: "override chart-testing upgrade",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
		{
			config: map[string]interface{}{
				"namespace": "ct-test-namespace",
			},
			description: "override chart-testing namespace",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
		{
			config: map[string]interface{}{
				"releaseLabel": "chart-verifier-app.kubernetes.io/instance",
			},
			description: "override chart-testing releaseLabel",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
	}

	for _, tc := range testCases {
		config := viper.New()
		settings := cli.New()

		_ = config.MergeConfigMap(tc.config)

		t.Run(tc.description, func(t *testing.T) {
			t.Skip()
			r, err := ChartTesting(
				&CheckOptions{
					URI:             tc.uri,
					ViperConfig:     config,
					HelmEnvSettings: settings,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
		})
	}
}

type ocVersionError struct{}

func (v ocVersionError) getVersion() (string, error) {
	return "", errors.New("error")
}

type ocVersionWithoutError struct {
	Version string
}

func (v ocVersionWithoutError) getVersion() (string, error) {
	return v.Version, nil
}

type testAnnotationHolder struct {
	OpenShiftVersion              string
	CertifiedOpenShiftVersionFlag string
}

func (holder *testAnnotationHolder) SetCertifiedOpenShiftVersion(version string) {
	holder.OpenShiftVersion = version
}

func (holder *testAnnotationHolder) GetCertifiedOpenShiftVersionFlag() string {
	return holder.CertifiedOpenShiftVersionFlag
}

func TestVersionSetting(t *testing.T) {
	type testCase struct {
		description string
		holder      *testAnnotationHolder
		versioner   Versioner
		version     string
		error       string
	}

	testCases := []testCase{
		{
			description: "oc.Version returns 4.7.9",
			holder:      &testAnnotationHolder{},
			versioner:   ocVersionWithoutError{Version: "4.7.9"},
			version:     "4.7.9",
		},
		{
			description: "oc.Version returns error, flag set to 4.7.8",
			holder:      &testAnnotationHolder{CertifiedOpenShiftVersionFlag: "4.7.8"},
			versioner:   ocVersionError{},
			version:     "4.7.8",
		},
		{
			description: "oc.Version returns semantic error, flag set to fourseveneight",
			holder:      &testAnnotationHolder{CertifiedOpenShiftVersionFlag: "fourseveneight"},
			versioner:   ocVersionError{},
			error:       "OpenShift version is not following SemVer spec. Invalid Semantic Version",
		},
		{
			description: "oc.Version returns error, flag not set",
			holder:      &testAnnotationHolder{},
			versioner:   ocVersionError{},
			error:       "Missing OpenShift version. error. And the 'openshift-version' flag has not set.",
		},
	}

	for _, tc := range testCases {

		t.Run(tc.description, func(t *testing.T) {

			err := setOCVersion(&CheckOptions{AnnotationHolder: tc.holder}, tc.versioner)

			if len(tc.error) > 0 {
				require.Error(t, err)
				require.Equal(t, tc.error, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.version, tc.holder.OpenShiftVersion)
			}

		})

	}

}
