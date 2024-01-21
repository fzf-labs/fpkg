import (
    "context"
	"encoding/json"
	"errors"

    "github.com/fzf-labs/fpkg/orm/custom"
    "github.com/fzf-labs/fpkg/orm/gen/cache"
    "{{.daoPkgPath}}"
    "{{.modelPkgPath}}"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)
