import{k as e,$ as t,l as a,a1 as s,n,o as l,a5 as i,a6 as o,a7 as r,s as p,X as d,t as c,w as g,Z as S,T as m,G as u,J as h,K as f}from"./preset_utils-BxnLH_sf.chunk.js";import{R as I}from"./suggest_reforges_action-BUOXREhe.chunk.js";import{aO as v,aP as P,T as y,aQ as A,S as R,a4 as T,a5 as k,a6 as O,a7 as w,G as M,aR as C,aS as E,aT as D,H,ac as b,af as W,ad as x,ae as F,b as G,ah as J,aj as j,ak as B,al as L,am as N,aE as V,aF as U,aG as _,F as K,R as $}from"./detailed_results-oqHDwgJK.chunk.js";import{s as q}from"./apl_utils-BS4fVgCj.chunk.js";import{s as z,P as Q,a as X,A as Z,N as Y}from"./shared-D-W5wHZG.chunk.js";const ee={inputs:[e({fieldName:"type",label:"Type",values:[{name:"Single Target",value:v.SingleTarget},{name:"AOE",value:v.Aoe}]}),e({fieldName:"sting",label:"Sting",labelTooltip:"Maintains the selected Sting on the primary target.",values:[{name:"None",value:P.NoSting},{name:"Scorpid Sting",value:P.ScorpidSting},{name:"Serpent Sting",value:P.SerpentSting}],showWhen:e=>e.getSimpleRotation().type==v.SingleTarget}),t({fieldName:"trapWeave",label:"Trap Weave",labelTooltip:"Uses Explosive Trap at appropriate times. Note that selecting this will disable Black Arrow because they share a CD."}),t({fieldName:"multiDotSerpentSting",label:"Multi-Dot Serpent Sting",labelTooltip:"Casts Serpent Sting on multiple targets",changeEmitter:e=>y.onAny([e.rotationChangeEmitter,e.talentsChangeEmitter])})]},te={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-1s"}}}],priorityList:[{action:{autocastOtherCooldowns:{}}},{action:{castSpell:{spellId:{spellId:13812}}}},{action:{castSpell:{spellId:{spellId:2643}}}},{action:{castSpell:{spellId:{spellId:56641}}}}]},ae={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:13812}}},doAtValue:{const:{val:"-25s"}}},{action:{castSpell:{spellId:{spellId:1130}}},doAtValue:{const:{val:"-11s"}}},{action:{castSpell:{spellId:{spellId:13165}}},doAtValue:{const:{val:"-10s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-3s"}}},{action:{castSpell:{spellId:{spellId:19434}}},doAtValue:{const:{val:"-3s"}}}],priorityList:[{action:{autocastOtherCooldowns:{}}},{action:{condition:{auraIsInactiveWithReactionTime:{auraId:{spellId:3045}}},castSpell:{spellId:{spellId:3045}}}},{action:{condition:{or:{vals:[{cmp:{op:"OpLe",lhs:{spellCastTime:{spellId:{spellId:19434}}},rhs:{const:{val:"1s"}}}},{isExecutePhase:{threshold:"E90"}}]}},castSpell:{spellId:{spellId:19434}}}},{action:{condition:{or:{vals:[{not:{val:{auraIsActive:{auraId:{spellId:53221,tag:1}}}}},{cmp:{op:"OpLe",lhs:{auraRemainingTime:{auraId:{spellId:53221,tag:1}}},rhs:{const:{val:"3s"}}}}]}},castSpell:{spellId:{spellId:56641}}}},{action:{condition:{and:{vals:[{not:{val:{isExecutePhase:{threshold:"E90"}}}},{spellCanCast:{spellId:{spellId:53209}}},{spellCanCast:{spellId:{spellId:23989}}}]}},strictSequence:{actions:[{castSpell:{spellId:{spellId:53209}}},{castSpell:{spellId:{spellId:23989}}}]}}},{action:{condition:{and:{vals:[{not:{val:{dotIsActive:{spellId:{spellId:1978}}}}},{not:{val:{isExecutePhase:{threshold:"E90"}}}}]}},castSpell:{spellId:{spellId:1978}}}},{action:{condition:{not:{val:{isExecutePhase:{threshold:"E90"}}}},castSpell:{spellId:{spellId:53209}}}},{action:{condition:{spellCanCast:{spellId:{spellId:53351}}},castSpell:{spellId:{spellId:53351}}}},{action:{condition:{or:{vals:[{cmp:{op:"OpGe",lhs:{currentFocus:{}},rhs:{const:{val:"66"}}}},{cmp:{op:"OpGe",lhs:{spellTimeToReady:{spellId:{spellId:53209}}},rhs:{const:{val:"4"}}}}]}},castSpell:{spellId:{spellId:3044}}}},{action:{castSpell:{spellId:{spellId:56641}}}}]},se={items:[{id:65206,enchant:4209,gems:[68778,52209],reforging:165},{id:69880,randomSuffix:-135,reforging:151},{id:65208,enchant:4204,gems:[52212],reforging:166},{id:69884,randomSuffix:-135,enchant:1099},{id:65204,enchant:4102,gems:[52212,52220]},{id:65028,enchant:4258,gems:[52212]},{id:65205,enchant:3222,gems:[52212,52212]},{id:65132,gems:[52212,52212]},{id:60230,enchant:3823,gems:[52212,52220,52209]},{id:65063,enchant:4105,gems:[52220]},{id:65082},{id:65367,randomSuffix:-133},{id:65140},{id:65026},{id:65139,enchant:4227,reforging:167},{},{id:65058,enchant:4175,reforging:167}]},ne={items:[{id:71503,enchant:4209,gems:[68778,52209],reforging:154},{id:71610,reforging:152},{id:71403,randomSuffix:-294,enchant:4204,gems:[52258]},{id:71415,enchant:4100,gems:[52258,52258],reforging:137},{id:71501,enchant:4102,gems:[52212,52209]},{id:71561,enchant:4258,gems:[52212],reforging:154},{id:71502,enchant:4107,gems:[52212,52212],reforging:154},{id:71255,gems:[52212,52212],reforging:151},{id:71504,enchant:4126,gems:[52212,52220],reforging:154},{id:71457,enchant:4105,gems:[52212]},{id:71216,gems:[52212],reforging:152},{id:71401,reforging:152},{id:69150},{id:69112},{id:71466,enchant:4227,reforging:144},{},{id:71611,enchant:4267,reforging:154}]},le=a("MM PreRaid Preset",{items:[{id:59456,enchant:4209,gems:[68778,59478,59493]},{id:52350,gems:[52212]},{id:64712,enchant:4204,gems:[52212],reforging:152},{id:56315,enchant:1099},{id:56564,enchant:4063,gems:[52258],reforging:152},{id:63479,enchant:4071,gems:[0]},{id:64709,enchant:3222,gems:[52212,0],reforging:137},{id:56539,gems:[52212,52212],reforging:165},{id:56386,enchant:4126,gems:[52258,52258]},{id:62385,enchant:4076,gems:[52212],reforging:166},{id:52348,gems:[52212],reforging:167},{id:62362,reforging:166},{id:68709,reforging:166},{id:56328,reforging:137},{id:55066,enchant:4227,reforging:165},{},{id:59367,enchant:4175,gems:[52212],reforging:151}]}),ie=a("MM P1 Preset",se),oe=a("MM T12 Preset",ne),re=A.create({type:v.SingleTarget,sting:P.SerpentSting,trapWeave:!0,multiDotSerpentSting:!0,allowExplosiveShotDownrank:!0});s("Simple Default",R.SpecMarksmanshipHunter,re);const pe=n("MM",ae),de=n("AOE",te),ce=l("MM P1",T.fromMap({[k.StatAgility]:3.05,[k.StatRangedAttackPower]:1,[k.StatHitRating]:2.25,[k.StatCritRating]:1.39,[k.StatHasteRating]:1.33,[k.StatMasteryRating]:1.15},{[O.PseudoStatRangedDps]:6.32})),ge=l("MM P3 (T12 4-set)",T.fromMap({[k.StatAgility]:3.05,[k.StatRangedAttackPower]:1,[k.StatHitRating]:2.79,[k.StatCritRating]:1.47,[k.StatHasteRating]:.9,[k.StatMasteryRating]:1.39},{[O.PseudoStatRangedDps]:7.33})),Se={name:"Marksman",data:w.create({talentsString:"032002-2302320232120231201-03",glyphs:M.create({prime1:C.GlyphOfArcaneShot,prime2:C.GlyphOfRapidFire,prime3:C.GlyphOfSteadyShot,major1:E.GlyphOfDisengage,major2:E.GlyphOfRaptorStrike,major3:E.GlyphOfTrapLauncher})})},me=r;me.wildHunt=1,me.sharkAttack=1;const ue=D.create({classOptions:{useHuntersMark:!0,petType:H.Wolf,petTalents:me,petUptime:1}}),he=b.create({defaultPotion:W.PotionOfTheTolvir,prepopPotion:W.PotionOfTheTolvir,flask:x.FlaskOfTheWinds,defaultConjured:i.value,food:F.FoodSeafoodFeast,tinkerHands:o.value}),fe={distanceFromTarget:24,profession1:G.Engineering,profession2:G.Jewelcrafting},Ie=p(R.SpecMarksmanshipHunter,{cssClass:"marksmanship-hunter-sim-ui",cssScheme:"hunter",knownIssues:[],warnings:[],epStats:[k.StatStamina,k.StatIntellect,k.StatAgility,k.StatRangedAttackPower,k.StatHitRating,k.StatCritRating,k.StatHasteRating,k.StatMP5,k.StatMasteryRating],epPseudoStats:[O.PseudoStatRangedDps],epReferenceStat:k.StatRangedAttackPower,displayStats:J.createDisplayStatArray([k.StatHealth,k.StatStamina,k.StatAgility,k.StatRangedAttackPower,k.StatMasteryRating],[O.PseudoStatPhysicalHitPercent,O.PseudoStatPhysicalCritPercent,O.PseudoStatRangedHastePercent]),modifyDisplayStats:e=>z(e),defaults:{gear:oe.gear,epWeights:ge.epWeights,statCaps:(new T).withPseudoStat(O.PseudoStatPhysicalHitPercent,8),other:fe,consumes:he,talents:Se.data,specOptions:ue,raidBuffs:j.create({arcaneBrilliance:!0,bloodlust:!0,markOfTheWild:!0,icyTalons:!0,moonkinForm:!0,leaderOfThePack:!0,powerWordFortitude:!0,strengthOfEarthTotem:!0,trueshotAura:!0,wrathOfAirTotem:!0,demonicPact:!0,blessingOfKings:!0,blessingOfMight:!0,communion:!0}),partyBuffs:B.create({}),individualBuffs:L.create({vampiricTouch:!0}),debuffs:N.create({sunderArmor:!0,faerieFire:!0,curseOfElements:!0,savageCombat:!0,bloodFrenzy:!0})},playerIconInputs:[Q()],rotationInputs:ee,petConsumeInputs:[],includeBuffDebuffInputs:[d,c],excludeBuffDebuffInputs:[],otherInputs:{inputs:[X(),Z(),Y(),g,S,m,u,h]},encounterPicker:{showExecuteProportion:!0},presets:{epWeights:[ce,ge],talents:[Se],rotations:[pe,de],gear:[oe,le,ie]},autoRotation:e=>e.sim.encounter.targets.length>=4?de.rotation.rotation:pe.rotation.rotation,simpleRotation:(e,t,a)=>{const[s,n]=q(a),l=V.fromJsonString('{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}'),i=V.fromJsonString(`{"condition":{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"6s"}}}},"multidot":{"spellId":{"spellId":49001},"maxDots":${t.multiDotSerpentSting?3:1},"maxOverlap":{"const":{"val":"0ms"}}}}`),o=V.fromJsonString('{"condition":{"auraShouldRefresh":{"auraId":{"spellId":3043},"maxOverlap":{"const":{"val":"0ms"}}}},"castSpell":{"spellId":{"spellId":3043}}}'),r=V.fromJsonString('{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49067}}}'),p=V.fromJsonString('{"castSpell":{"spellId":{"spellId":58434}}}'),d=V.fromJsonString('{"castSpell":{"spellId":{"spellId":61006}}}'),c=V.fromJsonString('{"castSpell":{"spellId":{"spellId":49050}}}'),g=V.fromJsonString('{"castSpell":{"spellId":{"spellId":49048}}}'),S=V.fromJsonString('{"castSpell":{"spellId":{"spellId":49052}}}'),m=V.fromJsonString('{"castSpell":{"spellId":{"spellId":34490}}}'),u=V.fromJsonString('{"castSpell":{"spellId":{"spellId":53209}}}');return t.type==v.Aoe?n.push(...[l,t.sting==P.ScorpidSting?o:null,t.sting==P.SerpentSting?i:null,t.trapWeave?r:null,p].filter((e=>e))):n.push(...[l,m,d,t.sting==P.ScorpidSting?o:null,t.sting==P.SerpentSting?i:null,t.trapWeave?r:null,u,c,g,S].filter((e=>e))),U.create({prepullActions:s,priorityList:n.map((e=>_.create({action:e})))})},raidSimPresets:[{spec:R.SpecMarksmanshipHunter,talents:Se.data,specOptions:ue,consumes:he,defaultFactionRaces:{[K.Unknown]:$.RaceUnknown,[K.Alliance]:$.RaceWorgen,[K.Horde]:$.RaceTroll},defaultGear:{[K.Unknown]:{},[K.Alliance]:{1:le.gear},[K.Horde]:{1:le.gear}},otherDefaults:fe}]});class ve extends f{constructor(e,t){super(e,t,Ie),t.sim.waitForInit().then((()=>{new I(this,{getEPDefaults:e=>e.getGear().getItemSetCount("Lightning-Charged Battlegear")>=4?ce.epWeights:e.getGear().getItemSetCount("Flamewaker's Battlegear")>=4?ge.epWeights:ce.epWeights})}))}}export{ve as M};