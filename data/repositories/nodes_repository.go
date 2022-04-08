package repositories

import (
	"loadbalancer/data"
	"loadbalancer/data/models"
	"time"
)

func GetAllNodes() ([]models.Node, error) {
	db := data.DbConnection

	var nodes []models.Node
	result := db.Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func GetAllNodeAddresses() ([]string, error) {
	db := data.DbConnection

	var nodes []models.Node
	result := db.Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	var addresses []string

	for _, n := range nodes {
		addresses = append(addresses, n.Address)
	}

	return addresses, nil
}

func CreateNode(address string) error {
	db := data.DbConnection

	newNode := models.Node{Address: address, CreatedAt: time.Now()}
	result := db.Create(&newNode)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteNodeByAddress(address string) error {
	db := data.DbConnection

	result := db.Delete(&models.Node{}, "address = ?", address)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
