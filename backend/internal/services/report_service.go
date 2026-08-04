package services

type LDAPService struct{}

func NewLDAPService() *LDAPService {
	return &LDAPService{}
}

func (s *LDAPService) SyncUsers() error {

	return nil
}