package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/x/margin/client/cli"
	"github.com/elys-network/elys/x/margin/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestShowMTPByPool(t *testing.T) {
	net, objs := networkWithMTPObjects(t, 2)

	ctx := net.Validators[0].ClientCtx

	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc      string
		ammPoolId uint64

		args []string
		err  error
		obj  *types.MTP
	}{
		{
			desc:      "found",
			ammPoolId: objs[0].AmmPoolId,

			args: common,
			obj:  objs[0],
		},
		{
			desc:      "not found",
			ammPoolId: (uint64)(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				fmt.Sprintf("%d", tc.ammPoolId),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdGetPositionsByPool(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.PositionsByPoolResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Mtps), 2)
				require.Subset(t,
					objs,
					resp.Mtps,
				)
			}
		})
	}
}