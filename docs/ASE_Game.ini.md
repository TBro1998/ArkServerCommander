# ASE Game.ini Parameters

This file contains the list of parameters for `Game.ini` applicable to ARK: Survival Evolved, extracted from the wiki. All options must be under the `[/script/shootergame.shootergamemode]` section.

| Parameter | Type | Default | Description |
|---|---|---|---|
| `AutoPvEStartTimeSeconds` | float | 0.0 | States when the [[PvE]] mode should start in a [[PvP|PvPvE]] server. Valid values are from 0 to 86400. Options `bAutoPvETimer`, `bAutoPvEUseSystemTime` and `AutoPvEStopTimeSeconds` must also be set. '''Note:''' although at code level it is defined as a floating point number, it is suggested to use an integer instead. |
| `AutoPvEStopTimeSeconds` | float | 0.0 | States when the [[PvE]] mode should end in a [[PvP|PvPvE]] server. Valid values are from 0 to 86400. Options `bAutoPvETimer`, `bAutoPvEUseSystemTime` and `AutoPvEStopTimeSeconds` must also be set. '''Note:''' although at code level it is defined as a floating point number, it is suggested to use an integer instead. |
| `BabyCuddleGracePeriodMultiplier` | float | 1.0 | Scales how long after delaying cuddling with the Baby before Imprinting Quality starts to decrease. |
| `BabyCuddleIntervalMultiplier` | float | 1.0 | Scales how often babies needs attention for imprinting. More often means you'll need to cuddle with them more frequently to gain Imprinting Quality. Scales always according to default `BabyMatureSpeedMultiplier` value: set at 1.0 the imprint request is every 8 hours. See also [[Imprinting#The_Imprinting_formula|The Imprinting formula]] how it affects the imprinting amount at each baby care/cuddle. |
| `BabyCuddleLoseImprintQualitySpeedMultiplier` | float | 1.0 | Scales how fast Imprinting Quality decreases after the grace period if you haven't yet cuddled with the Baby. |
| `BabyFoodConsumptionSpeedMultiplier` | float | 1.0 | Scales the speed that baby tames eat their food. A lower value decreases (by percentage) the food eaten by babies. |
| `BabyImprintAmountMultiplier` | float | 1.0 | Scales the percentage each imprint provides. A higher value, will rise the amount of imprinting % at each baby care/cuddle, a lower value will decrease it. This multiplier is global, meaning it will affect the imprinting progression of every species. See also [[Imprinting#The_Imprinting_formula|The Imprinting formula]] how it affects the imprinting amount at each baby care/cuddle. |
| `BabyImprintingStatScaleMultiplier` | float | 1.0 | Scales how much of an effect on stats the Imprinting Quality has. Set it to 0 to effectively disable the system. |
| `BabyMatureSpeedMultiplier` | float | 1.0 | Scales the maturation speed of babies. A higher number decreases (by percentage) time needed for baby tames to mature. See [[Breeding#Times_for_Breeding|Times for Breeding]] tables for values at 1.0, see [[Imprinting#The_Imprinting_formula|The Imprinting formula]] how it affects the imprinting amount at each baby care/cuddle. |
| `bAllowUnclaimDinos` | boolean | True | If `False`, prevents players to unclaim tame creatures. |
| `bAllowCustomRecipes` | boolean | True | If `False`, disabled custom RP-oriented Recipe/Cooking System (including Skill-Based results). |
| `bAllowFlyerSpeedLeveling` | boolean | False | Specifies whether flyer creatures can have their [[Movement Speed]] levelled up. |
| `bAllowPlatformSaddleMultiFloors` | boolean | False | If `True`, allows multiple platform floors. |
| `bAllowUnlimitedRespecs` | boolean | False | If `True`, allows more than one usage of [[Mindwipe Tonic]] without 24 hours cooldown. |
| `BaseTemperatureMultiplier` | float | 1.0 | Specifies the map base temperature scaling factor: lower value makes the environment colder, higher value makes the environment hotter. |
| `bAutoPvETimer` | boolean | False | If `True`, enabled [[PvE]] mode in a [[PvP|PvPvE]] server at pre-specified times. The option `bAutoPvEUseSystemTime` determinates what kind of time to use, while `AutoPvEStartTimeSeconds` and `AutoPvEStopTimeSeconds` set the begin and end time of PvE mode. |
| `bAutoPvEUseSystemTime` | boolean | False | If `True`, [[PvE]] mode begin and end times in a [[PvP|PvPvE]] server will refer to the server system time instead of in-game world time. Options `bAutoPvETimer`, `AutoPvEStartTimeSeconds` and `AutoPvEStopTimeSeconds` must also be set. |
| `bAutoUnlockAllEngrams` | boolean | False | If `True`, unlocks all Engrams available. Ignores OverrideEngramEntries and OverrideNamedEngramEntries entries. |
| `bDisableDinoBreeding` | boolean | False | If `True`, prevents tames to be bred. |
| `bDisableDinoRiding` | boolean | False | If `True`, prevents players to ride tames. |
| `bDisableDinoTaming` | boolean | False | If `True`, prevents players to tame wild creatures. |
| `bDisableFriendlyFire` | boolean | False | If `True`, prevents Friendly-Fire (among tribe mates/tames/structures). |
| `bDisableLootCrates` | boolean | False | If `True`, prevents spawning of Loot crates (artifact creates will still spawn). |
| `bDisableStructurePlacementCollision` | boolean | False | If `True`, allows for structures to clip into the terrain. |
| `bFlyerPlatformAllowUnalignedDinoBasing` | boolean | False | If `True`, Quetz platforms will allow any non-allied tame to base on them when they are flying. |
| `bIgnoreStructuresPreventionVolumes` | boolean | False | If `True`, enables building areas where normally it's not allowed, such around some maps' [[Obelisk|Obelisks]], in the [[The_Ancient_Device_(Aberration)|Aberration Portal]] and in Mission Volumes areas on [[Genesis: Part 1]]. '''Note:''' in Genesis: Part 1 this settings is enabled by default and there is an ad hoc settings called `bGenesisUseStructuresPreventionVolumes` to disable it. |
| `bIncreasePvPRespawnInterval` | boolean | True | If `False`, disables [[PvP]] additional re-spawn time (`IncreasePvPRespawnIntervalBaseAmount`) that scales (`IncreasePvPRespawnIntervalMultiplier`) when a player is killed by a team within a certain amount of time (`IncreasePvPRespawnIntervalCheckPeriod`). |
| `bOnlyAllowSpecifiedEngrams` | boolean | False | If `True`, any Engram not explicitly specified by `OverrideEngramEntries` or `OverrideNamedEngramEntries` list will be hidden. All Items and Blueprints based on hidden Engrams will be removed. |
| `bPassiveDefensesDamageRiderlessDinos` | boolean | False | If `True`, allows spike walls to damage wild/riderless creatures. |
| `bPvEAllowTribeWar` | boolean | True | If `False`, disables capability for Tribes to officially declare war on each other for mutually-agreed-upon period of time. |
| `bPvEAllowTribeWarCancel` | boolean | False | If `True`, allows cancellation of an agreed-upon war before it has actually started. |
| `bPvEDisableFriendlyFire` | boolean | False | If `True`, disabled Friendly-Fire (among tribe mates/tames/structures) in [[PvE]] servers. |
| `bShowCreativeMode` | boolean | False | If `True`, enables creative mode. |
| `bUseCorpseLocator` | boolean | True | If `False`, prevents survivors to see a green light beam at the location of their dead body. |
| `bUseDinoLevelUpAnimations` | boolean | True | If `False`, tame creatures on level-up will not perform the related animation. |
| `bUseSingleplayerSettings` | boolean | False | If `True`, all game settings will be more balanced for an individual player experience. Useful for dedicated server with a very small amount of players. See [[#Single_Player_Settings|Single Player Settings]] section for more details. |
| `bUseTameLimitForStructuresOnly` | boolean | False | If `True` will make Tame Units only be applied and used for Platforms with Structures and Rafts effectively disabling Tame Units for tames without Platform Structures. |
| `CraftingSkillBonusMultiplier` | float | 1.0 | Scales the bonus received from upgrading the [[Crafting Skill]]. |
| `CraftXPMultiplier` | float | 1.0 | Scales the amount of XP earned for crafting. |
| `CropDecaySpeedMultiplier` | float | 1.0 | Scales the speed of crop decay in plots. A higher value decrease (by percentage) speed of crop decay in plots. |
| `CropGrowthSpeedMultiplier` | float | 1.0 | Scales the speed of crop growth in plots. A higher value increases (by percentage) speed of crop growth. |
| `CustomRecipeEffectivenessMultiplier` | float | 1.0 | Scales the effectiveness of custom recipes. A higher value increases (by percentage) their effectiveness. |
| `CustomRecipeSkillMultiplier` | float | 1.0 | Scales the effect of the players crafting speed level that is used as a base for the formula in creating a custom recipe. A higher number increases (by percentage) the effect. |
| `DestroyTamesOverLevelClamp` | integer | 0 | Tames that exceed that level will be deleted on server start. Official servers have it set to `450`. |
| `DinoHarvestingDamageMultiplier` | float | 3.2 | Scales the damage done to a harvestable item/entity by a tame. A higher number increases (by percentage) the speed of harvesting. |
| `DinoTurretDamageMultiplier` | float | 1.0 | Scales the damage done by Turrets towards a creature. A higher values increases it (by percentage). |
| `EggHatchSpeedMultiplier` | float | 1.0 | Scales the time needed for a fertilised egg to hatch. A higher value decreases (by percentage) that time. |
| `FastDecayInterval` | integer | 43200 | Specifies the decay period for "Fast Decay" structures (such as pillars or lone foundations). Value is in seconds. `FastDecayUnsnappedCoreStructures` in ''GameUserSettings.ini'' must be set to `True` as well to take any effect. |
| `FishingLootQualityMultiplier` | float | 1.0 | Sets the quality of items that have a quality when fishing. Valid values are from 1.0 to 5.0. |
| `FuelConsumptionIntervalMultiplier` | float | 1.0 | Defines the interval of fuel consumption. |
| `GenericXPMultiplier` | float | 1.0 | Scales the amount of XP earned for generic XP (automatic over time). |
| `GlobalCorpseDecompositionTimeMultiplier` | float | 1.0 | Scales the decomposition time of corpses, (player and creature), globally. Higher values prolong the time. |
| `GlobalItemDecompositionTimeMultiplier` | float | 1.0 | Scales the decomposition time of dropped items, loot bags etc. globally. Higher values prolong the time. |
| `GlobalPoweredBatteryDurabilityDecreasePerSecond` | float | 3.0 | Specifies the rate at which charge batteries are used in electrical objects. |
| `GlobalSpoilingTimeMultiplier` | float | 1.0 | Scales the spoiling time of perishables globally. Higher values prolong the time. |
| `HairGrowthSpeedMultiplier` | float | 1.0 | Scales the hair growth. Higher values increase speed of growth. |
| `HarvestXPMultiplier` | float | 1.0 | Scales the amount of XP earned for harvesting. |
| `IncreasePvPRespawnIntervalBaseAmount` | float | 60.0 | If `bIncreasePvPRespawnInterval` is `True`, sets the additional [[PvP]] re-spawn time in seconds that scales (`IncreasePvPRespawnIntervalMultiplier`) when a player is killed by a team within a certain amount of time (`IncreasePvPRespawnIntervalCheckPeriod`). |
| `IncreasePvPRespawnIntervalCheckPeriod` | float | 300.0 | If `bIncreasePvPRespawnInterval` is `True`, sets the amount of time in seconds within a player re-spawn time increases (`IncreasePvPRespawnIntervalBaseAmount`) and scales (`IncreasePvPRespawnIntervalMultiplier`) when it is killed by a team in [[PvP]]. |
| `IncreasePvPRespawnIntervalMultiplier` | float | 2.0 | If `bIncreasePvPRespawnInterval` is `True`, scales the [[PvP]] additional re-spawn time (`IncreasePvPRespawnIntervalBaseAmount`) when a player is killed by a team within a certain amount of time (`IncreasePvPRespawnIntervalCheckPeriod`). |
| `KillXPMultiplier` | float | 1.0 | Scale the amount of XP earned for a kill. |
| `LayEggIntervalMultiplier` | float | 1.0 | Scales the time between eggs are spawning / being laid. Higher number increases it (by percentage). |
| `LimitNonPlayerDroppedItemsCount` | integer | 0 | Limits the number of dropped items in the area defined by `LimitNonPlayerDroppedItemsRange`. Official servers have it set to `600`. |
| `LimitNonPlayerDroppedItemsRange` | integer | 0 | Sets the area range (in Unreal Units) in which the option `LimitNonPlayerDroppedItemsCount` applies. Official servers have it set to `1600`. |
| `MatingIntervalMultiplier` | float | 1.0 | Scales the interval between tames can mate. A lower value decreases it (on a percentage scale). Example: a value of 0.5 would allow tames to mate 50% sooner. |
| `MatingSpeedMultiplier` | float | 1.0 | Scales the speed at which tames mate with each other. A higher value increases it (by percentage). Example: MatingSpeedMultiplier=2.0 would cause tames to complete mating in half the normal time. |
| `MaxAlliancesPerTribe` | integer | N/A | If set, defines the maximum alliances a tribe can form or be part of. |
| `MaxFallSpeedMultiplier` | float | 1.0 | Defines the falling speed multiplier at which players starts taking fall damage. The falling speed is based on the time players spent in the air while having a negated Z axis velocity meaning that the higher this setting is, the longer players can fall without taking fall damage. For example, having it set to `0.1` means players will no longer survive a regular jump while having it set very high such as to `100.0` means players will survive a fall from the sky limit, etc. This setting doesn't affect the gravity scale of the players so there won't be any physics difference to the character movements. |
| `MaxNumberOfPlayersInTribe` | integer | 0 | Sets the maximum survivors allowed in a tribe. A value of 1 effectively disables tribes. The default value of 0 means there is no limit about how many survivors can be in a tribe. |
| `MaxTribeLogs` | integer | 400 | Sets how many Tribe log entries are displayed for each tribe. |
| `MaxTribesPerAlliance` | integer | N/A | If set, defines the maximum of tribes in an alliance. |
| `OverrideMaxExperiencePointsDino` | integer | N/A | Overrides the max XP cap of tame characters by exact specified amount. |
| `OverrideMaxExperiencePointsPlayer` | integer | N/A | Overrides the max XP cap of players characters by exact specified amount. |
| `OverridePlayerLevelEngramPoints` | integer | N/A | Configures the number of engram points granted to players for each level gained. This option must be repeated for each player level set on the server. |
| `PassiveTameIntervalMultiplier` | float | 1.0 | Scales how often a survivor get tame requests for passive tame creatures. |
| `PlayerHarvestingDamageMultiplier` | float | 1.0 | Scales the damage done to a harvestable item/entity by a Player. A higher value increases it (by percentage): the higher number, the faster the survivors collects. |
| `PoopIntervalMultiplier` | float | 1.0 | Scales how frequently survivors can poop. Higher value decreases it (by percentage) |
| `PreventBreedingForClassNames` | "<string>" | N/A | Prevents breeding of specific creatures via classname. E.g. `PreventBreedingForClassNames="Argent_Character_BP_C"`. Creature classnames can be found on the [[Creature IDs]] page. |
| `PreventDinoTameClassNames` | "<string>" | N/A | Prevents taming of specific dinosaurs via classname. E.g. `PreventDinoTameClassNames="Argent_Character_BP_C"`. Dino classnames can be found on the [[Creature IDs]] page. |
| `PreventOfflinePvPConnectionInvincibleInterval` | float | 5.0 | Specifies the time in seconds a player cannot take damages after logged-in. |
| `PreventTransferForClassNames` | "<string>" | N/A | Prevents transfer of specific creatures via classname. E.g. `PreventTransferForClassNames="Argent_Character_BP_C"`Creature classnames can be found on the [[Creature IDs]] page. |
| `PvPZoneStructureDamageMultiplier` | float | 6.0 | Specifies the scaling factor for damage structures take within caves. The lower the value, the less damage the structure takes (i.e. setting to 1.0 will make structure built in or near a cave receive the same amount of damage as those built on the surface). |
| `ResourceNoReplenishRadiusPlayers` | float | 1.0 | Controls how resources regrow closer or farther away from players. Values higher than 1.0 increase the distance around players where resources are not allowed to grow back. Values between 0 and 1.0 will reduce it. |
| `ResourceNoReplenishRadiusStructures` | float | 1.0 | Controls how resources regrow closer or farther away from structures Values higher than 1.0 increase the distance around structures where resources are not allowed to grow back. Values between 0 and 1.0 will reduce it. |
| `SpecialXPMultiplier` | float | 1.0 | Scale the amount of XP earned for SpecialEvent. |
| `StructureDamageRepairCooldown` | integer | 180 | Option for cooldown period on structure repair from the last time damaged. Set to 180 seconds by default, 0 to disable it. |
| `SupplyCrateLootQualityMultiplier` | float | 1.0 | Increases the quality of items that have a quality in the supply crates. Valid values are from 1.0 to 5.0. The quality also depends on the Difficulty Offset. |
| `TamedDinoCharacterFoodDrainMultiplier` | float | 1.0 | Scales how fast tame creatures consume food. |
| `TamedDinoTorporDrainMultiplier` | float | 1.0 | Scales how fast tamed creatures lose torpor. |
| `TribeSlotReuseCooldown` | float | 0.0 | Locks a tribe slot for the value in seconds, e.g.: a value of 3600 would mean that if a survivor leaves the tribe, their place cannot be taken by another survivor (or re-join) for 1 hour. Used on Official Small Tribes Servers. |
| `UseCorpseLifeSpanMultiplier` | float | 1.0 | Modifies corpse and dropped box lifespan. |
| `WildDinoCharacterFoodDrainMultiplier` | float | 1.0 | Scales how fast wild creatures consume food. |
| `WildDinoTorporDrainMultiplier` | float | 1.0 | Scales how fast wild creatures lose torpor. |
| `bHardLimitTurretsInRange` | boolean | False | If `True`, enables the retroactive turret hard limit (100 turrets within a 10k unit radius). |
| `bLimitTurretsInRange` | boolean | True | If `False`, doesn't limit the maximum allowed automated turrets (including [[Plant Species X]]) in a certain range. |
| `LimitTurretsNum` | integer | 100 | Determines the maximum number of turrets that are allowed in the area. |
| `LimitTurretsRange` | float | 10000.0 | Determines the area in Unreal Unit in which turrets are added towards the limit. |
| `AdjustableMutagenSpawnDelayMultiplier` | float | 1.0 | Scales the Mutagen spawn rates. By default, The game attempts to spawn them every 8 hours on dedicated servers, and every hour on non-dedicated servers and single-player. Rising this value will rise the re-spawn interval, lowering will make it shorter. |
| `BaseHexagonRewardMultiplier` | float | 1.0 | Scales the missions score hexagon rewards. Also scales token rewards in Club Ark (ASA). |
| `bDisableHexagonStore` | boolean | False | If `True`, disables the Hexagon store |
| `bDisableDefaultMapItemSets` | boolean | False | If `True`, disables Genesis 2 Tek Suit on Spawn. |
| `bDisableGenesisMissions` | boolean | False | If `True`, disables missions on Genesis. |
| `bDisableWorldBuffs` | boolean | False | If `True`, disables world effects from [[Missions (Genesis: Part 2)]] altogether. To disable specific world buffs, see `DisableWorldBuffs` of [[#DynamicConfig|#DynamicConfig]]. |
| `bEnableWorldBuffScaling` | boolean | False | If `True`, makes world effects from [[Missions (Genesis: Part 2)]] scale from server settings, rather than add/subtract a flat amount to the value at runtime. |
| `bGenesisUseStructuresPreventionVolumes` | boolean | False | If `True`, disables building in mission areas on [[Genesis: Part 1]]. |
| `bHexStoreAllowOnlyEngramTradeOption` | boolean | False | If `True`, allows only Engrams to be sold on the Hex Store, disables everything else. |
| `HexagonCostMultiplier` | float | 1.0 | Scales the hexagon cost of items in the Hexagon store. Also scales token cost of items in Club Ark (ASA). |
| `WorldBuffScalingEfficacy` | float | 1.0 | Makes world effects from [[Missions (Genesis: Part 2)]] scaling more or less effective when setting `bEnableWorldBuffScaling=True`. 1 would be default, 0.5 would be 50% less effective, 100 would be 100x more effective. |

---

### [ModInstaller]

| Parameter | Type | Default | Description |
|---|---|---|---|
| `ModIDS` | ModID | | Specifies a single Steam Workshop Mods/Maps/TC ID to download/install/update on the server. To handle multiple IDs, multiple lines must be added with the same syntax, each one with the specific workshop ID. Requires `-automanagedmods` in the command line. |