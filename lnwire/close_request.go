package lnwire

import (
	"fmt"

	"github.com/roasbeef/btcd/btcec"
	"github.com/roasbeef/btcd/wire"
	"github.com/roasbeef/btcutil"

	"io"
)

// CloseRequest is sent by either side in order to initiate the cooperative
// closure of a channel. This message is rather sparse as both side implicitly
// know to craft a transaction sending the settled funds of both parties to the
// final delivery addresses negotiated during the funding workflow.
//
// NOTE: The requester is able to only send a signature to initiate the
// cooperative channel closure as all transactions are assembled observing
// BIP 69 which defines a cannonical ordering for input/outputs. Therefore,
// both sides are able to arrive at an identical closure transaction as they
// know the order of the inputs/outputs.
type CloseRequest struct {
	// ChannelPoint serves to identify which channel is to be closed.
	ChannelPoint *wire.OutPoint

	// RequesterCloseSig is the signature of the requester for the fully
	// assembled closing transaction.
	RequesterCloseSig *btcec.Signature

	// Fee is the required fee-per-KB the closing transaction must have.
	// It is recommended that a "sufficient" fee be paid in order to
	// achieve timely channel closure.
	// TODO(roasbeef): if initiator always pays fees, then no longer needed.
	Fee btcutil.Amount
}

// NewCloseRequest creates a new CloseRequest.
func NewCloseRequest(cp *wire.OutPoint, sig *btcec.Signature) *CloseRequest {
	// TODO(roasbeef): update once fees aren't hardcoded
	return &CloseRequest{
		ChannelPoint:      cp,
		RequesterCloseSig: sig,
	}
}

// A compile time check to ensure CloseRequest implements the lnwire.Message
// interface.
var _ Message = (*CloseRequest)(nil)

// Decode deserializes a serialized CloseRequest stored in the passed io.Reader
// observing the specified protocol version.
//
// This is part of the lnwire.Message interface.
func (c *CloseRequest) Decode(r io.Reader, pver uint32) error {
	// ChannelPoint (8)
	// RequesterCloseSig (73)
	// 	First byte length then sig
	// Fee (8)
	err := readElements(r,
		&c.ChannelPoint,
		&c.RequesterCloseSig,
		&c.Fee)
	if err != nil {
		return err
	}

	return nil
}

// Encode serializes the target CloseRequest into the passed io.Writer observing
// the protocol version specified.
//
// This is part of the lnwire.Message interface.
func (c *CloseRequest) Encode(w io.Writer, pver uint32) error {
	// ChannelID
	// RequesterCloseSig
	// Fee
	err := writeElements(w,
		c.ChannelPoint,
		c.RequesterCloseSig,
		c.Fee)
	if err != nil {
		return err
	}

	return nil
}

// Command returns the integer uniquely identifying this message type on the
// wire.
//
// This is part of the lnwire.Message interface.
func (c *CloseRequest) Command() uint32 {
	return CmdCloseRequest
}

// MaxPayloadLength returns the maximum allowed payload size for this message
// observing the specified protocol version.
//
// This is part of the lnwire.Message interface.
func (c *CloseRequest) MaxPayloadLength(pver uint32) uint32 {
	// 36 + 73 + 8
	return 117
}

// Validate performs any necessary sanity checks to ensure all fields present
// on the CloseRequest are valid.
//
// This is part of the lnwire.Message interface.
func (c *CloseRequest) Validate() error {
	// Fee must be greater than 0.
	if c.Fee < 0 {
		return fmt.Errorf("Fee must be greater than zero.")
	}

	// We're good!
	return nil
}
