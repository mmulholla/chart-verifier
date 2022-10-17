package tool

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

func TestPGPKeyEncoding(t *testing.T) {

	keyfileName := "../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.key"
	expectedDigest := "1cc31121e86388fad29e4cc6fc6660f102f43d8c52ce5f7d54e134c3cb94adc2"

	encodedKey, encodeErr := GetEncodedKey(keyfileName)
	require.NoError(t, encodeErr)
	require.True(t, len(encodedKey) > 0)

	keyDigest, digestErr := GetPublicKeyDigest(encodedKey)
	require.NoError(t, digestErr)
	require.Equal(t, expectedDigest, keyDigest)

	decodedKey, decodeErr := GetDecodedKey(encodedKey)
	require.NoError(t, decodeErr)
	require.True(t, len(decodedKey) > 0)

	keyBytes, readErr := ioutil.ReadFile(keyfileName)
	require.NoError(t, readErr)
	require.Equal(t, keyBytes, decodedKey)

	baseCmd := fmt.Sprintf("cat %s | base64 | sha256sum", keyfileName)
	out, baseCmdErr := exec.Command("bash", "-c", baseCmd).Output()
	require.NoError(t, baseCmdErr)

	shaResponseSplit := strings.Split(string(out), " ")
	require.Equal(t, expectedDigest, shaResponseSplit[0])

	//getShaCmd := fmt.Sprintf("%s | base64 | sha256sum", keyfileName)
	// cmdErr := exec.Command("base64", "-i", keyfileName, "-o", "./base64key.txt").Run()
	// require.NoError(t, cmdErr, fmt.Sprintf("Error: %v", cmdErr))

	//shaResponse, shaCmdErr := exec.Command("sha256sum", "./base64key.txt").Output()
	//removeErr := os.Remove("base64key.txt")
	//require.NoError(t, removeErr)
	//require.NoError(t, shaCmdErr, fmt.Sprintf("Error: %v : %s", shaCmdErr, string(shaResponse)))
	//shaResponseSplit := strings.Split(string(shaResponse), " ")
	//require.Equal(t, keyDigest, strings.TrimRight(shaResponseSplit[0], " -\n"))
}

func TestDigest(t *testing.T) {
	keyfileName := "../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.key"
	expectedDigest := "1cc31121e86388fad29e4cc6fc6660f102f43d8c52ce5f7d54e134c3cb94adc2"

	base64Cmd := exec.Command("base64", "-i", keyfileName)
	sha256cmd := exec.Command("sha256sum")

	sha256Value := bytes.NewBufferString("")
	sha256cmd.Stdin, _ = base64Cmd.StdoutPipe()
	sha256cmd.Stdout = sha256Value

	_ = sha256cmd.Start()
	_ = base64Cmd.Run()
	_ = sha256cmd.Wait

	t.Logf("sha256: %s", sha256Value.String())
	shaResponseSplit := strings.Split(sha256Value.String(), " ")
	require.Equal(t, expectedDigest, strings.TrimRight(shaResponseSplit[0], " -\n"))

}
