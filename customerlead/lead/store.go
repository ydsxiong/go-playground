package lead

type LeadStore interface {
	Save(lead Lead) error
	FindAll() (Leads, error)
	FindByEmail(email string) (*Lead, error)
}
