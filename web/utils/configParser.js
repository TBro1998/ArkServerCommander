// ARK 配置文件解析工具

/**
 * 解析 INI 格式的配置文件
 * @param {string} content - 配置文件内容
 * @returns {Object} 解析后的配置对象
 */
export function parseIniFile(content) {
  const result = {}
  let currentSection = 'default'
  
  if (!content) {
    console.warn('parseIniFile: content 为空')
    return result
  }
  
  console.log('parseIniFile: 开始解析 INI 文件')
  console.log('parseIniFile: 内容长度:', content.length)
  
  const lines = content.split('\n')
  console.log('parseIniFile: 总行数:', lines.length)
  
  lines.forEach((line, index) => {
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
      console.log(`parseIniFile: 发现段落 [${currentSection}]`)
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
      console.log(`parseIniFile: 解析键值对 ${cleanKey}=${value} (段落: ${currentSection})`)
    }
  })
  
  console.log('parseIniFile: 解析完成，结果:', result)
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
    console.log(`convertValue: 值为空，返回原值: ${value}`)
    return value
  }
  
  console.log(`convertValue: 转换 ${value} (类型: ${typeof value}) 到 ${type}`)
  
  switch (type) {
    case 'boolean':
      if (typeof value === 'string') {
        const lowerValue = value.toLowerCase().trim()
        // 支持多种布尔值表示
        const result = lowerValue === 'true' || lowerValue === '1' || lowerValue === 'yes' || lowerValue === 'on'
        console.log(`convertValue: 布尔转换 ${value} -> ${lowerValue} -> ${result}`)
        return result
      }
      const boolResult = Boolean(value)
      console.log(`convertValue: 布尔转换 ${value} -> ${boolResult}`)
      return boolResult
    
    case 'number':
      if (typeof value === 'string') {
        const num = parseFloat(value.trim())
        const result = isNaN(num) ? 0 : num
        console.log(`convertValue: 数字转换 ${value} -> ${result}`)
        return result
      }
      const numResult = Number(value)
      console.log(`convertValue: 数字转换 ${value} -> ${numResult}`)
      return numResult
    
    case 'text':
    case 'password':
    default:
      const strResult = String(value).trim()
      console.log(`convertValue: 字符串转换 ${value} -> ${strResult}`)
      return strResult
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
    console.warn('extractConfigValues: content 为空或不是字符串')
    return {}
  }
  
  if (!paramDefs || typeof paramDefs !== 'object') {
    console.warn('extractConfigValues: paramDefs 为空或不是对象')
    return {}
  }
  
  console.log('extractConfigValues: 开始解析配置文件')
  console.log('extractConfigValues: 配置文件内容长度:', content.length)
  const parsed = parseIniFile(content)
  console.log('extractConfigValues: 解析后的 INI 结构:', parsed)
  
  const result = {}
  let foundValuesCount = 0
  
  // 遍历所有参数定义
  Object.keys(paramDefs).forEach(sectionKey => {
    const section = paramDefs[sectionKey]
    if (!section || !section.params) {
      console.log(`extractConfigValues: 跳过无效的 section: ${sectionKey}`)
      return
    }
    
    Object.keys(section.params).forEach(paramKey => {
      const param = section.params[paramKey]
      console.log(`extractConfigValues: 查找参数 ${paramKey} (类型: ${param.type})`)
      
      // 在解析结果中查找值
      Object.keys(parsed).forEach(parsedSection => {
        if (parsed[parsedSection] && parsed[parsedSection][paramKey] !== undefined) {
          const rawValue = parsed[parsedSection][paramKey]
          const convertedValue = convertValue(rawValue, param.type)
          result[paramKey] = convertedValue
          foundValuesCount++
          console.log(`extractConfigValues: 找到参数 ${paramKey} = ${rawValue} (转换为 ${convertedValue})`)
        }
      })
    })
  })
  
  console.log(`extractConfigValues: 完成解析，找到 ${foundValuesCount} 个参数值`)
  console.log('extractConfigValues: 最终结果:', result)
  return result
} 