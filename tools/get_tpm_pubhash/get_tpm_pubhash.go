/*
 ** Copyright 2019 Bloomberg Finance L.P.
 **
 ** Licensed under the Apache License, Version 2.0 (the "License");
 ** you may not use this file except in compliance with the License.
 ** You may obtain a copy of the License at
 **
 **     http://www.apache.org/licenses/LICENSE-2.0
 **
 ** Unless required by applicable law or agreed to in writing, software
 ** distributed under the License is distributed on an "AS IS" BASIS,
 ** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 ** See the License for the specific language governing permissions and
 ** limitations under the License.
 */

package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/bloomberg/spire-tpm-plugin/pkg/common"
	"github.com/google/certificate-transparency-go/x509"
	"github.com/google/go-attestation/attest"
)

func main() {
	tpm, err := attest.OpenTPM(&attest.OpenConfig{
		TPMVersion: attest.TPMVersion20,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer tpm.Close()

	tpmPubHash, err := getTpmPubHash(tpm)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(tpmPubHash)
}

func getTpmPubHash(tpm *attest.TPM) (string, error) {
	eks, err := tpm.EKs()
	if err != nil {
		return "", err
	}

	var ekCert *x509.Certificate
	for _, ek := range eks {
		if ek.Certificate != nil && ek.Certificate.PublicKeyAlgorithm == x509.RSA {
			ekCert = ek.Certificate
			break
		}
	}
	if ekCert == nil {
		return "", errors.New("could not find RSA public key")
	}

	hashEncoded, err := common.GetPubHash(ekCert)
	if err != nil {
		return "", err
	}

	return hashEncoded, nil
}
