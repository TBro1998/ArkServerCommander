// ARK Server Startup Parameters Definition - TypeScript version with i18n support
export interface ServerParam {
  type: 'boolean' | 'number' | 'string' | 'select';
  category: string;
  default?: string | number | boolean;
  min?: number;
  max?: number;
  step?: number;
  options?: string[];
  hidden?: boolean;
}

// Query parameters (parameters starting with ?)
export const queryParams: Record<string, ServerParam> = {
  AltSaveDirectoryName: {
    type: 'string',
    category: 'maintenance',
    default: ''
  },
  EventColorsChanceOverride: {
    type: 'number',
    category: 'features',
    default: 0,
    min: 0,
    max: 1,
    step: 0.01
  },
  GameModIds: {
    type: 'string',
    category: 'mods',
    default: ''
  },
  NewYear1UTC: {
    type: 'number',
    category: 'features',
    default: 0
  },
  NewYear2UTC: {
    type: 'number',
    category: 'features',
    default: 0
  }
};

// Command line arguments (parameters starting with -)
export const commandLineArgs: Record<string, ServerParam> = {
  // Events and Features
  ActiveEvent: {
    type: 'string',
    category: 'features',
    default: ''
  },
  NewYearEvent: {
    type: 'boolean',
    category: 'features',
    default: false
  },

  // Mods and Steam
  automanagedmods: {
    type: 'boolean',
    category: 'mods',
    default: false
  },
  MapModID: {
    type: 'string',
    category: 'mods',
    default: ''
  },

  // Crossplay and Network
  crossplay: {
    type: 'boolean',
    category: 'core',
    default: false
  },
  epiconly: {
    type: 'boolean',
    category: 'core',
    default: false
  },
  PublicIPForEpic: {
    type: 'string',
    category: 'core',
    default: ''
  },
  MULTIHOME: {
    type: 'string',
    category: 'core',
    default: ''
  },

  // Server Management
  culture: {
    type: 'select',
    category: 'maintenance',
    default: '',
    options: ['ca', 'cs', 'da', 'de', 'en', 'es', 'eu', 'fi', 'fr', 'hu', 'it', 'ja', 'ka', 'ko', 'nl', 'pl', 'pt_BR', 'ru', 'sv', 'th', 'tr', 'zh', 'zh-Hans-CN', 'zh-TW']
  },
  exclusivejoin: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  EnableIdlePlayerKick: {
    type: 'boolean',
    category: 'maintenance',
    default: false
  },
  MaxNumOfSaveBackups: {
    type: 'number',
    category: 'maintenance',
    default: 20,
    min: 1,
    max: 100
  },
  newsaveformat: {
    type: 'boolean',
    category: 'maintenance',
    default: true
  },
  NoHangDetection: {
    type: 'boolean',
    category: 'maintenance',
    default: false
  },

  // Creatures and Gameplay
  DisableCustomFoldersInTributeInventories: {
    type: 'boolean',
    category: 'mechanics',
    default: false
  },
  ForceAllowCaveFlyers: {
    type: 'boolean',
    category: 'dinos',
    default: false
  },
  ForceRespawnDinos: {
    type: 'boolean',
    category: 'dinos',
    default: false
  },
  NoDinos: {
    type: 'boolean',
    category: 'dinos',
    default: false
  },
  imprintlimit: {
    type: 'number',
    category: 'dinos',
    default: 101,
    min: 0,
    max: 200
  },
  MinimumTimeBetweenInventoryRetrieval: {
    type: 'number',
    category: 'mechanics',
    default: 3600,
    min: 0
  },

  // PvP Settings
  DisableRailgunPVP: {
    type: 'boolean',
    category: 'pvp',
    default: false
  },
  pvedisallowtribewar: {
    type: 'boolean',
    category: 'pvp',
    default: false
  },
  pveallowtribewar: {
    type: 'boolean',
    category: 'pvp',
    default: false
  },

  // Security and Anti-Cheat
  insecure: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  NoBattlEye: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  noantispeedhack: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  speedhackbias: {
    type: 'number',
    category: 'security',
    default: 1.0,
    min: 0.1,
    max: 10.0,
    step: 0.1
  },
  noundermeshchecking: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  noundermeshkilling: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  SecureSendArKPayload: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  UseItemDupeCheck: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  UseSecureSpawnRules: {
    type: 'boolean',
    category: 'security',
    default: false
  },
  BattlEyeServerRecheck: {
    type: 'boolean',
    category: 'security',
    default: false
  },

  // Performance Optimization
  nocombineclientmoves: {
    type: 'boolean',
    category: 'performance',
    default: false
  },
  StasisKeepControllers: {
    type: 'boolean',
    category: 'performance',
    default: false
  },
  structurememopts: {
    type: 'boolean',
    category: 'performance',
    default: false
  },
  UseStructureStasisGrid: {
    type: 'boolean',
    category: 'performance',
    default: false
  },
  DormancyNetMultiplier: {
    type: 'number',
    category: 'performance',
    default: 1.0,
    min: 0.1,
    max: 10.0,
    step: 0.1
  },
  nodormancythrottling: {
    type: 'boolean',
    category: 'performance',
    default: false
  },
  nitradotest2: {
    type: 'boolean',
    category: 'performance',
    default: false
  },
  dedihibernation: {
    type: 'boolean',
    category: 'performance',
    default: false
  },

  // Graphics and Client
  ServerAllowAnsel: {
    type: 'boolean',
    category: 'graphics',
    default: false
  },

  // Logging and Admin
  servergamelog: {
    type: 'boolean',
    category: 'logging',
    default: false
  },
  servergamelogincludetribelogs: {
    type: 'boolean',
    category: 'logging',
    default: false
  },
  ServerRCONOutputTribeLogs: {
    type: 'boolean',
    category: 'logging',
    default: false
  },
  NotifyAdminCommandsInChat: {
    type: 'boolean',
    category: 'logging',
    default: false
  },

  // Transfer and Cluster
  ClusterDirOverride: {
    type: 'string',
    category: 'transfer',
    default: ''
  },
  clusterid: {
    type: 'string',
    category: 'transfer',
    default: ''
  },
  NoTransferFromFiltering: {
    type: 'boolean',
    category: 'transfer',
    default: false
  },
  usestore: {
    type: 'boolean',
    category: 'transfer',
    default: true
  },
  BackupTransferPlayerDatas: {
    type: 'boolean',
    category: 'transfer',
    default: true
  },
  converttostore: {
    type: 'boolean',
    category: 'transfer',
    default: true
  },

  // Communication
  UseVivox: {
    type: 'boolean',
    category: 'features',
    default: false
  },
  webalarm: {
    type: 'boolean',
    category: 'features',
    default: false
  },
  AllowChatSpam: {
    type: 'boolean',
    category: 'features',
    default: false
  },

  // Advanced/Undocumented
  CustomAdminCommandTrackingURL: {
    type: 'string',
    category: 'advanced',
    default: ''
  },
  CustomMerticsURL: {
    type: 'string',
    category: 'advanced',
    default: ''
  },
  CustomNotificationURL: {
    type: 'string',
    category: 'advanced',
    default: ''
  },
  DisableDupeLogDeletes: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  EnableOfficialOnlyVersioningCode: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  EnableVictoryCoreDupeCheck: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  forcedisablemeshchecking: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  ForceDupeLog: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  ignoredupeditems: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  MaxConnectionsPerIP: {
    type: 'number',
    category: 'advanced',
    default: 0,
    min: 0,
    max: 100
  },
  parseservertojson: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  pauseonddos: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  PreventTotalConversionSaveDir: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  ReloadedForBackup: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  UnstasisDinoObstructionCheck: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  UseTameEffectivenessClamp: {
    type: 'boolean',
    category: 'advanced',
    default: false
  },
  UseServerNetSpeedCheck: {
    type: 'boolean',
    category: 'advanced',
    default: false
  }
};

// Parameter categories
export const categories = {
  basic: 'basic',
  core: 'core',
  dinos: 'dinos',
  structures: 'structures',
  pvp: 'pvp',
  mechanics: 'mechanics',
  transfer: 'transfer',
  performance: 'performance',
  graphics: 'graphics',
  security: 'security',
  logging: 'logging',
  mods: 'mods',
  features: 'features',
  maintenance: 'maintenance',
  advanced: 'advanced',
  custom: 'custom'
} as const;

export type CategoryKey = keyof typeof categories;

// Get server startup parameters by category
export function getServerParamsByCategory() {
  const result: Record<CategoryKey, Array<{ key: string, param: ServerParam }>> = {
    basic: [],
    core: [],
    dinos: [],
    structures: [],
    pvp: [],
    mechanics: [],
    transfer: [],
    performance: [],
    graphics: [],
    security: [],
    logging: [],
    mods: [],
    features: [],
    maintenance: [],
    advanced: [],
    custom: []
  };

  // Categorize query parameters
  Object.entries(queryParams).forEach(([key, param]) => {
    if (!param.hidden) {
      result[param.category as CategoryKey].push({ key, param });
    }
  });

  // Categorize command line arguments
  Object.entries(commandLineArgs).forEach(([key, param]) => {
    if (!param.hidden) {
      result[param.category as CategoryKey].push({ key, param });
    }
  });

  return result;
}

// Get default parameter values
export function getDefaultValues() {
  const defaults = {
    query_params: {} as Record<string, string>,
    command_line_args: {} as Record<string, string | number | boolean | undefined>,
    custom_args: [] as string[]
  };

  Object.entries(queryParams).forEach(([key, param]) => {
    if (!param.hidden) {
      // Query parameter values must be strings
      if (param.type === 'boolean') {
        defaults.query_params[key] = param.default ? 'True' : 'False';
      } else {
        defaults.query_params[key] = String(param.default || '');
      }
    }
  });

  Object.entries(commandLineArgs).forEach(([key, param]) => {
    if (!param.hidden) {
      // Command line arguments can keep their original type
      defaults.command_line_args[key] = param.default;
    }
  });

  return defaults;
}