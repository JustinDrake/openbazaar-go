package addr

import (
	"bytes"

	ma "gx/ipfs/QmSWLfmj5frN9xVLMMN846dMDriy5wN5jeghUm7aTW3DAG/go-multiaddr"
	manet "gx/ipfs/QmVCNGTyD4EkvNYaAp253uMQ9Rjsjy2oGMvcdJJUoVRfja/go-multiaddr-net"
	mafmt "gx/ipfs/QmYjJnSTfXWhYL2cV1xFphPqjqowJqH7ZKLA1As8QrPHbn/mafmt"
)

func isFDCostlyTransport(a ma.Multiaddr) bool {
	return mafmt.TCP.Matches(a)
}

type AddrList []ma.Multiaddr

func (al AddrList) Len() int {
	return len(al)
}

func (al AddrList) Swap(i, j int) {
	al[i], al[j] = al[j], al[i]
}

func (al AddrList) Less(i, j int) bool {
	a := al[i]
	b := al[j]

	// dial localhost addresses next, they should fail immediately
	lba := manet.IsIPLoopback(a)
	lbb := manet.IsIPLoopback(b)
	if lba {
		if !lbb {
			return true
		}
	}

	// dial utp and similar 'non-fd-consuming' addresses first
	fda := isFDCostlyTransport(a)
	fdb := isFDCostlyTransport(b)
	if !fda {
		if fdb {
			return true
		}

		// if neither consume fd's, assume equal ordering
		return false
	}

	// if 'b' doesnt take a file descriptor
	if !fdb {
		return false
	}

	// if 'b' is loopback and both take file descriptors
	if lbb {
		return false
	}

	// for the rest, just sort by bytes
	return bytes.Compare(a.Bytes(), b.Bytes()) > 0
}
