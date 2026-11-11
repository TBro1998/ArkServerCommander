const messages = {
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
    passwordMinLength: '密码至少需要8位字符，最长72字节',
    passwordMinLengthError: '密码至少需要8位',
    passwordMaxLengthError: '密码过长，最多72字节（一个汉字占3字节）',
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
      download: '下载',
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
      noAffectedServers: '没有使用该镜像的服务器',
      runningCount: '其中 {count} 台正在运行',
      updateWarning: '更新须知',
      warningDownloadTime: '镜像较大，下载可能需要较长时间',
      warningContainerRecreate: '服务器会在下次启动时用新镜像重建容器，重建后首次启动会重新校验服务端文件并下载 Mod，耗时较长',
      warningDataSafety: '存档和配置保存在数据卷中，重建容器不会丢失',
      updateOptions: '更新选项',
      updateImageOnly: '仅更新镜像',
      updateImageOnlyDesc: '只下载新镜像。运行中的服务器继续使用旧镜像，下次启动时自动切换到新镜像。',
      updateAndRecreate: '更新镜像并重启运行中的服务器',
      updateAndRecreateDesc: '下载完成后，运行中的服务器会被停止、用新镜像重建容器并重新启动，期间玩家会掉线。',
      confirmUpdate: '确认更新',
      updateAvailable: '有新版本',
      update: '更新',
      upToDate: '所有镜像均为最新版本',
      updatesFound: '发现新版本镜像，可在下方更新',
      checkingUpdates: '正在检查更新...',
      checkUpdatesFailed: '检查更新失败，请确认服务器能访问 Docker Hub',
      updateStarted: '镜像更新已开始',
      updateFailed: '镜像更新失败',
      loadAffectedFailed: '获取受影响的服务器失败',
      loadingStatus: '加载镜像状态中...',
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
    serverUpdateError: '保存服务器配置失败',
    operationFailed: '操作失败，请稍后重试',
    deleteFailed: '删除失败，请稍后重试',
    recreateStarted: '容器重建已开始，完成后服务器处于已停止状态，需要手动启动',
    recreateFailed: '重建容器失败，请稍后重试',
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
      imageOutdated: '旧镜像',
      imageOutdatedTip: '镜像已更新，下次启动时会用新镜像重建容器',
      configPending: '配置待生效',
      configPendingTip: '配置文件修改将在下次启动时应用',
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
      lastError: '上次操作失败',
      recreateServer: '重建容器',
      confirmRecreateMessage: '将先停止运行中的服务器（玩家会掉线），再按当前配置重建容器，完成后需要手动启动。重建后的首次启动会重新校验服务端文件并下载 Mod。',
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
      loadServerInfoFailed: '加载服务器信息失败，请稍后重试',
      serverCreateError: '创建服务器失败',
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
      identifierRequired: '请填写服务器标识',
      configAppliesOnRestart: '服务器未停止：GameUserSettings.ini 和 Game.ini 的修改会先保存，在下次启动服务器时写入并生效（ARK 停服时会覆盖配置文件，所以不能立即写入）。',
      configPending: '有已保存但尚未生效的配置文件修改，将在下次启动服务器时应用。',
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
      mapId: '地图ID',
      range: '范围'
    },
    // 编辑器相关
    editor: {
      placeholders: {
        gameIni: '输入 Game.ini 内容...',
        gameUserSettings: '输入 GameUserSettings.ini 内容...'
      },
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
      basic: '基础',
      core: '核心',
      dinos: '恐龙',
      structures: '建筑',
      pvp: 'PvP',
      mechanics: '游戏机制',
      transfer: '传输与集群',
      performance: '性能',
      graphics: '图形',
      security: '安全',
      logging: '日志',
      mods: '模组',
      features: '功能',
      maintenance: '维护',
      advanced: '高级',
      custom: '自定义'
    },
    // 查询参数翻译
    queryParams: {
      AltSaveDirectoryName: '备用存档目录名称',
      EventColorsChanceOverride: '事件颜色概率覆盖',
      GameModIds: '游戏模组ID',
      NewYear1UTC: '新年事件开始时间 (UTC)',
      NewYear2UTC: '新年事件结束时间 (UTC)'
    },
    // 命令行参数翻译
    commandLineArgs: {
      // 事件和功能
      ActiveEvent: '激活事件',
      NewYearEvent: '新年事件',
      UseVivox: '使用 Vivox 语音聊天',
      webalarm: '网页警报',
      AllowChatSpam: '允许聊天刷屏',

      // 模组和 Steam
      automanagedmods: '自动管理模组',
      MapModID: '地图模组ID',

      // 跨平台和网络
      crossplay: '启用跨平台',
      epiconly: '仅限 Epic 游戏商店',
      PublicIPForEpic: 'Epic 游戏公网IP',
      MULTIHOME: '多宿主IP地址',

      // 服务器管理
      culture: '服务器语言文化',
      exclusivejoin: '独占加入（白名单）',
      EnableIdlePlayerKick: '启用闲置玩家踢出',
      MaxNumOfSaveBackups: '最大存档备份数量',
      newsaveformat: '新存档格式',
      NoHangDetection: '禁用挂起检测',

      // 生物和游戏玩法
      DisableCustomFoldersInTributeInventories: '禁用贡品库存中的自定义文件夹',
      ForceAllowCaveFlyers: '强制允许洞穴飞行',
      ForceRespawnDinos: '强制重生恐龙',
      NoDinos: '无恐龙',
      imprintlimit: '印记限制百分比',
      MinimumTimeBetweenInventoryRetrieval: '库存检索之间的最小时间',

      // PvP 设置
      DisableRailgunPVP: '在PvP中禁用轨道炮',
      pvedisallowtribewar: 'PvE 禁止部落战争',
      pveallowtribewar: 'PvE 允许部落战争',

      // 安全和反作弊
      insecure: '禁用 VAC（不安全）',
      NoBattlEye: '禁用 BattlEye',
      noantispeedhack: '禁用反加速作弊',
      speedhackbias: '加速作弊检测偏差',
      noundermeshchecking: '禁用穿地检测',
      noundermeshkilling: '禁用穿地击杀',
      SecureSendArKPayload: '安全发送 ARK 载荷',
      UseItemDupeCheck: '使用物品复制检查',
      UseSecureSpawnRules: '使用安全生成规则',
      BattlEyeServerRecheck: 'BattlEye 服务器重检',

      // 性能优化
      nocombineclientmoves: '禁用合并客户端移动',
      StasisKeepControllers: '静止保持控制器',
      structurememopts: '结构内存优化',
      UseStructureStasisGrid: '使用结构静止网格',
      DormancyNetMultiplier: '休眠网络倍数',
      nodormancythrottling: '禁用休眠节流',
      nitradotest2: 'Nitrado 测试模式 2',
      dedihibernation: '专用休眠',

      // 图形和客户端
      ServerAllowAnsel: '服务器允许 NVIDIA Ansel',

      // 日志和管理
      servergamelog: '服务器游戏日志',
      servergamelogincludetribelogs: '服务器游戏日志包含部落日志',
      ServerRCONOutputTribeLogs: '服务器 RCON 输出部落日志',
      NotifyAdminCommandsInChat: '在聊天中通知管理员命令',

      // 传输和集群
      ClusterDirOverride: '集群目录覆盖',
      clusterid: '集群ID',
      NoTransferFromFiltering: '无传输来源过滤',
      usestore: '使用备份',
      BackupTransferPlayerDatas: '备份传输玩家数据',
      converttostore: '转换为备份',

      // 高级/未记录
      CustomAdminCommandTrackingURL: '自定义管理员命令跟踪URL',
      CustomMerticsURL: '自定义指标URL',
      CustomNotificationURL: '自定义通知URL',
      DisableDupeLogDeletes: '禁用重复日志删除',
      EnableOfficialOnlyVersioningCode: '启用仅官方版本控制代码',
      EnableVictoryCoreDupeCheck: '启用胜利核心复制检查',
      forcedisablemeshchecking: '强制禁用网格检查',
      ForceDupeLog: '强制复制日志',
      ignoredupeditems: '忽略重复物品',
      MaxConnectionsPerIP: '每个IP的最大连接数',
      parseservertojson: '解析服务器为JSON',
      pauseonddos: 'DDoS时暂停',
      PreventTotalConversionSaveDir: '防止全面转换存档目录',
      ReloadedForBackup: '为备份重新加载',
      UnstasisDinoObstructionCheck: '解除静止恐龙阻塞检查',
      UseTameEffectivenessClamp: '使用驯服效果钳制',
      UseServerNetSpeedCheck: '使用服务器网络速度检查'
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
  },

  // 布局
  layout: {
    initializing: '初始化中...',
    verifying: '验证身份中...',
    menu: '菜单'
  },

  // 后端错误码对应的提示（中文界面直接显示后端返回的详细信息，这里主要给其他语言使用）
  apiErrors: {
    bad_request: '请求参数错误',
    backend_unavailable: '后端服务不可用',
    invalid_server_id: '无效的服务器ID',
    server_not_found: '服务器不存在',
    identifier_exists: '服务器标识已存在',
    server_running: '无法删除正在运行的服务器，请先停止服务器',
    already_running: '服务器已在运行中',
    already_starting: '服务器正在启动中',
    already_stopped: '服务器已经停止',
    already_stopping: '服务器正在停止中',
    image_pulling: '镜像正在下载中，请等待完成',
    port_conflict: '端口与其他服务器冲突',
    invalid_port: '端口无效',
    server_busy: '服务器正在启动或停止中，请稍后再试',
    image_not_allowed: '不允许操作该镜像',
    invalid_cluster_id: '集群ID格式无效',
    invalid_mod_ids: 'Mod ID 只能是逗号分隔的数字',
    server_not_running: '服务器未运行',
    server_not_stopped: '请先停止服务器',
    invalid_command: '命令无效',
    backup_not_found: '备份不存在',
    invalid_backup: '备份文件格式无效',
    backup_busy: '该服务器正在备份或恢复，请稍后再试',
    schedule_not_found: '定时任务不存在',
    invalid_schedule: '定时任务配置无效',
    rcon_unavailable: 'RCON 端口没有响应，服务器可能还在下载、安装 Mod 或加载地图',
    rcon_auth_failed: 'RCON 认证失败，重启服务器后会重新写入管理员密码',
    too_many_attempts: '登录失败次数过多，请稍后再试',
    invalid_credentials: '用户名或密码错误',
    invalid_password_length: '密码长度需要在 8 到 72 个字节之间',
    wrong_password: '当前密码错误'
  },

  // 服务器状态
  status: {
    running: '运行中',
    stopped: '已停止',
    starting: '启动中',
    stopping: '停止中',
    restarting: '重启中',
    unknown: '未知'
  },

  // 运行阶段
  phase: {
    downloading: '下载服务端',
    downloadingProgress: '下载服务端 {progress}%',
    installing_mods: '安装 Mod',
    installingModsProgress: '安装 Mod {current}/{total}',
    launching: '加载中',
    ready: '可连接',
    downloadingDesc: '正在下载或校验服务端文件。首次启动需要下载约 20 GB，耗时取决于网络速度。',
    installing_modsDesc: '正在下载并安装 Mod，已安装的 Mod 会跳过。',
    launchingDesc: '游戏进程已启动，正在加载地图和 Mod，通常需要几分钟，期间还不能连接。',
    readyDesc: '服务器已就绪，玩家可以连接。',
    unknownDesc: '正在判断服务器状态...'
  },

  // 服务器列表
  serverList: {
    summary: '共 {total} 台，{running} 台运行中，{players} 名玩家在线',
    refresh: '刷新',
    imageNotReadyGoHome: '启动服务器需要的镜像还没有下载，请到首页的"镜像管理"中下载。'
  },

  // 服务器卡片
  serverCard: {
    passwordFailed: '获取管理员密码失败',
    passwordCopied: '管理员密码已复制',
    copyFailed: '复制失败，请手动复制',
    imageOutdated: '旧镜像',
    imageOutdatedTip: '镜像已更新，下次启动时会用新镜像重建容器',
    configPending: '配置待生效',
    configPendingTip: '配置文件修改将在下次启动时应用',
    cluster: '集群 {id}',
    clusterTip: '与相同集群ID的服务器共享角色和物品传送数据',
    lastError: '上次操作失败',
    downloadHint: '首次启动需要下载服务端文件，请耐心等待',
    launchingHint: '正在加载地图，通常需要几分钟',
    players: '在线玩家',
    memory: '内存',
    map: '地图',
    ports: '端口（游戏/查询/RCON）',
    adminPassword: '管理员密码',
    showPassword: '显示密码',
    hidePassword: '隐藏密码',
    copyPassword: '复制密码'
  },

  // 服务器操作
  serverActions: {
    start: '启动',
    stop: '停止',
    restart: '重启',
    working: '处理中',
    more: '更多操作',
    details: '详情',
    logs: '查看日志',
    edit: '编辑配置',
    updateGame: '更新游戏/Mod',
    recreate: '重建容器',
    delete: '删除服务器',
    deleteNeedStop: '删除（需先停止）',
    imageNotReady: '镜像未就绪，无法启动',
    startSent: '已开始启动',
    startFailed: '启动服务器失败',
    stopSent: '已开始停止，停服前会保存世界',
    stopFailed: '停止服务器失败',
    restartSent: '已开始重启',
    restartFailed: '重启服务器失败',
    recreateSent: '容器重建已开始，完成后服务器处于已停止状态',
    recreateFailed: '重建容器失败',
    updateSent: '更新已开始，服务器启动时会下载更新',
    updateFailed: '更新失败',
    deleteSuccess: '服务器已删除',
    deleteFailed: '删除服务器失败',
    stopTitle: '停止 {name}？',
    stopDesc: '停服前会先保存世界。',
    stopWithPlayers: '当前有 {count} 名玩家在线，停止后他们会掉线。停服前会先保存世界。',
    restartTitle: '重启 {name}？',
    restartDesc: '会先保存世界并停止服务器，再按当前配置启动（配置有变化时会重建容器）。',
    restartStoppedDesc: '服务器已停止，将直接按当前配置启动。',
    playersOnline: '当前有 {count} 名玩家在线，他们会掉线。',
    updateTitle: '更新 {name}',
    updateDesc: '服务器会启动，并在启动时下载更新。',
    updateDescRunning: '服务器会先保存并停止，再启动并在启动时下载更新。',
    updateConfirm: '开始更新',
    updateGameOption: '更新游戏本体',
    updateGameOptionDesc: '用 steamcmd 校验并更新服务端文件，客户端更新后服务器版本不一致时需要执行',
    updateModsOption: '更新 Mod',
    updateModsOptionDesc: '重新下载全部 Mod',
    noMods: '没有配置 Mod',
    recreateTitle: '重建 {name} 的容器？',
    recreateDesc: '将先停止运行中的服务器（玩家会掉线），再按当前配置重建容器，完成后需要手动启动。存档和服务端文件保存在卷中，不会丢失。',
    deleteTitle: '删除 {name}？',
    deleteDesc: '服务器记录和容器会被删除。',
    deleteData: '同时删除全部数据',
    deleteDataOff: '保留数据卷（存档、服务端文件）和备份，可以之后手动清理',
    deleteDataOn: '存档、配置、服务端文件、插件和全部备份都会被永久删除，无法恢复',
    typeNameToConfirm: '输入服务器名称 "{name}" 确认'
  },

  // 连接信息
  connection: {
    address: '连接地址',
    steamConnect: 'Steam 连接',
    copyAddress: '复制连接地址',
    copyLink: '复制 Steam 链接',
    copied: '已复制 {text}',
    copyFailed: '复制失败，请手动复制',
    hint: '在游戏控制台输入 open 连接地址，或点击 Steam 链接加入。地址来自 PUBLIC_HOST 或自动探测的公网 IP，需要在路由器/防火墙放行端口。'
  },

  // Mod
  mods: {
    invalid: 'Mod ID 只能是逗号分隔的数字',
    loading: '查询名称中...',
    notFound: '创意工坊中找不到该 Mod',
    unavailable: '无法连接 Steam，暂时查不到名称'
  },

  // 服务器详情
  detail: {
    tabs: {
      overview: '概览',
      console: '控制台',
      logs: '日志',
      backups: '备份',
      schedules: '定时任务'
    },
    backToList: '返回服务器列表',
    notFound: '服务器不存在或已被删除',
    runtime: '运行状态',
    collecting: '正在采集运行信息...',
    notRunning: '服务器未运行',
    players: '在线玩家',
    memory: '内存',
    advertisedName: '服务器列表中的名称',
    connection: '连接信息',
    noPublicHost: '无法获取公网地址，可在环境变量 PUBLIC_HOST 中设置',
    gamePort: '游戏端口',
    queryPort: '查询端口',
    rconPort: 'RCON 端口',
    config: '配置',
    identifier: '服务器标识',
    map: '地图',
    maxPlayers: '最大玩家数',
    cluster: '集群ID',
    none: '无',
    autoRestart: '崩溃后自动重启',
    autoUpdate: '启动时自动更新',
    on: '开启',
    off: '关闭',
    createdAt: '创建时间',
    mods: 'Mod',
    noMods: '没有配置 Mod',
    imageOutdated: '旧镜像',
    configPending: '配置待生效',
    lastError: '上次操作失败',
    viewLogs: '查看日志'
  },

  // 日志
  logs: {
    connecting: '连接中...',
    live: '实时',
    paused: '已暂停',
    ended: '已结束',
    error: '连接失败',
    filter: '筛选日志...',
    pause: '暂停',
    resume: '继续',
    reconnect: '重新连接',
    clear: '清空',
    download: '下载日志',
    downloadFailed: '下载日志失败',
    empty: '暂无日志（容器尚未创建或刚刚启动）',
    scrollToBottom: '回到底部',
    hint: '这里是容器的输出：服务端下载进度、Mod 安装和游戏启动信息。游戏本身的详细日志在存档目录的 Logs 下。'
  },

  // 控制台
  console: {
    notRunning: '服务器未运行，启动后才能使用控制台。',
    notReady: '服务器还在准备中（下载、安装 Mod 或加载地图），RCON 命令可能暂时无法执行。',
    quickActions: '快捷操作',
    broadcast: '广播',
    broadcastPlaceholder: '输入要广播给所有玩家的消息',
    broadcastSent: '广播已发送',
    saveWorld: '保存世界',
    saved: '世界已保存',
    rconConsole: 'RCON 控制台',
    consoleHint: '输入 RCON 命令后回车执行，上下方向键浏览历史命令。',
    commandPlaceholder: '输入 RCON 命令，例如 ListPlayers',
    noOutput: '（无输出）',
    commandFailed: '命令执行失败',
    players: '在线玩家',
    playersFailed: '获取玩家列表失败',
    playersNotReady: '服务器就绪后才能查看玩家。',
    playersUnknown: '点击刷新获取玩家列表',
    noPlayers: '当前没有玩家在线',
    loading: '加载中...',
    steamIdCopied: 'SteamID 已复制',
    kick: '踢出',
    ban: '封禁',
    kickTitle: '踢出 {name}？',
    kickDesc: '玩家会被断开连接，但可以重新加入。',
    banTitle: '封禁 {name}？',
    banDesc: '玩家会被踢出并加入封禁列表，之后无法再加入。可以用 UnbanPlayer 命令解封。',
    kicked: '已踢出 {name}',
    banned: '已封禁 {name}'
  },

  // 备份
  backups: {
    title: '存档备份',
    summary: '共 {count} 个备份，占用 {size}',
    hint: '备份包含存档（SavedArks）和配置文件，不包含日志。运行中备份会先保存世界。恢复只能在服务器停止时进行，恢复前会自动备份当前存档。',
    create: '立即备份',
    upload: '上传备份',
    empty: '还没有备份',
    time: '时间',
    kind: '类型',
    size: '大小',
    actions: '操作',
    kinds: {
      manual: '手动',
      scheduled: '定时',
      'pre-restore': '恢复前',
      uploaded: '上传'
    },
    download: '下载',
    restore: '恢复',
    restoreNeedStop: '恢复（需先停止服务器）',
    delete: '删除',
    loadFailed: '获取备份列表失败',
    created: '备份完成',
    createFailed: '备份失败',
    uploaded: '备份已上传',
    uploadFailed: '上传失败',
    restored: '存档已恢复',
    restoreFailed: '恢复失败',
    deleted: '备份已删除',
    deleteFailed: '删除备份失败',
    restoreTitle: '从备份恢复？',
    restoreDesc: '当前的存档和配置文件会被 {time} 的备份覆盖。恢复前会自动备份当前存档。',
    deleteTitle: '删除备份？',
    deleteDesc: '{time} 的备份将被永久删除。'
  },

  // 定时任务
  schedules: {
    title: '定时任务',
    timezoneHint: '时间按服务器时区（CONTAINER_TZ）计算',
    add: '添加任务',
    edit: '编辑任务',
    empty: '还没有定时任务，可以添加每日重启、定时备份等',
    actions: {
      restart: '重启',
      backup: '备份存档',
      update: '更新游戏/Mod',
      command: 'RCON 命令',
      start: '启动',
      stop: '停止'
    },
    actionDesc: {
      restart: '运行中的服务器会在广播倒计时后保存并重启；已停止的不处理。',
      backup: '备份存档和配置，可以只保留最近几份。',
      update: '广播倒计时后停止，更新游戏本体和 Mod 后启动；已停止的服务器会在下次启动时更新。',
      command: '在运行中的服务器上执行一条 RCON 命令。',
      start: '服务器已停止时启动。',
      stop: '运行中的服务器会在广播倒计时后保存并停止。'
    },
    action: '任务类型',
    frequency: '频率',
    modes: {
      daily: '每天',
      weekly: '每周',
      interval: '每隔几小时',
      custom: '自定义 Cron'
    },
    weekdays: ['周日', '周一', '周二', '周三', '周四', '周五', '周六'],
    weekday: '星期',
    time: '时间',
    intervalHours: '间隔（小时）',
    cronHint: '5 段 Cron 表达式：分 时 日 月 周，例如 30 4 * * 1 表示每周一 4:30',
    warnMinutes: '提前广播（分钟）',
    warnMinutesHint: '执行前在游戏内广播倒计时，0 表示不广播',
    message: '广播内容',
    messagePlaceholder: '留空使用默认内容，{minutes} 会替换为剩余分钟数',
    command: '命令',
    keep: '保留份数',
    keepHint: '只保留最近几份定时备份，0 表示不清理',
    enabled: '启用',
    runNow: '立即执行',
    delete: '删除',
    deleteTitle: '删除定时任务？',
    cancel: '取消',
    save: '保存',
    everyDayAt: '每天 {time}',
    everyWeekAt: '每{day} {time}',
    everyHours: '每 {hours} 小时',
    nextRun: '下次 {time}',
    lastRun: '上次 {time}',
    ok: '成功',
    neverRun: '尚未执行',
    saved: '定时任务已保存',
    saveFailed: '保存定时任务失败',
    loadFailed: '获取定时任务失败',
    runStarted: '已开始执行',
    runFailed: '执行失败',
    deleted: '定时任务已删除',
    deleteFailed: '删除定时任务失败'
  },

  // 服务器表单
  serverForm: {
    sectionBasic: '基本信息',
    sectionNetwork: '端口和密码',
    sectionMods: 'Mod 和集群',
    sectionMaintenance: '运维',
    sessionNameHint: '玩家在服务器列表里看到的名称，可以包含空格和中文',
    identifierHint: '管理用的唯一标识，只在本系统内使用，不显示给玩家',
    gamePortHint: '还会占用 +1 端口',
    queryPortHint: 'Steam 服务器列表使用',
    rconPortHint: '远程管理使用',
    portsFirewallHint: '已自动填入未被其他服务器占用的端口。需要在路由器和防火墙放行游戏端口、游戏端口+1 和查询端口（UDP）。',
    passwordHint: '用于游戏内管理员命令（enablecheats）和 RCON',
    passwordKeepHint: '留空则保持原密码不变',
    passwordKeepPlaceholder: '留空则不修改',
    generate: '随机生成',
    modsHint: '创意工坊 Mod ID，多个用逗号分隔，下方会显示 Mod 名称',
    clusterHint: '填写相同集群ID的服务器之间可以传送角色、恐龙和物品（共享集群存储）。只能包含字母、数字、_ . -',
    autoRestart: '崩溃后自动重启',
    autoRestartHint: '游戏进程意外退出时由 Docker 自动重启',
    autoUpdate: '启动时自动更新',
    autoUpdateHint: '每次启动前校验更新游戏本体并重新下载 Mod，启动会变慢',
    identifierRequired: '请填写服务器标识',
    passwordRequired: '请设置管理员密码',
    passwordNoNewline: '密码不能包含换行',
    portRange: '端口需要在 1 到 65535 之间',
    portDuplicate: '端口重复：游戏端口、游戏端口+1、查询端口和 RCON 端口不能相同',
    maxPlayersRange: '最大玩家数需要在 1 到 200 之间',
    clusterInvalid: '集群ID只能包含字母、数字、_ . -，且以字母或数字开头',
    modsInvalid: 'Mod ID 只能是逗号分隔的数字',
    fixErrors: '请先修正"基本参数"中标红的字段',
    discardConfirm: '有未保存的修改，确定要离开吗？',
    saveAndRestart: '保存并重启',
    created: '服务器已创建',
    savedAndRestarting: '已保存，服务器正在重启',
    savedPending: '已保存，配置会在下次启动时生效'
  },

  // 首页概览
  dashboard: {
    title: '概览',
    publicHost: '公网地址',
    manageServers: '管理服务器',
    servers: '服务器',
    running: '运行中',
    readyCount: '{count} 台可连接',
    players: '在线玩家',
    cpuHint: '100% = 一个核心',
    memory: '内存',
    disk: '磁盘占用',
    recentErrors: '最近的错误',
    imageManagement: '镜像管理',
    diskUsage: '磁盘占用',
    diskHint: 'Docker 计算卷大小较慢，结果缓存 5 分钟',
    recalculate: '重新计算',
    noServers: '还没有服务器',
    server: '服务器',
    files: '服务端文件',
    saves: '存档/插件',
    backups: '备份',
    total: '合计',
    clusterVolume: '集群共享存储 {name}',
    downloadStarted: '已开始下载 {name}',
    downloadFailed: '下载镜像失败',
    diskFailed: '计算磁盘占用失败'
  },

  // 设置
  settings: {
    title: '设置',
    account: '账户',
    username: '用户名',
    changePassword: '修改登录密码',
    oldPassword: '当前密码',
    newPassword: '新密码',
    confirmPassword: '确认新密码',
    passwordHint: '8 到 72 个字节（一个汉字占 3 字节）',
    submit: '修改密码',
    sessionNote: '修改后当前登录保持有效，其他设备上已登录的会话在 24 小时内过期。',
    tooShort: '新密码至少需要 8 个字节',
    tooLong: '新密码最多 72 个字节',
    mismatch: '两次输入的新密码不一致',
    changed: '密码已修改',
    failed: '修改密码失败'
  }
}

export default messages;
