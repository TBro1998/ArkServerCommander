{{DisambigMsg|INI configuration files of {{gamelink|sa}} and {{gamelink|se}} and dedicated server command line options|the installation and fundamental setup of a dedicated server|Dedicated server setup|scripts that make it easier to run a dedicated server|Dedicated server scripts}}

This page discusses the extensive collection of behaviour and gameplay aspect-altering configuration settings of {{ItemLink|ARK: Survival Ascended}} and {{ItemLink|ARK: Survival Evolved}} servers and Single Player/Non-Dedicated sessions.

Options can be specified on the [[#Command_Line|command line]] when launching the server, or in the [[#Configuration_Files|configuration files]] loaded at start-up. Some of those options can be changed at runtime via the [[#DynamicConfig|dynamic configuration system]] as well.

== Command line ==
=== Syntax ===
{{ambox|color=blue|type=This section is only relevant to dedicated servers.}}
Dedicated servers are launched via command line using the following syntax to specify runtime options '''as a single line string''':

 '''<executable>''' '''<map_name>'''[?'''<option>'''='''<value>''']...[?'''<option>'''='''<value>'''] [-'''<option>'''[='''<value>''']] ... [-'''<option>'''[='''<value>''']]''

The <code><executable></code> (<code>ArkAscendedServer.exe</code> for ASA, <code>ShooterGameServer.exe</code> for ASE) is followed right after by the '''mandatory''' <code><map_name></code>. Afterwards, options are specified with each of them separated either first by a question mark (<code>?</code>) or by a space and a hyphen-minus (<code> -</code>) depending on their syntax. Some options require a value argument too. When an option is not specified, its default value is used automatically (refer to the option reference tables below for more details).

==== Example ====
To launch a dedicated ARK: Survival Evolved server that:
# shows each player a crosshair
# hides each player their position on the map
# enables players to view themselves in third person
# enforces a maximum amount of 1000 structures built by different players and tribes in a minimum distance
# disables BattlEye anti-cheat
On a Linux host:
 ./ShooterGameServer TheIsland?ServerCrosshair=true?ShowMapPlayerLocation=false?AllowThirdPersonPlayer=true?TheMaxStructuresInRange=1000 -NoBattlEye
On a Windows host:
 .\ShooterGameServer.exe TheIsland?ServerCrosshair=true?ShowMapPlayerLocation=false=AllowThirdPersonPlayer=true?TheMaxStructuresInRange=1000 -NoBattlEye
:Note: in both cases ShooterGameServer path had to be specified depending on launching environment.

=== Maps ===
{{ambox|color=blue|type=This section is only relevant to dedicated servers.}}
This table shows level names of all official maps ([[#Syntax|<code><map_name></code>]]):

{| class="wikitable"
! Map !! {{ItemLink|ARK: Survival Ascended}} !! {{ItemLink|ARK: Survival Evolved}}
|-
|[[The Island]] || <code>TheIsland_WP</code> || <code>TheIsland</code>
|-
|[[The Center]] || <code>TheCenter_WP</code> || <code>TheCenter</code>
|-
|[[Scorched Earth]] || <code>ScorchedEarth_WP</code> || <code>ScorchedEarth_P</code>
|-
|[[Ragnarok]] || <code>Ragnarok_WP</code>|| <code>Ragnarok</code>
|-
|[[Aberration]] || <code>Aberration_WP</code>|| <code>Aberration_P</code>
|-
|[[Extinction]] || <code>Extinction_WP</code>|| <code>Extinction</code>
|-
|[[Valguero]] || ''Unavailable'' || <code>Valguero_P</code>
|-
|[[Genesis: Part 1]] || ''Unavailable'' || <code>Genesis</code>
|-
|[[Crystal Isles]] || ''Unavailable'' || <code>CrystalIsles</code>
|-
|[[Genesis: Part 2]] || ''Unavailable'' || <code>Gen2</code>
|-
|[[Lost Island]] || ''Unavailable'' || <code>LostIsland</code>
|-
|[[Fjordur]] || ''Unavailable'' || <code>Fjordur</code>
|-
|[[Procedurally Generated Maps]] || ''Does not apply'' || <code>PGARK</code>
|-
|[[Astraeos]] || <code>Astraeos_WP</code> || ''Unavailable''
|-
|[[Club ARK]] || <code>BobsMissions_WP</code><br/>
Requires [https://www.curseforge.com/ark-survival-ascended/mods/club-ark Club ARK]<br/>mod <code>1005639</code> 
|| ''Does not apply''
|}

For [[Mod:Mods|modded]] maps check their [https://www.curseforge.com/ark-survival-ascended CurseForge] (for ASA) or [https://steamcommunity.com/app/346110/workshop/ Workshop] (for ASE) pages, or contact the author.

For {{ItemLink|ARK: Survival Evolved}}'s [[Procedurally Generated Maps|PGARK]]s refer to [https://survivetheark.com/index.php?/forums/topic/99062-quickie-procedurally-generated-arks-how-to-guide/|Quickie Procedurally Generated ARKs How-To Guide].

== Command line options ==
These options can be specified '''only on the command line'''. Options starting with a <code>?</code> (question mark) are appended directly one after another, options starting with a <code>-</code> (hyphen-minus) have to be separated by spaces.

Options having a red background are deprecated and may no longer work, most of time their presence when using an external tool (e.g.: ASM) is harmless unless different stated. Undocumented options are stated too: please, beware their behaviour has been inferred by reverse engineering of ''ShooterGameServer'' executable bits. Such options and are not officially supported by Wildcard and may result in undefined behaviour.

For '''Single Player''' and '''non-dedicated''' sessions, only options starting with a <code>-</code> (hyphen-minus) will work. All other options have to be configured with INI files.
* For '''Steam''' players: right click on game name in ''LIBRARY'' and choose ''Properties...'', in ''GENERAL'' tab under ''LAUNCH OPTIONS'', type the desired options.
* For '''Epic''' players: left click on account icon (top right of window) and choose ''Settings'', scroll down until game name (do not choose DLC names) and click on it, mark the ''Additional Command Line Arguments'', type the desired options.

Finally, although not suggested for easy of usage and command line readability, additional options may be set using the (<code>?</code>) syntax from [[#GameUserSettings.ini|GameUserSettings.ini]] (see [[#Configuration Files|Configuration Files]] for its location).

<!--start of command line-->
<div class="widetable">
{| class="wikitable config-table"
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Argument
! Description
! Since patch
{{Server config variable
| name = -ActiveEvent=<eventname>
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Enables a specified event or enables the event's colours palette on wild creatures. Only one can be specified and active at a time. Only the <code>WinterWonderland</code> event is still fully working in the game. Events marked red no longer even apply the colour modifications.
<table class="wikitable">
  <tr><th>eventname</th><th>Description</th></tr>
  <tr class="redish-background"><td><code>ark7th</code></td><td>Allows for [[ARK: 7th Anniversary]] to be activated.</td></tr>
  <tr class="redish-background"><td><code>ARKaeology</code></td><td>Allows for [[ARK: ARKaeology]] to be activated.</td></tr>
  <tr class="redish-background"><td><code>birthday</code></td><td>Allows for [[ARK: 5th Anniversary]] to be activated.</td></tr>
  <tr><td><code>Easter</code></td><td>Allows for [[ARK: Eggcellent Adventure 7]] to be activated.</td></tr>
  <tr><td><code>FearEvolved</code></td><td>Allows for [[ARK: Fear Evolved 6]] to be activated.</td></tr>
  <tr class="redish-background"><td><code>ExtinctionChronicles</code></td><td>Allows for [[Extinction Chronicles]] to be activated.</td></tr>
  <tr><td><code>PAX</code></td><td>Allows for [[ARK: PAX Party]] to be activated.</td></tr>
  <tr><td><code>Summer</code></td><td>Allows for [[ARK: Summer EVO]] to be activated.</td></tr>
  <tr><td><code>TurkeyTrial</code></td><td>Allows for [[ARK: Turkey Trial 6]] to be activated.</td></tr>
  <tr><td><code>vday</code></td><td>Allows for [[Valentine's EVO Event]] to be activated.</td></tr>
  <tr><td><code>WinterWonderland</code></td><td>Allows for [[ARK: Winter Wonderland 7]] to be activated.</td></tr>
  <tr><td><code>None</code></td><td>Disables normally active event. It is suggested to use this setting to completely turn off unwanted event references in server memory, like skin crafting.</td></tr>
</table>
''The last event to feature this option was [[ARK: Winter Wonderland 7]].''<br/>
'''Note:''' If you don't want to activate the event, but still want the wild creatures to have the event colors, use the option <code>ActiveEventColors</code> from [[#DynamicConfig|Dynamic Config]] instead.<br/>
{{gamelink|asa}}: Obsolete, use <code>-mods</code> argument instead. Adding <code>ActiveEvent=</code> in the [[#GameUserSettings.ini|GameUserSettings.ini]] file has no effect.
| version = [[280.114]]
}}
{{Server config variable
| name = -AllowFlyerSpeedLeveling
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Enables Level Up for Flyer Movement Speed. '''Note:''' alternatively, it can be enabled adding <code>bAllowFlyerSpeedLeveling=True</code> to the ''Game.ini'' file.
| version = [[321.1]]
}}
{{Server config variable
| name = ?AltSaveDirectoryName=<savedir_name>
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Allows to specify a custom directory name for server world-save. Usually used to manage clusters, read more about it in the [[#Cross-ARK_Data_Transfer|Cross-ARK Data Transfer]] section.
{{gamelink|asa}}: Note that it still uses the <code><savedir_name>/<map_name></code> directory structure.
}}
{{Server config variable
| name = -AlwaysTickDedicatedSkeletalMeshes
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Disables optimization to dynamically not tick idle/inactive creature animations on the dedicated server based on server frame rate (only becomes active when server is low FPS). '''Note''': side-effects can include inaccurate collisions for idle creatures on the client versus the server (such as inaccurate jumping position while standing on an inactive/idle Brontosaurus and jumping on its tail as it sits in your base).
| version = [[ARK: Survival Ascended/Patch/33.33|ASA 26.36]]
}}
{{Server config variable
| name = -AutoDestroyStructures
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Enables auto destruction of old structures. Timer can be adjusted with option <code>AutoDestroyOldStructuresMultiplier</code> in [[#GameUserSettings.ini|GameUserSettings.ini]] under <code>[ServerSettings]</code> section.
| version = [[245.987]]
}}
{{Server config variable
| name = -automanagedmods
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Steam only, automatic MOD download/installation/updating. Mod IDs are listed in [[#Game.ini|Game.ini]] under <code>[ModInstaller]</code> section.
| version = [[244.3]]
}}
{{Server config variable
| name = -crossplay
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Enables crossplay (i.e.: EPIC and Steam) on dedicated server. '''Note''': <code>PublicIPForEpic</code> must also be set.<br/>
{{gamelink|asa}}: Obsolete, use <code>-ServerPlatform</code> argument instead.
| version = [[311.74]]
}}
{{Server config variable
| name = -culture=<lang_code>
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Overrides the server output language. If set in server it doesn't override clients' language. List of currently supported language codes: <code>ca</code>, <code>cs</code>, <code>da</code>, <code>de</code>, <code>en</code>, <code>es</code>, <code>eu</code>, <code>fi</code>, <code>fr</code>, <code>hu</code>, <code>it</code>, <code>ja</code>, <code>ka</code>, <code>ko</code>, <code>nl</code>, <code>pl</code>, <code>pt_BR</code>, <code>ru</code>, <code>sv</code>, <code>th</code>, <code>tr</code>, <code>zh</code>, <code>zh-Hans-CN</code>, <code>zh-TW</code>. This is also needed to override some 
Genesis: Part 2 missions text outputs for clients that otherwise will follow the language settings of the server operating system.
| version = [[235.2]]
}}
{{Server config variable
| name = -DisableCustomCosmetics
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Disables the Custom Cosmetic system (allowing players to use and display special mods that should only be skins and will be downloaded automatically by the connected clients).
| version = [[ARK: Survival Ascended/Patch/34.49|ASA 34.49]]
}}
{{Server config variable
| name = -disabledinonetrangescaling
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Disables creature's network replication range optimization.
Detailed explanation : On dedicated servers the network replication range is scaled dynamically based on the amount of connected players (Note: On PvP servers this replication range is scaled only for small tamed creatures and wild creatures, rather than all tames). The scaling of replication range allows to improve performances for servers with a high amount of connected players. This optimization can be disabled by adding <code>-disabledinonetrangescaling</code> in server startup options.
| version = [[ARK: Survival Ascended/Patch/31.41|ASA 31.41]]
}}
{{Server config variable
| name = -DisableCustomFoldersInTributeInventories
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Disables creation of folders in Tribute inventories.
| version = [[349.1]]
}}
{{Server config variable
| name = -DisableRailgunPVP
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Disables the [[Tek Railgun (Aberration)|Tek Railgun]] on [[PvP]] servers (on official servers this argument is not added to the command line).
| version = [[349.1]]
}}
{{Server config variable
| name = -EasterColors
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Chance for Easter colors when creatures spawn.
| version = [[ARK: Survival Ascended/Patch/36.27|ASA 36.27]]
}}
{{Server config variable
| name = -EnableIdlePlayerKick
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Cause characters that have not moved or interacted within the <code>KickIdlePlayersPeriod</code> ([[#GameUserSettings.ini|GameUserSettings.ini]]) to be kicked.
| version = [[241.5]]
}}
{{Server config variable
| name = -epiconly
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Enables Epic Game Store only players to connect to the dedicated server.<br/>
{{Gamelink|sa}}: Obsolete since ASA is not available on Epic.
| version = [[311.74]]
}}
{{Server config variable
| name = ?EventColorsChanceOverride=<float>
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Takes a floating point value between 0 and 1.0 (1 = 100% chance for colour override). '''Note:''' This only affects freshly spawned wild creatures (you can use console command <code>Cheat DestroyWildDinos</code> to wipe wild creatures so that they all respawn with a chance of getting the event colors). 
| version = [[355.3]]
}}
{{Server config variable
| name = -exclusivejoin
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Activate a whitelist only mode on the server, allowing players to join if added to the allow list. For more info on whitelists, start with [[Server configuration#Ark_IDs|Ark IDs]].
}}
{{Server config variable
| name = -ForceAllowCaveFlyers
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Force flyer creatures to be allowed into caves (Flyers able to go into caves by default on custom maps).
| version = [[237.1]]
}}
{{Server config variable
| name = -ForceRespawnDinos
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Launch with this command to destroy all wild creatures on server start-up. '''Note:''' this will only destroy wild creatures that are not currently tamed.
| version = [[216.0]]
}}
{{Server config variable
| name = ?GameModIds=<ModID1>[,<ModID2>[...]]
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Steam only. Specifies the order and which mods are loaded, <code>ModIDs</code> need to be separated with commas (<code>,</code>). Mod priority is in descending order left to right (the left-most ID is the top priority mod). It is suggested to use instead the <code>ActiveMods</code> under <code>[ServerSettings]</code> of ''GameUserSettings.ini''.<br/>
{{Gamelink|sa}}: Obsolete, use argument <code>-mods</code> instead.
| version = [[190.0]]
}}
{{Server config variable
| name = -GBUsageToForceRestart=<value>
| inASA = Yes
| inASE = No
| default = 35
| cli = skip
| info = Enabled an out of memory protection for servers. In the event that a server reaches an excessive amount of memory, the world will be forced saved and the server will be restarted to prevent crashes or rollbacks. Replace <code><value></code> with the wanted amount in gigabytes. Setting it to 0 disables.
| version = [[ARK: Survival Ascended/Patch/26.32|ASA 26.32]]
}}
{{Server config variable
| name = -imprintlimit=101
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Automatically destroys creatures exceeding the imprinting bonus limit specified, in %. It is recommended going 1% higher (101) to catch potential floating-point errors, especially when working with custom imprinting rates.
| version = [[349.10]]
}}
{{Server config variable
| name = -insecure
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Disable Valve Anti-Cheat (VAC) system.
| version = [[241.5]]
}}
{{Server config variable
| name = -MapModID=<ModID>
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Dedicated servers can now optionally load custom maps via <code>ModID</code> directly, rather than having to specify the map name, using this syntax (where the <code>MapModID</code> is the Steam Workshop ID of your custom map, and the GameModIds are the Id’s of the stacked mods you wish to use, in order). <code>ActiveMods</code> must also be set in [[#GameUserSettings.ini|GameUserSettings.ini]].<br/>
{{gamelink|sa}}: Obsolete, use <code>-mods</code> argument instead. ASA no longer uses Steam Workshop mods and <code>ActiveMods</code> does not need to be set.
| version = [[193.0]]
}}
{{Server config variable
| name = -MaxNumOfSaveBackups=<integer>
| inASA = Unknown
| inASE = Yes
| default = 20
| cli = skip
| info = Set the max number of backup to keeps in the save folder. Whenever a new backup is created, the oldest ones will be deleted. Since such backups take place every 2 hours and they are uncompressed files, it is suggested to use more flexible third party solutions, especially with bigger saves.<br/>'''Unofficial admins running [[2023_official_server_save_files|Official 2023 Server saves]] may want to set this to 0 when using other backup mechanisms.'''
}}
{{Server config variable
| name = ?MapPlayerLocation=false
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Same as <code>ShowMapPlayerLocation</code> of [[#GameUserSettings.ini|GameUserSettings.ini]] (actually using this option will make the server write the <code>ShowMapPlayerLocation</code> entry).
| version = [[170.0]]
}}
{{Server config variable
| name = -MinimumTimeBetweenInventoryRetrieval=<seconds>
| inASA = Unknown
| inASE = Yes
| default = 3600
| cli = skip
| info = {{ItemLink|Fjordhawk}} cool-down in seconds on retrieval of player's inventory when killed.
| version = [[351.5]]
}}
{{Server config variable
| name = -mods=<ModId1>[,<ModId2>[...]]
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Specifies [https://www.curseforge.com/ark-survival-ascended CurseForge] Mod Project IDs. Mods are updated automatically when starting the server.
| version = {{ItemLink|ARK: Survival Ascended|&nbsp;}}
}}
{{Server config variable
| name = -MULTIHOME
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Enables multihoming. <code>MULTIHOME=<IP_ADDRESS></code> must be specified as well in [[#GameUserSettings.ini|GameUserSettings.ini]] under <code>[SessionSettings]</code> or in the command line.
}}
{{Server config variable
| name = -newsaveformat
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Enables the new save format (version 11), which is quicker at saving cryopods and requires less storage space for them. Use only if necessary (e.g.: the default save format fails as well the old save format).<br/>'''This option must be used to start the [[2023_official_server_save_files|Official 2023 Server saves]] backups''' along with <code>-usestore</code> to correctly load tribe ownership.
| version = [[244.5]]
}}
{{Server config variable
| name = ?NewYear1UTC=<epoch time>
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = You can use https://www.epochconverter.com to get your appropriate timestamp. Please keep in mind that this event cannot be run before Midnight EST on the 1st of January, so if you wish to change your time it would have to be set after that. Must be used with the with the <code>-NewYearEvent</code> command line option.
| version = [[320.38]]
}}
{{Server config variable
| name = ?NewYear2UTC=<epoch time>
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = You can use https://www.epochconverter.com to get your appropriate timestamp. Please keep in mind that this event cannot be run before Midnight EST on the 1st of January, so if you wish to change your time it would have to be set after that. Must be used with the with the <code>-NewYearEvent</code> command line option.
| version = [[320.38]]
}}
{{Server config variable
| name = -NewYearEvent
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Enables the [[ARK: Happy New Year!]] event; it will automatically start at Midnight EST and Noon EST on the 1st of January, unless specified with the following 2 options. This command was only available around the time it was scheduled.
| version = [[320.38]]
}}
{{Server config variable
| name = -noantispeedhack
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Anti speedhack detection is now enabled by default — you now use this option to disable it.
| version = [[218.5]]
}}
{{Server config variable
| name = -NoBattlEye
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Run server without BattleEye.
| version = [[253.0]]
}}
{{Server config variable
| name = -nocombineclientmoves
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Server player move physics optimization is now enabled by default (improves performance) — to disable it, use this server option.
| version = [[218.5]]
}}
{{Server config variable
| name = -NoDinos
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Prevents wild creatures from being spawned. In {{gamelink|sa}} enabling this option will make the following options unusable: <code>-NoDinosExceptForcedSpawn</code>, <code>-NoDinosExceptStreamingSpawn</code>, <code>-NoDinosExceptManualSpawn</code> and <code>-NoDinosExceptWaterSpawn</code>.<br/>
'''Note:''' wild creatures that were on the server before adding this argument will not get automatically destroyed (you can destroy them with console command <code>cheat destroywilddinos</code>).
}}
{{Server config variable
| name = -NoHangDetection
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Prevents the server from shutdown if its start-time takes more than 2700 seconds (45 minutes). Useful for servers with big saves or very slow machines. '''Note:''' this will disabled further server hangs detection.
| version = [[246.7]]
}}
{{Server config variable
| name = -NotifyAdminCommandsInChat
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Displays <code>cheat</code> and <code>admincheat</code> commands to admin players only in the global chat. Not to be confused with <code>AdminLogging</code> which displays such commands to all players.
| version = [[254.0]]
}}
{{Server config variable
| name = -noundermeshchecking
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Will turn off the anti-meshing system entirely.
| version = [[304.445]]
}}
{{Server config variable
| name = -noundermeshkilling
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Will turn off the anti-meshing player kills (but still allow teleporting).
| version = [[304.445]]
}}
{{Server config variable
| name = -NoWildBabies
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Disables spawning of [[Wild babies|wild babies]].
| version = {{ItemLink|ARK: Survival Ascended|&nbsp;}}
}}
{{Server config variable
| name = -passivemods=<ModId1>[,<ModId2>[...]]
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = This option disables a mod's functionality while still loading its data. This is useful if, for example, you want to have S-Dinos on Astraeos in your cluster but prevent them from spawning on The Island while still allowing them to be transferred back and forth.  

*Note: The -passivemods function works the same as the regular mod list and will download the specified mod. It does not require the passive mod to also be in the regular mods list. [https://ark.wiki.gg/wiki/Talk:Server_configuration#Passive_Mods_Do_not_need_-mods=]
| version = [[ARK: Survival Ascended/Patch/33.33|ASA 33.33]]
}}
{{Server config variable
| name = -port=<Server Port>
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = This is the server port used to host the dedicated server on. This is required if you want to run it on a port other than the default 7777.
| version = {{ItemLink|ARK: Survival Ascended|&nbsp;}}
}}
{{Server config variable
| name = -PublicIPForEpic=<IPAddress>
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = This is the public IP that EGS clients will attempt to connect to, if this option is missing and <code>-MULTIHOME</code> is specified, then EGS clients will attempt to connect to the multihome IP; note that if you're using multihome and specify a non-public IP address, then players will not be able to connect to your server using EGS. Make sure to set a public IP address (e.g.: WAN or external).
| version = [[312.13]]
}}
{{Server config variable
| name = -pvedisallowtribewar
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Prevents tribes in [[PvE]] mode to officially declare war on each other for mutually-agreed-upon period of time. Alternatively, it is possible to use the ''Game.ini'' <code>bPvEAllowTribeWar</code>.
| version = [[298.31]]
}}
{{Server config variable
| name = -SecureSendArKPayload
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Add this option to protect against the creature upload issue. '''Note:''' be aware this may have an impact on some mods.
| version = [[357.20]]
}}
{{Server config variable
| name = -ServerAllowAnsel
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Sets the server to allow clients to use [[NVIDIA Ansel]]. See link for details on its support.
| version = [[278.0]]
}}
{{Server config variable
| name = -servergamelog
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Enable server admin logs (including support in the RCON). Use RCON command <code>getgamelog</code> to print 100 entries at a time. Also, outputs to dated file in <code>\Logs</code>, adjust max length of the RCON buffer with the option <code>RCONServerGameLogBuffer=600</code>.
| version = [[224.0]]
}}
{{Server config variable
| name = -servergamelogincludetribelogs
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Include tribe logs in the normal server game logs. Note that this is required to get the logs from RCON (using <code>-ServerRCONOutputTribeLogs</code>). Additionally, it will display them in the application's console.
| version = [[254.0]]
}}
{{Server config variable
| name = -ServerRCONOutputTribeLogs
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Include tribe logs in the RCON output (which can be recovered with console command <code>GetGameLog</code>). Requires <code>-servergamelogincludetribelogs</code>.
| version = [[254.0]]
}}
{{Server config variable
| name = -speedhackbias=<value>
| inASA = Unknown
| inASE = Yes
| default = 1.0
| cli = skip
| info = Bias you can use on your server (can cause more rubber banding on laggy players).
| version = [[219.0]]
}}
{{Server config variable
| name = -StasisKeepControllers
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = If enabled, the server will keep Unreal Engine actors (e.g.: creatures) AI controller objects in memory when actors go in stasis. This favours performance when a portion of a map is rendered by a player and actors go out of stasis. However, especially on large maps, it will result in higher memory usage even with few or no players.
| version = [[252.96]]
}}
{{Server config variable
| name = -structurememopts
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Enables structure memory optimizations. '''Note:''' avoid using it when running structure-related mods (until they get updated) as it can break the snap points of these mod structures.
| version = [[295.108]]
}}
{{Server config variable
| name = -TotalConversionMod=<ModID>
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Used to specify a total conversion mod (e.g.: [[Primitive Plus]]). Alternatively, <code>ActiveTotalConversion=<ModID></code> in [[#GameUserSettings.ini|GameUserSettings.ini]] can also be set.
}}
{{Server config variable
| name = -UseDynamicConfig
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Enables the use of the dynamic config, if you don't provide a <code>CustomDynamicConfigUrl</code> the server will use the default dynamic config for that platform (the one used on official server). '''Note:''' unless you force an immediate update using the [[Console Commands#ForceUpdateDynamicConfig|ForceUpdateDynamicConfig]] command, the configuration is checked every time the world is (auto)saved. When you want to "undo" the config you should be changing it back to whatever your default is, so it can be read from the updated config and thus applied. Omitting any lines will not update those options.
| version = [[307.2]]
}}
{{Server config variable
| name = -UseItemDupeCheck
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Enables additional dupe protection. '''Note:''' this could have an impact on mods, so use with caution.
| version = [[333.13]]
}}
{{Server config variable
| name = -UseSecureSpawnRules
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Enables an additional protection for items and creatures spawning.<br/>
'''Note:''' This protection is disabled by default since version [[349.27]], due to server crashes caused by certain mods. The mods in question should now work correctly, so you can now use this option without too much fear.
| version = [[349.27]]
}}
{{Server config variable
| name = -usestore
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Same basic behaviour of official (non-legacy) servers on handling player characters data: character profile file (.arkprofile) is not saved separately from map save (.ark), nor its backup. When necessary (e.g.: a character just transferred or a new character has been created) a temporary character profile file is created and kept until next world-save. Tribe profiles are stored in a different format (.arktributetribe). Allows <code>-BackupTransferPlayerDatas</code> option (for fully mimic official servers save behaviour). Use the option <code>-converttostore</code> to convert previous non stored data to stored data.  '''This option must be used to correctly load the [[2023_official_server_save_files|Official 2023 Server saves]] backups''' along with <code>-newsaveformat</code>.
}}
{{Server config variable
| name = -UseStructureStasisGrid
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Enables the structure stasis grid to improve server performance on large bases with lots of players<ref name="life_tuning_system">StudioWildcard (October 9, 2020) [https://survivetheark.com/index.php?/articles.html/community-crunch-241-fear-evolved-4-live-tuning-system-network-improvements-and-more-r1652 "Live Tuning System"]. ''Community Crunch 241: Fear Evolved 4, Live Tuning System, Network Improvements, and More!''</ref>.
| version = [[314.5]]
}}
{{Server config variable
| name = -UseVivox
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Enables Vivox on Steam only servers. Default for EPIC servers but can be enabled for Steam server.
| version = [[311.74]]
}}
{{Server config variable
| name = -webalarm
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Activate Web alarms when important things happen to a tribe, such as Tripwire Alarms going off and Babies being born. See [[Web Notifications]] for details.
| version = [[243.0]]
}}
{{Server config variable
| name = -WinLiveMaxPlayers=<integer>
| inASA = Yes
| inASE = No
| default = 70
| cli = skip
| info = Sets the maximum of players for the server. This currently replaces the <code>MaxPlayers</code> setting from the [[#GameUserSettings.ini|GameUserSettings.ini]] option.
| version = {{ItemLink|ARK: Survival Ascended|&nbsp;}}
}}
|-
! colspan="6" | CrossARK Transfers
|-
| colspan="6" | If your server does not use the <code>-clusterid</code> argument nor the <code>-NoTransferFromFiltering</code>, make sure to setup correctly the following options in the <code>[ServerSettings]</code> section of file [[#GameUserSettings.ini|GameUserSettings.ini]]: <code>noTributeDownloads</code>, <code>PreventDownloadDinos</code>, <code>PreventDownloadItems</code>, <code>PreventDownloadSurvivors</code>, <code>PreventUploadDinos</code>, <code>PreventUploadItems</code> and <code>PreventUploadSurvivors</code>. ''Note: It is also possible to setup these 7 options in the server launch command line (for that, use syntax <code>?<option>=<value></code>).''
{{Server config variable
| name = -ClusterDirOverride=<PATH>
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Specifies cross-server storage location for [[#Cross-ARK_Data_Transfer|Cross-ARK Data Transfer]]. Requires also <code>-clusterid=<CLUSTER_NAME></code> option to specify the cluster ID.
| version = [[247.85]]
}}
{{Server config variable
| name = -clusterid=<CLUSTER_NAME>
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Specifies cluster name ID to allow [[#Cross-ARK_Data_Transfer|Cross-ARK Data Transfer]].
| version = [[247.85]]
}}
{{Server config variable
| name = -NoTransferFromFiltering
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Prevents ARK Data usage between single player and servers who do not have a cluster ID. Even with a cluster id set, it is suggested to keep this option to completely disable code path allowing it.
}}
|-
! colspan="6" | Deprecated options
{{Server config variable
| status = deprecated
| name = ?bRawSockets
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Direct UDP socket connections rather than Steam P2P. ''Deprecated since [[311.78]] and no longer present in the executable code''.
| version = [[213.0]]
}}
{{Server config variable
| status = deprecated
| name = -ClearOldItems
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Official [[PvP]] servers one-time Clearance of all old unequipped items (with the exception of blueprints, eatables, notes, and quest items), to ensure fairness after Item Duplication bug exploit. Server admins can enforce this once if they run with this command argument (will only work ''ONCE''' on pre-update save games). ''Deprecated and no longer present in the executable code''.
| version = [[178.0]]
}}
{{Server config variable
| status = deprecated
| name = -forcenetthreading
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Defaulted dedicated server <code>?bRawSockets</code> mode to '''not''' use threaded networking, seemed to generally be a net performance loss. Use this to forcefully enable it. ''Deprecated since [[311.78]].''
| version = [[271.17]]
}}
{{Server config variable
| status = deprecated
| name = -inlinesaveload
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Fixed a save game corruption case with large saves: this is experimental, so try with this command, if you have any saves that don’t load. We’ll formally roll this loader change out in a subsequent patch after ''(Wildcard is)'' 100% certain it has no side effects. ''It should be no more needed, and is considered deprecated''.
| version = [[214.0]]
}}
{{Server config variable
| status = deprecated
| name = -NoBiomeWalls
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Eliminates the upcoming biome change area wall effects as introduced in an undocumented addition of patch [[241.5]]. ''Deprecated since patch [[243.0]] and no longer present in the executable code''.
| version = [[242.7]]
}}
{{Server config variable
| status = deprecated
| name = -nofishloot
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Disable non-meat fishing loot when using [[Fishing Rod]]. ''Deprecated since patch [[260.0]] and no longer present in the executable code''.
| version = [[245.9]]
}}
{{Server config variable
| status = deprecated
| name = -nonetthreading
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Option for <code>?bRawSockets</code> servers to only utilize a single thread for networking (useful to improve performance for machines with more servers than CPU cores, on Linux in particular). ''Deprecated with patch [[271.17]] and no longer present in the executable code''.
| version = [[271.15]]
}}
{{Server config variable
| status = deprecated
| name = -noninlinesaveload
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Prevented <code>-inlinesaveload</code> to happen. ''Deprecated and no longer present in the executable code''.
| version = [[214.0]]
}}
{{Server config variable
| status = deprecated
| name = -oldsaveformat
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Used to force the old save format. The default save format is based on the ''newsaveformat'', which is approximately 4x faster and 50% smaller. If you want to use the old save format, launch with this option.
| version = [[244.6]]
}}
{{Server config variable
| status = deprecated
| name = -PVPDisablePenetratingHits
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Disable the [[Tek Railgun (Aberration)|Tek Railgun]]'s ability to shoot through walls. ''Replaced by <code>-DisableRailgunPVP</code> in [[349.1]] and no longer present in the executable code''.
| version = [[345.25]]
}}
{{Server config variable
| status = deprecated
| name = -StructureDestructionTag=<newBiomesStructuresZones>
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = One-Time Auto-Structure Demolish <code><newBiomesStructuresZones></code> Zones: To do this, which you can only execute once after updating to patch 216, run your server or game with this option. <code><newBiomesStructuresZones></code> could be <code>DestroySwampSnowStructures</code> once updating to version [[216.0]] or <code>DestroyRedwoodSnowStructures</code> once updating to version [[243.0]]. ''Deprecated and no longer present in the executable code''.
| version = [[216.0]]
}}
{{Server config variable
| status = deprecated
| name = -usecache
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = ~70% faster loading speed option. Choose "Experiment Fast Load Cache" launch option (use add "-usecache" to your launch command line manually). After the first & second times that you start the game & load will be still be slow, but the third time onwards will be fast. ''Deprecated since [[188.2]] is enabled by default on all maps and no longer present in the executable code.''
| version = [[188.1]]
}}
{{Server config variable
| status = deprecated
| name = -vday
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Activated Valentine’s Day Event. ''Deprecated, use <code>-ActiveEvent=vday</code> instead.''
| version = [[254.93]]
}}
|-
! colspan="6" | Undocumented options
{{Server config variable
| name = -AllowChatSpam
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Allows players to spam messages in server chat. '''Note:''' it may interfere with external plugins allowing cross-map chat. Undocumented by Wildcard.
}}
{{Server config variable
| name = -BackupTransferPlayerDatas
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Fully mimic official (non-legacy) servers on handling player characters data, adding permanent character profile backup files separately from the world-save. Requires <code>-usestore</code> option. '''This option must be used to fully re-enable character profile backups when playing the [[2023_official_server_save_files|Official 2023 Server saves]] backups.''' Undocumented by Wildcard.
}}
{{Server config variable
| name = -BattlEyeServerRecheck
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Adds at every server tick an extra check about BattleEye service status. Used in official clusters. Undocumented by Wildcard.
}}
{{Server config variable
| name = -converttostore
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Converts non stored data (default save behaviour of unofficial and legacy) to stored data (default save behaviour of official non-legacy) at next world-save. It will not delete older character profile saves. Requires <code>-usestore</code>. Undocumented by Wildcard.
}}
{{Server config variable
| name = -CustomAdminCommandTrackingURL=<URL>
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Allows server to specify custom URL to post admin commands tracking. Supports HTTP protocols only (HTTPS is not supported). If the URL is empty the metrics should be written in a file named ''ArkCustomAdminCommandTrackingURL.txt'' in the game save directory. Undocumented by Wildcard.
}}
{{Server config variable
| name = -CustomMerticsURL=<URL>
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Allows server to specify custom URL to post server metrics statistics. Supports HTTP protocols only (HTTPS is not supported). If the URL is empty the metrics should be written in a file named ''ArkCustomMetricsURL.txt'' in the game save directory. Undocumented by Wildcard.
}}
{{Server config variable
| name = -CustomNotificationURL=<URL>
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Allows server custom notification broadcast using the server message feature. Supports HTTP protocols only (HTTPS is not supported). Used in official clusters (e.g.: http://arkdedicated.com/pcnotification.html). Undocumented by Wildcard.
}}
{{Server config variable
| name = -dedihibernation
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Forces wild creatures to spawn in hibernation mode. Requires <code>-preventhibernation</code> to not be set (which doesn't have effects anyway in dedicated servers). Undocumented by Wildcard.
}}
{{Server config variable
| name = -disableCharacterTracker
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Used to disable character tracking. Alternatively, this option can be configured with <code>UseCharacterTracker</code> in file [[#GameUserSettings.ini|GameUserSettings.ini]] under <code>[ServerSettings]</code> section (note that the argument from command line has priority over the value set in [[#GameUserSettings.ini|GameUserSettings.ini]]). Undocumented by Wildcard.
}}
{{Server config variable
| name = -DisableDupeLogDeletes
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Prevents <code>-ForceDupeLog</code> to take effect. Undocumented by Wildcard.
}}
{{Server config variable
| name = -DormancyNetMultiplier=<float>
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Multiplier for dormancy rate, controlling when an NPC can enter in a low network bandwidth dormant mode. If value is less or equal to 1.0 requires <code>-nodormancythrottling</code>. Default is set at 1.0. Undocumented by Wildcard.
}}
{{Server config variable
| name = -EnableOfficialOnlyVersioningCode
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Mimics official versioning code-path on ARK data items upload/download behaviour, deleting them if there is a version mismatch. Undocumented by Wildcard.
}}
{{Server config variable
| name = -EnableSteelShield
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Used to enable [https://enterprise.nitrado.net/steelshield/ SteelShield] anti-DDOS protection. '''This only works when your server is hosted by Nitrado.''' Undocumented by Wildcard.
}}
{{Server config variable
| name = -EnableVictoryCoreDupeCheck
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Should enabled an extra anti-duping check on tamed creatures. Undocumented by Wildcard.
}}
{{Server config variable
| name = -forcedisablemeshchecking
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Disables mesh checking on foliage re-population, not set on official. Undocumented by Wildcard.
}}
{{Server config variable
| name = -ForceDupeLog
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Forces dupe logs. Requires <code>-DisableDupeLogDeletes</code> to be not set, enabled behaviour on official. Undocumented by Wildcard.
}}
{{Server config variable
| name = -forceuseperfthreads
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Forces the use of performance threads. Undocumented by Wildcard.
}}
{{Server config variable
| name = -ignoredupeditems
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = If a duped item is detected in an inventory, this will be ignored and not removed. This option does not disable duped items checks. Undocumented by Wildcard.
}}
{{Server config variable
| name = -ip=<ipv4_address>
| inASA = Yes
| inASE = no
| default = 
| cli = skip
| info = Option found in Nitrado server setup. Expected to work in a <code>-MULTIHOME</code> context. Undocumented by Wildcard.
}}
{{Server config variable
| name = -MaxConnectionsPerIP=<integer>
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Sets a limit of users connected to the server with the same IP at the same time. By default is -1 (virtually unlimited). If an IP reaches such value, any players trying to connect next with such IP will get a connect timeout. Undocumented by Wildcard.
}}
{{Server config variable
| name = -NitradoQueryPort
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Specifics unknown. Undocumented by Wildcard.
}}
{{Server config variable
| name = -nitradotest2
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Limits the tick-rate to 3 when there are no players connected. This can have a noticeable impact on CPU for large worlds since there will be less calls per second (less computation/memory allocation, etc.). But as its name clearly state, it seems to be an internal test so use at your own risks. Undocumented by Wildcard.
}}
{{Server config variable
| name = -NoAI
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = This will disable adding AI controller to creatures. Undocumented by Wildcard.
}}
{{Server config variable
| name = -NoDinosExceptForcedSpawn
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Prevents wild creatures from being spawned, except for forced spawns. You cannot use this option if <code>-NoDinos</code> is used. Enabling this option will make the following options unusable: <code>-NoDinosExceptStreamingSpawn</code>, <code>-NoDinosExceptManualSpawn</code> and <code>-NoDinosExceptWaterSpawn</code>. ''<b>Note:</b> wild creatures that were on the server before adding this argument will not get automatically destroyed (you can destroy them with console command <code>cheat destroywilddinos</code>)''. Undocumented by Wildcard.
}}
{{Server config variable
| name = -NoDinosExceptStreamingSpawn
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Prevents wild creatures from being spawned, except for streaming spawns. You cannot use this option if one of the following options is used: <code>-NoDinos</code> or <code>-NoDinosExceptForcedSpawn</code>. Enabling this option will make the following options unusable: <code>-NoDinosExceptManualSpawn</code> and <code>-NoDinosExceptWaterSpawn</code>. ''<b>Note:</b> wild creatures that were on the server before adding this argument will not get automatically destroyed (you can destroy them with console command <code>cheat destroywilddinos</code>)''. Undocumented by Wildcard.
}}
{{Server config variable
| name = -NoDinosExceptManualSpawn
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Prevents wild creatures from being spawned, except for manual spawns. You cannot use this option if one of the following options is used: <code>-NoDinos</code>, <code>-NoDinosExceptForcedSpawn</code> or <code>-NoDinosExceptStreamingSpawn</code>. Enabling this option will make option <code>-NoDinosExceptWaterSpawn</code> unusable. ''<b>Note:</b> wild creatures that were on the server before adding this argument will not get automatically destroyed (you can destroy them with console command <code>cheat destroywilddinos</code>)''. Undocumented by Wildcard.
}}
{{Server config variable
| name = -NoDinosExceptWaterSpawn
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Prevents wild creatures from being spawned, except for water spawns. You cannot use this option if one of the following options is used: <code>-NoDinos</code>, <code>-NoDinosExceptForcedSpawn</code>, <code>-NoDinosExceptStreamingSpawn</code> or <code>-NoDinosExceptManualSpawn</code>. ''<b>Note:</b> wild creatures that were on the server before adding this argument will not get automatically destroyed (you can destroy them with console command <code>cheat destroywilddinos</code>)''. Undocumented by Wildcard.
}}
{{Server config variable
| name = -nodormancythrottling
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Required to make <code>DormancyNetMultiplier</code> to work if its value is less than 1.0. Undocumented by Wildcard.
}}
{{Server config variable
| name = -noperfthreads
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Disables performance threads. Undocumented by Wildcard.
}}
{{Server config variable
| name = -nosound
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Disables sounds to improve performances. Undocumented by Wildcard.
}}
{{Server config variable
| name = -onethread
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Disables multithreading. Undocumented by Wildcard.
}}
{{Server config variable
| name = -parseservertojson
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Dumps at start-up the loading server save to a JSON file located at ''Saved/Logs/AuditLogs'' and named ''AuditLog_%s.json'' where ''%s'' is the current UTC date-time. Useful for debugging purposes. '''Note:''' this will result in a big slowdown of server start-up and a creation of a big JSON file. On big server saves this may result on server hang before start-up even for hours. Undocumented by Wildcard.
}}
{{Server config variable
| name = -pauseonddos
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Should pause/halt the game if a DDoS attack by a performed by players is detected. Undocumented by Wildcard.
}}
{{Server config variable
| name = -PreventTotalConversionSaveDir
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Prevents total conversion mods to use a dedicated save directory, forcing the default save folder. Undocumented by Wildcard.
}}
{{Server config variable
| name = -pveallowtribewar
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Opposite of <code>pvedisallowtribewar</code>, allows tribes to officially declare war on each other for mutually-agreed-upon period of time. It is suggested to use the ''Game.ini'' <code>bPvEAllowTribeWar</code> instead. Undocumented by Wildcard.
}}
{{Server config variable
| name = -ReloadedForBackup
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Marks the world-save to load as a trial backup save, preventing to properly load and continue the game from such world save. Useful to spare time when doing tests not needing to fully load the world-save. Undocumented by Wildcard.
}}
{{Server config variable
| name = -ServerIP=<ipv4_address>
| inASA = Yes
| inASE = no
| default = 
| cli = skip
| info = Option found in Nitrado server setup. Expected to work in a <code>-MULTIHOME</code> context. Undocumented by Wildcard.
}}
{{Server config variable
| name = -ServerPlatform=<plat1>[+<plat2>[...]]
| inASA = Yes
| inASE = No
| default = 
| cli = skip
| info = Allows the server to accept specified platforms. Options are <code>PC</code> for Steam, <code>PS5</code> for PlayStation 5, <code>XSX</code> for XBOX, <code>WINGDK</code> for Microsoft Store, <code>ALL</code> for crossplay between PC (Steam and Windows Store) and all consoles. With no option specified, crossplay is assumed and <code>ALL</code> platforms will be allowed. Undocumented by Wildcard.<br>Examples: <code>-ServerPlatform=PC</code> or <code>-ServerPlatform=PC+XSX+PS5</code>
}}
{{Server config variable
| name = -UnstasisDinoObstructionCheck
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = Should prevent creatures from ghosting through meshes/structures on re-render. Undocumented by Wildcard.
}}
{{Server config variable
| name = -UseTameEffectivenessClamp
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Prevents taming effectiveness never going over 100%. Undocumented by Wildcard.
}}
{{Server config variable
| name = -UseServerNetSpeedCheck
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| info = It should avoid players to accumulate too much movements data per server tick, discarding the last if those are too many. Can be turned on also with <code>-GUseServerNetSpeedCheck</code> in dynamic config. Enabled on official clusters. May be useful for servers having too many players and stressed on CPU side. Undocumented by Wildcard.
}}
|-
! colspan="6" | Clients only options (ineffective on servers)
{{Server config variable
| name = -allowansel
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Activates [[NVIDIA Ansel]] support in single player. Can be enabled on multiplayer servers with <code>-ServerAllowAnsel</code> since patch [[278.0]] for servers.''
| version = [[246.0]]
}}
{{Server config variable
| name = -d3d10 -dx10 -sm4
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]]
| info = Any of these options forces the use DirectX 10 using shader model 4. This will reduce the graphics engine to a lesser version, reducing some graphics, but raising the frame-rate on older and weaker hardware. '''Using this option on server will not override clients' behaviour, nor bring any benefit.'''
}}
{{Server config variable
| name = -d3d11 -dx11 -sm5
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]]
| info = Any of these options forces the use DirectX 11 using shader model 5. This is the default rendering mode. '''Using this option on server will not override clients' behaviour, nor bring any benefit.'''
}}
{{Server config variable
| status = deprecated
| name = -d3d12 -dx12
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]]
| info =  In {{ItemLink|ARK: Survival Evolved|ASE}} any of these options were used to force the DirectX 12 back-end. This mode no longer works due the removal of D3D12RHI module from the engine. '''Trying to use this option will crash the client at start-up.''' In {{ItemLink|ARK: Survival Ascended|ASA}} DirectX 12 back-end is enabled by default and using this option will take no effect.
}}
{{Server config variable
| status = deprecated
| name = -game
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]] Editor
| info = This is an option to start the Unreal Engine Editor (i.e.: the Ark DevKit) the uncooked assets. There is no use for the command on either the client or server.
}}
{{Server config variable
| status = deprecated
| name = -LANPLAY
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]] Server
| info = Unreal Engine option for servers. On ARK MaxClientRate is hardcoded to 60K and MaxInternetClientRate is limited to 15K per client. The Unreal Engine code that should double the server updates for LAN play is not present in executable bits. This options has no effect on ARK servers nor clients.
}}
{{Server config variable
| status = deprecated
| name = ?listen
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]] Server
| info = This is an option to start the Unreal Engine Editor (i.e.: the Ark DevKit) as an undedicated server using the uncooked assets. There is no use for the command on either the client or server.
}}
{{Server config variable
| status = deprecated
| name = -log
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]]
| info = This is an option to start the log window console in the Unreal Engine Editor (i.e.: the Ark DevKit) with <code>-game</code> option. There is no use for the command on either the client or server.
}}
{{Server config variable
| name = -lowmemory
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Launch options that reduces graphics and audio effects to save about 800 MB RAM, likely enabling 4GB RAM players to get past infinite-loading screens. '''Using this option on server will not override clients' behaviour, nor bring any benefit.'''
| version = [[170.0]]
}}
{{Server config variable
| status = deprecated
| name = -noaafonts
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Removed fonts anti-aliasing. ''Option is now deprecated.''
}}
{{Server config variable
| name = -nomansky
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Lots of detailed sky features are disabled, such as clouds and starry night sky. This decreases all of those but you can still have them. You still obtain the stars and sun. '''Using this option on server will not override clients' behaviour, nor bring any benefit.'''
| version = [[171.3]]
}}
{{Server config variable
| name = -nomemorybias
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Reduced client game memory usage by about 600 MB system and 600 MB GPU RAM! (all meshes now stream LOD's dynamically). This could potentially have a small runtime performance impact, so to use the old method (more RAM usage but no potential performance loss), launch with this option. '''Using this option on server will not override clients' behaviour, nor bring any benefit.'''
| version = [[244.6]]
}}
{{Server config variable
| name = -norhithread
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = Used to disable Direct3D 11 parallel batch commands. It may be useful only for clients that are having Device Lost crashes, at a performance cost. '''Using this option on server will not override clients' behaviour, nor bring any benefit.'''
}}
{{Server config variable
| name = -nosteamclient
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| info = Used automatically by client (''ShooterGame'') to launch the server from ''HOST\LOCAL'' menu. Launching the server (''ShooterGameServer'') directly or through a script using this option will have no effect and will be ignored.
}}
{{Server config variable
| status = deprecated
| name = -opengl -opengl3 -opengl4
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]]
| info = Any of these options forces the use OpenGL. Theses mode no longer works due the removal of OpenGL RHI mode from the engine. '''Trying to use this option will crash the client at start-up.'''
}}
{{Server config variable
| name = -PreventHibernation
| inASA = Unknown
| inASE = Yes
| default = 
| cli = skip
| info = In both '''Single Player''' and '''Non-Dedicated''' sessions, creatures in inactive zones are in hibernation instead of in stasis. Use this option to prevent hibernation at the cost of performance and memory usage. '''Using this option on server will not override clients' behaviour, nor bring any benefit.'''
| version = [[259.0]]
}}
{{Server config variable
| name = -ForceIgnoreSingleplayerSpawnRangeCheck
| inASA = Yes
| inASE = Unknown
| default = 
| cli = skip
| info = Forces '''Single Player''' and '''Non-Dedicated''' sessions to ignore spawn range checks, allowing them to behave like a dedicated server in terms of creature spawns. Use this option to enable creature spawns across the map, even when the player is not nearby. Enabling this option may impact performance and memory, including significant initial stuttering when creating a new world.
}}
{{Server config variable
| status = deprecated
| name = -server
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]] Editor
| info = This is an option to start the Unreal Engine Editor (i.e.: the Ark DevKit) as a dedicated server using the uncooked assets. There is no use for the command on either the client or server.
}}
{{Server config variable
| status = deprecated
| name = -USEALLAVAILABLECORES
| inASA = No
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]] Editor
| info = This is an Unreal Engine parameter, used while compiling shaders and cooking assets (i.e.: the ARK DevKit). This option is useless for dedicated servers and clients.
}}
{{Server config variable
| status = deprecated
| name = -vulkan
| inASA = Yes
| inASE = Yes
| default = 
| cli = skip
| version = [[File:Unreal.svg|20px]]
| info = Unreal Engine option to force Vulkan rendering path. In {{ItemLink|ARK: Survival Evolved|ASE}} was added first for the defunct Google Stadia version. In {{ItemLink|ARK: Survival Ascended|ASA}} is not supported as well. '''Trying to use this option will crash the client at start-up.'''
}}
|}
[[File:Unreal.svg|20px]] means this is a native feature of the Unreal Engine that is or was supported by the game at some point and has no specific "Ark version" since it was inherently supported in the Unreal Engine by default, unless proactive actions are used to disable it by the developers.
</div>

<!--end of command line-->

==Configuration Files==
Most options can also be specified in the game configuration files. The location of the configuration files varies by platform:
<div class="widetable">
{| class="wikitable"
|-
!Platform!!Configuration File!!Location
|-
| Linux||GameUserSettings.ini||<code>ShooterGame/Saved/Config/LinuxServer/</code>
|-
|Linux||Game.ini||<code>ShooterGame/Saved/Config/LinuxServer/</code>
|-
| Windows||GameUserSettings.ini||<code>ShooterGame\Saved\Config\WindowsServer\</code>
|-
|Windows||Game.ini||<code>ShooterGame\Saved\Config\WindowsServer\</code>
|}
</div>

Single Player and Non-Dedicated can access their game files through the Steam and Epic launchers:

Steam: Right click on the game name in ''Library'' and choose ''Properties...''. Select the ''Installed Files'' tab, then select ''Browse''
Epic: View the game from the ''Library'' tab (not Quick Launch). Click the three dots icon, then select ''Manage''. Under ''Installation'', click the folder icon.
The paths listed in the table above can then be followed from these locations. '''Ensure the game has been closed before making any changes to the INI files. The changes will be removed if edited while the game is running.'''

The '''GameUserSettings.ini''' file contains options for both the game client and the game server. Options for the game client are not used by the server.

The '''Game.ini''' file is used mostly for more advanced modifications, such as changing engram points or XP rewarded per level, disabling specific content, or rebalancing depending on player tastes.

For all supported platforms, options are listed one per line using the same basic syntax:

 '''<option>'''='''<value>'''

All options in the configuration file require a value. If an option is not listed in the configuration file, its default value is used automatically.

To configure a game server with the same configuration as shown in [[#Syntax|Command Line Syntax]] above:

 ServerCrosshair=True
 AllowThirdPersonPlayer=True
 ShowMapPlayerLocation=False
 TheMaxStructuresInRange=1000

With these options in the configuration file, the server can be launched with a much shorter command line. For Linux:

 ./ShooterGameServer TheIsland

For Windows:

 .\ShooterGameServer.exe TheIsland
<!--start of GameUserSettings.ini-->
=== GameUserSettings.ini===
''GameUserSettings.ini'' file holds most server as well most of basic and intermediate gameplay settings.

Although not suggested, settings marked with {{Checkmark|Y}} in the '''CMD''' column can be used in the command line too, appending them with the <code>?</code> syntax. In contrast, settings marked with {{Checkmark|N}} cannot be moved in the command line.

Beginners can easily create this file using the [http://ini.arkforum.de/index.php?lang=en&mode=all ini-Generator].
<!--start of [ServerSettings]-->
==== [ServerSettings] ====
The following options must come under the <code>[ServerSettings]</code> section of <code>GameUserSettings.ini</code>:

<div class="widetable">
{| class="wikitable config-table"
|-
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Variable
! Description
! {{HoverNote|CMD|If ticked, this option may be specified in the server's command line instead, after map name.}}
! Since patch
{{Server config variable
| name = ActiveMods
| type = list of mod IDs, comma-separated with no spaces, in a single line (for example: <code>ModID1,ModID2,ModID3</code>)
| default = 
| cli = No
| inASA = Yes
| inASE = Yes
| Specifies the order and which mods are loaded. ModIDs are comma separated and in one line. Priority is in descending order (the left-most ModID hast the highest priority). Alternatively, but not suggested, the command line <code>?GameModIds{{=}}</code> can be used.
| version = [[190.0]]
}}
{{Server config variable
| name = ActiveMapMod
| type = mod ID for currently active mod map
| default = 
| cli = Unknown
| inASA = Yes
| inASE = Unknown
| Specifies which mod map is loaded.
}}{{Server config variable
| name = ActiveTotalConversion
| type = ModID
| default = 
| cli = No
| inASA = No
| inASE = Yes
| Used to specify a total conversion mod (e.g.: [[Primitive Plus]]). Alternatively, the command line option <code>-TotalConversionMod{{=}}<ModID></code> can also be set.
}}
{{Server config variable
| name = AdminLogging
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, logs all admin commands to in-game chat.
| version = [[206.0]]
}}
{{Server config variable
| name = AllowAnyoneBabyImprintCuddle
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, allows anyone to "take care" of a baby creatures (cuddle etc.), not just whomever imprinted on it.
| version = [[242.0]]
}}
{{Server config variable
| name = AllowCaveBuildingPvE
| type = boolean
| default = False
| cli = No
| inASA = Yes
| inASE = Yes
| If <code>True</code>, allows building in caves when [[PvE]] mode is also enabled. '''Note:''' no more working in command-line options before patch [[241.5]].
| version = [[194.0]]
}}
{{Server config variable
| name = AllowCaveBuildingPvP
| type = boolean
| default = True
| cli = No
| inASA = Yes
| inASE = Yes
| If <code>False</code>, prevents building in caves when [[PvP]] mode is also enabled.
| version = [[326.13]]
}}
{{Server config variable
| name = AllowCrateSpawnsOnTopOfStructures
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, allows from-the-air Supply Crates to appear on top of Structures, rather than being prevented by Structures.
| version = [[253.92]]
}}
{{Server config variable
| name = AllowCryoFridgeOnSaddle
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = No
| If <code>True</code>, allows cryofridges to be built on platform saddles and rafts.
}}
{{Server config variable
| name = AllowFlyerCarryPvE
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, allows flying creatures to pick up wild creatures in [[PvE]].
| version = [[173.0]]
}}
{{Server config variable
| name = AllowFlyingStaminaRecovery
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, allows server to recover Stamina when standing on a Flyer.
| version = [[256.0]]
}}
{{Server config variable
| name = AllowHideDamageSourceFromLogs
| type = boolean
| default = True
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>False</code>, shows the damage sources in tribe logs.
| version = [[278.0]]
}}
{{Server config variable
| name = AllowHitMarkers
| type = boolean
| default = True
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>False</code>, disables optional markers for ranged attacks.
| version = [[245.0]]
}}
{{Server config variable
| name = AllowIntegratedSPlusStructures
| type = boolean
| default = True
| cli = Yes
| inASA = Unknown
| inASE = Yes
| if <code>False</code>, disables all of the new S+ structures (intended mainly for letting unofficial servers that want to keep using the S+ mod version to keep using that without a ton of extra duplicate structures).
| version = [[293.100]]
}}
{{Server config variable
| name = AllowMultipleAttachedC4
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, allows to attach more than one [[C4]] per creature.
| version = [[256.0]]
}}
{{Server config variable
| name = AllowRaidDinoFeeding
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, allows Titanosaurs to be permanently tamed (namely allow them to be fed). '''Note:''' in The Island only spawns a maximum of 3 Titanosaurs, so 3 tamed ones should ultimately block any more ones from spawning.
| version = [[243.0]]
}}
{{Server config variable
| name = AllowSharedConnections
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, allows family sharing players to connect to the server.
| version = [[312.73]]
}}
{{Server config variable
| name = AllowTekSuitPowersInGenesis
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, enables TEK suit powers in [[Genesis: Part 1]].
| version = [[306.79]]
}}
{{Server config variable
| name = AllowThirdPersonPlayer
| type = boolean
| default = True
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>False</code>, disables third person camera allowed by default on all dedicated servers.
| version = [[265.16]]
}}
{{Server config variable
| name = AlwaysAllowStructurePickup
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code> disables the timer on the quick pick-up system.
| version = [[293.100]]
}}
{{Server config variable
| name = AlwaysNotifyPlayerLeft
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, players will always get notified if someone leaves the server
}}
{{Server config variable
| name = ArmadoggoDeathCooldown
| type = float
| default = 3600
| cli = Yes
| inASA = Yes
| inASE = No
| Overrides the cooldown for Armadoggo to reappear after taking fatal damage (in seconds), default is set to 1 hour. Must be greater than 0.
| version = {{gameicon|sa}}<br>
[[ARK: Survival Ascended/Patch/57.12|57.12]]
}}
{{Server config variable
| name = AutoDestroyDecayedDinos
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, auto-destroys claimable decayed tames on load, rather than have them remain around as claimable. '''Note:''' after patch [[273.691]], in [[PvE]] mode the tame auto-unclaim after decay period has been disabled in official [[PvE]].
| version = [[255.0]]
}}
{{Server config variable
| name = AutoDestroyOldStructuresMultiplier
| type = float
| default = 0.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Allows auto-destruction of structures only after sufficient "no nearby tribe" time has passed (defined as a multiplier of the Allow Claim period). To enable it, set it to 1.0. Useful for servers to clear off abandoned structures automatically over time. Requires <code>-AutoDestroyStructures</code> command line option to work. The value scales with each structure decay time.
| version = [[222.0]]
}}
{{Server config variable
| name = AutoSavePeriodMinutes
| type = float
| default = 15.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Set interval for automatic saves. Setting this to 0 will cause constant saving.
| version = [[181.0]]
}}
{{Server config variable
| name = BanListURL
| type = string with a URL
| default = 
| cli = Yes
| inASA = Yes
| inASE = Yes
| Sets the global ban list. Must be enclosed in double quotes. The list is fetched every 10 minutes (to check if there are new banned IDs).<br/>
{{Gamelink|ase}}: Official ban list URL is [http://arkdedicated.com/banlist.txt http://arkdedicated.com/banlist.txt] (before [[279.233]] the URL was [http://playark.com/banlist.txt http://playark.com/banlist.txt]). '''Note:''' it supports the HTTP protocol only (HTTPS is not supported).<br/>
{{Gamelink|ase}}: Official ban list URL is [https://cdn2.arkdedicated.com/asa/BanList.txt https://cdn2.arkdedicated.com/asa/BanList.txt]. '''Note:''' it supports both HTTP and HTTPS protocols.
| version = [[201.0]]
}}
{{Server config variable
| name = bForceCanRideFliers
| type = boolean
| default = 
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, allows flyers to be used on maps where they normally are disabled. '''Note:''' if you set it to <code>False</code> it will disable flyers on any map.
| version = [[306.53]]
}}
{{Server config variable
| name = ClampItemSpoilingTimes
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, clamps all spoiling times to the items' maximum spoiling times. Useful if any infinite-spoiling exploits were used on the server and you wish to clean them up. Could potentially cause issues with mods that alter spoiling time, hence it is an option.
| version = [[254.944]]
}}
{{Server config variable
| name = ClampItemStats
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, enables stats clamping for items. See [[#ItemStatClamps|ItemStatClamps]] for more info.
| version = [[255.0]]
}}
{{Server config variable
| name = ClampResourceHarvestDamage
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, limit the damage caused by a tame to a resource on harvesting based on resource remaining health.  '''Note:''' enabling this setting may result in sensible resource harvesting reduction using high damage tools or creatures.
| version = [[182.0]]
}}
{{Server config variable
| name = CosmeticWhitelistOverride
| type = string with a URL
| default = 
| cli = Yes
| inASA = Yes
| inASE = No
| URL to a comma-separated list of whitelisted custom cosmetics, in this format: <code><nowiki>Mod ID|Enable Dynamic Download (0/1)|Allow non-dataonly blueprints(0/1)|Calculated package CRC.</nowiki></code> See this [https://survivetheark.com/index.php?/forums/topic/736592-enable-custom-cosmetics-validation-on-unofficial/ post] for details.
| version = [[34.49]]
}}
{{Server config variable
| name = CosmoWeaponAmmoReloadAmount
| type = float
| default = 1
| cli = Yes
| inASA = Yes
| inASE = No
| Determines how much ammo is given as the Cosmo's webslinger reloads over time.
| version = {{gameicon|sa}}<br>[[ARK: Survival Ascended/Patch/52.5|52.5]]
}}
{{Server config variable
| name = CustomDynamicConfigUrl
| type = string with a URL
| default = 
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Direct link to a live dynamicconfig.ini file (http://arkdedicated.com/dynamicconfig.ini), allowing live changes of the supported options without the need of server restart, as well defining custom colour set for wild creatures' spawns. '''Note:''' requires <code>-UseDynamicConfig</code> command line option, <code>string with a URL</code> must use the ''HTTP'' protocol (''HTTPS'' is not supported) and inside quotes if used in command line. Check the [[#DynamicConfig|DynamicConfig]] section for the supported settings.
| version = [[307.2]]
}}
{{Server config variable
| name = CustomLiveTuningUrl
| type = string with a URL
| default = 
| cli = Yes
| inASA = Yes
| inASE = Yes
| Direct link to the live tuning file. For more information on how to use this system check out the official announcement: https://survivetheark.com/index.php?/forums/topic/569366-server-configuration-live-tuning-system.<br/>
{{Gamelink|ase}}: official servers use [http://arkdedicated.com/DefaultOverloads.json http://arkdedicated.com/DefaultOverloads.json]. '''Note:''' <code>string with a URL</code> must use the ''HTTP'' protocol (''HTTPS'' is not supported) and inside quotes if used in command line.<br/>
{{Gamelink|asa}}: official servers use [https://cdn2.arkdedicated.com/asa/livetuningoverloads.json https://cdn2.arkdedicated.com/asa/livetuningoverloads.json]. '''Note:''' <code>string with a URL</code> must use inside quotes if used in command line.
| version = [[314.3]]
}}
{{Server config variable
| name = DayCycleSpeedScale
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the passage of time in the ARK, controlling how often day changes to night and night changes to day. The default value <code>1</code> provides the same cycle speed as the single player experience (and the official public servers). Values lower than 1 slow down the cycle; higher values accelerate it. Base time when value is 1 appears to be 1-minute real time equals approx. 28-minutes game time. Thus, for an approximate 24-hour day/night cycle in game, use .035 for the value.
| version = [[179.0]]
}}
{{Server config variable
| name = DayTimeSpeedScale
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the passage of time in the ARK during the day. This value determines the length of each day, relative to the length of each night (as specified by <code>NightTimeSpeedScale</code>). Lowering this value increases the length of each day.
| version = [[179.0]]
}}
{{Server config variable
| name = DestroyUnconnectedWaterPipes
| type = boolean
| default = False
| cli = Yes
| inASA = No
| inASE = Yes
| If <code>True</code>, after two days real-time the pipes will auto-destroy if unconnected to any non-pipe (directly or indirectly) and no allied player is nearby.
| version = [[252.4]]
}}
{{Server config variable
| name = DifficultyOffset
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the [[Difficulty|difficulty level]].
}}
{{Server config variable
| name = DinoCharacterFoodDrainMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for creatures' food consumption. Higher values increase food consumption (creatures get hungry faster). It also affects the taming-times.
| version = [[179.0]]
}}
{{Server config variable
| name = DinoCharacterHealthRecoveryMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for creatures' health recovery. Higher values increase the recovery rate (creatures heal faster).
| version = [[179.0]]
}}
{{Server config variable
| name = DinoCharacterStaminaDrainMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for creatures' stamina consumption. Higher values increase stamina consumption (creatures get tired faster).
| version = [[179.0]]
}}
{{Server config variable
| name = DinoCountMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Specifies the scaling factor for creature spawns. Higher values increase the number of creatures spawned throughout the ARK.
| version = [[179.0]]
}}
{{Server config variable
| name = DinoDamageMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the damage wild creatures deal with their attacks. The default value <code>1</code> provides normal damage. Higher values increase damage. Lower values decrease it.
| version = [[179.0]]
}}
{{Server config variable
| name = DinoResistanceMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the resistance to damage wild creatures receive when attacked. The default value <code>1</code> provides normal damage. Higher values decrease resistance, increasing damage per attack. Lower values increase it, reducing damage per attack. A value of 0.5 results in a creature taking half damage while a value of 2.0 would result in a creature taking double normal damage.
| version = [[179.0]]
}}
{{Server config variable
| name = DestroyTamesOverTheSoftTameLimit
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = No
| version = [[ARK: Survival Ascended/Patch/33.15|ASA 33.15]]
| Dinos above the Soft Tame Server Limit will be marked “For Cryo” and display an icon and a timer indicating how soon they need to be cryopodded before they are automatically destroyed. Dinos marked and dinos destroyed by this system will be logged in the tribe log. This differs from the standard tame limit, as it still allows new players to join servers and tame creatures above the 5000 limit, whereas previously, they would not have been able to tame or breed creatures at the server cap. You can define the soft tame limit with <code>MaxTamedDinos_SoftTameLimit</code>. You can define the time (in seconds) before tame's autodestruction with <code>MaxTamedDinos_SoftTameLimit_CountdownForDeletionDuration</code>.
}}
{{Server config variable
| name = DisableCryopodEnemyCheck
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = No
| If <code>True</code>, allows cryopods to be used while enemies are nearby.
}}
{{Server config variable
| name = DisableCryopodFridgeRequirement
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = No
| If <code>True</code>, allows cryopods to be used without needing to be in range of a powered cryofridge.
}}
{{Server config variable
| name = DisableDinoDecayPvE
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, disables the creature decay in [[PvE]] mode. '''Note:''' after patch [[273.691]], in [[PvE]] mode the creature auto-unclaim after decay period has been disabled.
| version = [[206.0]]
}}
{{Server config variable
| name = DisableImprintDinoBuff
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, disables the creature imprinting player Stat Bonus. Where whomever specifically imprinted on the creature, and raised it to have an Imprinting Quality, gets extra Damage/Resistance buff.
| version = [[242.0]]
}}
{{Server config variable
| name = DisablePvEGamma
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents use of console command "gamma" in [[PvE]] mode.
| version = [[207.0]]
}}
{{Server config variable
| name = DisableStructureDecayPvE
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, disables the gradual [[Building#Auto-decay|auto-decay]] of player structures.
| version = [[173.0]]
}}
{{Server config variable
| name = DisableWeatherFog
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, disables fog.
| version = [[278.0]]
}}
{{Server config variable
| name = DontAlwaysNotifyPlayerJoined
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, globally disables player joins notifications.
}}
{{Server config variable
| name = EnableExtraStructurePreventionVolumes
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, disables building in specific resource-rich areas, in particular setup on [[The Island]] around the major mountains.
| version = [[242.0]]
}}
{{Server config variable
| name = EnablePvPGamma
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, allows use of console command "gamma" in [[PvP]] mode.
| version = [[174.3]]
}}
{{Server config variable
| name = ExtinctionEventTimeInterval
| type = seconds
| default = 
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Used to enable the ''extinction mode'' ([[Game_Modes#ARKpocalypse_(PC)|ARKpocalypse]]). The number is the time in seconds. Use 2592000 value for 30 days.
}}
{{Server config variable
| name = FastDecayUnsnappedCoreStructures
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, unsnapped foundations/pillars/fences/Tek Dedicated Storage will decay after the time stated by <code>FastDecayInterval</code> in ''Game.ini'' (default is 12 hours). Before [[259.0]], it set the decay time for such structures straight to 5 times faster.
| version = [[245.987]]
}}
{{Server config variable
| name = ForceAllStructureLocking
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, will default lock all structures.
| version = [[222.0]]
}}
{{Server config variable
| name = ForceGachaUnhappyInCaves
| type = boolean
| default = True
| cli = Yes
| inASA = Yes
| inASE = Unknown
| If <code>True</code>, Gachas will become unhappy within caves.
| version = {{gameicon|sa}}<br>
[[ARK: Survival Ascended/Patch/57.12|57.12]]
}}
{{Server config variable
| name = globalVoiceChat
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, voice chat turns global.
}}
{{Server config variable
| name = HarvestAmountMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for yields from all harvesting activities (chopping down trees, picking berries, carving carcasses, mining rocks, etc.). Higher values increase the amount of materials harvested with each strike.
| version = [[179.0]]
}}
{{Server config variable
| name = HarvestHealthMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the "health" of items that can be harvested (trees, rocks, carcasses, etc.). Higher values increase the amount of damage (i.e., "number of strikes") such objects can withstand before being destroyed, which results in higher overall harvest yields.
| version = [[179.0]]
}}
{{Server config variable
| name = IgnoreLimitMaxStructuresInRangeTypeFlag
| type = boolean
| default = False
| cli = No
| inASA = Yes
| inASE = Yes
| If <code>True</code>, removes the limit of 150 decorative structures (flags, signs, dermis etc.).
}}
{{Server config variable
| name = ImplantSuicideCD
| type = float
| default = 28800
| cli = No
| inASA = Yes
| inASE = No
| Defines the time (in seconds) a player must wait between 2 uses of the implant's "Respawn" feature.
}}
{{Server config variable
| name = ItemStackSizeMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Allow increasing or decreasing global item stack size, this means all default stack sizes will be multiplied by the value given (excluding items that have a stack size of 1 by default).
| version = [[291.100]]
}}
{{Server config variable
| name = KickIdlePlayersPeriod
| type = float
| default = 3600.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Time in seconds after which characters that have not moved or interacted will be kicked (if -EnableIdlePlayerKick as command line parameter is set). '''Note:''' although at code level it is defined as a floating-point number, it is suggested to use an integer instead.
| version = [[241.5]]
}}
{{Server config variable
| name = MaxCosmoWeaponAmmo
| type = float
| default = -1
| cli = Yes
| inASA = Yes
| inASE = No
| This will make the maximum ammo amount for the Cosmo's webslinger to a set number instead of it scaling with the Cosmo's level. The default of <code>-1</code> will enable scaling with level.
| version = {{gameicon|sa}}<br>[[ARK: Survival Ascended/Patch/52.5|52.5]]
}}
{{Server config variable
| name = MaxGateFrameOnSaddles
| type = integer
| default = 0
| cli = No
| inASA = Unknown
| inASE = Yes
| Defines the maximum amount of gateways allowed on platform saddles. A value of 2 would prevent players from placing more than 2 gateways on their platform saddles (used in Official [[PvP]] servers). This setting is not retroactive, meaning existing builds won't be affected. Set to 0 to not allow players to place gateways on platform saddles (used in Official [[PvE]] servers).
| version = [[312.65]]
}}
{{Server config variable
| name = MaxHexagonsPerCharacter
| type = integer
| default = 2000000000
| cli = No
| inASA = Unknown
| inASE = Yes
| Sets the max amount of Hexagon a Character can accumulate. Official set it to 2500000.
| version = [[312.65]]
}}
{{Server config variable
| name = MaxPersonalTamedDinos
| type = integer
| default = 0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Sets a per-tribe creature tame limit (500 on official PvE servers, 300 in official PvP servers). The default value of 0 disables such limit.
| version = [[255.0]]
}}
{{Server config variable
| name = MaxPlatformSaddleStructureLimit
| type = integer
| default = 75
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Changes the maximum number of platformed-creatures/rafts allowed on the ARK (a potential performance cost). Example: <code>MaxPlatformSaddleStructureLimit{{=}}10</code> would only allow 10 platform saddles/rafts across the entire ARK.
| version = [[212.1]]
}}
{{Server config variable
| name = MaxTamedDinos
| type = float
| default = 5000.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Sets the maximum number of tame creatures on a server, this is a global cap. '''Note:''' although at code level it is defined as a floating-point number, it is suggested to use an integer instead.
| version = [[191.0]]
}}
{{Server config variable
| name = MaxTamedDinos_SoftTameLimit
| type = integer
| default = 5000
| cli = Yes
| inASA = Yes
| inASE = No
| version = [[ARK: Survival Ascended/Patch/33.15|ASA 33.15]]
| Defines the server-wide soft tame limit. See <code>DestroyTamesOverTheSoftTameLimit</code> for more info.
}}
{{Server config variable
| name = MaxTamedDinos_SoftTameLimit_CountdownForDeletionDuration
| type = integer
| default = 604800
| cli = Yes
| inASA = Yes
| inASE = No
| version = [[ARK: Survival Ascended/Patch/33.15|ASA 33.15]]
| Defines the time (in seconds) for tame to get destroyed. See <code>DestroyTamesOverTheSoftTameLimit</code> for more info.
}}
{{Server config variable
| name = MaxTrainCars
| type = integer
| default = 8
| cli = Yes
| inASA = Yes
| inASE = No
| version = [[ARK: Survival Ascended/Patch/38.18|ASA 38.18]]
| Defines the maximum amount of carts a train cave have.
}}
{{Server config variable
| name = MaxTributeCharacters
| type = integer
| default = 10
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Slots for uploaded characters. Any value less than default will be reverted. '''Note:''' rising it may corrupt player/cluster data and lead to lose of all stored characters.
}}
{{Server config variable
| name = MaxTributeDinos
| type = integer
| default = 20
| cli = Yes
| inASA = Yes
| inASE = Yes
| Slots for uploaded creatures. Any value less than default will be reverted. '''Note:''' Some player claimed maximum 273 to be safe cap and more will corrupt profile/cluster and lead to lose of all stored creatures but it need to be checked
}}
{{Server config variable
| name = MaxTributeItems
| type = integer
| default = 50
| cli = Yes
| inASA = Yes
| inASE = Yes
| Slots for uploaded items and resources. Any value less than default will be reverted. '''Note:''' Some player claimed maximum 154 to be safe cap and more will corrupt profile/cluster and lead to lose of all stored items and resources but it need to be checked
}}
{{Server config variable
| name = NightTimeSpeedScale
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the passage of time in the ARK during night time. This value determines the length of each night, relative to the length of each day (as specified by <code>DayTimeSpeedScale</code>) Lowering this value increases the length of each night.
| version = [[179.0]]
}}
{{Server config variable
| name = NonPermanentDiseases
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, makes permanent diseases not permanent. Players will lose them if on re-spawn.
| version = [[242.3]]
}}
{{Server config variable
| name = NPCNetworkStasisRangeScalePlayerCountStart
| type = integer
| default = 0
| cli = No
| inASA = Unknown
| inASE = Yes
| Minimum number of online players when the NPC Network Stasis Range Scale override is enabled (requires inputting into INI, not command line). Used to override the NPC Network Stasis Range Scale (to scale server performance for more player), on default is set to 0, disabling it. Official set it to 24.
| version = [[254.7]]
}}
{{Server config variable
| name = NPCNetworkStasisRangeScalePlayerCountEnd
| type = integer
| default = 0
| cli = No
| inASA = Unknown
| inASE = Yes
| Maximum number of online players when <code>NPCNetworkStasisRangeScalePercentEnd</code> is reached (requires inputting into INI, not command line). Used to override the NPC Network Stasis Range Scale (to scale server performance for more player), on default is set to 0, disabling it. Official set it to 70
| version = [[254.7]]
}}
{{Server config variable
| name = NPCNetworkStasisRangeScalePercentEnd
| type = float
| default = 0.55000001
| cli = No
| inASA = Unknown
| inASE = Yes
| The Maximum scale percentage used when <code>NPCNetworkStasisRangeScalePlayerCountEnd</code> is reached (requires inputting into INI, not command line). Used to override the NPC Network Stasis Range Scale (to scale server performance for more player). Official set it to 0.5.
| version = [[254.7]]
}}
{{Server config variable
| name = OnlyAutoDestroyCoreStructures
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, prevents any non-core/non-foundation structures from auto-destroying (however they'll still get auto-destroyed if a floor that they're on gets auto-destroyed). Official [[PvE]] Servers used this option.
| version = [[245.989]]
}}
{{Server config variable
| name = OnlyDecayUnsnappedCoreStructures
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, only unsnapped core structures will decay. Useful for eliminating lone pillar/foundation spam.
| version = [[245.986]]
}}
{{Server config variable
| name = OverrideOfficialDifficulty
| type = float
| default = 0.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Allows you to override the default server difficulty level of 4 with 5 to match the new official server difficulty level. Default value of 0.0 disables the override. A value of 5.0 will allow common creatures to spawn up to level 150. Originally ([[247.95]]) available only as command line option.
| version = [[252.5]]
}}
{{Server config variable
| name = OverrideStructurePlatformPrevention
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, [[Auto Turret|turrets]] becomes be buildable and functional on [[:Category:Platform Saddles|platform saddles]]. Since [[247.999]] applies on [[Wooden Spike Wall|spike]] structure too. '''Note:''' despite patch notes, in ''ShooterGameServer'' it's coded ''OverrideStructurePlatformPrevention'' with two ''r''.
| version = [[242.0]]
}}
{{Server config variable
| name = OxygenSwimSpeedStatMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Use this to set how swim speed is multiplied by level spent in oxygen. The value was reduced by 80% in [[256.0]].
| version = [[256.3]]
}}
{{Server config variable
| name = PerPlatformMaxStructuresMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Higher value increases (from a percentage scale) max number of items place-able on saddles and rafts.
| version = [[211.0]]
}}
{{Server config variable
| name = PersonalTamedDinosSaddleStructureCost
| type = integer
| default = 0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Determines the amount of "tame creature slots" a platform saddle (with structures) will use towards the tribe tame creature limit.
| version = [[265.0]]
}}
{{Server config variable
| name = PlatformSaddleBuildAreaBoundsMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Increasing the number allows structures being placed further away from the platform.
| version = [[295.102]]
}}
{{Server config variable
| name = PlayerCharacterFoodDrainMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for player characters' food consumption. Higher values increase food consumption (player characters get hungry faster).
| version = [[179.0]]
}}
{{Server config variable
| name = PlayerCharacterHealthRecoveryMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for player characters' health recovery. Higher values increase the recovery rate (player characters heal faster).
| version = [[179.0]]
}}
{{Server config variable
| name = PlayerCharacterStaminaDrainMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for player characters' stamina consumption. Higher values increase stamina consumption (player characters get tired faster).
| version = [[179.0]]
}}
{{Server config variable
| name = PlayerCharacterWaterDrainMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for player characters' water consumption. Higher values increase water consumption (player characters get thirsty faster).
| version = [[179.0]]
}}
{{Server config variable
| name = PlayerDamageMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the damage players deal with their attacks. The default value <code>1</code> provides normal damage. Higher values increase damage. Lower values decrease it.
| version = [[179.0]]
}}
{{Server config variable
| name = PlayerResistanceMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the resistance to damage players receive when attacked. The default value <code>1</code> provides normal damage. Higher values decrease resistance, increasing damage per attack. Lower values increase it, reducing damage per attack. A value of 0.5 results in a player taking half damage while a value of 2.0 would result in taking double normal damage.
| version = [[179.0]]
}}
{{Server config variable
| name = PreventDiseases
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, completely diseases on the server. Thus far just [[Swamp Fever]].
| version = [[242.3]]
}}
{{Server config variable
| name = PreventMateBoost
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, disables creature mate boosting.
| version = [[247.0]]
}}
{{Server config variable
| name = PreventOfflinePvP
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, enables the Offline Raiding Prevention (ORP). When all tribe members are logged off, tribe characters, creature and structures become invulnerable. Creature starvation still applies, moreover, characters and creature can still die if drowned. Despite the name, it works on both [[PvE]] and [[PvP]] game modes. Due to performance reason, it is recommended to set a minimum interval with <code>PreventOfflinePvPInterval</code> option before ORP becomes active. ORP also helps lowering memory and CPU usage on a server. Enabled by default on Official [[PvE]] since [[258.3]]
| version = [[242.0]]
}}
{{Server config variable
| name = PreventOfflinePvPInterval
| type = float
| default = 0.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Seconds to wait before a ORP becomes active for tribe/players and relative creatures/structures (10 seconds in official PvE servers). '''Note:''' although at code level it is defined as a floating-point number, it is suggested to use an integer instead.
| version = [[242.0]]
}}
{{Server config variable
| name = PreventSpawnAnimations
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, player characters (re)spawn without the wake-up animation.
| version = [[260.0]]
}}
{{Server config variable
| name = PreventTribeAlliances
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents tribes from creating Alliances.
| version = [[243.0]]
}}
{{Server config variable
| name = ProximityChat
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, only players near each other can see their chat messages
}}
{{Server config variable
| name = PvEAllowStructuresAtSupplyDrops
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, allows building near supply drop points in [[PvE]] mode.
| version = [[247.999]]
}}
{{Server config variable
| name = PvEDinoDecayPeriodMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Creature [[PvE]] auto-decay time multiplier. Requires <code>DisableDinoDecayPvE{{=}}false</code> in GameUserSettings.ini or <code>?DisableDinoDecayPvE{{=}}false</code> in command line to work.
| version = [[206.0]]
}}
{{Server config variable
| name = PvEStructureDecayPeriodMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Specifies the scaling factor for [[Building#Auto-decay| structures decay times]], e.g.: setting it at 2.0 will double all structure decay times, while setting at 0.5 will halve the timers. '''Note:''' despite the name, works in both [[PvP]] and [[PvE]] modes when structure decay is enabled.
}}
{{Server config variable
| name = PvPDinoDecay
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, enables creatures' decay in [[PvP]] while the Offline Raid Prevention is active.
| version = [[242.0]]
}}
{{Server config variable
| name = PvPStructureDecay
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, enables structures decay on [[PvP]] servers while the Offline Raid Prevention is active.
| version = [[206.0]]
}}
{{Server config variable
| name = RaidDinoCharacterFoodDrainMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Affects how quickly the food drains on such "raid dinos" (e.g.: Titanosaurus)
| version = [[243.0]]
}}
{{Server config variable
| name = RandomSupplyCratePoints
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, supply drops are in random locations. '''Note:''' This setting is known to cause artifacts becoming inaccessible on [[Ragnarok]] if active.
| version = [[278.0]]
}}
{{Server config variable
| name = RCONEnabled
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, enables RCON, needs <code>RCONPort{{=}}<TCP_PORT></code> and <code>ServerAdminPassword{{=}}<admin_password></code> to work.
| version = [[185.0]]
}}
{{Server config variable
| name = RCONPort
| type = integer
| default = 27020
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the optional TCP RCON Port. See [[Dedicated_server_setup#Network|Dedicated server setup]]
| version = [[185.0]]
}}
{{Server config variable
| name = RCONServerGameLogBuffer
| type = float
| default = 600.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Determines how many lines of game logs are send over the RCON. '''Note:''' despite being coded as a float it's suggested to treat it as integer.
| version = [[224.0]]
}}
{{Server config variable
| name = ResourcesRespawnPeriodMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the re-spawn rate for resource nodes (trees, rocks, bushes, etc.). Lower values cause nodes to re-spawn more frequently.
| version = [[181.0]]
}}
{{Server config variable
| name = ServerAdminPassword
| type = string
| default = 
| cli = Yes
| inASA = Yes
| inASE = Yes
| If specified, players must provide this password (via the in-game console) to gain access to administrator commands on the server. '''Note:''' no quotes are used.
}}
{{Server config variable
| name = ServerAutoForceRespawnWildDinosInterval
| type = float
| default = 0.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Force re-spawn of all wild creatures on server restart afters the value set in seconds. Default value of 0.0 disables it. Useful to prevent certain creature species (like the Basilo and Spino) from becoming depopulated on long running servers.
| version = [[265.0]]
}}
{{Server config variable
| name = ServerCrosshair
| type = boolean
| default = True
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>False</code>, disables the Crosshair on your server.
| version = [[245.0]]
}}
{{Server config variable
| name = ServerForceNoHUD
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, HUD is always disabled for non-tribe owned NPCs.
}}
{{Server config variable
| name = ServerHardcore
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, enables [[Hardcore]] mode (player characters revert to level 1 upon death)
}}
{{Server config variable
| name = ServerPassword
| type = string
| default = 
| cli = Yes
| inASA = Yes
| inASE = Yes
| If specified, players must provide this password to join the server. '''Note:''' no quotes are used.
}}
{{Server config variable
| name = serverPVE
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, disables [[PvP]] and enables [[PvE]]
}}
{{Server config variable
| name = ShowFloatingDamageText
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, enables RPG-style popup damage text mode.
| version = [[242.0]]
}}
{{Server config variable
| name = ShowMapPlayerLocation
| type = boolean
| default = True
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>False</code>, hides each player their own precise position when they view their map.
}}
{{Server config variable
| name = SpectatorPassword
| type = string
| default = 
| cli = Yes
| inASA = Unknown
| inASE = Yes
| To use non-admin spectator, the server must specify a spectator password. Then any client can use these console commands: <code>requestspectator &lt;password&gt;</code> and <code>stopspectating</code>. '''Note:''' no quotes are used.
| version = [[191.0]]
}}
{{Server config variable
| name = StructureDamageMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Specifies the scaling factor for the damage structures deal with their attacks (i.e., spiked walls). Higher values increase damage. Lower values decrease it.
| version = [[179.0]]
}}
{{Server config variable
| name = StructurePickupHoldDuration
| type = float
| default = 0.5
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the quick pick-up hold duration, a value of <code>0</code> results in instant pick-up.
| version = [[293.100]]
}}
{{Server config variable
| name = StructurePickupTimeAfterPlacement
| type = float
| default = 30.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Amount of time in seconds after placement that quick pick-up is available.
| version = [[293.100]]
}}
{{Server config variable
| name = StructurePreventResourceRadiusMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Same as <code>ResourceNoReplenishRadiusStructures</code> in ''Game.ini''. If both settings are set both multiplier will be applied. Can be useful when cannot change the ''Game.ini'' file as it works as a command line option too.
}}
{{Server config variable
| name = StructureResistanceMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the resistance to damage structures receive when attacked. The default value <code>1</code> provides normal damage. Higher values decrease resistance, increasing damage per attack. Lower values increase it, reducing damage per attack. A value of 0.5 results in a structure taking half damage while a value of 2.0 would result in a structure taking double normal damage.
| version = [[179.0]]
}}
{{Server config variable
| name = TamedDinoDamageMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Specifies the scaling factor for the damage tame creatures deal with their attacks. The default value <code>1</code> provides normal damage. Higher values increase damage. Lower values decrease it.
| version = [[179.0]]
}}
{{Server config variable
| name = TamedDinoResistanceMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Specifies the scaling factor for the resistance to damage tame creatures receive when attacked. The default value <code>1</code> provides normal damage. Higher values decrease resistance, increasing damage per attack. Lower values increase it, reducing damage per attack. A value of 0.5 results in a structure taking half damage while a value of 2.0 would result in a structure taking double normal damage.
| version = [[179.0]]
}}
{{Server config variable
| name = TamingSpeedMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for creature taming speed. Higher values make taming faster.
| version = [[179.0]]
}}
{{Server config variable
| name = TheMaxStructuresInRange
| type = integer
| default = 10500
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the maximum number of structures that can be constructed within a certain (currently hard-coded) range. Replaces the old value <code>NewMaxStructuresInRange</code>''
| version = [[252.1]]
}}
{{Server config variable
| name = TribeLogDestroyedEnemyStructures
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| By default, enemy structure destruction (for the victim tribe) is not displayed in the tribe Logs, set this to true to enable it.
| version = [[247.93]]
}}
{{Server config variable
| name = TribeNameChangeCooldown
| type = float
| default = 15.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Cool-down, in minutes, in between tribe name changes. Official server use a value of 172800.0 (2 days).
| version = [[278.0]]
}}
{{Server config variable
| name = UseFjordurTraversalBuff
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, enables the biome teleport in [[Fjordur]] when holding {{Key|R}} (enabled in official PvE servers).
| version = [[346.12]]
}}
{{Server config variable
| name = UseOptimizedHarvestingHealth
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, enables a server harvesting optimization with high <code>HarvestAmountMultiplier</code> (but less rare items). '''Note:''' on {{ItemLink|ARK: Survival Evolved}} it's suggested to enable this option if harvesting with Tek Stryder causes lag spikes.
| version = [[254.93]]
}}
{{Server config variable
| name = XPMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Specifies the scaling factor for the experience received by players, tribes and tames for various actions. The default value <code>1</code> provides the same amounts of experience as in the single player experience (and official public servers). Higher values increase XP amounts awarded for various actions; lower values decrease it. In [[313.5]] an additional hardcoded multiplier of 4 was activated.
| version = [[179.0]]
}}
{{Server config variable
| name = YoungIceFoxDeathCooldown
| type = float
| default = 3600
| cli = Yes
| inASA = Yes
| inASE = No
| Overrides the cooldown for Veilwyn to reappear after taking fatal damage (in seconds), default is set to 1 hour. Must be greater than 0.
}}
|-
! colspan="6" | CrossARK Transfers
|-
{{Server config variable
| name = CrossARKAllowForeignDinoDownloads
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, enables non-native creatures tribute download on [[Aberration]].
| version = [[275.0]]
}}
{{Server config variable
| name = MinimumDinoReuploadInterval
| type = float
| default = 0.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Number of seconds cool-down between allowed creature re-uploads (43200 on official Servers which is 12 hours).
| version = [[251.0]]
}}
{{Server config variable
| name = noTributeDownloads
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents CrossArk-data downloads in[[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]].
| version = [[246.0]]
}}
{{Server config variable
| name = PreventDownloadDinos
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents creatures download from ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]].
| version = [[246.0]]
}}
{{Server config variable
| name = PreventDownloadItems
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents items download from ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]].
| version = [[246.0]]
}}
{{Server config variable
| name = PreventDownloadSurvivors
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents survivors download from ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]].
| version = [[246.0]]
}}
{{Server config variable
| name = PreventUploadDinos
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents creatures upload to ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]].
| version = [[246.0]]
}}
{{Server config variable
| name = PreventUploadItems
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents items upload to ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]].
| version = [[246.0]]
}}
{{Server config variable
| name = PreventUploadSurvivors
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = Yes
| If <code>True</code>, prevents survivors upload to ARK Data in [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]].
| version = [[246.0]]
}}
{{Server config variable
| name = TributeCharacterExpirationSeconds
| type = integer
| default = 0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Set in seconds the expiration timer for uploaded survivors in ARK Data. With default or negative values there is no expiration time. Check [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]] for more details. '''Warning''': do not set this option to an insane high value, like more than 31536000 seconds (which is 1 year), as in ARK Data routines this is summed to upload time in Unix Epoch time format. Using really high values will result in overflow and may cause upload time checks to fail and ARK Data deleted. Finally, it is highly suggested to use the same value across all servers in a cluster, otherwise accessing ARK Data from a server with a lower value will result data deletion if in that server the timer is expired.
}}
{{Server config variable
| name = TributeDinoExpirationSeconds
| type = integer
| default = 86400
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Set in seconds the expiration timer for uploaded tames in ARK Data. If set to 0 or less will revert to default. Check [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]] for more details. '''Warning''': do not set this option to an insane high value, like more than 31536000 seconds (which is 1 year), as in ARK Data routines this is summed to upload time in Unix Epoch time format. Using really high values will result in overflow and may cause upload time checks to fail and ARK Data deleted. Finally, it is highly suggested to use the same value across all servers in a cluster, otherwise accessing ARK Data from a server with a lower value will result data deletion if in that server the timer is expired.
}}
{{Server config variable
| name = TributeItemExpirationSeconds
| type = integer
| default = 86400
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Set in seconds the expiration timer for uploaded items in ARK Data. If set to 0 or less will revert to default. Check [[#Cross-ARK Data Transfer|Cross-ARK Data Transfer]] for more details. '''Warning''': do not set this option to an insane high value, like more than 31536000 seconds (which is 1 year), as in ARK Data routines this is summed to upload time in Unix Epoch time format. Using really high values will result in overflow and may cause upload time checks to fail and ARK Data deleted. Finally, it is highly suggested to use the same value across all servers in a cluster, otherwise accessing ARK Data from a server with a lower value will result data deletion if in that server the timer is expired.
}}
|-
! colspan="6" | Cryo Sickness and Cryopod Nerf
|-
{{Server config variable
| name = CryopodNerfDamageMult
| type = float
| default = 0.0099999998
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Reduces the amount of damage dealt by the creature after it is deployed from the cryopod, as a percentage of total damage output, and for the length of time set by <code>CryopodNerfDuration</code>. <code>CryopodNerfDuration</code> needs to be set as well. <code>CryopodNerfDamageMult{{=}}0.01</code> means 99% of the damage is removed.  On official it is set to 0.1.
| version = [[309.53]]
}}
{{Server config variable
| name = CryopodNerfDuration
| type = integer
| default = 0.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Amount of time, in seconds, Cryosickness lasts after deploying a creature from a Cryopod. '''Note:''' although at code level it is defined as a floating-point number, it is suggested to use an integer instead. On official it is set to 10.0.
| version = [[309.53]]
}}
{{Server config variable
| name = CryopodNerfIncomingDamageMultPercent
| type = float
| default = 0.0
| cli = No
| inASA = Unknown
| inASE = Yes
| Increases the amount of damage taken by the creature after it is deployed from the cryopod, as a percentage of total damage received, and for the length of time set by <code>CryopodNerfDuration</code>. <code>CryopodNerfIncomingDamageMultPercent{{=}}0.25</code> means a released tame takes 25% more damage while the debuff lasts. On official it is set to 0.25.
| version = [[310.11]]
}}
{{Server config variable
| name = EnableCryopodNerf
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, there is no Cryopod cooldown timer, and creatures do not become unconscious. If this option is set, than <code>EnableCryopodNerf</code> and <code>CryopodNerfIncomingDamageMultPercent</code> must be set as well or they will default to 0.
| version = [[309.53]]
}}
{{Server config variable
| name = EnableCryoSicknessPVE
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, enables Cryopod cooldown timer when deploying a creature.
| version = [[307.56]]
}}
|-
! colspan="6" | Text Filtering
|-
{{Server config variable
| name = BadWordListURL
| type = string with a URL
| default = {{ItemLink|ARK: Survival Evolved|&nbsp;}}: "http://arkdedicated.com/badwords.txt"<br/>{{ItemLink|ARK: Survival Ascended|&nbsp;}}: "http://cdn2.arkdedicated.com/asa/badwords.txt"
| cli = No
| inASA = Yes
| inASE = Yes
| Add the <code>URL</code> to hosting your own bad words list. '''Note:''' on {{Gamelink|ase}} servers only the HTTP protocol is supported (an HTTPS URL will not work).
| version = [[356.3]]
}}
{{Server config variable
| name = BadWordWhiteListURL
| type = string with a URL
| default = {{ItemLink|ARK: Survival Evolved|&nbsp;}}: "http://arkdedicated.com/goodwords.txt"<br/>{{ItemLink|ARK: Survival Ascended|&nbsp;}}: "http://cdn2.arkdedicated.com/asa/goodwords.txt"
| cli = No
| inASA = Yes
| inASE = Yes
| Add the <code>URL</code> to hosting your own good words list. '''Note:''' on {{Gamelink|ase}} servers only the HTTP protocol is supported (an HTTPS URL will not work).
| version = [[356.3]]
}}
{{Server config variable
| name = bFilterCharacterNames
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, filters out character names based on the bad words/good words list.
| version = [[356.3]]
}}
{{Server config variable
| name = bFilterChat
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, filters out character names based on the bad word/good words list.
| version = [[356.3]]
}}
{{Server config variable
| name = bFilterTribeNames
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, filters out tribe names based on the badwords/goodwords list.
| version = [[356.3]]
}}
|-
! colspan="6" | Deprecated options
|-
{{Server config variable
| status = deprecated
| name = AllowDeprecatedStructures
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, allows servers to keep the Halloween Structures for a while after event ends (had to be used before relaunching the server with [[222.3]] update). Since no more events are planned to be activated or removed, this feature can be considered deprecated.
| version = [[222.3]]
}}
{{Server config variable
| status = deprecated
| name = bAllowFlyerCarryPVE
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| Use <code>AllowFlyerCarryPvE</code>.
| version = [[173.0]]
}}
{{Server config variable
| status = deprecated
| name = bDisableStructureDecayPvE
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| Use <code>DisableStructureDecayPvE</code>.
| version = [[173.0]]
}}
{{Server config variable
| status = deprecated
| name = ForceFlyerExplosives
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, allowed flyers (except [[Quetzal]] and [[Wyvern]]) to fly with [[C4]] attached to it. Deprecated since [[253.95]].
| version = [[252.83]]
}}
{{Server config variable
| status = deprecated
| name = MaxStructuresInRange
| type = integer
| default = 1300
| cli = No
| inASA = Unknown
| inASE = Yes
| Specifies the maximum number of structures that can be constructed within a certain (currently hard-coded) range. Deprecated with patch [[188.0]] by <code>NewMaxStructuresInRange</code>.
| version = [[173.0]]
}}
{{Server config variable
| status = deprecated
| name = NewMaxStructuresInRange
| type = integer
| default = 6000
| cli = No
| inASA = Unknown
| inASE = Yes
| Specifies the maximum number of structures that can be constructed within a certain (currently hard-coded) range. Deprecated with patch [[252.1]] by <code>TheMaxStructuresInRange</code>.
| version = [[188.0]]
}}
{{Server config variable
| status = deprecated
| name = PvEStructureDecayDestructionPeriod
| type = integer
| default = 0
| cli = No
| inASA = Unknown
| inASE = Yes
| Specified the time required for player structures to decay in [[PvE]] mode. Deprecated and no more present in executable bits with patch [[180.0]] where each type of structure has its own decay time, increasing with "tier quality".<ref name="extended_game_options">Drake (23 June 2015). [http://steamcommunity.com/app/346110/discussions/10/530646715633129364/ "Extended Game Options: How to Configure your Custom ARK server"]. ''Server Hosting & Advertisement''. Steam Community :: ARK: Survival Evolved. Retrieved 19 July 2015.</ref>
| version = [[179.0]]
}}
|-
! colspan="6" | Undocumented options
|-
{{Server config variable
| name = AdminListURL
| type = string with a URL
| default = N/A
| cli = Yes
| inASA = Yes
| inASE = No
| version = {{ItemLink|ARK: Survival Ascended|&nbsp;}}
| info = Alternative to ''AllowedCheaterAccountIDs.txt'' (see [[#Administrator_Whitelisting|Administrator Whitelisting]]) using a web resource. The interval at which the server queries the resource to check for admin list update is defined by <code>UpdateAllowedCheatersInterval</code>. Undocumented by Wildcard.
}}
{{Server config variable
| name = AllowedCheatersURL
| type = string with a URL
| default = N/A
| cli = Yes
| inASA = No
| inASE = Yes
| Alternative to ''AllowedCheaterSteamIDs.txt'' (see [[#Administrator_Whitelisting|Administrator Whitelisting]]) using a web resource (e.g.: official uses http://arkdedicated.com/globaladmins.txt). The interval at which the server queries the resource to check for admin list update is defined by <code>UpdateAllowedCheatersInterval</code> '''Note:''' it supports the HTTP protocol only (HTTPS is not supported). Undocumented by Wildcard.
}}
{{Server config variable
| name = AutoRestartIntervalSeconds
| type = float
| default = Unknown
| cli = Yes
| inASA = Yes
| inASE = Unknown
| Defines the time (in seconds) after which the server will automatically restart. Undocumented by Wildcard.
(Appears to shut off the server instead of restarting properly)}}
{{Server config variable
| name = ChatLogFileSplitIntervalSeconds
| type = integer
| default = 86400
| cli = No
| inASA = Unknown
| inASE = Yes
| Controls how to split the chat log file related to time in seconds. Cannot be set to a value lower than 45 (will default to 45 if the value is lower). Set to 0 only in official. Undocumented by Wildcard.
}}
{{Server config variable
| name = ChatLogFlushIntervalSeconds
| type = integer
| default = 86400
| cli = No
| inASA = Unknown
| inASE = Yes
| Controls in how many second the chat log is flushed to log file. Cannot be set to a value lower than 15 (will default to 15 if the value is lower). Set to 0 only in official. Undocumented by Wildcard.
}}
{{Server config variable
| name = ChatLogMaxAgeInDays
| type = integer
| default = 5
| cli = No
| inASA = Unknown
| inASE = Yes
| Controls how many days the chat log is long. Set it to a negative value will result it to set at -1 (virtually infinite). Set to 0 only in official. Undocumented by Wildcard.
}}
{{Server config variable
| name = DontRestoreBackup
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code> and <code>-DisableDupeLogDeletes</code> is present, prevents the server to automatically restore a backup in case of corrupted save. Undocumented by Wildcard.
}}
{{Server config variable
| name = EnableAFKKickPlayerCountPercent
| type = float
| default = 0.0
| cli = No
| inASA = Unknown
| inASE = Yes
| Enables the idle timeout to be applied only if the amount of online players reaches percentage value related to <code>MaxPlayers</code> argument. The percentage is expressed as normalised value between 0 and 1.0, where 1.0 means 100%. Official set it to 0.89999998. Undocumented by Wildcard.
}}
{{Server config variable
| name = EnableMeshBitingProtection
| type = boolean
| default = True
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>False</code>, disables mesh biting protection. Undocumented by Wildcard.
}}
{{Server config variable
| name = FreezeReaperPregnancy
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, freezes the {{ItemLink|Reaper King}} pregnancy timer and [[experience]] gain. Undocumented by Wildcard.
}}
{{Server config variable
| name = LogChatMessages
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, enables advanced chat logging. Chat logs will be saved in <code>ShooterGame/Saved/Logs/ChatLogs/<SessionName>/</code> in json format. Disabled on official. The file will be split according to <code>ChatLogFileSplitIntervalSeconds</code> value and flushed every <code>ChatLogFlushIntervalSeconds</code> seconds value. '''Note:''' <code><SessionName></code> will be written without spaces. Undocumented by Wildcard.
}}
{{Server config variable
| name = MaxStructuresInSmallRadius
| type = integer
| default = 0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Defines the amount of max structures allowed to be placed in a <code>RadiusStructuresInSmallRadius</code> from player position. Official set it to 40. Undocumented by Wildcard.
}}
{{Server config variable
| name = MaxStructuresToProcess
| type = integer
| default = 0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Controls the max batch size of structures to process (e.g.: culling, building graphs modifications, etc) at each server tick. Leaving at 0 (default behaviour) will force the server to process all structures in queue. It's a trade-off between how much can cost a server tick (worst case scenario) and simulation accuracy. Official set it to 5000. Undocumented by Wildcard.
}}
{{Server config variable
| name = PreventOutOfTribePinCodeUse
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, prevents out of tribe players to use pins on structures (doors, elevators, storage boxes, etc). Undocumented by Wildcard.
}}
{{Server config variable
| name = RadiusStructuresInSmallRadius
| type = float
| default = 0.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Defines the small radius dimension (in Unreal Units) used by <code>MaxStructuresInSmallRadius</code>. Official set it to 225.0. Undocumented by Wildcard.
}}
{{Server config variable
| name = ServerEnableMeshChecking
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Involved in foliage repopulation. Takes no effect if <code>-forcedisablemeshchecking</code> is set. Enabled on official. Undocumented.
}}
{{Server config variable
| name = TribeMergeAllowed
| type = boolean
| default = True
| cli = No
| inASA = Unknown
| inASE = Yes
| If <code>False</code>, prevents tribe to merge. Undocumented by Wildcard.
}}
{{Server config variable
| name = TribeMergeCooldown
| type = float
| default = 0.0
| cli = No
| inASA = Unknown
| inASE = Yes
| Tribe merge cool-down in seconds. Official uses 86400.0. Undocumented by Wildcard.
}}
{{Server config variable
| name = UpdateAllowedCheatersInterval
| type = float
| default = 600.0
| cli = Yes
| inASA = Yes
| inASE = Yes
| Times in seconds at which the remote admin list linked by <code>AllowedCheatersURL</code> is queried for updates. Any value less than 3.0 will be reverted to 3.0. Undocumented by Wildcard.
}}
{{Server config variable
| name = UseCharacterTracker
| type = boolean
| default = False
| cli = Yes
| inASA = Yes
| inASE = No
| Used to enable character tracking. Alternatively, this option can be configured with <code>-disableCharacterTracker</code> argument in the command line (note that the argument from command line has priority over the value set in [[#GameUserSettings.ini|GameUserSettings.ini]]). Undocumented by Wildcard.
}}
{{Server config variable
| name = UseExclusiveList
| type = boolean
| default = False
| cli = Yes
| inASA = Unknown
| inASE = Yes
| If <code>True</code>, allows same behaviour as <code>-exclusivejoin</code>. Undocumented by Wildcard.
}}
|-
! colspan="6" | Clients only options (ineffective on dedicated servers)
|-
{{Server config variable
| name = ListenServerTetherDistanceMultiplier
| type = float
| default = 1.0
| cli = Yes
| inASA = Unknown
| inASE = Yes
| Scales the tether distance between host and other players on non-dedicated sessions only. '''Note:''' despite being readable from command line, this option affects non-dedicated sessions only, thus it has to be set in the ''GameUserSettings.ini'' file or through the in-game local host graphics menu.
}}
|}
</div>
<!--end of [ServerSettings]-->
<!--start of [SessionSettings]-->

==== [SessionSettings] ====
The following options must come under the <code>[SessionSettings]</code> section of <code>GameUserSettings.ini</code>.
<div class="widetable">
{| class="wikitable config-table"
|-
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Variable
! Description
! {{HoverNote|CMD|If ticked, this option may be specified in the server's command line instead, after map name.}}
! Since patch
{{Server config variable
| name = MultiHome
| type = IP_ADDRESS
| default = N/A
| cli = Yes
| inASA = Unknown
| inASE = Yes
| info = Specifies MultiHome IP Address. Boolean <code>Multihome</code> option must be set to <code>True</code> as well (command line or <code>[MultiHome]</code> section). Leave it empty if not using multihoming. Can be specified in command line too.
}}
{{Server config variable
| name = Port
| type = integer
| default = 7777
| cli = Yes
| inASA = Yes
| inASE = Yes
| info = Specifies the ''UDP Game Port''. See [[Dedicated_server_setup#Network|Dedicated server setup]]<br/>'''Note''': command line append syntax is not supported by {{gamelink|sa}}
}}
{{Server config variable
| name = QueryPort
| type = integer
| default = 27015
| cli = Yes
| inASA = Unknown
| inASE = Yes
| info = Specifies the ''UDP Steam Query Port''. See [[Dedicated_server_setup#Network|Dedicated server setup]]
}}
{{Server config variable
| name = SessionName
| type = string
| default = ARK #123456
| cli = Yes
| inASA = Yes
| inASE = Yes
| info = Specifies the Server name advertised in the Game Server Browser as well in Steam Server browser. If no name is provide, the default name will be ''ARK #'' followed by a random 6 digit number. '''Note:''' Name must '''not''' be typed between quotes unless it is launched from command line.
}}
|}
</div>
<!--end of [SessionSettings]-->
<!--start of [MultiHome]-->

==== [MultiHome] ====
The following options must come under the <code>[/MultiHome]</code> section of <code>GameUserSettings.ini</code>.
<div class="widetable">
{| class="wikitable"
|-a
! Option !! Default !! Effect !! CMD
|-
|<code>MultiHome='''<boolean>'''</code> || False || If <code>True</code>, enables multihoming. <code>MultiHome</code> IP must be specified in <code>[SessionSettings]</code> or in command line as well. || {{Checkmark|N}}
|}
</div>
<!--start of [MultiHome]-->
<!--start of [/Script/Engine.GameSession]-->

==== [/Script/Engine.GameSession] ====
The following options must come under the <code>[/Script/Engine.GameSession]</code> section of <code>GameUserSettings.ini</code>.
<div class="widetable">
{| class="wikitable"
|-
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Variable
! Description
! {{HoverNote|CMD|If ticked, this option may be specified in the server's command line instead, after map name.}}
! Since patch
{{Server config variable
| name = MaxPlayers
| type = integer
| default = 70
| cli = Yes
| inASA = No
| inASE = Yes
| info = Specifies the maximum number of players that can play on the server simultaneously.<br/>{{ItemLink|ARK: Survival Ascended|ASA}}: This setting is replaced with <code>-WinLiveMaxPlayers</code> in the [[#options|command line options]], as otherwise, it will force it back to the default value.
}}
|}
</div>
<!--end of [/Script/Engine.GameSession]-->
<!--start of [Ragnarok]-->

==== [Ragnarok] ====
The following options must come under the <code>[Ragnarok]</code> section of <code>GameUserSettings.ini</code>.
<div class="widetable">
{| class="wikitable config-table"
|-
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Variable
! Description
! {{HoverNote|CMD|If ticked, this option may be specified in the server's command line instead, after map name.}}
! Since patch
{{Server config variable
| name = AllowMultipleTamedUnicorns
| type = boolean
| default = False
| cli = No
| inASA = Unknown
| inASE = Yes
| info = <code>False</code> = one unicorn on the map at a time, <code>True</code> = one wild and unlimited tamed Unicorns on the map.
}}
{{Server config variable
| name = EnableVolcano
| type = boolean
| default = True
| cli = No
| inASA = Unknown
| inASE = Yes
| info = <code>False</code> = disabled (the volcano will '''not''' become active), <code>True</code> = enabled
}}
{{Server config variable
| name = UnicornSpawnInterval
| type = integer
| default = 24
| cli = No
| inASA = Unknown
| inASE = Yes
| info = How long in hours the game should wait before spawning a new [[Unicorn]] if the wild one is killed (or tamed, if <code>AllowMultipleTamedUnicorns</code> is enabled). This value sets the minimum amount of time (in hours), and the maximum is equal to 2x this value.
}}
{{Server config variable
| name = VolcanoIntensity
| type = float
| default = 1
| cli = No
| inASA = Unknown
| inASE = Yes
| info = The lower the value, the more intense the volcano's eruption will be. Recommended to leave at 1. The minimum value is 0.25, and for multiplayer games, it should not go below 0.5. Very high numbers will basically disable the flaming rocks flung out of the volcano.
}}
{{Server config variable
| name = VolcanoInterval
| type = integer
| default = 0
| cli = No
| inASA = Unknown
| inASE = Yes
| info = 0 = 5000 (min) - 15000 (max) seconds between instances of the volcano becoming active. Any number above 0 acts as a multiplier, with a minimum value of .1
}}
|}
</div>
<!--end of [Ragnarok]-->
<!--start of [/Script/Engine.GameSession]-->

==== [MessageOfTheDay] ====
The following options must come under the <code>[MessageOfTheDay]</code> section of <code>GameUserSettings.ini</code>.
<div class="widetable">
{| class="wikitable"
|-
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Variable
! Description
! {{HoverNote|CMD|If ticked, this option may be specified in the server's command line instead, after map name.}}
! Since patch
{{Server config variable
| name = Duration
| type = integer
| default = 20
| cli = No
| inASA = Yes
| inASE = Yes
| info = Specifies in seconds the duration of the displayed message on player log-in.
| version = [[171.0]]
}}
{{Server config variable
| name = Message
| type = string
| default = N/A
| cli = No
| inASA = Yes
| inASE = Yes
| info = A single line string for a message displayed to played once logged-in. No quotes needed. Use <code>\n</code> to start a new line in the message.
| version = [[171.0]]
}}
|}
</div>
<!--end of [/Script/Engine.GameSession]-->
<!--end of GameUserSettings.ini-->

===Game.ini===
All the following options can only be set in the <code>[/script/shootergame.shootergamemode]</code> section of <code>Game.ini</code>, located in the same folder as <code>GameUserSettings.ini</code> (see [[#Configuration Files|Configuration Files]] for its location). Specifying them on the command line will have no effect.
<!--start of Game.ini-->
<div class="widetable">
{| class="wikitable config-table"
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Variable
! Description
! Since patch
{{Server config variable
| name = AutoPvEStartTimeSeconds
| type = float
| default = 0.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = States when the [[PvE]] mode should start in a [[PvP|PvPvE]] server. Valid values are from 0 to 86400. Options <code>bAutoPvETimer</code>, <code>bAutoPvEUseSystemTime</code> and <code>AutoPvEStopTimeSeconds</code> must also be set. '''Note:''' although at code level it is defined as a floating point number, it is suggested to use an integer instead.
| version = [[196.0]]
}}
{{Server config variable
| name = AutoPvEStopTimeSeconds
| type = float
| default = 0.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = States when the [[PvE]] mode should end in a [[PvP|PvPvE]] server. Valid values are from 0 to 86400. Options <code>bAutoPvETimer</code>, <code>bAutoPvEUseSystemTime</code> and <code>AutoPvEStopTimeSeconds</code> must also be set. '''Note:''' although at code level it is defined as a floating point number, it is suggested to use an integer instead.
| version = [[196.0]]
}}
{{Server config variable
| name = BabyCuddleGracePeriodMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales how long after delaying cuddling with the Baby before Imprinting Quality starts to decrease.
| version = [[242.0]]
}}
{{Server config variable
| name = BabyCuddleIntervalMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales how often babies needs attention for imprinting. More often means you'll need to cuddle with them more frequently to gain Imprinting Quality. Scales always according to default <code>BabyMatureSpeedMultiplier</code> value: set at 1.0 the imprint request is every 8 hours. See also [[Imprinting#The_Imprinting_formula|The Imprinting formula]] how it affects the imprinting amount at each baby care/cuddle.
| version = [[219.0]]
| version = [[242.0]]
}}
{{Server config variable
| name = BabyCuddleLoseImprintQualitySpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales how fast Imprinting Quality decreases after the grace period if you haven't yet cuddled with the Baby.
| version = [[242.0]]
}}
{{Server config variable
| name = BabyFoodConsumptionSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the speed that baby tames eat their food. A lower value decreases (by percentage) the food eaten by babies.
| version = [[222.3]]
}}
{{Server config variable
| name = BabyImprintAmountMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the percentage each imprint provides. A higher value, will rise the amount of imprinting % at each baby care/cuddle, a lower value will decrease it. This multiplier is global, meaning it will affect the imprinting progression of every species. See also [[Imprinting#The_Imprinting_formula|The Imprinting formula]] how it affects the imprinting amount at each baby care/cuddle.
| version = [[312.35]]
}}
{{Server config variable
| name = BabyImprintingStatScaleMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales how much of an effect on stats the Imprinting Quality has. Set it to 0 to effectively disable the system.
| version = [[242.0]]
}}
{{Server config variable
| name = BabyMatureSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the maturation speed of babies. A higher number decreases (by percentage) time needed for baby tames to mature. See [[Breeding#Times_for_Breeding|Times for Breeding]] tables for values at 1.0, see [[Imprinting#The_Imprinting_formula|The Imprinting formula]] how it affects the imprinting amount at each baby care/cuddle.
| version = [[219.0]]
}}
{{Server config variable
| name = bAllowUnclaimDinos
| type = boolean
| default = True
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>False</code>, prevents players to unclaim tame creatures.
}}
{{Server config variable
| name = bAllowCustomRecipes
| type = boolean
| default = True
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>False</code>, disabled custom RP-oriented Recipe/Cooking System (including Skill-Based results).
| version = [[224.0]]
}}
{{Server config variable
| name = bAllowFlyerSpeedLeveling
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Specifies whether flyer creatures can have their [[Movement Speed]] levelled up.
| version = [[321.1]]
}}
{{Server config variable
| name = bAllowPlatformSaddleMultiFloors
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, allows multiple platform floors.
| version = [[260.0]]
}}
{{Server config variable
| name = bAllowUnlimitedRespecs
| type = boolean
| default = False
| cli = skip
| inASA = Yes
| inASE = Yes
| info = If <code>True</code>, allows more than one usage of [[Mindwipe Tonic]] without 24 hours cooldown.
| version = [[260.0]]
}}
{{Server config variable
| name = BaseTemperatureMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Specifies the map base temperature scaling factor: lower value makes the environment colder, higher value makes the environment hotter.
| version = [[273.5]]
}}
{{Server config variable
| name = bAutoPvETimer
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, enabled [[PvE]] mode in a [[PvP|PvPvE]] server at pre-specified times. The option <code>bAutoPvEUseSystemTime</code> determinates what kind of time to use, while <code>AutoPvEStartTimeSeconds</code> and <code>AutoPvEStopTimeSeconds</code> set the begin and end time of PvE mode.
| version = [[196.0]]
}}
{{Server config variable
| name = bAutoPvEUseSystemTime
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, [[PvE]] mode begin and end times in a [[PvP|PvPvE]] server will refer to the server system time instead of in-game world time. Options <code>bAutoPvETimer</code>, <code>AutoPvEStartTimeSeconds</code> and <code>AutoPvEStopTimeSeconds</code> must also be set.
| version = [[196.0]]
}}
{{Server config variable
| name = bAutoUnlockAllEngrams
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, unlocks all Engrams available. Ignores OverrideEngramEntries and OverrideNamedEngramEntries entries.
| version = [[273.62]]
}}
{{Server config variable
| name = bDisableDinoBreeding
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, prevents tames to be bred.
}}
{{Server config variable
| name = bDisableDinoRiding
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, prevents players to ride tames.
| version = [[254.5]]
}}
{{Server config variable
| name = bDisableDinoTaming
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, prevents players to tame wild creatures.
| version = [[254.5]]
}}
{{Server config variable
| name = bDisableFriendlyFire
| type = boolean
| default = False
| cli = skip
| inASA = Yes
| inASE = Yes
| info = If <code>True</code>, prevents Friendly-Fire (among tribe mates/tames/structures).
| version = [[228.4]]
}}
{{Server config variable
| name = bDisableLootCrates
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, prevents spawning of Loot crates (artifact creates will still spawn).
| version = [[231.7]]
}}
{{Server config variable
| name = bDisablePhotoMode
| type = boolean
| default = False
| cli = skip
| inASA = Yes
| inASE = No
| info = Defines if photo mode is allowed (<code>False</code>) or not (<code>True</code>).
}}
{{Server config variable
| name = bDisableStructurePlacementCollision
| type = boolean
| default = False
| cli = skip
| inASA = Yes
| inASE = Yes
| info = If <code>True</code>, allows for structures to clip into the terrain.
| version = [[273.62]]
}}
{{Server config variable
| name = bFlyerPlatformAllowUnalignedDinoBasing
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, Quetz platforms will allow any non-allied tame to base on them when they are flying.
| version = [[218.3]]
}}
{{Server config variable
| name = bIgnoreStructuresPreventionVolumes
| type = boolean
| default = False
| cli = skip
| inASA = Yes
| inASE = Yes
| info = If <code>True</code>, enables building areas where normally it's not allowed, such around some maps' [[Obelisk|Obelisks]], in the [[The_Ancient_Device_(Aberration)|Aberration Portal]] and in Mission Volumes areas on [[Genesis: Part 1]]. '''Note:''' in Genesis: Part 1 this settings is enabled by default and there is an ad hoc settings called <code>bGenesisUseStructuresPreventionVolumes</code> to disable it.
}}
{{Server config variable
| name = bIncreasePvPRespawnInterval
| type = boolean
| default = True
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>False</code>, disables [[PvP]] additional re-spawn time (<code>IncreasePvPRespawnIntervalBaseAmount</code>) that scales (<code>IncreasePvPRespawnIntervalMultiplier</code>) when a player is killed by a team within a certain amount of time (<code>IncreasePvPRespawnIntervalCheckPeriod</code>).
| version = [[196.0]]
}}
{{Server config variable
| name = bOnlyAllowSpecifiedEngrams
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, any Engram not explicitly specified by <code>OverrideEngramEntries</code> or <code>OverrideNamedEngramEntries</code> list will be hidden. All Items and Blueprints based on hidden Engrams will be removed.
| version = [[187.0]]
}}
{{Server config variable
| name = bPassiveDefensesDamageRiderlessDinos
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, allows spike walls to damage wild/riderless creatures.
| version = [[224.0]]
}}
{{Server config variable
| name = bPvEAllowTribeWar
| type = boolean
| default = True
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>False</code>, disables capability for Tribes to officially declare war on each other for mutually-agreed-upon period of time.
| version = [[223.0]]
}}
{{Server config variable
| name = bPvEAllowTribeWarCancel
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, allows cancellation of an agreed-upon war before it has actually started.
| version = [[223.0]]
}}
{{Server config variable
| name = bPvEDisableFriendlyFire
| type = boolean
| default = False
| cli = skip
| inASA = Yes
| inASE = Yes
| info = If <code>True</code>, disabled Friendly-Fire (among tribe mates/tames/structures) in [[PvE]] servers.
| version = [[202.0]]
}}
{{Server config variable
| name = bShowCreativeMode
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, enables creative mode.
| version = [[278.0]]
}}
{{Server config variable
| name = bUseCorpseLocator
| type = boolean
| default = True
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>False</code>, prevents survivors to see a green light beam at the location of their dead body.
| version = [[259.0]]
}}
{{Server config variable
| name = bUseDinoLevelUpAnimations
| type = boolean
| default = True
| cli = skip
| inASA = Yes
| inASE = Yes
| info = If <code>False</code>, tame creatures on level-up will not perform the related animation.
}}
{{Server config variable
| name = bUseSingleplayerSettings
| type = boolean
| default = False
| cli = skip
| inASA = Yes
| inASE = Yes
| info = If <code>True</code>, all game settings will be more balanced for an individual player experience. Useful for dedicated server with a very small amount of players. See [[#Single_Player_Settings|Single Player Settings]] section for more details.
| version = [[259.0]]
}}
{{Server config variable
| name = bUseTameLimitForStructuresOnly
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code> will make Tame Units only be applied and used for Platforms with Structures and Rafts effectively disabling Tame Units for tames without Platform Structures.
| version = [[259.0]]
}}
{{Server config variable
| name = CheatTeleportLocations=(TeleportName="<string>",TeleportLocation=(X=<float>,Y=-<float>,Z=<float>))
| type = (...)
| cli = skip
| inASA = Yes
| inASE = Unknown
| info = Creates a named teleport location that can be used with the [[Console commands#TP|TP]] command. The coordinates must be listed in Unreal units, not in-game gps coordinates. Example:
  CheatTeleportLocations=(TeleportName="Hightower",TeleportLocation=(X=467967.0,Y=-359082.0,Z=6879.0))
  With this example, executing the command <code>cheat TP Hightower</code> will teleport you to the defined coordinates, which (if the map is Ragnarok{{IconLink|ASA|size=25px}}) is the highlands lighthouse. 

}}
{{Server config variable
| name = ConfigAddNPCSpawnEntriesContainer
| type = (...)
| default = N/A
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Adds specific creatures in spawn areas. See [[#ConfigAddNPCSpawnEntriesContainer|Creature Spawn related]] section for more detail.
| version = [[248.0]]
}}
{{Server config variable
| name = ConfigOverrideItemCraftingCosts
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Overrides items crafting resource requirements. See [[#ConfigOverrideItemCraftingCosts|Item related]] section for more details.
| version = [[242.0]]
}}
{{Server config variable
| name = ConfigOverrideItemMaxQuantity
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Overrides items stack size on a per-item basis. See [[#ConfigOverrideItemMaxQuantity|Item related]] section for more details.
| version = [[292.100]]
}}
{{Server config variable
| name = ConfigOverrideNPCSpawnEntriesContainer
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Overrides specific creatures in spawn areas. See [[#ConfigOverrideNPCSpawnEntriesContainer|Creature Spawn related]] section for more details.
| version = [[248.0]]
}}
{{Server config variable
| name = ConfigOverrideSupplyCrateItems
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Overrides items contained in loot crates. See [[#ConfigOverrideSupplyCrateItems|Items related]] section for more details.
| version = [[242.0]]
}}
{{Server config variable
| name = ConfigSubtractNPCSpawnEntriesContainer
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Removes specific creatures in spawn areas. See [[#ConfigSubtractNPCSpawnEntriesContainer|Creature Spawn related]] section for more detail.
| version = [[248.0]]
}}
{{Server config variable
| name = CraftingSkillBonusMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the bonus received from upgrading the [[Crafting Skill]].
| version = [[259.32]]
}}
{{Server config variable
| name = CraftXPMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the amount of XP earned for crafting.
| version = [[243.0]]
}}
{{Server config variable
| name = CropDecaySpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the speed of crop decay in plots. A higher value decrease (by percentage) speed of crop decay in plots.
| version = [[218.0]]
}}
{{Server config variable
| name = CropGrowthSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the speed of crop growth in plots. A higher value increases (by percentage) speed of crop growth.
| version = [[218.0]]
}}
{{Server config variable
| name = CustomRecipeEffectivenessMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the effectiveness of custom recipes. A higher value increases (by percentage) their effectiveness.
| version = [[226.0]]
}}
{{Server config variable
| name = CustomRecipeSkillMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the effect of the players crafting speed level that is used as a base for the formula in creating a custom recipe. A higher number increases (by percentage) the effect.
| version = [[226.0]]
}}
{{Server config variable
| name = DestroyTamesOverLevelClamp
| type = integer
| default = 0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Tames that exceed that level will be deleted on server start. Official servers have it set to <code>450</code>.
| version = [[255.0]]
}}
{{Server config variable
| name = DinoClassDamageMultipliers
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Globally overrides wild creatures damages. See [[#DinoClassDamageMultipliers|Creature Stats related]] section for more detail.
| version = [[194.0]]
}}
{{Server config variable
| name = DinoClassResistanceMultipliers
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Globally overrides wild creatures resistance. See [[#DinoClassResistanceMultipliers|Creature Stats related]] section for more detail.
| version = [[194.0]]
}}
{{Server config variable
| name = DinoHarvestingDamageMultiplier
| type = float
| default = 3.2
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales the damage done to a harvestable item/entity by a tame. A higher number increases (by percentage) the speed of harvesting.
| version = [[321.1]]
}}
{{Server config variable
| name = DinoSpawnWeightMultipliers
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Globally overrides creatures spawns likelihood. See [[#DinoSpawnWeightMultipliers|Creature Spawn related]] section for more detail.
| version = [[179.0]]
}}
{{Server config variable
| name = DinoTurretDamageMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales the damage done by Turrets towards a creature. A higher values increases it (by percentage).
| version = [[231.4]]
}}
{{Server config variable
| name = EggHatchSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the time needed for a fertilised egg to hatch. A higher value decreases (by percentage) that time.
| version = [[219.0]]
}}
{{Server config variable
| name = EngramEntryAutoUnlocks
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Automatically unlocks the specified Engram when reaching the level specified. See [[#EngramEntryAutoUnlocks|Engram Entries related]] section for more detail.
| version = [[273.62]]
}}
{{Server config variable
| name = ExcludeItemIndices
| type = integer
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Excludes an item from supply crates specifying its [[Item IDs|Item ID]]. You can have multiple lines of this option.
}}
{{Server config variable
| name = FastDecayInterval
| type = integer
| default = 43200
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Specifies the decay period for "Fast Decay" structures (such as pillars or lone foundations). Value is in seconds. <code>FastDecayUnsnappedCoreStructures</code> in ''GameUserSettings.ini'' must be set to <code>True</code> as well to take any effect.
| version = [[259.0]]
}}
{{Server config variable
| name = FishingLootQualityMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Sets the quality of items that have a quality when fishing. Valid values are from 1.0 to 5.0.
| version = [[260.0]]
}}
{{Server config variable
| name = FuelConsumptionIntervalMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Defines the interval of fuel consumption.
| version = [[264.0]]
}}
{{Server config variable
| name = GenericXPMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the amount of XP earned for generic XP (automatic over time).
| version = [[243.0]]
}}
{{Server config variable
| name = GlobalCorpseDecompositionTimeMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales the decomposition time of corpses, (player and creature), globally. Higher values prolong the time.
| version = [[189.0]]
}}
{{Server config variable
| name = GlobalItemDecompositionTimeMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the decomposition time of dropped items, loot bags etc. globally. Higher values prolong the time.
| version = [[189.0]]
}}
{{Server config variable
| name = GlobalPoweredBatteryDurabilityDecreasePerSecond
| type = float
| default = 3.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Specifies the rate at which charge batteries are used in electrical objects.
| version = [[275.0]]
}}
{{Server config variable
| name = GlobalSpoilingTimeMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the spoiling time of perishables globally. Higher values prolong the time.
| version = [[189.0]]
}}
{{Server config variable
| name = HairGrowthSpeedMultiplier
| type = float
| default = 1.0 (ASE), 0 (ASA)
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the hair growth. Higher values increase speed of growth.
| version = [[254.0]]
}}
{{Server config variable
| name = HarvestResourceItemAmountClassMultipliers
| type = (...)
| default = N/A
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales on a per-resource type basis, the amount of resources harvested. See [[#HarvestResourceItemAmountClassMultipliers|Items related]] section for more details.
| version = [[189.0]]
}}
{{Server config variable
| name = HarvestXPMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the amount of XP earned for harvesting.
| version = [[243.0]]
}}
{{Server config variable
| name = IncreasePvPRespawnIntervalBaseAmount
| type = float
| default = 60.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>bIncreasePvPRespawnInterval</code> is <code>True</code>, sets the additional [[PvP]] re-spawn time in seconds that scales (<code>IncreasePvPRespawnIntervalMultiplier</code>) when a player is killed by a team within a certain amount of time (<code>IncreasePvPRespawnIntervalCheckPeriod</code>).
| version = [[196.0]]
}}
{{Server config variable
| name = IncreasePvPRespawnIntervalCheckPeriod
| type = float
| default = 300.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>bIncreasePvPRespawnInterval</code> is <code>True</code>, sets the amount of time in seconds within a player re-spawn time increases (<code>IncreasePvPRespawnIntervalBaseAmount</code>) and scales (<code>IncreasePvPRespawnIntervalMultiplier</code>) when it is killed by a team in [[PvP]].
| version = [[196.0]]
}}
{{Server config variable
| name = IncreasePvPRespawnIntervalMultiplier
| type = float
| default = 2.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>bIncreasePvPRespawnInterval</code> is <code>True</code>, scales the [[PvP]] additional re-spawn time (<code>IncreasePvPRespawnIntervalBaseAmount</code>) when a player is killed by a team within a certain amount of time (<code>IncreasePvPRespawnIntervalCheckPeriod</code>).
| version = [[196.0]]
}}
{{Server config variable
| name = ItemStatClamps[<attribute>]
| type = value
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Globally clamps item stats. See [[#ItemStatClamps|Items related]] section for more details. Requires <code>?ClampItemStats=true</code> option in the command line.
| version = [[255.0]]
}}
{{Server config variable
| name = KillXPMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scale the amount of XP earned for a kill.
| version = [[243.0]]
}}
{{Server config variable
| name = LayEggIntervalMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the time between eggs are spawning / being laid. Higher number increases it (by percentage).
| version = [[218.0]]
}}
{{Server config variable
| name = LevelExperienceRampOverrides
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Configures the total number of levels available to players and tame creatures and the experience points required to reach each level. See [[#Players and tames levels override related|Players and tames levels override]] section for more details.
| version = [[179.0]]
}}
{{Server config variable
| name = LimitNonPlayerDroppedItemsCount
| type = integer
| default = 0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Limits the number of dropped items in the area defined by <code>LimitNonPlayerDroppedItemsRange</code>. Official servers have it set to <code>600</code>.
| version = [[302.4]]
}}
{{Server config variable
| name = LimitNonPlayerDroppedItemsRange
| type = integer
| default = 0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Sets the area range (in Unreal Units) in which the option <code>LimitNonPlayerDroppedItemsCount</code> applies. Official servers have it set to <code>1600</code>.
| version = [[302.4]]
}}
{{Server config variable
| name = MatingIntervalMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the interval between tames can mate. A lower value decreases it (on a percentage scale). Example: a value of 0.5 would allow tames to mate 50% sooner.
| version = [[219.0]]
}}
{{Server config variable
| name = MatingSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the speed at which tames mate with each other. A higher value increases it (by percentage). Example: MatingSpeedMultiplier=2.0 would cause tames to complete mating in half the normal time.
}}
{{Server config variable
| name = MaxAlliancesPerTribe
| type = integer
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If set, defines the maximum alliances a tribe can form or be part of.
| version = [[265.0]]
}}
{{Server config variable
| name = MaxFallSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Defines the falling speed multiplier at which players starts taking fall damage. The falling speed is based on the time players spent in the air while having a negated Z axis velocity meaning that the higher this setting is, the longer players can fall without taking fall damage. For example, having it set to <code>0.1</code> means players will no longer survive a regular jump while having it set very high such as to <code>100.0</code> means players will survive a fall from the sky limit, etc. This setting doesn't affect the gravity scale of the players so there won't be any physics difference to the character movements.
| version = [[279.224]]
}}
{{Server config variable
| name = MaxNumberOfPlayersInTribe
| type = integer
| default = 0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Sets the maximum survivors allowed in a tribe. A value of 1 effectively disables tribes. The default value of 0 means there is no limit about how many survivors can be in a tribe.
| version = [[242.0]]
}}
{{Server config variable
| name = MaxTribeLogs
| type = integer
| default = 400
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Sets how many Tribe log entries are displayed for each tribe.
| version = [[242.0]]
}}
{{Server config variable
| name = MaxTribesPerAlliance
| type = integer
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If set, defines the maximum of tribes in an alliance.
| version = [[265.0]]
}}
{{Server config variable
| name = NPCReplacements
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Globally replaces specific creatures with another using class names. See [[#NPCReplacements|Creature Spawn related]] section for more detail.
| version = [[196.0]]
}}
{{Server config variable
| name = OverrideMaxExperiencePointsDino
| type = integer
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Overrides the max XP cap of tame characters by exact specified amount.
| version = [[189.0]]
}}
{{Server config variable
| name = OverrideMaxExperiencePointsPlayer
| type = integer
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Overrides the max XP cap of players characters by exact specified amount.
| version = [[189.0]]
}}
{{Server config variable
| name = OverrideEngramEntries
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Configures the status and requirements for learning an engram, specified by its index. See [[#OverrideEngramEntries_and_OverrideNamedEngramEntries|Engram Entries related]] section for more detail.
| version = [[179.0]]
}}
{{Server config variable
| name = OverrideNamedEngramEntries
| type = (...)
| default = N/A
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Configures the status and requirements for learning an engram, specified by its name. See [[#OverrideEngramEntries_and_OverrideNamedEngramEntries|Engram Entries related]] section for more detail.
| version = [[179.0]]
}}
{{Server config variable
| name = OverridePlayerLevelEngramPoints
| type = integer
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Configures the number of engram points granted to players for each level gained. This option must be repeated for each player level set on the server, e.g.: if there are 65 player levels available this option should appear 65 times, with the first setting the engram points for reaching level 1, the next one setting engram points for level 2 and so on, all the way to the 65th which configures engram points for level 65.
 OverridePlayerLevelEngramPoints=5
 OverridePlayerLevelEngramPoints=20
 ...
 OverridePlayerLevelEngramPoints=300
| version = [[179.0]]
}}
{{Server config variable
| name = PassiveTameIntervalMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales how often a survivor get tame requests for passive tame creatures.
| version = [[278.0]]
}}
{{Server config variable
| name = PerLevelStatsMultiplier_Player['''<integer>''']
| type = float
| default = N/A
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales Player stats. See [[#PerLevelStatsMultiplier|Level stats related]] section for more detail.
| version = [[202.0]]
}}
{{Server config variable
| name = PerLevelStatsMultiplier_DinoTamed'''<_type>'''['''<integer>''']
| type = float
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales tamed creature stats. See [[#PerLevelStatsMultiplier|Level stats related]] section for more detail.
| version = [[202.0]]
}}
{{Server config variable
| name = PerLevelStatsMultiplier_DinoWild['''<integer>''']
| type = float
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales wild creatures stats. See [[#PerLevelStatsMultiplier|Level stats related]] section for more detail.
| version = [[202.0]]
}}
{{Server config variable
| name = PhotoModeRangeLimit
| type = integer
| default = 3000
| cli = skip
| inASA = Yes
| inASE = No
| info = Defines the maximum distance between photo mode camera position and player position.
}}
{{Server config variable
| name = PlayerBaseStatMultipliers[<attribute>]
| type = multiplier
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Changes the base stats of a player by multiplying with the default value. Meaning the start stats of a new spawned character. See [[#PlayerBaseStatMultipliers|Stats related]] section for more detail.
| version = [[254.6]]
}}
{{Server config variable
| name = PlayerHarvestingDamageMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales the damage done to a harvestable item/entity by a Player. A higher value increases it (by percentage): the higher number, the faster the survivors collects.
| version = [[231.1]]
}}
{{Server config variable
| name = PoopIntervalMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales how frequently survivors can poop. Higher value decreases it (by percentage)
| version = [[218.0]]
}}
{{Server config variable
| name = PreventBreedingForClassNames
| type = "<string>"
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Prevents breeding of specific creatures via classname. E.g. <code>PreventBreedingForClassNames="Argent_Character_BP_C"</code>. Creature classnames can be found on the [[Creature IDs]] page.
}}
{{Server config variable
| name = PreventDinoTameClassNames
| type = "<string>"
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Prevents taming of specific dinosaurs via classname. E.g. <code>PreventDinoTameClassNames="Argent_Character_BP_C"</code>. Dino classnames can be found on the [[Creature IDs]] page.
| version = [[194.0]]
}}
{{Server config variable
| name = PreventOfflinePvPConnectionInvincibleInterval
| type = float
| default = 5.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Specifies the time in seconds a player cannot take damages after logged-in.
| version = [[278.0]]
}}
{{Server config variable
| name = PreventTransferForClassNames
| type = "<string>"
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Prevents transfer of specific creatures via classname. E.g. <code>PreventTransferForClassNames="Argent_Character_BP_C"</code>Creature classnames can be found on the [[Creature IDs]] page.
| version = [[326.13]]
}}
{{Server config variable
| name = PvPZoneStructureDamageMultiplier
| type = float
| default = 6.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Specifies the scaling factor for damage structures take within caves. The lower the value, the less damage the structure takes (i.e. setting to 1.0 will make structure built in or near a cave receive the same amount of damage as those built on the surface).
| version = [[187.0]]
}}
{{Server config variable
| name = ResourceNoReplenishRadiusPlayers
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Controls how resources regrow closer or farther away from players. Values higher than 1.0 increase the distance around players where resources are not allowed to grow back. Values between 0 and 1.0 will reduce it.
| version = [[196.0]]
}}
{{Server config variable
| name = ResourceNoReplenishRadiusStructures
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Controls how resources regrow closer or farther away from structures Values higher than 1.0 increase the distance around structures where resources are not allowed to grow back. Values between 0 and 1.0 will reduce it.
| version = [[196.0]]
}}
{{Server config variable
| name = SpecialXPMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scale the amount of XP earned for SpecialEvent.
| version = [[243.0]]
}}
{{Server config variable
| name = StructureDamageRepairCooldown
| type = integer
| default = 180
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Option for cooldown period on structure repair from the last time damaged. Set to 180 seconds by default, 0 to disable it.
| version = [[222.0]]
}}
{{Server config variable
| name = SupplyCrateLootQualityMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Increases the quality of items that have a quality in the supply crates. Valid values are from 1.0 to 5.0. The quality also depends on the Difficulty Offset.
| version = [[260.0]]
}}
{{Server config variable
| name = TamedDinoCharacterFoodDrainMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales how fast tame creatures consume food.
| version = [[278.0]]
}}
{{Server config variable
| name = TamedDinoClassDamageMultipliers
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Globally overrides tamed creatures damages. See [[#DinoClassDamageMultipliers|Creature Stats related]] section for more details.
| version = [[194.0]]
}}
{{Server config variable
| name = TamedDinoClassResistanceMultipliers
| type = (...)
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Globally overrides tamed creatures resistance. See [[#DinoClassResistanceMultipliers|Creature Stats related]] section for more details.
| version = [[194.0]]
}}
{{Server config variable
| name = TamedDinoTorporDrainMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales how fast tamed creatures lose torpor.
| version = [[278.0]]
}}
{{Server config variable
| name = TribeSlotReuseCooldown
| type = float
| default = 0.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Locks a tribe slot for the value in seconds, e.g.: a value of 3600 would mean that if a survivor leaves the tribe, their place cannot be taken by another survivor (or re-join) for 1 hour. Used on Official Small Tribes Servers.
| version = [[280.114]]
}}
{{Server config variable
| name = UseCorpseLifeSpanMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Modifies corpse and dropped box lifespan.
| version = [[275.0]]
}}
{{Server config variable
| name = WildDinoCharacterFoodDrainMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales how fast wild creatures consume food.
| version = [[278.0]]
}}
{{Server config variable
| name = WildDinoTorporDrainMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales how fast wild creatures lose torpor.
| version = [[278.0]]
}}
|-
! colspan=5 | Turrets limit
{{Server config variable
| name = bHardLimitTurretsInRange
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, enables the retroactive turret hard limit (100 turrets within a 10k unit radius).
| version = [[278.0]]
}}
{{Server config variable
| name = bLimitTurretsInRange
| type = boolean
| default = True
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>False</code>, doesn't limit the maximum allowed automated turrets (including [[Plant Species X]]) in a certain range.
| version = [[274.0]]
}}
{{Server config variable
| name = LimitTurretsNum
| type = integer
| default = 100
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Determines the maximum number of turrets that are allowed in the area.
| version = [[274.0]]
}}
{{Server config variable
| name = LimitTurretsRange
| type = float
| default = 10000.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Determines the area in Unreal Unit in which turrets are added towards the limit.
| version = [[274.0]]
}}
|-
! colspan=5 | Genesis
{{Server config variable
| name = AdjustableMutagenSpawnDelayMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Scales the Mutagen spawn rates. By default, The game attempts to spawn them every 8 hours on dedicated servers, and every hour on non-dedicated servers and single-player. Rising this value will rise the re-spawn interval, lowering will make it shorter.
| version = [[329.5]]
}}
{{Server config variable
| name = BaseHexagonRewardMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the missions score hexagon rewards. Also scales token rewards in Club Ark (ASA).
| version = [[329.5]]
}}
{{Server config variable
| name = bDisableHexagonStore
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, disables the Hexagon store
| version = [[329.5]]
}}
{{Server config variable
| name = bDisableDefaultMapItemSets
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, disables Genesis 2 Tek Suit on Spawn.
| version = [[329.7]]
}}
{{Server config variable
| name = bDisableGenesisMissions
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, disables missions on Genesis.
| version = [[306.79]]
}}
{{Server config variable
| name = bDisableWorldBuffs
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, disables world effects from [[Missions (Genesis: Part 2)]] altogether. To disable specific world buffs, see <code>DisableWorldBuffs</code> of [[#DynamicConfig|#DynamicConfig]].
| version = [[329.51]]
}}
{{Server config variable
| name = bEnableWorldBuffScaling
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, makes world effects from [[Missions (Genesis: Part 2)]] scale from server settings, rather than add/subtract a flat amount to the value at runtime.
| version = [[329.25]]
}}
{{Server config variable
| name = bGenesisUseStructuresPreventionVolumes
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, disables building in mission areas on [[Genesis: Part 1]].
| version = [[306.53]]
}}
{{Server config variable
| name = bHexStoreAllowOnlyEngramTradeOption
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = If <code>True</code>, allows only Engrams to be sold on the Hex Store, disables everything else.
| version = [[329.5]]
}}
{{Server config variable
| name = HexagonCostMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Scales the hexagon cost of items in the Hexagon store. Also scales token cost of items in Club Ark (ASA).
| version = [[329.5]]
}}
{{Server config variable
| name = MutagenLevelBoost['''<Stat_ID>''']
| type = integer
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = States the number of levels {{ItemLink|Mutagen}} adds to tames with wild ancestry. See [[#MutagenLevelBoost|Stats related]] section for more details.
| version = [[329.51]]
}}
{{Server config variable
| name = MutagenLevelBoost_Bred['''<Stat_ID>''']
| type = integer
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Alike as <code>MutagenLevelBoost</code>, but for bred tames. See [[#MutagenLevelBoost_Bred|Stats related]] section for more details.
| version = [[329.51]]
}}
{{Server config variable
| name = WorldBuffScalingEfficacy
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Makes world effects from [[Missions (Genesis: Part 2)]] scaling more or less effective when setting <code>bEnableWorldBuffScaling=True</code>. 1 would be default, 0.5 would be 50% less effective, 100 would be 100x more effective.
| version = [[329.51]]
}}
|}
</div>

====[ModInstaller]====
The section <code>[ModInstaller]</code> handles each extra Steam Workshop Mods/Maps/TC IDs in the ''Game.ini''.
<div class="widetable">
{| class="wikitable config-table"
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Variable
! Description
! Since patch
{{Server config variable
| name = ModIDS
| type = ModID
| default = 
| cli = skip
| inASA = No
| inASE = Yes
| info = Specifies a single Steam Workshop Mods/Maps/TC ID to download/install/update on the server. To handle multiple IDs, multiple lines must be added with the same syntax, each one with the specific workshop ID. Requires <code>-automanagedmods</code> in the command line. Example:
 [ModInstaller]
 ModIDS=123456789
 ModIDS=987654321
 ModIDS=147852369
'''Note:''' to active mods it's required to use the <code>ActiveMods</code> options in <code>[ServerSettings]</code> section of ''GameUserSettings.ini'' or in the command line.
| version = [[244.3]]
}}
|}
<!--end of Game.ini-->

====Creature Spawn related====
<div style="border: 1px solid silver; padding: 0 1em;">
=====ConfigAddNPCSpawnEntriesContainer=====
 ConfigAddNPCSpawnEntriesContainer=(
 [NPCSpawnEntriesContainerClassString="'''<[[Spawn entries|spawn_class]]>'''"],
 [NPCSpawnEntries=(([AnEntryName="'''<spawn_name>'''"],
 [EntryWeight='''<factor>'''],
 [NPCsToSpawnStrings=("'''<[[Creature IDs|entity_id]]>'''")]))],
 [NPCSpawnLimits=(([NPCClassString="'''<[[Creature IDs|entity_id]]>'''"],
 [MaxPercentageOfDesiredNumToAllow=<percentage>]))])
Arguments:
{| class="wikitable"
|-
|<code>[[Spawn entries|spawn_class]]</code> || ''string'' || Spawn Group Container Class Name, see [[Spawn entries]]
|-
|<code>spawn_name</code> || ''string'' || Spawn Name
|-
|<code>factor</code> || ''float'' || Weight Factor for this spawn
|-
|<code>[[Creature IDs|entity_id]]</code> || ''string'' || Entity ID of the creature to be added, see [[Creature IDs]]
|-
|<code>percentage</code> || ''float'' || Sets the maximum allowed creatures for this instance
|}
Allows specified creatures to spawn in specified locations directly through the spawners limited to what is specified.

Examples:

Adding Giganotosaurus to the beach spawn area:
 ConfigAddNPCSpawnEntriesContainer=(
   NPCSpawnEntriesContainerClassString="DinoSpawnEntriesBeach_C",
   NPCSpawnEntries=((AnEntryName="GigaSpawner", EntryWeight=1.0, NPCsToSpawnStrings=("Gigant_Character_BP_C"))),
   NPCSpawnLimits=((NPCClassString="Gigant_Character_BP_C", MaxPercentageOfDesiredNumToAllow=0.01))
 )
Adding a pack of two Dodos and a Rex to the [[Herbivore Island]]:
 ConfigAddNPCSpawnEntriesContainer=(
   NPCSpawnEntriesContainerClassString="DinoSpawnEntriesDamiensAtoll_C",
   NPCSpawnEntries=(
     (AnEntryName="Dodos (2)", EntryWeight=1.0, NPCsToSpawnStrings=("Dodo_Character_BP_C","Dodo_Character_BP_C")),
     (AnEntryName="Rex (1)", EntryWeight=0.5, NPCsToSpawnStrings=("Rex_Character_BP_C"))
   ),
   NPCSpawnLimits=(
     (NPCClassString="Dodo_Character_BP_C", MaxPercentageOfDesiredNumToAllow=0.5),
     (NPCClassString="Rex_Character_BP_C", MaxPercentageOfDesiredNumToAllow=0.01)
   )
 )
''The examples provided here are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

More advanced arguments are treated in [[#Advanced_*NPCSpawnEntriesContainer_usage|Advanced *NPCSpawnEntriesContainer usage]].

=====ConfigSubtractNPCSpawnEntriesContainer=====
 ConfigSubtractNPCSpawnEntriesContainer=(
 [NPCSpawnEntriesContainerClassString="'''<[[Spawn entries|spawn_class]]>'''"],
 [NPCSpawnEntries=((NPCsToSpawnStrings=("'''<[[Creature IDs|entity_id]]>'''")))],
 [NPCSpawnLimits=((NPCClassString="'''<[[Creature IDs|entity_id]]>'''"))])
Arguments:
{| class="wikitable"
|-
| <code>[[Spawn entries|spawn_class]]</code> || ''string'' || Spawn Class String, see [[Spawn entries]]
|-
|<code>[[Creature IDs|entity_id]]</code> || ''string'' || Entity ID of the creature to be removed, see [[Creature IDs]]
|}
Completely removes specified creatures from specified locations. More than one creature can be specified.

You cannot reference the same Spawn Entry in multiple lines and have them all take effect, even if the Entity ID or Spawn Class referenced is different between each line. If removing multiple creatures from the same Spawn Entry, use one line to do so.

Example: removing Trike and Ptera from the Beach.
 ConfigSubtractNPCSpawnEntriesContainer=(
   NPCSpawnEntriesContainerClassString="DinoSpawnEntriesBeach_C",
   NPCSpawnEntries=(
     (NPCsToSpawnStrings=("Trike_Character_BP_C")),
     (NPCsToSpawnStrings=("Ptero_Character_BP_C"))
  ),
  NPCSpawnLimits=(
     (NPCClassString="Trike_Character_BP_C"),
     (NPCClassString="Ptero_Character_BP_C")
   )
 )
''The example provided is are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

More advanced arguments are treated in [[#Advanced_*NPCSpawnEntriesContainer_usage|Advanced *NPCSpawnEntriesContainer usage]].

=====ConfigOverrideNPCSpawnEntriesContainer=====
 ConfigOverrideNPCSpawnEntriesContainer=(
 [NPCSpawnEntriesContainerClassString="'''<[[Spawn entries|spawn_class]]>'''"],
 [NPCSpawnEntries=(([AnEntryName="'''<spawn_name>'''"],
 [EntryWeight='''<factor>'''],
 [NPCsToSpawnStrings=("'''<[[Creature IDs|entity_id]]>'''")]))],
 [NPCSpawnLimits=(([NPCClassString="'''<[[Creature IDs|entity_id]]>'''"],
 [MaxPercentageOfDesiredNumToAllow='''<percentage>''']))])
Arguments:
{| class="wikitable"
|-
|<code>[[Spawn entries|spawn_class]]</code> || ''string'' || Spawn Class String, see [[Spawn entries]]
|-
|<code>spawn_name</code> || ''string'' || Spawn Name
|-
|<code>factor</code> || ''float'' || Weight Factor for this spawn
|-
|<code>[[Creature IDs|entity_id]]</code> || ''string'' || Entity ID of the creature to override other spawns, see [[Creature IDs]]
|-
|<code>percentage</code> || ''float'' || Sets the maximum allowed creatures for this instance
|}
Overrides all creatures within a specified area with specified creatures, allowing areas to be dedicated for the spawning of specific creatures. Does not prevent wanderers from entering said location.

Example: only Rex and Carno spawns in the Mountain Region
 ConfigOverrideNPCSpawnEntriesContainer=(
   NPCSpawnEntriesContainerClassString="DinoSpawnEntriesMountain_C",
   NPCSpawnEntries=(
     (AnEntryName="CarnoSpawner", EntryWeight=1.0, NPCsToSpawnStrings=("Carno_Character_BP_C")),
     (AnEntryName="RexSpawner", EntryWeight=1.0, NPCsToSpawnStrings=("Rex_Character_BP_C"))
   ),
   NPCSpawnLimits=(
     (NPCClassString="Carno_Character_BP_C", MaxPercentageOfDesiredNumToAllow=0.5),
     (NPCClassString="Rex_Character_BP_C", MaxPercentageOfDesiredNumToAllow=0.5)
   )
 )
''The example provided here is split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

More advanced arguments are treated in [[#Advanced_*NPCSpawnEntriesContainer_usage|Advanced *NPCSpawnEntriesContainer usage]].

=====Advanced *NPCSpawnEntriesContainer usage=====
The following arguments can be used to the <code>*NPCSpawnEntriesContainer</code> options:
*To set the distance for each spawn entry:
 [NPCsSpawnOffsets=((X='''<float>''',Y='''<float>''',Z='''<float>'''),...,(X='''<float>''',Y='''<float>''',Z='''<float>'''))]
:where each ''X,Y,Z'' sets the distance in unreal units from the spawn point for each entry.
*To set the spawn percentage change for each spawn entry:
 [NPCsToSpawnPercentageChance=('''<float>''',...,'''<float>''')]
:where each value sets chance from 0 to 1.0 (100%) for an entry to be in that spawn.
*To set a spread radius from the spawn point:
 [ManualSpawnPointSpreadRadius='''<float>''']
:where the values sets the distance in unreal unit from the from the spawn point.
*To set specific difficulty level ranges for each spawn entry:
 [NPCDifficultyLevelRanges=(
 ([EnemyLevelsMin=('''<float>''',...,'''<float>''')],
 [EnemyLevelsMax=('''<float>''',...,'''<float>''')],
 [GameDifficulties=('''<float>''',...,'''<float>''')]),
 ...
 ([EnemyLevelsMin=('''<float>''',...,'''<float>''')],
 [EnemyLevelsMax=('''<float>''',...,'''<float>''')],
 [GameDifficulties=('''<float>''',...,'''<float>''')]))]
:where for each spawn entry there is a group of <code>EnemyLevelsMin</code>, <code>EnemyLevelsMax</code> and  <code>GameDifficulties</code>. Each of those settings define respectively base levels, the max levels and the difficult offsets on which wild creatures scale. See [[Difficulty|Difficulty Level]] for more details.
*To set the minimum and maximum spawn offsets in height, useful for flyers and water creatures:
 [RandGroupSpawnOffsetZMin='''<float>'''],
 [RandGroupSpawnOffsetZMax='''<float>''']
:where they sets respectively the minimum and maximum offsets in unreal unit.
Example: add small groups (from 1 to 2 samples) of [[Anglerfish]] in the <code>DinoSpawnEntries_Ocean_C</code> spawn container, with a linear level distribution up to <code>DifficultyOffset=5.0</code>
 ConfigAddNPCSpawnEntriesContainer=(
   NPCSpawnEntriesContainerClassString="DinoSpawnEntries_DarkWaterAngler_C",
   NPCSpawnEntries=(
     (AnEntryName="Anglers(1-2)",
      EntryWeight=0.9,
      NPCsToSpawnStrings=("Angler_Character_BP_C","Angler_Character_BP_C"),
      NPCsSpawnOffsets=(
        (X=0.0,Y=0.0,Z=0.0),
        (X=0.0,Y=250.0,Z=0.0)
      ),
      NPCsToSpawnPercentageChance=(1.0,0.5),
      ManualSpawnPointSpreadRadius=650.0,
      NPCDifficultyLevelRanges=(
        (EnemyLevelsMin=(1.0,1.0,1.0,1.0,1.0,1.0),
         EnemyLevelsMax=(30.0,30.0,30.0,30.0,30.0,30.0),
         GameDifficulties=(0.0,1.0,2.0,3.0,4.0,5.0)),
        (EnemyLevelsMin=(1.0,1.0,1.0,1.0,1.0,1.0),
         EnemyLevelsMax=(30.0,30.0,30.0,30.0,30.0,30.0),
         GameDifficulties=(0.0,1.0,2.0,3.0,4.0,5.0))
      )
      RandGroupSpawnOffsetZMin=1200.0,
      RandGroupSpawnOffsetZMax=5000.0)),
   NPCSpawnLimits=(
     (NPCClassString="Angler_Character_BP_C",
      MaxPercentageOfDesiredNumToAllow=0.35)
   )
 )
''The example provided here is split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

*More arguments about levels:
 [NPCOverrideLevel='''<integer>''']
 [NPCMaxLevelOffset='''<float>''']
 [NPCMaxLevelOffset='''<float>'''] 
 [NPCMinLevelMultiplier='''<float>''']
 [NPCMaxLevelMultiplier='''<float>''']
 [bAddLevelOffsetBeforeMultiplier='''<1>''']
:it is possible to combine such arguments in a similar way to <code>NPCDifficultyLevelRanges</code> to override base levels of creatures, their minimum and max levels as well their relative multiplier.
:In a game where <code>DifficultyOffset=5.0</code>, to set only Parasaur at level 1 to spawn in [[The Dead Island]]:
 ConfigOverrideNPCSpawnEntriesContainer=(
   NPCSpawnEntriesContainerClassString="DinoSpawnEntriesMonsterIsland_C",
   NPCSpawnEntries=(
     (AnEntryName="Parasaur lvl 1",
      NPCsToSpawnStrings=("Para_Character_BP_C"),
      NPCOverrideLevel=(1))
   )
 )
:To set them all spawn at max level add:
      NPCOverrideLevel=(1),
      NPCMinLevelOffset=(30.0),
      NPCMaxLevelOffset=(30.0))))
:To double their max level change in:
     (AnEntryName="Parasaur lvl 1",
      NPCMinLevelMultiplier=(1.0),
      NPCMaxLevelMultiplier=(2.0))))
:To increase their level by 150 change in:
      NPCsToSpawnStrings=("Para_Character_BP_C"),
      NPCMinLevelOffset=(30.0),
      NPCMaxLevelOffset=(30.0))))
:To double their max level and set a minimum level of 150 change in:
      NPCsToSpawnStrings=("Para_Character_BP_C"),
      NPCMinLevelOffset=(30.0),
      NPCMaxLevelOffset=(30.0),
      NPCMinLevelMultiplier=(1.0),
      NPCMaxLevelMultiplier=(2.0))))
:To increase their minimum level by 150 and then double them change in:
      NPCsToSpawnStrings=("Para_Character_BP_C"),
      bAddLevelOffsetBeforeMultiplier=1,
      NPCMinLevelOffset=(30.0),
      NPCMaxLevelOffset=(30.0),
      NPCMinLevelMultiplier=(1.0),
      NPCMaxLevelMultiplier=(2.0))))
''The examples provided here are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====DinoSpawnWeightMultipliers=====
 DinoSpawnWeightMultipliers=(
 DinoNameTag='''<[[Creature IDs|tag]]>'''
 [,SpawnWeightMultiplier='''<factor>''']
 [,OverrideSpawnLimitPercentage='''<override>''']
 [,SpawnLimitPercentage='''<limit>'''])
Arguments:
{| class="wikitable"
|-
|<code>[[Creature IDs|tag]]</code> || ''string'' || Creature type to adjust, see [[Creature IDs]] (column Name Tags)
|-
|<code>factor</code> || ''float'' || Weight factor for this type
|-
|<code>override</code> || ''booleaon'' || If <code>True</code>, use the specified <code>SpawnLimitPercentage</code>
|-
|<code>limit</code> || ''float'' || Maximum percentage (among all spawns) for this type
|}
Customizes the spawning rate for a given creature type (at all creature spawn points). Types with a larger <code>SpawnWeightMultiplier</code> are selected more often when spawning new creatures than types with lower multipliers. When <code>OverrideSpawnLimitPercentage</code> is specified (and <code>True</code>), the type will never be spawned more than <code>SpawnLimitPercentage * 100</code> percent of the time, regardless of multiplier. For example, a <code>SpawnLimitPercentage</code> of <code>0.25</code> specifies that the type will be selected for spawning no more than 25% of the time. Note that the <code>SpawnLimitPercentage</code> is calculated per spawn region (aka 'Spawn Entry Container') and not against the global spawn population. A limit of 0.1 means 10% of the creatures in any given region and *NOT* 10% of the map-size population. 

Multiple <code>DinoSpawnWeightMultipliers</code> entries can be specified in the file, but <code>DinoNameTag</code> values should not be repeated across multiple entries.
 DinoSpawnWeightMultipliers=(
   DinoNameTag=Bronto,
   SpawnWeightMultiplier=10.0,
   OverrideSpawnLimitPercentage=True,
   SpawnLimitPercentage=0.5
 )
''The example provided here is split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====NPCReplacements=====
 NPCReplacements=(FromClassName="'''<classname>'''", ToClassName="'''<classname>'''")
Arguments:
{| class="wikitable"
|-
|<code>[[Creature IDs|classname]]</code> || ''string'' || Spawn Class String, see [[Creature IDs]] (column Entity ID)
|}
This can be used to disable specific Alpha Creatures, replace the spawns of a particular NPC (NPC=Non-player character, i.e., a creature) with that of a different NPC, or completely disable any specific NPC spawn. 
 NPCReplacements=(FromClassName="MegaRaptor_Character_BP_C", ToClassName="Dodo_Character_BP_C")
Dynamic Config:

As of [[319.14]], NPC replacements can be defined through the dynamic config. However, the syntax is slightly different. Rather than having multiple entries each defining a single replacement, it must be a single-entry mapping all of them. For example (note the extra brackets):
 NPCReplacements=((FromClassName="MegaRaptor_Character_BP_C", ToClassName="Dodo_Character_BP_C"))
Any additional replacements must be mapped inside that set separated by a comma like so:
 NPCReplacements=((FromClassName="MegaRaptor_Character_BP_C", ToClassName="Dodo_Character_BP_C"), (FromClassName="Coel_Character_BP_C", ToClassName="Piranha_Character_BP_C"))
</div>

====Creature stats related====
<div style="border: 1px solid silver; padding: 0 1em;">
=====DinoClassDamageMultipliers=====
 DinoClassDamageMultipliers=(ClassName="'''<string>'''",Multiplier='''<float>''')
 TamedDinoClassDamageMultipliers=(ClassName="'''<string>'''",Multiplier='''<float>''')
Arguments:
{| class="wikitable"
|-
|<code>string</code> = creature ''classname''<br/><code>float</code> = ''multiplier'', default: 1.0
|}
Multiplies damage dealt of specific creatures via ''classname''. Higher values increase the damage dealt.

Creature classnames can be found on the [[Creature IDs]] page. Multiple <code>DinoClassDamageMultipliers</code> and <code>TamedDinoClassDamageMultipliers</code> entries can be specified in the file, but ''classname'' values should not be repeated across multiple entries.

Examples:
 DinoClassDamageMultipliers=(
   ClassName="MegaRex_Character_BP_C",
   Multiplier=0.1
 )

 TamedDinoClassDamageMultipliers=(
   ClassName="Rex_Character_BP_C",
   Multiplier=10.0
 )
''The examples provided here are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====DinoClassResistanceMultipliers=====
 DinoClassResistanceMultipliers=(ClassName="'''<string>'''",Multiplier='''<float>''')
 TamedDinoClassResistanceMultipliers=(ClassName="'''<string>'''",Multiplier='''<float>''')
Arguments:
{| class="wikitable"
|-
|<code>string</code> = creature ''classname''<br/><code>float</code> = ''multiplier'', default = 1.0
|}
Multiplies resistance of specific creatures via ''classname''. Higher values decrease the damage received.

Creature classnames can be found on the [[Creature IDs]] page.

The examples provided here are split into multiple lines for space considerations. In the configuration file, an entry must be placed on a single line. Multiple <code>DinoClassResistanceMultipliers</code> and <code>TamedDinoClassResistanceMultipliers</code> entries can be specified in the file, but ''classname'' values should not be repeated across multiple entries.

Examples:
 DinoClassResistanceMultipliers=(
   ClassName="MegaRex_Character_BP_C",
   Multiplier=0.1
 )

 TamedDinoClassResistanceMultipliers=(
   ClassName="Rex_Character_BP_C",
   Multiplier=10.0
 )
''The examples provided here are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====TamedDinoClassSpeedMultipliers=====
 TamedDinoClassSpeedMultipliers=(ClassName="'''<string>'''",Multiplier='''<float>''')
Arguments:
{| class="wikitable"
|-
|<code>string</code> = creature ''classname''<br/><code>float</code> = ''multiplier'', default: 1.0
|}
Multiplies base speed of specific creatures via ''classname''. Higher values increase the speed.

Creature classnames can be found on the [[Creature IDs]] page. Multiple <code>TamedDinoClassSpeedMultipliers</code> entries can be specified in the file, but ''classname'' values should not be repeated across multiple entries.

 TamedDinoClassSpeedMultipliers=(
   ClassName="Argent_Character_BP_C",
   Multiplier=2.0
 )
''The examples provided here are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====TamedDinoClassStaminaMultipliers=====
 TamedDinoClassStaminaMultipliers=(ClassName="'''<string>'''",Multiplier='''<float>''')
Arguments:
{| class="wikitable"
|-
|<code>string</code> = creature ''classname''<br/><code>float</code> = ''multiplier'', default: 1.0
|}
Multiplies the Stamina of specific creatures via ''classname''. Higher values increase the stamina.

Creature classnames can be found on the [[Creature IDs]] page. Multiple <code>TamedDinoClassStaminaMultipliers</code> entries can be specified in the file, but ''classname'' values should not be repeated across multiple entries.

 TamedDinoClassStaminaMultipliers=(
   ClassName="Ptero_Character_BP_C",
   Multiplier=5.0
 )
''The examples provided here are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''
</div>

====Engram Entries related====
<div style="border:1px solid silver; padding:0 1em;">
=====OverrideEngramEntries and OverrideNamedEngramEntries=====
 OverrideEngramEntries=(EngramIndex='''<[[Engram classnames|index]]>'''
 [,EngramHidden='''<hidden>''']
 [,EngramPointsCost='''<cost>''']
 [,EngramLevelRequirement='''<level>''']
 [,RemoveEngramPreReq='''<remove_prereq>'''])

 OverrideNamedEngramEntries=(EngramClassName="'''<[[Engram classnames|class_name]]>'''"
 [,EngramHidden='''<hidden>''']
 [,EngramPointsCost='''<cost>''']
 [,EngramLevelRequirement='''<level>''']
 [,RemoveEngramPreReq='''<remove_prereq>'''])
Arguments:
{| class="wikitable"
|-
|<code>[[Engram classnames|index]]</code> || ''integer'' || Index of the engram
|-
|<code>[[Engram classnames|class_name]]</code> || ''integer'' || Class name of the engram
|-
|<code>hidden</code> || ''boolean'' || If <code>True</code>, hides the engram in the players' Engrams panel
|-
|<code>cost</code> || ''integer'' || Engram points needed to learn engram
|-
|<code>level</code> || ''integer'' || Minimum level needed to learn engram
|-
|<code>remove_prereq</code> || ''boolean'' || If <code>True</code>, removes the need of prerequisite engrams to learn this engram.
|}
Configures the status and requirements for learning an engram. For OverrideEngramEntries the <code>EngramIndex</code> argument is always required, for OverrideNamedEngramEntries the <code>EngramClassName</code> argument is always required; the rest are optional, but at least one must be provided in order for the option to have any effect. The option may be repeated in ''Game.ini'' once for each engram to be configured.

For the Engram Index and Engram Class Name see [[Engram classnames]].

The examples provided here are split into multiple lines for space considerations. In the configuration file, an entry must be placed entirely on a single line and without spaces. Multiple <code>OverrideEngramEntries</code> and <code>OverrideNamedEngramEntries</code> entries can be specified in the file, but <code>EngramIndex</code> and <code>EngramClassName</code> values should not be repeated across multiple entries.

Examples:
 OverrideEngramEntries=(EngramIndex=0,EngramHidden=False)

 OverrideEngramEntries=(EngramIndex=1,EngramHidden=False,EngramPointsCost=3,EngramLevelRequirement=3,RemoveEngramPreReq=True)

 OverrideNamedEngramEntries=(EngramClassName="EngramEntry_Campfire_C",EngramHidden=False)

 OverrideNamedEngramEntries=(EngramClassName="EngramEntry_StoneHatchet_C",EngramHidden=False,EngramPointsCost=3,EngramLevelRequirement=3, RemoveEngramPreReq=True)

=====EngramEntryAutoUnlocks=====
<code>EngramEntryAutoUnlocks=(EngramClassName="'''<string>'''",LevelToAutoUnlock='''<integer>''')</code>
{| class="wikitable"
|-
|<code>string</code> || [[Engram class names|Engram Classname]]
|-
|<code>integer</code> || Level you need to gain to unlock the Engram automatically
|}
Automatically unlocks the specified Engram when reaching the level specified.

This example unlocks the [[Tek Teleporter]] with level 0:
 EngramEntryAutoUnlocks=(
   EngramClassName="EngramEntry_TekTeleporter_C",
   LevelToAutoUnlock=0
 )
''The example provided here is split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''
</div>

====Items related====
<div style="border: 1px solid silver; padding: 0 1em;">
Every Item Class String can be found in the Dev Kit. Currently doesn't change the repair cost and demolish refund of edited structures. This can result in potential exploit for lowered crafting costs and may make structures unrepairable.

'''Note:''' if using stack mods, refer to the mod new resources instead of vanilla ones (i.e.: PrimalItemResource_Electronics_Child_C instead of PrimalItemResource_Electronics_C).
=====ConfigOverrideItemCraftingCosts=====
This is an example of how to make the Hatchet require 1 thatch and 2 stone arrows to craft.
 ConfigOverrideItemCraftingCosts=(
   ItemClassString="PrimalItem_WeaponStoneHatchet_C",
   BaseCraftingResourceRequirements=(
     (ResourceItemTypeString="PrimalItemResource_Thatch_C",
      BaseResourceRequirement=1.0,
      bCraftingRequireExactResourceType=False),
     (ResourceItemTypeString="PrimalItemAmmo_ArrowStone_C",
      BaseResourceRequirement=2.0,
      bCraftingRequireExactResourceType=False)
   )
 )
Here is an example to make the torch requiring 3 raw meat and 2 cooked meat to craft (because ya know, Meat Torches are the best torches!)
 ConfigOverrideItemCraftingCosts=(
   ItemClassString="PrimalItem_WeaponTorch_C",
   BaseCraftingResourceRequirements=(
     (ResourceItemTypeString="PrimalItemConsumable_RawMeat_C",
      BaseResourceRequirement=3.0,
      bCraftingRequireExactResourceType=False),
     (ResourceItemTypeString="PrimalItemConsumable_CookedMeat_C",
      BaseResourceRequirement=2.0,
      bCraftingRequireExactResourceType=False)
   )
 )
''The examples provided here are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====ConfigOverrideItemMaxQuantity=====
 ConfigOverrideItemMaxQuantity=(ItemClassString="'''<string>'''",Quantity=(MaxItemQuantity='''<integer>''', bIgnoreMultiplier='''<boolean>'''))
Arguments:
{| class="wikitable"
|-
|<code>string</code> || [[Item IDs|Class Name]] of the item that will be overridden to new stack size
|-
|<code>integer</code> || New stack size of the specified item
|-
|<code>boolean</code> || If <code>False</code>, that means that the real stack size is ''ItemStackSizeMultiplier * MaxItemQuantity'' for that item<br/>
if <code>True</code>, it uses ''MaxItemQuantity'' directly and ignores the multiplier
|}
Allows manually overriding item stack size on a per-item basis.

Example to make the Tranq Arrows stack to 543 items per stack.
 ConfigOverrideItemMaxQuantity=(
   ItemClassString="PrimalItemAmmo_ArrowTranq_C",
   Quantity=(MaxItemQuantity=543,bIgnoreMultiplier=True)
 )
''The example provided here is split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====ConfigOverrideSupplyCrateItems=====
Allows manually overriding items contained in loot crates. Each loot crate can have one or multiple item sets. Each set can have one or multiple items grouped in one or multiple entries.
 ConfigOverrideSupplyCrateItems=(SupplyCrateClassString='''"<string>"''',MinItemSets='''<integer>''',MaxItemSets='''<integer>''',NumItemSetsPower='''<float>''',bSetsRandomWithoutReplacement='''<boolean>'''[,bAppendItemSets='''<boolean>'''][,bAppendPreventIncreasingMinMaxItemSets='''<boolean>'''],ItemSets=([SetName='''"<string>"''',]MinNumItems='''<integer>''',MaxNumItems='''<integer>''',NumItemsPower='''<float>''',SetWeight='''<float>''',bItemsRandomWithoutReplacement='''<boolean>''',ItemEntries=([ItemEntryName='''"<string>"''',]EntryWeight='''<float>''',ItemClassStrings=('''"<string>"[,...n])''',ItemsWeights=('''<float>[,...n]'''),MinQuantity='''<float>''',MaxQuantity='''<float>''',MinQuality='''<float>''',MaxQuality='''<float>''',bForceBlueprint='''<boolean>''',ChanceToBeBlueprintOverride='''<float>'''[,bForcePreventGrinding='''<boolean>'''][,ItemStatClampsMultiplier='''<float>''']))'''[,...m]''')
Arguments:
{| class="wikitable"
! colspan = 3 | Crate options
|-
| <code>SupplyCrateClassString</code>
| ''string''
| Quoted [[Item IDs|Class Name]] of the loot crate that will be overridden.
|-
| <code>MinItemSets</code>
| ''integer''
| The minimum number of item sets that will be in the loot crate.
|-
| <code>MaxItemSets</code>
| ''integer''
| The maximum number of item sets that will be in the loot crate.
|-
| <code>NumItemSetsPower</code>
| ''float''
| Item quality for all sets in the loot crate. Default: 1.0.
|-
| <code>bSetsRandomWithoutReplacement</code>
| ''boolean''
| If <code>True</code>, ensures any item of any sets will never be added multiple times to the same loot crate.
|-
| <code>bAppendItemSets</code>
| ''boolean''
| Optional, if <code>True</code>, the item sets will be added to the loot crate instead of totally replace its content. Default: <code>False</code>
|-
| <code>bAppendPreventIncreasingMinMaxItemSets</code>
| ''boolean''
| Optional, if <code>True</code>, dynamically increases the amount of items dropped by however many additional item-sets added. If <code>False</code>, the <code>MinItemSets</code> and <code>MaxItemSets</code> original crate values will be increased by the provided values. Requires <code>bAppendItemSets=True</code>. Default: <code>False</code>
|-
| <code>ItemSets</code>
| ''string''
| Defines one or multiple item sets.
|-
! colspan = 3 | Sets options
|-
| <code>SetName</code>
| ''string''
| Optional, quoted name of the set.
|-
| <code>MinNumItems</code>
| ''integer''
| The minimum number of items in this set that can be in the loot crate.
|-
| <code>MaxNumItems</code>
| ''integer''
| The maximum number of items in this set that can be in the loot crate.
|-
| <code>NumItemsPower</code>
| ''float''
| Item quality multiplier for this set.
|-
| <code>SetWeight</code>
| ''float''
| Weight factor for this set to be in the crate. Higher the weight higher the chance. Default: 1.0.
|-
| <code>bItemsRandomWithoutReplacement</code>
| ''boolean''
| If <code>True</code>, ensures each item of this set will never be added twice to the same loot crate.
|-
| <code>ItemEntries</code>
| ''string''
| Defines one or multiple item entries in this set.
|-
! colspan = 3 | Entries options
|-
| <code>ItemEntryName</code>
| ''string''
| Optional, quoted name of the item entry.
|- 
| <code>EntryWeight</code>
| ''float''
| Weight factor for this entry to be in its set. Higher the weight higher the chance. Default: 1.0.
|- 
| <code>ItemClassStrings</code>
| ''string''
| Comma separated [[Item IDs|Class Names]] list of items added in this entry. Each class name must be between quotes. Its elements number must match the <code>ItemsWeights</code> list. Despite its length, the list must be between parenthesis.
|- 
| <code>ItemsWeights</code>
| ''float''
| Comma separated weight factors list of items in this entry. Its elements number must match the <code>ItemClassStrings</code> list. Despite its length, the list must be between parenthesis.
|- 
| <code>MinQuantity</code>
| ''float''
| The minimum number of items if this entry that can be in the loot crate. Note: this value is still limited by its set <code>MinNumItems</code> and <code>MaxNumItems</code> values.
|- 
| <code>MaxQuantity</code>
| ''float''
| The maximum number of items if this entry that can be in the loot crate. Note: this value is still limited by its set <code>MinNumItems</code> and <code>MaxNumItems</code> values.
|- 
| <code>MinQuality</code>
| ''float''
| Items minimum quality in this entry. Default: 1.0.
|- 
| <code>MaxQuality</code>
| ''float''
| Items maximum quality in this entry. Default: 1.0.
|- 
| <code>bForceBlueprint</code>
| ''boolean''
| If <code>True</code>, the items of this entry will be always a blueprint. Default: false.
|- 
| <code>ChanceToBeBlueprintOverride</code>
| ''float''
| Normalized probability of items of this entry to be a blueprint. Note: 1.0 means 100%. Default: 0.0.
|- 
| <code>bForcePreventGrinding</code>
| ''boolean''
| Optional, if <code>True</code>, the items of this entry cannot be ground. Default: false.
|- 
| <code>ItemStatClampsMultiplier</code>
| ''float''
| Optional, normalized value about how much the items can differ from clamps values. Note: 1.0 means 100%. Default: 1.0.
|-
|}
Since patch [[273.7|273.7]] <code>SupplyCrateClassString</code> also takes a part of the class name<ref name="modify_loot_crate_contents>Tom (June 4, 2016). [https://survivetheark.com/index.php?/forums/topic/65930-tutorial-modify-loot-crate-contents-v242/ "Modify loot crate contents"]. [+TUTORIAL+] Modify loot crate contents [v242] </ref>. E.g.: "SupplyCrate" will override all supply crates. 

All of the Item class strings and Supply Crate names can be found in the ARK Dev Kit and on page [[Beacon IDs]].

This example completely overrides the items contained in the regular Level 3 supply crate, to contain just some stone and thatch:
 ConfigOverrideSupplyCrateItems=(
   SupplyCrateClassString="SupplyCrate_Level03_C",
   MinItemSets=1,
   MaxItemSets=1,
   NumItemSetsPower=1.0,
   bSetsRandomWithoutReplacement=True,
   ItemSets=(
     (MinNumItems=2,
      MaxNumItems=2,
      NumItemsPower=1.0,
      SetWeight=1.0,
      bItemsRandomWithoutReplacement=True,
      ItemEntries=(
       (EntryWeight=1.0,
        ItemClassStrings=("PrimalItemResource_Stone_C"),
        ItemsWeights=(1.0),
        MinQuantity=10.0,
        MaxQuantity=10.0,
        MinQuality=1.0,
        MaxQuality=1.0,
        bForceBlueprint=False,
        ChanceToBeBlueprintOverride=0.0),
       (EntryWeight=1.0,
        ItemClassStrings=("PrimalItemResource_Thatch_C"),
        ItemsWeights=(1.0),
        MinQuantity=10.0,
        MaxQuantity=10.0,
        MinQuality=1.0,
        MaxQuality=1.0,
        bForceBlueprint=False,
        ChanceToBeBlueprintOverride=0.0)
       )
     )
   )
 )
''The example provided here is split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====HarvestResourceItemAmountClassMultipliers=====
<code>HarvestResourceItemAmountClassMultipliers=(ClassName="'''<string>'''",Multiplier='''<float>''')</code>
{| class="wikitable"
|-
|<code>'''string'''</code> || Class Name of resource,<br/>see [[Item IDs]]
|-
|<code>'''float'''</code> || Default: 1.0
|}
Scales on a per-resource type basis, the amount of resources harvested. Higher values increase the amount per swing/attack. Resource classnames can be found at [[Item IDs]]. It works in the same way as the global setting <code>HarvestAmountMultiplier</code> but for only the type of resource named on this line. Additional lines can be added with other resource types, such as Wood, Stone etc.

Example that provides twice the amount harvested when harvesting thatch from a tree:
 HarvestResourceItemAmountClassMultipliers=(
   ClassName="PrimalItemResource_Thatch_C",
   Multiplier=2.0
 )
''The example provided here is split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''

=====ItemStatClamps=====
'''NOTE:''' The command line argument <code>ClampItemStats</code> need to be set to <code>True</code> for the clamping to be enabled on your server. See [[#Command_Line_Syntax|Command Line Syntax]].
 ItemStatClamps[<attribute>]='''<value>'''

Arguments:
{| class="wikitable"
|<code>attribute</code> || ''integer'' ||
0: Generic Quality<br/>
1: Armor<br/>
2: Max Durability<br/>
3: Weapon Damage Percent<br/>
4: Weapon Clip Ammo<br/>
5: Hypothermal Insulation<br/>
6: Weight<br/>
7: Hyperthermal Insulation
|-
| <code>value</code> || ''integer'' || The algorithm used appears to be the following:

<code><Initial Value Constant> + ((<ItemStatClamps[<attribute>]> * <State Modifier Scale>) * (<Randomizer Range Multiplier> * <Initial Value Constant>))</code>

Each items have their own specific data which can be found in the Dev Kit.
|}

For example, here are the values needed to have the same clamping as official servers for <code>Armor</code> and <code>Weapon Damage Percent</code>:
 ItemStatClamps[1]=19800
 ItemStatClamps[3]=19800

This would clamp [[Saddles]] to <code>124.0</code> armour (<code>74.5</code> for the 'tank' creatures such as [[Doedicurus]], [[Rock Elemental]], etc.), [[Flak Armor]] pieces to <code>496.0</code> armour, [[Longneck Rifle]] to <code>298.0%</code> damage, etc.

'''WARNING:''' This will permanently change the stats of any existing items so make sure to backup your current save before modifying and playing with the clamping values.
</div>

====Players and tames levels override related====
<div style="border: 1px solid silver; padding: 0 1em;">
 LevelExperienceRampOverrides=(ExperiencePointsForLevel['''<n>''']='''<points>''',[ExperiencePointsForLevel['''<n>''']='''<points>'''],...,[ExperiencePointsForLevel['''<n>''']='''<points>'''])
Arguments:
{| class="wikitable"
|-
|<code>n</code> || ''integer'' || Level to configure
|-
|<code>points</code> || ''integer'' || Points needed to reach level
|}
Configures the total number of levels available to players and creatures and the experience points required to reach each level.

This directive can be specified twice in the configuration file. The first time it appears, the values provided will configure player levels. The second time it appears, the values provided will configure tame creature levels.

Because of this, each time the directive is used, it must list ''all'' the levels players/tames can reach on the server. One <code>ExperiencePointsForLevel</code> argument must appear for each desired level. Values for <code>'''<n>'''</code> must be sequential, starting from zero. Keep in mind that the last 100 levels are used for ascension, chibi experience, explorer notes and rune rewards, meaning that you have to put 100 extra levels in your configuration file.

This first example specifies 50 player levels and 15 ascension levels.
 LevelExperienceRampOverrides=(
 ExperiencePointsForLevel[0]=1,
 ExperiencePointsForLevel[1]=5,
 ...
 ExperiencePointsForLevel[64]=1000)
The second example (when placed in the configuration file after the first) specifies 35 tame creature levels.
 LevelExperienceRampOverrides=(
 ExperiencePointsForLevel[0]=1,
 ExperiencePointsForLevel[1]=5,
 ...
 ExperiencePointsForLevel[34]=1000)
''The examples provided here are split into multiple lines for space considerations. In the configuration file, each entry must be placed entirely on a single line and without spaces.''
</div>

====Stats related====
<div style="border: 1px solid silver; padding: 0 1em;">
=====Attribute index table=====
This table shows relationships between each stats attribute and its coded index
{| class="wikitable"
|-
! Index !! Stats
|-
|0 || {{ItemLink|Health}}
|-
|1 || {{ItemLink|Stamina}} / {{ItemLink|Charge Capacity}}
|-
|2 || {{ItemLink|Torpidity}}
|-
|3 || {{ItemLink|Oxygen}} / {{ItemLink|Charge Regeneration}}
|-
|4 || {{ItemLink|Food}}
|-
|5 || {{ItemLink|Water}}
|-
|6 || Temperature
|-
|7 || {{ItemLink|Weight}}
|-
|8 || {{ItemLink|Melee Damage}} / {{ItemLink|Charge Emission Range}}
|-
|9 || {{ItemLink|Movement Speed}} / {{ItemLink|Maewing}}'s Nursing effectiveness
|-
|10 || {{ItemLink|Fortitude}}
|-
|11 || {{ItemLink|Crafting Speed}}
|}

=====MutagenLevelBoost=====
<code>MutagenLevelBoost['''<Stat_ID>''']='''<integer>'''</code>
{| class="wikitable"
|-
|<code>'''Stat_ID'''</code> || Index of the attribute to override. See the [[#Attribute_index_table|Attribute index]] table.
|-
|<code>'''integer'''</code> || Level points. Default values: 5, 5, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0
|}
Number of levels {{ItemLink|Mutagen}} adds to tames with wild ancestry.

The example provided doubles the amounts of level points a mutagen adds on Health and Damage stats but removes the extra level gains into Stamina and Weight:
 MutagenLevelBoost[0]=10 
 MutagenLevelBoost[1]=0
 MutagenLevelBoost[7]=0
 MutagenLevelBoost[8]=10

=====MutagenLevelBoost_Bred=====
<code>MutagenLevelBoost_Bred['''<Stat_ID>''']='''<integer>'''</code>
{| class="wikitable"
|-
|<code>'''Stat_ID'''</code> || Index of the attribute to override. See the [[#Attribute_index_table|Attribute index]] table.
|-
|<code>'''integer'''</code> || Level points. Default values: 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0
|}
Like <code>MutagenLevelBoost</code>, but for bred tame creatures.

The example provided doubles the amounts of level points a mutagen adds on Health and Damage stats but removes the extra level gains into Stamina and Weight:
 MutagenLevelBoost_Bred[0]=2 
 MutagenLevelBoost_Bred[1]=0
 MutagenLevelBoost_Bred[7]=0
 MutagenLevelBoost_Bred[8]=2

=====PlayerBaseStatMultipliers=====
 PlayerBaseStatMultipliers[<attribute>]='''<multiplier>'''
Arguments:
{| class="wikitable"
|<code>attribute</code> || ''integer'' || See the [[#Attribute_index_table|Attribute index]] table.
|-
|<code>multiplier</code> || ''float'' || Default: 1.0, see table below for the default value
|}
Changes the base stats of a player by multiplying with the default value. Meaning the start stats of a new spawned character.

Default values:
{| class="wikitable"
|-
!Attribute!!default!!output
|-
| 0 Health || 1.0 || 100.0
|-
|1 Stamina || 1.0 || 100.0
|-
|2 Torpidity || 1.0 || 200.0 (you still become unconscious at 50 regardless of total amount)
|-
|3 Oxygen || 1.0 || 100.0
|-
|4 Food || 1.0 || 100.0
|-
| 5 Water || 1.0 || 100.0
|-
|6 Temperature || 0.0 || 0.0 (Unused stat)
|-
|7 Weight || 1.0 || 100.0
|-
|8 MeleeDamageMultiplier || 0.0 || 100% (Base cannot be increased)
|-
|9 SpeedMultiplier || 0.0 || 100% (Base cannot be increased)
|-
| 10 TemperatureFortitude || 0.0 || 0 (Base cannot be increased)
|-
|11 CraftingSpeedMultiplier || 0.0 || 100% (Base cannot be increased)
|}

=====PerLevelStatsMultiplier=====
 PerLevelStatsMultiplier_DinoTamed<type>[<attribute>]='''<multiplier>'''
 PerLevelStatsMultiplier_DinoWild[<attribute>]='''<multiplier>'''
 PerLevelStatsMultiplier_Player[<attribute>]='''<multiplier>'''

Arguments:
{| class="wikitable"
|<code>type</code> || ''text'' ||
no type given: Multiplier applied for each tamed level-up point.<br/>
_Add: Multiplier immediately added for tamed creature.<br/>
_Affinity: Multiplier applied dependant on affinity.
|-
|<code>attribute</code> || ''integer'' || See the [[#Attribute_index_table|Attribute index]] table.
|-
|<code>multiplier</code> || ''float'' || Default: 1.0 or see table below
|}
Allows changing the amount of stats gained for each level:

<code>PerLevelStatsMultiplier_Player</code> changes the amount for players.

<code>PerLevelStatsMultiplier_DinoTamed</code> changes the amount for tamed creatures.

<code>PerLevelStatsMultiplier_DinoWild</code> changes the amount for wild creatures.

To nearly disable gaining stats use 0.01 because setting the value to 0 makes it default to 1.0.

Default Values:
{| class="wikitable"
|-
! Attribute!!Wild!!Tamed !!Tamed_Add!!Tamed_Affinity
|-
|0 Health || 1 || 0.2 || 0.14 || 0.44
|-
|1 Stamina || 1 || 1 || 1 || 1
|-
|2 Torpidity || 1 || 1 || 1 || 1
|-
|3 Oxygen || 1 || 1 || 1 || 1
|-
|4 Food || 1 || 1 || 1 || 1
|-
|7 Weight || 1 || 1 || 1 || 1
|-
| 8 Damage || 1 || 0.17 || 0.14 || 0.44
|-
| 9 Speed || 1 || 1 || 1 || 1
|}
Example how to double the weight-increase per level for players:
 PerLevelStatsMultiplier_Player[7]=2.0
Example for different types effecting Health of a tamed creatures: 
 PerLevelStatsMultiplier_DinoTamed[0]=1.0
 PerLevelStatsMultiplier_DinoTamed_Add[0]=1.0
 PerLevelStatsMultiplier_DinoTamed_Affinity[0]=1.0
</div>

====Single Player Settings====
<div style="border: 1px solid silver; padding: 0 1em;">
If <code>bUseSingleplayerSettings=True</code> than the following options are applied additionally to the configured (or default) values:
{| class="wikitable"
|-
! Option !! Base value !! Additional multiplier
|-
| <code>BabyCuddleIntervalMultiplier</code> || ''ini'' || x 0.17
|-
| <code>BabyMatureSpeedMultiplier</code> || ''ini'' || x 35.0
|-
| <code>AllowRaidDinoFeeding</code> || True || ''N/A''
|-
| <code>bAllowUnlimitedRespecs</code> || True || ''N/A''
|-
| <code>CropGrowthSpeedMultiplier</code> || ''ini'' || x 4.0
|-
| <code>EggHatchSpeedMultiplier</code> || ''ini'' || x 9.0
|-
| <code>HairGrowthSpeedMultiplier</code> || ''ini'' || x 0.69999999
|-
| <code>MatingIntervalMultiplier</code> || ''ini'' || x 0.15000001
|-
| <code>PerLevelStatsMultiplier_DinoTamed_Add[0]</code> || ''ini'' || x 3.5714285
|-
| <code>PerLevelStatsMultiplier_DinoTamed_Add[8]</code> || ''ini'' || x 3.5714285
|-
| <code>PerLevelStatsMultiplier_DinoTamed_Affinity[0]</code> || ''ini'' || x 2.2727273
|-
| <code>PerLevelStatsMultiplier_DinoTamed_Affinity[8]</code> || ''ini'' || x 2.2727273
|-
| <code>PerLevelStatsMultiplier_DinoTamed[0]</code> || ''ini'' || x 2.125
|-
| <code>PerLevelStatsMultiplier_DinoTamed[8]</code> || ''ini'' || x 2.3529413
|-
| <code>RaidDinoCharacterFoodDrainMultiplier</code> || ''ini'' || x 0.40000001
|-
| <code>TamingSpeedMultiplier</code> || ''ini'' || x 2.5
|-
| <code>UseCorpseLifeSpanMultiplier</code> || 6.0 || ''N/A''
|}
Where ''ini'' means the default or already set option is scaled by the value in the ''Additional multiplier'' column, e.g.:
:If <code>BabyMatureSpeedMultiplier</code> has not been set or left at default value of 1.0, with single player settings the final value would be 35.0.
:If <code>BabyMatureSpeedMultiplier</code> has been set at 2.0 instead, with single player settings the final value would be 70.0.
Finally, the [[Tek ATV]] engram will be available.
</div>

==Ark IDs==
An ID is necessary to identify any user conencting to a server whether it is for [[Server configuration#Administrator_whitelisting|Administrator Whitelisting]],[[Server configuration#Player_Whitelisting|Player Whitelisting]], or some other purpose.
<tabber>
|-|Ark: Survival Ascended=
For {{gamelink|sa}}, the game uses Epic Online Services (EOS). The following methods can be used to acquire the necessary Ark ID (a 32 character alphanumeric string, aka EOS ID):
* Start a Single Player map and after creating a survivor, open and expand the console (on PC double tap on {{Key|Tab}}), and type the command <code>whoami</code>. Record the OSS Id, the 32 character alphanumeric string.
* While connected to a Dedicated Server after creating a character, open and expand the console (on PC double tap on {{Key|Tab}}), and type in the command <code>whoami</code>. Record the OSS Id, the 32 character alphanumeric string. Also, after [[Console_commands#EnableCheats|enabling cheats]], it is also possible to use the command <code>cheat ListPlayers</code> to display all currently connected players and their Ark IDs.
* In the log of the running server, when a connection attempt is made, the user Ark ID is listed as the <code>UniqueNetId</code>.
|-|Ark: Survival Evolved=
For  {{gamelink|se}}, the game uses the Steam API. The following methods can be used to acquire the necessary Ark ID (a 17 digits string for Steam players, aka the Steam ID, and a 19 digits string for Epic Game Store players):
* Start a Local Host map and after creating a Survivor, open the console and expand it with two taps on {{Key|Tab}}, type the command <code>whoami</code> or <code>ListPlayers</code>.
* While connected to a Dedicated Server after creating a character, open the console and expand it with two taps on {{Key|Tab}},  type the command <code>whoami</code>. After [[Console_commands#EnableCheats|enabling cheats]], it is also possible to use the command <code>cheat ListPlayers</code> to display all currently connected players and their ARK IDs.
* Steam players only can also use SteamDB's Calculator to determine their Steam ID:
** Open Steam and click on your profile name in the upper bar, to display the the Steam account profile page.
** Click the address bar just below that to copy the profile URL.
** Open [https://steamdb.info/calculator/ SteamDB's Calculator] and paste the profile URL in into the field provided and hit enter. In the table below the currency amounts, it will list the <code>SteamID</code>.
'''Note''': <code>whoami</code> command will display two values, the Online SubSystem ID (OSS Id) and the ARK ID: Epic Game Store players need to take note of the ARK ID only (for Steam players it is literally the same value).
</tabber>

==Administrator Whitelisting==
First, start by getting necessary [[Server configuration#Ark_IDs|Ark IDs]].
<tabber>
|-|Ark: Survival Ascended=
In {{gamelink|sa}} players can be whitelisted as administrators on the server via their EOS Id (the 32 character alphanumeric string). These players can use cheat commands on the server automatically, as if they had authenticated themselves via the <code>enablecheats</code> command.
| * To whitelist administrators, create the file <code>ShooterGame/Saved/AllowedCheaterAccountIDs.txt</code>. In the file, list each player's [[Server configuration#Ark_IDs | Ark ID]], one per line. |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |Ark: Survival Evolved=
In {{gamelink|se}} players can be whitelisted as administrators on the server via their Steam Id (the 17 character numeric string). These players can use cheat commands on the server automatically, as if they had authenticated themselves via the <code>enablecheats</code> command.
* To whitelist administrators, create the file <code>ShooterGame/Saved/AllowedCheaterSteamIDs.txt</code>. In the file, list each player's [[Server configuration#Ark_IDs|Ark IDs]], one per line.
</tabber>

'''Note:''' When this method is used, it is not necessary to specify a server administrator password. A password can still be specified, and can be used by players not on the whitelist to gain administrator privileges, but the server will function without it and will still automatically grant privileges to whitelisted administrators.

==Player Whitelisting==
First, start by getting necessary [[Server configuration#Ark_IDs|Ark IDs]].

A Player Whitelist can come in two forms:
* A Whitelist which allows players to connect to the server when <code>-exclusivejoin</code> is added to the command line (confirmed working in both {{Gamelink|sa}} and {{Gamelink|se}}), or <code>UseExclusiveList=true</code> is present in <code>GameUserSettings.ini</code> under <code>[ServerSettings]</code>, or attached to the command line parameters with the <code>?</code> syntax (note: not documented by WildCard and has not been confirmed to work in {{Gamelink|sa}}). Players that are not on this whitelist will not be able to join the server. The regular server password can be discarded if this whitelist is in use. If the server password is left in place, whitelisted players will still need to enter it. This whitelist uses the file <code>/ShooterGame/Binaries/<PLATFORM>/PlayersExclusiveJoinList.txt</code>. To use the whitelist, insert the necessary [[Server configuration#Ark_IDs|Ark IDs]] into this file, one per line and use one of the three methods above to turn on the whitelist.
* A Whitelist which allows players to bypass the maximum number of players on the server. Players not on this whitelist are still subject to the maximum player count. This whitelist uses the file <code>/ShooterGame/Binaries/<PLATFORM>/PlayersJoinNoCheckList.txt</code>. To use the whitelist, insert the necessary [[Server configuration#Ark_IDs|Ark IDs]] into this file, one per line.
'''Note:''' In the file paths, <code><PLATFORM></code> is either <code>Linux</code> or <code>Win64</code>, depending on host operating system. 

For any ID addition or removal in the file, the server needs to restart. However, server admins can use the command <code>Cheat AllowPlayerToJoinNoCheck <ARK_ID></code> in the in-game console to add a new players to <code>PlayersJoinNoCheckList.txt</code> without restarting the server.

==Cross-ARK Data Transfer==
In '''Officials''', survivors can 'upload' or 'transfer' from any server at any Supply Crate, Terminal, Obelisk and {{ItemLink|Tek Transmitter}}. Survivors can then be 'downloaded' onto any Official Server of that same game mode.

For '''Unofficial Servers''', to allow dynamic Cross-ARK Travel it is needed to run at least two Servers having the same cluster ID and cluster directory, launched with similar command lines:

<code>ShooterGameServer TheIsland?SessionName=MySession1 -clusterid=<CLUSTER_NAME> -ClusterDirOverride="<PATH>"</code>

<code>ShooterGameServer ScorchedEarth_P?SessionName=MySession2 -clusterid=<CLUSTER_NAME> -ClusterDirOverride="<PATH>"</code>

Terminals will than list all the servers with the same <code>clusterid</code>. The Survivor will then be automatically transferred onto the server selected.

The <code>clusterid</code> must be the same between the servers, otherwise they would not be allowed to see each other.

It is also needed to specify the cross-server storage location with the <code>-ClusterDirOverride=<PATH></code> command line option, e.g.:
:On Linux: <code>-ClusterDirOverride="/MyARKClusterStorage"</code>
:On Windows: <code>-ClusterDirOverride="C:\MyARKClusterStorage"</code>
Otherwise, the server would default to <code>ShooterGame\Saved\clusters</code> preventing each server running at the same time and from different path to see each other.<ref name="How_Cross-ARK_Data_Transfer_works">Jeremy Stieglitz (September 2, 2016) [https://survivetheark.com/index.php?/forums/topic/85463-scorched-earth-technicaldetail-faq-ongoing "How Cross-ARK Data Transfer works"]. ''Scorched Earth Technical/Detail FAQ! (Ongoing)'' </ref>

===ARK data settings===
It is possible to control what uploads/downloads are allowed in tributes with the following options (command line or ''GameUserSettngns.ini''):
:<code>noTributeDownloads=<boolean></code>
:<code>PreventDownloadSurvivors=<boolean></code>
:<code>PreventDownloadItems=<boolean></code>
:<code>PreventDownloadDinos=<boolean></code>
:<code>PreventUploadSurvivors=<boolean></code>
:<code>PreventUploadItems=<boolean></code>
:<code>PreventUploadDinos=<boolean></code>
Tames re-upload timer can be controlled with the following option:
:<code>MinimumDinoReuploadInterval=<float></code>
It is also possible to change the expiration timer for uploaded tames, items (cryopod included) and survivors.
:<code>TributeCharacterExpirationSeconds=<integer></code>
:<code>TributeDinoExpirationSeconds=<integer></code>
:<code>TributeItemExpirationSeconds=<integer></code>
Every access to ARK Data will take in account each server settings, and will eventually delete the data if this has expired according to the server settings the data is accessed from. Therefore, it is suggested to keep the same expiration times on every server of the cluster.

Another care to take in consideration is to not use extremely high values for such settings, since these are summed to upload time in Unix Epoch time format and may results in an overflow.

Finally, if <code>TributeCharacterExpirationSeconds</code> is not set, uploaded survivors will have no expiration timer (the same result is achieved using a value equal or less than 0). This is not true for items and creature expiration settings.

===More About Cluster Files and Running Multiple Servers===
The cluster files will always be saved under the <code><ClusterDirOverride>/<clusterid></code> directory. This implies a few things:
* For a specific player, downloading (thus, also, transferring) a Survivor onto a server always overwrites whatever Survivor data was previously there. Inversely, there can only one Survivor 'Uploaded'.
* If two servers share the same <code>clusterid</code> but they do not share the same cluster directory (different <code><ClusterDirOverride></code>), players will still see on transfer list such servers but they will not be able to download (and thus also transfer) anything uploaded from a server to another.
* If you select a <code>clusterid</code> that is common, there is a high chance that your server list eventually becomes "polluted" by other servers that are configured to use the same <code>clusterid</code> which can confuse players since the data is not shared with those other servers. Given this issue, it's suggested to not use a common <code>clusterid</code> name as other unconnected servers could be listed while the players are using the transfers.

===Advanced cluster managements clarifications===
The following considerations are left for experts administrators only and require proper testing before putting such setups in production:

It is possible to set a shared location between different physical or virtual machines to use a common <code>-clusterid</code>. If over a network, it might slightly increase the chance of failing transfers given network issues, which typically ends with the player losing the entire data. For such more advanced setup, it is usually suggested to use a third-party applications to dynamically handle cluster and player data backups.

To spare some drive space, increase cache hits, allow better resource sharing and installation/update time, it is possible to use the same installation to run multiple servers. This can also be combined with <code>?AltSaveDirectoryName</code>.

*Note in this way they will also share the same <code>INI</code> settings because they load the same files unless using complicated file linking tricks. To still have them differ it is possible to set most of <code>GameUserSettings.ini</code> options in the command line using the <code>?</code> syntax. You could also use a dynamically generated <code>INI</code> file to change some of the <code>Game.ini</code> settings. Manual editing of <code>INI</code> files is usually not recommended. The server no longer loads any <code>INI</code> data once the <code>PID</code> of the process is displayed in the output console. It is thus needed to be careful to load the server with proper intervals as not to use the wrong INI options from a previous server launch. Also, it must be remembered that, when quitting the server normally, the process overwrites some content of the INI files with its current settings.
*The server deletes old save game files using the age of the file as long as they are named <code>.ark</code>. This means that regardless of the map you are using, the savegame files will eventually get deleted from the directory. If you stop a map, no more autosaves will be made so leaving any <code>.ark</code> files there will get them eventually deleted. Note that if this happens, you can at least still use your <code>.bak</code> backup files when the server was last started, which do not get deleted, given you don't have another backup solution running on the server.

<code>-clusterid</code> / <code>-NoTransferFromFiltering</code>:
*If <code>-clusterid</code> is '''not set''' and <code>-NoTransferFromFiltering</code> is '''not set''': it is still possible for players to use the 'download'/'upload' buttons on terminals. Tribute data would then be stored in players' own machines, the same way as their single player data. This would also allow them to download it into another server with such setup.
*If <code>-clusterid</code> is '''not set''' and <code>-NoTransferFromFiltering</code> is '''also set''': player transfers from single player is disabled. Given the Tribute Data is saved in the <code>SavedArks</code> directory, it is still possible to create a "fake cluster" by sharing the <code>SavedArks</code> directory between servers, allowing players to use the 'download'/'upload' buttons on terminals. Finally, it's required by players to known what servers are linked.

It is not possible to mix and match these types of servers between one another, as the Tribute Data format will not be the same: with the <code>-clusterid</code>, the Tribute Data is saved in <code>-clusterid</code> directory and if <code>-NoTransferFromFiltering</code> is not set, it will save on the players' own machines.

For {{ItemLink|ARK: Survival Ascended|ASA}}, each map is now in it's own directory, even when using the <code>?AltSaveDirectoryName</code>. To reproduce this functionality, you will need to link/junction each of the individual map directories into a single one. Note a special consideration about [[Club ARK]] as it should not be linked/junctioned. The mod actually creates a new user specifically for Club ARK and is apparently only designed to update the shared cluster data instead. Thus if using a link/junction for the map, it will delete the existing user data.

== DynamicConfig ==
The '''DynamicConfig''' allows to change some server settings without the need of restarting the server. The file must be hosted on a web server (e.g.: IIS on Windows, Apache etc...) using the '''HTTP''' protocol (HTTPS is not supported) and the ''.ini'' file extension must be exposed as ''mime'' type as ''text/plain''. It cannot be directly linked using a system path.

The configuration is checked every time the world is (auto)saved, unless an admin force an immediate update using the <code>ForceUpdateDynamicConfig</code> command.

Once <code>-UseDynamicConfig</code> and <code>CustomDynamicConfigUrl</code> options are set, the following settings are available to dynamically configure: 
{| class="wikitable config-table"
|-
! data-filter-name="Supported by ARK: Survival Ascended" | {{ItemLink|ARK: Survival Ascended|ASA}}
! data-filter-name="Supported by ARK: Survival Evolved" | {{ItemLink|ARK: Survival Evolved|ASE}}
! Variable
! Description
! Since patch
{{Server config variable
| name = ActiveEventColors
| type = string
| default = N/A
| cli = skip
| inASA = No
| inASE = Yes
| info = This will activate the associated colours of an event:
<table class="wikitable">
<tr><th>eventname</th><th>Description</th></tr>
<tr><td><code>easter</code></td><td>[[ARK: Eggcellent Adventure 7]]</td></tr>
<tr><td><code>FearEvolved</code></td><td>[[ARK: Fear Evolved 6]]</td></tr>
<tr><td><code>PAX</code></td><td>[[ARK: PAX Party]]'', Re-added in [[357.12]]''</td></tr>
<tr><td><code>Summer</code></td><td>[[ARK: Summer EVO]]</td></tr>
<tr><td><code>TurkeyTrial</code></td><td>[[ARK: Turkey Trial 6]]</td></tr>
<tr><td><code>vday</code></td><td>[[Valentine's EVO Event]]</td></tr>
<tr><td><code>WinterWonderland</code></td><td>[[ARK: Winter Wonderland 7]]</td></tr>
<tr><td><code>custom</code></td><td>If specified, you have to provide the colours in <code>DynamicColorset</code>.</td></tr>
</table>
| version = [[356.11]]
}}
{{Server config variable
| name = BabyCuddleIntervalMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = See same option in ''Game.ini''.
| version = [[316.18]]
}}
{{Server config variable
| name = BabyFoodConsumptionSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''Game.ini''.
}}
{{Server config variable
| name = BabyImprintAmountMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = See same option in ''Game.ini''.
| version = [[316.18]]
}}
{{Server config variable
| name = BabyMatureSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = See same option in ''Game.ini''.
| version = [[307.2]]
}}
{{Server config variable
| name = bDisableDinoDecayPvE
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Same as <code>DisableDinoDecayPvE</code> in ''GameUserSettings.ini''.
}}
{{Server config variable
| name = bDisableStructureDecayPvE
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Same as <code>DisableStructureDecayPvE</code> in ''GameUserSettings.ini''.
}}
{{Server config variable
| name = bPvPDinoDecay
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Same as <code>PvPDinoDecay</code> in ''GameUserSettings.ini''.
}}
{{Server config variable
| name = bPvPStructureDecay
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Same as <code>PvPStructureDecay</code>as in ''GameUserSettings.ini''.
}}
{{Server config variable
| status = deprecated
| name = bUseAlarmNotifications
| type = boolean
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Toggles web alarm notifications. See [[Web_Notifications|Web Notifications]] for more details. Undocumented by Wildcard.
}}
{{Server config variable
| name = CropGrowthSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''Game.ini''.
}}
{{Server config variable
| name = CustomRecipeEffectivenessMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''Game.ini''.
}}
{{Server config variable
| name = DinoCharacterFoodDrainMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''GameUserSettings.ini''.
}}
{{Server config variable
| status = deprecated
| name = DisableTimestampVerification
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Undocumented.
}}
{{Server config variable
| name = DisableWorldBuffs
| type = string
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Disabling specific world buffs on [[Genesis: Part 2]] that would be increased beyond normal values on those Arks with other multiplier changes.<br/>'''Note:''' buffs must be specified in a single row comma separated list without spaces.<br/>Known buff identifiers: <code>MATINGINTERVAL_DOWN_HARD</code>, <code>MATINGINTERVAL_DOWN_MEDIUM</code>, <code>MATINGINTERVAL_DOWN_EASY</code>, <code>BABYMATURE_BOON_EASY</code>, <code>BABYMATURE_BOON_MEDIUM</code>, <code>BABYMATURE_BOON_HARD</code>.
}}
{{Server config variable
| status = deprecated
| name = DynamicUndermeshRegions
| type = string
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Forces dynamic undermesh volume update. Undocumented.
}}
{{Server config variable
| name = DynamicColorset
| type = string
| default = N/A
| cli = skip
| inASA = Yes
| inASE = Yes
| info = Comma-separated list of [[Color_IDs#List_of_colors|colour names]], names must be **exact**, accounting for special characters and spaces. <code>DynamicColorset</code> will only be used if using <code>ActiveEventColors=custom</code>.
| version = [[356.11]]
}}
{{Server config variable
| name = DynamicColorsetChanceOverride
| type = float
| default = N/A
| cli = skip
| inASA = Yes
| inASE = Unknown
| info = Sets the chance of the Dynamic Colorset being applied from 0.0 to 1.0 (0-100%). With nothing here the chance is 0.25 (25%). 
| version = [[56.13]]
}}
{{Server config variable
| name = EggHatchSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = See same option in ''Game.ini''.
| version = [[307.2]]
}}
{{Server config variable
| status = deprecated
| name = EnableFullDump
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Forces full server memory dump on server crash. Undocumented by Wildcard.
}}
{{Server config variable
| name = EnableWorldBuffScaling
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Same as <code>bEnableWorldBuffScaling</code> in ''Game.ini''.
}}
{{Server config variable
| status = deprecated
| name = GMaxFlameThrowerServerTicksPerFrame
| type = integer
| default = 5
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Controls the tick rate of [[Flamethrower]] per server tick, i.e.: the shooting speed of Flamethrower per server tick instead of per fixed time. Higher values will make the Flamethrower using more ammos per server tick. Using too high values may result in server performance issues. Undocumented by Wildcard.
}}
{{Server config variable
| name = GlobalSpoilingTimeMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''Game.ini''.
}}
{{Server config variable
| status = deprecated
| name = GUseServerNetSpeedCheck
| type = boolean
| default = False
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = It should avoid players to accumulate too much movements data per server tick, discarding the last if those are too many. Can be turned on also with <code>-UseServerNetSpeedCheck</code> command line option. Enabled on official clusters. May be useful for servers having too many players and stressed on CPU side. Undocumented by Wildcard.
}}
{{Server config variable
| name = HarvestAmountMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = See same option in ''GameUserSettings.ini''.
| version = [[307.2]]
}}
{{Server config variable
| name = HexagonRewardMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = See same option in ''Game.ini''.
| version = [[316.18]]
}}
{{Server config variable
| name = MatingIntervalMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = See same option in ''Game.ini''.
| version = [[307.2]]
}}
{{Server config variable
| name = MatingSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''Game.ini''.
}}
{{Server config variable
| name = NPCReplacements
| type = string
| default = N/A
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = Globally replaces specific creatures with another using class names. See [[#NPCReplacements|Creature Spawn related]] section.
}}
{{Server config variable
| name = PvEDinoDecayPeriodMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''GameUserSettings.ini''.
}}
{{Server config variable
| name = PvEStructureDecayPeriodMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''GameUserSettings.ini''.
}}
{{Server config variable
| name = StructureDamageMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''GameUserSettings.ini''.
}}
{{Server config variable
| name = TamingSpeedMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''GameUserSettings.ini''.
| version = [[307.2]]
}}
{{Server config variable
| name = TributeDinoExpirationSeconds
| type = integer
| default = 86400
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''GameUserSettings.ini''.
}}
{{Server config variable
| name = TributeItemExpirationSeconds
| type = integer
| default = 86400
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''GameUserSettings.ini''.
}}
{{Server config variable
| name = XPMultiplier
| type = float
| default = 1.0
| cli = skip
| inASA = Yes
| inASE = Yes
| info = See same option in ''Game.ini''.
| version = [[307.2]]
}}
{{Server config variable
| name = WorldBuffScalingEfficacy
| type = float
| default = 1.0
| cli = skip
| inASA = Unknown
| inASE = Yes
| info = See same option in ''Game.ini''.
}}
|}
</div>

== The Survival of the Fittest ==
[[The Survival of the Fittest]] currently supports the below server options officially, aside from those settings and options from above should work if applicable to the game mode.

=== Command line ===
{| class="wikitable"
! Option!!Effect!!Since Patch
|-
|<code>-gameplaylogging</code>||Survival of the Fittest only servers can launch with this option to output a dated log file to <code>\Saved</code> folder, which will contain a timestamped kill & winners log listing steam id, steam name, character name, etc. Handy for automatic tournament records.||[[205.0]]
|}

=== GameUserSettings.ini ===
put the following Settings below the <code>[TSOTF]</code> tag or add the tag if not present in the file
{| class="wikitable"
|-
! Option !! Value !! Effect
|-
|<code>EnableDeathTeamSpectator='''<boolean>'''</code> || false || If <code>True</code>, enables player spectator mode after team death. Can be set in command line to with the <code>?</code> syntax.
|-
|<code>MaxPlayersPerTribe='''<integer>'''</code> || 4 || Sets the maximum players allowed per tribe.
|-
|<code>MinPlayersToQuickStartMatch='''<integer>'''</code> || 32 || Controls how many players you need to start the quick timer.
|-
|<code>MinTribesToStartMatch='''<integer>'''</code> || 5 || Controls how many tribes you need to start the normal timer.
|-
|<code>PreventBosses='''<boolean>'''</code> || 0 || Set to 1 to disable bosses .
|-
|<code>PreventDinoSupplyCrates='''<boolean>'''</code> || 0 || Set to 1 to disable the Dino Crate POI system. 
|-
|<code>PreventRetameBoss='''<boolean>'''</code> || 1 || Set to 0 to enable the retaming of Bosses.
|-
|<code>PreventSpectator='''<boolean>'''</code> || 0 || set to 1 disables spectator as an eliminated player.
|-
|<code>PreventTribes='''<boolean>'''</code> || 0 || Set to 1 disables tribes.
|-
|<code>TimeToQuickStartMatch='''<integer>'''</code> || 120 || Quick timer to auto start the match in seconds.
|-
|<code>TimeToStartMatch='''<integer>'''</code> || 600 || Normal timer to auto start the match in seconds.
|}

put the following Settings below the <code>[ServerSettings]</code> tag or add the tag if not present in the file
{| class="wikitable"
|-
! Option !! Value !! Effect
|-
|<code>AutoCreateTribes='''<boolean>'''</code> || true || Forces players into a tribe automatically (good for prevent tribe servers (1 player)).
|-
|<code>BadWordListURL='''<string>'''</code> || ''N/A'' || Adds the url to hosting your own bad words list .
|-
|<code>BadWordWhiteListURL='''<string>'''</code> || ''N/A'' || Adds the url to hosting your own good words list.
|-
|<code>bEnablePlayerMoveThroughAlly='''<boolean>'''</code> || false || If set to true allows players to walk through ally dinos.
|-
|<code>bEnablePlayerMoveThroughSleeping='''<boolean>'''</code> || false || If set to true allows players to walk through unconscious dinos.
|-
|<code>bFilterCharacterNames='''<boolean>'''</code> || true || Filters out character names based on the bad words/good words list.
|-
|<code>bFilterChat='''<boolean>'''</code> || true || Filters out character names based on the bad word/good words list .
|-
|<code>bFilterTribeNames='''<boolean>'''</code> || true || Filters out tribe names based on the badwords/goodwords list. 
|-
|<code>RiderDinoCollision='''<boolean>'''</code> || false || This is a new config that has been added to the game mode to minimize some pain of creatures stacking together when having follow armies, therefore making them easier to kill -- however this also means that they won't follow as easily due to having collision with their allies -- hence why we're tinkering to find the best/optimal way. For now we are not planning on enabling this on our official servers. If enabled the following occurs:
*Ridden and Possessed creatures can go through corpses, unconscious creatures, and ally tames
*AI driven creatures (so commander mode, unridden, and unpossessed) cannot go through ally tames
*AI driven creatures (so commander mode, unridden, and unpossessed) can go through sleeping creatures and corpses.
|}

=== Game.ini ===
put the following Settings below the <code>[/script/shootergame.shootergamemode]</code> tag or add the tag if not present in the file
{| class="wikitable"
|-
! Option!!Arguments!!Effects
|-
|<code>DinoClassDamageMultipliers=(<br/>ClassName="'''<classname>'''",<br/>Multiplier='''<multiplier>'''<br/>)</code>
<code>TamedDinoClassDamageMultipliers=(<br/>ClassName="'''<classname>'''",<br/>Multiplier='''<multiplier>'''<br/>)</code>
|classname = ''string''
multiplier = ''float'' Default 1.0
|Multiplies damage dealt of specific dinosaurs via classname. Higher values increase the damage dealt.
Dino classnames can be found on the [[Creature IDs]] page.
The examples provided here are split into multiple lines for space considerations. In the configuration file, an entry must be placed on a single line. Multiple <code>DinoClassDamageMultipliers</code> and <code>TamedDinoClassDamageMultipliers</code> entries can be specified in the file, but <code>ClassName</code> values should not be repeated across multiple entries.

Examples:
 DinoClassDamageMultipliers=(
   ClassName="MegaRex_Character_BP_C",
   Multiplier=0.1
 )

 TamedDinoClassDamageMultipliers=(
   ClassName="Rex_Character_BP_C",
   Multiplier=10.0
 )
Line breaks and spaces are here for better readability of the example. Keep it as one line in your configuration file.
|-
|<code>DinoClassSpeedMultipliers=(ClassName="'''<classname>'''",Multiplier=<value>)</code>
|classname = ''string'' 
multiplier = ''float'' Default 1.0
|allows overriding a wild dino's speed.
Dino classnames can be found on the [[Creature IDs]] page.
|-
|<code>PreventDinoTameClassNames='''<classname>'''</code>
|classname = ''string''
|prevent dinos from being spawned E.g <code>PreventDinoTameClassNames=Yutyrannus_Character_BP_TSOTF_C </code>
Dino classnames can be found on the [[Creature IDs]] page.
|-
|<code>TamedDinoClassSpeedMultipliers=(ClassName="'''<classname>'''",Multiplier=<value>)</code>
|classname = ''string'' 
multiplier = ''float'' Default 1.0
|allows overriding a tamed dino's speed 
Dino classnames can be found on the [[Creature IDs]] page.
|-
|<code>TamedDinoClassStaminaMultipliers=(ClassName="'''<classname>'''",Multiplier=<value>)</code>
|classname = ''string'' 
multiplier = ''float'' Default 1.0
|allows overriding a tamed dino's stamina
Dino classnames can be found on the [[Creature IDs]] page.
|}

==References==
<references />

== External links ==
* {{gamelink|sa}} PC patch notes: [https://survivetheark.com/index.php?/forums/topic/708761-asa-pc-patch-notes-client-v3134-server-v3137-updated-11232023/]
* {{gamelink|se}} PC patch notes:
:: New Server Commands/Options for version [[242.0]]: [https://steamcommunity.com/app/346110/discussions/10/350532536100692514/]
:: Notes between [[1.556]] and [[171.1]]: [https://steamcommunity.com/app/346110/discussions/0/615086038671453707/]
:: Notes between [[171.1]] and [[171.3]]: [https://steamcommunity.com/app/346110/discussions/0/615086038674193312/]
:: Notes between [[171.31]] and [[171.8]]: [https://steamcommunity.com/app/346110/discussions/0/615086038683522620/]
:: Notes between [[172.3]] and [[172.5]]: [https://steamcommunity.com/app/346110/discussions/0/594820473977643250/]
:: Notes between [[173.0]] and [[287.103]] (excluding from 232.0 to 240.6): [https://survivetheark.com/index.php?/forums/topic/166421-archived-pc-patch-notes/]
:: Notes between [[232.0]] and [[278.73]]: [https://steamcommunity.com/app/346110/discussions/0/594820656447032287/]
:: Notes since [[285.104]]: [https://survivetheark.com/index.php?/forums/topic/388254-pc-patch-notes-client-v3573-server-v35718-updated-05052023/]

{{Nav technical}}

[[Category:Guides]]
[[Category:Servers]]

[[es:Configuración del servidor]]
[[fr:Configuration serveur]]
[[ja:サーバー構成]]
[[ru:Настройка сервера]]
[[th:การกำหนดค่าของเซิร์ฟเวอร์]]
{{MissingTranslations|de|it|pl|pt-br}}
