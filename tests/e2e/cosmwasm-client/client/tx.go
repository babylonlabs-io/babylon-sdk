package client

import (
	"context"
	"fmt"
	"sync"

	"cosmossdk.io/errors"
	"github.com/avast/retry-go/v4"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/zap"

	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmwasm-client/wasmclient"
)

// ToProviderMsgs converts a list of sdk.Msg to a list of provider.RelayerMessage
func ToProviderMsgs(msgs []sdk.Msg) []wasmclient.RelayerMessage {
	relayerMsgs := []wasmclient.RelayerMessage{}
	for _, m := range msgs {
		relayerMsgs = append(relayerMsgs, wasmclient.NewCosmosMessage(m, func(signer string) {}))
	}
	return relayerMsgs
}

// SendMsgToMempool sends a message to the mempool.
// It does not wait for the messages to be included.
func (c *Client) SendMsgToMempool(ctx context.Context, msg sdk.Msg) error {
	return c.SendMsgsToMempool(ctx, []sdk.Msg{msg})
}

// SendMsgsToMempool sends a set of messages to the mempool.
// It does not wait for the messages to be included.
func (c *Client) SendMsgsToMempool(ctx context.Context, msgs []sdk.Msg) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	relayerMsgs := ToProviderMsgs(msgs)
	if err := retry.Do(func() error {
		var sendMsgErr error
		krErr := c.accessKeyWithLock(func() {
			sendMsgErr = c.provider.SendMessagesToMempool(ctx, relayerMsgs, "", ctx, []func(*wasmclient.RelayerTxResponse, error){})
		})
		if krErr != nil {
			c.logger.Error("unrecoverable err when submitting the tx, skip retrying", zap.Error(krErr))
			return retry.Unrecoverable(krErr)
		}
		return sendMsgErr
	}, retry.Context(ctx), rtyAtt, rtyDel, rtyErr, retry.OnRetry(func(n uint, err error) {
		c.logger.Debug("retrying", zap.Uint("attemp", n+1), zap.Uint("max_attempts", rtyAttNum), zap.Error(err))
	})); err != nil {
		return err
	}

	return nil
}

// SendMsg sends a message to the chain.
func (c *Client) SendMsg(ctx context.Context, msg sdk.Msg, expectedErrors []*errors.Error, unrecoverableErrors []*errors.Error) (*wasmclient.RelayerTxResponse, error) {
	return c.SendMsgs(ctx, []sdk.Msg{msg}, expectedErrors, unrecoverableErrors)
}

// SendMsgs sends a list of messages to the chain.
func (c *Client) SendMsgs(ctx context.Context, msgs []sdk.Msg, expectedErrors []*errors.Error, unrecoverableErrors []*errors.Error) (*wasmclient.RelayerTxResponse, error) {
	return c.ReliablySendMsgs(ctx, msgs, expectedErrors, unrecoverableErrors, 1)
}

// ReliablySendMsg reliable sends a message to the chain.
// It utilizes a file lock as well as a keyring lock to ensure atomic access.
// TODO: needs tests
func (c *Client) ReliablySendMsg(ctx context.Context, msg sdk.Msg, expectedErrors []*errors.Error, unrecoverableErrors []*errors.Error, retries ...uint) (*wasmclient.RelayerTxResponse, error) {
	return c.ReliablySendMsgs(ctx, []sdk.Msg{msg}, expectedErrors, unrecoverableErrors, retries...)
}

// ReliablySendMsgs reliably sends a list of messages to the chain.
// It utilizes a file lock as well as a keyring lock to ensure atomic access.
// TODO: needs tests
func (c *Client) ReliablySendMsgs(ctx context.Context, msgs []sdk.Msg, expectedErrors []*errors.Error, unrecoverableErrors []*errors.Error, retries ...uint) (*wasmclient.RelayerTxResponse, error) {
	rty := rtyAttNum
	rtyAttempts := rtyAtt
	if len(retries) > 0 {
		rty = retries[0]
		rtyAttempts = retry.Attempts(rty)
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		rlyResp     *wasmclient.RelayerTxResponse
		callbackErr error
		wg          sync.WaitGroup
	)

	callback := func(rtr *wasmclient.RelayerTxResponse, err error) {
		rlyResp = rtr
		callbackErr = err
		wg.Done()
	}

	wg.Add(1)

	// convert message type
	relayerMsgs := ToProviderMsgs(msgs)

	// TODO: consider using Babylon's retry package
	if err := retry.Do(func() error {
		var sendMsgErr error
		krErr := c.accessKeyWithLock(func() {
			sendMsgErr = c.provider.SendMessagesToMempool(ctx, relayerMsgs, "", ctx, []func(*wasmclient.RelayerTxResponse, error){callback})
		})
		if krErr != nil {
			c.logger.Error("unrecoverable err when submitting the tx, skip retrying", zap.Error(krErr))
			return retry.Unrecoverable(krErr)
		}
		if sendMsgErr != nil {
			if errorContained(sendMsgErr, unrecoverableErrors) {
				c.logger.Error("unrecoverable err when submitting the tx, skip retrying", zap.Error(sendMsgErr))
				return retry.Unrecoverable(sendMsgErr)
			}
			if errorContained(sendMsgErr, expectedErrors) {
				// this is necessary because if err is returned
				// the callback function will not be executed so
				// that the inside wg.Done will not be executed
				wg.Done()
				c.logger.Error("expected err when submitting the tx, skip retrying", zap.Error(sendMsgErr))
				return nil
			}
			return sendMsgErr
		}
		return nil
	}, retry.Context(ctx), rtyAttempts, rtyDel, rtyErr, retry.OnRetry(func(n uint, err error) {
		c.logger.Debug("retrying", zap.Uint("attemp", n+1), zap.Uint("max_attempts", rty), zap.Error(err))
	})); err != nil {
		return nil, err
	}

	wg.Wait()

	if callbackErr != nil {
		if errorContained(callbackErr, expectedErrors) {
			return nil, nil
		}
		return nil, callbackErr
	}

	if rlyResp == nil {
		// this case could happen if the error within the retry is an expected error
		return nil, nil
	}

	if rlyResp.Code != 0 {
		return rlyResp, fmt.Errorf("transaction failed with code: %d", rlyResp.Code)
	}

	return rlyResp, nil
}
