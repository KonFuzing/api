package domain

import (
	"fmt"
	"math/rand"
	"time"
	"api/services/arena/internal/core/domain/entity"
)

// Value Object: เก็บผลลัพธ์ (ไม่มี logic)
type BattleResult struct {
	Winner string
	Logs   []string
}

// Domain Service: ควบคุมกฏการต่อสู้ (Battle Logic)
// รับ Entity เข้ามา และสั่งงานผ่าน Method ของ Entity
func SimulateFight(c1, c2 *entity.Cowboy) BattleResult {
	// Seed random
	rand.Seed(time.Now().UnixNano())

	var logs []string
	logs = append(logs, fmt.Sprintf("🔥 Match Start: %s (HP:%d) VS %s (HP:%d)", c1.Name, c1.Health, c2.Name, c2.Health))

	// สร้างตัวแปร pointer ชั่วคราวเพื่อสลับเทิร์น (Attacker / Defender)
	// เราใช้ตัวจริงเลยเพราะ Cowboy เป็น Pointer อยู่แล้ว และเรามี Method TakeDamage คุม State
	var attacker, defender *entity.Cowboy

	if c1.Speed >= c2.Speed {
		attacker, defender = c1, c2
		logs = append(logs, fmt.Sprintf("⚡ %s is faster!", c1.Name))
	} else {
		attacker, defender = c2, c1
		logs = append(logs, fmt.Sprintf("⚡ %s is faster!", c2.Name))
	}

	turn := 1

	// วนลูปจนกว่าจะมีฝ่ายใดฝ่ายหนึ่งตาย (ใช้ Method IsDead เช็ค)
	for !c1.IsDead() && !c2.IsDead() {
		logs = append(logs, fmt.Sprintf("--- Turn %d ---", turn))

		// คำนวณโอกาสแม่นยำ
		if rand.Float64() <= attacker.Accuracy {
			// คำนวณ Damage (Variance +/- 20%)
			variance := float64(attacker.Damage) * 0.2
			dmg := attacker.Damage + rand.Intn(int(variance)*2+1) - int(variance)

			// 💥 เรียกใช้ Logic ภายใน Entity ให้รับดาเมจ
			defender.TakeDamage(dmg)

			logs = append(logs, fmt.Sprintf("💥 %s hits %s for %d (HP left: %d)",
				attacker.Name, defender.Name, dmg, defender.Health))
		} else {
			logs = append(logs, fmt.Sprintf("💨 %s missed!", attacker.Name))
		}

		// เช็คจบเกมทันทีหลังโดนยิง
		if defender.IsDead() {
			break
		}

		// สลับฝั่ง
		attacker, defender = defender, attacker
		turn++
	}

	return BattleResult{Winner: attacker.Name, Logs: logs}
}
