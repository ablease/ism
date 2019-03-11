package osbapi

type Broker struct {
	ID        string
	CreatedAt string
	Name      string
	URL       string
	Username  string
	Password  string
}

type Service struct {
	ID          string
	Name        string
	Description string
	BrokerID    string
}

type Plan struct {
	Name      string
	ServiceID string
}

type Platform struct {
	ID   string
	Name string
	URL  string
}
