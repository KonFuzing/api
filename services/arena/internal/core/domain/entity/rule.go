package entity

func (c *Cowboy) IsDead() bool {
	return c.Health <= 0
}

func (c *Cowboy) TakeDamage(dmg int) {
	c.Health -= dmg
	if c.Health < 0 {
		c.Health = 0
	}
}
