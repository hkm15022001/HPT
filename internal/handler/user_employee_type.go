package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hkm15022001/Supply-Chain-Event-Management/internal/model"
	"gopkg.in/validator.v2"
)

// -------------------- EMPLOYEE TYPE HANDLER FUNTION --------------------

// GetEmployeeTypeListHandler in database
func GetEmployeeTypeListHandler(c *gin.Context) {
	store := c.MustGet("store").(*Store)

	employeeTypes := []model.EmployeeType{}
	store.GetReadOnlyConnection().Order("id asc").Find(&employeeTypes)
	c.JSON(http.StatusOK, gin.H{"employee_type_list": &employeeTypes})
	return
}

func getEmployeeTypeOrNotFound(c *gin.Context) (*model.EmployeeType, error) {
	store := c.MustGet("store").(*Store)

	employeeType := &model.EmployeeType{}
	if err := store.GetReadOnlyConnection().First(employeeType, c.Param("id")).Error; err != nil {
		return employeeType, err
	}
	return employeeType, nil
}

// GetEmployeeTypeHandler in database
func GetEmployeeTypeHandler(c *gin.Context) {

	employeeType, err := getEmployeeTypeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"employee_type_info": &employeeType})
	return
}

// CreateEmployeeTypeHandler in database
func CreateEmployeeTypeHandler(c *gin.Context) {
	store := c.MustGet("store").(*Store)

	employeeType := &model.EmployeeType{}
	if err := c.ShouldBindJSON(&employeeType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&employeeType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := store.GetReadWriteConnection().Create(employeeType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "An employee type has been created!"})
	return
}

// UpdateEmployeeTypeHandler in database
func UpdateEmployeeTypeHandler(c *gin.Context) {
	store := c.MustGet("store").(*Store)

	employeeType, err := getEmployeeTypeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&employeeType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&employeeType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	employeeType.ID = getIDFromParam(c)
	if err = store.GetReadWriteConnection().Model(&employeeType).Updates(&employeeType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteEmployeeTypeHandler in database
func DeleteEmployeeTypeHandler(c *gin.Context) {
	store := c.MustGet("store").(*Store)

	if _, err := getEmployeeTypeOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := store.GetReadWriteConnection().Delete(&model.EmployeeType{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
