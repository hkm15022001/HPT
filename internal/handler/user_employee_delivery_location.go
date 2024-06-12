package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hkm15022001/Supply-Chain-Event-Management/internal/model"
	"gopkg.in/validator.v2"
)

// -------------------- DELIVERY LOCATION HANDLER FUNTION --------------------

// GetDeliveryLocationListHandler in database
func GetDeliveryLocationListHandler(c *gin.Context) {
	store := c.MustGet("store").(*Store)

	deliveryLocations := []model.DeliveryLocation{}
	store.GetReadOnlyConnection().Order("id asc").Find(&deliveryLocations)
	c.JSON(http.StatusOK, gin.H{"delivery_location_list": &deliveryLocations})
	return
}

func getDeliveryLocationOrNotFound(c *gin.Context) (*model.DeliveryLocation, error) {
	store := c.MustGet("store").(*Store)

	deliveryLocation := &model.DeliveryLocation{}
	if err := store.GetReadOnlyConnection().First(deliveryLocation, c.Param("id")).Error; err != nil {
		return deliveryLocation, err
	}
	return deliveryLocation, nil
}

// GetDeliveryLocationHandler in database
func GetDeliveryLocationHandler(c *gin.Context) {

	deliveryLocation, err := getDeliveryLocationOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"delivery_location_info": &deliveryLocation})
	return
}

// CreateDeliveryLocationHandler in database
func CreateDeliveryLocationHandler(c *gin.Context) {
	store := c.MustGet("store").(*Store)

	deliveryLocation := &model.DeliveryLocation{}
	if err := c.ShouldBindJSON(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := store.GetReadOnlyConnection().Where("city = ?", deliveryLocation.City).First(&deliveryLocation).Error
	if err == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := store.GetReadWriteConnection().Create(&deliveryLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "A delivery location has been created!"})
	return
}

// UpdateDeliveryLocationHandler in database
func UpdateDeliveryLocationHandler(c *gin.Context) {
	store := c.MustGet("store").(*Store)

	deliveryLocation, err := getDeliveryLocationOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	deliveryLocation.ID = getIDFromParam(c)
	if err = store.GetReadWriteConnection().Model(&deliveryLocation).Updates(&deliveryLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteDeliveryLocationHandler in database
func DeleteDeliveryLocationHandler(c *gin.Context) {
	store := c.MustGet("store").(*Store)

	if _, err := getDeliveryLocationOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := store.GetReadWriteConnection().Delete(&model.DeliveryLocation{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
