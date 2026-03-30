import { onUnmounted, ref } from "vue";

export interface LongPressOptions {
  duration?: number;
  moveThreshold?: number;
  onStart?: () => void;
  onFinish?: () => void;
  onCancel?: () => void;
}

export function useLongPress(options: LongPressOptions = {}) {
  const {
    duration = 500,
    moveThreshold = 10,
    onStart,
    onFinish,
    onCancel,
  } = options;

  const isPressed = ref(false);
  let timeoutId: ReturnType<typeof setTimeout> | null = null;
  let startX = 0;
  let startY = 0;

  const clearTimer = () => {
    if (timeoutId) {
      clearTimeout(timeoutId);
      timeoutId = null;
    }
  };

  const onPointerDown = (event: PointerEvent) => {
    if (event.button !== 0) return;

    startX = event.clientX;
    startY = event.clientY;
    isPressed.value = true;
    onStart?.();

    clearTimer();
    timeoutId = setTimeout(() => {
      isPressed.value = false;
      onFinish?.();
    }, duration);
  };

  const onPointerMove = (event: PointerEvent) => {
    if (!isPressed.value) return;

    const deltaX = Math.abs(event.clientX - startX);
    const deltaY = Math.abs(event.clientY - startY);

    if (deltaX > moveThreshold || deltaY > moveThreshold) {
      isPressed.value = false;
      clearTimer();
      onCancel?.();
    }
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
    onPointerMove,
    onPointerUp,
    onPointerLeave,
    onPointerCancel,
  };
}
