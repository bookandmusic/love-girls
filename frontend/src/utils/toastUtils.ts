import { inject } from 'vue'

/**
 * 使用注入的 showToast 函数
 *
 * 这个钩子函数用于在子组件中显示 toast 通知
 * 它会从父级布局组件（AdminLayout 或 DefaultLayout）中注入 showToast 函数
 *
 * 使用方法：
 * ```typescript
 * import { useToast } from '@/utils/toastUtils'
 *
 * const showToast = useToast()
 * showToast('操作成功！', 'success')
 * showToast('发生错误！', 'error')
 * showToast('提示信息', 'info')
 * ```
 *
 * @returns showToast function that can be called to show toast notifications
 */
export const useToast = () => {
  const showToast =
    inject<(message: string, type?: 'success' | 'error' | 'info' | undefined) => void>('showToast')

  if (!showToast) {
    // 如果没有注入的showToast函数，则抛出错误
    // 这意味着组件不在提供 showToast 函数的组件树下
    throw new Error(
      'useToast must be used within a component树 that provides showToast, such as under AdminLayout or DefaultLayout'
    )
  }

  return showToast
}
