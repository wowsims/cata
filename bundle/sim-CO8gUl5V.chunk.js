import{l as e,n as s,o as a,s as l,t,w as n,T as i,G as r,K as o}from"./preset_utils-BxnLH_sf.chunk.js";import{R as p}from"./suggest_reforges_action-BUOXREhe.chunk.js";import{a4 as d,a5 as c,a6 as I,a7 as m,G as h,bQ as g,bR as u,bS as v,bT as f,ac as S,ad as T,ae as O,af as y,ag as P,b as R,a as A,ah as E,ai as G,aj as x,ak as w,al as k,am as C,S as L,F as b,R as H}from"./detailed_results-oqHDwgJK.chunk.js";import{S as B,a as K}from"./inputs-B_4R5LN7.chunk.js";const W={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:2457}}},doAtValue:{const:{val:"-2s"}}},{action:{move:{rangeFromTarget:{const:{val:"9"}}}},doAtValue:{const:{val:"-2s"}}},{action:{castSpell:{spellId:{spellId:6673}}},doAtValue:{const:{val:"-1.5s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-1s"}}},{action:{castSpell:{spellId:{spellId:100}}},doAtValue:{const:{val:"-0.2s"}}}],priorityList:[{hide:!0,action:{castSpell:{spellId:{spellId:100}}}},{action:{autocastOtherCooldowns:{}}},{action:{condition:{or:{vals:[{and:{vals:[{isExecutePhase:{threshold:"E20"}},{or:{vals:[{spellIsReady:{spellId:{itemId:69113}}},{spellIsReady:{spellId:{itemId:68972}}}]}}]}},{and:{vals:[{isExecutePhase:{threshold:"E20"}},{not:{val:{or:{vals:[{spellIsKnown:{spellId:{itemId:69113}}},{spellIsKnown:{spellId:{itemId:68972}}}]}}}},{auraIsActive:{auraId:{spellId:57519}}}]}},{cmp:{op:"OpLt",lhs:{remainingTime:{}},rhs:{const:{val:"26.5s"}}}}]}},castSpell:{spellId:{otherId:"OtherActionPotion"}}}},{action:{condition:{or:{vals:[{and:{vals:[{cmp:{op:"OpEq",lhs:{auraNumStacks:{auraId:{spellId:96923}}},rhs:{const:{val:"5"}}}},{or:{vals:[{spellIsReady:{spellId:{itemId:69113}}},{spellIsReady:{spellId:{itemId:68972}}}]}},{or:{vals:[{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}},{cmp:{op:"OpGt",lhs:{spellTimeToReady:{spellId:{spellId:86346}}},rhs:{const:{val:"6s"}}}}]}},{or:{vals:[{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"125s"}}}},{isExecutePhase:{threshold:"E20"}}]}}]}},{cmp:{op:"OpLe",lhs:{remainingTime:{}},rhs:{const:{val:"16.5s"}}}},{and:{vals:[{cmp:{op:"OpGe",lhs:{currentTime:{}},rhs:{const:{val:"1.5s"}}}},{spellIsReady:{spellId:{spellId:1719}}}]}}]}},castSpell:{spellId:{spellId:33697}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpEq",lhs:{auraNumStacks:{auraId:{spellId:96923}}},rhs:{const:{val:"5"}}}},{or:{vals:[{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}},{cmp:{op:"OpGt",lhs:{spellTimeToReady:{spellId:{spellId:86346}}},rhs:{const:{val:"6s"}}}},{cmp:{op:"OpLe",lhs:{remainingTime:{}},rhs:{const:{val:"16.5s"}}}}]}},{or:{vals:[{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"125s"}}}},{isExecutePhase:{threshold:"E20"}}]}}]}},castSpell:{spellId:{itemId:69113}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpEq",lhs:{auraNumStacks:{auraId:{spellId:96923}}},rhs:{const:{val:"5"}}}},{or:{vals:[{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}},{cmp:{op:"OpGt",lhs:{spellTimeToReady:{spellId:{spellId:86346}}},rhs:{const:{val:"6s"}}}},{cmp:{op:"OpLe",lhs:{remainingTime:{}},rhs:{const:{val:"16.5s"}}}}]}},{or:{vals:[{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"125s"}}}},{isExecutePhase:{threshold:"E20"}}]}}]}},castSpell:{spellId:{itemId:68972}}}},{action:{condition:{or:{vals:[{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"120s"}}}},{isExecutePhase:{threshold:"E20"}}]}},castSpell:{spellId:{itemId:59461}}}},{action:{condition:{or:{vals:[{and:{vals:[{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"130s"}}}},{cmp:{op:"OpGe",lhs:{currentTime:{}},rhs:{const:{val:"1.5s"}}}}]}},{isExecutePhase:{threshold:"E20"}}]}},castSpell:{spellId:{itemId:62464}}}},{action:{condition:{or:{vals:[{and:{vals:[{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"130s"}}}},{cmp:{op:"OpGe",lhs:{currentTime:{}},rhs:{const:{val:"1.5s"}}}}]}},{isExecutePhase:{threshold:"E20"}}]}},castSpell:{spellId:{itemId:62469}}}},{action:{condition:{or:{vals:[{and:{vals:[{cmp:{op:"OpGe",lhs:{currentTime:{}},rhs:{const:{val:"1.5s"}}}},{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"90s"}}}}]}},{isExecutePhase:{threshold:"E20"}}]}},castSpell:{spellId:{itemId:77116}}}},{action:{condition:{or:{vals:[{and:{vals:[{cmp:{op:"OpGe",lhs:{currentTime:{}},rhs:{const:{val:"1.5s"}}}},{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"65s"}}}}]}},{isExecutePhase:{threshold:"E20"}}]}},castSpell:{spellId:{itemId:69002}}}},{action:{condition:{or:{vals:[{and:{vals:[{cmp:{op:"OpGe",lhs:{currentTime:{}},rhs:{const:{val:"1.5s"}}}},{not:{val:{or:{vals:[{spellIsKnown:{spellId:{itemId:69113}}},{spellIsKnown:{spellId:{itemId:68972}}}]}}}}]}},{or:{vals:[{spellIsKnown:{spellId:{itemId:69113}}},{spellIsKnown:{spellId:{itemId:68972}}}]}}]}},castSpell:{spellId:{spellId:82174}}}},{action:{condition:{cmp:{op:"OpGt",lhs:{numberTargets:{}},rhs:{const:{val:"1"}}}},castSpell:{spellId:{spellId:46924}}}},{action:{condition:{cmp:{op:"OpGt",lhs:{numberTargets:{}},rhs:{const:{val:"1"}}}},castSpell:{spellId:{spellId:12328}}}},{action:{condition:{const:{val:"false"}},castSpell:{spellId:{spellId:64382}}}},{action:{condition:{or:{vals:[{and:{vals:[{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"125s"}}}},{not:{val:{or:{vals:[{spellIsKnown:{spellId:{itemId:69113}}},{spellIsKnown:{spellId:{itemId:68972}}},{cmp:{op:"OpEq",lhs:{numEquippedStatProcTrinkets:{statType1:6,statType2:-1,statType3:-1,minIcdSeconds:110}},rhs:{const:{val:"1"}}}}]}}}}]}},{and:{vals:[{cmp:{op:"OpLt",lhs:{remainingTime:{}},rhs:{const:{val:"125s"}}}},{isExecutePhase:{threshold:"E20"}},{not:{val:{or:{vals:[{spellIsKnown:{spellId:{itemId:69113}}},{spellIsKnown:{spellId:{itemId:68972}}},{cmp:{op:"OpEq",lhs:{numEquippedStatProcTrinkets:{statType1:6,statType2:-1,statType3:-1,minIcdSeconds:110}},rhs:{const:{val:"1"}}}}]}}}}]}},{anyTrinketStatProcsActive:{statType1:6,statType2:-1,statType3:-1,minIcdSeconds:110}}]}},castSpell:{spellId:{spellId:33697}}}},{action:{condition:{cmp:{op:"OpGe",lhs:{currentTime:{}},rhs:{const:{val:"1.5s"}}}},castSpell:{spellId:{spellId:1719}}}},{action:{condition:{or:{vals:[{and:{vals:[{cmp:{op:"OpGe",lhs:{currentTime:{}},rhs:{const:{val:"1.5s"}}}},{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"122s"}}}},{cmp:{op:"OpLt",lhs:{spellTimeToReady:{spellId:{spellId:78}}},rhs:{const:{val:"1.5s"}}}}]}},{and:{vals:[{cmp:{op:"OpLt",lhs:{remainingTime:{}},rhs:{const:{val:"122s"}}}},{isExecutePhase:{threshold:"E20"}},{or:{vals:[{cmp:{op:"OpLt",lhs:{spellTimeToReady:{spellId:{spellId:78}}},rhs:{const:{val:"1.5s"}}}},{cmp:{op:"OpLt",lhs:{remainingTime:{}},rhs:{const:{val:"11s"}}}}]}}]}}]}},castSpell:{spellId:{spellId:85730}}}},{action:{castSpell:{spellId:{spellId:1134}}}},{action:{condition:{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}},castSpell:{spellId:{spellId:6544}}}},{action:{condition:{and:{vals:[{not:{val:{dotIsActive:{spellId:{spellId:772}}}}},{not:{val:{auraIsActive:{auraId:{spellId:2457}}}}}]}},strictSequence:{actions:[{castSpell:{spellId:{spellId:2457}}},{castSpell:{spellId:{spellId:772}}}]}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:772}}}}},castSpell:{spellId:{spellId:772}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:60503}}},{not:{val:{isExecutePhase:{threshold:"E20"}}}},{cmp:{op:"OpLe",lhs:{auraRemainingTime:{auraId:{spellId:60503}}},rhs:{const:{val:"5s"}}}}]}},castSpell:{spellId:{spellId:2457}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:60503}}},{not:{val:{isExecutePhase:{threshold:"E20"}}}},{cmp:{op:"OpLe",lhs:{auraRemainingTime:{auraId:{spellId:60503}}},rhs:{const:{val:"5s"}}}}]}},castSpell:{spellId:{spellId:7384}}}},{action:{condition:{or:{vals:[{spellCanCast:{spellId:{spellId:12294}}},{and:{vals:[{spellCanCast:{spellId:{spellId:86346}}},{not:{val:{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}}}}]}}]}},castSpell:{spellId:{spellId:2458}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:99233}}},{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"6s"}}}}]}},castSpell:{spellId:{spellId:6673}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"1.5s"}}}},{not:{val:{and:{vals:[{isExecutePhase:{threshold:"E20"}},{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}},{auraIsActive:{auraId:{spellId:57519}}}]}}}}]}},castSpell:{spellId:{spellId:12294}}}},{hide:!0,action:{condition:{and:{vals:[{cmp:{op:"OpLt",lhs:{currentRage:{}},rhs:{const:{val:"75"}}}},{cmp:{op:"OpGe",lhs:{autoTimeToNext:{}},rhs:{const:{val:"2.5s"}}}},{not:{val:{or:{vals:[{auraIsActive:{auraId:{spellId:2825,tag:-1}}},{auraIsActive:{auraId:{spellId:85730}}},{auraIsActive:{auraId:{spellId:1719}}}]}}}},{spellIsReady:{spellId:{spellId:100}}}]}},move:{rangeFromTarget:{const:{val:"9"}}}}},{action:{condition:{and:{vals:[{not:{val:{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}}}},{cmp:{op:"OpGt",lhs:{remainingTime:{}},rhs:{const:{val:"3s"}}}}]}},castSpell:{spellId:{spellId:86346}}}},{action:{condition:{or:{vals:[{or:{vals:[{cmp:{op:"OpGt",lhs:{currentRage:{}},rhs:{const:{val:"75"}}}},{and:{vals:[{auraIsActive:{auraId:{spellId:86627}}},{cmp:{op:"OpGe",lhs:{currentRage:{}},rhs:{const:{val:"65"}}}}]}}]}},{auraIsActive:{auraId:{spellId:85730}}},{auraIsActiveWithReactionTime:{auraId:{spellId:12964}}},{and:{vals:[{cmp:{op:"OpLt",lhs:{spellTimeToReady:{spellId:{spellId:85730}}},rhs:{const:{val:"1s"}}}},{not:{val:{cmp:{op:"OpLt",lhs:{spellTimeToReady:{spellId:{spellId:1719}}},rhs:{const:{val:"1s"}}}}}},{cmp:{op:"OpGt",lhs:{currentRage:{}},rhs:{const:{val:"30"}}}},{not:{val:{cmp:{op:"OpLt",lhs:{remainingTime:{}},rhs:{const:{val:"123s"}}}}}}]}},{and:{vals:[{cmp:{op:"OpGe",lhs:{currentRage:{}},rhs:{const:{val:"50"}}}},{or:{vals:[{cmp:{op:"OpGe",lhs:{spellTimeToReady:{spellId:{spellId:12294}}},rhs:{const:{val:"1.5s"}}}},{auraIsActive:{auraId:{spellId:86627}}}]}},{isExecutePhase:{threshold:"E20"}}]}},{and:{vals:[{cmp:{op:"OpLt",lhs:{remainingTime:{}},rhs:{const:{val:"2s"}}}},{cmp:{op:"OpGe",lhs:{currentRage:{}},rhs:{const:{val:"40"}}}}]}}]}},castSpell:{spellId:{spellId:78}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:60503}}},{gcdIsReady:{}},{not:{val:{spellCanCast:{spellId:{spellId:12294}}}}},{not:{val:{and:{vals:[{spellCanCast:{spellId:{spellId:86346}}},{cmp:{op:"OpLt",lhs:{auraRemainingTime:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}},rhs:{const:{val:"1s"}}}}]}}}},{cmp:{op:"OpGe",lhs:{currentRage:{}},rhs:{const:{val:"5"}}}},{or:{vals:[{not:{val:{isExecutePhase:{threshold:"E20"}}}},{and:{vals:[{not:{val:{spellCanCast:{spellId:{spellId:86346}}}}},{auraIsInactiveWithReactionTime:{auraId:{spellId:85730}}},{auraIsInactiveWithReactionTime:{auraId:{spellId:1719}}},{cmp:{op:"OpLt",lhs:{currentRage:{}},rhs:{const:{val:"50"}}}}]}}]}},{auraIsInactiveWithReactionTime:{auraId:{spellId:85730}}}]}},castSpell:{spellId:{spellId:2457}}}},{action:{condition:{auraIsActive:{auraId:{spellId:60503}}},castSpell:{spellId:{spellId:7384}}}},{action:{condition:{or:{vals:[{not:{val:{auraIsActive:{auraId:{spellId:60503}}}}},{isExecutePhase:{threshold:"E20"}}]}},castSpell:{spellId:{spellId:2458}}}},{action:{castSpell:{spellId:{spellId:5308}}}},{action:{condition:{cmp:{op:"OpLe",lhs:{auraRemainingTime:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:86346}}},rhs:{const:{val:"1.5s"}}}},castSpell:{spellId:{spellId:86346}}}},{action:{castSpell:{spellId:{spellId:1464}}}},{action:{castSpell:{spellId:{spellId:18499}}}},{action:{condition:{not:{val:{auraIsKnown:{auraId:{spellId:99233}}}}},castSpell:{spellId:{spellId:6673}}}}]},D={items:[{id:65266,enchant:4208,gems:[68779,52222],reforging:152},{id:69885,randomSuffix:-120,reforging:140},{id:65268,enchant:4202,gems:[52206],reforging:152},{id:69879,randomSuffix:-120,enchant:4100,reforging:139},{id:65264,enchant:4102,gems:[52206,52213],reforging:151},{id:60228,enchant:4256,gems:[52222,52206],reforging:154},{id:65265,enchant:4106,gems:[52206,52206],reforging:152},{id:65040,gems:[52206,52206],reforging:140},{id:65379,randomSuffix:-173,enchant:4126,gems:[52206,52213],reforging:139},{id:65075,enchant:4094,gems:[52206],reforging:166},{id:60226,gems:[52222]},{id:65382,randomSuffix:-120,reforging:140},{id:65072,reforging:159},{id:59461},{id:65003,enchant:4099,reforging:140},{},{id:60210,reforging:159}]},U={items:[{id:65266,enchant:4208,gems:[68779,52222],reforging:152},{id:69885,randomSuffix:-120,reforging:139},{id:65268,enchant:4202,gems:[52206],reforging:152},{id:65117,enchant:4100,reforging:154},{id:65264,enchant:4102,gems:[52206,52213],reforging:154},{id:60228,enchant:4256,gems:[52222,52206],reforging:151},{id:65071,enchant:4106,gems:[52206,52206],reforging:151},{id:65040,gems:[52206,52206],reforging:139},{id:65267,enchant:4126,gems:[52206,52206],reforging:166},{id:65075,enchant:4104,gems:[52206],reforging:166},{id:60226,gems:[52222]},{id:65382,randomSuffix:-120,reforging:140},{id:65072,reforging:159},{id:59461},{id:65003,enchant:4099,reforging:140},{},{id:60210,reforging:166}]},F={items:[{id:71430,enchant:4208,gems:[68779,52213],reforging:165},{id:71446,reforging:139},{id:71603,enchant:4202,gems:[52206],reforging:154},{id:69879,randomSuffix:-120,enchant:4100,reforging:140},{id:71600,enchant:4102,gems:[52222,52222]},{id:71418,enchant:4256,gems:[52206],reforging:154},{id:71601,enchant:4106,gems:[52222,52206],reforging:166},{id:71443,gems:[52206,52206],reforging:166},{id:71602,enchant:4126,gems:[52206,52206],reforging:161},{id:71404,enchant:4104,gems:[52206],reforging:153},{id:71215,gems:[52206]},{id:71433,gems:[52206],reforging:159},{id:69113},{id:69167},{id:70723,enchant:4099,gems:[52206,52206],reforging:168},{},{id:71593,reforging:137}]},M=e("Preraid",{items:[{id:60325,enchant:4208,gems:[68779,52222],reforging:152},{id:69885,randomSuffix:-122,reforging:154},{id:60327,enchant:4202,gems:[52206],reforging:152},{id:69879,randomSuffix:-120,enchant:4100,reforging:140},{id:71068,enchant:4102,gems:[52222,52222],reforging:140},{id:60228,enchant:4256,gems:[52222,52206],reforging:151},{id:71069,enchant:4106,gems:[52206,52206],reforging:166},{id:65369,randomSuffix:-222,gems:[52206,52206],reforging:153},{id:71071,enchant:4126,gems:[52206,52206],reforging:161},{id:69946,enchant:4094,gems:[52206],reforging:139},{id:71208,reforging:151},{id:60226,gems:[52206]},{id:65072,reforging:159},{id:59461},{id:63679,enchant:4099,reforging:140},{},{id:71154,reforging:137}]}),j=e("P1 - BIS",D),q=e("P1 - Realistic",U),V=e("P3 - BIS",F),N=s("Default",W),_=a("Default",d.fromMap({[c.StatStrength]:2.21,[c.StatAgility]:1.12,[c.StatAttackPower]:1,[c.StatExpertiseRating]:1.75,[c.StatHitRating]:2.77,[c.StatCritRating]:1.45,[c.StatHasteRating]:.68,[c.StatMasteryRating]:.89},{[I.PseudoStatMainHandDps]:9.22,[I.PseudoStatOffHandDps]:0})),z={name:"Default",data:m.create({talentsString:"30220303120212312211-0322-3",glyphs:h.create({prime1:g.GlyphOfMortalStrike,prime2:g.GlyphOfOverpower,prime3:g.GlyphOfSlam,major1:u.GlyphOfCleaving,major2:u.GlyphOfSweepingStrikes,major3:u.GlyphOfThunderClap,minor1:v.GlyphOfBerserkerRage,minor2:v.GlyphOfCommand,minor3:v.GlyphOfBattle})})},Q=f.create({classOptions:{startingRage:0}}),J=S.create({flask:T.FlaskOfTitanicStrength,food:O.FoodBeerBasedCrocolisk,defaultPotion:y.GolembloodPotion,prepopPotion:y.GolembloodPotion,tinkerHands:P.TinkerHandsSynapseSprings}),X={profession1:R.Engineering,profession2:R.Blacksmithing,distanceFromTarget:9},Y=l(L.SpecArmsWarrior,{cssClass:"arms-warrior-sim-ui",cssScheme:A.getCssClass(A.Warrior),knownIssues:[],epStats:[c.StatStrength,c.StatAgility,c.StatAttackPower,c.StatExpertiseRating,c.StatHitRating,c.StatCritRating,c.StatHasteRating,c.StatMasteryRating],epPseudoStats:[I.PseudoStatMainHandDps,I.PseudoStatOffHandDps],epReferenceStat:c.StatAttackPower,displayStats:E.createDisplayStatArray([c.StatHealth,c.StatStamina,c.StatStrength,c.StatAgility,c.StatAttackPower,c.StatExpertiseRating,c.StatMasteryRating],[I.PseudoStatPhysicalHitPercent,I.PseudoStatPhysicalCritPercent,I.PseudoStatMeleeHastePercent]),defaults:{gear:V.gear,epWeights:_.epWeights,statCaps:(()=>{const e=(new d).withPseudoStat(I.PseudoStatPhysicalHitPercent,8),s=(new d).withStat(c.StatExpertiseRating,26*G);return e.add(s)})(),other:X,consumes:J,talents:z.data,specOptions:Q,raidBuffs:x.create({arcaneBrilliance:!0,bloodlust:!0,markOfTheWild:!0,icyTalons:!0,moonkinForm:!0,leaderOfThePack:!0,powerWordFortitude:!0,strengthOfEarthTotem:!0,trueshotAura:!0,wrathOfAirTotem:!0,demonicPact:!0,blessingOfKings:!0,blessingOfMight:!0,communion:!0}),partyBuffs:w.create({}),individualBuffs:k.create({}),debuffs:C.create({bloodFrenzy:!0,mangle:!0,sunderArmor:!0,curseOfWeakness:!0,ebonPlaguebringer:!0})},playerIconInputs:[],includeBuffDebuffInputs:[t],excludeBuffDebuffInputs:[],otherInputs:{inputs:[B(),K(),n,i,r]},encounterPicker:{showExecuteProportion:!0},presets:{epWeights:[_],talents:[z],rotations:[N],gear:[M,j,q,V]},autoRotation:e=>N.rotation.rotation,raidSimPresets:[{spec:L.SpecArmsWarrior,talents:z.data,specOptions:Q,consumes:J,defaultFactionRaces:{[b.Unknown]:H.RaceUnknown,[b.Alliance]:H.RaceWorgen,[b.Horde]:H.RaceOrc},defaultGear:{[b.Unknown]:{},[b.Alliance]:{1:V.gear},[b.Horde]:{1:V.gear}},otherDefaults:X}]});class Z extends o{constructor(e,s){super(e,s,Y),s.sim.waitForInit().then((()=>{new p(this)}))}}export{Z as A};