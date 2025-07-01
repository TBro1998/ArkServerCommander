// ARK 配置文件解析工具

/**
 * 解析 INI 格式的配置文件
 * @param {string} content - 配置文件内容
 * @returns {Object} 解析后的配置对象
 */
export function parseIniFile(content) {
  const result = {}
  let currentSection = 'default'
  
  if (!content) return result
  
  const lines = content.split('\n')
  
  lines.forEach(line => {
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
 * @param {Object} config - 配置对象
 * @returns {string} INI 格式的字符串
 */
export function generateIniFile(config) {
  let content = ''
  
  Object.keys(config).forEach(sectionName => {
    if (sectionName !== 'default') {
      content += `[${sectionName}]\n`
    }
    
    const section = config[sectionName]
    Object.keys(section).forEach(key => {
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
 * 转换值类型
 * @param {any} value - 原始值
 * @param {string} type - 目标类型
 * @returns {any} 转换后的值
 */
export function convertValue(value, type) {
  if (value === undefined || value === null) {
    return value
  }
  
  switch (type) {
    case 'boolean':
      if (typeof value === 'string') {
        return value.toLowerCase() === 'true'
      }
      return Boolean(value)
    
    case 'number':
      if (typeof value === 'string') {
        const num = parseFloat(value)
        return isNaN(num) ? 0 : num
      }
      return Number(value)
    
    case 'text':
    case 'password':
    default:
      return String(value)
  }
}

/**
 * 验证配置值
 * @param {any} value - 要验证的值
 * @param {Object} param - 参数定义
 * @returns {boolean} 是否有效
 */
export function validateValue(value, param) {
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
 * @param {Object} paramDef - 参数定义
 * @returns {any} 默认值
 */
export function getDefaultValue(paramDef) {
  return paramDef.default
}

/**
 * 格式化配置文件内容
 * @param {string} content - 原始内容
 * @returns {string} 格式化后的内容
 */
export function formatConfigContent(content) {
  if (!content) return ''
  
  const lines = content.split('\n')
  const formatted = []
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
 * @param {Object} defaultConfig - 默认配置
 * @param {Object} userConfig - 用户配置
 * @returns {Object} 合并后的配置
 */
export function mergeConfigs(defaultConfig, userConfig) {
  const result = { ...defaultConfig }
  
  Object.keys(userConfig).forEach(section => {
    if (!result[section]) {
      result[section] = {}
    }
    
    Object.keys(userConfig[section]).forEach(key => {
      result[section][key] = userConfig[section][key]
    })
  })
  
  return result
}

/**
 * 从现有配置文件提取已设置的值
 * @param {string} content - 配置文件内容
 * @param {Object} paramDefs - 参数定义
 * @returns {Object} 提取的配置值
 */
export function extractConfigValues(content, paramDefs) {
  // 安全检查
  if (!content || typeof content !== 'string') {
    return {}
  }
  
  if (!paramDefs || typeof paramDefs !== 'object') {
    return {}
  }
  
  const parsed = parseIniFile(content)
  const result = {}
  
  // 遍历所有参数定义
  Object.keys(paramDefs).forEach(sectionKey => {
    const section = paramDefs[sectionKey]
    if (!section || !section.params) {
      return
    }
    
    Object.keys(section.params).forEach(paramKey => {
      const param = section.params[paramKey]
      
      // 在解析结果中查找值
      Object.keys(parsed).forEach(parsedSection => {
        if (parsed[parsedSection] && parsed[parsedSection][paramKey] !== undefined) {
          result[paramKey] = convertValue(parsed[parsedSection][paramKey], param.type)
        }
      })
    })
  })
  
  return result
} 