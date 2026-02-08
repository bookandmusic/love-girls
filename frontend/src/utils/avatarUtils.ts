/**
 * 根据用户名生成个性化头像
 * @param name 用户名
 * @param size 头像尺寸
 * @returns base64格式的PNG图像数据URL
 */
export const generateAvatar = (name: string, size = 144): string => {
  const canvas = document.createElement('canvas')
  canvas.width = size
  canvas.height = size
  const ctx = canvas.getContext('2d')

  if (!ctx) return ''

  // 创建一个简单的hash函数
  const hashCode = (str: string): number => {
    let hash = 0
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash)
    }
    return hash
  }

  // 从hash生成颜色
  const intToRGB = (i: number): string => {
    const c = (i & 0x00ffffff).toString(16).toUpperCase()

    return '#' + '00000'.substring(0, 6 - c.length) + c
  }

  // 生成基于名字的hash值
  const hash = hashCode(name)

  // 生成背景色和图案色
  const backgroundColor = intToRGB(hash)
  const patternColor = intToRGB(~hash)

  // 绘制背景
  ctx.fillStyle = backgroundColor
  ctx.fillRect(0, 0, size, size)

  // 绘制对称的几何图案 (5x5网格)
  const grid = 5
  const cellSize = size / grid
  const hashStr = Math.abs(hash).toString(16).padStart(8, '0')

  // 使用hash值决定哪些格子需要填充
  for (let i = 0; i < Math.ceil((grid * grid) / 2); i++) {
    const charIndex = i % hashStr.length
    const char = hashStr[charIndex] || '0'
    const shouldFill = parseInt(char, 16) % 2 === 0

    if (shouldFill) {
      const row = Math.floor(i / grid)
      const col = i % grid

      // 左半部分
      ctx.fillStyle = patternColor
      ctx.fillRect(col * cellSize, row * cellSize, cellSize, cellSize)

      // 右半部分镜像
      ctx.fillRect((grid - 1 - col) * cellSize, row * cellSize, cellSize, cellSize)
    }
  }

  return canvas.toDataURL('image/png')
}

/**
 * 处理头像加载错误的函数
 * @param event 错误事件
 * @param name 用户名
 * @param size 头像尺寸
 */
export const handleAvatarError = (event: Event, name: string, size: number): void => {
  const target = event.target as HTMLImageElement
  target.src = generateAvatar(name, size)
  target.onerror = null // 防止无限循环
}
