import{m as e,k as t,l as a,n as s,o as n,q as i,a7 as l,a5 as o,a6 as r,s as d,X as p,t as c,a8 as g,w as m,Z as h,T as u,G as S,J as f,K as I}from"./preset_utils-BxnLH_sf.chunk.js";import{R as v}from"./suggest_reforges_action-BUOXREhe.chunk.js";import{aO as P,aP as y,T as A,a7 as T,G as R,aR as O,aS as k,a4 as w,a5 as C,a6 as H,aq as G,aU as b,H as E,ac as W,af as x,ad as F,ae as V,b as D,a as L,ah as M,at as j,au as B,aV as U,aj as N,ak as _,al as K,am as q,aF as J,aG as z,S as X,F as Z,R as Q}from"./detailed_results-oqHDwgJK.chunk.js";import{s as Y}from"./apl_utils-BS4fVgCj.chunk.js";import{s as $,P as ee,a as te,A as ae,N as se}from"./shared-D-W5wHZG.chunk.js";const ne=e({fieldName:"sniperTrainingUptime",label:"ST Uptime (%)",labelTooltip:"Uptime for the Sniper Training talent, as a percent of the fight duration.",percent:!0,showWhen:e=>e.getTalents().sniperTraining>0,changeEmitter:e=>A.onAny([e.specOptionsChangeEmitter,e.talentsChangeEmitter])}),ie={inputs:[t({fieldName:"type",label:"Type",values:[{name:"Single Target",value:P.SingleTarget},{name:"AOE",value:P.Aoe}]}),t({fieldName:"sting",label:"Sting",labelTooltip:"Maintains the selected Sting on the primary target.",values:[{name:"None",value:y.NoSting},{name:"Serpent Sting",value:y.SerpentSting}],showWhen:e=>e.getSimpleRotation().type==P.SingleTarget})]},le={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:13165}}},doAtValue:{const:{val:"-10s"}}},{action:{castSpell:{spellId:{spellId:1130}}},doAtValue:{const:{val:"-5s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-1s"}}},{action:{castSpell:{spellId:{spellId:13812}}},doAtValue:{const:{val:"-1s"}}}],priorityList:[{action:{condition:{cmp:{op:"OpGt",lhs:{currentTime:{}},rhs:{const:{val:"1s"}}}},autocastOtherCooldowns:{}}},{action:{condition:{auraIsActive:{auraId:{spellId:77769}}},castSpell:{spellId:{spellId:13812}}}},{action:{condition:{spellIsReady:{spellId:{spellId:13812}}},castSpell:{spellId:{spellId:77769}}}},{action:{castSpell:{spellId:{spellId:2643}}}},{action:{castSpell:{spellId:{spellId:53351}}}},{action:{condition:{and:{vals:[{auraIsActive:{auraId:{spellId:56343}}},{or:{vals:[{not:{val:{dotIsActive:{spellId:{spellId:53301}}}}},{cmp:{op:"OpLt",lhs:{dotRemainingTime:{spellId:{spellId:53301}}},rhs:{math:{op:"OpAdd",lhs:{spellTravelTime:{spellId:{spellId:53301}}},rhs:{const:{val:"1s"}}}}}}]}}]}},castSpell:{spellId:{spellId:53301}}}},{action:{castSpell:{spellId:{spellId:77767}}}}]},oe={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:13812}}},doAtValue:{const:{val:"-25s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-1.4s"}}},{action:{castSpell:{spellId:{spellId:77767}}},doAtValue:{const:{val:"-1.4s"}}},{action:{castSpell:{spellId:{spellId:13165}}},doAtValue:{const:{val:"-10s"}}},{action:{castSpell:{spellId:{spellId:1130}}},doAtValue:{const:{val:"-3s"}}},{action:{castSpell:{spellId:{spellId:53517}}},doAtValue:{const:{val:"-4s"}},hide:!0},{action:{triggerIcd:{auraId:{spellId:97125}}},doAtValue:{const:{val:"-40s"}},hide:!0}],priorityList:[{action:{condition:{cmp:{op:"OpGt",lhs:{currentTime:{}},rhs:{const:{val:"3s"}}}},autocastOtherCooldowns:{}}},{action:{condition:{or:{vals:[{isExecutePhase:{threshold:"E20"}},{cmp:{op:"OpLe",lhs:{remainingTime:{}},rhs:{const:{val:"25s"}}}}]}},castSpell:{spellId:{itemId:58145}}}},{hide:!0,action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:1978}}}}},castSpell:{spellId:{spellId:2643}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:1978}}}}},castSpell:{spellId:{spellId:1978}}}},{action:{castSpell:{spellId:{spellId:53301}}}},{action:{castSpell:{spellId:{spellId:53351}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{remainingTime:{}},rhs:{const:{val:"8s"}}}},{spellIsReady:{spellId:{spellId:3674}}}]}},castSpell:{spellId:{spellId:3674}}}},{action:{condition:{cmp:{op:"OpGe",lhs:{currentFocus:{}},rhs:{const:{val:"80"}}}},castSpell:{spellId:{spellId:3044}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{currentFocus:{}},rhs:{const:{val:"40"}}}},{cmp:{op:"OpLe",lhs:{remainingTime:{}},rhs:{const:{val:"8s"}}}},{cmp:{op:"OpGe",lhs:{spellTimeToReady:{spellId:{spellId:53301}}},rhs:{const:{val:"1s"}}}}]}},castSpell:{spellId:{spellId:3044}}}},{hide:!0,action:{condition:{or:{vals:[{isExecutePhase:{threshold:"E20"}},{cmp:{op:"OpLe",lhs:{remainingTime:{}},rhs:{const:{val:"25s"}}}}]}},castSpell:{spellId:{spellId:3045}}}},{action:{castSpell:{spellId:{spellId:77767}}}}]},re={items:[{id:65206,enchant:4209,gems:[68778,52209],reforging:165},{id:69880,randomSuffix:-135,reforging:151},{id:65208,enchant:4204,gems:[52212],reforging:166},{id:69884,randomSuffix:-135,enchant:1099},{id:65204,enchant:4102,gems:[52212,52220]},{id:65028,enchant:4258,gems:[52212]},{id:65205,enchant:3222,gems:[52212,52212]},{id:65132,gems:[52212,52212]},{id:60230,enchant:3823,gems:[52212,52220,52209]},{id:65063,enchant:4105,gems:[52220]},{id:65082},{id:65367,randomSuffix:-133},{id:65140},{id:65026},{id:65139,enchant:4227,reforging:167},{},{id:65058,enchant:4175,reforging:167}]},de={items:[{id:71503,enchant:4209,gems:[68778,52209],reforging:154},{id:71610,reforging:152},{id:71403,randomSuffix:-294,enchant:4204,gems:[52258],reforging:165},{id:71415,enchant:4100,gems:[52258,52258],reforging:137},{id:71501,enchant:4102,gems:[52212,52209]},{id:71561,enchant:4258,gems:[52212],reforging:152},{id:71502,enchant:4107,gems:[52212,52212],reforging:151},{id:71255,gems:[52212,52212],reforging:152},{id:71504,enchant:4126,gems:[52212,52220],reforging:152},{id:71457,enchant:4105,gems:[52212]},{id:71216,gems:[52212],reforging:152},{id:71401,reforging:152},{id:69150},{id:69112},{id:71466,enchant:4227},{},{id:71611,enchant:4267,reforging:151}]},pe={items:[{id:78698,enchant:4209,gems:[68778,71840]},{id:77091,reforging:152},{id:78737,enchant:4204,gems:[71879,71879],reforging:154},{id:71415,enchant:4100,gems:[71879,71879],reforging:137},{id:78661,enchant:4102,gems:[71879,71879,71840],reforging:152},{id:78430,enchant:4258,gems:[71879,52212],reforging:165},{id:78362,enchant:4107,gems:[71879,71879,52212],reforging:151},{id:78447,gems:[52212,71879,71879],reforging:151},{id:78709,enchant:4126,gems:[71879,71879,71879]},{id:78415,enchant:4105,gems:[71879,71879],reforging:151},{id:78413,gems:[71879],reforging:154},{id:77111,gems:[71879],reforging:152},{id:77994},{id:77999},{id:78473,enchant:4227},{},{id:78471,enchant:4267}]},ce=a("Pre-raid",{items:[{id:60303,enchant:4209,gems:[68778,52209],reforging:165},{id:69880,randomSuffix:-135,reforging:151},{id:60306,enchant:4204,gems:[52212],reforging:166},{id:69884,randomSuffix:-135,enchant:4087,reforging:151},{id:60304,enchant:4102,gems:[52212,52220],reforging:166},{id:65028,enchant:4071,gems:[0],reforging:138},{id:60307,enchant:3222,gems:[52212,0]},{id:71255,gems:[52212,52212],reforging:152},{id:60230,enchant:4126,gems:[52258,52258,52212],reforging:165},{id:70123,enchant:4076,gems:[52258],reforging:152},{id:70105},{id:65367,randomSuffix:-135,reforging:151},{id:69001,reforging:166},{id:65140},{id:70165,enchant:4227,reforging:165},{},{id:71077,enchant:4267,gems:[52212],reforging:165}]}),ge=a("P2",re),me=a("P3",de),he=a("P4",pe),ue=s("SV",oe),Se=s("AOE",le),fe={name:"Survival",data:T.create({talentsString:"03-2302-03203203023022121311",glyphs:R.create({prime1:O.GlyphOfExplosiveShot,prime2:O.GlyphOfKillShot,prime3:O.GlyphOfArcaneShot,major1:k.GlyphOfDisengage,major2:k.GlyphOfRaptorStrike,major3:k.GlyphOfTrapLauncher})})},Ie=n("P2",w.fromMap({[C.StatStamina]:.5,[C.StatAgility]:3.27,[C.StatRangedAttackPower]:1,[C.StatHitRating]:2.16,[C.StatCritRating]:1.17,[C.StatHasteRating]:.89,[C.StatMasteryRating]:.88},{[H.PseudoStatRangedDps]:3.75})),ve=n("P3",w.fromMap({[C.StatStamina]:.5,[C.StatAgility]:3.37,[C.StatRangedAttackPower]:1,[C.StatHitRating]:2.56,[C.StatCritRating]:1.27,[C.StatHasteRating]:1.09,[C.StatMasteryRating]:1.04},{[H.PseudoStatRangedDps]:4.16})),Pe=n("P4",w.fromMap({[C.StatStamina]:.5,[C.StatAgility]:3.47,[C.StatRangedAttackPower]:1,[C.StatHitRating]:2.56,[C.StatCritRating]:1.45,[C.StatHasteRating]:1.09,[C.StatMasteryRating]:1.04},{[H.PseudoStatRangedDps]:4.16})),ye=i("P2",{gear:ge,epWeights:Ie,talents:fe,rotationType:G.TypeAuto}),Ae=i("P3",{gear:me,epWeights:ve,talents:fe,rotationType:G.TypeAuto}),Te=i("P4",{gear:he,epWeights:Pe,talents:fe,rotationType:G.TypeAuto}),Re=b.create({classOptions:{useHuntersMark:!0,petType:E.Wolf,petTalents:l,petUptime:1},sniperTrainingUptime:.9}),Oe=W.create({defaultPotion:x.PotionOfTheTolvir,prepopPotion:x.PotionOfTheTolvir,flask:F.FlaskOfTheWinds,defaultConjured:o.value,food:V.FoodSeafoodFeast,tinkerHands:r.value}),ke={distanceFromTarget:24,profession1:D.Engineering,profession2:D.Jewelcrafting},we=d(X.SpecSurvivalHunter,{cssClass:"survival-hunter-sim-ui",cssScheme:L.getCssClass(L.Hunter),knownIssues:[],warnings:[],epStats:[C.StatStamina,C.StatAgility,C.StatRangedAttackPower,C.StatHitRating,C.StatCritRating,C.StatHasteRating,C.StatMasteryRating],epPseudoStats:[H.PseudoStatRangedDps],epReferenceStat:C.StatRangedAttackPower,displayStats:M.createDisplayStatArray([C.StatHealth,C.StatStamina,C.StatAgility,C.StatRangedAttackPower,C.StatMasteryRating],[H.PseudoStatPhysicalHitPercent,H.PseudoStatPhysicalCritPercent,H.PseudoStatRangedHastePercent]),modifyDisplayStats:e=>$(e),defaults:{gear:me.gear,epWeights:ve.epWeights,statCaps:(new w).withPseudoStat(H.PseudoStatPhysicalHitPercent,8),softCapBreakpoints:[j.fromPseudoStat(H.PseudoStatRangedHastePercent,{breakpoints:[20],capType:B.TypeSoftCap,postCapEPs:[.89*U]})],other:ke,consumes:Oe,talents:fe.data,specOptions:Re,raidBuffs:N.create({arcaneBrilliance:!0,bloodlust:!0,markOfTheWild:!0,icyTalons:!0,moonkinForm:!0,leaderOfThePack:!0,powerWordFortitude:!0,strengthOfEarthTotem:!0,trueshotAura:!0,wrathOfAirTotem:!0,demonicPact:!0,blessingOfKings:!0,blessingOfMight:!0,communion:!0}),partyBuffs:_.create({}),individualBuffs:K.create({vampiricTouch:!0}),debuffs:q.create({sunderArmor:!0,curseOfElements:!0,savageCombat:!0,bloodFrenzy:!0})},playerIconInputs:[ee()],rotationInputs:ie,petConsumeInputs:[],includeBuffDebuffInputs:[p,c,g],excludeBuffDebuffInputs:[],otherInputs:{inputs:[te(),ae(),se(),ne,m,h,u,S,f]},encounterPicker:{showExecuteProportion:!1},presets:{epWeights:[Ie,ve,Pe],talents:[fe],rotations:[ue,Se],builds:[ye,Ae,Te],gear:[ce,ge,me,he]},autoRotation:e=>e.sim.encounter.targets.length>=4?Se.rotation.rotation:ue.rotation.rotation,simpleRotation:(e,t,a)=>{const[s,n]=Y(a);return J.create({prepullActions:s,priorityList:n.map((e=>z.create({action:e})))})},raidSimPresets:[{spec:X.SpecSurvivalHunter,talents:fe.data,specOptions:Re,consumes:Oe,defaultFactionRaces:{[Z.Unknown]:Q.RaceUnknown,[Z.Alliance]:Q.RaceWorgen,[Z.Horde]:Q.RaceTroll},defaultGear:{[Z.Unknown]:{},[Z.Alliance]:{1:ce.gear},[Z.Horde]:{1:ce.gear}},otherDefaults:ke}]});class Ce extends I{constructor(e,t){super(e,t,we),t.sim.waitForInit().then((()=>{new v(this,{getEPDefaults:e=>e.getGear().getItemSetCount("Lightning-Charged Battlegear")>=4?Ie.epWeights:(e.getGear().getItemSetCount("Flamewaker's Battlegear"),ve.epWeights)})}))}}export{Ce as S};