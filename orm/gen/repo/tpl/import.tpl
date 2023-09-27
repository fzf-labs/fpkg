import (
	"context"
	"errors"
	"encoding/json"
	"time"
    "gorm.io/gorm"
    "github.com/fzf-labs/fpkg/cache/cachekey"
    "github.com/fzf-labs/fpkg/conv"
    "github.com/fzf-labs/fpkg/orm/gen/cache"
    "{{.daoPkgPath}}"
    "{{.modelPkgPath}}"
)
