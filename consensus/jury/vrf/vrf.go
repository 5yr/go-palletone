package vrf

type Vrf interface {
	VrfProve(priKey interface{}, msg []byte) (proof, selData []byte, err error)
	VrfVerify(pubKey, msg, proof []byte) (verify bool, selData []byte, err error)
}
