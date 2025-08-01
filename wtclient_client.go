package lndclient

import (
	"context"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc/wtclientrpc"
	"google.golang.org/grpc"
)

type WtClientClient interface {
	ServiceClient[wtclientrpc.WatchtowerClientClient]

	AddTower(ctx context.Context, pubkey []byte, address string) error

	DeactivateTower(ctx context.Context, pubkey []byte) (string, error)

	GetTowerInfo(ctx context.Context, pubkey []byte, includeSessions,
		excludeExhaustedSessions bool) (*wtclientrpc.Tower, error)

	ListTowers(ctx context.Context, includeSessions,
		excludeExhaustedSessions bool) ([]*wtclientrpc.Tower, error)

	Policy(ctx context.Context, policyType wtclientrpc.PolicyType) (*PolicyResponse, error)

	RemoveTower(ctx context.Context, pubkey []byte, address string) error

	Stats(ctx context.Context) (*StatsResponse, error)

	TerminateSession(ctx context.Context, sessionId []byte) (string, error)
}

type PolicyResponse struct {
	maxUpdates       uint32
	sweepSatPerVbyte uint32
}

type StatsResponse struct {
	numBackups           uint32
	numPendingBackups    uint32
	numFailedBackups     uint32
	numSessionsAcquired  uint32
	numSessionsExhausted uint32
}

type wtClientClient struct {
	client      wtclientrpc.WatchtowerClientClient
	wtClientMac serializedMacaroon
	timeout     time.Duration
}

// A compile time check to ensure that  wtClientClient implements the
// WtclientClient interface.
var _ WtClientClient = (*wtClientClient)(nil)

func newWtClientClient(conn grpc.ClientConnInterface,
	wtClientMac serializedMacaroon, timeout time.Duration) *wtClientClient {

	return &wtClientClient{
		client:      wtclientrpc.NewWatchtowerClientClient(conn),
		wtClientMac: wtClientMac,
		timeout:     timeout,
	}
}

// RawClientWithMacAuth returns a context with the proper macaroon
// authentication, the default RPC timeout, and the raw client.
func (m *wtClientClient) RawClientWithMacAuth(
	parentCtx context.Context) (context.Context, time.Duration,
	wtclientrpc.WatchtowerClientClient) {

	return m.wtClientMac.WithMacaroonAuth(parentCtx), m.timeout, m.client
}

func (m *wtClientClient) AddTower(ctx context.Context, pubkey []byte, address string) error {
	rpcCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	rpcReq := &wtclientrpc.AddTowerRequest{
		Pubkey:  pubkey,
		Address: address,
	}

	rpcCtx = m.wtClientMac.WithMacaroonAuth(rpcCtx)
	_, err := m.client.AddTower(rpcCtx, rpcReq)
	if err != nil {
		return err
	}

	return nil
}

func (m *wtClientClient) DeactivateTower(ctx context.Context, pubkey []byte) (string, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	rpcCtx = m.wtClientMac.WithMacaroonAuth(rpcCtx)
	resp, err := m.client.DeactivateTower(rpcCtx, &wtclientrpc.DeactivateTowerRequest{
		Pubkey: pubkey,
	})
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

func (m *wtClientClient) GetTowerInfo(ctx context.Context, pubkey []byte, includeSessions,
	excludeExhaustedSessions bool) (*wtclientrpc.Tower, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	rpcCtx = m.wtClientMac.WithMacaroonAuth(rpcCtx)
	resp, err := m.client.GetTowerInfo(rpcCtx, &wtclientrpc.GetTowerInfoRequest{
		Pubkey:                   pubkey,
		IncludeSessions:          includeSessions,
		ExcludeExhaustedSessions: excludeExhaustedSessions,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *wtClientClient) ListTowers(ctx context.Context, includeSessions,
	excludeExhaustedSessions bool) ([]*wtclientrpc.Tower, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	rpcCtx = m.wtClientMac.WithMacaroonAuth(rpcCtx)
	resp, err := m.client.ListTowers(rpcCtx, &wtclientrpc.ListTowersRequest{
		IncludeSessions:          includeSessions,
		ExcludeExhaustedSessions: excludeExhaustedSessions,
	})
	if err != nil {
		return nil, err
	}

	return resp.Towers, nil
}

func (m *wtClientClient) Policy(ctx context.Context, policyType wtclientrpc.PolicyType) (*PolicyResponse, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	rpcCtx = m.wtClientMac.WithMacaroonAuth(rpcCtx)
	resp, err := m.client.Policy(rpcCtx, &wtclientrpc.PolicyRequest{
		PolicyType: policyType,
	})
	if err != nil {
		return nil, err
	}

	return &PolicyResponse{
		maxUpdates:       resp.MaxUpdates,
		sweepSatPerVbyte: resp.SweepSatPerVbyte,
	}, nil
}
func (m *wtClientClient) RemoveTower(ctx context.Context, pubkey []byte, address string) error {
	rpcCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	rpcCtx = m.wtClientMac.WithMacaroonAuth(rpcCtx)
	_, err := m.client.RemoveTower(rpcCtx, &wtclientrpc.RemoveTowerRequest{
		Pubkey:  pubkey,
		Address: address,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *wtClientClient) Stats(ctx context.Context) (*StatsResponse, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	rpcCtx = m.wtClientMac.WithMacaroonAuth(rpcCtx)
	resp, err := m.client.Stats(rpcCtx, &wtclientrpc.StatsRequest{})
	if err != nil {
		return nil, err
	}

	return &StatsResponse{
		numBackups:           resp.NumBackups,
		numPendingBackups:    resp.NumPendingBackups,
		numFailedBackups:     resp.NumFailedBackups,
		numSessionsAcquired:  resp.NumSessionsAcquired,
		numSessionsExhausted: resp.NumSessionsExhausted,
	}, nil
}

func (m *wtClientClient) TerminateSession(ctx context.Context, sessionId []byte) (string, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	rpcCtx = m.wtClientMac.WithMacaroonAuth(rpcCtx)
	resp, err := m.client.TerminateSession(rpcCtx, &wtclientrpc.TerminateSessionRequest{
		SessionId: sessionId,
	})
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}
