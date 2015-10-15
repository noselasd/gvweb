package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os/exec"
	"path"
)

//Result from running Graphviz
type graphvizResult struct {
	//name of generated file
	fileName string
	//non nil if an error occured, in which case the fileName is
	//empty/meaningless
	err error
}

func genUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}

func runGraphviz(tool, code, imgType string) graphvizResult {
	uuid, err := genUUID()
	if err != nil {
		return graphvizResult{"", err}
	}

	fileName := path.Join(g_DataDir, uuid)

	err = ioutil.WriteFile(fileName, []byte(code), 0644)
	outputFile := fileName + "." + imgType
	output, err := exec.Command(tool, "-T", imgType, "-o", outputFile, fileName).CombinedOutput()
	if err != nil {
		if len(output) > 0 { //graphviz outputted an error to us.
			return graphvizResult{"", errors.New(string(output))}
		}
		return graphvizResult{"", err}
	}

	return graphvizResult{outputFile, nil}
}
