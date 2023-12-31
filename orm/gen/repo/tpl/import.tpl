import (
    "context"
	"encoding/json"
	"errors"

    "github.com/fzf-labs/fpkg/orm/gen/cache"
    "{{.daoPkgPath}}"
    "{{.modelPkgPath}}"
    "github.com/fzf-labs/fpkg/orm/paginator"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)
