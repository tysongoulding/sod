package core

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceDwarf:
		// Dwarf: Mace Specialization (+1% crit with maces)
		character.AddStat(stats.MeleeCrit, 1*CritRatingPerCritChance)
		character.AddStat(stats.SpellCrit, 1*SpellCritRatingPerCritChance)

		// Big Game Hunter (+5% damage to beasts)
		character.Env.RegisterPostFinalizeEffect(func() {
			for _, t := range character.Env.Encounter.Targets {
				if t.MobType == proto.MobType_MobTypeBeast {
					for _, at := range character.AttackTables[t.UnitIndex] {
						at.DamageDealtMultiplier *= 1.05
						at.CritMultiplier *= 1.05
					}
				}
			}
		})

		// Stoneform (-10% physical damage taken for 8 sec, 3 min CD)
		actionID := ActionID{SpellID: 20594}
		stoneFormAura := character.RegisterAura(Aura{
			Label:    "Stoneform",
			ActionID: actionID,
			Duration: time.Second * 8,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.DamageTakenMultiplier *= 0.90
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.DamageTakenMultiplier /= 0.90
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				stoneFormAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeSurvival,
			ShouldActivate: func(s *Simulation, c *Character) bool {
				return false
			},
		})

	case proto.Race_RaceGnome:
		// Expansive Mind:
		// Priest, Mage, Warlock: Maximum Mana +5%
		// Rogue: Maximum Energy +5%
		// Warrior: Maximum Rage +5%
		if slices.Contains([]proto.Class{proto.Class_ClassPriest, proto.Class_ClassMage, proto.Class_ClassWarlock}, character.Class) {
			character.MultiplyStat(stats.Mana, 1.05)
		} else if character.Class == proto.Class_ClassRogue {
			character.Env.RegisterPostFinalizeEffect(func() {
				character.Unit.energyBar.maxEnergy *= 1.05
			})
		}

		// Eureka! (Next 3 damaging abilities deal +10% damage and cost -10% resource, 2 min CD)
		eurekaActionID := ActionID{SpellID: 462101}
		eurekaAura := character.RegisterAura(Aura{
			Label:     "Eureka!",
			ActionID:  eurekaActionID,
			Duration:  time.Second * 30,
			MaxStacks: 3,
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.SetStacks(sim, 3)
			},
			OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
				if aura.GetStacks() > 0 && result.Landed() && result.Damage > 0 {
					aura.RemoveStack(sim)
				}
			},
		})
		eurekaAura.AttachSpellMod(SpellModConfig{
			Kind:       SpellMod_DamageDone_Pct,
			FloatValue: 1.10,
		})
		eurekaAura.AttachSpellMod(SpellModConfig{
			Kind:     SpellMod_PowerCost_Pct,
			IntValue: -10,
		})

		eurekaSpell := character.RegisterSpell(SpellConfig{
			ActionID: eurekaActionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				eurekaAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: eurekaSpell,
			Type:  CooldownTypeDPS,
		})

	case proto.Race_RaceHuman:
		character.MultiplyStat(stats.Spirit, 1.05)

		// Sword Specialization (+2% crit with swords or 2H swords)
		character.AddStat(stats.MeleeCrit, 2*CritRatingPerCritChance)
		character.AddStat(stats.SpellCrit, 2*SpellCritRatingPerCritChance)

	case proto.Race_RaceNightElf:
		character.AddStat(stats.Dodge, 1)

		// Elune's Light (+10% crit chance for 15 sec, 3 min CD)
		eluneActionID := ActionID{SpellID: 462102}
		critVal := 10.0 * CritRatingPerCritChance
		elunesLightAura := character.RegisterAura(Aura{
			Label:    "Elune's Light",
			ActionID: eluneActionID,
			Duration: time.Second * 15,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.AddStatDynamic(sim, stats.MeleeCrit, critVal)
				character.AddStatDynamic(sim, stats.SpellCrit, critVal)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.AddStatDynamic(sim, stats.MeleeCrit, -critVal)
				character.AddStatDynamic(sim, stats.SpellCrit, -critVal)
			},
		})

		eluneSpell := character.RegisterSpell(SpellConfig{
			ActionID: eluneActionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				elunesLightAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: eluneSpell,
			Type:  CooldownTypeDPS,
		})

	case proto.Race_RaceOrc:
		// Axe Specialization (+1% crit with axes or 2H axes)
		character.AddStat(stats.MeleeCrit, 1*CritRatingPerCritChance)
		character.AddStat(stats.SpellCrit, 1*SpellCritRatingPerCritChance)

		// Blood Fury (+10% AP, +10% RAP, and +10% SP for 15 sec, 2 min CD)
		actionID := ActionID{SpellID: 20572}
		var bfAP, bfRAP, bfSP float64
		bloodFuryAura := character.RegisterAura(Aura{
			Label:    "Blood Fury",
			ActionID: actionID,
			Duration: time.Second * 15,
			OnGain: func(aura *Aura, sim *Simulation) {
				bfAP = character.GetStat(stats.AttackPower) * 0.10
				bfRAP = character.GetStat(stats.RangedAttackPower) * 0.10
				bfSP = character.GetStat(stats.SpellPower) * 0.10
				character.AddStatDynamic(sim, stats.AttackPower, bfAP)
				character.AddStatDynamic(sim, stats.RangedAttackPower, bfRAP)
				character.AddStatDynamic(sim, stats.SpellPower, bfSP)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.AddStatDynamic(sim, stats.AttackPower, -bfAP)
				character.AddStatDynamic(sim, stats.RangedAttackPower, -bfRAP)
				character.AddStatDynamic(sim, stats.SpellPower, -bfSP)
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				bloodFuryAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})

	case proto.Race_RaceTauren:
		// Endurance: Total Health +5%, Hit +1% (physical and spell hit)
		character.MultiplyStat(stats.Health, 1.05)
		character.AddStat(stats.MeleeHit, 1*MeleeHitRatingPerHitChance)
		character.AddStat(stats.SpellHit, 1*SpellHitRatingPerHitChance)

	case proto.Race_RaceTroll:
		// Beast Slaying (+5% damage to beasts)
		character.Env.RegisterPostFinalizeEffect(func() {
			for _, t := range character.Env.Encounter.Targets {
				if t.MobType == proto.MobType_MobTypeBeast {
					for _, at := range character.AttackTables[t.UnitIndex] {
						at.DamageDealtMultiplier *= 1.05
						at.CritMultiplier *= 1.05
					}
				}
			}
		})

		// Berserking (Increases attack and casting speed by 10% for 10 sec, 3 min CD)
		berserkingActionID := ActionID{SpellID: 26297}
		berserkingAura := character.RegisterAura(Aura{
			Label:    "Berserking",
			ActionID: berserkingActionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1.10)
				character.MultiplyAttackSpeed(sim, 1.10)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1 / 1.10)
				character.MultiplyAttackSpeed(sim, 1/1.10)
			},
		})

		berserkingSpell := character.RegisterSpell(SpellConfig{
			ActionID: berserkingActionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				berserkingAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: berserkingSpell,
			Type:  CooldownTypeDPS,
		})

	case proto.Race_RaceUndead:
		// Touch of the Grave:
		// Your spells and attacks have a 5% chance to drain Health from the target, up to 5% of your maximum Health.
		procChance := 0.05
		actionID := ActionID{SpellID: 462103}
		drainSpell := character.RegisterSpell(SpellConfig{
			ActionID:         actionID,
			SpellSchool:      SpellSchoolShadow,
			DefenseType:      DefenseTypeMagic,
			ProcMask:         ProcMaskSpellDamage,
			Flags:            SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				drainAmount := character.MaxHealth() * 0.05
				res := spell.CalcAndDealDamage(sim, target, drainAmount, spell.OutcomeMagicHitAndCrit)
				if res.Landed() {
					character.GainHealth(sim, drainAmount, spell.HealthMetrics(&character.Unit))
				}
			},
		})
		MakeProcTriggerAura(&character.Unit, ProcTrigger{
			Name:       "Touch of the Grave Trigger",
			Callback:   CallbackOnSpellHitDealt,
			Outcome:    OutcomeLanded,
			ProcMask:   ProcMaskDirect | ProcMaskSpellDamage,
			ProcChance: procChance,
			ICD:        0,
			Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
				drainSpell.Cast(sim, result.Target)
			},
		})

	case proto.Race_RaceSkyborneHighOrder, proto.Race_RaceSkyborneWindshaper:
		// Wind Blessed (+1% spellcasting, melee, and ranged Haste)
		character.MultiplyCastSpeed(1.01)
		character.PseudoStats.MeleeSpeedMultiplier *= 1.01
		character.PseudoStats.RangedSpeedMultiplier *= 1.01

		// Elemental Insight (+5% damage to elementals)
		character.Env.RegisterPostFinalizeEffect(func() {
			for _, t := range character.Env.Encounter.Targets {
				if t.MobType == proto.MobType_MobTypeElemental {
					for _, at := range character.AttackTables[t.UnitIndex] {
						at.DamageDealtMultiplier *= 1.05
						at.CritMultiplier *= 1.05
					}
				}
			}
		})
	}
}

func (character *Character) GetFaction() proto.Faction {
	if slices.Contains([]proto.Race{proto.Race_RaceHuman, proto.Race_RaceDwarf, proto.Race_RaceGnome, proto.Race_RaceNightElf, proto.Race_RaceSkyborneHighOrder}, character.Race) {
		return proto.Faction_Alliance
	} else if slices.Contains([]proto.Race{proto.Race_RaceOrc, proto.Race_RaceTroll, proto.Race_RaceTauren, proto.Race_RaceUndead, proto.Race_RaceSkyborneWindshaper}, character.Race) {
		return proto.Faction_Horde
	} else {
		return proto.Faction_Unknown
	}
}
