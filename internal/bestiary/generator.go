package bestiary

import (
	"fmt"
	"math/rand"

	"lotgd/internal/engine"
	"lotgd/internal/i18n"
)

// MonsterGenerator produz criaturas procedurais com modificadores estocásticos e escalonamento por nível.
//
// Didática Go: Injetamos o gerador `*rand.Rand` via campo de struct para isolar o estado de aleatoriedade,
// permitindo o uso de sementes fixas em testes unitários para reprodutibilidade total.
type MonsterGenerator struct {
	rng *rand.Rand
}

// NewMonsterGenerator instancia um novo gerador de monstros procedurais.
// Se `rng` for nil, inicializa um gerador com semente aleatória baseada no relógio do sistema.
func NewMonsterGenerator(rng *rand.Rand) *MonsterGenerator {
	if rng == nil {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	return &MonsterGenerator{rng: rng}
}

// GenerateForPlayer seleciona e gera um monstro balanceado para o nível atual do herói.
//
// Regras de Distribuição por Tier:
// - Nível 1       -> Tier 1 (estritamente sem encontros em Tier 2 para viabilidade inicial)
// - Nível 2       -> Tier 1 (com 10% de chance de encontro desafiador em Tier 2)
// - Níveis 3 a 4  -> Tier 2 (com 15% de chance de Tier 3)
// - Níveis 5 a 7  -> Tier 3 (com 15% de chance de Tier 4)
// - Níveis 8 a 10 -> Tier 4
func (mg *MonsterGenerator) GenerateForPlayer(playerLevel int) engine.Monster {
	var tier int
	roll := mg.rng.Float64()

	switch {
	case playerLevel <= 2:
		tier = 1
		if playerLevel > 1 && roll < 0.10 {
			tier = 2
		}
	case playerLevel <= 4:
		tier = 2
		if roll < 0.15 {
			tier = 3
		}
	case playerLevel <= 7:
		tier = 3
		if roll < 0.15 {
			tier = 4
		}
	default:
		tier = 4
	}

	return mg.GenerateByTier(tier)
}

// GenerateByTier gera um monstro procedural pertencente estritamente ao Tier requisitado.
//
// Fluxo do Algoritmo:
// 1. Clampa o parâmetro Tier entre 1 e 4.
// 2. Sorteia um ID de monstro da lista de criaturas do Tier correspondente.
// 3. Aplica uma chance estocástica de 50% de sortear e aplicar um afixo procedural (`AvailableAffixes`).
// 4. Modifica os atributos (HP, ATK, DEF, XP, Ouro) através dos multiplicadores do afixo.
// 5. Garante valores mínimos de viabilidade (HP >= 5, ATK >= 1, XP >= 1, Ouro >= 1).
func (mg *MonsterGenerator) GenerateByTier(tier int) engine.Monster {
	if tier < 1 {
		tier = 1
	}
	if tier > 4 {
		tier = 4
	}

	monsterList := TierMonstersLists[tier]
	chosenID := monsterList[mg.rng.Intn(len(monsterList))]
	tpl := CanonicalTemplates[chosenID]

	baseName := i18n.GetMonsterName(chosenID)

	// 50% de chance de receber um afixo especial procedural
	var prefix string
	hp := tpl.BaseHP
	atk := tpl.BaseATK
	def := tpl.BaseDEF
	xp := tpl.BaseXP
	gold := tpl.BaseGold

	if mg.rng.Float64() < 0.50 {
		affix := AvailableAffixes[mg.rng.Intn(len(AvailableAffixes))]
		prefix = affix.NamePTBR
		hp = int(float64(hp) * affix.HPMult)
		atk = int(float64(atk) * affix.ATKMult)
		def = int(float64(def) * affix.DEFMult)
		xp = int(float64(xp) * affix.XPMult)
		gold = int(float64(gold) * affix.GoldMult)
	}

	fullName := baseName
	if prefix != "" {
		fullName = fmt.Sprintf("%s %s", prefix, baseName)
	}

	// Trava os atributos em limites mínimos viáveis para evitar monstros inválidos
	if hp < 5 {
		hp = 5
	}
	if atk < 1 {
		atk = 1
	}
	if xp < 1 {
		xp = 1
	}
	if gold < 1 {
		gold = 1
	}

	return engine.Monster{
		ID:         chosenID,
		Name:       fullName,
		Tier:       tier,
		Health:     hp,
		MaxHealth:  hp,
		Attack:     atk,
		Defense:    def,
		XPReward:   xp,
		GoldReward: gold,
		Prefix:     prefix,
		IsDragon:   false,
	}
}
