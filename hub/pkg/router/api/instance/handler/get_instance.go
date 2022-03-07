package instance

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetInstanceRequestUri struct {
	Instance string `uri:"instance" binding:"required"`
}

// GetInstance returns the instance information for the given authority.
//
// @Summary      Get instance information
// @Description  get instance information by authority
// @Tags         authority
// @Accept       json
// @Produce      json
// @Param        authority  path      string  true  "Authority"
// @Success      200        {object}  web.Response{data=GetInstanceResponseData}
// @Router       /{authority} [get]
func GetInstance(c *gin.Context) {
	var uri GetInstanceRequestUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, status.Error(status.ErrorInvalidURI))

		return
	}

	instance, err := rss3uri.ParseInstance(uri.Instance)
	if err != nil {
		c.JSON(http.StatusBadRequest, status.Error(status.ErrorInvalidURI))

		return
	}

	// TODO Development environment has only account data
	if instance.GetPrefix() != constants.PrefixNameAccount {
		c.JSON(http.StatusAccepted, status.Error(errors.New("not an account")))

		return
	}

	accountID := fmt.Sprintf("%s@%s", instance.GetIdentity(), instance.GetSuffix())
	account := model.Account{}

	if err := db.DB.Where("account_id = ?", accountID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, status.Error(status.ErrorAccountNotFound))

			return
		}
	}

	c.JSON(http.StatusOK, status.Data(account))
}
