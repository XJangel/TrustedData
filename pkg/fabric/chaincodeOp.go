package fabric

import (
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
)

func (c *Client) InvokeCC(peers []string) (fab.TransactionID, error) {
	// new channel request for invoke
	args := packArgs([]string{"a", "b", "10"})
	req := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "invoke",
		Args:        args,
	}

	// send request and handle response
	// peers is needed
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.cc.Execute(req, reqPeers)
	log.Printf("Invoke chaincode response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)
	if err != nil {
		return "", errors.WithMessage(err, "invoke chaincode error")
	}

	return resp.TransactionID, nil
}

func (c *Client) InvokeCCDelete(peers []string) (fab.TransactionID, error) {
	log.Println("Invoke delete")
	// new channel request for invoke
	args := packArgs([]string{"c"})
	req := channel.Request{
		ChaincodeID: c.CCID,
		Fcn:         "delete",
		Args:        args,
	}

	// send request and handle response
	// peers is needed
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.cc.Execute(req, reqPeers)
	log.Printf("Invoke chaincode delete response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)
	if err != nil {
		return "", errors.WithMessage(err, "invoke chaincode error")
	}

	return resp.TransactionID, nil
}

func (c *Client) QueryCC(peer, keys string) error {
	// new channel request for query
	req := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "query",
		Args:        packArgs([]string{keys}),
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer)
	resp, err := c.cc.Query(req, reqPeers)
	if err != nil {
		return errors.WithMessage(err, "query chaincode error")
	}

	log.Printf("Query chaincode tx response:\ntx: %s\nresult: %v\n\n",
		resp.TransactionID,
		string(resp.Payload))
	return nil
}
func (c *Client) QueryCCInfo(v string, peer string) {

}

func (c *Client) Close() {
	c.SDK.Close()
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}
