package main

type Customer struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func (c *Customer) Register()                 {}
func (c *Customer) Login()                    {}
func (c *Customer) UpdateProfile()            {}
func (c *Customer) ViewOrderHistory() []Order { return nil }
