import { onUnmounted, ref } from "vue";

export interface LongPressOptions {
  duration?: number;
  onStart?: () => void;
  onFinish?: () => void;
  onCancel?: () => void;
}

export function useLongPress(options: LongPressOptions = {}) {
  const { duration = 500, onStart, onFinish, onCancel } = options;

  const isPressed = ref(false);
  let timeoutId: ReturnType<typeof setTimeout> | null = null;

  const clearTimer = () => {
    if (timeoutId) {
      clearTimeout(timeoutId);
      timeoutId = null;
    }
  };

  const onPointerDown = (event: PointerEvent) => {
    if (event.button !== 0) return;

    isPressed.value = true;
    onStart?.();

    clearTimer();
    timeoutId = setTimeout(() => {
      isPressed.value = false;
      onFinish?.();
    }, duration);
  };

  const onPointerUp = () => {
    if (!isPressed.value) return;

    isPressed.value = false;
    clearTimer();
    onCancel?.();
  };

  const onPointerLeave = () => {
    if (!isPressed.value) return;

    isPressed.value = false;
    clearTimer();
    onCancel?.();
  };

  const onPointerCancel = () => {
    isPressed.value = false;
    clearTimer();
    onCancel?.();
  };

  onUnmounted(() => {
    clearTimer();
  });

  return {
    isPressed,
    onPointerDown,
    onPointerUp,
    onPointerLeave,
    onPointerCancel,
  };
}
