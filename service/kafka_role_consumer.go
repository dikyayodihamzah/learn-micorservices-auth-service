package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"gitlab.com/learn-micorservices/auth-service/model/domain"
	"gitlab.com/learn-micorservices/auth-service/repository"
)

type KafkaRoleConsumerService interface {
	Insert(message []byte) error
	Update(message []byte) error
	Delete(message []byte) error
}

type kafkaRoleConsumerService struct {
	roleRepository repository.RoleRepository
}

func NewKafkaRoleConsumerService(roleRepository repository.RoleRepository) KafkaRoleConsumerService {
	return &kafkaRoleConsumerService{
		roleRepository: roleRepository,
	}
}

func (service *kafkaRoleConsumerService) Insert(message []byte) error {
	var roleDTO map[string]interface{}
	
	if err := json.Unmarshal(message, &roleDTO); err != nil {
		log.Println("error kafka insert role consumer:", err.Error())
	}
	
	id := roleDTO["id"]
	name := roleDTO["name"]
	createdat := roleDTO["created_at"]
	updatedat := roleDTO["updated_at"]
	created_at, _ := time.Parse(time.RFC3339, createdat.(string))
	update_at, _ := time.Parse(time.RFC3339, updatedat.(string))
	
	role := domain.Role{
		ID : id.(string),
		Name : name.(string),
		CreatedAt : created_at,
		UpdatedAt : update_at,
	}

	if err := service.roleRepository.Create(context.Background(), role); err != nil {
		log.Println("error kafka insert role consumer:", err.Error())
		return err
	}

	return nil
}

func (service *kafkaRoleConsumerService) Update(message []byte) error {
	var roleDTO map[string]interface{}

	if err := json.Unmarshal(message, &roleDTO); err != nil {
		log.Println("error kafka update role consumer:", err.Error())
	}
	
	id := roleDTO["id"]
	name := roleDTO["name"]
	createdat := roleDTO["created_at"]
	updatedat := roleDTO["updated_at"]
	created_at, _ := time.Parse(time.RFC3339, createdat.(string))
	update_at, _ := time.Parse(time.RFC3339, updatedat.(string))
	
	role := domain.Role{
		ID : id.(string),
		Name : name.(string),
		CreatedAt : created_at,
		UpdatedAt : update_at,
	}

	if err := service.roleRepository.Update(context.Background(), role.ID, role); err != nil {
		log.Println("error kafka update role consumer:", err.Error())
		return err
	}

	return nil
}

func (service *kafkaRoleConsumerService) Delete(message []byte) error {
	var roleDTO map[string]interface{}

	if err := json.Unmarshal(message, &roleDTO); err != nil {
		log.Println("error kafka update user consumer:", err.Error())
	}

	user := domain.Role{
		ID: roleDTO["id"].(string),
	}
	
	if err := service.roleRepository.Delete(context.Background(), user.ID); err != nil {
		log.Println("error kafka update user consumer:", err.Error())
		return err
	}

	return nil
}