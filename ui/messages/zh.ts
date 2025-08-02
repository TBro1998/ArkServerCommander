export default {
  // 通用
  common: {
    confirm: '确定',
    cancel: '取消',
    save: '保存',
    delete: '删除',
    edit: '编辑',
    add: '添加',
    close: '关闭',
    loading: '加载中...',
    error: '错误',
    success: '成功',
    warning: '警告',
    info: '信息',
    yes: '是',
    no: '否',
    back: '返回',
    next: '下一步',
    previous: '上一步',
    submit: '提交',
    reset: '重置',
    search: '搜索',
    filter: '筛选',
    sort: '排序',
    refresh: '刷新',
    copy: '复制',
    download: '下载',
    upload: '上传',
    export: '导出',
    import: '导入',
    settings: '设置',
    help: '帮助',
    about: '关于',
    version: '版本',
    language: '语言',
    theme: '主题',
    dark: '深色',
    light: '浅色',
    auto: '自动'
  },

  // 导航
  navigation: {
    home: '首页',
    dashboard: '控制台',
    servers: '服务器管理',
    players: '玩家管理',
    logs: '日志监控',
    settings: '设置',
    logout: '退出登录',
    welcome: '欢迎',
    user: '用户'
  },

  // 认证
  auth: {
    login: '登录',
    logout: '退出登录',
    username: '用户名',
    password: '密码',
    loginTitle: 'ARK 服务器管理器',
    loginSubtitle: '安全登录您的管理账户',
    loginButton: '登录',
    loginLoading: '登录中...',
    loginError: '登录失败',
    loginSuccess: '登录成功',
    logoutSuccess: '已退出登录',
    initCheck: '检查系统初始化状态',
    initRequired: '系统需要初始化',
    alreadyLoggedIn: '您已经登录',
    enterUsername: '请输入用户名',
    enterPassword: '请输入密码',
    firstTimeTip: '首次使用？系统将自动引导您完成初始化',
    secureLogin: '安全登录系统',
    // 初始化相关
    initTitle: '系统初始化',
    initSubtitle: '首次使用，请设置管理员账户',
    adminUsername: '管理员用户名',
    enterAdminUsername: '请输入管理员用户名',
    confirmPassword: '确认密码',
    enterConfirmPassword: '请再次输入密码',
    passwordMinLength: '密码至少需要6位字符',
    passwordMinLengthError: '密码至少需要6位',
    passwordMismatch: '两次输入的密码不一致',
    initButton: '初始化系统',
    initLoading: '初始化中...',
    initSuccess: '系统初始化成功',
    initError: '系统初始化失败',
    initTip: '初始化完成后将自动跳转到主页面',
    initWizard: '系统初始化向导'
  },

  // 首页
  home: {
    title: '欢迎使用 ARK 服务器管理器',
    subtitle: '您已成功登录系统，可以开始管理您的 ARK 服务器了。',
    systemInfo: '系统信息',
    username: '用户名',
    userID: '用户ID',
    imageManagement: '镜像管理',
    features: '功能模块',
    serverManagement: '服务器管理',
    serverManagementDesc: '添加、配置和管理您的ARK服务器，支持一键启动、停止和监控',
    startManage: '开始管理',
    playerManagement: '玩家管理',
    playerManagementDesc: '管理服务器玩家，查看在线状态和权限设置',
    logMonitoring: '日志监控',
    logMonitoringDesc: '实时监控服务器日志，查看系统状态和性能指标',
    comingSoon: '即将推出',
    tip: '点击上方卡片开始管理您的ARK服务器'
  },

  // 服务器管理
  servers: {
    title: '服务器管理',
    serverManagementDesc: '管理和监控您的ARK服务器实例',
    addServer: '添加服务器',
    editServer: '编辑服务器',
    deleteServer: '删除服务器',
    serverName: '服务器名称',
    serverPort: '服务器端口',
    serverPath: '服务器路径',
    serverStatus: '服务器状态',
    serverActions: '操作',
    startServer: '启动服务器',
    stopServer: '停止服务器',
    restartServer: '重启服务器',
    viewLogs: '查看日志',
    serverConfig: '服务器配置',
    gameIni: 'Game.ini 配置',
    gameUserSettings: 'GameUserSettings.ini 配置',
    serverArgs: '启动参数',
    running: '运行中',
    stopped: '已停止',
    starting: '启动中',
    stopping: '停止中',
    error: '错误',
    unknown: '未知',
    confirmDelete: '确定要删除这个服务器吗？',
    deleteWarning: '此操作不可撤销',
    serverAdded: '服务器添加成功',
    serverUpdated: '服务器更新成功',
    serverDeleted: '服务器删除成功',
    serverStartSuccess: '服务器启动成功',
    serverStopSuccess: '服务器停止成功',
    serverRestartSuccess: '服务器重启成功',
    serverStartError: '服务器启动失败',
    serverStopError: '服务器停止失败',
    serverRestartError: '服务器重启失败',
    noServers: '暂无服务器',
    noServersDesc: '点击"添加服务器"开始创建您的第一个ARK服务器',
    serverConfigSaved: '服务器配置已保存',
    serverConfigError: '服务器配置保存失败',
    invalidPort: '端口号无效',
    invalidPath: '服务器路径无效',
    portInUse: '端口已被占用',
    pathNotExists: '服务器路径不存在',
    serverNameRequired: '服务器名称不能为空',
    serverPortRequired: '服务器端口不能为空',
    serverPathRequired: '服务器路径不能为空',
    imageStatus: '镜像状态',
    imageDownloading: '镜像下载中',
    imageNotReady: '镜像未就绪',
    imageDownloadingDesc: '正在下载镜像，请稍后创建服务器',
    imageNotReadyDesc: '镜像未就绪，无法创建服务器',
    // 镜像状态详细翻译
    dockerImages: {
      title: '镜像下载状态',
      overallStatus: '总体状态',
      imageReady: '镜像就绪',
      imageNotReady: '镜像未就绪（无法启动服务器）',
      imageMissingManualDownload: '镜像缺失，请手动下载',
      downloading: '下载中',
      ready: '就绪',
      notReady: '未就绪',
      waitingDownload: '等待下载',
      layerProgress: '层级下载进度',
      totalImages: '镜像总数',
      downloadingCount: '下载中',
      refreshStatus: '刷新镜像状态',
      manualDownload: '手动下载',
      checkUpdates: '检查更新',
      updateConfirm: '镜像更新确认',
      imageInfo: '镜像信息',
      imageName: '镜像名称',
      affectedServers: '受影响的服务器',
      updateWarning: '更新风险提示',
      warningDownloadTime: '镜像下载可能需要较长时间，请耐心等待',
      warningContainerRecreate: '容器重建将导致服务器短暂停机',
      warningDataSafety: '请确保已备份重要数据，避免数据丢失',
      updateOptions: '更新选项',
      updateImageOnly: '仅更新镜像',
      updateImageOnlyDesc: '只下载新镜像，不重建容器。需要手动重建容器以使用新镜像。',
      updateAndRecreate: '更新镜像并重建容器',
      updateAndRecreateDesc: '下载新镜像并自动重建所有受影响的容器。服务器将短暂停机。',
      confirmUpdate: '确认更新',
      unknownSize: '未知大小',
      // 镜像名称
      arkServer: 'ARK服务器',
      alpineSystem: 'Alpine系统',
      // 层级信息
      layerDetails: '层级详情',
      layers: '层级',
      // 层级状态
      layerStatus: {
        pending: '等待中',
        downloading: '下载中',
        extracting: '解压中',
        verifying: '验证中',
        complete: '已完成'
      }
    },
    cannotDeleteRunning: '无法删除正在运行的服务器，请先停止服务器',
    serverCreateSuccess: '服务器创建成功',
    serverUpdateSuccess: '服务器更新成功',
    serverDeleteSuccess: '服务器删除成功',
    serverStartInProgress: '服务器启动中...',
    serverStopInProgress: '服务器停止中...',
    copyToClipboard: '已复制到剪贴板',
    copyFailed: '复制失败，请手动复制',
    authenticationFailed: '认证失败，请重新登录',
    getServerListFailed: '获取服务器列表失败，请稍后重试',
    loadServerInfoFailed: '加载服务器信息失败，请稍后重试',
    operationFailed: '操作失败，请稍后重试',
    deleteFailed: '删除失败，请稍后重试',
    startServerFailed: '启动服务器失败，请稍后重试',
    stopServerFailed: '停止服务器失败，请稍后重试',
    imageStatusError: '获取镜像状态失败',
    // 服务器卡片相关
    card: {
      startServer: '启动服务器',
      stopServer: '停止服务器',
      starting: '启动中...',
      stopping: '停止中...',
      unknownStatus: '未知状态',
      cannotStartImageNotReady: '镜像未就绪，无法启动',
      rconInfo: 'RCON信息',
      rconConnectionInfo: 'RCON连接信息',
      serverIdentifier: '服务器标识',
      rconPort: 'RCON端口',
      adminPassword: '管理员密码',
      editServer: '编辑服务器',
      deleteServer: '删除服务器',
      confirmDelete: '确认删除',
      confirmDeleteMessage: '您确定要删除服务器 "{identifier}" 吗？此操作无法撤销。',
      status: '状态',
      serverName: '服务器名称',
      clusterId: '集群ID',
      map: '地图',
      maxPlayers: '最大玩家数',
      portConfig: '端口配置',
      gamePort: '游戏端口',
      queryPort: '查询端口',
      rconPortLabel: 'RCON端口',
      authInfo: '认证信息',
      timeInfo: '时间信息',
      createdAt: '创建时间',
      updatedAt: '更新时间',
      serverId: '服务器ID',
      copy: '复制',
      close: '关闭',
      showPassword: '显示密码',
      hidePassword: '隐藏密码'
    },
    // 服务器编辑相关
    edit: {
      title: '服务器编辑',
      createTitle: '新增服务器',
      editTitle: '编辑服务器',
      createServerDesc: '配置并创建一个新的ARK服务器实例',
      basicParams: '基本参数',
      gameUserSettings: 'GameUserSettings.ini',
      gameIni: 'Game.ini',
      serverArgs: '启动参数',
      serverIdentifier: '服务器标识',
      serverIdentifierRequired: '服务器标识 *',
      serverIdentifierPlaceholder: '输入服务器标识',
      serverName: '服务器名称',
      serverNamePlaceholder: '输入服务器名称',
      serverNameDesc: '显示在游戏服务器列表中的名称',
      clusterId: '集群ID',
      clusterIdPlaceholder: '输入集群ID（可选）',
      clusterIdDesc: '用于集群服务器之间的数据共享',
      gamePort: '游戏端口',
      gamePortRequired: '游戏端口 *',
      gamePortPlaceholder: '7777',
      queryPort: '查询端口',
      queryPortRequired: '查询端口 *',
      queryPortPlaceholder: '27015',
      rconPort: 'RCON端口',
      rconPortRequired: 'RCON端口 *',
      rconPortPlaceholder: '32330',
      map: '地图',
      mapPlaceholder: '选择地图',
      maxPlayers: '最大玩家数',
      maxPlayersPlaceholder: '70',
      maxPlayersDesc: '服务器最大玩家数量（1-200）',
      modIds: '模组ID',
      modIdsPlaceholder: '输入模组ID，多个用逗号分隔（如：123456,789012）',
      modIdsDesc: 'Steam创意工坊模组ID，多个模组用逗号分隔',
      adminPassword: '管理员密码',
      adminPasswordRequired: '管理员密码 *',
      adminPasswordPlaceholder: '输入管理员密码（同时作为RCON密码）',
      showPassword: '显示密码',
      hidePassword: '隐藏密码',
      saveChanges: '保存更改',
      createServer: '创建服务器',
      saving: '保存中...',
      preparing: '准备中...',
      loadingServerInfo: '加载服务器信息中...',
      closeConfirm: '确定要关闭吗？未保存的数据将会丢失。',
      // 地图选项
      maps: {
        TheIsland: '孤岛',
        TheCenter: '中心岛',
        ScorchedEarth_P: '焦土',
        Aberration_P: '畸变',
        Extinction: '灭绝',
        Valguero_P: '瓦尔盖罗',
        Genesis: '创世纪',
        Genesis2: '创世纪2',
        CrystalIsles: '水晶群岛',
        LostIsland: '失落岛',
        Fjordur: '冰封群岛'
      },
      selectMapPlaceholder: '选择地图或输入地图名称',
      searchMaps: '搜索地图',
      officialMaps: '官方地图',
      customMapPlaceholder: '输入自定义地图名称',
      customMap: '自定义地图',
      noMatchingMaps: '没有匹配的地图',
      mapId: '地图ID'
    },
    // 编辑器相关
    editor: {
      visualEdit: '可视化编辑',
      textEdit: '文本编辑',
      visualEditMode: '可视化编辑模式',
      textEditMode: '文本编辑',
      syncing: '正在同步...',
      synced: '已同步',
      pendingSync: '待同步',
      resetToDefault: '重置为默认',
      format: '格式化',
      content: '内容',
      placeholder: '输入配置内容...',
      description: '此文件包含服务器的基本设置，如端口、密码、最大玩家数等',
      visualEditModeDesc: '可视化编辑模式Desc',
      visualEditModeTip: '通过表单控件修改参数，鼠标悬停在参数名称旁的图标可查看详细说明。',
      gameIniTextEditDesc: '直接编辑 Game.ini 配置文件内容。',
      showPassword: '显示密码',
      hidePassword: '隐藏密码',
      enabled: '启用',
      disabled: '禁用',
      parametersCount: '个参数',
      defaultValue: '默认值',
      parseGameIniError: '解析Game.ini文本失败',
      syncVisualToTextError: '同步可视化配置到文本失败',
      range: '范围',
      syncTip: {
        visual: '在可视化模式下修改参数，切换到文本模式时会将配置同步到文本。',
        text: '在文本模式下直接编辑配置文件，切换到可视化模式时会解析文本内容。'
      }
    },
    // 启动参数编辑器相关
    argsEditor: {
      title: '启动参数配置',
      switchParams: '开关参数',
      numberParams: '数值参数',
      textParams: '文本参数',
      selectParams: '选择参数',
      range: '范围',
      pleaseSelect: '请选择',
      customArgs: '自定义参数',
      customArgsDesc: '添加自定义启动参数，这些参数将直接添加到启动命令中',
      addCustomArg: '添加自定义参数',
      removeCustomArg: '删除自定义参数',
      customArgPlaceholder: '输入自定义启动参数，如 -nosteam',
      enabled: '启用',
      disabled: '禁用'
    },
    // 启动参数分类
    paramCategories: {
      basic: '基础设置',
      core: '核心游戏',
      dinos: '恐龙管理',
      structures: '建筑系统',
      pvp: 'PvP设置',
      mechanics: '游戏机制',
      transfer: '数据传输',
      performance: '性能优化',
      graphics: '图形优化',
      security: '安全设置',
      logging: '日志记录',
      mods: '模组平台',
      features: '游戏功能',
      maintenance: '维护工具',
      advanced: '高级配置',
      custom: '自定义参数',
      // 新增缺失的分类
      events: '事件设置',
      network: '网络设置',
      compatibility: '兼容性设置',
      localization: '本地化设置',
      notifications: '通知设置',
      voice: '语音设置',
      gameplay: '游戏玩法'
    },
    // 查询参数翻译
    queryParams: {
      // 核心游戏设置
      OverrideOfficialDifficulty: '覆盖官方难度等级',
      AllowThirdPersonPlayer: '启用第三人称视角',
      ServerCrosshair: '禁用准星',
      ShowFloatingDamageText: '启用RPG风格伤害弹出文本',
      MapPlayerLocation: '在地图上显示玩家精确位置',
      DisableWeatherFog: '禁用雾天',
      // 恐龙管理
      MaxPersonalTamedDinos: '每个部落的恐龙驯服上限',
      MaxTamedDinos: '服务器驯服恐龙总数上限',
      PersonalTamedDinosSaddleStructureCost: '平台鞍占用恐龙槽位数量',
      MinimumDinoReuploadInterval: '恐龙重新上传冷却时间（秒）',
      ServerAutoForceRespawnWildDinosInterval: '强制野生恐龙重生间隔',
      AllowAnyoneBabyImprintCuddle: '允许任何人照顾幼龙',
      DisableImprintDinoBuff: '禁用恐龙印记玩家状态加成',
      PreventMateBoost: '禁用恐龙配偶加成',
      AllowFlyerCarryPvE: 'PvE服务器允许飞行生物携带野生恐龙',
      AllowFlyingStaminaRecovery: '飞行生物在玩家站立时恢复耐力',
      ForceFlyerExplosives: '允许飞行生物携带C4飞行',
      AllowMultipleAttachedC4: '允许在恐龙上附加多个C4',
      AutoDestroyDecayedDinos: '自动销毁已腐烂的可认领恐龙',
      // 建筑系统
      TheMaxStructuresInRange: '服务器最大建筑数量',
      ForceAllStructureLocking: '默认锁定所有建筑',
      EnableExtraStructurePreventionVolumes: '在特定资源丰富区域禁用建筑',
      OverrideStructurePlatformPrevention: '允许在平台鞍上建造炮塔',
      DestroyUnconnectedWaterPipes: '自动销毁未连接的水管',
      FastDecayUnsnappedCoreStructures: '未连接核心建筑快速腐烂',
      OnlyAutoDestroyCoreStructures: '仅自动销毁核心建筑',
      OnlyDecayUnsnappedCoreStructures: '仅未连接核心建筑腐烂',
      // PvP设置
      PreventOfflinePvP: '启用离线PvP保护',
      PreventOfflinePvPInterval: '离线PvP保护等待时间（秒）',
      PreventTribeAlliances: '禁用部落联盟',
      // 游戏机制
      AllowCrateSpawnsOnTopOfStructures: '允许空投箱出现在建筑上',
      RandomSupplyCratePoints: '空投箱随机位置',
      NonPermanentDiseases: '疾病非永久性（重生后消失）',
      PreventDiseases: '完全禁用疾病',
      PreventSpawnAnimations: '禁用重生动画',
      UseOptimizedHarvestingHealth: '优化收获健康值',
      ClampItemSpoilingTimes: '限制物品腐烂时间到最大值',
      ClampItemStats: '启用物品属性限制',
      ExtinctionEventTimeInterval: '灭绝事件时间间隔（秒）',
      OxygenSwimSpeedStatMultiplier: '氧气游泳速度倍数',
      // 数据传输控制
      PreventDownloadSurvivors: '禁止下载幸存者数据',
      PreventDownloadItems: '禁止下载物品数据',
      PreventDownloadDinos: '禁止下载恐龙数据',
      PreventUploadSurvivors: '禁止上传幸存者数据',
      PreventUploadItems: '禁止上传物品数据',
      PreventUploadDinos: '禁止上传恐龙数据',
      // 高级配置
      CustomLiveTuningUrl: '自定义实时调优文件URL',
      customdynamicconfigurl: '自定义动态配置文件URL',
      EnableFullDump: '启用完整服务器日志转储',
      TribeLogDestroyedEnemyStructures: '在部落日志中显示敌方建筑破坏'
    },
    // 命令行参数翻译
    commandLineArgs: {
      // 性能优化
      allcores: '使用所有CPU核心',
      USEALLAVAILABLECORES: '强制使用所有可用CPU核心',
      lowmemory: '低内存模式（节省800MB RAM）',
      nomemorybias: '减少内存使用（约600MB）',
      usecache: '使用缓存（提升70%加载速度）',
      high: '高优先级进程运行',
      nocombineclientmoves: '禁用客户端移动优化',
      // 图形优化
      allowansel: '启用NVIDIA Ansel支持',
      d3d10: '强制使用DX10（仅Windows）',
      sm4: '强制使用Shader Model 4（仅Windows）',
      nomansky: '禁用天空细节（云层、星空）',
      noaafonts: '移除字体抗锯齿',
      // 安全设置
      NoBattlEye: '禁用BattleEye反作弊',
      noantispeedhack: '禁用反速度黑客检测',
      insecure: '禁用VAC系统',
      UseItemDupeCheck: '启用额外复制保护',
      exclusivejoin: '白名单模式',
      // 日志记录
      servergamelog: '启用服务器游戏日志',
      servergamelogincludetribelogs: '在游戏日志中包含部落日志',
      ServerRCONOutputTribeLogs: '在RCON中显示部落聊天',
      gameplaylogging: '游戏玩法日志记录',
      webalarm: '启用Web警报',
      // 模组和平台
      automanagedmods: '自动管理MOD下载/安装/更新',
      crossplay: '启用跨平台（EPIC和Steam）',
      epiconly: '仅允许Epic玩家连接',
      UseVivox: '启用Vivox语音',
      MapModID: '自定义地图MOD ID',
      PublicIPForEpic: 'Epic客户端的公共IP地址',
      // 游戏功能
      ForceAllowCaveFlyers: '强制允许飞行生物进入洞穴',
      ForceRespawnDinos: '启动时销毁所有野生生物',
      NoBiomeWalls: '消除生物群落变化区域墙效果',
      PreventHibernation: '防止休眠（消耗性能）',
      StasisKeepControllers: '在静止状态保持AI控制器',
      EnableIdlePlayerKick: '踢出空闲玩家',
      nofishloot: '禁用钓鱼非肉类战利品',
      ActiveEvent: '启用指定事件',
      // 维护工具
      ClearOldItems: '清除旧物品（一次性）',
      noninlinesaveload: '修复大型存档损坏问题',
      oldsaveformat: '使用旧存档格式',
      StructureDestructionTag: '结构销毁标签',
      culture: '语言代码覆盖'
    },

    // Game.ini 参数分类
    gameIniCategories: {
      gameBasic: '基础设置',
      experienceSettings: '经验值和等级',
      breedingSettings: '繁殖',
      itemSettings: '物品和资源',
      dinoSettings: '恐龙',
      tribeSettings: '部落和玩家',
      pvpSettings: 'PvP',
      structureSettings: '建筑和结构',
      advancedSettings: '高级功能',
      customSettings: '自定义配置'
    },

    // GameUserSettings.ini 参数分类
    gameUserSettingsCategories: {
      serverBasic: '基本设置',
      gameMode: '游戏模式',
      communication: '聊天和通讯',
      gameMultipliers: '游戏倍率',
      characterSettings: '角色',
      dinoSettings: '恐龙',
      environmentSettings: '环境',
      structureSettings: '建筑',
      tribeSettings: '部落和联盟',
      breedingSettings: '繁殖和印记',
      itemSettings: '物品和补给',
      performanceSettings: '服务器性能',
      diseaseSettings: '疾病和状态',
      offlineRaidSettings: '离线突袭保护',
      crossArkSettings: '跨服传输',
      flyerSettings: '飞行载具',
      advancedSettings: '高级功能'
    },

    // Game.ini 参数翻译
    gameIniParams: {
      // 基础设置
      bUseSingleplayerSettings: '使用单人游戏设置',
      bDisableStructurePlacementCollision: '禁用建筑放置碰撞检测',
      bAllowFlyerCarryPvE: 'PvE模式允许飞行生物抓取',
      bDisableStructureDecayPvE: '禁用PvE建筑腐朽',
      bAllowUnlimitedRespecs: '允许无限重置技能点',
      bAllowPlatformSaddleMultiFloors: '允许平台鞍多层建筑',
      bPassiveDefensesDamageRiderlessDinos: '被动防御伤害无骑手恐龙',
      bPvEDisableFriendlyFire: 'PvE禁用友军伤害',
      bDisableFriendlyFire: '禁用友军伤害',
      bEnablePvPGamma: '启用PvP伽马调节',
      DifficultyOffset: '难度偏移值',
      OverrideOfficialDifficulty: '覆盖官方难度',

      // 经验值和等级设置
      XPMultiplier: '经验值倍率',
      PlayerCharacterWaterDrainMultiplier: '玩家口渴消耗倍率',
      PlayerCharacterFoodDrainMultiplier: '玩家饥饿消耗倍率',
      PlayerCharacterStaminaDrainMultiplier: '玩家耐力消耗倍率',
      PlayerCharacterHealthRecoveryMultiplier: '玩家生命恢复倍率',

      // 繁殖设置
      MatingIntervalMultiplier: '交配间隔倍率',
      EggHatchSpeedMultiplier: '蛋孵化速度倍率',
      BabyMatureSpeedMultiplier: '幼体成长速度倍率',
      BabyFoodConsumptionSpeedMultiplier: '幼体食物消耗速度倍率',
      BabyCuddleIntervalMultiplier: '幼体照顾间隔倍率',
      BabyCuddleGracePeriodMultiplier: '幼体照顾宽限期倍率',
      BabyCuddleLoseImprintQualitySpeedMultiplier: '幼体失去印记质量速度倍率',

      // 物品和资源设置
      HarvestAmountMultiplier: '采集数量倍率',
      HarvestHealthMultiplier: '资源血量倍率',
      ResourcesRespawnPeriodMultiplier: '资源重生周期倍率',
      ItemStackSizeMultiplier: '物品堆叠数量倍率',
      CropGrowthSpeedMultiplier: '作物生长速度倍率',
      GlobalItemDecompositionTimeMultiplier: '全局物品分解时间倍率',
      GlobalCorpseDecompositionTimeMultiplier: '全局尸体分解时间倍率',

      // 恐龙设置
      TamingSpeedMultiplier: '驯服速度倍率',
      DinoCharacterFoodDrainMultiplier: '恐龙饥饿消耗倍率',
      DinoCharacterStaminaDrainMultiplier: '恐龙耐力消耗倍率',
      DinoCharacterHealthRecoveryMultiplier: '恐龙生命恢复倍率',
      DinoCountMultiplier: '恐龙数量倍率',
      WildDinoCharacterFoodDrainMultiplier: '野生恐龙饥饿消耗倍率',
      WildDinoTorporDrainMultiplier: '野生恐龙眩晕值消耗倍率',

      // 部落和玩家设置
      MaxNumberOfPlayersInTribe: '部落最大玩家数量',
      TribeNameChangeCooldown: '部落改名冷却时间（分钟）',
      bPvEAllowTribeWar: 'PvE允许部落战争',
      bPvEAllowTribeWarCancel: 'PvE允许取消部落战争',

      // PvP设置
      bIncreasePvPRespawnInterval: '启用PvP重生间隔递增',
      IncreasePvPRespawnIntervalCheckPeriod: 'PvP重生间隔检查周期（秒）',
      IncreasePvPRespawnIntervalMultiplier: 'PvP重生间隔递增倍率',
      IncreasePvPRespawnIntervalBaseAmount: 'PvP重生间隔基础时间（秒）',

      // 建筑和结构设置
      StructureDamageMultiplier: '建筑伤害倍率',
      StructureResistanceMultiplier: '建筑抗性倍率',
      StructureDamageRepairCooldown: '建筑损伤修复冷却时间（秒）',
      PvEStructureDecayPeriodMultiplier: 'PvE建筑腐朽周期倍率',
      MaxStructuresInRange: '范围内最大建筑数量',

      // 高级功能设置
      bAutoPvETimer: '启用自动PvE计时器',
      bAutoPvEUseSystemTime: '自动PvE使用系统时间',
      AutoPvEStartTimeSeconds: '自动PvE开始时间（秒）',
      AutoPvEStopTimeSeconds: '自动PvE结束时间（秒）',
      bOnlyAllowSpecifiedEngrams: '仅允许指定印痕',
      bAutoUnlockAllEngrams: '自动解锁所有印痕',
      bShowCreativeMode: '显示创造模式',
      bUseCorpseLocator: '使用尸体定位器',
      bDisableLootCrates: '禁用战利品箱',
      bDisableDinoRiding: '禁用恐龙骑乘',
      bDisableDinoTaming: '禁用恐龙驯服',
      bAllowCustomRecipes: '允许自定义配方',

      // 自定义配置
      DayCycleSpeedScale: '昼夜循环速度倍率',
      NightTimeSpeedScale: '夜晚时间速度倍率',
      DayTimeSpeedScale: '白天时间速度倍率'
    },

    // GameUserSettings.ini 参数翻译
    gameUserSettingsParams: {
      // 服务器基本设置
      ServerPassword: '服务器密码',
      SpectatorPassword: '观察者密码',
      AdminLogging: '管理员日志',

      // 游戏模式设置
      serverPVE: 'PvE模式',
      serverHardcore: '硬核模式',
      ShowMapPlayerLocation: '显示玩家位置',
      allowThirdPersonPlayer: '允许第三人称',
      ServerCrosshair: '显示准星',
      EnablePvPGamma: 'PvP伽马调节',
      DisablePvEGamma: '禁用PvE伽马调节',
      serverForceNoHud: '强制隐藏HUD',
      ShowFloatingDamageText: '显示浮动伤害文字',
      AllowHitMarkers: '允许命中标记',

      // 聊天和通讯设置
      globalVoiceChat: '全局语音聊天',
      proximityChat: '附近聊天',
      alwaysNotifyPlayerJoined: '总是通知玩家加入',
      alwaysNotifyPlayerLeft: '总是通知玩家离开',
      DontAlwaysNotifyPlayerJoined: '禁用玩家加入通知',

      // 游戏倍率设置
      XPMultiplier: '经验倍率',
      TamingSpeedMultiplier: '驯服速度倍率',
      HarvestAmountMultiplier: '采集倍率',
      HarvestHealthMultiplier: '资源血量倍率',
      ResourcesRespawnPeriodMultiplier: '资源重生倍率',
      ItemStackSizeMultiplier: '物品堆叠倍率',

      // 角色设置
      PlayerCharacterHealthRecoveryMultiplier: '玩家回血倍率',
      PlayerCharacterFoodDrainMultiplier: '玩家饥饿倍率',
      PlayerCharacterWaterDrainMultiplier: '玩家口渴倍率',
      PlayerCharacterStaminaDrainMultiplier: '玩家耐力倍率',
      PlayerDamageMultiplier: '玩家伤害倍率',
      PlayerResistanceMultiplier: '玩家抗性倍率',
      OxygenSwimSpeedStatMultiplier: '氧气游泳速度倍率',
      ImplantSuicideCD: '植入体自杀冷却',

      // 恐龙设置
      DinoCountMultiplier: '恐龙数量倍率',
      DinoCharacterHealthRecoveryMultiplier: '恐龙回血倍率',
      DinoCharacterFoodDrainMultiplier: '恐龙饥饿倍率',
      DinoCharacterStaminaDrainMultiplier: '恐龙耐力倍率',
      DinoDamageMultiplier: '恐龙伤害倍率',
      TamedDinoDamageMultiplier: '驯服恐龙伤害倍率',
      DinoResistanceMultiplier: '恐龙抗性倍率',
      TamedDinoResistanceMultiplier: '驯服恐龙抗性倍率',
      MaxTamedDinos: '最大驯服恐龙数',
      MaxPersonalTamedDinos: '个人最大驯服数',
      DisableDinoDecayPvE: '禁用PvE恐龙腐朽',
      AutoDestroyDecayedDinos: '自动销毁腐朽恐龙',
      PvEDinoDecayPeriodMultiplier: 'PvE恐龙腐朽倍率',
      PvPDinoDecay: 'PvP恐龙腐朽',
      AllowRaidDinoFeeding: '允许突袭恐龙喂食',
      RaidDinoCharacterFoodDrainMultiplier: '突袭恐龙饥饿倍率',
      AllowFlyerCarryPvE: 'PvE飞行载具抓取',
      bForceCanRideFliers: '强制允许骑乘飞行生物',

      // 环境设置
      DayCycleSpeedScale: '昼夜循环速度',
      DayTimeSpeedScale: '白天时长倍率',
      NightTimeSpeedScale: '夜晚时长倍率',
      DisableWeatherFog: '禁用雾效',
      DifficultyOffset: '难度偏移',
      OverrideOfficialDifficulty: '覆盖官方难度',
      RandomSupplyCratePoints: '随机补给箱位置',

      // 建筑设置
      StructureDamageMultiplier: '建筑伤害倍率',
      StructureResistanceMultiplier: '建筑抗性倍率',
      TheMaxStructuresInRange: '范围内最大建筑数',
      NewMaxStructuresInRange: '新范围内最大建筑数',
      MaxStructuresInRange: '最大建筑范围数',
      DisableStructureDecayPvE: '禁用PvE建筑腐朽',
      PvEStructureDecayPeriodMultiplier: 'PvE建筑腐朽倍率',
      PvEStructureDecayDestructionPeriod: 'PvE建筑腐朽销毁期',
      PvPStructureDecay: 'PvP建筑腐朽',
      StructurePickupTimeAfterPlacement: '建筑快捷拾取时间',
      StructurePickupHoldDuration: '拾取按键持续时间',
      AlwaysAllowStructurePickup: '总是允许拾取建筑',
      OnlyAutoDestroyCoreStructures: '仅自动销毁核心建筑',
      OnlyDecayUnsnappedCoreStructures: '仅腐朽未连接核心建筑',
      FastDecayUnsnappedCoreStructures: '快速腐朽未连接建筑',
      DestroyUnconnectedWaterPipes: '销毁未连接水管',
      StructurePreventResourceRadiusMultiplier: '建筑阻止资源半径倍率',
      MaxPlatformSaddleStructureLimit: '平台鞍最大建筑数',
      PerPlatformMaxStructuresMultiplier: '平台建筑倍率',
      PlatformSaddleBuildAreaBoundsMultiplier: '平台建造区域倍率',
      OverrideStructurePlatformPrevention: '覆盖平台限制',
      EnableExtraStructurePreventionVolumes: '启用额外建筑限制区域',
      AllowCaveBuildingPvE: '允许PvE洞穴建造',
      AllowCaveBuildingPvP: '允许PvP洞穴建造',
      PvEAllowStructuresAtSupplyDrops: 'PvE允许补给点建造',
      AllowCrateSpawnsOnTopOfStructures: '允许补给箱生成在建筑上',
      bAllowPlatformSaddleMultiFloors: '允许平台鞍多层',
      MaxGateFrameOnSaddles: '鞍座最大门框数',

      // 部落和联盟设置
      MaxNumberOfPlayersInTribe: '部落最大人数',
      TribeNameChangeCooldown: '部落改名冷却',
      PreventTribeAlliances: '防止部落联盟',
      MaxAlliancesPerTribe: '每个部落最大联盟数',
      MaxTribesPerAlliance: '每个联盟最大部落数',

      // 繁殖和印记设置
      AllowAnyoneBabyImprintCuddle: '允许任何人照顾幼体',
      DisableImprintDinoBuff: '禁用印记恐龙加成',
      BabyImprintingStatScaleMultiplier: '幼体印记属性倍率',

      // 物品和补给设置
      ClampItemSpoilingTimes: '限制物品腐坏时间',
      ClampResourceHarvestDamage: '限制资源采集伤害',
      UseOptimizedHarvestingHealth: '使用优化采集血量',
      BanListURL: '封禁列表URL',

      // 服务器性能设置
      AutoSavePeriodMinutes: '自动保存间隔',
      KickIdlePlayersPeriod: '踢出空闲玩家时间',
      ListenServerTetherDistanceMultiplier: '监听服务器距离倍率',
      RCONServerGameLogBuffer: 'RCON游戏日志缓冲',
      NPCNetworkStasisRangeScalePlayerCountStart: 'NPC网络停滞范围玩家数起始',
      NPCNetworkStasisRangeScalePlayerCountEnd: 'NPC网络停滞范围玩家数结束',
      NPCNetworkStasisRangeScalePercentEnd: 'NPC网络停滞范围缩放百分比结束',

      // 疾病和状态设置
      PreventDiseases: '防止疾病',
      NonPermanentDiseases: '非永久性疾病',
      PreventSpawnAnimations: '防止生成动画',

      // 离线突袭保护设置
      PreventOfflinePvP: '防止离线PvP',
      PreventOfflinePvPInterval: '离线保护等待时间',

      // 跨服传输设置
      NoTributeDownloads: '禁用贡品下载',
      PreventDownloadSurvivors: '防止下载角色',
      PreventDownloadItems: '防止下载物品',
      PreventDownloadDinos: '防止下载恐龙',
      PreventUploadSurvivors: '防止上传角色',
      PreventUploadItems: '防止上传物品',
      PreventUploadDinos: '防止上传恐龙',
      CrossARKAllowForeignDinoDownloads: '允许外来恐龙下载',
      MaxTributeDinos: '最大贡品恐龙数',
      MaxTributeItems: '最大贡品物品数',
      MinimumDinoReuploadInterval: '最小恐龙重新上传间隔',
      TributeItemExpirationSeconds: '贡品物品过期时间',
      TributeDinoExpirationSeconds: '贡品恐龙过期时间',
      TributeCharacterExpirationSeconds: '贡品角色过期时间',

      // 飞行载具设置
      AllowFlyingStaminaRecovery: '飞行耐力恢复',
      ForceFlyerExplosives: '强制飞行器爆炸物',

      // 高级功能设置
      AllowMultipleAttachedC4: '允许多个C4',
      AllowIntegratedSPlusStructures: '允许集成S+建筑',
      AllowHideDamageSourceFromLogs: '隐藏伤害来源日志',
      AllowSharedConnections: '允许共享连接',
      bFilterTribeNames: '过滤部落名称',
      bFilterCharacterNames: '过滤角色名称',
      bFilterChat: '过滤聊天',
      EnableCryoSicknessPVE: '启用PvE冷冻舱疾病',
      EnableCryopodNerf: '启用冷冻舱削弱',
      CryopodNerfDuration: '冷冻舱削弱持续时间',
      CryopodNerfDamageMult: '冷冻舱削弱伤害倍率',
      CryopodNerfIncomingDamageMultPercent: '冷冻舱削弱受到伤害倍率',
      DisableCryopodEnemyCheck: '禁用冷冻舱敌人检查',
      DisableCryopodFridgeRequirement: '禁用冷冻舱冰箱需求',
      AllowCryoFridgeOnSaddle: '允许鞍座冷冻冰箱',
      MaxTrainCars: '最大火车车厢数',
      MaxHexagonsPerCharacter: '每角色最大六边形数',
      AllowTekSuitPowersInGenesis: '允许创世纪TEK套装能力',
      CustomDynamicConfigUrl: '自定义动态配置URL'
    },

    // GameUserSettings.ini editor description
    gameUserSettingsTextEditDesc: '直接编辑 GameUserSettings.ini 配置文件内容。修改会自动解析并同步到可视化界面，切换到可视化模式可看到解析后的参数设置。',

    // Default values for GameUserSettings parameters
    defaultValues: {
      sessionName: '我的ARK服务器',
      message: '欢迎来到ARK服务器！'
    },

    // Placeholders
    placeholders: {
      gameUserSettings: `[ServerSettings]
SessionName=我的ARK服务器
ServerPassword=
MaxPlayers=70

[SessionSettings]
SessionName=我的ARK服务器

[MessageOfTheDay]
Message=欢迎来到ARK服务器！

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
    privacyPolicy: '隐私政策',
    termsOfService: '服务条款',
    privacyPolicyContent: '我们重视您的隐私。本应用仅收集必要的服务器管理信息，不会泄露您的个人数据。',
    termsOfServiceContent: '使用本ARK服务器管理器即表示您同意遵守相关使用条款。请合理使用本工具进行服务器管理。'
  },

  // 页脚
  footer: {
    copyright: '© {year} ARK服务器管理器 由 {company} 开发',
    privacyPolicy: '隐私政策',
    termsOfService: '服务条款',
    support: '技术支持',
    github: 'GitHub'
  },

  // 错误信息
  errors: {
    networkError: '网络连接错误',
    serverError: '服务器错误',
    unauthorized: '未授权访问',
    forbidden: '禁止访问',
    notFound: '页面不存在',
    validationError: '输入验证失败',
    unknownError: '未知错误',
    tryAgain: '请重试',
    contactSupport: '请联系技术支持'
  },

  // 成功信息
  success: {
    operationSuccess: '操作成功',
    dataSaved: '数据已保存',
    dataDeleted: '数据已删除',
    dataUpdated: '数据已更新',
    dataCreated: '数据已创建'
  },

  // 表单
  form: {
    required: '此字段为必填项',
    invalidFormat: '格式无效',
    minLength: '最少需要 {min} 个字符',
    maxLength: '最多允许 {max} 个字符',
    invalidEmail: '邮箱格式无效',
    invalidUrl: 'URL格式无效',
    invalidNumber: '请输入有效的数字',
    invalidPort: '端口号必须在 1-65535 之间',
    invalidPath: '请输入有效的路径',
    passwordMismatch: '两次输入的密码不一致',
    usernameExists: '用户名已存在',
    serverNameExists: '服务器名称已存在',
    portInUse: '端口已被占用'
  }
}