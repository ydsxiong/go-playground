package lead

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jinzhu/gorm"
)

type Lead struct {
	gorm.Model
	Fname         string `json:"first_name"`
	Lname         string `json:"last_name"`
	Email         string `gorm:"unique" json:"email"`
	Company       string `json:"company"`
	Postcode      string `json:"postcode"`
	TermsAccepted bool   `json:"terms_accepted"`
}

/**
A customerized json unmarshaller for Lead data object to add an extra layer of validating the presence of all mandatory fields,
in a context whenever data is to be unmarshalled into a lead object.
*/
func (l *Lead) UnmarshalJSON(data []byte) (err error) {
	var auxObj leadAux
	err = json.Unmarshal(data, &auxObj)
	if err != nil {
		return
	}

	err = auxObj.validate()
	if err != nil {
		return
	}

	l.Fname = auxObj.Fname
	l.Lname = auxObj.Lname
	l.Email = auxObj.Email
	l.Company = auxObj.Company
	l.Postcode = auxObj.Postcode
	l.TermsAccepted = auxObj.TermsAccepted
	return
}

type leadAux struct {
	Fname         string `json:"first_name"`
	Lname         string `json:"last_name"`
	Email         string `json:"email"`
	Company       string `json:"company"`
	Postcode      string `json:"postcode"`
	TermsAccepted bool   `json:"terms_accepted"`
}

func (l *leadAux) validate() (err error) {
	var missedFields []string

	if l.Fname == "" {
		missedFields = append(missedFields, "first_name")
	}
	if l.Lname == "" {
		missedFields = append(missedFields, "last_name")
	}
	if l.Email == "" {
		missedFields = append(missedFields, "email")
	}
	if !l.TermsAccepted {
		missedFields = append(missedFields, "terms_accepted")
	}
	if len(missedFields) > 0 {
		err = fmt.Errorf("Required fields missing: %v", missedFields)
	}
	return
}

type Leads []Lead

/**
return a copy of the leads data
*/
func (l Leads) FindAll() Leads {
	results := make(Leads, 0)
	for _, v := range l {
		results = append(results, v)
	}
	return results
}

/**
return a copy of lead data
*/
func (l Leads) FindByEmail(email string) *Lead {
	for _, v := range l {
		if v.Email == email {
			return &v
		}
	}
	return nil
}

func (l Leads) Sort() {
	sort.Slice(l, func(i, j int) bool {
		return l[i].ID < l[j].ID
	})
}
