// ARK 配置文件解析工具

/**
 * 配置对象接口
 */
export interface ConfigObject {
  [sectionName: string]: {
    [key: string]: string;
  };
}

/**
 * 参数定义接口
 */
export interface ParamDefinition {
  type: 'boolean' | 'number' | 'text' | 'password';
  default?: any;
  min?: number;
  max?: number;
  description?: string;
}

/**
 * 参数定义集合接口
 */
export interface ParamDefinitions {
  [sectionKey: string]: {
    [paramKey: string]: ParamDefinition;
  };
}

/**
 * 配置值集合接口
 */
export interface ConfigValues {
  [key: string]: any;
}

/**
 * 解析 INI 格式的配置文件
 * @param content - 配置文件内容
 * @returns 解析后的配置对象
 */
export function parseIniFile(content: string): ConfigObject {
  const result: ConfigObject = {}
  let currentSection = 'default'
  
  if (!content) {
    console.warn('parseIniFile: content 为空')
    return result
  }
  
  const lines = content.split('\n')
  
  lines.forEach((line: string, index: number) => {
    const trimmedLine = line.trim()
    
    // 跳过空行和注释
    if (!trimmedLine || trimmedLine.startsWith('#') || trimmedLine.startsWith(';')) {
      return
    }
    
    // 处理段落标题 [SectionName]
    if (trimmedLine.startsWith('[') && trimmedLine.endsWith(']')) {
      currentSection = trimmedLine.slice(1, -1)
      if (!result[currentSection]) {
        result[currentSection] = {}
      }
      return
    }
    
    // 处理键值对
    if (trimmedLine.includes('=')) {
      const [key, ...valueParts] = trimmedLine.split('=')
      const value = valueParts.join('=').trim()
      const cleanKey = key.trim()
      
      if (!result[currentSection]) {
        result[currentSection] = {}
      }
      
      result[currentSection][cleanKey] = value
    }
  })
  
  return result
}

/**
 * 将配置对象转换为 INI 格式
 * @param config - 配置对象
 * @returns INI 格式的字符串
 */
export function generateIniFile(config: ConfigObject): string {
  let content = ''
  
  Object.keys(config).forEach((sectionName: string) => {
    if (sectionName !== 'default') {
      content += `[${sectionName}]\n`
    }
    
    const section = config[sectionName]
    Object.keys(section).forEach((key: string) => {
      const value = section[key]
      if (value !== undefined && value !== null && value !== '') {
        content += `${key}=${value}\n`
      }
    })
    
    content += '\n'
  })
  
  return content.trim()
}

/**
 * 智能合并配置 - 保留用户自定义配置，只更新可视化编辑器管理的参数
 * @param originalContent - 原始配置文件内容
 * @param visualConfig - 可视化编辑器的配置值
 * @param paramDefs - 参数定义（用于确定哪些参数是可视化编辑器管理的）
 * @param defaultSections - 默认的段落名称数组
 * @returns 合并后的配置内容
 */
export function mergeConfigContent(
  originalContent: string, 
  visualConfig: ConfigValues, 
  paramDefs: ParamDefinitions,
  defaultSections: string[] = []
): string {
  // 解析原始配置
  const originalParsed = parseIniFile(originalContent || '')
  
  // 收集所有可视化编辑器管理的参数名
  const managedParams = new Set<string>()
  Object.keys(paramDefs).forEach(sectionKey => {
    const section = paramDefs[sectionKey]
    if (section) {
      Object.keys(section).forEach(paramKey => {
        managedParams.add(paramKey)
      })
    }
  })
  
  // 创建结果配置对象
  const result: ConfigObject = {}
  
  // 首先复制原始配置中的所有内容
  Object.keys(originalParsed).forEach(sectionName => {
    result[sectionName] = { ...originalParsed[sectionName] }
  })
  
  // 确保默认段落存在
  defaultSections.forEach(sectionName => {
    if (!result[sectionName]) {
      result[sectionName] = {}
    }
  })
  
  // 更新可视化编辑器管理的参数
  // 对于每个可视化配置的参数，在适当的段落中更新或添加
  Object.keys(visualConfig).forEach(paramKey => {
    const value = visualConfig[paramKey]
    
    // 只处理有值且是管理的参数
    if (managedParams.has(paramKey) && value !== undefined && value !== null && value !== '') {
      // 查找该参数应该放在哪个段落
      let targetSection: string | null = null
      
      // 首先检查该参数是否已经存在于某个段落中
      Object.keys(result).forEach(sectionName => {
        if (result[sectionName][paramKey] !== undefined) {
          targetSection = sectionName
        }
      })
      
      // 如果参数不存在，使用默认段落（第一个）
      if (!targetSection && defaultSections.length > 0) {
        targetSection = defaultSections[0]
      }
      
      // 如果仍然没有目标段落，创建一个默认段落
      if (!targetSection) {
        targetSection = 'default'
        if (!result[targetSection]) {
          result[targetSection] = {}
        }
      }
      
      // 更新参数值
      result[targetSection][paramKey] = String(value)
    }
  })
  
  // 移除空段落
  Object.keys(result).forEach(sectionName => {
    if (Object.keys(result[sectionName]).length === 0) {
      delete result[sectionName]
    }
  })
  
  return generateIniFileWithComments(result, originalContent)
}

/**
 * 生成带注释的 INI 文件 - 保留原始文件中的注释
 * @param config - 配置对象
 * @param originalContent - 原始文件内容（用于提取注释）
 * @returns 带注释的 INI 格式字符串
 */
export function generateIniFileWithComments(config: ConfigObject, originalContent: string): string {
  if (!originalContent) {
    return generateIniFile(config)
  }
  
  const originalLines = originalContent.split('\n')
  const result: string[] = []
  let currentSection = 'default'
  const processedParams = new Set<string>()
  
  // 第一遍：处理原始文件的结构，保留注释和已存在的参数
  for (let i = 0; i < originalLines.length; i++) {
    const line = originalLines[i]
    const trimmedLine = line.trim()
    
    // 保留空行和注释
    if (!trimmedLine || trimmedLine.startsWith('#') || trimmedLine.startsWith(';')) {
      result.push(line)
      continue
    }
    
    // 处理段落标题
    if (trimmedLine.startsWith('[') && trimmedLine.endsWith(']')) {
      currentSection = trimmedLine.slice(1, -1)
      result.push(line)
      continue
    }
    
    // 处理键值对
    if (trimmedLine.includes('=')) {
      const [key] = trimmedLine.split('=')
      const cleanKey = key.trim()
      
      // 如果新配置中有这个参数，使用新值；否则保留原值
      if (config[currentSection] && config[currentSection][cleanKey] !== undefined) {
        result.push(`${cleanKey}=${config[currentSection][cleanKey]}`)
        processedParams.add(`${currentSection}.${cleanKey}`)
      } else {
        // 保留原始行
        result.push(line)
      }
    } else {
      // 保留无法识别的行
      result.push(line)
    }
  }
  
  // 第二遍：添加新配置中有但原始文件中没有的参数
  Object.keys(config).forEach(sectionName => {
    const section = config[sectionName]
    let sectionAdded = false
    
    Object.keys(section).forEach(paramKey => {
      const paramPath = `${sectionName}.${paramKey}`
      
      if (!processedParams.has(paramPath)) {
        // 如果段落标题还没有添加，先添加段落标题
        if (!sectionAdded) {
          // 检查是否需要添加空行
          if (result.length > 0 && result[result.length - 1].trim() !== '') {
            result.push('')
          }
          if (sectionName !== 'default') {
            result.push(`[${sectionName}]`)
          }
          sectionAdded = true
        }
        
        result.push(`${paramKey}=${section[paramKey]}`)
      }
    })
  })
  
  return result.join('\n')
}

/**
 * 转换值类型
 * @param value - 原始值
 * @param type - 目标类型
 * @returns 转换后的值
 */
export function convertValue(value: any, type: string): any {
  if (value === undefined || value === null) {
    return value
  }
  
  switch (type) {
    case 'boolean':
      if (typeof value === 'string') {
        const lowerValue = value.toLowerCase().trim()
        // 支持多种布尔值表示
        const result = lowerValue === 'true' || lowerValue === '1' || lowerValue === 'yes' || lowerValue === 'on'
        return result
      }
      return Boolean(value)
    
    case 'number':
      if (typeof value === 'string') {
        const num = parseFloat(value.trim())
        return isNaN(num) ? 0 : num
      }
      return Number(value)
    
    case 'text':
    case 'password':
    default:
      return String(value).trim()
  }
}

/**
 * 验证配置值
 * @param value - 要验证的值
 * @param param - 参数定义
 * @returns 是否有效
 */
export function validateValue(value: any, param: ParamDefinition): boolean {
  if (value === undefined || value === null) {
    return true // 允许空值
  }
  
  switch (param.type) {
    case 'number':
      const num = Number(value)
      if (isNaN(num)) return false
      if (param.min !== undefined && num < param.min) return false
      if (param.max !== undefined && num > param.max) return false
      return true
    
    case 'text':
    case 'password':
      return typeof value === 'string'
    
    case 'boolean':
      return typeof value === 'boolean'
    
    default:
      return true
  }
}

/**
 * 获取参数的默认值
 * @param paramDef - 参数定义
 * @returns 默认值
 */
export function getDefaultValue(paramDef: ParamDefinition): any {
  return paramDef.default
}

/**
 * 格式化配置文件内容
 * @param content - 原始内容
 * @returns 格式化后的内容
 */
export function formatConfigContent(content: string): string {
  if (!content) return ''
  
  const lines = content.split('\n')
  const formatted: string[] = []
  let lastWasEmpty = false
  
  for (const line of lines) {
    const trimmed = line.trim()
    
    if (trimmed === '') {
      if (!lastWasEmpty) {
        formatted.push('')
        lastWasEmpty = true
      }
    } else {
      formatted.push(trimmed)
      lastWasEmpty = false
    }
  }
  
  return formatted.join('\n')
}

/**
 * 合并配置
 * @param defaultConfig - 默认配置
 * @param userConfig - 用户配置
 * @returns 合并后的配置
 */
export function mergeConfigs(defaultConfig: ConfigObject, userConfig: ConfigObject): ConfigObject {
  const result: ConfigObject = { ...defaultConfig }
  
  Object.keys(userConfig).forEach((section: string) => {
    if (!result[section]) {
      result[section] = {}
    }
    
    Object.keys(userConfig[section]).forEach((key: string) => {
      result[section][key] = userConfig[section][key]
    })
  })
  
  return result
}

/**
 * 从现有配置文件提取已设置的值
 * @param content - 配置文件内容
 * @param paramDefs - 参数定义
 * @returns 提取的配置值
 */
export function extractConfigValues(content: string, paramDefs: ParamDefinitions): ConfigValues {
  // 安全检查
  if (!content || typeof content !== 'string') {
    console.warn('extractConfigValues: content 为空或不是字符串')
    return {}
  }
  
  if (!paramDefs || typeof paramDefs !== 'object') {
    console.warn('extractConfigValues: paramDefs 为空或不是对象')
    return {}
  }
  
  const parsed = parseIniFile(content)
  const result: ConfigValues = {}
  let foundValuesCount = 0
  
  // 遍历所有参数定义
  Object.keys(paramDefs).forEach((sectionKey: string) => {
    const section = paramDefs[sectionKey]
    if (!section) {
      return
    }
    
    Object.keys(section).forEach((paramKey: string) => {
      const param = section[paramKey]
      
      // 在解析结果中查找值
      Object.keys(parsed).forEach((parsedSection: string) => {
        if (parsed[parsedSection] && parsed[parsedSection][paramKey] !== undefined) {
          const rawValue = parsed[parsedSection][paramKey]
          const convertedValue = convertValue(rawValue, param.type)
          result[paramKey] = convertedValue
          foundValuesCount++
        }
      })
    })
  })
  
  return result
} 