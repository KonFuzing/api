package domain

import (
	"fmt"
	"math/rand"
	"time"
)

type Cowboy struct {
	ID       string
	Name     string
	Health   int
	Damage   int
	Speed    int
	Accuracy float64
}

type BattleResult struct {
	Winner string
	Logs   []string
}

// Business Logic: การต่อสู้ (Pure Logic)
func SimulateFight(c1, c2 *Cowboy) BattleResult {
	rand.Seed(time.Now().UnixNano())
	var logs []string
	logs = append(logs, fmt.Sprintf("🔥 Match Start: %s (HP:%d) VS %s (HP:%d)", c1.Name, c1.Health, c2.Name, c2.Health))

	hp1, hp2 := c1.Health, c2.Health
	var attacker, defender *Cowboy
	var attackerHP, defenderHP *int

	if c1.Speed >= c2.Speed {
		attacker, defender = c1, c2
		attackerHP, defenderHP = &hp1, &hp2
		logs = append(logs, fmt.Sprintf("⚡ %s is faster!", c1.Name))
	} else {
		attacker, defender = c2, c1
		attackerHP, defenderHP = &hp2, &hp1
		logs = append(logs, fmt.Sprintf("⚡ %s is faster!", c2.Name))
	}

	turn := 1
	for *attackerHP > 0 && *defenderHP > 0 {
		logs = append(logs, fmt.Sprintf("--- Turn %d ---", turn))
		if rand.Float64() <= attacker.Accuracy {
			variance := float64(attacker.Damage) * 0.2
			dmg := attacker.Damage + rand.Intn(int(variance)*2+1) - int(variance)
			*defenderHP -= dmg
			if *defenderHP < 0 { *defenderHP = 0 }
			logs = append(logs, fmt.Sprintf("💥 %s hits %s for %d (HP left: %d)", attacker.Name, defender.Name, dmg, *defenderHP))
		} else {
			logs = append(logs, fmt.Sprintf("💨 %s missed!", attacker.Name))
		}

		if *defenderHP <= 0 { break }
		attacker, defender = defender, attacker
		attackerHP, defenderHP = defenderHP, attackerHP
		turn++
	}

	return BattleResult{Winner: attacker.Name, Logs: logs}
}