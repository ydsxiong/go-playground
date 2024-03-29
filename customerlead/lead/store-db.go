package lead

import "github.com/jinzhu/gorm"

type databaseStore struct {
	db *gorm.DB
}

func NewDatabaseStore(gormdb *gorm.DB) LeadStore {
	return &databaseStore{gormdb}
}

/**
id and createdAt will have been generated by the database insert operation.
*/
func (ds *databaseStore) Save(lead Lead) error {
	return ds.db.Save(lead).Error
}

func (ds *databaseStore) FindAll() (Leads, error) {
	var leads Leads
	if err := ds.db.Find(&leads).Error; err != nil {
		return nil, err
	}
	return leads, nil
}

func (ds *databaseStore) FindByEmail(email string) (*Lead, error) {
	lead := Lead{}
	if err := ds.db.First(&lead, Lead{Email: email}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &lead, nil

}
