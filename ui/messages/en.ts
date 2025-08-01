export default {
  // Common
  common: {
    confirm: 'Confirm',
    cancel: 'Cancel',
    save: 'Save',
    delete: 'Delete',
    edit: 'Edit',
    add: 'Add',
    close: 'Close',
    loading: 'Loading...',
    error: 'Error',
    success: 'Success',
    warning: 'Warning',
    info: 'Info',
    yes: 'Yes',
    no: 'No',
    back: 'Back',
    next: 'Next',
    previous: 'Previous',
    submit: 'Submit',
    reset: 'Reset',
    search: 'Search',
    filter: 'Filter',
    sort: 'Sort',
    refresh: 'Refresh',
    copy: 'Copy',
    download: 'Download',
    upload: 'Upload',
    export: 'Export',
    import: 'Import',
    settings: 'Settings',
    help: 'Help',
    about: 'About',
    version: 'Version',
    language: 'Language',
    theme: 'Theme',
    dark: 'Dark',
    light: 'Light',
    auto: 'Auto'
  },

  // Navigation
  navigation: {
    home: 'Home',
    dashboard: 'Dashboard',
    servers: 'Server Management',
    players: 'Player Management',
    logs: 'Log Monitoring',
    settings: 'Settings',
    logout: 'Logout',
    welcome: 'Welcome',
    user: 'User'
  },

  // Authentication
  auth: {
    login: 'Login',
    logout: 'Logout',
    username: 'Username',
    password: 'Password',
    loginTitle: 'Ark Server Commander',
    loginSubtitle: 'Secure login to your management account',
    loginButton: 'Login',
    loginLoading: 'Logging in...',
    loginError: 'Login failed',
    loginSuccess: 'Login successful',
    logoutSuccess: 'Logged out successfully',
    initCheck: 'Checking system initialization status',
    initRequired: 'System initialization required',
    alreadyLoggedIn: 'You are already logged in',
    enterUsername: 'Please enter username',
    enterPassword: 'Please enter password',
    firstTimeTip: 'First time? The system will automatically guide you through initialization',
    secureLogin: 'Secure Login System',
    // Initialization related
    initTitle: 'System Initialization',
    initSubtitle: 'First time use, please set up administrator account',
    adminUsername: 'Administrator Username',
    enterAdminUsername: 'Please enter administrator username',
    confirmPassword: 'Confirm Password',
    enterConfirmPassword: 'Please enter password again',
    passwordMinLength: 'Password must be at least 6 characters',
    passwordMinLengthError: 'Password must be at least 6 characters',
    passwordMismatch: 'Passwords do not match',
    initButton: 'Initialize System',
    initLoading: 'Initializing...',
    initSuccess: 'System initialized successfully',
    initError: 'System initialization failed',
    initTip: 'Will automatically redirect to main page after initialization',
    initWizard: 'System Initialization Wizard'
  },

  // Home page
  home: {
    title: 'Welcome to ARK Server Commander',
    subtitle: 'You have successfully logged in and can start managing your ARK servers.',
    systemInfo: 'System Information',
    username: 'Username',
    userID: 'User ID',
    imageManagement: 'Image Management',
    features: 'Feature Modules',
    serverManagement: 'Server Management',
    serverManagementDesc: 'Add, configure and manage your ARK servers with one-click start, stop and monitoring',
    startManage: 'Start Managing',
    playerManagement: 'Player Management',
    playerManagementDesc: 'Manage server players, view online status and permission settings',
    logMonitoring: 'Log Monitoring',
    logMonitoringDesc: 'Real-time server log monitoring, view system status and performance metrics',
    comingSoon: 'Coming Soon',
    tip: 'Click the cards above to start managing your ARK servers'
  },

  // Server management
  servers: {
    title: 'Server Management',
    serverManagementDesc: 'Manage and monitor your ARK server instances',
    addServer: 'Add Server',
    editServer: 'Edit Server',
    deleteServer: 'Delete Server',
    serverName: 'Server Name',
    serverPort: 'Server Port',
    serverPath: 'Server Path',
    serverStatus: 'Server Status',
    serverActions: 'Actions',
    startServer: 'Start Server',
    stopServer: 'Stop Server',
    restartServer: 'Restart Server',
    viewLogs: 'View Logs',
    serverConfig: 'Server Configuration',
    gameIni: 'Game.ini Configuration',
    gameUserSettings: 'GameUserSettings.ini Configuration',
    serverArgs: 'Startup Arguments',
    running: 'Running',
    stopped: 'Stopped',
    starting: 'Starting',
    stopping: 'Stopping',
    error: 'Error',
    unknown: 'Unknown',
    confirmDelete: 'Are you sure you want to delete this server?',
    deleteWarning: 'This action cannot be undone',
    serverAdded: 'Server added successfully',
    serverUpdated: 'Server updated successfully',
    serverDeleted: 'Server deleted successfully',
    serverStartSuccess: 'Server started successfully',
    serverStopSuccess: 'Server stopped successfully',
    serverRestartSuccess: 'Server restarted successfully',
    serverStartError: 'Failed to start server',
    serverStopError: 'Failed to stop server',
    serverRestartError: 'Failed to restart server',
    noServers: 'No servers yet',
    noServersDesc: 'Click "Add Server" to create your first ARK server',
    serverConfigSaved: 'Server configuration saved',
    serverConfigError: 'Failed to save server configuration',
    invalidPort: 'Invalid port number',
    invalidPath: 'Invalid server path',
    portInUse: 'Port is already in use',
    pathNotExists: 'Server path does not exist',
    serverNameRequired: 'Server name is required',
    serverPortRequired: 'Server port is required',
    serverPathRequired: 'Server path is required',
    imageStatus: 'Image Status',
    imageDownloading: 'Image Downloading',
    imageNotReady: 'Image Not Ready',
    imageDownloadingDesc: 'Downloading image, please wait before creating server',
    imageNotReadyDesc: 'Image not ready, cannot create server',
    // Docker images detailed translations
    dockerImages: {
      title: 'Image Download Status',
      overallStatus: 'Overall Status',
      imageReady: 'Images Ready',
      imageNotReady: 'Images Not Ready (Cannot Start Server)',
      imageMissingManualDownload: 'Images Missing, Please Download Manually',
      downloading: 'Downloading',
      ready: 'Ready',
      notReady: 'Not Ready',
      waitingDownload: 'Waiting Download',
      layerProgress: 'Layer Download Progress',
      totalImages: 'Total Images',
      downloadingCount: 'Downloading',
      refreshStatus: 'Refresh Image Status',
      manualDownload: 'Manual Download',
      checkUpdates: 'Check Updates',
      updateConfirm: 'Image Update Confirmation',
      imageInfo: 'Image Information',
      imageName: 'Image Name',
      affectedServers: 'Affected Servers',
      updateWarning: 'Update Risk Warning',
      warningDownloadTime: 'Image download may take a long time, please be patient',
      warningContainerRecreate: 'Container recreation will cause brief server downtime',
      warningDataSafety: 'Please ensure important data is backed up to avoid data loss',
      updateOptions: 'Update Options',
      updateImageOnly: 'Update Image Only',
      updateImageOnlyDesc: 'Only download new image, do not recreate containers. Manual container recreation required to use new image.',
      updateAndRecreate: 'Update Image and Recreate Containers',
      updateAndRecreateDesc: 'Download new image and automatically recreate all affected containers. Servers will be briefly offline.',
      confirmUpdate: 'Confirm Update',
      unknownSize: 'Unknown Size',
      // Image names
      arkServer: 'ARK Server',
      alpineSystem: 'Alpine System',
      // Layer information
      layerDetails: 'Layer Details',
      layers: 'Layers',
      // Layer status
      layerStatus: {
        pending: 'Pending',
        downloading: 'Downloading',
        extracting: 'Extracting',
        verifying: 'Verifying',
        complete: 'Complete'
      }
    },
    cannotDeleteRunning: 'Cannot delete running server, please stop it first',
    serverCreateSuccess: 'Server created successfully',
    serverUpdateSuccess: 'Server updated successfully',
    serverDeleteSuccess: 'Server deleted successfully',
    serverStartInProgress: 'Server starting...',
    serverStopInProgress: 'Server stopping...',
    copyToClipboard: 'Copied to clipboard',
    copyFailed: 'Copy failed, please copy manually',
    authenticationFailed: 'Authentication failed, please login again',
    getServerListFailed: 'Failed to get server list, please try again later',
    loadServerInfoFailed: 'Failed to load server info, please try again later',
    operationFailed: 'Operation failed, please try again later',
    deleteFailed: 'Delete failed, please try again later',
    startServerFailed: 'Failed to start server, please try again later',
    stopServerFailed: 'Failed to stop server, please try again later',
    imageStatusError: 'Failed to get image status',
    // Server card related
    card: {
      startServer: 'Start Server',
      stopServer: 'Stop Server',
      starting: 'Starting...',
      stopping: 'Stopping...',
      unknownStatus: 'Unknown Status',
      cannotStartImageNotReady: 'Image not ready, cannot start',
      rconInfo: 'RCON Info',
      rconConnectionInfo: 'RCON Connection Info',
      serverIdentifier: 'Server Identifier',
      rconPort: 'RCON Port',
      adminPassword: 'Admin Password',
      editServer: 'Edit Server',
      deleteServer: 'Delete Server',
      confirmDelete: 'Confirm Delete',
      confirmDeleteMessage: 'Are you sure you want to delete server "{identifier}"? This action cannot be undone.',
      status: 'Status',
      serverName: 'Server Name',
      clusterId: 'Cluster ID',
      map: 'Map',
      maxPlayers: 'Max Players',
      portConfig: 'Port Configuration',
      gamePort: 'Game Port',
      queryPort: 'Query Port',
      rconPortLabel: 'RCON Port',
      authInfo: 'Authentication Info',
      timeInfo: 'Time Information',
      createdAt: 'Created At',
      updatedAt: 'Updated At',
      serverId: 'Server ID',
      copy: 'Copy',
      close: 'Close',
      showPassword: 'Show Password',
      hidePassword: 'Hide Password'
    },
    // Server edit related
    edit: {
      title: 'Server Edit',
      createTitle: 'Add Server',
      editTitle: 'Edit Server',
      createServerDesc: 'Configure and create a new ARK server instance',
      basicParams: 'Basic Parameters',
      gameUserSettings: 'GameUserSettings.ini',
      gameIni: 'Game.ini',
      serverArgs: 'Startup Arguments (SERVER_ARGS)',
      serverIdentifier: 'Server Identifier',
      serverIdentifierRequired: 'Server Identifier *',
      serverIdentifierPlaceholder: 'Enter server identifier',
      serverName: 'Server Name',
      serverNamePlaceholder: 'Enter server name',
      serverNameDesc: 'Name displayed in game server list',
      clusterId: 'Cluster ID',
      clusterIdPlaceholder: 'Enter cluster ID (optional)',
      clusterIdDesc: 'For data sharing between cluster servers',
      gamePort: 'Game Port',
      gamePortRequired: 'Game Port *',
      gamePortPlaceholder: '7777',
      queryPort: 'Query Port',
      queryPortRequired: 'Query Port *',
      queryPortPlaceholder: '27015',
      rconPort: 'RCON Port',
      rconPortRequired: 'RCON Port *',
      rconPortPlaceholder: '32330',
      map: 'Map',
      mapPlaceholder: 'Select map',
      maxPlayers: 'Max Players',
      maxPlayersPlaceholder: '70',
      maxPlayersDesc: 'Maximum number of players (1-200)',
      modIds: 'Mod IDs',
      modIdsPlaceholder: 'Enter mod IDs, separated by commas (e.g., 123456,789012)',
      modIdsDesc: 'Steam Workshop mod IDs, multiple mods separated by commas',
      adminPassword: 'Admin Password',
      adminPasswordRequired: 'Admin Password *',
      adminPasswordPlaceholder: 'Enter admin password (also used as RCON password)',
      showPassword: 'Show Password',
      hidePassword: 'Hide Password',
      saveChanges: 'Save Changes',
      createServer: 'Create Server',
      saving: 'Saving...',
      preparing: 'Preparing...',
      loadingServerInfo: 'Loading server information...',
      closeConfirm: 'Are you sure you want to close? Unsaved data will be lost.',
      range: 'Range',
      // Map options
      maps: {
        TheIsland: 'The Island',
        TheCenter: 'The Center',
        ScorchedEarth_P: 'Scorched Earth',
        Aberration_P: 'Aberration',
        Extinction: 'Extinction',
        Valguero_P: 'Valguero',
        Genesis: 'Genesis',
        Genesis2: 'Genesis 2',
        CrystalIsles: 'Crystal Isles',
        LostIsland: 'Lost Island',
        Fjordur: 'Fjordur'
      },
      selectMapPlaceholder: 'Select map or enter map name',
      customMapPlaceholder: 'Enter custom map name',
      customMap: 'Custom Map',
      noMatchingMaps: 'No matching maps',
      mapId: 'Map ID'
    },
    // Editor related
    editor: {
      visualEdit: 'Visual Edit',
      textEdit: 'Text Edit',
      visualEditMode: 'Visual Edit',
      textEditMode: 'Text Edit',
      syncing: 'Syncing...',
      synced: 'Synced',
      pendingSync: 'Pending Sync',
      resetToDefault: 'Reset to Default',
      format: 'Format',
      content: 'Content',
      placeholder: 'Enter configuration content...',
      description: 'This file contains basic server settings such as port, password, max players, etc.',
      visualEditModeDesc: 'Visual Edit Mode',
      visualEditModeTip: 'Modify parameters through form controls. Hover over the icon next to parameter names to view detailed descriptions.',
      gameIniTextEditDesc: 'Directly edit Game.ini configuration file content.',
      showPassword: 'Show Password',
      hidePassword: 'Hide Password',
      enabled: 'Enabled',
      disabled: 'Disabled',
      parametersCount: ' parameters',
      defaultValue: 'Default',
      parseGameIniError: 'Failed to parse Game.ini text',
      syncVisualToTextError: 'Failed to sync visual config to text',
      range: 'Range',
      syncTip: {
        visual: 'Modify parameters in visual mode, configuration will be synced to text when switching to text mode.',
        text: 'Edit configuration file directly in text mode, content will be parsed when switching to visual mode.'
      }
    },
    // Server args editor related
    argsEditor: {
      title: 'Startup Arguments Configuration',
      switchParams: 'Switch Parameters',
      numberParams: 'Number Parameters',
      textParams: 'Text Parameters',
      selectParams: 'Select Parameters',
      range: 'Range',
      pleaseSelect: 'Please Select',
      customArgs: 'Custom Arguments',
      customArgsDesc: 'Add custom startup arguments that will be directly added to the startup command',
      addCustomArg: 'Add Custom Argument',
      removeCustomArg: 'Remove Custom Argument',
      customArgPlaceholder: 'Enter custom startup argument, e.g. -nosteam',
      enabled: 'Enabled',
      disabled: 'Disabled'
    },
    // Startup parameter categories
    paramCategories: {
      basic: 'Basic Settings',
      core: 'Core Game',
      dinos: 'Dinosaur Management',
      structures: 'Building System',
      pvp: 'PvP Settings',
      mechanics: 'Game Mechanics',
      transfer: 'Data Transfer',
      performance: 'Performance Optimization',
      graphics: 'Graphics Optimization',
      security: 'Security Settings',
      logging: 'Logging',
      mods: 'Mods & Platform',
      features: 'Game Features',
      maintenance: 'Maintenance Tools',
      advanced: 'Advanced Configuration',
      custom: 'Custom Parameters',
      // Add missing categories
      events: 'Event Settings',
      network: 'Network Settings',
      compatibility: 'Compatibility Settings',
      localization: 'Localization Settings',
      notifications: 'Notification Settings',
      voice: 'Voice Settings',
      gameplay: 'Gameplay Settings'
    },
    // Query parameters translation
    queryParams: {
      // Core game settings
      OverrideOfficialDifficulty: 'Override Official Difficulty Level',
      AllowThirdPersonPlayer: 'Enable Third Person View',
      ServerCrosshair: 'Disable Crosshair',
      ShowFloatingDamageText: 'Enable RPG Style Floating Damage Text',
      MapPlayerLocation: 'Show Player Precise Location on Map',
      DisableWeatherFog: 'Disable Fog Weather',
      // Dinosaur management
      MaxPersonalTamedDinos: 'Tamed Dinosaur Limit per Tribe',
      MaxTamedDinos: 'Server Total Tamed Dinosaur Limit',
      PersonalTamedDinosSaddleStructureCost: 'Platform Saddle Dinosaur Slot Cost',
      MinimumDinoReuploadInterval: 'Dinosaur Re-upload Cooldown (seconds)',
      ServerAutoForceRespawnWildDinosInterval: 'Force Wild Dinosaur Respawn Interval',
      AllowAnyoneBabyImprintCuddle: 'Allow Anyone to Care for Baby Dinosaurs',
      DisableImprintDinoBuff: 'Disable Dinosaur Imprint Player Stat Bonus',
      PreventMateBoost: 'Disable Dinosaur Mate Boost',
      AllowFlyerCarryPvE: 'PvE Server Allow Flying Creatures to Carry Wild Dinosaurs',
      AllowFlyingStaminaRecovery: 'Flying Creatures Recover Stamina When Player Standing',
      ForceFlyerExplosives: 'Allow Flying Creatures to Carry C4 While Flying',
      AllowMultipleAttachedC4: 'Allow Multiple C4 Attached to Dinosaurs',
      AutoDestroyDecayedDinos: 'Auto Destroy Decayed Claimable Dinosaurs',
      // Building system
      TheMaxStructuresInRange: 'Server Maximum Structure Count',
      ForceAllStructureLocking: 'Default Lock All Structures',
      EnableExtraStructurePreventionVolumes: 'Disable Building in Specific Resource-rich Areas',
      OverrideStructurePlatformPrevention: 'Allow Building Turrets on Platform Saddles',
      DestroyUnconnectedWaterPipes: 'Auto Destroy Unconnected Water Pipes',
      FastDecayUnsnappedCoreStructures: 'Fast Decay Unsnapped Core Structures',
      OnlyAutoDestroyCoreStructures: 'Only Auto Destroy Core Structures',
      OnlyDecayUnsnappedCoreStructures: 'Only Decay Unsnapped Core Structures',
      // PvP settings
      PreventOfflinePvP: 'Enable Offline PvP Protection',
      PreventOfflinePvPInterval: 'Offline PvP Protection Wait Time (seconds)',
      PreventTribeAlliances: 'Disable Tribe Alliances',
      // Game mechanics
      AllowCrateSpawnsOnTopOfStructures: 'Allow Supply Crates to Spawn on Structures',
      RandomSupplyCratePoints: 'Random Supply Crate Locations',
      NonPermanentDiseases: 'Non-permanent Diseases (Disappear After Respawn)',
      PreventDiseases: 'Completely Disable Diseases',
      PreventSpawnAnimations: 'Disable Spawn Animations',
      UseOptimizedHarvestingHealth: 'Optimize Harvesting Health',
      ClampItemSpoilingTimes: 'Clamp Item Spoiling Times to Maximum',
      ClampItemStats: 'Enable Item Stat Limits',
      ExtinctionEventTimeInterval: 'Extinction Event Time Interval (seconds)',
      OxygenSwimSpeedStatMultiplier: 'Oxygen Swimming Speed Stat Multiplier',
      // Data transfer control
      PreventDownloadSurvivors: 'Prevent Download Survivor Data',
      PreventDownloadItems: 'Prevent Download Item Data',
      PreventDownloadDinos: 'Prevent Download Dinosaur Data',
      PreventUploadSurvivors: 'Prevent Upload Survivor Data',
      PreventUploadItems: 'Prevent Upload Item Data',
      PreventUploadDinos: 'Prevent Upload Dinosaur Data',
      // Advanced configuration
      CustomLiveTuningUrl: 'Custom Live Tuning File URL',
      customdynamicconfigurl: 'Custom Dynamic Configuration File URL',
      EnableFullDump: 'Enable Full Server Log Dump',
      TribeLogDestroyedEnemyStructures: 'Show Enemy Structure Destruction in Tribe Log'
    },
    // Command line parameters translation
    commandLineArgs: {
      // Performance optimization
      allcores: 'Use All CPU Cores',
      USEALLAVAILABLECORES: 'Force Use All Available CPU Cores',
      lowmemory: 'Low Memory Mode (Save 800MB RAM)',
      nomemorybias: 'Reduce Memory Usage (About 600MB)',
      usecache: 'Use Cache (70% Loading Speed Improvement)',
      high: 'High Priority Process',
      nocombineclientmoves: 'Disable Client Move Optimization',
      // Graphics optimization
      allowansel: 'Enable NVIDIA Ansel Support',
      d3d10: 'Force Use DX10 (Windows Only)',
      sm4: 'Force Use Shader Model 4 (Windows Only)',
      nomansky: 'Disable Sky Details (Clouds, Stars)',
      noaafonts: 'Remove Font Anti-aliasing',
      // Security settings
      NoBattlEye: 'Disable BattleEye Anti-cheat',
      noantispeedhack: 'Disable Anti-speedhack Detection',
      insecure: 'Disable VAC System',
      UseItemDupeCheck: 'Enable Extra Duplication Protection',
      exclusivejoin: 'Whitelist Mode',
      // Logging
      servergamelog: 'Enable Server Game Log',
      servergamelogincludetribelogs: 'Include Tribe Logs in Game Log',
      ServerRCONOutputTribeLogs: 'Show Tribe Chat in RCON',
      gameplaylogging: 'Gameplay Logging',
      webalarm: 'Enable Web Alerts',
      // Mods and platform
      automanagedmods: 'Auto Manage MOD Download/Install/Update',
      crossplay: 'Enable Crossplay (EPIC and Steam)',
      epiconly: 'Epic Players Only',
      UseVivox: 'Enable Vivox Voice',
      MapModID: 'Custom Map MOD ID',
      PublicIPForEpic: 'Public IP Address for Epic Clients',
      // Game features
      ForceAllowCaveFlyers: 'Force Allow Flying Creatures in Caves',
      ForceRespawnDinos: 'Destroy All Wild Creatures on Startup',
      NoBiomeWalls: 'Remove Biome Transition Wall Effects',
      PreventHibernation: 'Prevent Hibernation (Performance Cost)',
      StasisKeepControllers: 'Keep AI Controllers in Stasis',
      EnableIdlePlayerKick: 'Kick Idle Players',
      nofishloot: 'Disable Fishing Non-meat Loot',
      ActiveEvent: 'Enable Specified Event',
      // Maintenance tools
      ClearOldItems: 'Clear Old Items (One-time)',
      noninlinesaveload: 'Fix Large Save Corruption Issues',
      oldsaveformat: 'Use Old Save Format',
      StructureDestructionTag: 'Structure Destruction Tag',
      culture: 'Language Code Override'
    },

    // Game.ini parameter categories
    gameIniCategories: {
      gameBasic: 'Basic Game',
      experienceSettings: 'Experience and Level',
      breedingSettings: 'Breeding',
      itemSettings: 'Item and Resource',
      dinoSettings: 'Dinosaur',
      tribeSettings: 'Tribe and Player',
      pvpSettings: 'PvP',
      structureSettings: 'Building and Structure',
      advancedSettings: 'Advanced',
      customSettings: 'Custom Configuration'
    },

    // GameUserSettings.ini parameter categories
    gameUserSettingsCategories: {
      serverBasic: 'Server Basic',
      gameMode: 'Game Mode',
      communication: 'Chat and Communication',
      gameMultipliers: 'Game Multiplier',
      characterSettings: 'Character',
      dinoSettings: 'Dinosaur',
      environmentSettings: 'Environment',
      structureSettings: 'Structure',
      tribeSettings: 'Tribe and Alliance',
      breedingSettings: 'Breeding and Imprinting',
      itemSettings: 'Item and Supply',
      performanceSettings: 'Server Performance',
      diseaseSettings: 'Disease and Status',
      offlineRaidSettings: 'Offline Raid Protection',
      crossArkSettings: 'Cross-ARK Transfer',
      flyerSettings: 'Flyer',
      advancedSettings: 'Advanced Feature'
    },

    // Game.ini parameter translations
    gameIniParams: {
      // Basic settings
      bUseSingleplayerSettings: 'Use Singleplayer Settings',
      bDisableStructurePlacementCollision: 'Disable Structure Placement Collision',
      bAllowFlyerCarryPvE: 'Allow Flyer Carry PvE',
      bDisableStructureDecayPvE: 'Disable Structure Decay PvE',
      bAllowUnlimitedRespecs: 'Allow Unlimited Respecs',
      bAllowPlatformSaddleMultiFloors: 'Allow Platform Saddle Multi Floors',
      bPassiveDefensesDamageRiderlessDinos: 'Passive Defenses Damage Riderless Dinos',
      bPvEDisableFriendlyFire: 'PvE Disable Friendly Fire',
      bDisableFriendlyFire: 'Disable Friendly Fire',
      bEnablePvPGamma: 'Enable PvP Gamma',
      DifficultyOffset: 'Difficulty Offset',
      OverrideOfficialDifficulty: 'Override Official Difficulty',

      // Experience and level settings
      XPMultiplier: 'XP Multiplier',
      PlayerCharacterWaterDrainMultiplier: 'Player Water Drain Multiplier',
      PlayerCharacterFoodDrainMultiplier: 'Player Food Drain Multiplier',
      PlayerCharacterStaminaDrainMultiplier: 'Player Stamina Drain Multiplier',
      PlayerCharacterHealthRecoveryMultiplier: 'Player Health Recovery Multiplier',

      // Breeding settings
      MatingIntervalMultiplier: 'Mating Interval Multiplier',
      EggHatchSpeedMultiplier: 'Egg Hatch Speed Multiplier',
      BabyMatureSpeedMultiplier: 'Baby Mature Speed Multiplier',
      BabyFoodConsumptionSpeedMultiplier: 'Baby Food Consumption Speed Multiplier',
      BabyCuddleIntervalMultiplier: 'Baby Cuddle Interval Multiplier',
      BabyCuddleGracePeriodMultiplier: 'Baby Cuddle Grace Period Multiplier',
      BabyCuddleLoseImprintQualitySpeedMultiplier: 'Baby Cuddle Lose Imprint Quality Speed Multiplier',

      // Item and resource settings
      HarvestAmountMultiplier: 'Harvest Amount Multiplier',
      HarvestHealthMultiplier: 'Harvest Health Multiplier',
      ResourcesRespawnPeriodMultiplier: 'Resources Respawn Period Multiplier',
      ItemStackSizeMultiplier: 'Item Stack Size Multiplier',
      CropGrowthSpeedMultiplier: 'Crop Growth Speed Multiplier',
      GlobalItemDecompositionTimeMultiplier: 'Global Item Decomposition Time Multiplier',
      GlobalCorpseDecompositionTimeMultiplier: 'Global Corpse Decomposition Time Multiplier',

      // Dinosaur settings
      TamingSpeedMultiplier: 'Taming Speed Multiplier',
      DinoCharacterFoodDrainMultiplier: 'Dino Food Drain Multiplier',
      DinoCharacterStaminaDrainMultiplier: 'Dino Stamina Drain Multiplier',
      DinoCharacterHealthRecoveryMultiplier: 'Dino Health Recovery Multiplier',
      DinoCountMultiplier: 'Dino Count Multiplier',
      WildDinoCharacterFoodDrainMultiplier: 'Wild Dino Food Drain Multiplier',
      WildDinoTorporDrainMultiplier: 'Wild Dino Torpor Drain Multiplier',

      // Tribe and player settings
      MaxNumberOfPlayersInTribe: 'Max Number Of Players In Tribe',
      TribeNameChangeCooldown: 'Tribe Name Change Cooldown (Minutes)',
      bPvEAllowTribeWar: 'PvE Allow Tribe War',
      bPvEAllowTribeWarCancel: 'PvE Allow Tribe War Cancel',

      // PvP settings
      bIncreasePvPRespawnInterval: 'Increase PvP Respawn Interval',
      IncreasePvPRespawnIntervalCheckPeriod: 'PvP Respawn Interval Check Period (Seconds)',
      IncreasePvPRespawnIntervalMultiplier: 'PvP Respawn Interval Multiplier',
      IncreasePvPRespawnIntervalBaseAmount: 'PvP Respawn Interval Base Amount (Seconds)',

      // Structure and building settings
      StructureDamageMultiplier: 'Structure Damage Multiplier',
      StructureResistanceMultiplier: 'Structure Resistance Multiplier',
      StructureDamageRepairCooldown: 'Structure Damage Repair Cooldown (Seconds)',
      PvEStructureDecayPeriodMultiplier: 'PvE Structure Decay Period Multiplier',
      MaxStructuresInRange: 'Max Structures In Range',

      // Advanced feature settings
      bAutoPvETimer: 'Auto PvE Timer',
      bAutoPvEUseSystemTime: 'Auto PvE Use System Time',
      AutoPvEStartTimeSeconds: 'Auto PvE Start Time (Seconds)',
      AutoPvEStopTimeSeconds: 'Auto PvE Stop Time (Seconds)',
      bOnlyAllowSpecifiedEngrams: 'Only Allow Specified Engrams',
      bAutoUnlockAllEngrams: 'Auto Unlock All Engrams',
      bShowCreativeMode: 'Show Creative Mode',
      bUseCorpseLocator: 'Use Corpse Locator',
      bDisableLootCrates: 'Disable Loot Crates',
      bDisableDinoRiding: 'Disable Dino Riding',
      bDisableDinoTaming: 'Disable Dino Taming',
      bAllowCustomRecipes: 'Allow Custom Recipes',

      // Custom configuration
      DayCycleSpeedScale: 'Day Cycle Speed Scale',
      NightTimeSpeedScale: 'Night Time Speed Scale',
      DayTimeSpeedScale: 'Day Time Speed Scale'
    },

    // GameUserSettings.ini parameter translations
    gameUserSettingsParams: {
      // Server basic settings
      ServerPassword: 'Server Password',
      SpectatorPassword: 'Spectator Password',
      AdminLogging: 'Admin Logging',

      // Game mode settings
      serverPVE: 'PvE Mode',
      serverHardcore: 'Hardcore Mode',
      ShowMapPlayerLocation: 'Show Player Location',
      allowThirdPersonPlayer: 'Allow Third Person',
      ServerCrosshair: 'Show Crosshair',
      EnablePvPGamma: 'PvP Gamma Adjustment',
      DisablePvEGamma: 'Disable PvE Gamma Adjustment',
      serverForceNoHud: 'Force Hide HUD',
      ShowFloatingDamageText: 'Show Floating Damage Text',
      AllowHitMarkers: 'Allow Hit Markers',

      // Chat and communication settings
      globalVoiceChat: 'Global Voice Chat',
      proximityChat: 'Proximity Chat',
      alwaysNotifyPlayerJoined: 'Always Notify Player Joined',
      alwaysNotifyPlayerLeft: 'Always Notify Player Left',
      DontAlwaysNotifyPlayerJoined: 'Disable Player Join Notification',

      // Game multiplier settings
      XPMultiplier: 'XP Multiplier',
      TamingSpeedMultiplier: 'Taming Speed Multiplier',
      HarvestAmountMultiplier: 'Harvest Amount Multiplier',
      HarvestHealthMultiplier: 'Harvest Health Multiplier',
      ResourcesRespawnPeriodMultiplier: 'Resources Respawn Period Multiplier',
      ItemStackSizeMultiplier: 'Item Stack Size Multiplier',

      // Character settings
      PlayerCharacterHealthRecoveryMultiplier: 'Player Health Recovery Multiplier',
      PlayerCharacterFoodDrainMultiplier: 'Player Food Drain Multiplier',
      PlayerCharacterWaterDrainMultiplier: 'Player Water Drain Multiplier',
      PlayerCharacterStaminaDrainMultiplier: 'Player Stamina Drain Multiplier',
      PlayerDamageMultiplier: 'Player Damage Multiplier',
      PlayerResistanceMultiplier: 'Player Resistance Multiplier',
      OxygenSwimSpeedStatMultiplier: 'Oxygen Swim Speed Stat Multiplier',
      ImplantSuicideCD: 'Implant Suicide Cooldown',

      // Dinosaur settings
      DinoCountMultiplier: 'Dino Count Multiplier',
      DinoCharacterHealthRecoveryMultiplier: 'Dino Health Recovery Multiplier',
      DinoCharacterFoodDrainMultiplier: 'Dino Food Drain Multiplier',
      DinoCharacterStaminaDrainMultiplier: 'Dino Stamina Drain Multiplier',
      DinoDamageMultiplier: 'Dino Damage Multiplier',
      TamedDinoDamageMultiplier: 'Tamed Dino Damage Multiplier',
      DinoResistanceMultiplier: 'Dino Resistance Multiplier',
      TamedDinoResistanceMultiplier: 'Tamed Dino Resistance Multiplier',
      MaxTamedDinos: 'Max Tamed Dinos',
      MaxPersonalTamedDinos: 'Max Personal Tamed Dinos',
      DisableDinoDecayPvE: 'Disable Dino Decay PvE',
      AutoDestroyDecayedDinos: 'Auto Destroy Decayed Dinos',
      PvEDinoDecayPeriodMultiplier: 'PvE Dino Decay Period Multiplier',
      PvPDinoDecay: 'PvP Dino Decay',
      AllowRaidDinoFeeding: 'Allow Raid Dino Feeding',
      RaidDinoCharacterFoodDrainMultiplier: 'Raid Dino Food Drain Multiplier',
      AllowFlyerCarryPvE: 'PvE Allow Flyer Carry',
      bForceCanRideFliers: 'Force Can Ride Fliers',

      // Environment settings
      DayCycleSpeedScale: 'Day Cycle Speed Scale',
      DayTimeSpeedScale: 'Day Time Speed Scale',
      NightTimeSpeedScale: 'Night Time Speed Scale',
      DisableWeatherFog: 'Disable Weather Fog',
      DifficultyOffset: 'Difficulty Offset',
      OverrideOfficialDifficulty: 'Override Official Difficulty',
      RandomSupplyCratePoints: 'Random Supply Crate Points',

      // Structure settings
      StructureDamageMultiplier: 'Structure Damage Multiplier',
      StructureResistanceMultiplier: 'Structure Resistance Multiplier',
      TheMaxStructuresInRange: 'Max Structures In Range',
      NewMaxStructuresInRange: 'New Max Structures In Range',
      MaxStructuresInRange: 'Max Structures In Range',
      DisableStructureDecayPvE: 'Disable Structure Decay PvE',
      PvEStructureDecayPeriodMultiplier: 'PvE Structure Decay Period Multiplier',
      PvEStructureDecayDestructionPeriod: 'PvE Structure Decay Destruction Period',
      PvPStructureDecay: 'PvP Structure Decay',
      StructurePickupTimeAfterPlacement: 'Structure Pickup Time After Placement',
      StructurePickupHoldDuration: 'Structure Pickup Hold Duration',
      AlwaysAllowStructurePickup: 'Always Allow Structure Pickup',
      OnlyAutoDestroyCoreStructures: 'Only Auto Destroy Core Structures',
      OnlyDecayUnsnappedCoreStructures: 'Only Decay Unsnapped Core Structures',
      FastDecayUnsnappedCoreStructures: 'Fast Decay Unsnapped Core Structures',
      DestroyUnconnectedWaterPipes: 'Destroy Unconnected Water Pipes',
      StructurePreventResourceRadiusMultiplier: 'Structure Prevent Resource Radius Multiplier',
      MaxPlatformSaddleStructureLimit: 'Max Platform Saddle Structure Limit',
      PerPlatformMaxStructuresMultiplier: 'Per Platform Max Structures Multiplier',
      PlatformSaddleBuildAreaBoundsMultiplier: 'Platform Saddle Build Area Bounds Multiplier',
      OverrideStructurePlatformPrevention: 'Override Structure Platform Prevention',
      EnableExtraStructurePreventionVolumes: 'Enable Extra Structure Prevention Volumes',
      AllowCaveBuildingPvE: 'Allow Cave Building PvE',
      AllowCaveBuildingPvP: 'Allow Cave Building PvP',
      PvEAllowStructuresAtSupplyDrops: 'PvE Allow Structures At Supply Drops',
      AllowCrateSpawnsOnTopOfStructures: 'Allow Crate Spawns On Top Of Structures',
      bAllowPlatformSaddleMultiFloors: 'Allow Platform Saddle Multi Floors',
      MaxGateFrameOnSaddles: 'Max Gate Frame On Saddles',

      // Tribe and alliance settings
      MaxNumberOfPlayersInTribe: 'Max Number Of Players In Tribe',
      TribeNameChangeCooldown: 'Tribe Name Change Cooldown',
      PreventTribeAlliances: 'Prevent Tribe Alliances',
      MaxAlliancesPerTribe: 'Max Alliances Per Tribe',
      MaxTribesPerAlliance: 'Max Tribes Per Alliance',

      // Breeding and imprinting settings
      AllowAnyoneBabyImprintCuddle: 'Allow Anyone Baby Imprint Cuddle',
      DisableImprintDinoBuff: 'Disable Imprint Dino Buff',
      BabyImprintingStatScaleMultiplier: 'Baby Imprinting Stat Scale Multiplier',

      // Item and supply settings
      ClampItemSpoilingTimes: 'Clamp Item Spoiling Times',
      ClampResourceHarvestDamage: 'Clamp Resource Harvest Damage',
      UseOptimizedHarvestingHealth: 'Use Optimized Harvesting Health',
      BanListURL: 'Ban List URL',

      // Server performance settings
      AutoSavePeriodMinutes: 'Auto Save Period Minutes',
      KickIdlePlayersPeriod: 'Kick Idle Players Period',
      ListenServerTetherDistanceMultiplier: 'Listen Server Tether Distance Multiplier',
      RCONServerGameLogBuffer: 'RCON Server Game Log Buffer',
      NPCNetworkStasisRangeScalePlayerCountStart: 'NPC Network Stasis Range Scale Player Count Start',
      NPCNetworkStasisRangeScalePlayerCountEnd: 'NPC Network Stasis Range Scale Player Count End',
      NPCNetworkStasisRangeScalePercentEnd: 'NPC Network Stasis Range Scale Percent End',

      // Disease and status settings
      PreventDiseases: 'Prevent Diseases',
      NonPermanentDiseases: 'Non Permanent Diseases',
      PreventSpawnAnimations: 'Prevent Spawn Animations',

      // Offline raid protection settings
      PreventOfflinePvP: 'Prevent Offline PvP',
      PreventOfflinePvPInterval: 'Prevent Offline PvP Interval',

      // Cross-ARK transfer settings
      NoTributeDownloads: 'No Tribute Downloads',
      PreventDownloadSurvivors: 'Prevent Download Survivors',
      PreventDownloadItems: 'Prevent Download Items',
      PreventDownloadDinos: 'Prevent Download Dinos',
      PreventUploadSurvivors: 'Prevent Upload Survivors',
      PreventUploadItems: 'Prevent Upload Items',
      PreventUploadDinos: 'Prevent Upload Dinos',
      CrossARKAllowForeignDinoDownloads: 'Cross ARK Allow Foreign Dino Downloads',
      MaxTributeDinos: 'Max Tribute Dinos',
      MaxTributeItems: 'Max Tribute Items',
      MinimumDinoReuploadInterval: 'Minimum Dino Reupload Interval',
      TributeItemExpirationSeconds: 'Tribute Item Expiration Seconds',
      TributeDinoExpirationSeconds: 'Tribute Dino Expiration Seconds',
      TributeCharacterExpirationSeconds: 'Tribute Character Expiration Seconds',

      // Flyer settings
      AllowFlyingStaminaRecovery: 'Allow Flying Stamina Recovery',
      ForceFlyerExplosives: 'Force Flyer Explosives',

      // Advanced feature settings
      AllowMultipleAttachedC4: 'Allow Multiple Attached C4',
      AllowIntegratedSPlusStructures: 'Allow Integrated S+ Structures',
      AllowHideDamageSourceFromLogs: 'Allow Hide Damage Source From Logs',
      AllowSharedConnections: 'Allow Shared Connections',
      bFilterTribeNames: 'Filter Tribe Names',
      bFilterCharacterNames: 'Filter Character Names',
      bFilterChat: 'Filter Chat',
      EnableCryoSicknessPVE: 'Enable Cryo Sickness PVE',
      EnableCryopodNerf: 'Enable Cryopod Nerf',
      CryopodNerfDuration: 'Cryopod Nerf Duration',
      CryopodNerfDamageMult: 'Cryopod Nerf Damage Multiplier',
      CryopodNerfIncomingDamageMultPercent: 'Cryopod Nerf Incoming Damage Multiplier Percent',
      DisableCryopodEnemyCheck: 'Disable Cryopod Enemy Check',
      DisableCryopodFridgeRequirement: 'Disable Cryopod Fridge Requirement',
      AllowCryoFridgeOnSaddle: 'Allow Cryo Fridge On Saddle',
      MaxTrainCars: 'Max Train Cars',
      MaxHexagonsPerCharacter: 'Max Hexagons Per Character',
      AllowTekSuitPowersInGenesis: 'Allow Tek Suit Powers In Genesis',
      CustomDynamicConfigUrl: 'Custom Dynamic Config URL'
    },

    // GameUserSettings.ini editor description
    gameUserSettingsTextEditDesc: 'Edit GameUserSettings.ini configuration file content directly. Changes will be automatically parsed and synchronized to the visual interface. Switch to visual mode to see the parsed parameter settings.',

    // Default values for GameUserSettings parameters
    defaultValues: {
      sessionName: 'My ARK Server',
      message: 'Welcome to ARK Server!'
    },

    // Placeholders
    placeholders: {
      gameUserSettings: `[ServerSettings]
SessionName=My ARK Server
ServerPassword=
MaxPlayers=70

[SessionSettings]
SessionName=My ARK Server

[MessageOfTheDay]
Message=Welcome to ARK Server!

[/Script/Engine.GameSession]
MaxPlayers=70`,
      gameIni: `[/Script/ShooterGame.ShooterGameMode]
DifficultyOffset=0.2
OverrideOfficialDifficulty=5.0
XPMultiplier=1.0
TamingSpeedMultiplier=1.0
HarvestAmountMultiplier=1.0
ResourcesRespawnPeriodMultiplier=1.0
PlayerCharacterWaterDrainMultiplier=1.0
PlayerCharacterFoodDrainMultiplier=1.0
DinoCharacterFoodDrainMultiplier=1.0
PlayerCharacterStaminaDrainMultiplier=1.0
DinoCharacterStaminaDrainMultiplier=1.0
PlayerCharacterHealthRecoveryMultiplier=1.0
DinoCharacterHealthRecoveryMultiplier=1.0
DinoCountMultiplier=1.0
AllowFlyerCarryPvE=False
MaxTamedDinos=4000
StructureDamageMultiplier=1.0
StructureResistanceMultiplier=1.0
TheMaxStructuresInRange=10500
BabyMatureSpeedMultiplier=1.0
EggHatchSpeedMultiplier=1.0
BabyCuddleIntervalMultiplier=1.0
BabyCuddleGracePeriodMultiplier=1.0
BabyImprintAmountMultiplier=1.0`
    }
  },

  // Modals
  modals: {
    privacyPolicy: 'Privacy Policy',
    termsOfService: 'Terms of Service',
    privacyPolicyContent: 'We value your privacy. This application only collects necessary server management information and will not leak your personal data.',
    termsOfServiceContent: 'Using this Ark Server Commander means you agree to comply with the relevant terms of use. Please use this tool reasonably for server management.'
  },

  // Footer
  footer: {
    copyright: '© {year} Ark Server Commander developed by {company}',
    privacyPolicy: 'Privacy Policy',
    termsOfService: 'Terms of Service',
    support: 'Technical Support',
    github: 'GitHub'
  },

  // Error messages
  errors: {
    networkError: 'Network connection error',
    serverError: 'Server error',
    unauthorized: 'Unauthorized access',
    forbidden: 'Access forbidden',
    notFound: 'Page not found',
    validationError: 'Input validation failed',
    unknownError: 'Unknown error',
    tryAgain: 'Please try again',
    contactSupport: 'Please contact technical support'
  },

  // Success messages
  success: {
    operationSuccess: 'Operation successful',
    dataSaved: 'Data saved',
    dataDeleted: 'Data deleted',
    dataUpdated: 'Data updated',
    dataCreated: 'Data created'
  },

  // Form
  form: {
    required: 'This field is required',
    invalidFormat: 'Invalid format',
    minLength: 'Minimum {min} characters required',
    maxLength: 'Maximum {max} characters allowed',
    invalidEmail: 'Invalid email format',
    invalidUrl: 'Invalid URL format',
    invalidNumber: 'Please enter a valid number',
    invalidPort: 'Port number must be between 1-65535',
    invalidPath: 'Please enter a valid path',
    passwordMismatch: 'Passwords do not match',
    usernameExists: 'Username already exists',
    serverNameExists: 'Server name already exists',
    portInUse: 'Port is already in use'
  }
}