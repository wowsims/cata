import{m as e,l as s,n as l,o as a,q as t,s as n,R as r,U as o,V as i,W as c,X as p,Y as d,T as h,w as g,Z as u,J as I,K as m}from"./preset_utils-BxnLH_sf.chunk.js";import{R as S}from"./suggest_reforges_action-BUOXREhe.chunk.js";import{a4 as v,a5 as f,a7 as O,G as y,aw as P,ax as E,ay as A,az as T,U as L,ac as k,ad as R,ae as G,af as w,aj as C,al as F,ak as N,am as b,b as q,a as V,ah as M,a6 as x,ao as B,S as H,F as W,R as D}from"./detailed_results-oqHDwgJK.chunk.js";import{S as j}from"./inputs-DtBUrljf.chunk.js";const U=e({fieldName:"okfUptime",label:"Owlkin Frenzy Uptime (%)",labelTooltip:"Percentage of fight uptime for Owlkin Frenzy",percent:!0}),z={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:88747}}},doAtValue:{const:{val:"-5s"}}},{action:{castSpell:{spellId:{spellId:88747}}},doAtValue:{const:{val:"-4s"}}},{action:{castSpell:{spellId:{spellId:88747}}},doAtValue:{const:{val:"-3s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-2s"}}},{action:{castSpell:{spellId:{spellId:2912}}},doAtValue:{const:{val:"-2s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-1.5s"}},hide:!0},{action:{castSpell:{spellId:{spellId:5176}}},doAtValue:{const:{val:"-1.5s"}},hide:!0}],priorityList:[{action:{condition:{cmp:{op:"OpGt",lhs:{currentTime:{}},rhs:{const:{val:"2s"}}}},castSpell:{spellId:{spellId:2825,tag:-1}}}},{action:{condition:{and:{vals:[{druidCurrentEclipsePhase:{eclipsePhase:"NeutralPhase"}},{not:{val:{auraIsActive:{auraId:{spellId:61345}}}}},{cmp:{op:"OpLt",lhs:{currentTime:{}},rhs:{const:{val:"2s"}}}}]}},castSpell:{spellId:{spellId:8921}}}},{action:{castSpell:{spellId:{spellId:33831}}}},{action:{condition:{or:{vals:[{cmp:{op:"OpGe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"100"}}}}]}},castSpell:{spellId:{otherId:"OtherActionPotion"}}}},{action:{condition:{or:{vals:[{cmp:{op:"OpGe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"100"}}}},{cmp:{op:"OpGe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"100"}}}}]}},autocastOtherCooldowns:{}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"100"}}}}]}},castSpell:{spellId:{spellId:48505}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:48517}}},{or:{vals:[{and:{vals:[{cmp:{op:"OpLe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"15"}}}},{cmp:{op:"OpLe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"0"}}}}]}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:5570}}},rhs:{dotTickFrequency:{spellId:{spellId:5570}}}}}]}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:5570}}},rhs:{math:{op:"OpMul",lhs:{dotTickFrequency:{spellId:{spellId:5570}}},rhs:{const:{val:"4"}}}}}}]}},strictSequence:{actions:[{castSpell:{spellId:{spellId:5570}}},{castSpell:{spellId:{spellId:93402}}}]}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:48518}}},{or:{vals:[{and:{vals:[{cmp:{op:"OpLe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"20"}}}},{cmp:{op:"OpLe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"0"}}}}]}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:5570}}},rhs:{dotTickFrequency:{spellId:{spellId:5570}}}}}]}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:5570}}},rhs:{math:{op:"OpMul",lhs:{dotTickFrequency:{spellId:{spellId:5570}}},rhs:{const:{val:"8"}}}}}}]}},strictSequence:{actions:[{castSpell:{spellId:{spellId:5570}}},{castSpell:{spellId:{spellId:8921}}}]}}},{action:{condition:{and:{vals:[{not:{val:{druidCurrentEclipsePhase:{eclipsePhase:"SolarPhase"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"100"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"80"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"60"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"40"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"20"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"0"}}}}]}},castSpell:{spellId:{spellId:78674}}}},{action:{condition:{and:{vals:[{druidCurrentEclipsePhase:{eclipsePhase:"SolarPhase"}},{cmp:{op:"OpLe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"60"}}}}]}},castSpell:{spellId:{spellId:78674}}}},{action:{condition:{or:{vals:[{auraIsActive:{auraId:{spellId:48518}}},{auraIsActive:{auraId:{spellId:48517}}}]}},castSpell:{spellId:{spellId:78674}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:48517}}},{cmp:{op:"OpEq",lhs:{auraNumStacks:{auraId:{spellId:88747}}},rhs:{const:{val:"3"}}}}]}},castSpell:{spellId:{spellId:88751}}}},{hide:!0,action:{condition:{not:{val:{druidCurrentEclipsePhase:{eclipsePhase:"LunarPhase"}}}},castSpell:{spellId:{spellId:5176}}}},{action:{condition:{druidCurrentEclipsePhase:{eclipsePhase:"SolarPhase"}},castSpell:{spellId:{spellId:5176}}}},{action:{castSpell:{spellId:{spellId:2912}}}}]},_={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:88747}}},doAtValue:{const:{val:"-5s"}}},{action:{castSpell:{spellId:{spellId:88747}}},doAtValue:{const:{val:"-4s"}}},{action:{castSpell:{spellId:{spellId:88747}}},doAtValue:{const:{val:"-3s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-2s"}}},{action:{castSpell:{spellId:{spellId:2912}}},doAtValue:{const:{val:"-2s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-1.5s"}},hide:!0},{action:{castSpell:{spellId:{spellId:5176}}},doAtValue:{const:{val:"-1.5s"}},hide:!0}],priorityList:[{action:{condition:{cmp:{op:"OpGt",lhs:{currentTime:{}},rhs:{const:{val:"2s"}}}},castSpell:{spellId:{spellId:2825,tag:-1}}}},{action:{condition:{and:{vals:[{druidCurrentEclipsePhase:{eclipsePhase:"NeutralPhase"}},{not:{val:{auraIsActive:{auraId:{spellId:61345}}}}},{cmp:{op:"OpLt",lhs:{currentTime:{}},rhs:{const:{val:"2s"}}}}]}},castSpell:{spellId:{spellId:8921}}}},{action:{castSpell:{spellId:{spellId:33831}}}},{action:{condition:{or:{vals:[{cmp:{op:"OpGe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"100"}}}}]}},castSpell:{spellId:{otherId:"OtherActionPotion"}}}},{action:{condition:{or:{vals:[{cmp:{op:"OpGe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"100"}}}},{cmp:{op:"OpGe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"100"}}}}]}},autocastOtherCooldowns:{}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"100"}}}}]}},castSpell:{spellId:{spellId:48505}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:48517}}},{or:{vals:[{and:{vals:[{cmp:{op:"OpLe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"15"}}}},{cmp:{op:"OpGe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"1"}}}}]}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:5570}}},rhs:{dotTickFrequency:{spellId:{spellId:5570}}}}}]}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:5570}}},rhs:{math:{op:"OpMul",lhs:{dotTickFrequency:{spellId:{spellId:5570}}},rhs:{const:{val:"4"}}}}}}]}},strictSequence:{actions:[{castSpell:{spellId:{spellId:5570}}},{castSpell:{spellId:{spellId:93402}}}]}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:48518}}},{or:{vals:[{and:{vals:[{cmp:{op:"OpLe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"20"}}}},{cmp:{op:"OpGe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"1"}}}}]}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:5570}}},rhs:{dotTickFrequency:{spellId:{spellId:5570}}}}}]}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:5570}}},rhs:{math:{op:"OpMul",lhs:{dotTickFrequency:{spellId:{spellId:5570}}},rhs:{const:{val:"8"}}}}}}]}},strictSequence:{actions:[{castSpell:{spellId:{spellId:5570}}},{castSpell:{spellId:{spellId:8921}}}]}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:48518}}},{auraIsActive:{auraId:{spellId:93399}}},{or:{vals:[{cmp:{op:"OpGt",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"10"}}}},{cmp:{op:"OpEq",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"0"}}}}]}}]}},castSpell:{spellId:{spellId:78674}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:48517}}},{auraIsActive:{auraId:{spellId:93399}}},{or:{vals:[{cmp:{op:"OpGt",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"8"}}}},{cmp:{op:"OpEq",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"0"}}}}]}}]}},castSpell:{spellId:{spellId:78674}}}},{action:{condition:{and:{vals:[{not:{val:{druidCurrentEclipsePhase:{eclipsePhase:"SolarPhase"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"100"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"80"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"60"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"40"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"20"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"0"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"25"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"50"}}}},{cmp:{op:"OpNe",lhs:{currentSolarEnergy:{}},rhs:{const:{val:"75"}}}}]}},castSpell:{spellId:{spellId:78674}}}},{action:{condition:{and:{vals:[{druidCurrentEclipsePhase:{eclipsePhase:"SolarPhase"}},{cmp:{op:"OpLe",lhs:{currentLunarEnergy:{}},rhs:{const:{val:"60"}}}},{not:{val:{auraIsActive:{auraId:{spellId:48517}}}}}]}},castSpell:{spellId:{spellId:78674}}}},{action:{condition:{or:{vals:[{and:{vals:[{auraIsActive:{auraId:{spellId:48518}}},{not:{val:{auraIsActive:{auraId:{spellId:93399}}}}}]}},{and:{vals:[{auraIsActive:{auraId:{spellId:48517}}},{not:{val:{auraIsActive:{auraId:{spellId:93399}}}}}]}}]}},castSpell:{spellId:{spellId:78674}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:48517}}},{cmp:{op:"OpEq",lhs:{auraNumStacks:{auraId:{spellId:88747}}},rhs:{const:{val:"3"}}}}]}},castSpell:{spellId:{spellId:88751}}}},{hide:!0,action:{condition:{not:{val:{druidCurrentEclipsePhase:{eclipsePhase:"LunarPhase"}}}},castSpell:{spellId:{spellId:5176}}}},{action:{condition:{druidCurrentEclipsePhase:{eclipsePhase:"SolarPhase"}},castSpell:{spellId:{spellId:5176}}}},{action:{castSpell:{spellId:{spellId:2912}}}}]},K={items:[{id:65200,enchant:4207,gems:[68780,52236],reforging:141},{id:65112,reforging:162},{id:65203,enchant:4200,gems:[52207],reforging:162},{id:60232,enchant:4115,gems:[52207],reforging:162},{id:65045,enchant:4102,gems:[52207,52207],reforging:167},{id:65021,enchant:4257,gems:[0],reforging:167},{id:65199,enchant:4068,gems:[52207,0],reforging:141},{id:65374,randomSuffix:-231,gems:[52208,52207]},{id:65201,enchant:4110,gems:[52207,52236]},{id:60236,enchant:4104,gems:[52236,52207],reforging:167},{id:65123,reforging:166},{id:65373,randomSuffix:-131},{id:65105},{id:62047,reforging:167},{id:65041,enchant:4097},{id:65133,enchant:4091,reforging:134},{id:64672,gems:[52207],reforging:141}]},J={items:[{id:71497,enchant:4207,gems:[68780,52208],reforging:162},{id:71472,gems:[52207],reforging:162},{id:71450,randomSuffix:-285,enchant:4200,gems:[52208],reforging:162},{id:71434,enchant:4115,reforging:145},{id:71499,enchant:4102,gems:[52207,52207],reforging:162},{id:71463,enchant:4257,gems:[0]},{id:71496,enchant:4068,gems:[52208,0]},{id:71249,gems:[52207,52207],reforging:141},{id:71498,enchant:4110,gems:[52207,52207],reforging:145},{id:71436,enchant:4104,gems:[52208],reforging:117},{id:71217,gems:[52207],reforging:140},{id:71449,reforging:145},{id:69110},{id:62047,reforging:167},{id:71086,enchant:4097,gems:[52207,52207,52207],reforging:134},{},{id:71580,gems:[52208]}]},X={items:[{id:78696,enchant:4207,gems:[68780,71881],reforging:119},{id:78364,reforging:119},{id:78744,enchant:4200,gems:[71881,71881],reforging:145},{id:77096,enchant:4115,gems:[0],reforging:119},{id:78662,enchant:4102,gems:[71881,71881,71850],reforging:117},{id:78372,enchant:4257,gems:[71881,0],reforging:117},{id:78676,enchant:4068,gems:[71881,0]},{id:78420,gems:[71881,71881,71881],reforging:117},{id:78714,enchant:4110,gems:[71881,71881,71881]},{id:78434,enchant:4104,gems:[71881,71881],reforging:145},{id:78491,gems:[71881]},{id:78419,gems:[71881]},{id:77995},{id:77991},{id:78363,enchant:4097,gems:[71881],reforging:147},{id:78433,enchant:4091,gems:[71881],reforging:119},{id:77082,gems:[52207],reforging:147}]},Y=s("Pre-raid",{items:[{id:60282,enchant:4207,gems:[68780,52236],reforging:141},{id:69882,randomSuffix:-131},{id:60284,enchant:4200,gems:[52207],reforging:162},{id:60232,enchant:4115,gems:[52207],reforging:162},{id:60281,enchant:4102,gems:[52207,52207],reforging:145},{id:65021,enchant:4257,gems:[0],reforging:167},{id:60285,enchant:4068,gems:[52207,0],reforging:141},{id:65374,randomSuffix:-231,gems:[52208,52207]},{id:65384,randomSuffix:-193,enchant:4110,gems:[52207,52236]},{id:60236,enchant:4104,gems:[52236,52207],reforging:167},{id:65373,randomSuffix:-131},{id:70124,reforging:145},{id:65105},{id:62047,reforging:167},{id:70157,enchant:4097,reforging:167},{id:59484,enchant:4091,reforging:148},{id:70111,gems:[52207],reforging:162}]}),Z=s("T11",K),Q=s("T12",J),$=s("T13 (WIP)",X),ee=l("T11 4P",z),se=l("T12",_),le=a("Standard",v.fromMap({[f.StatIntellect]:1.3,[f.StatSpirit]:1.27,[f.StatSpellPower]:1,[f.StatHitRating]:1.27,[f.StatCritRating]:.41,[f.StatHasteRating]:.8,[f.StatMasteryRating]:.56})),ae={name:"Standard",data:O.create({talentsString:"33230221123212111001-01-020331",glyphs:y.create({prime1:P.GlyphOfInsectSwarm,prime2:P.GlyphOfMoonfire,prime3:P.GlyphOfWrath,major1:E.GlyphOfStarfall,major2:E.GlyphOfRebirth,major3:E.GlyphOfMonsoon,minor1:A.GlyphOfTyphoon,minor2:A.GlyphOfUnburdenedRebirth,minor3:A.GlyphOfMarkOfTheWild})})},te=T.create({classOptions:{innervateTarget:L.create()}}),ne=k.create({flask:R.FlaskOfTheDraconicMind,food:G.FoodSeafoodFeast,defaultPotion:w.VolcanicPotion,prepopPotion:w.VolcanicPotion}),re=C.create({arcaneBrilliance:!0,bloodlust:!0,markOfTheWild:!0,icyTalons:!0,moonkinForm:!0,leaderOfThePack:!0,powerWordFortitude:!0,strengthOfEarthTotem:!0,trueshotAura:!0,wrathOfAirTotem:!0,demonicPact:!0,blessingOfKings:!0,blessingOfMight:!0,communion:!0}),oe=F.create({vampiricTouch:!0,darkIntent:!0}),ie=N.create({}),ce=b.create({bloodFrenzy:!0,sunderArmor:!0,ebonPlaguebringer:!0,mangle:!0,criticalMass:!0,demoralizingShout:!0,frostFever:!0}),pe={distanceFromTarget:20,profession1:q.Engineering,profession2:q.Tailoring,darkIntentUptime:100},de=t("Balance T11",{gear:Z,talents:ae,rotation:ee,epWeights:le}),he=t("Balance T12",{gear:Q,talents:ae,rotation:se,epWeights:le}),ge=n(H.SpecBalanceDruid,{cssClass:"balance-druid-sim-ui",cssScheme:V.getCssClass(V.Druid),knownIssues:[],epStats:[f.StatIntellect,f.StatSpirit,f.StatSpellPower,f.StatHitRating,f.StatCritRating,f.StatHasteRating,f.StatMasteryRating],epReferenceStat:f.StatSpellPower,displayStats:M.createDisplayStatArray([f.StatHealth,f.StatMana,f.StatStamina,f.StatIntellect,f.StatSpirit,f.StatSpellPower,f.StatMasteryRating],[x.PseudoStatSpellHitPercent,x.PseudoStatSpellCritPercent,x.PseudoStatSpellHastePercent]),modifyDisplayStats:e=>{const s=e.getCurrentStats(),l=v.fromProto(s.gearStats),a=v.fromProto(s.talentsStats).subtract(l);let t=(new v).withPseudoStat(x.PseudoStatSpellCritPercent,2*e.getTalents().naturesMajesty);return t=t.addStat(f.StatHitRating,a.getPseudoStat(x.PseudoStatSpellHitPercent)*B),{talents:t}},defaults:{gear:Q.gear,epWeights:le.epWeights,statCaps:(new v).withPseudoStat(x.PseudoStatSpellHitPercent,17),consumes:ne,talents:ae.data,specOptions:te,raidBuffs:re,partyBuffs:ie,individualBuffs:oe,debuffs:ce,other:pe},playerIconInputs:[j()],includeBuffDebuffInputs:[r,o,i,c,p,d],excludeBuffDebuffInputs:[],otherInputs:{inputs:[U,h,g,u,I]},encounterPicker:{showExecuteProportion:!1},presets:{epWeights:[le],talents:[ae],rotations:[ee,se],gear:[Y,Z,Q,$],builds:[de,he]},autoRotation:e=>se.rotation.rotation,raidSimPresets:[{spec:H.SpecBalanceDruid,talents:ae.data,specOptions:te,consumes:ne,otherDefaults:pe,defaultFactionRaces:{[W.Unknown]:D.RaceUnknown,[W.Alliance]:D.RaceWorgen,[W.Horde]:D.RaceTroll},defaultGear:{[W.Unknown]:{},[W.Alliance]:{1:$.gear},[W.Horde]:{1:$.gear}}}]});class ue extends m{constructor(e,s){super(e,s,ge),s.sim.waitForInit().then((()=>{new S(this)}))}}export{ue as B};