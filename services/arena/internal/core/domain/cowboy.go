package domain

// Cowboy Entity: เก็บข้อมูลและ Logic ภายในตัว
type Cowboy struct {
	ID       string
	Name     string
	Health   int
	Damage   int
	Speed    int
	Accuracy float64
}

// IsDead: เช็คว่าตายหรือยัง (Logic ส่วนตัว)
func (c *Cowboy) IsDead() bool {
	return c.Health <= 0
}

// TakeDamage: คำนวณการรับดาเมจ (Logic ส่วนตัว)
func (c *Cowboy) TakeDamage(dmg int) {
	c.Health -= dmg
	if c.Health < 0 {
		c.Health = 0
	}
}
