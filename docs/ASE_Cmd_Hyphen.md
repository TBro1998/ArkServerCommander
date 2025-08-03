# ASE Command-Line Options: Hyphen (-)

This file contains the list of command-line parameters for ARK: Survival Evolved that start with a hyphen (`-`), extracted from the wiki. Parameters that already exist in `Game.ini` or `GameUserSettings.ini` have been excluded.

| Argument | Default | Description |
|---|---|---|
| `-ActiveEvent=<eventname>` | | Enables a specified event or enables the event's colours palette on wild creatures. |
| `-automanagedmods` | | Steam only, automatic MOD download/installation/updating. Mod IDs are listed in `Game.ini` under `[ModInstaller]` section. |
| `-crossplay` | | Enables crossplay (i.e.: EPIC and Steam) on dedicated server. '''Note''': `PublicIPForEpic` must also be set. |
| `-culture=<lang_code>` | | Overrides the server output language. |
| `-DisableCustomFoldersInTributeInventories` | | Disables creation of folders in Tribute inventories. |
| `-DisableRailgunPVP` | | Disables the [[Tek Railgun (Aberration)|Tek Railgun]] on [[PvP]] servers. |
| `-EnableIdlePlayerKick` | | Cause characters that have not moved or interacted within the `KickIdlePlayersPeriod` to be kicked. |
| `-epiconly` | | Enables Epic Game Store only players to connect to the dedicated server. |
| `-exclusivejoin` | | Activate a whitelist only mode on the server. |
| `-ForceAllowCaveFlyers` | | Force flyer creatures to be allowed into caves. |
| `-ForceRespawnDinos` | | Launch with this command to destroy all wild creatures on server start-up. |
| `-imprintlimit=101` | | Automatically destroys creatures exceeding the imprinting bonus limit specified, in %. |
| `-insecure` | | Disable Valve Anti-Cheat (VAC) system. |
| `-MapModID=<ModID>` | | Dedicated servers can now optionally load custom maps via `ModID` directly. |
| `-MaxNumOfSaveBackups=<integer>` | 20 | Set the max number of backup to keeps in the save folder. |
| `-MinimumTimeBetweenInventoryRetrieval=<seconds>` | 3600 | {{ItemLink|Fjordhawk}} cool-down in seconds on retrieval of player's inventory when killed. |
| `-MULTIHOME` | | Enables multihoming. `MULTIHOME=<IP_ADDRESS>` must be specified as well. |
| `-newsaveformat` | | Enables the new save format (version 11). |
| `-NewYearEvent` | | Enables the [[ARK: Happy New Year!]] event. |
| `-noantispeedhack` | | Disables anti speedhack detection. |
| `-NoBattlEye` | | Run server without BattleEye. |
| `-nocombineclientmoves` | | Disables server player move physics optimization. |
| `-NoDinos` | | Prevents wild creatures from being spawned. |
| `-NoHangDetection` | | Prevents the server from shutdown if its start-time takes more than 45 minutes. |
| `-NotifyAdminCommandsInChat` | | Displays `cheat` and `admincheat` commands to admin players only in the global chat. |
| `-noundermeshchecking` | | Will turn off the anti-meshing system entirely. |
| `-noundermeshkilling` | | Will turn off the anti-meshing player kills (but still allow teleporting). |
| `-PublicIPForEpic=<IPAddress>` | | This is the public IP that EGS clients will attempt to connect to. |
| `-pvedisallowtribewar` | | Prevents tribes in [[PvE]] mode to officially declare war on each other. |
| `-SecureSendArKPayload` | | Add this option to protect against the creature upload issue. |
| `-ServerAllowAnsel` | | Sets the server to allow clients to use [[NVIDIA Ansel]]. |
| `-servergamelog` | | Enable server admin logs (including support in the RCON). |
| `-servergamelogincludetribelogs` | | Include tribe logs in the normal server game logs. |
| `-ServerRCONOutputTribeLogs` | | Include tribe logs in the RCON output. |
| `-speedhackbias=<value>` | 1.0 | Bias you can use on your server for speedhack detection. |
| `-StasisKeepControllers` | | If enabled, the server will keep AI controller objects in memory when actors go in stasis. |
| `-structurememopts` | | Enables structure memory optimizations. |
| `-UseItemDupeCheck` | | Enables additional dupe protection. |
| `-UseSecureSpawnRules` | | Enables an additional protection for items and creatures spawning. |
| `-usestore` | | Handles player characters data like official non-legacy servers. |
| `-UseStructureStasisGrid` | | Enables the structure stasis grid to improve server performance. |
| `-UseVivox` | | Enables Vivox on Steam only servers. |
| `-webalarm` | | Activate Web alarms for important tribe events. |
| `-ClusterDirOverride=<PATH>` | | Specifies cross-server storage location for Cross-ARK Data Transfer. |
| `-clusterid=<CLUSTER_NAME>` | | Specifies cluster name ID to allow Cross-ARK Data Transfer. |
| `-NoTransferFromFiltering` | | Prevents ARK Data usage between single player and servers who do not have a cluster ID. |
| `-AllowChatSpam` | | Allows players to spam messages in server chat. Undocumented. |
| `-BackupTransferPlayerDatas` | | Adds permanent character profile backup files separately from the world-save. Requires `-usestore`. Undocumented. |
| `-BattlEyeServerRecheck` | | Adds at every server tick an extra check about BattleEye service status. Undocumented. |
| `-converttostore` | | Converts non-stored data to stored data at next world-save. Requires `-usestore`. Undocumented. |
| `-CustomAdminCommandTrackingURL=<URL>` | | Allows server to specify custom URL to post admin commands tracking. Undocumented. |
| `-CustomMerticsURL=<URL>` | | Allows server to specify custom URL to post server metrics statistics. Undocumented. |
| `-CustomNotificationURL=<URL>` | | Allows server custom notification broadcast using the server message feature. Undocumented. |
| `-dedihibernation` | | Forces wild creatures to spawn in hibernation mode. Undocumented. |
| `-DisableDupeLogDeletes` | | Prevents `-ForceDupeLog` to take effect. Undocumented. |
| `-DormancyNetMultiplier=<float>` | | Multiplier for dormancy rate. Undocumented. |
| `-EnableOfficialOnlyVersioningCode` | | Mimics official versioning code-path on ARK data items upload/download behaviour. Undocumented. |
| `-EnableVictoryCoreDupeCheck` | | Should enabled an extra anti-duping check on tamed creatures. Undocumented. |
| `-forcedisablemeshchecking` | | Disables mesh checking on foliage re-population. Undocumented. |
| `-ForceDupeLog` | | Forces dupe logs. Undocumented. |
| `-ignoredupeditems` | | If a duped item is detected in an inventory, it will be ignored and not removed. Undocumented. |
| `-MaxConnectionsPerIP=<integer>` | | Sets a limit of users connected to the server with the same IP at the same time. Undocumented. |
| `-nitradotest2` | | Limits the tick-rate to 3 when there are no players connected. Undocumented. |
| `-nodormancythrottling` | | Required to make `DormancyNetMultiplier` to work if its value is less than 1.0. Undocumented. |
| `-parseservertojson` | | Dumps the loading server save to a JSON file at start-up. Undocumented. |
| `-pauseonddos` | | Should pause/halt the game if a DDoS attack by a performed by players is detected. Undocumented. |
| `-PreventTotalConversionSaveDir` | | Prevents total conversion mods to use a dedicated save directory. Undocumented. |
| `-pveallowtribewar` | | Allows tribes to officially declare war on each other. Undocumented. |
| `-ReloadedForBackup` | | Marks the world-save to load as a trial backup save. Undocumented. |
| `-UnstasisDinoObstructionCheck` | | Should prevent creatures from ghosting through meshes/structures on re-render. Undocumented. |
| `-UseTameEffectivenessClamp` | | Prevents taming effectiveness never going over 100%. Undocumented. |
| `-UseServerNetSpeedCheck` | | Avoids players accumulating too much movement data per server tick. Undocumented. |