package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	_ "github.com/wowsims/sod/sim/common"
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	feraldruid "github.com/wowsims/sod/sim/druid/feral"
	dpshunter "github.com/wowsims/sod/sim/hunter/dps_hunter"
	dpsmage "github.com/wowsims/sod/sim/mage/dps_mage"
	retributionpaladin "github.com/wowsims/sod/sim/paladin/retribution"
	shadowpriest "github.com/wowsims/sod/sim/priest/shadow"
	dpsrogue "github.com/wowsims/sod/sim/rogue/dps_rogue"
	enhancementshaman "github.com/wowsims/sod/sim/shaman/enhancement"
	dpswarlock "github.com/wowsims/sod/sim/warlock/dps"
	dpswarrior "github.com/wowsims/sod/sim/warrior/dps_warrior"
	"google.golang.org/protobuf/encoding/protojson"
)

type RaceResult struct {
	Race string  `json:"race"`
	DPS  float64 `json:"dps"`
}

type ClassResults struct {
	Class    string       `json:"class"`
	Duration float64      `json:"duration"`
	Results  []RaceResult `json:"results"`
}

type ClassConfig struct {
	ClassName   string
	Class       proto.Class
	GearSet     string
	Apl         string
	Talents     string
	SpecOptions any
	Races       []proto.Race
}

func main() {
	dpswarrior.RegisterDpsWarrior()
	dpsrogue.RegisterDpsRogue()
	dpsmage.RegisterDPSMage()
	dpshunter.RegisterDPSHunter()
	dpswarlock.RegisterDpsWarlock()
	shadowpriest.RegisterShadowPriest()
	retributionpaladin.RegisterRetributionPaladin()
	enhancementshaman.RegisterEnhancementShaman()
	feraldruid.RegisterFeralDruid()

	configs := []ClassConfig{
		{
			ClassName: "Warrior",
			Class:     proto.Class_ClassWarrior,
			GearSet:   "ui/warrior/gear_sets/phase_4_dw.gear.json",
			Apl:       "ui/warrior/apls/phase_4_fury.apl.json",
			Talents:   "20305020302-05050005525010051",
			SpecOptions: &proto.Player_Warrior{
				Warrior: &proto.Warrior{
					Options: &proto.Warrior_Options{},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceHuman,
				proto.Race_RaceDwarf,
				proto.Race_RaceNightElf,
				proto.Race_RaceGnome,
				proto.Race_RaceSkyborneHighOrder,
				proto.Race_RaceOrc,
				proto.Race_RaceUndead,
				proto.Race_RaceTauren,
				proto.Race_RaceTroll,
				proto.Race_RaceSkyborneWindshaper,
			},
		},
		{
			ClassName: "Rogue",
			Class:     proto.Class_ClassRogue,
			GearSet:   "ui/rogue/gear_sets/p4_saber.gear.json",
			Apl:       "ui/rogue/apls/Saber_DPS_60.apl.json",
			Talents:   "00532310155104-02330520000501",
			SpecOptions: &proto.Player_Rogue{
				Rogue: &proto.Rogue{
					Options: &proto.RogueOptions{
						HonorAmongThievesCritRate: 100,
					},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceHuman,
				proto.Race_RaceDwarf,
				proto.Race_RaceNightElf,
				proto.Race_RaceGnome,
				proto.Race_RaceSkyborneHighOrder,
				proto.Race_RaceOrc,
				proto.Race_RaceUndead,
				proto.Race_RaceTroll,
				proto.Race_RaceSkyborneWindshaper,
			},
		},
		{
			ClassName: "Mage",
			Class:     proto.Class_ClassMage,
			GearSet:   "ui/mage/gear_sets/p4_fire.gear.json",
			Apl:       "ui/mage/apls/p4_fire.apl.json",
			Talents:   "21-5052300123033151-203500031",
			SpecOptions: &proto.Player_Mage{
				Mage: &proto.Mage{
					Options: &proto.Mage_Options{
						Armor: proto.Mage_Options_MageArmor,
					},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceHuman,
				proto.Race_RaceGnome,
				proto.Race_RaceSkyborneHighOrder,
				proto.Race_RaceOrc,
				proto.Race_RaceUndead,
				proto.Race_RaceTroll,
			},
		},
		{
			ClassName: "Hunter",
			Class:     proto.Class_ClassHunter,
			GearSet:   "ui/hunter/gear_sets/p4_ranged.gear.json",
			Apl:       "ui/hunter/apls/p4_ranged.apl.json",
			Talents:   "-055500005-3305202202303051",
			SpecOptions: &proto.Player_Hunter{
				Hunter: &proto.Hunter{
					Options: &proto.Hunter_Options{
						PetType: proto.Hunter_Options_Cat,
					},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceHuman,
				proto.Race_RaceDwarf,
				proto.Race_RaceNightElf,
				proto.Race_RaceSkyborneHighOrder,
				proto.Race_RaceOrc,
				proto.Race_RaceTauren,
				proto.Race_RaceTroll,
				proto.Race_RaceSkyborneWindshaper,
			},
		},
		{
			ClassName: "Warlock",
			Class:     proto.Class_ClassWarlock,
			GearSet:   "ui/warlock/gear_sets/p4/destruction.gear.json",
			Apl:       "ui/warlock/apls/p4/destruction.apl.json",
			Talents:   "-01-055020512000415",
			SpecOptions: &proto.Player_Warlock{
				Warlock: &proto.Warlock{
					Options: &proto.WarlockOptions{
						Armor:  proto.WarlockOptions_DemonArmor,
						Summon: proto.WarlockOptions_Imp,
					},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceHuman,
				proto.Race_RaceGnome,
				proto.Race_RaceOrc,
				proto.Race_RaceUndead,
				proto.Race_RaceTroll,
			},
		},
		{
			ClassName: "Priest",
			Class:     proto.Class_ClassPriest,
			GearSet:   "ui/shadow_priest/gear_sets/phase_4.gear.json",
			Apl:       "ui/shadow_priest/apls/phase_4.apl.json",
			Talents:   "--5022204002501251",
			SpecOptions: &proto.Player_ShadowPriest{
				ShadowPriest: &proto.ShadowPriest{
					Options: &proto.ShadowPriest_Options{
						UseShadowfiend: true,
					},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceHuman,
				proto.Race_RaceDwarf,
				proto.Race_RaceNightElf,
				proto.Race_RaceGnome,
				proto.Race_RaceUndead,
				proto.Race_RaceTroll,
			},
		},
		{
			ClassName: "Paladin",
			Class:     proto.Class_ClassPaladin,
			GearSet:   "ui/retribution_paladin/gear_sets/p4-twist.gear.json",
			Apl:       "ui/retribution_paladin/apls/p4-twist.apl.json",
			Talents:   "--532300512003151",
			SpecOptions: &proto.Player_RetributionPaladin{
				RetributionPaladin: &proto.RetributionPaladin{
					Options: &proto.PaladinOptions{
						PrimarySeal: proto.PaladinSeal_Martyrdom,
					},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceHuman,
				proto.Race_RaceDwarf,
				proto.Race_RaceUndead,
			},
		},
		{
			ClassName: "Shaman",
			Class:     proto.Class_ClassShaman,
			GearSet:   "ui/enhancement_shaman/gear_sets/phase_4_dw.gear.json",
			Apl:       "ui/enhancement_shaman/apls/phase_4.apl.json",
			Talents:   "-5005202105023051",
			SpecOptions: &proto.Player_EnhancementShaman{
				EnhancementShaman: &proto.EnhancementShaman{
					Options: &proto.EnhancementShaman_Options{},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceDwarf,
				proto.Race_RaceOrc,
				proto.Race_RaceTauren,
				proto.Race_RaceTroll,
				proto.Race_RaceSkyborneWindshaper,
			},
		},
		{
			ClassName: "Druid",
			Class:     proto.Class_ClassDruid,
			GearSet:   "ui/feral_druid/gear_sets/phase_4.gear.json",
			Apl:       "ui/feral_druid/apls/phase_4.apl.json",
			Talents:   "-550002032320211-05",
			SpecOptions: &proto.Player_FeralDruid{
				FeralDruid: &proto.FeralDruid{
					Options: &proto.FeralDruid_Options{
						LatencyMs: 100,
					},
				},
			},
			Races: []proto.Race{
				proto.Race_RaceNightElf,
				proto.Race_RaceSkyborneHighOrder,
				proto.Race_RaceTauren,
				proto.Race_RaceSkyborneWindshaper,
			},
		},
	}

	buffs := core.FullBuffsPhase4
	durations := []float64{60, 120}

	allResults := []ClassResults{}

	for _, cfg := range configs {
		fmt.Printf("\n========================================\n")
		fmt.Printf("SIMULATING CLASS: %s\n", cfg.ClassName)
		fmt.Printf("========================================\n")

		gearFile, err := os.ReadFile(cfg.GearSet)
		if err != nil {
			fmt.Printf("Failed reading gear set %s: %v\n", cfg.GearSet, err)
			continue
		}
		gearSpec := &proto.EquipmentSpec{}
		if err := (protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(gearFile, gearSpec)); err != nil {
			fmt.Printf("Failed unmarshaling gear set %s: %v\n", cfg.GearSet, err)
			continue
		}

		aplFile, err := os.ReadFile(cfg.Apl)
		if err != nil {
			fmt.Printf("Failed reading apl %s: %v\n", cfg.Apl, err)
			continue
		}
		aplSpec := &proto.APLRotation{}
		if err := (protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(aplFile, aplSpec)); err != nil {
			fmt.Printf("Failed unmarshaling apl %s: %v\n", cfg.Apl, err)
			continue
		}

		for _, d := range durations {
			fmt.Printf("\n--- Duration: %.0fs ---\n", d)
			classRes := ClassResults{
				Class:    cfg.ClassName,
				Duration: d,
				Results:  []RaceResult{},
			}

			for _, r := range cfg.Races {
				player := core.WithSpec(&proto.Player{
					Class:              cfg.Class,
					Level:              60,
					Race:               r,
					Equipment:          gearSpec,
					Consumes:           &proto.Consumes{},
					Buffs:              buffs.Player,
					TalentsString:      cfg.Talents,
					Profession1:        proto.Profession_Engineering,
					Rotation:           aplSpec,
					ReactionTimeMs:     150,
					ChannelClipDelayMs: 50,
				}, cfg.SpecOptions)

				raid := core.SinglePlayerRaidProto(player, buffs.Party, buffs.Raid, buffs.Debuffs)

				rsr := &proto.RaidSimRequest{
					Raid: raid,
					Encounter: &proto.Encounter{
						Duration: d,
						Targets: []*proto.Target{
							core.DefaultTargetProtoLvl60,
						},
					},
					SimOptions: &proto.SimOptions{
						Iterations: 400,
						IsTest:     true,
					},
				}

				result := core.RunRaidSim(rsr)
				dps := 0.0
				totgDps := 0.0
				totgCasts := 0.0
				if result.Error != nil {
					fmt.Printf("  Race %-25s: ERROR: %s\n", r.String(), result.Error.Message)
				} else {
					p := result.RaidMetrics.Parties[0].Players[0]
					dps = p.Dps.Avg
					for _, action := range p.Actions {
						if action.Id.GetSpellId() == 462103 {
							for _, t := range action.Targets {
								totgCasts += float64(t.Hits + t.Crits)
								totgDps += t.Damage / (d * 400.0)
							}
						}
					}
					if r == proto.Race_RaceUndead {
						fmt.Printf("  Race %-25s: DPS = %8.2f (TotG Procs/Fight: %4.1f, TotG DPS: %5.2f)\n", r.String(), dps, totgCasts/400.0, totgDps)
					} else {
						fmt.Printf("  Race %-25s: DPS = %8.2f\n", r.String(), dps)
					}
				}

				classRes.Results = append(classRes.Results, RaceResult{
					Race: r.String(),
					DPS:  dps,
				})
			}

			sort.Slice(classRes.Results, func(i, j int) bool {
				return classRes.Results[i].DPS > classRes.Results[j].DPS
			})

			allResults = append(allResults, classRes)
		}
	}

	jsonData, _ := json.MarshalIndent(allResults, "", "  ")
	os.WriteFile("sim_results_60_120.json", jsonData, 0644)
	fmt.Println("\nAll results saved to sim_results_60_120.json!")
}
