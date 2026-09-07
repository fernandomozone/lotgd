package bestiary

import (
	"lotgd/internal/i18n"
)

// MonsterTemplate define a estrutura de dados base e estatísticas canônicas de uma espécie de criatura.
//
// Didática Go: A struct agrupa os atributos básicos sem afixos procedurais. As tags `json:"..."`
// facilitam a serialização para exportação de dados ou depuração.
type MonsterTemplate struct {
	ID       i18n.MonsterID `json:"id"`
	Tier     int            `json:"tier"`
	BaseHP   int            `json:"base_hp"`
	BaseATK  int            `json:"base_atk"`
	BaseDEF  int            `json:"base_def"`
	BaseXP   int            `json:"base_xp"`
	BaseGold int            `json:"base_gold"`
}

// AffixModifier define alterações estocásticas em atributos e recompensas aplicadas dinamicamente aos monstros.
//
// Didática Go: Os multiplicadores de ponto flutuante (`float64`) aplicam variações percentuais
// sobre os atributos base da criatura no momento da geração estocástica.
type AffixModifier struct {
	NamePTBR string  `json:"name_pt_br"`
	HPMult   float64 `json:"hp_mult"`
	ATKMult  float64 `json:"atk_mult"`
	DEFMult  float64 `json:"def_mult"`
	XPMult   float64 `json:"xp_mult"`
	GoldMult float64 `json:"gold_mult"`
}

// AvailableAffixes contém a lista de afixos procedurais que podem ser sorteados aleatoriamente ao gerar um monstro.
var AvailableAffixes = []AffixModifier{
	{NamePTBR: "Feroz", HPMult: 1.1, ATKMult: 1.15, DEFMult: 1.0, XPMult: 1.3, GoldMult: 1.2},
	{NamePTBR: "Covarde", HPMult: 0.8, ATKMult: 0.8, DEFMult: 0.9, XPMult: 0.8, GoldMult: 1.1},
	{NamePTBR: "Enfurecido", HPMult: 1.2, ATKMult: 1.25, DEFMult: 0.8, XPMult: 1.4, GoldMult: 1.2},
	{NamePTBR: "Sortudo", HPMult: 1.0, ATKMult: 1.0, DEFMult: 1.0, XPMult: 1.1, GoldMult: 2.5},
	{NamePTBR: "Faminto", HPMult: 1.1, ATKMult: 1.2, DEFMult: 1.0, XPMult: 1.2, GoldMult: 1.0},
	{NamePTBR: "Preguiçoso", HPMult: 1.2, ATKMult: 0.7, DEFMult: 1.2, XPMult: 0.9, GoldMult: 0.9},
	{NamePTBR: "Astuto", HPMult: 1.0, ATKMult: 1.2, DEFMult: 1.3, XPMult: 1.3, GoldMult: 1.4},
	{NamePTBR: "Gigantesco", HPMult: 1.6, ATKMult: 1.3, DEFMult: 1.2, XPMult: 1.6, GoldMult: 1.5},
}

// CanonicalTemplates mapeia cada `MonsterID` para o seu modelo canônico contendo atributos base.
//
// Didática Go: Mapeamento direto via `map[i18n.MonsterID]MonsterTemplate` provê acesso imediato em tempo O(1)
// para obtenção das estatísticas padrão dos 20 monstros organizados do Tier 1 ao Tier 4.
var CanonicalTemplates = map[i18n.MonsterID]MonsterTemplate{
	// Tier 1 (Iniciante - Níveis 1 a 2)
	i18n.MonsterSewerRat:     {ID: i18n.MonsterSewerRat, Tier: 1, BaseHP: 12, BaseATK: 3, BaseDEF: 1, BaseXP: 15, BaseGold: 8},
	i18n.MonsterMossSpider:   {ID: i18n.MonsterMossSpider, Tier: 1, BaseHP: 15, BaseATK: 4, BaseDEF: 1, BaseXP: 20, BaseGold: 12},
	i18n.MonsterClumsyKobold: {ID: i18n.MonsterClumsyKobold, Tier: 1, BaseHP: 18, BaseATK: 5, BaseDEF: 2, BaseXP: 25, BaseGold: 16},
	i18n.MonsterGreenSlime:   {ID: i18n.MonsterGreenSlime, Tier: 1, BaseHP: 22, BaseATK: 4, BaseDEF: 3, BaseXP: 28, BaseGold: 14},
	i18n.MonsterSinisterCrow: {ID: i18n.MonsterSinisterCrow, Tier: 1, BaseHP: 14, BaseATK: 6, BaseDEF: 1, BaseXP: 22, BaseGold: 10},

	// Tier 2 (Aventureiro - Níveis 3 a 4)
	i18n.MonsterGoblinScout:   {ID: i18n.MonsterGoblinScout, Tier: 2, BaseHP: 35, BaseATK: 8, BaseDEF: 4, BaseXP: 55, BaseGold: 35},
	i18n.MonsterOrcRecruit:    {ID: i18n.MonsterOrcRecruit, Tier: 2, BaseHP: 45, BaseATK: 10, BaseDEF: 5, BaseXP: 75, BaseGold: 45},
	i18n.MonsterShadowWolf:    {ID: i18n.MonsterShadowWolf, Tier: 2, BaseHP: 38, BaseATK: 11, BaseDEF: 3, BaseXP: 70, BaseGold: 40},
	i18n.MonsterRoadBandit:    {ID: i18n.MonsterRoadBandit, Tier: 2, BaseHP: 42, BaseATK: 9, BaseDEF: 4, BaseXP: 65, BaseGold: 60},
	i18n.MonsterRustySkeleton: {ID: i18n.MonsterRustySkeleton, Tier: 2, BaseHP: 40, BaseATK: 10, BaseDEF: 6, BaseXP: 80, BaseGold: 30},

	// Tier 3 (Veterano - Níveis 5 a 7)
	i18n.MonsterBridgeTroll:       {ID: i18n.MonsterBridgeTroll, Tier: 3, BaseHP: 85, BaseATK: 16, BaseDEF: 8, BaseXP: 180, BaseGold: 110},
	i18n.MonsterRudeOgre:          {ID: i18n.MonsterRudeOgre, Tier: 3, BaseHP: 100, BaseATK: 18, BaseDEF: 7, BaseXP: 210, BaseGold: 130},
	i18n.MonsterNoviceNecromancer: {ID: i18n.MonsterNoviceNecromancer, Tier: 3, BaseHP: 75, BaseATK: 21, BaseDEF: 6, BaseXP: 230, BaseGold: 160},
	i18n.MonsterSingingHarpy:      {ID: i18n.MonsterSingingHarpy, Tier: 3, BaseHP: 70, BaseATK: 20, BaseDEF: 8, BaseXP: 190, BaseGold: 140},
	i18n.MonsterCrackedStoneGolem: {ID: i18n.MonsterCrackedStoneGolem, Tier: 3, BaseHP: 120, BaseATK: 15, BaseDEF: 14, BaseXP: 250, BaseGold: 120},

	// Tier 4 (Lendário - Níveis 8 a 10)
	i18n.MonsterFallenKnight:     {ID: i18n.MonsterFallenKnight, Tier: 4, BaseHP: 160, BaseATK: 28, BaseDEF: 16, BaseXP: 500, BaseGold: 300},
	i18n.MonsterMountainWyvern:   {ID: i18n.MonsterMountainWyvern, Tier: 4, BaseHP: 190, BaseATK: 32, BaseDEF: 14, BaseXP: 620, BaseGold: 380},
	i18n.MonsterLesserLich:       {ID: i18n.MonsterLesserLich, Tier: 4, BaseHP: 150, BaseATK: 36, BaseDEF: 12, BaseXP: 700, BaseGold: 450},
	i18n.MonsterColdGazeBasilisk: {ID: i18n.MonsterColdGazeBasilisk, Tier: 4, BaseHP: 180, BaseATK: 30, BaseDEF: 18, BaseXP: 580, BaseGold: 350},
	i18n.MonsterSwampSpecter:     {ID: i18n.MonsterSwampSpecter, Tier: 4, BaseHP: 140, BaseATK: 34, BaseDEF: 15, BaseXP: 550, BaseGold: 320},
}

// TierMonstersLists agrupa os IDs de monstros em slices indexados pelo seu número de Tier (1 a 4).
//
// Didática Go: Permite que o algoritmo de geração escolha aleatoriamente qualquer monstro pertencente
// ao Tier selecionado via indexação direta de fatias (slices).
var TierMonstersLists = map[int][]i18n.MonsterID{
	1: {i18n.MonsterSewerRat, i18n.MonsterMossSpider, i18n.MonsterClumsyKobold, i18n.MonsterGreenSlime, i18n.MonsterSinisterCrow},
	2: {i18n.MonsterGoblinScout, i18n.MonsterOrcRecruit, i18n.MonsterShadowWolf, i18n.MonsterRoadBandit, i18n.MonsterRustySkeleton},
	3: {i18n.MonsterBridgeTroll, i18n.MonsterRudeOgre, i18n.MonsterNoviceNecromancer, i18n.MonsterSingingHarpy, i18n.MonsterCrackedStoneGolem},
	4: {i18n.MonsterFallenKnight, i18n.MonsterMountainWyvern, i18n.MonsterLesserLich, i18n.MonsterColdGazeBasilisk, i18n.MonsterSwampSpecter},
}
