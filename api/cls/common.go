package cls

import (
	"github.com/eryajf/cloudctl/public"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func GetSecret() *common.Credential {
	return common.NewCredential(
		public.Config.TcSecretID,
		public.Config.TcSecretKey,
	)
}
