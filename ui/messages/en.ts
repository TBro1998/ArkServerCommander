const messages = {
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
    tagline: 'Create, run and maintain your ARK: Survival Evolved servers from one place.',
    login: 'Sign in',
    logout: 'Logout',
    username: 'Username',
    password: 'Password',
    loginTitle: 'Ark Server Commander',
    loginSubtitle: 'Sign in with your admin account',
    loginButton: 'Sign in',
    loginLoading: 'Signing in...',
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
    passwordMinLength: 'Password must be 8 to 72 bytes long',
    passwordMinLengthError: 'Password must be at least 8 characters',
    passwordMaxLengthError: 'Password is too long: at most 72 bytes (non-ASCII characters take more than one byte)',
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
      download: 'Download',
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
      noAffectedServers: 'No servers use this image',
      runningCount: '{count} of them running',
      updateWarning: 'Before You Update',
      warningDownloadTime: 'The image is large and may take a long time to download',
      warningContainerRecreate: 'Servers recreate their containers with the new image on next start. The first start after that re-validates server files and re-downloads mods, which takes a while',
      warningDataSafety: 'Saves and configs live in data volumes and survive container recreation',
      updateOptions: 'Update Options',
      updateImageOnly: 'Update Image Only',
      updateImageOnlyDesc: 'Only download the new image. Running servers keep the old image and switch to the new one on next start.',
      updateAndRecreate: 'Update Image and Restart Running Servers',
      updateAndRecreateDesc: 'After download, running servers are stopped, recreated with the new image and started again. Players will be disconnected.',
      confirmUpdate: 'Confirm Update',
      updateAvailable: 'Update Available',
      update: 'Update',
      upToDate: 'All images are up to date',
      updatesFound: 'New image versions found, update them below',
      checkingUpdates: 'Checking for updates...',
      checkUpdatesFailed: 'Failed to check for updates. Make sure the host can reach Docker Hub',
      updateStarted: 'Image update started',
      updateFailed: 'Image update failed',
      loadAffectedFailed: 'Failed to load affected servers',
      loadingStatus: 'Loading image status...',
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
    serverUpdateError: 'Failed to save the server config',
    operationFailed: 'Operation failed, please try again later',
    deleteFailed: 'Delete failed, please try again later',
    recreateStarted: 'Container recreation started. The server will be stopped when it finishes and must be started manually',
    recreateFailed: 'Failed to recreate the container, please try again later',
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
      imageOutdated: 'Old Image',
      imageOutdatedTip: 'The image has been updated. The container will be recreated with the new image on next start',
      configPending: 'Config pending',
      configPendingTip: 'Config file changes will be applied on next start',
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
      lastError: 'Last operation failed',
      recreateServer: 'Recreate container',
      confirmRecreateMessage: 'A running server is stopped first (players will be disconnected), then its container is recreated from the current configuration. You must start it again manually. The first start afterwards re-validates server files and re-downloads mods.',
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
      loadServerInfoFailed: 'Failed to load server info, please try again later',
      serverCreateError: 'Failed to create the server',
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
      identifierRequired: 'Server identifier is required',
      configAppliesOnRestart: 'The server is not stopped: changes to GameUserSettings.ini and Game.ini are saved now and written when the server next starts (ARK overwrites these files on shutdown, so they cannot be written immediately).',
      configPending: 'There are saved config file changes that have not been applied yet. They will be applied on the next server start.',
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
        Fjordur: 'Fjordur',
        Ragnarok: 'Ragnarok'
      },
      selectMapPlaceholder: 'Select map or enter map name',
      clearMap: 'Clear map',
      customMapPlaceholder: 'Enter custom map name',
      customMap: 'Custom Map',
      noMatchingMaps: 'No matching maps',
      mapId: 'Map ID',
      searchMaps: 'Search maps',
      officialMaps: 'Official maps'
    },
    // Editor related
    editor: {
      placeholders: {
        gameIni: 'Enter Game.ini content...',
        gameUserSettings: 'Enter GameUserSettings.ini content...'
      },
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
      basic: 'Basic',
      core: 'Core',
      dinos: 'Dinosaurs',
      structures: 'Structures',
      pvp: 'PvP',
      mechanics: 'Game Mechanics',
      transfer: 'Transfer & Cluster',
      performance: 'Performance',
      graphics: 'Graphics',
      security: 'Security',
      logging: 'Logging',
      mods: 'Mods',
      features: 'Features',
      maintenance: 'Maintenance',
      advanced: 'Advanced',
      custom: 'Custom'
    },
    // Query parameters translation
    queryParams: {
      AltSaveDirectoryName: 'Alternative Save Directory Name',
      EventColorsChanceOverride: 'Event Colors Chance Override',
      GameModIds: 'Game Mod IDs',
      NewYear1UTC: 'New Year Event Start Time (UTC)',
      NewYear2UTC: 'New Year Event End Time (UTC)'
    },
    // Command line parameters translation
    commandLineArgs: {
      // Events and Features
      ActiveEvent: 'Active Event',
      NewYearEvent: 'New Year Event',
      UseVivox: 'Use Vivox Voice Chat',
      webalarm: 'Web Alarms',
      AllowChatSpam: 'Allow Chat Spam',
      
      // Mods and Steam
      automanagedmods: 'Auto-Managed Mods',
      MapModID: 'Map Mod ID',
      
      // Crossplay and Network
      crossplay: 'Enable Crossplay',
      epiconly: 'Epic Games Store Only',
      PublicIPForEpic: 'Public IP for Epic Games',
      MULTIHOME: 'Multi-Home IP Address',
      
      // Server Management
      culture: 'Server Language Culture',
      exclusivejoin: 'Exclusive Join (Whitelist)',
      EnableIdlePlayerKick: 'Enable Idle Player Kick',
      MaxNumOfSaveBackups: 'Maximum Number of Save Backups',
      newsaveformat: 'New Save Format',
      NoHangDetection: 'Disable Hang Detection',
      
      // Creatures and Gameplay
      DisableCustomFoldersInTributeInventories: 'Disable Custom Folders in Tribute Inventories',
      ForceAllowCaveFlyers: 'Force Allow Cave Flyers',
      ForceRespawnDinos: 'Force Respawn Dinos',
      NoDinos: 'No Dinosaurs',
      imprintlimit: 'Imprint Limit Percentage',
      MinimumTimeBetweenInventoryRetrieval: 'Minimum Time Between Inventory Retrieval',
      
      // PvP Settings
      DisableRailgunPVP: 'Disable Railgun in PvP',
      pvedisallowtribewar: 'PvE Disallow Tribe War',
      pveallowtribewar: 'PvE Allow Tribe War',
      
      // Security and Anti-Cheat
      insecure: 'Disable VAC (Insecure)',
      NoBattlEye: 'Disable BattlEye',
      noantispeedhack: 'Disable Anti-Speedhack',
      speedhackbias: 'Speedhack Detection Bias',
      noundermeshchecking: 'Disable Under-Mesh Checking',
      noundermeshkilling: 'Disable Under-Mesh Killing',
      SecureSendArKPayload: 'Secure Send ARK Payload',
      UseItemDupeCheck: 'Use Item Duplication Check',
      UseSecureSpawnRules: 'Use Secure Spawn Rules',
      BattlEyeServerRecheck: 'BattlEye Server Recheck',
      
      // Performance Optimization
      nocombineclientmoves: 'Disable Combine Client Moves',
      StasisKeepControllers: 'Stasis Keep Controllers',
      structurememopts: 'Structure Memory Optimizations',
      UseStructureStasisGrid: 'Use Structure Stasis Grid',
      DormancyNetMultiplier: 'Dormancy Network Multiplier',
      nodormancythrottling: 'Disable Dormancy Throttling',
      nitradotest2: 'Nitrado Test Mode 2',
      dedihibernation: 'Dedicated Hibernation',
      
      // Graphics and Client
      ServerAllowAnsel: 'Server Allow NVIDIA Ansel',
      
      // Logging and Admin
      servergamelog: 'Server Game Log',
      servergamelogincludetribelogs: 'Server Game Log Include Tribe Logs',
      ServerRCONOutputTribeLogs: 'Server RCON Output Tribe Logs',
      NotifyAdminCommandsInChat: 'Notify Admin Commands in Chat',
      
      // Transfer and Cluster
      ClusterDirOverride: 'Cluster Directory Override',
      clusterid: 'Cluster ID',
      NoTransferFromFiltering: 'No Transfer From Filtering',
      usestore: 'Use Store',
      BackupTransferPlayerDatas: 'Backup Transfer Player Data',
      converttostore: 'Convert to Store',
      
      // Advanced/Undocumented
      CustomAdminCommandTrackingURL: 'Custom Admin Command Tracking URL',
      CustomMerticsURL: 'Custom Metrics URL',
      CustomNotificationURL: 'Custom Notification URL',
      DisableDupeLogDeletes: 'Disable Duplicate Log Deletes',
      EnableOfficialOnlyVersioningCode: 'Enable Official Only Versioning Code',
      EnableVictoryCoreDupeCheck: 'Enable Victory Core Duplication Check',
      forcedisablemeshchecking: 'Force Disable Mesh Checking',
      ForceDupeLog: 'Force Duplication Log',
      ignoredupeditems: 'Ignore Duplicated Items',
      MaxConnectionsPerIP: 'Maximum Connections Per IP',
      parseservertojson: 'Parse Server to JSON',
      pauseonddos: 'Pause on DDoS',
      PreventTotalConversionSaveDir: 'Prevent Total Conversion Save Directory',
      ReloadedForBackup: 'Reloaded for Backup',
      UnstasisDinoObstructionCheck: 'Unstasis Dino Obstruction Check',
      UseTameEffectivenessClamp: 'Use Tame Effectiveness Clamp',
      UseServerNetSpeedCheck: 'Use Server Network Speed Check'
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
  },

  // Layout
  layout: {
    skipToContent: 'Skip to content',
    redirecting: 'Redirecting...',
    initializing: 'Initializing...',
    verifying: 'Verifying session...',
    menu: 'Menu'
  },

  // Messages for backend error codes (the Chinese UI shows the detailed backend message instead)
  apiErrors: {
    bad_request: 'Invalid request',
    backend_unavailable: 'Backend service unavailable',
    invalid_server_id: 'Invalid server ID',
    server_not_found: 'Server not found',
    identifier_exists: 'Server identifier already exists',
    server_running: 'Cannot delete a running server. Stop it first.',
    already_running: 'Server is already running',
    already_starting: 'Server is already starting',
    already_stopped: 'Server is already stopped',
    already_stopping: 'Server is already stopping',
    image_pulling: 'Image is downloading, please wait',
    port_conflict: 'Ports conflict with another server',
    invalid_port: 'Invalid port',
    server_busy: 'Server is starting or stopping, please try again later',
    image_not_allowed: 'This image cannot be managed',
    invalid_cluster_id: 'Invalid cluster ID',
    invalid_mod_ids: 'Mod IDs must be comma-separated numbers',
    server_not_running: 'Server is not running',
    server_not_stopped: 'Stop the server first',
    invalid_command: 'Invalid command',
    backup_not_found: 'Backup not found',
    invalid_backup: 'Invalid backup file',
    backup_busy: 'A backup or restore is already in progress for this server',
    schedule_not_found: 'Scheduled task not found',
    invalid_schedule: 'Invalid scheduled task',
    rcon_unavailable: 'RCON is not responding; the server may still be downloading, installing mods or loading the map',
    rcon_auth_failed: 'RCON authentication failed; restart the server to rewrite the admin password',
    too_many_attempts: 'Too many failed attempts. Please try again later.',
    invalid_credentials: 'Invalid username or password',
    invalid_password_length: 'Password must be 8 to 72 bytes',
    wrong_password: 'Current password is incorrect'
  },

  // Server status
  status: {
    running: 'Running',
    stopped: 'Stopped',
    starting: 'Starting',
    stopping: 'Stopping',
    restarting: 'Restarting',
    unknown: 'Unknown'
  },

  // Runtime phase
  phase: {
    copying: 'Copying server files',
    downloading: 'Downloading server',
    downloadingProgress: 'Downloading {progress}%',
    installing_mods: 'Installing mods',
    installingModsProgress: 'Mods {current}/{total}',
    launching: 'Loading',
    ready: 'Joinable',
    copyingDesc: 'Copying server files and mods from an installed server, usually a few minutes. On start the files are validated and only differences from the latest version are downloaded.',
    downloadingDesc: 'Downloading or validating server files. The first start downloads about 20 GB, depending on your network speed.',
    installing_modsDesc: 'Downloading and installing mods. Mods that are already installed are skipped.',
    launchingDesc: 'The game process has started and is loading the map and mods. This usually takes a few minutes; players cannot join yet.',
    readyDesc: 'The server is ready and players can join.',
    unknownDesc: 'Checking server state...'
  },

  // Server list
  serverList: {
    filter: 'Filter by status',
    filterEmpty: 'No servers match this filter',
    filters: {
      all: 'All',
      active: 'Running',
      stopped: 'Stopped'
    },
    summary: '{total} servers, {running} running, {players} players online',
    refresh: 'Refresh',
    imageNotReadyGoHome: 'The images required to start servers have not been downloaded yet. Download them under "Image management" on the home page.'
  },

  // Server card
  serverCard: {
    resources: 'Usage',
    viewLogs: 'View logs',
    passwordFailed: 'Failed to load the admin password',
    passwordCopied: 'Admin password copied',
    copyFailed: 'Copy failed, please copy manually',
    imageOutdated: 'Old image',
    imageOutdatedTip: 'The image has been updated; the container will be recreated with the new image on next start',
    configPending: 'Config pending',
    configPendingTip: 'Config file changes will be applied on next start',
    cluster: 'Cluster {id}',
    clusterTip: 'Shares character and item transfers with servers using the same cluster ID',
    lastError: 'Last operation failed',
    downloadHint: 'The first start downloads the server files, please be patient',
    copyingHint: 'Copying server files from "{source}" instead of downloading them',
    launchingHint: 'Loading the map, this usually takes a few minutes',
    players: 'Players online',
    memory: 'Memory',
    map: 'Map',
    ports: 'Ports (game/query/RCON)',
    adminPassword: 'Admin password',
    showPassword: 'Show password',
    hidePassword: 'Hide password',
    copyPassword: 'Copy password'
  },

  // Server actions
  serverActions: {
    start: 'Start',
    stop: 'Stop',
    restart: 'Restart',
    working: 'Working',
    more: 'More actions',
    details: 'Details',
    logs: 'View logs',
    edit: 'Edit config',
    updateGame: 'Update game/mods',
    recreate: 'Recreate container',
    delete: 'Delete server',
    deleteNeedStop: 'Delete (stop first)',
    imageNotReady: 'Images not ready, cannot start',
    startSent: 'Starting',
    startFailed: 'Failed to start the server',
    stopSent: 'Stopping. The world is saved before shutdown.',
    stopFailed: 'Failed to stop the server',
    restartSent: 'Restarting',
    restartFailed: 'Failed to restart the server',
    recreateSent: 'Recreating the container. The server will be stopped when done.',
    recreateFailed: 'Failed to recreate the container',
    updateSent: 'Update started. The update is downloaded when the server starts.',
    updateFailed: 'Update failed',
    deleteSuccess: 'Server deleted',
    deleteFailed: 'Failed to delete the server',
    stopTitle: 'Stop {name}?',
    stopDesc: 'The world is saved before shutdown.',
    stopWithPlayers: '{count} players are online and will be disconnected. The world is saved before shutdown.',
    restartTitle: 'Restart {name}?',
    restartDesc: 'The world is saved and the server stopped, then started with the current config (the container is recreated if the config changed).',
    restartStoppedDesc: 'The server is stopped and will be started with the current config.',
    playersOnline: '{count} players are online and will be disconnected.',
    updateTitle: 'Update {name}',
    updateDesc: 'The server will start and download updates while starting.',
    updateDescRunning: 'The server will be saved and stopped, then started and download updates while starting.',
    updateConfirm: 'Start update',
    updateGameOption: 'Update game files',
    updateGameOptionDesc: 'Validate and update server files with steamcmd. Needed when clients have updated and the server version no longer matches.',
    updateModsOption: 'Update mods',
    updateModsOptionDesc: 'Download all mods again',
    noMods: 'No mods configured',
    recreateTitle: 'Recreate the container of {name}?',
    recreateDesc: 'A running server is stopped first (players will be disconnected), then the container is recreated with the current config. Start it manually afterwards. Saves and server files are kept in the data directory.',
    deleteTitle: 'Delete {name}?',
    deleteDesc: 'The server record and container will be deleted.',
    deleteData: 'Also delete all data',
    deleteDataOff: 'Files in the data directory (saves, server files) and backups are kept and can be cleaned up manually later',
    deleteDataOn: 'Saves, config, server files, plugins and all backups will be permanently deleted',
    typeNameToConfirm: 'Type the server name "{name}" to confirm'
  },

  // Connection info
  connection: {
    address: 'Address',
    steamConnect: 'Steam connect',
    copyAddress: 'Copy address',
    copyLink: 'Copy Steam link',
    copied: 'Copied {text}',
    copyFailed: 'Copy failed, please copy manually',
    hint: 'Type "open" followed by the address in the game console, or click the Steam link to join. The address comes from PUBLIC_HOST or the detected public IP; the ports must be open on your router/firewall.'
  },

  // Mods
  mods: {
    invalid: 'Mod IDs must be comma-separated numbers',
    loading: 'Looking up name...',
    notFound: 'Mod not found on the Workshop',
    unavailable: 'Could not reach Steam to look up the name'
  },

  // Server details
  detail: {
    tabs: {
      overview: 'Overview',
      console: 'Console',
      logs: 'Logs',
      backups: 'Backups',
      schedules: 'Schedules'
    },
    backToList: 'Back to servers',
    notFound: 'Server not found or deleted',
    runtime: 'Runtime',
    collecting: 'Collecting runtime info...',
    notRunning: 'Server is not running',
    players: 'Players',
    memory: 'Memory',
    advertisedName: 'Name in server list',
    connection: 'Connection',
    noPublicHost: 'Could not detect the public address. Set PUBLIC_HOST to configure it.',
    gamePort: 'Game port',
    queryPort: 'Query port',
    rconPort: 'RCON port',
    config: 'Configuration',
    identifier: 'Identifier',
    map: 'Map',
    maxPlayers: 'Max players',
    cluster: 'Cluster ID',
    none: 'None',
    autoRestart: 'Restart on crash',
    autoUpdate: 'Update on start',
    on: 'On',
    off: 'Off',
    createdAt: 'Created',
    mods: 'Mods',
    noMods: 'No mods configured',
    imageOutdated: 'Old image',
    configPending: 'Config pending',
    lastError: 'Last operation failed',
    viewLogs: 'View logs'
  },

  // Logs
  logs: {
    connecting: 'Connecting...',
    live: 'Live',
    paused: 'Paused',
    ended: 'Ended',
    error: 'Connection failed',
    filter: 'Filter logs...',
    pause: 'Pause',
    resume: 'Resume',
    reconnect: 'Reconnect',
    clear: 'Clear',
    download: 'Download logs',
    downloadFailed: 'Failed to download logs',
    empty: 'No logs yet (the container has not been created or just started)',
    scrollToBottom: 'Jump to bottom',
    hint: 'This is the container output: server download progress, mod installation and game startup. Detailed game logs are in the Logs folder of the save directory.'
  },

  // Console
  console: {
    notRunning: 'The server is not running. Start it to use the console.',
    notReady: 'The server is still preparing (downloading, installing mods or loading the map); RCON commands may not work yet.',
    quickActions: 'Quick actions',
    broadcast: 'Broadcast',
    broadcastPlaceholder: 'Message to broadcast to all players',
    broadcastSent: 'Broadcast sent',
    saveWorld: 'Save world',
    saved: 'World saved',
    rconConsole: 'RCON console',
    consoleHint: 'Type an RCON command and press Enter. Use the arrow keys to browse history.',
    commandPlaceholder: 'RCON command, e.g. ListPlayers',
    noOutput: '(no output)',
    commandFailed: 'Command failed',
    players: 'Players online',
    playersFailed: 'Failed to load players',
    playersNotReady: 'Players are available once the server is ready.',
    playersUnknown: 'Click refresh to load players',
    noPlayers: 'No players online',
    loading: 'Loading...',
    steamIdCopied: 'SteamID copied',
    kick: 'Kick',
    ban: 'Ban',
    kickTitle: 'Kick {name}?',
    kickDesc: 'The player is disconnected but can rejoin.',
    banTitle: 'Ban {name}?',
    banDesc: 'The player is kicked and added to the ban list. Use the UnbanPlayer command to unban.',
    kicked: 'Kicked {name}',
    banned: 'Banned {name}'
  },

  // Backups
  backups: {
    title: 'Save backups',
    summary: '{count} backups, {size}',
    hint: 'Backups contain the saves (SavedArks) and config files, not logs. The world is saved first when the server is running. Restoring requires the server to be stopped, and the current save is backed up automatically before restoring.',
    create: 'Back up now',
    upload: 'Upload backup',
    empty: 'No backups yet',
    time: 'Time',
    kind: 'Type',
    size: 'Size',
    actions: 'Actions',
    kinds: {
      manual: 'Manual',
      scheduled: 'Scheduled',
      'pre-restore': 'Pre-restore',
      uploaded: 'Uploaded'
    },
    download: 'Download',
    restore: 'Restore',
    restoreNeedStop: 'Restore (stop the server first)',
    delete: 'Delete',
    loadFailed: 'Failed to load backups',
    created: 'Backup completed',
    createFailed: 'Backup failed',
    uploaded: 'Backup uploaded',
    uploadFailed: 'Upload failed',
    restored: 'Save restored',
    restoreFailed: 'Restore failed',
    deleted: 'Backup deleted',
    deleteFailed: 'Failed to delete the backup',
    restoreTitle: 'Restore from backup?',
    restoreDesc: 'The current saves and config files will be replaced by the backup from {time}. The current save is backed up first.',
    deleteTitle: 'Delete backup?',
    deleteDesc: 'The backup from {time} will be permanently deleted.'
  },

  // Schedules
  schedules: {
    title: 'Scheduled tasks',
    timezoneHint: 'Times use the server time zone (CONTAINER_TZ)',
    add: 'Add task',
    edit: 'Edit task',
    empty: 'No scheduled tasks yet. Add daily restarts, regular backups and more.',
    actions: {
      restart: 'Restart',
      backup: 'Back up saves',
      update: 'Update game/mods',
      command: 'RCON command',
      start: 'Start',
      stop: 'Stop'
    },
    actionDesc: {
      restart: 'A running server is saved and restarted after the broadcast countdown; stopped servers are skipped.',
      backup: 'Back up saves and config, optionally keeping only the latest few.',
      update: 'Stops after the countdown, updates game files and mods, then starts. Stopped servers update on their next start.',
      command: 'Run an RCON command on the running server.',
      start: 'Start the server if it is stopped.',
      stop: 'A running server is saved and stopped after the broadcast countdown.'
    },
    action: 'Task',
    frequency: 'Frequency',
    modes: {
      daily: 'Daily',
      weekly: 'Weekly',
      interval: 'Every few hours',
      custom: 'Custom cron'
    },
    weekdays: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
    weekday: 'Day',
    time: 'Time',
    intervalHours: 'Interval (hours)',
    cronHint: '5-field cron: minute hour day month weekday, e.g. 30 4 * * 1 means every Monday at 4:30',
    warnMinutes: 'Broadcast warning (minutes)',
    warnMinutesHint: 'Broadcast an in-game countdown before running; 0 disables it',
    message: 'Broadcast message',
    messagePlaceholder: 'Leave empty for the default; {minutes} is replaced with the remaining minutes',
    command: 'Command',
    keep: 'Backups to keep',
    keepHint: 'Keep only the latest scheduled backups; 0 keeps all',
    enabled: 'Enabled',
    runNow: 'Run now',
    delete: 'Delete',
    deleteTitle: 'Delete scheduled task?',
    cancel: 'Cancel',
    save: 'Save',
    everyDayAt: 'Every day at {time}',
    everyWeekAt: 'Every {day} at {time}',
    everyHours: 'Every {hours} hours',
    nextRun: 'Next {time}',
    lastRun: 'Last {time}',
    ok: 'OK',
    neverRun: 'Never run',
    saved: 'Scheduled task saved',
    saveFailed: 'Failed to save the scheduled task',
    loadFailed: 'Failed to load scheduled tasks',
    runStarted: 'Started',
    runFailed: 'Failed to run',
    deleted: 'Scheduled task deleted',
    deleteFailed: 'Failed to delete the scheduled task'
  },

  // Server form
  serverForm: {
    unsaved: 'Unsaved changes',
    saveShortcut: 'Save (Ctrl+S)',
    sectionBasic: 'Basics',
    sectionNetwork: 'Ports and password',
    sectionMods: 'Mods and cluster',
    sectionMaintenance: 'Maintenance',
    sessionNameHint: 'The name players see in the server list',
    identifierHint: 'Unique identifier used only inside this tool; not shown to players',
    gamePortHint: 'Also uses port +1',
    queryPortHint: 'Used by the Steam server list',
    rconPortHint: 'Used for remote administration',
    portsFirewallHint: 'Ports not used by other servers were filled in automatically. Open the game port, game port +1 and query port (UDP) on your router and firewall.',
    passwordHint: 'Used for in-game admin commands (enablecheats) and RCON',
    passwordKeepHint: 'Leave empty to keep the current password',
    passwordKeepPlaceholder: 'Leave empty to keep unchanged',
    generate: 'Generate',
    modsHint: 'Workshop mod IDs separated by commas; mod names are shown below',
    clusterHint: 'Servers with the same cluster ID can transfer characters, dinos and items (shared cluster storage). Letters, digits, _ . - only.',
    autoRestart: 'Restart on crash',
    autoRestartHint: 'Docker restarts the server automatically if the game process exits unexpectedly',
    autoUpdate: 'Update on start',
    autoUpdateHint: 'Validate game files and download mods again before every start. Starting becomes slower.',
    identifierRequired: 'Server identifier is required',
    passwordRequired: 'Admin password is required',
    passwordNoNewline: 'Password cannot contain line breaks',
    portRange: 'Port must be between 1 and 65535',
    portDuplicate: 'Duplicate ports: game port, game port +1, query port and RCON port must all differ',
    maxPlayersRange: 'Max players must be between 1 and 200',
    clusterInvalid: 'Cluster ID may only contain letters, digits, _ . - and must start with a letter or digit',
    modsInvalid: 'Mod IDs must be comma-separated numbers',
    fixErrors: 'Please fix the highlighted fields under "Basic"',
    discardConfirm: 'You have unsaved changes. Leave anyway?',
    saveAndRestart: 'Save and restart',
    created: 'Server created',
    savedAndRestarting: 'Saved, the server is restarting',
    savedPending: 'Saved. Config changes take effect on next start.'
  },

  // Dashboard
  dashboard: {
    hostCopied: 'Public address copied',
    copyHost: 'Copy public address',
    viewAll: 'All servers',
    createFirst: 'Create your first server',
    title: 'Overview',
    publicHost: 'Public address',
    manageServers: 'Manage servers',
    servers: 'Servers',
    running: 'Running',
    readyCount: '{count} joinable',
    players: 'Players online',
    cpuHint: '100% = one core',
    memory: 'Memory',
    disk: 'Disk usage',
    recentErrors: 'Recent errors',
    imageManagement: 'Image management',
    diskUsage: 'Disk usage',
    diskHint: 'Docker computes volume sizes slowly; results are cached for 5 minutes',
    recalculate: 'Recalculate',
    noServers: 'No servers yet',
    server: 'Server',
    files: 'Server files',
    saves: 'Saves/plugins',
    backups: 'Backups',
    total: 'Total',
    clusterVolume: 'Cluster storage {name}',
    downloadStarted: 'Started downloading {name}',
    downloadFailed: 'Failed to download the image',
    diskFailed: 'Failed to compute disk usage'
  },

  // Settings
  settings: {
    title: 'Settings',
    account: 'Account',
    username: 'Username',
    changePassword: 'Change login password',
    oldPassword: 'Current password',
    newPassword: 'New password',
    confirmPassword: 'Confirm new password',
    passwordHint: '8 to 72 bytes',
    submit: 'Change password',
    sessionNote: 'You stay signed in. Sessions on other devices expire within 24 hours.',
    tooShort: 'New password must be at least 8 bytes',
    tooLong: 'New password can be at most 72 bytes',
    mismatch: 'The new passwords do not match',
    changed: 'Password changed',
    failed: 'Failed to change the password'
  }
}

export default messages;
