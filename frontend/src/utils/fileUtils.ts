/**
 * 计算文件的MD5哈希值
 * @param file 要计算哈希值的文件
 * @returns Promise<string> 文件的MD5哈希值
 */
export const calculateFileHash = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    // 检查浏览器是否支持crypto API
    if (!('crypto' in window) || !('subtle' in window.crypto)) {
      reject(new Error('浏览器不支持哈希计算'))
      return
    }

    const reader = new FileReader()
    reader.onload = async e => {
      try {
        const arrayBuffer = e.target?.result as ArrayBuffer

        try {
          // 尝试使用MD5算法
          const hashBuffer = await crypto.subtle.digest('MD5', arrayBuffer)
          const hashArray = Array.from(new Uint8Array(hashBuffer))
          const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
          resolve(hashHex)
        } catch {
          // 如果MD5不可用，使用SHA-256作为备选
          const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer)
          const hashArray = Array.from(new Uint8Array(hashBuffer))
          const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
          resolve(hashHex)
        }
      } catch {
        reject(new Error('哈希计算失败'))
      }
    }
    reader.onerror = () => {
      reject(new Error('文件读取失败'))
    }
    reader.readAsArrayBuffer(file)
  })
}
