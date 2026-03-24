import type { ComputedRef, Ref } from "vue";
import { nextTick } from "vue";

/**
 * 自动填充页面 - 无限滚动辅助工具
 *
 * 当内容不足以填满容器时自动加载更多数据，解决大屏幕下无法触发滚动加载的问题
 *
 * @param scrollContainer 滚动容器的 ref
 * @param hasMore 是否还有更多数据的计算属性
 * @param isLoading 加载状态的 ref 或计算属性
 * @param loadMoreFn 加载更多的函数（可以是同步或异步）
 * @returns checkAndAutoLoadMore - 在数据加载完成后调用的检查函数
 */
export function useAutoFillPage(
  scrollContainer: Ref<HTMLElement | null>,
  hasMore: ComputedRef<boolean>,
  isLoading: Ref<boolean> | ComputedRef<boolean>,
  loadMoreFn: () => void | Promise<void>,
) {
  /**
   * 检查是否需要自动加载更多
   * 应在每次数据加载完成后调用
   */
  const checkAndAutoLoadMore = async () => {
    await nextTick();

    // 如果正在加载或没有更多数据，直接返回
    if (isLoading.value || !hasMore.value) {
      return;
    }

    // 检查内容高度是否小于容器高度
    const container = scrollContainer.value;
    if (container) {
      // 使用小的阈值来判断是否填满，避免浮点数精度问题
      const isNotFilled = container.scrollHeight <= container.clientHeight + 10;

      if (isNotFilled) {
        await loadMoreFn();
      }
    }
  };

  return {
    checkAndAutoLoadMore,
  };
}
