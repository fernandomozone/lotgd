package bestiary_test

import (
	"math/rand"
	"testing"

	"lotgd/internal/bestiary"
	"lotgd/internal/engine"
	"lotgd/internal/i18n"
)

func TestMonsterGenerator_TierBounds(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	gen := bestiary.NewMonsterGenerator(rng)

	// Gera 20 monstros do Tier 1 e valida bounds
	for i := 0; i < 20; i++ {
		m := gen.GenerateByTier(1)
		if m.Tier != 1 {
			t.Fatalf("esperado monstro de Tier 1, obtido %d", m.Tier)
		}
		if m.Health <= 0 || m.Attack <= 0 || m.Defense < 0 {
			t.Fatalf("atributos inválidos para o monstro: %+v", m)
		}
	}
}

func TestMonsterGenerator_PlayerScaling(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	gen := bestiary.NewMonsterGenerator(rng)

	// Jogador de nível alto deve receber monstros de Tier 4
	m := gen.GenerateForPlayer(10)
	if m.Tier != 4 {
		t.Fatalf("jogador de nível 10 deve enfrentar monstros de Tier 4, obtido Tier %d", m.Tier)
	}
}

func TestMonsterGenerator_LevelOneNoTierTwo(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	gen := bestiary.NewMonsterGenerator(rng)

	// Jogador de nível 1 nunca deve encontrar monstros de Tier 2
	for i := 0; i < 1000; i++ {
		m := gen.GenerateForPlayer(1)
		if m.Tier != 1 {
			t.Fatalf("jogador de nível 1 não deve encontrar monstro de Tier %d", m.Tier)
		}
	}
}

func TestLevelOneWinRateVsFerozMonsters(t *testing.T) {
	rng := rand.New(rand.NewSource(2026))
	ce := engine.NewCombatEngine(rng)

	wins := 0
	simulations := 1000

	for i := 0; i < simulations; i++ {
		// Jogador nível 1 padrão (ATK 6 = 5+1 stick, DEF 2 = 2+0 clothes, HP 20) com 1 poção inicial
		p := &engine.Player{
			Username:     "HeroLevel1",
			Level:        1,
			Health:       20,
			MaxHealth:    20,
			BaseAttack:   5,
			BaseDefense:  2,
			Weapon:       engine.WeaponsCatalog[0], // Stick (+1)
			Armor:        engine.ArmorsCatalog[0],  // Clothes (+0)
			PotionsCount: 1,
		}

		// Seleciona um monstro Tier 1 aleatório e aplica o afixo "Feroz"
		monsterList := bestiary.TierMonstersLists[1]
		chosenID := monsterList[rng.Intn(len(monsterList))]
		tpl := bestiary.CanonicalTemplates[chosenID]

		ferozAffix := bestiary.AffixModifier{
			NamePTBR: "Feroz", HPMult: 1.1, ATKMult: 1.15, DEFMult: 1.0, XPMult: 1.3, GoldMult: 1.2,
		}

		hp := int(float64(tpl.BaseHP) * ferozAffix.HPMult)
		atk := int(float64(tpl.BaseATK) * ferozAffix.ATKMult)
		def := int(float64(tpl.BaseDEF) * ferozAffix.DEFMult)

		m := &engine.Monster{
			ID:        chosenID,
			Name:      "Feroz " + i18n.GetMonsterName(chosenID),
			Tier:      1,
			Health:    hp,
			MaxHealth: hp,
			Attack:    atk,
			Defense:   def,
			Prefix:    "Feroz",
		}

		// Simula o combate até a vitória ou derrota
		for p.IsAlive() && m.IsAlive() {
			// Se o jogador estiver com vida crítica (<= 8 HP) e possuir poção, usa a poção
			if p.Health <= 8 && p.PotionsCount > 0 {
				_, err := ce.UsePotion(p)
				if err == nil {
					continue
				}
			}

			res := ce.Attack(p, m)
			if res.MonsterDefeated {
				wins++
				break
			}
			if res.PlayerDefeated {
				break
			}
		}
	}

	winRate := float64(wins) / float64(simulations)
	t.Logf("Taxa de vitória de jogador Nível 1 vs monstros Feroz do Tier 1: %.2f%% (%d/%d)", winRate*100, wins, simulations)

	if winRate < 0.70 {
		t.Fatalf("esperado taxa de vitória >= 70%%, obtido: %.2f%%", winRate*100)
	}
}

func TestGenerateDragonOfDay_Deterministic(t *testing.T) {
	day1 := "2026-08-25"
	day2 := "2026-08-26"

	d1a := bestiary.GenerateDragonOfDay(day1)
	d1b := bestiary.GenerateDragonOfDay(day1)
	d2 := bestiary.GenerateDragonOfDay(day2)

	// Mesma data deve gerar exatamente o mesmo dragão
	if d1a.Name != d1b.Name || d1a.Health != d1b.Health || d1a.Attack != d1b.Attack || d1a.Defense != d1b.Defense {
		t.Fatalf("o Dragão do Dia deve ser 100%% determinístico para a mesma data. d1a=%+v, d1b=%+v", d1a, d1b)
	}

	if d1a.ID != i18n.MonsterDragon || !d1a.IsDragon {
		t.Fatalf("o Dragão deve ter ID canônica e IsDragon=true")
	}

	// Datas diferentes devem gerar variações
	if d1a.Health == d2.Health && d1a.Attack == d2.Attack && d1a.Name == d2.Name {
		t.Logf("Aviso: stats coincidentemente iguais entre dias diferentes (raro mas possível)")
	}
}
