# ASE GameUserSetting.ini Parameters

This file contains the list of parameters for `GameUserSettings.ini` applicable to ARK: Survival Evolved, extracted from the wiki.

## [ServerSettings]

| Parameter | Type | Default | Description |
|---|---|---|---|
| `ActiveMods` | list of mod IDs, comma-separated with no spaces, in a single line (for example: `ModID1,ModID2,ModID3`) | | Specifies the order and which mods are loaded. ModIDs are comma separated and in one line. Priority is in descending order (the left-most ModID hast the highest priority). Alternatively, but not suggested, the command line `?GameModIds{{=}}` can be used. |
| `ActiveTotalConversion` | ModID | | Used to specify a total conversion mod (e.g.: [[Primitive Plus]]). Alternatively, the command line option `-TotalConversionMod{{=}}<ModID>` can also be set. |
| `AdminLogging` | boolean | False | If `True`, logs all admin commands to in-game chat. |
| `AllowAnyoneBabyImprintCuddle` | boolean | False | If `True`, allows anyone to "take care" of a baby creatures (cuddle etc.), not just whomever imprinted on it. |
| `AllowCaveBuildingPvE` | boolean | False | If `True`, allows building in caves when [[PvE]] mode is also enabled. '''Note:''' no more working in command-line options before patch [[241.5]]. |
| `AllowCaveBuildingPvP` | boolean | True | If `False`, prevents building in caves when [[PvP]] mode is also enabled. |
| `AllowCrateSpawnsOnTopOfStructures` | boolean | False | If `True`, allows from-the-air Supply Crates to appear on top of Structures, rather than being prevented by Structures. |
| `AllowFlyerCarryPvE` | boolean | False | If `True`, allows flying creatures to pick up wild creatures in [[PvE]]. |
| `AllowFlyingStaminaRecovery` | boolean | False | If `True`, allows server to recover Stamina when standing on a Flyer. |
| `AllowHideDamageSourceFromLogs` | boolean | True | If `False`, shows the damage sources in tribe logs. |
| `AllowHitMarkers` | boolean | True | If `False`, disables optional markers for ranged attacks. |
| `AllowIntegratedSPlusStructures` | boolean | True | if `False`, disables all of the new S+ structures (intended mainly for letting unofficial servers that want to keep using the S+ mod version to keep using that without a ton of extra duplicate structures). |
| `AllowMultipleAttachedC4` | boolean | False | If `True`, allows to attach more than one [[C4]] per creature. |
| `AllowRaidDinoFeeding` | boolean | False | If `True`, allows Titanosaurs to be permanently tamed (namely allow them to be fed). '''Note:''' in The Island only spawns a maximum of 3 Titanosaurs, so 3 tamed ones should ultimately block any more ones from spawning. |
| `AllowSharedConnections` | boolean | False | If `True`, allows family sharing players to connect to the server. |
| `AllowTekSuitPowersInGenesis` | boolean | False | If `True`, enables TEK suit powers in [[Genesis: Part 1]]. |
| `AllowThirdPersonPlayer` | boolean | True | If `False`, disables third person camera allowed by default on all dedicated servers. |
| `AlwaysAllowStructurePickup` | boolean | False | If `True` disables the timer on the quick pick-up system. |
| `AlwaysNotifyPlayerLeft` | boolean | False | If `True`, players will always get notified if someone leaves the server |
| `AutoDestroyDecayedDinos` | boolean | False | If `True`, auto-destroys claimable decayed tames on load, rather than have them remain around as claimable. '''Note:''' after patch [[273.691]], in [[PvE]] mode the tame auto-unclaim after decay period has been disabled in official [[PvE]]. |
| `AutoDestroyOldStructuresMultiplier` | float | 0.0 | Allows auto-destruction of structures only after sufficient "no nearby tribe" time has passed (defined as a multiplier of the Allow Claim period). To enable it, set it to 1.0. Useful for servers to clear off abandoned structures automatically over time. Requires `-AutoDestroyStructures` command line option to work. The value scales with each structure decay time. |
| `AutoSavePeriodMinutes` | float | 15.0 | Set interval for automatic saves. Setting this to 0 will cause constant saving. |
| `BanListURL` | string with a URL | | Sets the global ban list. Must be enclosed in double quotes. The list is fetched every 10 minutes (to check if there are new banned IDs).<br/>{{Gamelink|ase}}: Official ban list URL is [http://arkdedicated.com/banlist.txt http://arkdedicated.com/banlist.txt] (before [[279.233]] the URL was [http://playark.com/banlist.txt http://playark.com/banlist.txt]). '''Note:''' it supports the HTTP protocol only (HTTPS is not supported).<br/>{{Gamelink|ase}}: Official ban list URL is [https://cdn2.arkdedicated.com/asa/BanList.txt https://cdn2.arkdedicated.com/asa/BanList.txt]. '''Note:''' it supports both HTTP and HTTPS protocols. |
| `bForceCanRideFliers` | boolean | | If `True`, allows flyers to be used on maps where they normally are disabled. '''Note:''' if you set it to `False` it will disable flyers on any map. |
| `ClampItemSpoilingTimes` | boolean | False | If `True`, clamps all spoiling times to the items' maximum spoiling times. Useful if any infinite-spoiling exploits were used on the server and you wish to clean them up. Could potentially cause issues with mods that alter spoiling time, hence it is an option. |
| `ClampItemStats` | boolean | False | If `True`, enables stats clamping for items. See [[#ItemStatClamps|ItemStatClamps]] for more info. |
| `ClampResourceHarvestDamage` | boolean | False | If `True`, limit the damage caused by a tame to a resource on harvesting based on resource remaining health.  '''Note:''' enabling this setting may result in sensible resource harvesting reduction using high damage tools or creatures. |
| `CustomDynamicConfigUrl` | string with a URL | | Direct link to a live dynamicconfig.ini file (http://arkdedicated.com/dynamicconfig.ini), allowing live changes of the supported options without the need of server restart, as well defining custom colour set for wild creatures' spawns. '''Note:''' requires `-UseDynamicConfig` command line option, `string with a URL` must use the ''HTTP'' protocol (''HTTPS'' is not supported) and inside quotes if used in command line. Check the [[#DynamicConfig|DynamicConfig]] section for the supported settings. |
| `CustomLiveTuningUrl` | string with a URL | | Direct link to the live tuning file. For more information on how to use this system check out the official announcement: https://survivetheark.com/index.php?/forums/topic/569366-server-configuration-live-tuning-system.<br/>{{Gamelink|ase}}: official servers use [http://arkdedicated.com/DefaultOverloads.json http://arkdedicated.com/DefaultOverloads.json]. '''Note:''' `string with a URL` must use the ''HTTP'' protocol (''HTTPS'' is not supported) and inside quotes if used in command line.<br/>{{Gamelink|asa}}: official servers use [https://cdn2.arkdedicated.com/asa/livetuningoverloads.json https://cdn2.arkdedicated.com/asa/livetuningoverloads.json]. '''Note:''' `string with a URL` must use inside quotes if used in command line. |
| `DayCycleSpeedScale` | float | 1.0 | Specifies the scaling factor for the passage of time in the ARK, controlling how often day changes to night and night changes to day. The default value `1` provides the same cycle speed as the single player experience (and the official public servers). Values lower than 1 slow down the cycle; higher values accelerate it. Base time when value is 1 appears to be 1-minute real time equals approx. 28-minutes game time. Thus, for an approximate 24-hour day/night cycle in game, use .035 for the value. |
| `DayTimeSpeedScale` | float | 1.0 | Specifies the scaling factor for the passage of time in the ARK during the day. This value determines the length of each day, relative to the length of each night (as specified by `NightTimeSpeedScale`). Lowering this value increases the length of each day. |
| `DestroyUnconnectedWaterPipes` | boolean | False | If `True`, after two days real-time the pipes will auto-destroy if unconnected to any non-pipe (directly or indirectly) and no allied player is nearby. |
| `DifficultyOffset` | float | 1.0 | Specifies the [[Difficulty|difficulty level]]. |
| `DinoCharacterFoodDrainMultiplier` | float | 1.0 | Specifies the scaling factor for creatures' food consumption. Higher values increase food consumption (creatures get hungry faster). It also affects the taming-times. |
| `DinoCharacterHealthRecoveryMultiplier` | float | 1.0 | Specifies the scaling factor for creatures' health recovery. Higher values increase the recovery rate (creatures heal faster). |
| `DinoCharacterStaminaDrainMultiplier` | float | 1.0 | Specifies the scaling factor for creatures' stamina consumption. Higher values increase stamina consumption (creatures get tired faster). |
| `DinoCountMultiplier` | float | 1.0 | Specifies the scaling factor for creature spawns. Higher values increase the number of creatures spawned throughout the ARK. |
| `DinoDamageMultiplier` | float | 1.0 | Specifies the scaling factor for the damage wild creatures deal with their attacks. The default value `1` provides normal damage. Higher values increase damage. Lower values decrease it. |
| `DinoResistanceMultiplier` | float | 1.0 | Specifies the scaling factor for the resistance to damage wild creatures receive when attacked. The default value `1` provides normal damage. Higher values decrease resistance, increasing damage per attack. Lower values increase it, reducing damage per attack. A value of 0.5 results in a creature taking half damage while a value of 2.0 would result in a creature taking double normal damage. |
| `DisableDinoDecayPvE` | boolean | False | If `True`, disables the creature decay in [[PvE]] mode. '''Note:''' after patch [[273.691]], in [[PvE]] mode the creature auto-unclaim after decay period has been disabled. |
| `DisableImprintDinoBuff` | boolean | False | If `True`, disables the creature imprinting player Stat Bonus. Where whomever specifically imprinted on the creature, and raised it to have an Imprinting Quality, gets extra Damage/Resistance buff. |
| `DisablePvEGamma` | boolean | False | If `True`, prevents use of console command "gamma" in [[PvE]] mode. |
| `DisableStructureDecayPvE` | boolean | False | If `True`, disables the gradual [[Building#Auto-decay|auto-decay]] of player structures. |
| `DisableWeatherFog` | boolean | False | If `True`, disables fog. |
| `DontAlwaysNotifyPlayerJoined` | boolean | False | If `True`, globally disables player joins notifications. |
| `EnableExtraStructurePreventionVolumes` | boolean | False | If `True`, disables building in specific resource-rich areas, in particular setup on [[The Island]] around the major mountains. |
| `EnablePvPGamma` | boolean | False | If `True`, allows use of console command "gamma" in [[PvP]] mode. |
| `ExtinctionEventTimeInterval` | seconds | | Used to enable the ''extinction mode'' ([[Game_Modes#ARKpocalypse_(PC)|ARKpocalypse]]). The number is the time in seconds. Use 2592000 value for 30 days. |
| `FastDecayUnsnappedCoreStructures` | boolean | False | If `True`, unsnapped foundations/pillars/fences/Tek Dedicated Storage will decay after the time stated by `FastDecayInterval` in ''Game.ini'' (default is 12 hours). Before [[259.0]], it set the decay time for such structures straight to 5 times faster. |
| `ForceAllStructureLocking` | boolean | False | If `True`, will default lock all structures. |
| `globalVoiceChat` | boolean | False | If `True`, voice chat turns global. |
| `HarvestAmountMultiplier` | float | 1.0 | Specifies the scaling factor for yields from all harvesting activities (chopping down trees, picking berries, carving carcasses, mining rocks, etc.). Higher values increase the amount of materials harvested with each strike. |
| `HarvestHealthMultiplier` | float | 1.0 | Specifies the scaling factor for the "health" of items that can be harvested (trees, rocks, carcasses, etc.). Higher values increase the amount of damage (i.e., "number of strikes") such objects can withstand before being destroyed, which results in higher overall harvest yields. |
| `IgnoreLimitMaxStructuresInRangeTypeFlag` | boolean | False | If `True`, removes the limit of 150 decorative structures (flags, signs, dermis etc.). |
| `ItemStackSizeMultiplier` | float | 1.0 | Allow increasing or decreasing global item stack size, this means all default stack sizes will be multiplied by the value given (excluding items that have a stack size of 1 by default). |
| `KickIdlePlayersPeriod` | float | 3600.0 | Time in seconds after which characters that have not moved or interacted will be kicked (if -EnableIdlePlayerKick as command line parameter is set). '''Note:''' although at code level it is defined as a floating-point number, it is suggested to use an integer instead. |
| `MaxGateFrameOnSaddles` | integer | 0 | Defines the maximum amount of gateways allowed on platform saddles. A value of 2 would prevent players from placing more than 2 gateways on their platform saddles (used in Official [[PvP]] servers). This setting is not retroactive, meaning existing builds won't be affected. Set to 0 to not allow players to place gateways on platform saddles (used in Official [[PvE]] servers). |
| `MaxHexagonsPerCharacter` | integer | 2000000000 | Sets the max amount of Hexagon a Character can accumulate. Official set it to 2500000. |
| `MaxPersonalTamedDinos` | integer | 0 | Sets a per-tribe creature tame limit (500 on official PvE servers, 300 in official PvP servers). The default value of 0 disables such limit. |
| `MaxPlatformSaddleStructureLimit` | integer | 75 | Changes the maximum number of platformed-creatures/rafts allowed on the ARK (a potential performance cost). Example: `MaxPlatformSaddleStructureLimit{{=}}10` would only allow 10 platform saddles/rafts across the entire ARK. |
| `MaxTamedDinos` | float | 5000.0 | Sets the maximum number of tame creatures on a server, this is a global cap. '''Note:''' although at code level it is defined as a floating-point number, it is suggested to use an integer instead. |
| `MaxTributeCharacters` | integer | 10 | Slots for uploaded characters. Any value less than default will be reverted. '''Note:''' rising it may corrupt player/cluster data and lead to lose of all stored characters. |
| `MaxTributeDinos` | integer | 20 | Slots for uploaded creatures. Any value less than default will be reverted. '''Note:''' Some player claimed maximum 273 to be safe cap and more will corrupt profile/cluster and lead to lose of all stored creatures but it need to be checked |
| `MaxTributeItems` | integer | 50 | Slots for uploaded items and resources. Any value less than default will be reverted. '''Note:''' Some player claimed maximum 154 to be safe cap and more will corrupt profile/cluster and lead to lose of all stored items and resources but it need to be checked |
| `NightTimeSpeedScale` | float | 1.0 | Specifies the scaling factor for the passage of time in the ARK during night time. This value determines the length of each night, relative to the length of each day (as specified by `DayTimeSpeedScale`) Lowering this value increases the length of each night. |
| `NonPermanentDiseases` | boolean | False | If `True`, makes permanent diseases not permanent. Players will lose them if on re-spawn. |
| `NPCNetworkStasisRangeScalePlayerCountStart` | integer | 0 | Minimum number of online players when the NPC Network Stasis Range Scale override is enabled (requires inputting into INI, not command line). Used to override the NPC Network Stasis Range Scale (to scale server performance for more player), on default is set to 0, disabling it. Official set it to 24. |
| `NPCNetworkStasisRangeScalePlayerCountEnd` | integer | 0 | Maximum number of online players when `NPCNetworkStasisRangeScalePercentEnd` is reached (requires inputting into INI, not command line). Used to override the NPC Network Stasis Range Scale (to scale server performance for more player), on default is set to 0, disabling it. Official set it to 70 |
| `NPCNetworkStasisRangeScalePercentEnd` | float | 0.55000001 | The Maximum scale percentage used when `NPCNetworkStasisRangeScalePlayerCountEnd` is reached (requires inputting into INI, not command line). Used to override the NPC Network Stasis Range Scale (to scale server performance for more player). Official set it to 0.5. |
| `OnlyAutoDestroyCoreStructures` | boolean | False | If `True`, prevents any non-core/non-foundation structures from auto-destroying (however they'll still get auto-destroyed if a floor that they're on gets auto-destroyed). Official [[PvE]] Servers used this option. |
| `OnlyDecayUnsnappedCoreStructures` | boolean | False | If `True`, only unsnapped core structures will decay. Useful for eliminating lone pillar/foundation spam. |
| `OverrideOfficialDifficulty` | float | 0.0 | Allows you to override the default server difficulty level of 4 with 5 to match the new official server difficulty level. Default value of 0.0 disables the override. A value of 5.0 will allow common creatures to spawn up to level 150. Originally ([[247.95]]) available only as command line option. |
| `OverrideStructurePlatformPrevention` | boolean | False | If `True`, [[Auto Turret|turrets]] becomes be buildable and functional on [[:Category:Platform Saddles|platform saddles]]. Since [[247.999]] applies on [[Wooden Spike Wall|spike]] structure too. '''Note:''' despite patch notes, in ''ShooterGameServer'' it's coded ''OverrideStructurePlatformPrevention'' with two ''r''. |
| `OxygenSwimSpeedStatMultiplier` | float | 1.0 | Use this to set how swim speed is multiplied by level spent in oxygen. The value was reduced by 80% in [[256.0]]. |
| `PerPlatformMaxStructuresMultiplier` | float | 1.0 | Higher value increases (from a percentage scale) max number of items place-able on saddles and rafts. |
| `PersonalTamedDinosSaddleStructureCost` | integer | 0 | Determines the amount of "tame creature slots" a platform saddle (with structures) will use towards the tribe tame creature limit. |
| `PlatformSaddleBuildAreaBoundsMultiplier` | float | 1.0 | Increasing the number allows structures being placed further away from the platform. |
| `PlayerCharacterFoodDrainMultiplier` | float | 1.0 | Specifies the scaling factor for player characters' food consumption. Higher values increase food consumption (player characters get hungry faster). |
| `PlayerCharacterHealthRecoveryMultiplier` | float | 1.0 | Specifies the scaling factor for player characters' health recovery. Higher values increase the recovery rate (player characters heal faster). |
| `PlayerCharacterStaminaDrainMultiplier` | float | 1.0 | Specifies the scaling factor for player characters' stamina consumption. Higher values increase stamina consumption (player characters get tired faster). |
| `PlayerCharacterWaterDrainMultiplier` | float | 1.0 | Specifies the scaling factor for player characters' water consumption. Higher values increase water consumption (player characters get thirsty faster). |
| `PlayerDamageMultiplier` | float | 1.0 | Specifies the scaling factor for the damage players deal with their attacks. The default value `1` provides normal damage. Higher values increase damage. Lower values decrease it. |
| `PlayerResistanceMultiplier` | float | 1.0 | Specifies the scaling factor for the resistance to damage players receive when attacked. The default value `1` provides normal damage. Higher values decrease resistance, increasing damage per attack. Lower values increase it, reducing damage per attack. A value of 0.5 results in a player taking half damage while a value of 2.0 would result in taking double normal damage. |
| `PreventDiseases` | boolean | False | If `True`, completely diseases on the server. Thus far just [[Swamp Fever]]. |
| `PreventMateBoost` | boolean | False | If `True`, disables creature mate boosting. |
| `PreventOfflinePvP` | boolean | False | If `True`, enables the Offline Raiding Prevention (ORP). When all tribe members are logged off, tribe characters, creature and structures become invulnerable. Creature starvation still applies, moreover, characters and creature can still die if drowned. Despite the name, it works on both [[PvE]] and [[PvP]] game modes. Due to performance reason, it is recommended to set a minimum interval with `PreventOfflinePvPInterval` option before ORP becomes active. ORP also helps lowering memory and CPU usage on a server. Enabled by default on Official [[PvE]] since [[258.3]] |
| `PreventOfflinePvPInterval` | float | 0.0 | Seconds to wait before a ORP becomes active for tribe/players and relative creatures/structures (10 seconds in official PvE servers). '''Note:''' although at code level it is defined as a floating-point number, it is suggested to use an integer instead. |
| `PreventSpawnAnimations` | boolean | False | If `True`, player characters (re)spawn without the wake-up animation. |
| `PreventTribeAlliances` | boolean | False | If `True`, prevents tribes from creating Alliances. |
| `ProximityChat` | boolean | False | If `True`, only players near each other can see their chat messages |
| `PvEAllowStructuresAtSupplyDrops` | boolean | False | If `True`, allows building near supply drop points in [[PvE]] mode. |
| `PvEDinoDecayPeriodMultiplier` | float | 1.0 | Creature [[PvE]] auto-decay time multiplier. Requires `DisableDinoDecayPvE{{=}}false` in GameUserSettings.ini or `?DisableDinoDecayPvE{{=}}false` in command line to work. |
| `PvEStructureDecayPeriodMultiplier` | float | 1.0 | Specifies the scaling factor for [[Building#Auto-decay| structures decay times]], e.g.: setting it at 2.0 will double all structure decay times, while setting at 0.5 will halve the timers. '''Note:''' despite the name, works in both [[PvP]] and [[PvE]] modes when structure decay is enabled. |
| `PvPDinoDecay` | boolean | False | If `True`, enables creatures' decay in [[PvP]] while the Offline Raid Prevention is active. |
| `PvPStructureDecay` | boolean | False | If `True`, enables structures decay on [[PvP]] servers while the Offline Raid Prevention is active. |
| `RaidDinoCharacterFoodDrainMultiplier` | float | 1.0 | Affects how quickly the food drains on such "raid dinos" (e.g.: Titanosaurus) |
| `RandomSupplyCratePoints` | boolean | False | If `True`, supply drops are in random locations. '''Note:''' This setting is known to cause artifacts becoming inaccessible on [[Ragnarok]] if active. |
| `RCONEnabled` | boolean | False | If `True`, enables RCON, needs `RCONPort{{=}}<TCP_PORT>` and `ServerAdminPassword{{=}}<admin_password>` to work. |
| `RCONPort` | integer | 27020 | Specifies the optional TCP RCON Port. See [[Dedicated_server_setup#Network|Dedicated server setup]] |
| `RCONServerGameLogBuffer` | float | 600.0 | Determines how many lines of game logs are send over the RCON. '''Note:''' despite being coded as a float it's suggested to treat it as integer. |
| `ResourcesRespawnPeriodMultiplier` | float | 1.0 | Specifies the scaling factor for the re-spawn rate for resource nodes (trees, rocks, bushes, etc.). Lower values cause nodes to re-spawn more frequently. |
| `ServerAdminPassword` | string | | If specified, players must provide this password (via the in-game console) to gain access to administrator commands on the server. '''Note:''' no quotes are used. |
| `ServerAutoForceRespawnWildDinosInterval` | float | 0.0 | Force re-spawn of all wild creatures on server restart afters the value set in seconds. Default value of 0.0 disables it. Useful to prevent certain creature species (like the Basilo and Spino) from becoming depopulated on long running servers. |
| `ServerCrosshair` | boolean | True | If `False`, disables the Crosshair on your server. |
| `ServerForceNoHUD` | boolean | False | If `True`, HUD is always disabled for non-tribe owned NPCs. |
| `ServerHardcore` | boolean | False | If `True`, enables [[Hardcore]] mode (player characters revert to level 1 upon death) |
| `ServerPassword` | string | | If specified, players must provide this password to join the server. '''Note:''' no quotes are used. |
| `serverPVE` | boolean | False | If `True`, disables [[PvP]] and enables [[PvE]] |
| `ShowFloatingDamageText` | boolean | False | If `True`, enables RPG-style popup damage text mode. |
| `ShowMapPlayerLocation` | boolean | True | If `False`, hides each player their own precise position when they view their map. |
| `SpectatorPassword` | string | | To use non-admin spectator, the server must specify a spectator password. Then any client can use these console commands: `requestspectator <password>` and `stopspectating`. '''Note:''' no quotes are used. |
| `StructureDamageMultiplier` | float | 1.0 | Specifies the scaling factor for the damage structures deal with their attacks (i.e., spiked walls). Higher values increase damage. Lower values decrease it. |
| `StructurePickupHoldDuration` | float | 0.5 | Specifies the quick pick-up hold duration, a value of `0` results in instant pick-up. |
| `StructurePickupTimeAfterPlacement` | float | 30.0 | Amount of time in seconds after placement that quick pick-up is available. |
| `StructurePreventResourceRadiusMultiplier` | float | 1.0 | Same as `ResourceNoReplenishRadiusStructures` in ''Game.ini''. If both settings are set both multiplier will be applied. Can be useful when cannot change the ''Game.ini'' file as it works as a command line option too. |
| `StructureResistanceMultiplier` | float | 1.0 | Specifies the scaling factor for the resistance to damage structures receive when attacked. The default value `1` provides normal damage. Higher values decrease resistance, increasing damage per attack. Lower values increase it, reducing damage per attack. A value of 0.5 results in a structure taking half damage while a value of 2.0 would result in a structure taking double normal damage. |
| `TamedDinoDamageMultiplier` | float | 1.0 | Specifies the scaling factor for the damage tame creatures deal with their attacks. The default value `1` provides normal damage. Higher values increase damage. Lower values decrease it. |
| `TamedDinoResistanceMultiplier` | float | 1.0 | Specifies the scaling factor for the resistance to damage tame creatures receive when attacked. The default value `1` provides normal damage. Higher values decrease resistance, increasing damage per attack. Lower values increase it, reducing damage per attack. A value of 0.5 results in a structure taking half damage while a value of 2.0 would result in a structure taking double normal damage. |
| `TamingSpeedMultiplier` | float | 1.0 | Specifies the scaling factor for creature taming speed. Higher values make taming faster. |
| `TheMaxStructuresInRange` | integer | 10500 | Specifies the maximum number of structures that can be constructed within a certain (currently hard-coded) range. Replaces the old value `NewMaxStructuresInRange`'' |
| `TribeLogDestroyedEnemyStructures` | boolean | False | By default, enemy structure destruction (for the victim tribe) is not displayed in the tribe Logs, set this to true to enable it. |
| `TribeNameChangeCooldown` | float | 15.0 | Cool-down, in minutes, in between tribe name changes. Official server use a value of 172800.0 (2 days). |
| `UseFjordurTraversalBuff` | boolean | False | If `True`, enables the biome teleport in [[Fjordur]] when holding {{Key|R}} (enabled in official PvE servers). |
| `UseOptimizedHarvestingHealth` | boolean | False | If `True`, enables a server harvesting optimization with high `HarvestAmountMultiplier` (but less rare items). '''Note:''' on {{ItemLink|ARK: Survival Evolved}} it's suggested to enable this option if harvesting with Tek Stryder causes lag spikes. |
| `XPMultiplier` | float | 1.0 | Specifies the scaling factor for the experience received by players, tribes and tames for various actions. The default value `1` provides the same amounts of experience as in the single player experience (and official public servers). Higher values increase XP amounts awarded for various actions; lower values decrease it. In [[313.5]] an additional hardcoded multiplier of 4 was activated. |
| `CrossARKAllowForeignDinoDownloads` | boolean | False | If `True`, enables non-native creatures tribute download on [[Aberration]]. |
| `MinimumDinoReuploadInterval` | float | 0.0 | Number of seconds cool-down between allowed creature re-uploads (43200 on official Servers which is 12 hours). |
| `noTributeDownloads` | boolean | False | If `True`, prevents CrossArk-data downloads in[[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]]. |
| `PreventDownloadDinos` | boolean | False | If `True`, prevents creatures download from ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]]. |
| `PreventDownloadItems` | boolean | False | If `True`, prevents items download from ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]]. |
| `PreventDownloadSurvivors` | boolean | False | If `True`, prevents survivors download from ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]]. |
| `PreventUploadDinos` | boolean | False | If `True`, prevents creatures upload to ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]]. |
| `PreventUploadItems` | boolean | False | If `True`, prevents items upload to ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]]. |
| `PreventUploadSurvivors` | boolean | False | If `True`, prevents survivors upload to ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]]. |
| `TributeCharacterExpirationSeconds` | integer | 0 | Set in seconds the expiration timer for uploaded survivors in ARK Data. With default or negative values there is no expiration time. Check [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]] for more details. '''Warning''': do not set this option to an insane high value, like more than 31536000 seconds (which is 1 year), as in ARK Data routines this is summed to upload time in Unix Epoch time format. Using really high values will result in overflow and may cause upload time checks to fail and ARK Data deleted. Finally, it is highly suggested to use the same value across all servers in a cluster, otherwise accessing ARK Data from a server with a lower value will result data deletion if in that server the timer is expired. |
| `TributeDinoExpirationSeconds` | integer | 86400 | Set in seconds the expiration timer for uploaded tames in ARK Data. If set to 0 or less will revert to default. Check [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]] for more details. '''Warning''': do not set this option to an insane high value, like more than 31536000 seconds (which is 1 year), as in ARK Data routines this is summed to upload time in Unix Epoch time format. Using really high values will result in overflow and may cause upload time checks to fail and ARK Data deleted. Finally, it is highly suggested to use the same value across all servers in a cluster, otherwise accessing ARK Data from a server with a lower value will result data deletion if in that server the timer is expired. |
| `TributeItemExpirationSeconds` | integer | 86400 | Set in seconds the expiration timer for uploaded items in ARK Data. If set to 0 or less will revert to default. Check [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]] for more details. '''Warning''': do not set this option to an insane high value, like more than 31536000 seconds (which is 1 year), as in ARK Data routines this is summed to upload time in Unix Epoch time format. Using really high values will result in overflow and may cause upload time checks to fail and ARK Data deleted. Finally, it is highly suggested to use the same value across all servers in a cluster, otherwise accessing ARK Data from a server with a lower value will result data deletion if in that server the timer is expired. |
| `CryopodNerfDamageMult` | float | 0.0099999998 | Reduces the amount of damage dealt by the creature after it is deployed from the cryopod, as a percentage of total damage output, and for the length of time set by `CryopodNerfDuration`. `CryopodNerfDuration` needs to be set as well. `CryopodNerfDamageMult{{=}}0.01` means 99% of the damage is removed.  On official it is set to 0.1. |
| `CryopodNerfDuration` | integer | 0.0 | Amount of time, in seconds, Cryosickness lasts after deploying a creature from a Cryopod. '''Note:''' although at code level it is defined as a floating-point number, it is suggested to use an integer instead. On official it is set to 10.0. |
| `CryopodNerfIncomingDamageMultPercent` | float | 0.0 | Increases the amount of damage taken by the creature after it is deployed from the cryopod, as a percentage of total damage received, and for the length of time set by `CryopodNerfDuration`. `CryopodNerfIncomingDamageMultPercent{{=}}0.25` means a released tame takes 25% more damage while the debuff lasts. On official it is set to 0.25. |
| `EnableCryopodNerf` | boolean | False | If `True`, there is no Cryopod cooldown timer, and creatures do not become unconscious. If this option is set, than `EnableCryopodNerf` and `CryopodNerfIncomingDamageMultPercent` must be set as well or they will default to 0. |
| `EnableCryoSicknessPVE` | boolean | False | If `True`, enables Cryopod cooldown timer when deploying a creature. |
| `BadWordListURL` | string with a URL | "http://arkdedicated.com/badwords.txt" | Add the `URL` to hosting your own bad words list. '''Note:''' on {{Gamelink|ase}} servers only the HTTP protocol is supported (an HTTPS URL will not work). |
| `BadWordWhiteListURL` | string with a URL | "http://arkdedicated.com/goodwords.txt" | Add the `URL` to hosting your own good words list. '''Note:''' on {{Gamelink|ase}} servers only the HTTP protocol is supported (an HTTPS URL will not work). |
| `bFilterCharacterNames` | boolean | False | If `True`, filters out character names based on the bad words/good words list. |
| `bFilterChat` | boolean | False | If `True`, filters out character names based on the bad word/good words list. |
| `bFilterTribeNames` | boolean | False | If `True`, filters out tribe names based on the badwords/goodwords list. |
| `AllowedCheatersURL` | string with a URL | N/A | Alternative to ''AllowedCheaterSteamIDs.txt'' (see [[#Administrator_Whitelisting|Administrator Whitelisting]]) using a web resource (e.g.: official uses http://arkdedicated.com/globaladmins.txt). The interval at which the server queries the resource to check for admin list update is defined by `UpdateAllowedCheatersInterval` '''Note:''' it supports the HTTP protocol only (HTTPS is not supported). Undocumented by Wildcard. |
| `ChatLogFileSplitIntervalSeconds` | integer | 86400 | Controls how to split the chat log file related to time in seconds. Cannot be set to a value lower than 45 (will default to 45 if the value is lower). Set to 0 only in official. Undocumented by Wildcard. |
| `ChatLogFlushIntervalSeconds` | integer | 86400 | Controls in how many second the chat log is flushed to log file. Cannot be set to a value lower than 15 (will default to 15 if the value is lower). Set to 0 only in official. Undocumented by Wildcard. |
| `ChatLogMaxAgeInDays` | integer | 5 | Controls how many days the chat log is long. Set it to a negative value will result it to set at -1 (virtually infinite). Set to 0 only in official. Undocumented by Wildcard. |
| `DontRestoreBackup` | boolean | False | If `True` and `-DisableDupeLogDeletes` is present, prevents the server to automatically restore a backup in case of corrupted save. Undocumented by Wildcard. |
| `EnableAFKKickPlayerCountPercent` | float | 0.0 | Enables the idle timeout to be applied only if the amount of online players reaches percentage value related to `MaxPlayers` argument. The percentage is expressed as normalised value between 0 and 1.0, where 1.0 means 100%. Official set it to 0.89999998. Undocumented by Wildcard. |
| `EnableMeshBitingProtection` | boolean | True | If `False`, disables mesh biting protection. Undocumented by Wildcard. |
| `FreezeReaperPregnancy` | boolean | False | If `True`, freezes the {{ItemLink|Reaper King}} pregnancy timer and [[experience]] gain. Undocumented by Wildcard. |
| `LogChatMessages` | boolean | False | If `True`, enables advanced chat logging. Chat logs will be saved in `ShooterGame/Saved/Logs/ChatLogs/<SessionName>/` in json format. Disabled on official. The file will be split according to `ChatLogFileSplitIntervalSeconds` value and flushed every `ChatLogFlushIntervalSeconds` seconds value. '''Note:''' `<SessionName>` will be written without spaces. Undocumented by Wildcard. |
| `MaxStructuresInSmallRadius` | integer | 0 | Defines the amount of max structures allowed to be placed in a `RadiusStructuresInSmallRadius` from player position. Official set it to 40. Undocumented by Wildcard. |
| `MaxStructuresToProcess` | integer | 0 | Controls the max batch size of structures to process (e.g.: culling, building graphs modifications, etc) at each server tick. Leaving at 0 (default behaviour) will force the server to process all structures in queue. It's a trade-off between how much can cost a server tick (worst case scenario) and simulation accuracy. Official set it to 5000. Undocumented by Wildcard. |
| `PreventOutOfTribePinCodeUse` | boolean | False | If `True`, prevents out of tribe players to use pins on structures (doors, elevators, storage boxes, etc). Undocumented by Wildcard. |
| `RadiusStructuresInSmallRadius` | float | 0.0 | Defines the small radius dimension (in Unreal Units) used by `MaxStructuresInSmallRadius`. Official set it to 225.0. Undocumented by Wildcard. |
| `ServerEnableMeshChecking` | boolean | False | Involved in foliage repopulation. Takes no effect if `-forcedisablemeshchecking` is set. Enabled on official. Undocumented. |
| `TribeMergeAllowed` | boolean | True | If `False`, prevents tribe to merge. Undocumented by Wildcard. |
| `TribeMergeCooldown` | float | 0.0 | Tribe merge cool-down in seconds. Official uses 86400.0. Undocumented by Wildcard. |
| `UpdateAllowedCheatersInterval` | float | 600.0 | Times in seconds at which the remote admin list linked by `AllowedCheatersURL` is queried for updates. Any value less than 3.0 will be reverted to 3.0. Undocumented by Wildcard. |
| `UseExclusiveList` | boolean | False | If `True`, allows same behaviour as `-exclusivejoin`. Undocumented by Wildcard. |
| `ListenServerTetherDistanceMultiplier` | float | 1.0 | Scales the tether distance between host and other players on non-dedicated sessions only. '''Note:''' despite being readable from command line, this option affects non-dedicated sessions only, thus it has to be set in the ''GameUserSettings.ini'' file or through the in-game local host graphics menu. |

## [SessionSettings]

| Parameter | Type | Default | Description |
|---|---|---|---|
| `MultiHome` | IP_ADDRESS | N/A | Specifies MultiHome IP Address. Boolean `Multihome` option must be set to `True` as well (command line or `[MultiHome]` section). Leave it empty if not using multihoming. Can be specified in command line too. |
| `Port` | integer | 7777 | Specifies the ''UDP Game Port''. See [[Dedicated_server_setup#Network|Dedicated server setup]]<br/>'''Note''': command line append syntax is not supported by {{gamelink|sa}} |
| `QueryPort` | integer | 27015 | Specifies the ''UDP Steam Query Port''. See [[Dedicated_server_setup#Network|Dedicated server setup]] |
| `SessionName` | string | ARK #123456 | Specifies the Server name advertised in the Game Server Browser as well in Steam Server browser. If no name is provide, the default name will be ''ARK #'' followed by a random 6 digit number. '''Note:''' Name must '''not''' be typed between quotes unless it is launched from command line. |

## [MultiHome]

| Option | Default | Effect |
|---|---|---|
| `MultiHome='''<boolean>'''` | False | If `True`, enables multihoming. `MultiHome` IP must be specified in `[SessionSettings]` or in command line as well. |

## [/Script/Engine.GameSession]

| Parameter | Type | Default | Description |
|---|---|---|---|
| `MaxPlayers` | integer | 70 | Specifies the maximum number of players that can play on the server simultaneously.<br/>{{ItemLink|ARK: Survival Ascended|ASA}}: This setting is replaced with `-WinLiveMaxPlayers` in the [[#options|command line options]], as otherwise, it will force it back to the default value. |

## [Ragnarok]

| Parameter | Type | Default | Description |
|---|---|---|---|
| `AllowMultipleTamedUnicorns` | boolean | False | `False` = one unicorn on the map at a time, `True` = one wild and unlimited tamed Unicorns on the map. |
| `EnableVolcano` | boolean | True | `False` = disabled (the volcano will '''not''' become active), `True` = enabled |
| `UnicornSpawnInterval` | integer | 24 | How long in hours the game should wait before spawning a new [[Unicorn]] if the wild one is killed (or tamed, if `AllowMultipleTamedUnicorns` is enabled). This value sets the minimum amount of time (in hours), and the maximum is equal to 2x this value. |
| `VolcanoIntensity` | float | 1 | The lower the value, the more intense the volcano's eruption will be. Recommended to leave at 1. The minimum value is 0.25, and for multiplayer games, it should not go below 0.5. Very high numbers will basically disable the flaming rocks flung out of the volcano. |
| `VolcanoInterval` | integer | 0 | 0 = 5000 (min) - 15000 (max) seconds between instances of the volcano becoming active. Any number above 0 acts as a multiplier, with a minimum value of .1 |

## [MessageOfTheDay]

| Parameter | Type | Default | Description |
|---|---|---|---|
| `Duration` | integer | 20 | Specifies in seconds the duration of the displayed message on player log-in. |
| `Message` | string | N/A | A single line string for a message displayed to played once logged-in. No quotes needed. Use `\n` to start a new line in the message. |