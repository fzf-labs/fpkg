import (
	"context"
	"encoding/json"
	"time"
	"github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "github.com/fzf-labs/fpkg/cache/cachekey"
    "{{.mod}}/{{.relativePath}}/{{.lowerDbName}}_dao"
    "{{.mod}}/{{.relativePath}}/{{.lowerDbName}}_model"
)
