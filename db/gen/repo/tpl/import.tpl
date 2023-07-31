import (
	"context"
	"encoding/json"
	"time"
	"github.com/redis/go-redis/v9"
    "gorm.io/gorm"
    "github.com/fzf-labs/fpkg/cache/cachekey"
    "{{.FillDaoPkgPath}}"
    "{{.FillModelPkgPath}}"
)
