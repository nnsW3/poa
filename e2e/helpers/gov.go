package helpers

import (
	"context"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

// Modified from ictest
func VoteOnProposalAllValidators(ctx context.Context, c *cosmos.CosmosChain, proposalID string, vote string) error {
	var eg errgroup.Group
	valKey := "validator"
	for _, n := range c.Nodes() {
		if n.Validator {
			n := n
			eg.Go(func() error {
				// gas-adjustment was using 1.3 default instead of the setup's 2.0+ for some reason.
				// return n.VoteOnProposal(ctx, valKey, proposalID, vote)

				_, err := n.ExecTx(ctx, valKey,
					"gov", "vote",
					proposalID, vote, "--gas", "auto", "--gas-adjustment", "2.0",
				)
				return err
			})
		}
	}
	return eg.Wait()
}

func ValidatorVote(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, proposalID string, searchHeightDelta uint64) {
	err := VoteOnProposalAllValidators(ctx, chain, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to vote on proposal")

	height, err := chain.Height(ctx)
	require.NoError(t, err, "failed to get height")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+searchHeightDelta, proposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")
}
